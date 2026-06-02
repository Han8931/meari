package editor

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func key(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func apply(m Model, msgs ...tea.Msg) Model {
	for _, msg := range msgs {
		tm, _ := m.Update(msg)
		m = tm.(Model)
	}
	return m
}

func TestVimModeSwitchingAndInsert(t *testing.T) {
	m := New("", true, nil)
	m.SetSize(40, 10)
	if m.mode != modeNormal {
		t.Fatalf("vim editor should open in Normal mode")
	}

	// i enters Insert; typing lands in the buffer; Esc returns to Normal.
	m = apply(m, key("i"), key("h"), key("e"), key("y"))
	if m.mode != modeInsert {
		t.Fatalf("expected Insert mode after 'i'")
	}
	if got := m.Value(); got != "hey" {
		t.Fatalf("buffer = %q, want %q", got, "hey")
	}
	m = apply(m, tea.KeyMsg{Type: tea.KeyEsc})
	if m.mode != modeNormal {
		t.Fatalf("Esc should return to Normal mode")
	}
}

func TestVimDeleteLineAndOpenBelow(t *testing.T) {
	m := New("alpha\nbeta\ngamma", true, nil)
	m.SetSize(40, 10)

	// Cursor starts at top (line 0). dd removes the first line.
	m = apply(m, key("d"), key("d"))
	if got := m.Value(); got != "beta\ngamma" {
		t.Fatalf("after dd buffer = %q, want %q", got, "beta\ngamma")
	}

	// o opens a line below the current one and enters Insert.
	m = apply(m, key("o"), key("X"))
	if m.mode != modeInsert {
		t.Fatalf("o should enter Insert mode")
	}
	if got := m.Value(); got != "beta\nX\ngamma" {
		t.Fatalf("after o X buffer = %q, want %q", got, "beta\nX\ngamma")
	}
}

func TestVimReplaceChar(t *testing.T) {
	m := New("cat", true, nil)
	m.SetSize(40, 10)
	// Cursor on 'c'; r b replaces it with 'b'.
	m = apply(m, key("r"), key("b"))
	if got := m.Value(); got != "bat" {
		t.Fatalf("after rb buffer = %q, want %q", got, "bat")
	}
	if m.mode != modeNormal {
		t.Fatalf("replace should stay in Normal mode")
	}
}

func TestStatusBadgeReflectsMode(t *testing.T) {
	m := New("", true, nil)
	m.SetSize(40, 10)
	if !strings.Contains(m.statusLine(), "NORMAL") {
		t.Errorf("status line should show NORMAL in normal mode: %q", m.statusLine())
	}
	m = apply(m, key("i"))
	if !strings.Contains(m.statusLine(), "INSERT") {
		t.Errorf("status line should show INSERT after 'i': %q", m.statusLine())
	}
}

func TestPendingOperatorClearedByEsc(t *testing.T) {
	m := New("hello", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("d")) // arm operator
	if m.pending == 0 {
		t.Fatal("'d' should arm a pending operator")
	}
	m = apply(m, tea.KeyMsg{Type: tea.KeyEsc}) // cancel
	if m.pending != 0 {
		t.Fatal("Esc should cancel the pending operator")
	}
	if got := m.Value(); got != "hello" {
		t.Fatalf("buffer should be unchanged after d<Esc>, got %q", got)
	}
}
