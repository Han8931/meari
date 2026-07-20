package tui

// chat_vim.go extends the chat input's light Vim layer with a character-wise
// Visual mode (v → motions → y/d/x/c) plus gg/G first/last-line jumps. Like
// the in-app editor, it drives the bubbles textarea from outside: motions move
// the real cursor (sent as keys or via inputMoveTo), the selection span is
// tracked as flat rune indices into Value(), and operations rewrite the value.
//
// The selection highlight is painted over the rendered input (see
// paintInputSelection). Painting needs the textarea's soft-wrap layout, which
// is only reproducible when the whole value fits in the visible rows — the
// overwhelmingly common case for a 3-row chat input. Taller values keep
// correct motions and operations; only the paint is skipped.

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	runewidth "github.com/mattn/go-runewidth"
)

// --- flat-index math over the input value (rows separated by '\n') ---

func inputFlatIndex(runes []rune, row, col int) int {
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

func inputRowColOf(runes []rune, idx int) (row, col int) {
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

// cursorFlat returns the input cursor as a flat rune index. LineInfo reports
// the visual position within a soft-wrapped row; StartColumn + ColumnOffset
// recovers the logical column (same recipe as the editor's cursorPos).
func (c *chatModel) cursorFlat() int {
	li := c.input.LineInfo()
	return inputFlatIndex([]rune(c.input.Value()), c.input.Line(), li.StartColumn+li.ColumnOffset)
}

// chatReposition is a no-op message pumped through the textarea so it scrolls
// its viewport to the cursor after a long jump (it only does that in Update).
type chatReposition struct{}

// inputMoveTo places the input cursor at logical (row, col), stepping visual
// lines like the editor's moveTo.
func (c *chatModel) inputMoveTo(row, col int) {
	for guard := 0; c.input.Line() < row && guard < 10000; guard++ {
		before := c.input.Line()
		c.input.CursorDown()
		if c.input.Line() == before && c.input.LineInfo().RowOffset == 0 {
			break
		}
	}
	for guard := 0; c.input.Line() > row && guard < 10000; guard++ {
		before := c.input.Line()
		c.input.CursorUp()
		if c.input.Line() == before && c.input.LineInfo().RowOffset == 0 {
			break
		}
	}
	c.input.SetCursor(col)
	c.input, _ = c.input.Update(chatReposition{})
}

// inputFirstLine / inputLastLine are Vim's gg and G over the input buffer.
func (c *chatModel) inputFirstLine() { c.inputMoveTo(0, 0) }

func (c *chatModel) inputLastLine() {
	row, _ := inputRowColOf([]rune(c.input.Value()), len([]rune(c.input.Value())))
	c.inputMoveTo(row, 0)
}

// --- Visual mode ---

// enterVisual anchors a character-wise selection at the cursor.
func (c *chatModel) enterVisual() {
	c.visual = true
	c.vAnchor = c.cursorFlat()
}

// visualSpanInput returns the value's runes and the selection as a half-open
// range [start, cut) — inclusive of the rune under the cursor, Vim-style.
func (c *chatModel) visualSpanInput() (runes []rune, start, cut int) {
	runes = []rune(c.input.Value())
	a, b := c.vAnchor, c.cursorFlat()
	if a > b {
		a, b = b, a
	}
	if a < 0 {
		a = 0
	}
	if b < len(runes) {
		b++
	}
	if b > len(runes) {
		b = len(runes)
	}
	return runes, a, b
}

// visualYankInput copies the selection to the clipboard and, as in Vim,
// leaves the cursor at the selection start.
func (c *chatModel) visualYankInput() string {
	runes, start, cut := c.visualSpanInput()
	c.visual = false
	if start >= cut {
		return ""
	}
	text := string(runes[start:cut])
	notice := "✓ yanked " + itoa(len(runes[start:cut])) + " chars"
	if err := copyToClipboard(text); err != nil {
		notice = "✓ yanked to the terminal clipboard (OSC 52)"
	}
	r, col := inputRowColOf(runes, start)
	c.inputMoveTo(r, col)
	return notice
}

// visualDeleteInput removes the selection (d/x); change=true re-enters Insert
// mode afterwards (c).
func (c *chatModel) visualDeleteInput(change bool) {
	runes, start, cut := c.visualSpanInput()
	c.visual = false
	if start >= cut {
		return
	}
	rest := append(append([]rune{}, runes[:start]...), runes[cut:]...)
	c.input.SetValue(string(rest))
	r, col := inputRowColOf(rest, min(start, len(rest)))
	c.inputMoveTo(r, col)
	if change {
		c.exitNormal()
	}
}

// visualKey handles one keystroke of Visual mode. Motions move the real
// cursor so the selection end is always visible.
func (c *chatModel) visualKey(msg tea.KeyMsg) string {
	send := func(t tea.KeyType) { c.input, _ = c.input.Update(tea.KeyMsg{Type: t}) }
	alt := func(r rune) {
		c.input, _ = c.input.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}, Alt: true})
	}
	if c.pendingG {
		c.pendingG = false
		if msg.String() == "g" {
			c.inputFirstLine()
		}
		return ""
	}
	switch msg.String() {
	case "esc", "v":
		c.visual = false
	// --- motions (extend the selection) ---
	case "h", "left":
		send(tea.KeyLeft)
	case "l", "right":
		send(tea.KeyRight)
	case "j", "down":
		send(tea.KeyDown)
	case "k", "up":
		send(tea.KeyUp)
	case "w", "e":
		alt('f')
	case "b":
		alt('b')
	case "0", "^":
		send(tea.KeyHome)
	case "$":
		send(tea.KeyEnd)
	case "g":
		c.pendingG = true
	case "G":
		c.inputLastLine()
	case "o": // swap selection ends
		a := c.vAnchor
		c.vAnchor = c.cursorFlat()
		r, col := inputRowColOf([]rune(c.input.Value()), a)
		c.inputMoveTo(r, col)
	// --- operations ---
	case "y":
		return c.visualYankInput()
	case "d", "x":
		c.visualDeleteInput(false)
	case "c":
		c.visualDeleteInput(true)
	}
	return ""
}

// --- selection painting ---

// wrapInput replicates the textarea's word-wrap closely enough to locate each
// wrapped row's rune span (cell math via go-runewidth, which the rest of the
// TUI already standardizes on). Trailing-space bookkeeping at wrap boundaries
// may drift a cell in pathological cases; operations never depend on this.
func wrapInput(runes []rune, width int) [][]rune {
	if width <= 0 {
		return [][]rune{runes}
	}
	lines := [][]rune{{}}
	var word []rune
	row, spaces := 0, 0
	cellW := func(rs []rune) int { return runewidth.StringWidth(string(rs)) }

	for _, r := range runes {
		if r == ' ' {
			spaces++
		} else {
			word = append(word, r)
		}
		if spaces > 0 {
			if cellW(lines[row])+cellW(word)+spaces > width {
				row++
				lines = append(lines, append(append([]rune{}, word...), spacesOf(spaces)...))
			} else {
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], spacesOf(spaces)...)
			}
			spaces, word = 0, nil
		} else if len(word) > 0 && cellW(lines[row])+cellW(word) > width {
			if len(lines[row]) > 0 {
				row++
				lines = append(lines, []rune{})
			}
			lines[row] = append(lines[row], word...)
			word = nil
		}
	}
	lines[row] = append(lines[row], word...)
	lines[row] = append(lines[row], spacesOf(spaces)...)
	return lines
}

func spacesOf(n int) []rune {
	s := make([]rune, n)
	for i := range s {
		s[i] = ' '
	}
	return s
}

// inputScreenRows lays the whole value out as the textarea renders it:
// logical lines split by '\n', each soft-wrapped, each row carrying its flat
// rune span [start, end).
type inputScreenRow struct {
	start, end int // flat rune span (newline separators excluded)
	runes      []rune
}

func (c *chatModel) inputScreenRows() []inputScreenRow {
	width := c.input.Width()
	var rows []inputScreenRow
	flat := 0
	for i, line := range strings.Split(c.input.Value(), "\n") {
		if i > 0 {
			flat++ // the '\n' separator between logical lines
		}
		lr := []rune(line)
		off := 0
		for _, wr := range wrapInput(lr, width) {
			rows = append(rows, inputScreenRow{start: flat + off, end: flat + off + len(wr), runes: wr})
			off += len(wr)
		}
		flat += len(lr)
	}
	return rows
}

// paintInputSelection overlays the Visual selection on the rendered input
// lines (the transcript overlay's ansi.Cut technique). Only possible when the
// whole value fits the visible rows — then the textarea's viewport is at the
// top and screen rows map 1:1 to inputScreenRows.
func (c chatModel) paintInputSelection(lines []string) {
	rows := c.inputScreenRows()
	if len(rows) > c.input.Height() {
		return
	}
	_, selStart, selCut := (&c).visualSpanInput()
	if selStart >= selCut {
		return
	}
	const promptW = 2 // the "> " / "  " prompt column
	for i := range lines {
		if i >= len(rows) {
			break
		}
		row := rows[i]
		lo, hi := max(row.start, selStart), min(row.end, selCut)
		if lo >= hi {
			continue
		}
		loCol := promptW + runewidth.StringWidth(string(row.runes[:lo-row.start]))
		hiCol := promptW + runewidth.StringWidth(string(row.runes[:hi-row.start]))
		w := ansi.StringWidth(lines[i])
		if loCol >= w {
			continue
		}
		if hiCol > w {
			hiCol = w
		}
		lines[i] = ansi.Cut(lines[i], 0, loCol) +
			chatSelStyle.Render(ansi.Strip(ansi.Cut(lines[i], loCol, hiCol))) +
			ansi.Cut(lines[i], hiCol, w)
	}
}
