package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/curriculum"
)

func rkey(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

// A charwise Visual selection driven by the keyboard yanks the covered text,
// reusing the mouse drag-select machinery.
func TestReaderVisualYank(t *testing.T) {
	var copied string
	prev := copyToClipboard
	copyToClipboard = func(s string) error { copied = s; return nil }
	t.Cleanup(func() { copyToClipboard = prev })

	c := newLessonPane()
	c.focused = true
	c.setSize(40, 8)
	c.setLesson("alpha beta gamma")
	// contentLines: 0 = " lesson " badge, 1 = the body line. Put the cursor on
	// the body's first character.
	c.curLine, c.curCol = 1, 0

	c, _ = c.Update(rkey("v"))
	if !c.docVisual {
		t.Fatal("v should enter document Visual mode")
	}
	// Extend across "alpha" (4 more cells) and yank.
	for i := 0; i < 4; i++ {
		c, _ = c.Update(rkey("l"))
	}
	c, _ = c.Update(rkey("y"))
	if c.docVisual {
		t.Fatal("y should leave Visual mode")
	}
	if copied != "alpha" {
		t.Fatalf("yanked %q, want %q", copied, "alpha")
	}
}

// gd under a [[wikilink]] reports its target (alias stripped); off a link it
// reports nothing.
func TestWikilinkAtCursor(t *testing.T) {
	c := newLessonPane()
	c.setSize(60, 8)
	c.setLesson("see [[Ownership & Moves|ownership]] for more")
	body := 1 // contentLines[1] is the body

	// Locate the "[[" so we can land the cursor inside the link.
	line := c.lineText(body)
	open := strings.Index(line, "[[")
	if open < 0 {
		t.Fatalf("no wikilink rendered in %q", line)
	}
	c.curLine, c.curCol = body, open+3 // inside the target text

	got, ok := c.wikilinkAtCursor()
	if !ok || got != "Ownership & Moves" {
		t.Fatalf("wikilinkAtCursor = %q ok=%v, want %q", got, ok, "Ownership & Moves")
	}

	// Off the link (column 0 = "s" of "see") there is nothing to follow.
	c.curCol = 0
	if _, ok := c.wikilinkAtCursor(); ok {
		t.Fatal("cursor off a link should report no target")
	}
}

// "/" search moves the cursor to the next match and n / N cycle matches.
func TestReaderSearch(t *testing.T) {
	c := newLessonPane()
	c.focused = true
	c.setSize(40, 6)
	c.setLesson("apple\nbanana\napricot\ncherry")

	if !c.search("ap", true) {
		t.Fatal("search should find the first 'ap'")
	}
	firstLine := c.curLine
	if got := c.lineText(c.curLine); !strings.HasPrefix(strings.TrimSpace(got), "apple") {
		t.Fatalf("first match landed on %q", got)
	}
	// n advances to the next match ("apricot" line).
	if !c.search("", true) {
		t.Fatal("n should find the next 'ap'")
	}
	if c.curLine <= firstLine {
		t.Fatalf("n should move forward past line %d, got %d", firstLine, c.curLine)
	}
	// N goes back to the first match.
	if !c.search("", false) {
		t.Fatal("N should find the previous match")
	}
	if c.curLine != firstLine {
		t.Fatalf("N should return to line %d, got %d", firstLine, c.curLine)
	}
}

// recordLessonJump keeps an undo-stack of visited lectures: it appends, dedupes
// consecutive positions, and never records while the jumplist is driving.
func TestRecordLessonJump(t *testing.T) {
	var m Model
	m.currentTopicID, m.lesson.curLine = "a", 5
	m.recordLessonJump()
	m.currentTopicID, m.lesson.curLine = "b", 10
	m.recordLessonJump()
	if len(m.lessonJumps) != 2 || m.lessonJumpIdx != 2 {
		t.Fatalf("after two jumps: len=%d idx=%d", len(m.lessonJumps), m.lessonJumpIdx)
	}
	if m.lessonJumps[0] != (lessonJump{"a", 5, 0}) {
		t.Fatalf("first entry = %+v", m.lessonJumps[0])
	}
	// A consecutive duplicate position doesn't grow the list.
	m.recordLessonJump()
	if len(m.lessonJumps) != 2 {
		t.Fatalf("consecutive dup should not grow, len=%d", len(m.lessonJumps))
	}
	// No recording while navigating the jumplist itself.
	m.inLessonJump = true
	m.currentTopicID, m.lesson.curLine = "c", 1
	m.recordLessonJump()
	if len(m.lessonJumps) != 2 {
		t.Fatalf("should not record during a jump, len=%d", len(m.lessonJumps))
	}
}

// The fuzzy finder lists the course lectures and Enter jumps to one, recording
// a jumplist entry so ⌃o can return.
func TestLessonFinderJump(t *testing.T) {
	m := readyModel(t)
	tm, _ := m.openLessonFinder()
	m = tm.(Model)
	if !m.finderMode {
		t.Fatal("openLessonFinder should open the overlay")
	}
	if len(m.topicResults) == 0 {
		t.Fatal("finder should list the course lectures")
	}
	target := m.topicResults[0].id
	tm, _ = m.updateLessonFinder(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(Model)
	if m.finderMode {
		t.Fatal("Enter should close the finder")
	}
	if m.currentTopicID != target {
		t.Fatalf("finder jump opened %q, want %q", m.currentTopicID, target)
	}
	if len(m.lessonJumps) == 0 {
		t.Fatal("a finder jump should record a jumplist entry")
	}
}

// Editing is offered only for file-backed lectures; the vault-course lecture is
// editable, and clearing its path falls back to a read-only notice.
func TestEnterLessonEdit(t *testing.T) {
	m := readyModel(t)
	id := m.currentTopicID
	if m.topicPathByID[id] == "" {
		t.Skip("current lecture is not file-backed in this environment")
	}
	tm, _ := m.enterLessonEdit("i")
	m = tm.(Model)
	if !m.lessonEditing {
		t.Fatal("i should enter edit mode on a file-backed lecture")
	}
	if m.lessonEditor.Value() != m.topicByID[id].Lesson {
		t.Fatal("the editor should load the lecture source verbatim")
	}
	tm, _ = m.exitLessonEdit()
	m = tm.(Model)
	if m.lessonEditing {
		t.Fatal("exitLessonEdit should return to the reader")
	}

	// a / o / O also enter edit mode on a file-backed lecture.
	for _, k := range []string{"a", "o", "O"} {
		tm, _ = m.enterLessonEdit(k)
		m = tm.(Model)
		if !m.lessonEditing {
			t.Fatalf("%q should enter edit mode", k)
		}
		tm, _ = m.exitLessonEdit()
		m = tm.(Model)
	}

	// With no source path, edit mode is refused.
	m.topicPathByID = map[string]string{}
	tm, _ = m.enterLessonEdit("i")
	m = tm.(Model)
	if m.lessonEditing {
		t.Fatal("a non-file-backed lecture must not enter edit mode")
	}
}

func TestResolveTopicLink(t *testing.T) {
	topics := map[string]curriculum.Topic{
		"py-b-vars": {ID: "py-b-vars", Title: "Variables & Types"},
		"py-b-loop": {ID: "py-b-loop", Title: "Loops"},
	}
	if id, ok := resolveTopicLink("variables & types", topics); !ok || id != "py-b-vars" {
		t.Fatalf("title match = %q ok=%v", id, ok)
	}
	if id, ok := resolveTopicLink("py-b-loop", topics); !ok || id != "py-b-loop" {
		t.Fatalf("id match = %q ok=%v", id, ok)
	}
	if _, ok := resolveTopicLink("nonexistent", topics); ok {
		t.Fatal("unknown target should not resolve")
	}
	if _, ok := resolveTopicLink("  ", topics); ok {
		t.Fatal("blank target should not resolve")
	}
}
