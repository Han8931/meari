package core

import (
	"context"
	"strings"
	"testing"

	"meari/internal/config"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// newTestService builds a Service over a temp vault and an offline tutor (no API
// key, non-Ollama provider => offline built-in content).
func newTestService(t *testing.T) *Service {
	t.Helper()
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tut := tutor.New(config.AIConfig{Provider: "openai"})
	if !tut.Offline() {
		t.Fatal("expected offline tutor for the test")
	}
	return New(v, tut)
}

func TestSaveOpenList(t *testing.T) {
	s := newTestService(t)

	meta, err := s.SaveNote("history/Cold War.md", "# Cold War\n\nA rivalry.\n")
	if err != nil {
		t.Fatalf("SaveNote: %v", err)
	}
	if meta.Title != "Cold War" {
		t.Fatalf("derived title = %q", meta.Title)
	}

	n, err := s.OpenNote("history/Cold War.md")
	if err != nil {
		t.Fatalf("OpenNote: %v", err)
	}
	if n.Source != "user" || n.Body == "" {
		t.Fatalf("opened note wrong: %+v", n)
	}

	list, err := s.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(list) != 1 || list[0].Path != "history/Cold War.md" {
		t.Fatalf("ListNotes = %+v", list)
	}
}

func TestGenerateLessonWritesNote(t *testing.T) {
	s := newTestService(t)
	meta, err := s.GenerateLesson(context.Background(), "the french revolution")
	if err != nil {
		t.Fatalf("GenerateLesson: %v", err)
	}
	if meta.Path == "" || meta.Title == "" {
		t.Fatalf("generated meta empty: %+v", meta)
	}
	// It must actually persist to the vault.
	if _, err := s.OpenNote(meta.Path); err != nil {
		t.Fatalf("generated note not readable: %v", err)
	}
}

func TestBacklinks(t *testing.T) {
	s := newTestService(t)
	if _, err := s.SaveNote("Limits.md", "# Limits\n\nfoundational.\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SaveNote("Derivatives.md", "# Derivatives\n\nbuilds on [[Limits]].\n"); err != nil {
		t.Fatal(err)
	}

	back, err := s.Backlinks("Limits.md")
	if err != nil {
		t.Fatalf("Backlinks: %v", err)
	}
	if len(back) != 1 || back[0].Title != "Derivatives" {
		t.Fatalf("Backlinks = %+v", back)
	}

	// A note nobody links to has no backlinks.
	none, err := s.Backlinks("Derivatives.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(none) != 0 {
		t.Fatalf("expected no backlinks, got %+v", none)
	}
}

func TestSearch(t *testing.T) {
	s := newTestService(t)
	_, _ = s.SaveNote("math/Algebra.md", "# Algebra\n\nsolving equations\n")
	_, _ = s.SaveNote("bio/Cells.md", "# Cells\n\nmitochondria\n")

	hits, err := s.Search("equations")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].Title != "Algebra" {
		t.Fatalf("body search = %+v", hits)
	}
	if all, _ := s.Search(""); len(all) != 2 {
		t.Fatalf("empty query should return all, got %d", len(all))
	}
}

func TestEssayAndChatOffline(t *testing.T) {
	s := newTestService(t)
	g, err := s.GradeEssay(context.Background(), "explain X", "my answer")
	if err != nil {
		t.Fatal(err)
	}
	if g.Score != 1 { // offline: any non-empty answer passes
		t.Fatalf("offline essay score = %v", g.Score)
	}
	reply, err := s.Chat(context.Background(), "", []tutor.ChatTurn{{Role: "user", Content: "hi"}})
	if err != nil {
		t.Fatal(err)
	}
	if reply == "" {
		t.Fatal("expected an offline chat reply")
	}
}

func TestTrimTurnsAndClampContext(t *testing.T) {
	long := make([]tutor.ChatTurn, 40)
	for i := range long {
		long[i] = tutor.ChatTurn{Role: "user", Content: "turn"}
	}
	if got := len(TrimTurns(long)); got != maxChatTurns {
		t.Fatalf("trimmed to %d, want %d", got, maxChatTurns)
	}
	short := long[:3]
	if got := len(TrimTurns(short)); got != 3 {
		t.Fatalf("short history must pass through, got %d", got)
	}

	big := strings.Repeat("x", maxContextChars+500)
	clamped := ClampContext(big)
	if len([]rune(clamped)) > maxContextChars+50 {
		t.Fatalf("context not clamped: %d runes", len([]rune(clamped)))
	}
	if !strings.HasSuffix(clamped, "(truncated)") {
		t.Fatal("clamped context should be marked truncated")
	}
	if ClampContext("small") != "small" {
		t.Fatal("small context must pass through")
	}
}

func TestModelAnswerOffline(t *testing.T) {
	s := newTestService(t)
	ans, err := s.ModelAnswer(context.Background(), "explain X")
	if err != nil {
		t.Fatal(err)
	}
	if ans == "" {
		t.Fatal("expected an offline model-answer fallback")
	}
}
