package tui

import (
	"strings"
	"testing"
)

func TestApplyThemeKnownAndUnknown(t *testing.T) {
	t.Cleanup(func() { applyTheme(defaultThemeName) })

	if applyTheme("no-such-theme") {
		t.Fatal("applyTheme accepted an unknown name")
	}
	if currentThemeName != defaultThemeName {
		t.Fatalf("unknown theme changed currentThemeName to %q", currentThemeName)
	}
	if !applyTheme("nord") {
		t.Fatal("applyTheme rejected nord")
	}
	if currentThemeName != "nord" {
		t.Fatalf("currentThemeName = %q, want nord", currentThemeName)
	}
	if chatInputBGSeq == "" || !strings.HasPrefix(chatInputBGSeq, "\x1b[") {
		t.Fatalf("chatInputBGSeq not a SGR sequence: %q", chatInputBGSeq)
	}
}

func TestThemeCommandPersistsAndReloads(t *testing.T) {
	t.Cleanup(func() { applyTheme(defaultThemeName) })
	dir := t.TempDir()

	// A bare :theme lists the choices without switching.
	msg := themeCommand(dir, "")
	if !strings.Contains(msg, "nord") || !strings.Contains(msg, defaultThemeName) {
		t.Fatalf("listing misses themes: %q", msg)
	}
	if currentThemeName != defaultThemeName {
		t.Fatalf("bare :theme switched theme to %q", currentThemeName)
	}

	if msg := themeCommand(dir, "Gruvbox"); msg != "theme gruvbox" { // case-insensitive
		t.Fatalf("switch message = %q", msg)
	}

	// A fresh session picks the persisted choice back up.
	applyTheme(defaultThemeName)
	loadTheme(dir)
	if currentThemeName != "gruvbox" {
		t.Fatalf("loadTheme restored %q, want gruvbox", currentThemeName)
	}

	if msg := themeCommand(dir, "sparkle-pony"); !strings.Contains(msg, "unknown theme") {
		t.Fatalf("unknown theme message = %q", msg)
	}
}

func TestThemeArgCandidates(t *testing.T) {
	if c := themeArgCandidates("theme "); len(c) != len(themes) || c[0] != "theme "+themeNames()[0] {
		t.Fatalf("candidates = %v", c)
	}
	if c := themeArgCandidates("topic go"); c != nil {
		t.Fatalf("non-theme input produced candidates: %v", c)
	}
}
