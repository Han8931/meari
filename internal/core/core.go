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
}

// New builds a Service over a vault and tutor.
func New(v *vault.Vault, t *tutor.Tutor) *Service {
	return &Service{vault: v, tutor: t}
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

// ListNotes returns every note's metadata, sorted by path.
func (s *Service) ListNotes() ([]NoteMeta, error) {
	notes, err := s.vault.List()
	if err != nil {
		return nil, err
	}
	out := make([]NoteMeta, 0, len(notes))
	for _, n := range notes {
		out = append(out, metaOf(n))
	}
	return out, nil
}

// OpenNote loads a single note by its vault-relative path.
func (s *Service) OpenNote(path string) (Note, error) {
	n, err := s.vault.Read(path)
	if err != nil {
		return Note{}, err
	}
	return Note{NoteMeta: metaOf(n), Body: n.Body, Source: n.Source}, nil
}

// SaveNote writes body to the note at path, preserving existing frontmatter or
// creating a fresh note (deriving a sensible title) when none exists.
func (s *Service) SaveNote(path, body string) (NoteMeta, error) {
	n, err := s.vault.Read(path)
	if err != nil {
		n = vault.Note{
			RelPath: path,
			Title:   vault.DeriveTitle(body, path),
			Source:  "user",
		}
	}
	n.Body = body
	saved, err := s.vault.Write(n)
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
	target, err := s.vault.Read(path)
	if err != nil {
		return nil, err
	}
	notes, err := s.vault.List()
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
	notes, err := s.vault.List()
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

// Chat continues a free-form tutoring conversation.
func (s *Service) Chat(ctx context.Context, history []tutor.ChatTurn) (string, error) {
	return s.tutor.Chat(ctx, history)
}

// GradeEssay grades a free-text answer to a study prompt.
func (s *Service) GradeEssay(ctx context.Context, prompt, answer string) (EssayResult, error) {
	g, err := s.tutor.GradeEssay(ctx, prompt, answer)
	if err != nil {
		return EssayResult{}, err
	}
	return EssayResult{Score: g.Score, Feedback: g.Feedback}, nil
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
