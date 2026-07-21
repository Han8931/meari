package main

// app.go — the Wails binding layer: a thin App over core.Service. Every
// exported method here is callable from the TypeScript front-end (Wails
// generates the typed bindings). No feature logic lives here — core is the
// single brain — only DTO shaping, markdown rendering, and stream bookkeeping.

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/tutor"
)

// App is the bound object. emit is injectable so tests can record events
// without a live Wails context.
type App struct {
	ctx context.Context
	svc *core.Service
	cfg config.Config
	md  goldmark.Markdown

	emit func(event string, data ...any)

	mu      sync.Mutex
	streams map[string]context.CancelFunc
	nextID  int
}

func newApp(svc *core.Service, cfg config.Config) *App {
	a := &App{
		svc:     svc,
		cfg:     cfg,
		md:      goldmark.New(goldmark.WithExtensions(extension.GFM)),
		streams: map[string]context.CancelFunc{},
	}
	a.emit = func(event string, data ...any) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, event, data...)
		}
	}
	return a
}

// startup is Wails' OnStartup hook.
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// --- DTOs (JSON-shaped for the front-end) ---

// ChatTurn is one message of the tutor conversation.
type ChatTurn struct {
	Role    string `json:"role"` // "user" | "assistant"
	Content string `json:"content"`
}

// AppInfo is the one-shot startup snapshot the front-end reads.
type AppInfo struct {
	Offline  bool   `json:"offline"`
	VaultDir string `json:"vaultDir"`
	Version  string `json:"version"`
}

// --- read-only queries ---

// Info reports whether the tutor is offline plus where the vault lives.
func (a *App) Info() AppInfo {
	return AppInfo{
		Offline:  a.svc.Offline(),
		VaultDir: a.cfg.VaultDir,
		Version:  version,
	}
}

// Tree returns the combined vault + courses file structure.
func (a *App) Tree() ([]core.TreeEntry, error) {
	return a.svc.Tree()
}

// OpenNote loads a note (vault or course) by path.
func (a *App) OpenNote(path string) (core.Note, error) {
	return a.svc.OpenNote(path)
}

// Preview renders markdown to HTML with wikilinks turned into clickable spans
// — identical output to the web server's /api/preview.
func (a *App) Preview(body string) (string, error) {
	var buf strings.Builder
	if err := a.md.Convert([]byte(body), &buf); err != nil {
		return "", err
	}
	return renderWikilinks(buf.String()), nil
}

// Backlinks lists notes that link to the note at path.
func (a *App) Backlinks(path string) ([]core.NoteMeta, error) {
	return a.svc.Backlinks(path)
}

// Search returns notes matching a full-text query.
func (a *App) Search(query string) ([]core.NoteMeta, error) {
	return a.svc.Search(query)
}

// Courses lists every course manifest.
func (a *App) Courses() ([]core.CourseMeta, error) {
	return a.svc.ListCourses()
}

// Course loads one course by id/title/path.
func (a *App) Course(key string) (core.Course, error) {
	return a.svc.LoadCourse(key)
}

// --- mutations ---

// SaveNote writes a note's body and returns the refreshed metadata.
func (a *App) SaveNote(path, body string) (core.NoteMeta, error) {
	return a.svc.SaveNote(path, body)
}

// NewNote creates a note at path (title.md) with a starter heading and returns
// its metadata.
func (a *App) NewNote(path, title string) (core.NoteMeta, error) {
	return a.svc.SaveNote(path, "# "+title+"\n\n")
}

// Rename moves a note or directory.
func (a *App) Rename(oldPath, newPath string) error {
	return a.svc.Rename(oldPath, newPath)
}

// Delete removes a note or directory.
func (a *App) Delete(path string) error {
	return a.svc.Delete(path)
}

// --- study ---

// GradeEssay scores a free-text answer against a prompt.
func (a *App) GradeEssay(prompt, answer string) (core.EssayResult, error) {
	return a.svc.GradeEssay(a.streamCtx(), prompt, answer)
}

// ModelAnswer returns a reference answer for a prompt.
func (a *App) ModelAnswer(prompt string) (string, error) {
	return a.svc.ModelAnswer(a.streamCtx(), prompt)
}

// --- streaming ---

// StartChat streams a tutor reply grounded on the open note. Deltas arrive as
// "stream:delta" events; completion as "stream:done", failure "stream:error".
func (a *App) StartChat(path string, history []ChatTurn) string {
	studyCtx := ""
	if path != "" {
		if n, err := a.svc.OpenNote(path); err == nil {
			studyCtx = n.Title + "\n\n" + n.Body
		}
	}
	turns := make([]tutor.ChatTurn, 0, len(history))
	for _, t := range history {
		turns = append(turns, tutor.ChatTurn{Role: t.Role, Content: t.Content})
	}
	return a.stream(func(ctx context.Context, onDelta func(string)) (string, error) {
		return a.svc.ChatStream(ctx, studyCtx, turns, onDelta)
	})
}

// StartExplain streams a plain-words explanation of the selected text — the
// GUI counterpart of the TUI's :explain.
func (a *App) StartExplain(selection string) string {
	const prompt = "Explain the selected text in simple words, as if to someone new to the " +
		"subject. Keep it short, avoid jargon, and give one concrete example if it helps."
	turns := []tutor.ChatTurn{{
		Role:    "user",
		Content: prompt + "\n\nSelected text:\n" + selection,
	}}
	return a.stream(func(ctx context.Context, onDelta func(string)) (string, error) {
		return a.svc.ChatStream(ctx, "", turns, onDelta)
	})
}

// StartPolish streams an AI edit of a note body (or a selection), following an
// instruction; the front-end diffs the result before applying.
func (a *App) StartPolish(body, instruction string) string {
	return a.stream(func(ctx context.Context, onDelta func(string)) (string, error) {
		return a.svc.PolishNote(ctx, body, instruction, onDelta)
	})
}

// CancelStream stops an in-flight stream; partial text already delivered stays.
func (a *App) CancelStream(id string) {
	a.mu.Lock()
	cancel := a.streams[id]
	a.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// stream runs an AI call in the background, emitting deltas as they arrive and
// a terminal done/error event. It returns the stream id immediately so the
// front-end can subscribe and cancel. Mirrors page_sage's proven helper.
func (a *App) stream(run func(ctx context.Context, onDelta func(string)) (string, error)) string {
	a.mu.Lock()
	a.nextID++
	id := fmt.Sprintf("s%d", a.nextID)
	base := a.ctx
	if base == nil {
		base = context.Background()
	}
	ctx, cancel := context.WithCancel(base)
	a.streams[id] = cancel
	a.mu.Unlock()

	go func() {
		defer func() {
			a.mu.Lock()
			delete(a.streams, id)
			a.mu.Unlock()
			cancel()
		}()
		full, err := run(ctx, func(delta string) { a.emit("stream:delta", id, delta) })
		if err != nil && ctx.Err() == nil {
			a.emit("stream:error", id, err.Error())
			return
		}
		a.emit("stream:done", id, full)
	}()
	return id
}

// streamCtx returns the app context (or Background before startup), for the
// synchronous AI calls (grade, model answer).
func (a *App) streamCtx() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

// --- markdown helpers ---

var wikilinkHTML = regexp.MustCompile(`\[\[([^\]]+)\]\]`)

// renderWikilinks turns [[Target]] / [[Target|alias]] left in the rendered
// HTML into clickable spans the front-end intercepts.
func renderWikilinks(htmlStr string) string {
	return wikilinkHTML.ReplaceAllStringFunc(htmlStr, func(m string) string {
		inner := strings.TrimSuffix(strings.TrimPrefix(m, "[["), "]]")
		target, alias := inner, inner
		if i := strings.IndexByte(inner, '|'); i >= 0 {
			target = strings.TrimSpace(inner[:i])
			alias = strings.TrimSpace(inner[i+1:])
		}
		return fmt.Sprintf(`<a href="#" class="wikilink" data-target="%s">%s</a>`,
			htmlAttr(target), htmlText(alias))
	})
}

func htmlAttr(s string) string {
	r := strings.NewReplacer(`&`, "&amp;", `"`, "&quot;", `'`, "&#39;", `<`, "&lt;", `>`, "&gt;")
	return r.Replace(s)
}

func htmlText(s string) string {
	r := strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")
	return r.Replace(s)
}
