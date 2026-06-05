package editor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func ctrlR() tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyCtrlR} }

func TestUndoRedoDeleteLine(t *testing.T) {
	m := New("one\ntwo\nthree", true, nil)
	m.SetSize(40, 10)

	m = apply(m, key("d"), key("d"))
	if got := m.Value(); got != "two\nthree" {
		t.Fatalf("after dd: %q", got)
	}
	m = apply(m, key("u"))
	if got := m.Value(); got != "one\ntwo\nthree" {
		t.Fatalf("after u: %q", got)
	}
	m = apply(m, ctrlR())
	if got := m.Value(); got != "two\nthree" {
		t.Fatalf("after ctrl+r: %q", got)
	}
}

func TestUndoInsertSessionIsOneUnit(t *testing.T) {
	m := New("base", true, nil)
	m.SetSize(40, 10)
	// One insert session types several characters; a single u removes them all.
	m = apply(m, key("i"), key("X"), key("Y"), key("Z"), tea.KeyMsg{Type: tea.KeyEsc})
	if got := m.Value(); got != "XYZbase" {
		t.Fatalf("after insert: %q", got)
	}
	m = apply(m, key("u"))
	if got := m.Value(); got != "base" {
		t.Fatalf("one undo should revert the whole session: %q", got)
	}
}

func TestUndoAtOldestIsSafe(t *testing.T) {
	m := New("text", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("u"), key("u"))
	if got := m.Value(); got != "text" {
		t.Fatalf("undo with no history must not change the buffer: %q", got)
	}
}

func TestUndoVisualDeleteAndIndent(t *testing.T) {
	m := New("alpha beta", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("v"), key("e"), key("d"))
	if got := m.Value(); got != " beta" {
		t.Fatalf("after v e d: %q", got)
	}
	m = apply(m, key("u"))
	if got := m.Value(); got != "alpha beta" {
		t.Fatalf("undo visual delete: %q", got)
	}

	m = apply(m, key(">"), key(">"), key("u"))
	if got := m.Value(); got != "alpha beta" {
		t.Fatalf("undo >>: %q", got)
	}
}

func TestNewChangeInvalidatesRedo(t *testing.T) {
	m := New("one\ntwo", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("d"), key("d"), key("u")) // delete, undo
	m = apply(m, key("x"))                     // new change
	m = apply(m, ctrlR())                      // redo must be dead now
	if got := m.Value(); got != "ne\ntwo" {
		t.Fatalf("redo after a new change must do nothing: %q", got)
	}
}

func TestSetValueClearsHistory(t *testing.T) {
	m := New("first", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("x")) // create history on the first buffer
	m.SetValue("second")
	m = apply(m, key("u"))
	if got := m.Value(); got != "second" {
		t.Fatalf("undo must not cross buffers: %q", got)
	}
}

func TestOpenBelowAutoIndents(t *testing.T) {
	m := New("    indented", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("o"), key("x"))
	if got := m.Value(); got != "    indented\n    x" {
		t.Fatalf("o should open at the same indentation: %q", got)
	}
	// And the whole o+typing session is one undo unit.
	m = apply(m, tea.KeyMsg{Type: tea.KeyEsc}, key("u"))
	if got := m.Value(); got != "    indented" {
		t.Fatalf("undo after o: %q", got)
	}
}

func TestOpenAboveAutoIndents(t *testing.T) {
	// Note: the textarea sanitizes tabs to spaces on input, so indentation is
	// space-based by the time o/O sees it.
	m := New("  deep", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("O"), key("y"))
	if got := m.Value(); got != "  y\n  deep" {
		t.Fatalf("O should open above at the same indentation: %q", got)
	}
}

func TestOpenBelowNoIndentUnchanged(t *testing.T) {
	m := New("plain", true, nil)
	m.SetSize(40, 10)
	m = apply(m, key("o"), key("x"))
	if got := m.Value(); got != "plain\nx" {
		t.Fatalf("o on an unindented line: %q", got)
	}
}
