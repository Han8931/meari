package main

// app_test.go — the bindings against a temp vault, no Wails runtime needed
// (emit is a no-op without a context). Mirrors internal/tui's offline fixtures.

import (
	"strings"
	"testing"
	"time"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/tutor"
	"meari/internal/vault"
)

func newTestApp(t *testing.T) *App {
	t.Helper()
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tut := tutor.New(config.AIConfig{Provider: "openai"}) // offline (no key)
	if !tut.Offline() {
		t.Fatal("expected an offline tutor")
	}
	svc := core.New(v, tut)
	svc.SetCourseDir(t.TempDir())
	return newApp(svc, config.Config{VaultDir: v.Root()})
}

func TestInfoReportsOffline(t *testing.T) {
	a := newTestApp(t)
	info := a.Info()
	if !info.Offline {
		t.Fatal("Info should report the offline tutor")
	}
	if info.VaultDir == "" {
		t.Fatal("Info should report the vault dir")
	}
}

func TestNewSaveOpenTreeRoundtrip(t *testing.T) {
	a := newTestApp(t)

	if _, err := a.NewNote("Derivatives.md", "Derivatives"); err != nil {
		t.Fatalf("NewNote: %v", err)
	}
	if _, err := a.SaveNote("Derivatives.md", "# Derivatives\n\nA derivative measures change.\n"); err != nil {
		t.Fatalf("SaveNote: %v", err)
	}

	n, err := a.OpenNote("Derivatives.md")
	if err != nil {
		t.Fatalf("OpenNote: %v", err)
	}
	if !strings.Contains(n.Body, "measures change") {
		t.Fatalf("note body not saved: %q", n.Body)
	}

	tree, err := a.Tree()
	if err != nil {
		t.Fatalf("Tree: %v", err)
	}
	found := false
	for _, e := range tree {
		if e.Path == "Derivatives.md" {
			found = true
		}
	}
	if !found {
		t.Fatal("saved note missing from the tree")
	}
}

func TestPreviewRendersMarkdownAndWikilinks(t *testing.T) {
	a := newTestApp(t)
	html, err := a.Preview("# Title\n\nSee [[Limits]] and [[Chain Rule|the chain rule]].\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<h1") {
		t.Fatalf("markdown not rendered: %q", html)
	}
	if !strings.Contains(html, `data-target="Limits"`) {
		t.Fatalf("wikilink not rendered: %q", html)
	}
	if !strings.Contains(html, "the chain rule") {
		t.Fatalf("wikilink alias not rendered: %q", html)
	}
}

// The offline tutor still streams (a canned reply), so StartChat must emit a
// terminal done event with an id.
func TestStartChatEmitsStreamEvents(t *testing.T) {
	a := newTestApp(t)

	done := make(chan string, 1)
	a.emit = func(event string, data ...any) {
		if event == "stream:done" && len(data) >= 1 {
			if id, ok := data[0].(string); ok {
				select {
				case done <- id:
				default:
				}
			}
		}
	}

	id := a.StartChat("", []ChatTurn{{Role: "user", Content: "hi"}})
	if id == "" {
		t.Fatal("StartChat should return a stream id")
	}
	select {
	case gotID := <-done:
		if gotID != id {
			t.Fatalf("done event id = %q, want %q", gotID, id)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no stream:done event")
	}
}
