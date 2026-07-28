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

// Esc leaves the cursor one past the last rune; v there must still grab the
// character under it (Vim's inclusive cursor) rather than an empty span.
func TestChatInputVisualAtEndOfLine(t *testing.T) {
	forceColorTUI(t)
	c := visualChat(t, "hello world")
	c.input.SetCursor(len("hello world"))
	pressChat(&c, "v")
	_, start, cut := c.visualSpanInput()
	if start != 10 || cut != 11 {
		t.Fatalf("selection [%d,%d), want [10,11) — the final rune", start, cut)
	}
	if !strings.Contains(c.inputView(), "48;5;") {
		t.Fatal("no selection background painted at end of line")
	}
}

// V selects whole lines, j extends by a line, and d removes them outright.
func TestChatInputVisualLinewise(t *testing.T) {
	forceColorTUI(t)
	c := visualChat(t, "one\ntwo\nthree")
	pressChat(&c, "V")
	if !c.visual || !c.vLine {
		t.Fatal("V should enter linewise Visual")
	}
	_, start, cut := c.visualSpanInput()
	if start != 0 || cut != 4 { // "one\n"
		t.Fatalf("V selection [%d,%d), want [0,4)", start, cut)
	}
	if !strings.Contains(c.inputView(), "48;5;") {
		t.Fatal("no selection background painted in linewise Visual")
	}
	pressChat(&c, "j", "d")
	if got := c.input.Value(); got != "three" {
		t.Fatalf("after Vjd: %q, want %q", got, "three")
	}

	// Deleting the last line takes the newline above it, leaving no blank.
	c = visualChat(t, "one\ntwo")
	pressChat(&c, "G", "V", "d")
	if got := c.input.Value(); got != "one" {
		t.Fatalf("after GVd: %q, want %q", got, "one")
	}
}

// Vc empties the line but keeps it, landing in Insert mode (as Vim does).
func TestChatInputVisualLinewiseChange(t *testing.T) {
	c := visualChat(t, "one\ntwo")
	pressChat(&c, "V", "c")
	if got := c.input.Value(); got != "\ntwo" {
		t.Fatalf("after Vc: %q, want %q", got, "\ntwo")
	}
	if c.normal {
		t.Fatal("Vc should end in Insert mode")
	}
}

// v and V switch flavors while keeping the anchor; the same key twice exits.
func TestChatInputVisualToggle(t *testing.T) {
	c := visualChat(t, "one\ntwo")
	pressChat(&c, "v", "V")
	if !c.visual || !c.vLine {
		t.Fatal("V from character-wise Visual should switch to linewise")
	}
	pressChat(&c, "v")
	if !c.visual || c.vLine {
		t.Fatal("v from linewise Visual should switch back to character-wise")
	}
	pressChat(&c, "v")
	if c.visual {
		t.Fatal("v in character-wise Visual should exit")
	}
	pressChat(&c, "V", "V")
	if c.visual || c.vLine {
		t.Fatal("V in linewise Visual should exit")
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
