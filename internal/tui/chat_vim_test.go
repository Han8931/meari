package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/editor"
)

func visualChat(t *testing.T, value string) chatModel {
	t.Helper()
	c := newChat()
	c.setSize(40, 12)
	c.focus()
	c.input.SetValue(value)
	c.inputMoveTo(0, 0)
	c.enterNormal()
	return c
}

func pressChat(c *chatModel, keys ...string) {
	for _, k := range keys {
		var msg tea.KeyMsg
		switch k {
		case "esc":
			msg = tea.KeyMsg{Type: tea.KeyEsc}
		default:
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
		}
		*c, _ = c.Update(msg)
	}
}

// v + motions select; y yanks to the clipboard and exits Visual.
func TestChatInputVisualYank(t *testing.T) {
	var copied string
	prev := copyToClipboard
	copyToClipboard = func(s string) error { copied = s; return nil }
	defer func() { copyToClipboard = prev }()

	c := visualChat(t, "hello world")
	pressChat(&c, "v", "l", "l", "l", "l", "y") // select "hello"
	if copied != "hello" {
		t.Fatalf("yanked %q, want %q", copied, "hello")
	}
	if c.visual {
		t.Fatal("y should leave Visual mode")
	}
	if !c.normal {
		t.Fatal("y should return to Normal mode")
	}
}

// d deletes the selection; c deletes and re-enters Insert.
func TestChatInputVisualDeleteAndChange(t *testing.T) {
	c := visualChat(t, "hello world")
	pressChat(&c, "v", "l", "l", "l", "l", "l", "d") // delete "hello "
	if got := c.input.Value(); got != "world" {
		t.Fatalf("after d: %q, want %q", got, "world")
	}

	c = visualChat(t, "hello world")
	pressChat(&c, "$", "v", "b", "c") // change the last word
	if got := c.input.Value(); !strings.HasPrefix(got, "hello ") || strings.Contains(got, "world") {
		t.Fatalf("after c: %q", got)
	}
	if c.normal {
		t.Fatal("c should end in Insert mode")
	}
}

// gg and G jump between the input's first and last line, in Normal and Visual.
func TestChatInputGGAndG(t *testing.T) {
	c := visualChat(t, "one\ntwo\nthree")
	pressChat(&c, "G")
	if c.input.Line() != 2 {
		t.Fatalf("G: line %d, want 2", c.input.Line())
	}
	pressChat(&c, "g", "g")
	if c.input.Line() != 0 {
		t.Fatalf("gg: line %d, want 0", c.input.Line())
	}
	// Visual: G extends the selection to the last line.
	pressChat(&c, "v", "G")
	if !c.visual {
		t.Fatal("G should keep Visual mode")
	}
	_, start, cut := c.visualSpanInput()
	if start != 0 || cut < len("one\ntwo\n") {
		t.Fatalf("selection [%d,%d) should span to the last line", start, cut)
	}
	pressChat(&c, "d")
	if got := c.input.Value(); got != "three" && got != "hree" && got != "" {
		t.Fatalf("after V-G-d: %q", got)
	}
}

// The Visual selection is painted into the rendered input.
func TestChatInputVisualPaint(t *testing.T) {
	forceColorTUI(t)
	c := visualChat(t, "hello world")
	pressChat(&c, "v", "l", "l") // select "hel"
	view := c.inputView()
	if !strings.Contains(view, "48;5;") { // a background span from chatSelStyle
		t.Fatalf("no selection background painted:\n%q", view)
	}
}

// :explain uses the editor selection (or chat drag) and asks in simple words.
func TestVaultExplainSelection(t *testing.T) {
	m := newTestVaultModel(t)
	m = openNote(t, m, "N.md", "# n\n\nMonetary policy is the control of money supply.\n")

	// No selection: a hint, no request.
	tm, _ := m.runEx("explain")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "select text first") {
		t.Fatalf("notice = %q", m.notice)
	}

	m.cmdSel = &editor.Selection{Text: "Monetary policy"}
	tm, cmd := m.runEx("explain")
	m = tm.(VaultModel)
	if m.focusExcerpt != "Monetary policy" {
		t.Fatalf("excerpt = %q", m.focusExcerpt)
	}
	if !m.streaming || cmd == nil {
		t.Fatal(":explain should start a tutor reply")
	}
	last := m.chatHist[len(m.chatHist)-1]
	if last.Role != "user" || !strings.Contains(last.Content, "simple words") {
		t.Fatalf("last turn = %+v", last)
	}
}
