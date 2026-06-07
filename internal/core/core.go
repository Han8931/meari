// Package core is Meari's headless engine: the shared brain that both the
// terminal UI and the web GUI drive. It owns the vault-and-tutor orchestration
// (list/open/save notes, generate a lesson as a note, compute backlinks, search,
// chat, grade an essay) and returns plain data — no Bubble Tea, no net/http, no
// presentation. Front-ends are thin layers that call these methods and render
// the result.
//
// The rule that keeps two front-ends in parity: business logic lives here, never
// in a UI handler. As new capabilities land (a SQLite index, more study modes),
// they are wired in behind these methods so both front-ends gain them at once.
package core

import (
	"context"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"meari/internal/tutor"
	"meari/internal/vault"
)

// Service is the engine. Construct it with New and share one instance across a
// process; its collaborators are safe for the single-trusted-user model.
type Service struct {
	vault *vault.Vault
	tutor *tutor.Tutor
	// coursesDir backs the meari-course/ mount (see mount.go): generated
	// course material lives here, outside the notes vault.
	coursesDir string
}

// New builds a Service over a vault and tutor. Courses default to living
// inside the vault (the standalone layout); the app re-points them at the
// project directory with SetCourseDir.
func New(v *vault.Vault, t *tutor.Tutor) *Service {
	return &Service{vault: v, tutor: t, coursesDir: filepath.Join(v.Root(), CourseDir)}
}

// Offline reports whether the tutor is running on built-in content (no provider).
func (s *Service) Offline() bool { return s.tutor.Offline() }

// NoteMeta is the lightweight view of a note used in lists, trees, and links.
type NoteMeta struct {
	Path    string   `json:"path"`
	Title   string   `json:"title"`
	Subject string   `json:"subject"`
	Tags    []string `json:"tags"`
}

// Note is a full note: its metadata plus the markdown body and provenance.
type Note struct {
	NoteMeta
	Body   string `json:"body"`
	Source string `json:"source"`
}

// EssayResult is the outcome of grading a free-text answer.
type EssayResult struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// --- notes ---

// ListNotes returns every note's metadata — vault and course material — sorted
// by path.
func (s *Service) ListNotes() ([]NoteMeta, error) {
	notes, err := s.allNotes()
	if err != nil {
		return nil, err
	}
	out := make([]NoteMeta, 0, len(notes))
	for _, n := range notes {
		out = append(out, metaOf(n))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

// TreeEntry is one node of the vault's on-disk structure: a directory or a
// markdown note. Paths are vault-relative and "/"-separated; Name is the
// display name (the base name, without ".md" for notes), so file-tree UIs
// mirror the learner's real layout.
type TreeEntry struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Dir  bool   `json:"dir"`
}

// Tree returns the combined file structure — the vault's directories and
// notes plus the meari-course/ mount — sorted by path.
func (s *Service) Tree() ([]TreeEntry, error) {
	notes, err := s.vault.List()
	if err != nil {
		return nil, err
	}
	dirs, err := s.vault.Dirs()
	if err != nil {
		return nil, err
	}
	out := make([]TreeEntry, 0, len(dirs)+len(notes))
	for _, d := range dirs {
		if d == CourseDir || strings.HasPrefix(d, CourseDir+"/") {
			continue // the mount serves these
		}
		out = append(out, TreeEntry{Path: d, Name: path.Base(d), Dir: true})
	}
	for _, n := range notes {
		if _, ok := inMount(n.RelPath); ok {
			continue
		}
		name := strings.TrimSuffix(path.Base(n.RelPath), path.Ext(n.RelPath))
		out = append(out, TreeEntry{Path: n.RelPath, Name: name})
	}
	out = append(out, s.mountTree()...)
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

// Delete removes a note, or a directory and everything in it.
func (s *Service) Delete(path string) error { return s.deletePath(path) }

// Rename moves a note or directory to a new vault-relative path.
func (s *Service) Rename(oldPath, newPath string) error { return s.renamePath(oldPath, newPath) }

// MakeDir creates a directory, so structure can be laid out before notes exist.
func (s *Service) MakeDir(path string) error { return s.makeDirPath(path) }

// OpenNote loads a single note by its vault-relative path.
func (s *Service) OpenNote(path string) (Note, error) {
	n, err := s.readNote(path)
	if err != nil {
		return Note{}, err
	}
	return Note{NoteMeta: metaOf(n), Body: n.Body, Source: n.Source}, nil
}

// SaveNote writes body to the note at path, preserving existing frontmatter or
// creating a fresh note (deriving a sensible title) when none exists.
func (s *Service) SaveNote(path, body string) (NoteMeta, error) {
	n, err := s.readNote(path)
	if err != nil {
		n = vault.Note{
			RelPath: path,
			Title:   vault.DeriveTitle(body, path),
			Source:  "user",
		}
	}
	n.Body = body
	saved, err := s.writeNote(n)
	if err != nil {
		return NoteMeta{}, err
	}
	return metaOf(saved), nil
}

// GenerateLesson turns a learn-request into a new AI-authored note saved in the
// vault, and returns its metadata. This is the headline "I want to learn X →
// owned, linkable note" flow.
func (s *Service) GenerateLesson(ctx context.Context, request string) (NoteMeta, error) {
	nc, err := s.tutor.GenerateNote(ctx, request)
	if err != nil {
		return NoteMeta{}, err
	}
	saved, err := s.vault.Write(vault.Note{
		Title:   nc.Title,
		Subject: nc.Subject,
		Tags:    nc.Tags,
		Source:  "ai-generated",
		Body:    nc.Body,
	})
	if err != nil {
		return NoteMeta{}, err
	}
	return metaOf(saved), nil
}

// Backlinks returns the notes whose body links (via [[wikilink]]) to the note at
// path. Currently an in-memory scan of the vault; a SQLite index will back this
// later without changing the signature.
func (s *Service) Backlinks(path string) ([]NoteMeta, error) {
	target, err := s.readNote(path)
	if err != nil {
		return nil, err
	}
	notes, err := s.allNotes()
	if err != nil {
		return nil, err
	}
	var out []NoteMeta
	for _, n := range notes {
		if n.RelPath == target.RelPath {
			continue
		}
		for _, l := range vault.ParseLinks(n.Body) {
			if linkMatches(l.Target, n2target(target)) {
				out = append(out, metaOf(n))
				break
			}
		}
	}
	return out, nil
}

// Search returns notes whose title, subject, or body contains query
// (case-insensitive), sorted by path. A simple substring scan for now.
func (s *Service) Search(query string) ([]NoteMeta, error) {
	q := strings.ToLower(strings.TrimSpace(query))
	notes, err := s.allNotes()
	if err != nil {
		return nil, err
	}
	var out []NoteMeta
	for _, n := range notes {
		if q == "" ||
			strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.Subject), q) ||
			strings.Contains(strings.ToLower(n.Body), q) {
			out = append(out, metaOf(n))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

// --- tutor ---

// maxChatTurns bounds how much conversation is SENT to the model (the visible
// transcript keeps everything). Long sessions otherwise grow tokens without
// bound and overflow small local models' context windows.
const maxChatTurns = 24

// maxContextChars bounds the study-context block sent with each chat call.
const maxContextChars = 6000

// TrimTurns bounds a conversation to the most recent turns (see maxChatTurns).
func TrimTurns(history []tutor.ChatTurn) []tutor.ChatTurn {
	if len(history) > maxChatTurns {
		return history[len(history)-maxChatTurns:]
	}
	return history
}

// ClampContext bounds a study-context string to a model-friendly size.
func ClampContext(s string) string {
	r := []rune(s)
	if len(r) <= maxContextChars {
		return s
	}
	return string(r[:maxContextChars]) + "\n…(truncated)"
}

// Chat continues a free-form tutoring conversation. studyContext is what the
// learner is currently looking at (note body, challenge, code) so replies stay
// grounded; "" sends conversation only.
func (s *Service) Chat(ctx context.Context, studyContext string, history []tutor.ChatTurn) (string, error) {
	return s.ChatStream(ctx, studyContext, history, nil)
}

// ChatStream is Chat with incremental delivery: onDelta receives each reply
// chunk as the model produces it; the assembled reply is returned at the end.
func (s *Service) ChatStream(ctx context.Context, studyContext string, history []tutor.ChatTurn, onDelta func(string)) (string, error) {
	return s.tutor.ChatStream(ctx, ClampContext(studyContext), TrimTurns(history), onDelta)
}

// GradeEssay grades a free-text answer to a study prompt.
func (s *Service) GradeEssay(ctx context.Context, prompt, answer string) (EssayResult, error) {
	g, err := s.tutor.GradeEssay(ctx, prompt, answer)
	if err != nil {
		return EssayResult{}, err
	}
	return EssayResult{Score: g.Score, Feedback: g.Feedback}, nil
}

// ModelAnswer reveals a reference answer to a study prompt.
func (s *Service) ModelAnswer(ctx context.Context, prompt string) (string, error) {
	return s.tutor.ModelAnswer(ctx, prompt)
}

// --- helpers ---

func metaOf(n vault.Note) NoteMeta {
	return NoteMeta{Path: n.RelPath, Title: n.Title, Subject: n.Subject, Tags: n.Tags}
}

// linkTarget is the minimal identity a wikilink can resolve against.
type linkTarget struct {
	title, id, relPath string
}

func n2target(n vault.Note) linkTarget {
	return linkTarget{title: n.Title, id: n.ID, relPath: n.RelPath}
}

// linkMatches reports whether a wikilink target refers to the given note, by
// title (case-insensitive), id, or filename stem.
func linkMatches(target string, n linkTarget) bool {
	t := strings.TrimSpace(strings.ToLower(target))
	if t == "" {
		return false
	}
	if strings.EqualFold(t, n.title) || strings.EqualFold(t, n.id) {
		return true
	}
	stem := n.relPath
	if i := strings.LastIndexByte(stem, '/'); i >= 0 {
		stem = stem[i+1:]
	}
	stem = strings.TrimSuffix(stem, ".md")
	return strings.EqualFold(t, stem)
}
