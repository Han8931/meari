package editor

// vim.go implements the parts of Vim that need real text math rather than
// synthetic textarea keys: word motions with Vim semantics (w = start of the
// NEXT word, e = end of word, b = start of the previous word) and an internal
// unnamed register so deletes/yanks can be pasted with p/P. Motions treat any
// whitespace run (including line breaks) as a separator — i.e. big-WORD
// semantics for both w and W, a deliberate simplification.

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// tabIndent is inserted when Tab is pressed in Insert (or default) mode. The
// textarea has no native Tab handling — KeyTab is not a rune key — so without
// this, Tab silently does nothing.
const tabIndent = "    "

// --- cursor position & movement over the textarea ---

// cursorPos returns the cursor's logical (row, col) in rune terms. LineInfo
// reports the visual position within a soft-wrapped row; StartColumn +
// ColumnOffset recovers the logical column.
func (m *Model) cursorPos() (row, col int) {
	li := m.ta.LineInfo()
	return m.ta.Line(), li.StartColumn + li.ColumnOffset
}

// moveTo places the cursor at logical (row, col). Vertical movement goes
// through CursorUp/Down (which step visual lines), so a guard bounds the loop
// on soft-wrapped rows; SetCursor clamps the column to the line length.
func (m *Model) moveTo(row, col int) {
	for guard := 0; m.ta.Line() < row && guard < 100000; guard++ {
		before := m.ta.Line()
		m.ta.CursorDown()
		if m.ta.Line() == before && m.ta.LineInfo().RowOffset == 0 {
			break // cannot advance further
		}
	}
	for guard := 0; m.ta.Line() > row && guard < 100000; guard++ {
		before := m.ta.Line()
		m.ta.CursorUp()
		if m.ta.Line() == before && m.ta.LineInfo().RowOffset == 0 {
			break
		}
	}
	m.ta.SetCursor(col)
}

// --- Vim word motions (pure text math, unit-testable) ---

// flatIndex converts (row, col) to an index into text's rune slice, where rows
// are separated by '\n'. col is clamped to the row length.
func flatIndex(runes []rune, row, col int) int {
	lineStart, r := 0, 0
	for i := 0; i < len(runes) && r < row; i++ {
		if runes[i] == '\n' {
			r++
			lineStart = i + 1
		}
	}
	lineLen := 0
	for i := lineStart; i < len(runes) && runes[i] != '\n'; i++ {
		lineLen++
	}
	if col > lineLen {
		col = lineLen
	}
	return lineStart + col
}

// rowColOf converts a rune index back to (row, col).
func rowColOf(runes []rune, idx int) (row, col int) {
	if idx > len(runes) {
		idx = len(runes)
	}
	for i := 0; i < idx; i++ {
		if runes[i] == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}
	return row, col
}

func isWordSpace(r rune) bool { return unicode.IsSpace(r) }

// nextWordStart returns the index of the start of the next word after idx
// (Vim's w): skip the rest of the current word, then any whitespace.
func nextWordStart(runes []rune, idx int) int {
	n := len(runes)
	if n == 0 {
		return 0
	}
	i := idx
	for i < n && !isWordSpace(runes[i]) {
		i++
	}
	for i < n && isWordSpace(runes[i]) {
		i++
	}
	if i >= n {
		i = n - 1
	}
	return i
}

// nextWordEnd returns the index of the last rune of the next word end after
// idx (Vim's e): step off the current position, skip whitespace, then run to
// the final rune of that word.
func nextWordEnd(runes []rune, idx int) int {
	n := len(runes)
	if n == 0 {
		return 0
	}
	i := idx + 1
	for i < n && isWordSpace(runes[i]) {
		i++
	}
	if i >= n {
		return n - 1
	}
	for i+1 < n && !isWordSpace(runes[i+1]) {
		i++
	}
	return i
}

// prevWordStart returns the index of the start of the previous word before idx
// (Vim's b): step back, skip whitespace, then run to the word's first rune.
func prevWordStart(runes []rune, idx int) int {
	i := idx - 1
	for i > 0 && isWordSpace(runes[i]) {
		i--
	}
	for i > 0 && !isWordSpace(runes[i-1]) {
		i--
	}
	if i < 0 {
		i = 0
	}
	return i
}

// motion applies one of the pure word motions to the live cursor.
func (m *Model) motion(move func([]rune, int) int) {
	runes := []rune(m.ta.Value())
	row, col := m.cursorPos()
	target := move(runes, flatIndex(runes, row, col))
	r, c := rowColOf(runes, target)
	m.moveTo(r, c)
}

// applyMotion handles a cursor-movement key shared by Normal and Visual mode.
// It reports whether the key was a motion (and was applied).
func (m *Model) applyMotion(key string) bool {
	switch key {
	case "h", "left":
		m.arrow(tea.KeyLeft)
	case "l", "right":
		m.arrow(tea.KeyRight)
	case "j", "down":
		m.arrow(tea.KeyDown)
	case "k", "up":
		m.arrow(tea.KeyUp)
	case "w", "W":
		m.motion(nextWordStart)
	case "e", "E":
		m.motion(nextWordEnd)
	case "b", "B":
		m.motion(prevWordStart)
	case "0", "^", "_":
		m.ta.CursorStart()
	case "$":
		m.ta.CursorEnd()
	case "G":
		m.send(tea.KeyCtrlEnd) // jump to end of buffer
	default:
		return false
	}
	return true
}

// --- undo / redo ---

// editState is one undo snapshot: the whole buffer plus the cursor.
type editState struct {
	value    string
	row, col int
}

const maxUndo = 100

// pushUndo snapshots the current buffer before a mutating operation. Insert
// mode pushes once on entry, so a whole typing session undoes as one unit
// (like Vim). Starting a new change invalidates the redo stack.
func (m *Model) pushUndo() {
	row, col := m.cursorPos()
	st := editState{value: m.ta.Value(), row: row, col: col}
	if n := len(m.undoStack); n > 0 && m.undoStack[n-1].value == st.value {
		m.undoStack[n-1] = st // same text: just refresh the cursor
	} else {
		m.undoStack = append(m.undoStack, st)
		if len(m.undoStack) > maxUndo {
			m.undoStack = m.undoStack[1:]
		}
	}
	m.redoStack = nil
}

// undo restores the most recent snapshot that differs from the buffer.
func (m *Model) undo() {
	cur := m.ta.Value()
	for len(m.undoStack) > 0 {
		st := m.undoStack[len(m.undoStack)-1]
		m.undoStack = m.undoStack[:len(m.undoStack)-1]
		if st.value == cur {
			continue // no-op snapshot (e.g. an operation that changed nothing)
		}
		row, col := m.cursorPos()
		m.redoStack = append(m.redoStack, editState{value: cur, row: row, col: col})
		m.ta.SetValue(st.value)
		m.moveTo(st.row, st.col)
		return
	}
	m.status = "already at oldest change"
}

// redo re-applies the change most recently undone.
func (m *Model) redo() {
	cur := m.ta.Value()
	for len(m.redoStack) > 0 {
		st := m.redoStack[len(m.redoStack)-1]
		m.redoStack = m.redoStack[:len(m.redoStack)-1]
		if st.value == cur {
			continue
		}
		row, col := m.cursorPos()
		m.undoStack = append(m.undoStack, editState{value: cur, row: row, col: col})
		m.ta.SetValue(st.value)
		m.moveTo(st.row, st.col)
		return
	}
	m.status = "already at newest change"
}

// clearHistory drops the undo/redo stacks — called when the buffer is replaced
// wholesale (switching notes/challenges), so undo never crosses files.
func (m *Model) clearHistory() {
	m.undoStack, m.redoStack = nil, nil
}

// --- auto-indent for o/O ---

// lineIndent returns the leading whitespace of the cursor's current line, so
// o/O can open the new line at the same indentation.
func (m *Model) lineIndent() string {
	lines := strings.Split(m.ta.Value(), "\n")
	r := m.ta.Line()
	if r < 0 || r >= len(lines) {
		return ""
	}
	line := lines[r]
	i := 0
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	return line[:i]
}

// --- Visual mode ---

// enterVisual starts a selection anchored at the cursor. linewise selects whole
// lines (V); otherwise the selection is charwise (v).
func (m *Model) enterVisual(linewise bool) {
	m.anchorRow, m.anchorCol = m.cursorPos()
	m.visualLine = linewise
	m.mode = modeVisual
}

// updateVisual routes a key while a selection is active: motions extend it,
// operators act on it, Esc abandons it.
func (m Model) updateVisual(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// A pending 'g' prefix (for gg) is completed or cancelled by the next key.
	if m.pending != 0 {
		op := m.pending
		m.pending = 0
		if op == 'g' && msg.String() == "g" {
			m.send(tea.KeyCtrlHome)
		}
		return m, nil
	}

	if msg.Type == tea.KeyEsc {
		m.mode = modeNormal
		return m, nil
	}
	if m.applyMotion(msg.String()) {
		return m, nil
	}

	switch msg.String() {
	case "v":
		if m.visualLine {
			m.visualLine = false // switch V -> v, keeping the selection
		} else {
			m.mode = modeNormal
		}
	case "V":
		if m.visualLine {
			m.mode = modeNormal
		} else {
			m.visualLine = true // switch v -> V
		}
	case "o":
		// Swap the cursor and the anchor, Vim-style, to grow the other end.
		row, col := m.cursorPos()
		ar, ac := m.anchorRow, m.anchorCol
		m.anchorRow, m.anchorCol = row, col
		m.moveTo(ar, ac)
	case "d", "x":
		m.pushUndo()
		m.visualDelete()
		m.mode = modeNormal
	case "y":
		m.visualYank()
		m.mode = modeNormal
	case "c":
		m.pushUndo()
		linewise := m.visualLine
		m.visualDelete()
		if linewise {
			// cc-style change: keep an empty line to type into.
			m.ta.CursorStart()
			m.ta.InsertString("\n")
			m.arrow(tea.KeyUp)
		}
		m.mode = modeInsert
	case "<":
		m.pushUndo()
		m.visualIndent(-1)
		m.mode = modeNormal
	case ">":
		m.pushUndo()
		m.visualIndent(1)
		m.mode = modeNormal
	case "g":
		m.pending = 'g'
	}
	return m, nil
}

// visualSpan returns the buffer runes and the selection as a half-open range
// [start, cut). Charwise selections include the rune under both ends; linewise
// selections cover whole lines including one bounding newline so deleting the
// span removes the lines entirely.
func (m *Model) visualSpan() (runes []rune, start, cut int) {
	runes = []rune(m.ta.Value())
	row, col := m.cursorPos()
	a := flatIndex(runes, m.anchorRow, m.anchorCol)
	b := flatIndex(runes, row, col)
	if a > b {
		a, b = b, a
	}
	if !m.visualLine {
		cut = b + 1
		if cut > len(runes) {
			cut = len(runes)
		}
		return runes, a, cut
	}
	start = a
	for start > 0 && runes[start-1] != '\n' {
		start--
	}
	cut = b
	for cut < len(runes) && runes[cut] != '\n' {
		cut++
	}
	switch {
	case cut < len(runes):
		cut++ // include the trailing newline
	case start > 0:
		start-- // last line: take the preceding newline instead
	}
	return runes, start, cut
}

// visualDelete removes the selection into the register and leaves the cursor at
// the selection start.
func (m *Model) visualDelete() {
	runes, start, cut := m.visualSpan()
	if start >= cut {
		return
	}
	text := string(runes[start:cut])
	if m.visualLine {
		text = strings.Trim(text, "\n")
	}
	m.register, m.regLinewise = text, m.visualLine

	rest := append(append([]rune{}, runes[:start]...), runes[cut:]...)
	m.ta.SetValue(string(rest))
	idx := start
	if idx > len(rest) {
		idx = len(rest)
	}
	r, c := rowColOf(rest, idx)
	m.moveTo(r, c)
}

// visualYank copies the selection into the register (buffer unchanged) and, as
// in Vim, moves the cursor to the selection start.
func (m *Model) visualYank() {
	runes, start, cut := m.visualSpan()
	if start >= cut {
		return
	}
	text := string(runes[start:cut])
	if m.visualLine {
		text = strings.Trim(text, "\n")
	}
	m.register, m.regLinewise = text, m.visualLine
	r, c := rowColOf(runes, start)
	m.moveTo(r, c)
}

// visualIndent shifts every line the selection touches by one level and leaves
// the cursor on the first of them.
func (m *Model) visualIndent(delta int) {
	row, _ := m.cursorPos()
	r1, r2 := m.anchorRow, row
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	lines := strings.Split(m.ta.Value(), "\n")
	for i := r1; i <= r2 && i < len(lines); i++ {
		lines[i] = shiftLine(lines[i], delta)
	}
	m.ta.SetValue(strings.Join(lines, "\n"))
	m.moveTo(r1, 0)
}

// shiftLine indents (+1) or dedents (-1) a single line by one level.
func shiftLine(line string, delta int) string {
	if delta > 0 {
		if line == "" {
			return line // Vim leaves empty lines alone when indenting
		}
		return tabIndent + line
	}
	if strings.HasPrefix(line, "\t") {
		return line[1:]
	}
	for i := 0; i < len(tabIndent); i++ {
		if !strings.HasPrefix(line, " ") {
			break
		}
		line = line[1:]
	}
	return line
}

// --- line indentation (Vim's << and >>) ---

// indentLine shifts the current line one indent level: +1 prepends tabIndent
// (>>), -1 removes up to one level of leading whitespace — a tab or up to
// len(tabIndent) spaces (<<). The cursor stays on the same character.
func (m *Model) indentLine(delta int) {
	row, col := m.cursorPos()
	lines := strings.Split(m.ta.Value(), "\n")
	if row < 0 || row >= len(lines) {
		return
	}
	old := lines[row]
	lines[row] = shiftLine(old, delta)
	d := len(lines[row]) - len(old) // indent chars are ASCII, so bytes == runes
	if d == 0 {
		return
	}
	if col += d; col < 0 {
		col = 0
	}
	m.ta.SetValue(strings.Join(lines, "\n"))
	m.moveTo(row, col)
}

// --- Visual-mode rendering ---

// Visual-mode styles: the selection bar, the cursor block, and the gutter.
var (
	visualSelStyle    = lipgloss.NewStyle().Background(lipgloss.Color("24")).Foreground(lipgloss.Color("231"))
	visualCursorStyle = lipgloss.NewStyle().Reverse(true)
	visualGutterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	visualBadge       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("135")).Padding(0, 1)
)

// visualView renders the buffer with the active selection highlighted. The
// textarea cannot draw selections, so Visual mode takes over rendering: line
// numbers + width-aware soft wrap + a scrolled window that keeps the cursor
// visible. Editing is impossible while selecting, so the simpler renderer can't
// drift from the real buffer.
func (m Model) visualView() string {
	lines := strings.Split(m.ta.Value(), "\n")
	curRow, curCol := m.cursorPos()
	_, selStart, selCut := m.visualSpan()

	height := m.ta.Height()
	if height <= 0 {
		height = 10
	}
	gutterW := len(fmt.Sprintf("%d", len(lines))) + 2
	contentW := m.width - gutterW
	if contentW < 8 {
		contentW = 8
	}

	// Lay every logical line out into wrapped segments, tracking which visual
	// row holds the cursor so the window can follow it.
	type seg struct {
		lineIdx    int // logical line
		start, end int // rune range within the line
		first      bool
	}
	var segs []seg
	cursorSeg := 0
	for li, line := range lines {
		parts := wrapWidths([]rune(line), contentW)
		for si, p := range parts {
			s := seg{lineIdx: li, start: p[0], end: p[1], first: si == 0}
			if li == curRow && curCol >= p[0] && (curCol < p[1] || si == len(parts)-1) {
				cursorSeg = len(segs)
			}
			segs = append(segs, s)
		}
	}

	// Scroll: keep the cursor's row inside the window.
	top := 0
	if cursorSeg >= height {
		top = cursorSeg - height + 1
	}
	bottom := top + height
	if bottom > len(segs) {
		bottom = len(segs)
	}

	flatStarts := flatLineStarts(lines)
	var b strings.Builder
	for i := top; i < bottom; i++ {
		s := segs[i]
		if i > top {
			b.WriteString("\n")
		}
		gutter := strings.Repeat(" ", gutterW)
		if s.first {
			gutter = fmt.Sprintf(" %*d ", gutterW-2, s.lineIdx+1)
		}
		b.WriteString(visualGutterStyle.Render(gutter))
		b.WriteString(renderSelSegment(
			[]rune(lines[s.lineIdx]), s.start, s.end, flatStarts[s.lineIdx],
			selStart, selCut, s.lineIdx == curRow, curCol,
		))
	}
	// Pad to the textarea's height so the surrounding layout doesn't shift.
	for i := bottom - top; i < height; i++ {
		b.WriteString("\n")
	}
	return b.String()
}

// renderSelSegment renders one wrapped segment ([segStart, segEnd) of line),
// batching runs of equally-styled runes (normal / selected / cursor).
func renderSelSegment(line []rune, segStart, segEnd, flatStart, selStart, selCut int, cursorLine bool, curCol int) string {
	classOf := func(col int) int {
		if cursorLine && col == curCol {
			return 2 // cursor
		}
		if idx := flatStart + col; idx >= selStart && idx < selCut {
			return 1 // selected
		}
		return 0
	}
	styleOf := func(class int) *lipgloss.Style {
		switch class {
		case 1:
			return &visualSelStyle
		case 2:
			return &visualCursorStyle
		}
		return nil
	}

	var b strings.Builder
	runStart, runClass := segStart, -1
	flush := func(end int) {
		if runClass < 0 || end <= runStart {
			return
		}
		text := string(line[runStart:end])
		if st := styleOf(runClass); st != nil {
			text = st.Render(text)
		}
		b.WriteString(text)
	}
	for col := segStart; col < segEnd && col < len(line); col++ {
		c := classOf(col)
		if c != runClass {
			flush(col)
			runStart, runClass = col, c
		}
	}
	end := segEnd
	if end > len(line) {
		end = len(line)
	}
	flush(end)

	// The cursor (or a linewise selection) on the line-end position is past the
	// last rune; show it as a styled space so it stays visible.
	switch {
	case cursorLine && curCol >= len(line) && segEnd >= len(line):
		b.WriteString(visualCursorStyle.Render(" "))
	case len(line) == 0 && flatStart >= selStart && flatStart < selCut:
		b.WriteString(visualSelStyle.Render(" "))
	}
	return b.String()
}

// wrapWidths splits a line of runes into display-width-bounded [start, end)
// segments (soft wrap). A line always yields at least one (possibly empty)
// segment so empty lines still occupy a row.
func wrapWidths(line []rune, width int) [][2]int {
	if len(line) == 0 {
		return [][2]int{{0, 0}}
	}
	var out [][2]int
	start, w := 0, 0
	for i, r := range line {
		rw := runewidth.RuneWidth(r)
		if w+rw > width && i > start {
			out = append(out, [2]int{start, i})
			start, w = i, 0
		}
		w += rw
	}
	out = append(out, [2]int{start, len(line)})
	return out
}

// flatLineStarts returns each line's starting index in the flat rune view of
// the buffer (lines joined by '\n').
func flatLineStarts(lines []string) []int {
	starts := make([]int, len(lines))
	idx := 0
	for i, l := range lines {
		starts[i] = idx
		idx += len([]rune(l)) + 1
	}
	return starts
}

// --- unnamed register: delete/yank then paste ---

// captureDelete runs op (which must remove a contiguous chunk from the buffer)
// and stores the removed text in the register. linewise marks the register as
// whole-line content (pasted on its own line by p/P).
func (m *Model) captureDelete(linewise bool, op func()) {
	before := m.ta.Value()
	op()
	removed, ok := removedChunk(before, m.ta.Value())
	if !ok || removed == "" {
		return
	}
	if linewise {
		removed = strings.Trim(removed, "\n")
	}
	m.register = removed
	m.regLinewise = linewise
}

// removedChunk diffs two buffer states where after is before with one
// contiguous chunk removed, and returns that chunk.
func removedChunk(before, after string) (string, bool) {
	if len(after) >= len(before) {
		return "", false
	}
	p := 0
	for p < len(after) && before[p] == after[p] {
		p++
	}
	s := 0
	for s < len(after)-p && before[len(before)-1-s] == after[len(after)-1-s] {
		s++
	}
	return before[p : len(before)-s], true
}

// yankLine copies the current line into the register (Vim's yy).
func (m *Model) yankLine() {
	lines := strings.Split(m.ta.Value(), "\n")
	if r := m.ta.Line(); r >= 0 && r < len(lines) {
		m.register = lines[r]
		m.regLinewise = true
	}
}

// paste inserts the register at the cursor (p pastes after, P before). With an
// empty register it falls back to the system clipboard via the textarea's
// paste command — propagating the returned tea.Cmd, which a bare synthetic
// key-send would discard.
func (m Model) paste(after bool) (tea.Model, tea.Cmd) {
	if m.register == "" {
		var cmd tea.Cmd
		m.ta, cmd = m.ta.Update(tea.KeyMsg{Type: tea.KeyCtrlV})
		return m, cmd
	}
	if m.regLinewise {
		if after {
			m.ta.CursorEnd()
			m.ta.InsertString("\n" + m.register)
		} else {
			m.ta.CursorStart()
			m.ta.InsertString(m.register + "\n")
		}
		return m, nil
	}
	if after {
		m.arrow(tea.KeyRight)
	}
	m.ta.InsertString(m.register)
	return m, nil
}
