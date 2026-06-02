package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sidebarItem is one row: a challenge and its progress, which (since each
// challenge maps to one workspace/drafts/<id>.py) doubles as the learner's file
// list.
type sidebarItem struct {
	id     string
	title  string // first line of the prompt, the topic title, or the id
	status string // "done" | "in_progress" | ""
	header bool   // a non-selectable section heading (e.g. a curriculum module)
}

// sidebarModel is the left pane: a hand-rolled cursor list. Hand-rolled rather
// than bubbles/list so it doesn't fight the parent for j/k and filter keys.
type sidebarModel struct {
	items   []sidebarItem
	cursor  int
	w, h    int
	focused bool
}

func newSidebar() sidebarModel { return sidebarModel{} }

func (s *sidebarModel) setSize(w, h int) { s.w, s.h = w, h }

// setItems replaces the list, keeping the cursor on the same id when possible
// (so a status refresh doesn't jump the selection). The cursor always lands on a
// selectable (non-header) row.
func (s *sidebarModel) setItems(items []sidebarItem) {
	var keepID string
	if s.cursor >= 0 && s.cursor < len(s.items) {
		keepID = s.items[s.cursor].id
	}
	s.items = items
	s.cursor = s.firstSelectable()
	for i, it := range items {
		if !it.header && it.id == keepID {
			s.cursor = i
			break
		}
	}
}

func (s sidebarModel) firstSelectable() int {
	for i, it := range s.items {
		if !it.header {
			return i
		}
	}
	return 0
}

// selected returns the highlighted item, or ok=false when there is no selectable
// item under the cursor.
func (s sidebarModel) selected() (sidebarItem, bool) {
	if s.cursor < 0 || s.cursor >= len(s.items) || s.items[s.cursor].header {
		return sidebarItem{}, false
	}
	return s.items[s.cursor], true
}

// Update handles cursor movement when focused, skipping header rows. Returns
// whether Enter was pressed so the parent can act on the selection.
func (s sidebarModel) Update(msg tea.Msg) (sidebarModel, bool) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return s, false
	}
	switch key.String() {
	case "j", "down":
		s.move(1)
	case "k", "up":
		s.move(-1)
	case "g", "home":
		s.cursor = s.firstSelectable()
	case "G", "end":
		s.cursor = s.lastSelectable()
	case "enter":
		return s, true
	}
	return s, false
}

// move steps the cursor by dir (±1), skipping headers and stopping at the ends.
func (s *sidebarModel) move(dir int) {
	for i := s.cursor + dir; i >= 0 && i < len(s.items); i += dir {
		if !s.items[i].header {
			s.cursor = i
			return
		}
	}
}

func (s sidebarModel) lastSelectable() int {
	for i := len(s.items) - 1; i >= 0; i-- {
		if !s.items[i].header {
			return i
		}
	}
	return 0
}

func (s sidebarModel) view() string {
	if len(s.items) == 0 {
		return hintStyle.Render("(nothing here yet)")
	}

	// Show a window of items that fits the pane height, scrolled to keep the
	// cursor visible. Returning more lines than s.h would overflow the box.
	lo, hi := s.window()

	var b strings.Builder
	for i := lo; i < hi; i++ {
		it := s.items[i]
		row := s.renderRow(it, i == s.cursor)
		b.WriteString(row)
		if i < hi-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// window returns the [lo, hi) slice of items to display so the cursor stays in
// view within s.h rows.
func (s sidebarModel) window() (int, int) {
	h := s.h
	if h <= 0 || h >= len(s.items) {
		return 0, len(s.items)
	}
	lo := s.cursor - h/2
	if lo < 0 {
		lo = 0
	}
	hi := lo + h
	if hi > len(s.items) {
		hi = len(s.items)
		lo = hi - h
	}
	return lo, hi
}

func (s sidebarModel) renderRow(it sidebarItem, selected bool) string {
	if it.header {
		return headerRow.Width(s.w).Render(it.title)
	}

	glyph := "  "
	switch it.status {
	case "done":
		glyph = doneGlyph.Render("✓ ")
	case "in_progress":
		glyph = wipGlyph.Render("… ")
	}

	title := it.title
	max := s.w - 2 // reserve room for the 2-cell glyph
	if max > 0 && lipgloss.Width(title) > max {
		title = truncate(title, max)
	}

	row := glyph + title
	if selected {
		row = selectedRow.Width(s.w).Render(row)
	}
	return row
}

// truncate cuts s to at most w display cells, adding an ellipsis.
func truncate(s string, w int) string {
	if w <= 1 {
		return "…"
	}
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	return string(r[:w-1]) + "…"
}
