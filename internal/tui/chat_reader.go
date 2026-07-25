package tui

// chat_reader.go gives the read-only lecture pane (a chatModel with noInput) a
// Vim-style keyboard cursor: Normal-mode motions, a Visual selection that reuses
// the mouse drag-select machinery (chatSelection / copySelection), an incremental
// "/" search with n/N, and gd to follow a [[wikilink]] under the cursor. The
// cursor lives in content coordinates — curLine indexes contentLines (the
// rendered transcript) and curCol is a rune index into the ANSI-stripped display
// line — so it survives scrolling exactly like the drag selection does.

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	runewidth "github.com/mattn/go-runewidth"
)

var wikilinkRe = regexp.MustCompile(`\[\[[^\]\n]+\]\]`)

func isSpaceRune(r rune) bool { return unicode.IsSpace(r) }

// lineText is the plain (ANSI-stripped) text of a content line, or "" if out of
// range — the surface all cursor math runs over.
func (c chatModel) lineText(i int) string {
	if i < 0 || i >= len(c.contentLines) {
		return ""
	}
	return ansi.Strip(c.contentLines[i])
}

func (c chatModel) curRunes() []rune { return []rune(c.lineText(c.curLine)) }

// cursorCell maps the cursor's (line, rune-col) to (line, display-cell col), the
// coordinate space chatSelection / overlaySelection speak.
func (c chatModel) cursorCell() (line, cell int) {
	runes := c.curRunes()
	col := c.curCol
	if col > len(runes) {
		col = len(runes)
	}
	return c.curLine, runewidth.StringWidth(string(runes[:col]))
}

// clampCursor keeps the cursor on a real cell: a valid line, and a column within
// the line (on the last character, Vim-style; column 0 on an empty line).
func (c *chatModel) clampCursor() {
	if len(c.contentLines) == 0 {
		c.curLine, c.curCol = 0, 0
		return
	}
	c.curLine = clampInt(c.curLine, 0, len(c.contentLines)-1)
	n := len(c.curRunes())
	if n == 0 {
		c.curCol = 0
		return
	}
	c.curCol = clampInt(c.curCol, 0, n-1)
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// --- key dispatch ---

// readerKey handles one keystroke of the lecture pane's reader. Motions move the
// cursor; v toggles a Visual selection; y copies it; gd follows a wikilink; n/N
// repeat the last search. Anything unrecognized falls through to page scrolling.
func (c *chatModel) readerKey(msg tea.KeyMsg) {
	c.readerNotice = ""
	key := msg.String()
	// A pending f/F/t/T consumes the next key as its target character.
	if c.readerFindPending != 0 {
		op := c.readerFindPending
		c.readerFindPending = 0
		if len(msg.Runes) == 1 {
			c.readerFind(op, msg.Runes[0], max(1, c.readerCount))
			c.readerLastFindOp, c.readerLastFindCh = op, msg.Runes[0]
		}
		c.readerCount = 0
		c.afterMotion()
		return
	}
	if c.pendingReaderG {
		c.pendingReaderG = false
		c.readerCount = 0
		switch key {
		case "g": // gg
			c.gotoDocTop()
			c.afterMotion()
		case "d": // gd: follow the [[wikilink]] under the cursor
			if t, ok := c.wikilinkAtCursor(); ok {
				c.followTarget = t
			} else {
				c.readerNotice = "no [[link]] under the cursor"
			}
		}
		return
	}
	// Numeric prefix: 1-9 always, 0 only to extend an in-progress count (else 0
	// is the line-start motion).
	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' || (key == "0" && c.readerCount > 0) {
		c.readerCount = c.readerCount*10 + int(key[0]-'0')
		return
	}
	n := max(1, c.readerCount)
	c.readerCount = 0
	switch key {
	case "h", "left":
		c.curCol -= n
	case "l", "right":
		c.curCol += n
	case "j", "down":
		c.curLine += n
	case "k", "up":
		c.curLine -= n
	case "w":
		for i := 0; i < n; i++ {
			c.wordForward()
		}
	case "e":
		for i := 0; i < n; i++ {
			c.wordEnd()
		}
	case "b":
		for i := 0; i < n; i++ {
			c.wordBack()
		}
	case "f", "F", "t", "T":
		c.readerFindPending = []rune(key)[0]
		c.readerCount = n // preserve the count for the target key
		return
	case ";":
		c.repeatReaderFind(false, n)
	case ",":
		c.repeatReaderFind(true, n)
	case "0":
		c.curCol = 0
	case "^":
		c.firstNonBlank()
	case "$":
		c.curCol = len(c.curRunes()) // clamped to last char in afterMotion
	case "g":
		c.pendingReaderG = true
		return
	case "G", "end":
		c.gotoDocBottom()
	case "home":
		c.gotoDocTop()
	case "{":
		c.paragraphJump(-1)
		c.curLine, c.curCol = c.vp.YOffset, 0
	case "}":
		c.paragraphJump(1)
		c.curLine, c.curCol = c.vp.YOffset, 0
	case "v":
		if c.docVisual {
			c.docVisual = false
			c.clearSelect()
		} else {
			c.startDocVisual()
		}
	case "esc":
		if c.docVisual {
			c.docVisual = false
			c.clearSelect()
		}
		return
	case "y":
		if c.docVisual {
			c.readerNotice = copySelection(c)
			c.docVisual = false
			c.clearSelect()
		}
		return
	case "n":
		c.search("", true)
		return
	case "N":
		c.search("", false)
		return
	case "ctrl+e": // scroll the view down a line, cursor following (Vim ⌃e)
		for i := 0; i < n; i++ {
			c.vp.ScrollDown(1)
		}
		c.clampCursorToView()
		return
	case "ctrl+y": // scroll up a line (Vim ⌃y)
		for i := 0; i < n; i++ {
			c.vp.ScrollUp(1)
		}
		c.clampCursorToView()
		return
	default:
		if c.scrollKey(msg) { // ⌃d/⌃u/⌃f/⌃b, PgUp/PgDn, Shift-arrows
			c.clampCursorToView()
		}
		return
	}
	c.afterMotion()
}

// afterMotion normalizes state after a cursor move: clamp, drop any stale search
// highlight, extend the Visual selection, and scroll to keep the cursor visible.
func (c *chatModel) afterMotion() {
	c.clampCursor()
	c.searchMatchCol, c.searchMatchLen = -1, 0
	if c.docVisual {
		c.updateDocSelectionHead()
	}
	c.scrollToCursor()
}

// readerFind moves the cursor to the n-th f/F/t/T target on the current line;
// on a miss the cursor stays put, as in Vim.
func (c *chatModel) readerFind(op, ch rune, n int) {
	runes := c.curRunes()
	col := c.curCol
	for k := 0; k < n; k++ {
		next := findCharInLine(runes, col, op, ch)
		if next < 0 {
			return
		}
		col = next
	}
	c.curCol = col
}

// findCharInLine returns the column an f/F/t/T motion lands on from col, or -1
// when the character isn't found on this line.
func findCharInLine(runes []rune, col int, op, ch rune) int {
	switch op {
	case 'f':
		for i := col + 1; i < len(runes); i++ {
			if runes[i] == ch {
				return i
			}
		}
	case 't':
		for i := col + 1; i < len(runes); i++ {
			if runes[i] == ch {
				return i - 1
			}
		}
	case 'F':
		for i := col - 1; i >= 0; i-- {
			if runes[i] == ch {
				return i
			}
		}
	case 'T':
		for i := col - 1; i >= 0; i-- {
			if runes[i] == ch {
				return i + 1
			}
		}
	}
	return -1
}

// repeatReaderFind repeats the last f/F/t/T (; keeps direction, , reverses it).
func (c *chatModel) repeatReaderFind(reverse bool, n int) {
	if c.readerLastFindOp == 0 {
		return
	}
	op := c.readerLastFindOp
	if reverse {
		op = flipFind(op)
	}
	c.readerFind(op, c.readerLastFindCh, n)
}

func flipFind(op rune) rune {
	switch op {
	case 'f':
		return 'F'
	case 'F':
		return 'f'
	case 't':
		return 'T'
	case 'T':
		return 't'
	}
	return op
}

func (c *chatModel) firstNonBlank() {
	runes := c.curRunes()
	i := 0
	for i < len(runes) && isSpaceRune(runes[i]) {
		i++
	}
	c.curCol = i
}

func (c *chatModel) gotoDocTop()    { c.curLine, c.curCol = 0, 0 }
func (c *chatModel) gotoDocBottom() { c.curLine, c.curCol = len(c.contentLines)-1, 0 }

// --- word motions (whitespace-delimited, spanning line breaks) ---

func (c *chatModel) wordForward() {
	runes := c.curRunes()
	i := c.curCol
	for i < len(runes) && !isSpaceRune(runes[i]) {
		i++
	}
	for i < len(runes) && isSpaceRune(runes[i]) {
		i++
	}
	if i < len(runes) {
		c.curCol = i
		return
	}
	for li := c.curLine + 1; li < len(c.contentLines); li++ {
		nr := []rune(c.lineText(li))
		if len(nr) == 0 {
			c.curLine, c.curCol = li, 0
			return
		}
		j := 0
		for j < len(nr) && isSpaceRune(nr[j]) {
			j++
		}
		if j < len(nr) {
			c.curLine, c.curCol = li, j
			return
		}
	}
	c.curCol = len(runes)
}

func (c *chatModel) wordBack() {
	runes := c.curRunes()
	i := c.curCol - 1
	for i >= 0 && isSpaceRune(runes[i]) {
		i--
	}
	if i >= 0 {
		for i > 0 && !isSpaceRune(runes[i-1]) {
			i--
		}
		c.curCol = i
		return
	}
	for li := c.curLine - 1; li >= 0; li-- {
		pr := []rune(c.lineText(li))
		if len(pr) == 0 {
			c.curLine, c.curCol = li, 0
			return
		}
		k := len(pr) - 1
		for k >= 0 && isSpaceRune(pr[k]) {
			k--
		}
		if k >= 0 {
			for k > 0 && !isSpaceRune(pr[k-1]) {
				k--
			}
			c.curLine, c.curCol = li, k
			return
		}
	}
	c.curCol = 0
}

func (c *chatModel) wordEnd() {
	runes := c.curRunes()
	i := c.curCol + 1
	for i < len(runes) && isSpaceRune(runes[i]) {
		i++
	}
	if i < len(runes) {
		for i+1 < len(runes) && !isSpaceRune(runes[i+1]) {
			i++
		}
		c.curCol = i
		return
	}
	for li := c.curLine + 1; li < len(c.contentLines); li++ {
		nr := []rune(c.lineText(li))
		j := 0
		for j < len(nr) && isSpaceRune(nr[j]) {
			j++
		}
		if j < len(nr) {
			for j+1 < len(nr) && !isSpaceRune(nr[j+1]) {
				j++
			}
			c.curLine, c.curCol = li, j
			return
		}
	}
	c.curCol = len(runes)
}

// --- scrolling to follow the cursor ---

func (c *chatModel) scrollToCursor() {
	h := c.vp.Height
	if h <= 0 {
		return
	}
	top := c.vp.YOffset
	if c.curLine < top {
		c.vp.SetYOffset(c.curLine)
	} else if c.curLine > top+h-1 {
		c.vp.SetYOffset(c.curLine - h + 1)
	}
}

// clampCursorToView pulls the cursor back onto the visible page after a scroll
// key moved the viewport out from under it.
func (c *chatModel) clampCursorToView() {
	top := c.vp.YOffset
	bot := top + c.vp.Height - 1
	c.curLine = clampInt(c.curLine, top, bot)
	c.clampCursor()
	if c.docVisual {
		c.updateDocSelectionHead()
	}
}

// --- Visual selection (reuses the mouse chatSelection) ---

func (c *chatModel) startDocVisual() {
	c.docVisual = true
	line, cell := c.cursorCell()
	c.sel = chatSelection{anchorLine: line, anchorCol: cell, headLine: line, headCol: cell, active: true}
}

func (c *chatModel) updateDocSelectionHead() {
	line, cell := c.cursorCell()
	c.sel.headLine, c.sel.headCol = line, cell
	c.sel.active = true
}

// --- "/" search over the rendered lines ---

// search moves the cursor to the next (or previous) case-insensitive match of q,
// wrapping around the document. An empty q reuses the last pattern (n / N).
// Reports whether a match was found.
func (c *chatModel) search(q string, forward bool) bool {
	if strings.TrimSpace(q) != "" {
		c.searchQuery = q
	}
	c.searchMatchCol, c.searchMatchLen = -1, 0
	if c.searchQuery == "" || len(c.contentLines) == 0 {
		return false
	}
	needle := strings.ToLower(c.searchQuery)
	n := len(c.contentLines)
	if forward {
		for step := 0; step <= n; step++ {
			li := (c.curLine + step) % n
			line := strings.ToLower(c.lineText(li))
			from := 0
			if step == 0 {
				from = byteOfRune(c.lineText(li), c.curCol+1)
			}
			if from <= len(line) {
				if rel := strings.Index(line[from:], needle); rel >= 0 {
					c.setMatch(li, from+rel)
					return true
				}
			}
		}
	} else {
		for step := 0; step <= n; step++ {
			li := ((c.curLine-step)%n + n) % n
			line := strings.ToLower(c.lineText(li))
			limit := len(line)
			if step == 0 {
				limit = byteOfRune(c.lineText(li), c.curCol)
			}
			if limit >= 0 && limit <= len(line) {
				if idx := strings.LastIndex(line[:limit], needle); idx >= 0 {
					c.setMatch(li, idx)
					return true
				}
			}
		}
	}
	return false
}

// setMatch places the cursor on a match found at byte offset byteIdx of line li
// and records its span for highlighting.
func (c *chatModel) setMatch(li, byteIdx int) {
	line := c.lineText(li)
	c.curLine = li
	c.curCol = utf8.RuneCountInString(line[:byteIdx])
	c.searchMatchCol = c.curCol
	c.searchMatchLen = utf8.RuneCountInString(c.searchQuery)
	c.clampCursor()
	if c.docVisual {
		c.updateDocSelectionHead()
	}
	c.scrollToCursor()
}

// byteOfRune returns the byte offset of the runeCol-th rune of s (clamped to the
// string bounds).
func byteOfRune(s string, runeCol int) int {
	if runeCol <= 0 {
		return 0
	}
	i, seen := 0, 0
	for i < len(s) {
		if seen == runeCol {
			return i
		}
		_, sz := utf8.DecodeRuneInString(s[i:])
		i += sz
		seen++
	}
	return len(s)
}

// --- gd: wikilink under the cursor ---

// wikilinkAtCursor returns the target of the [[wikilink]] the cursor sits inside
// (alias after "|" stripped), if any. Same-line only: a link soft-wrapped across
// two rendered rows won't resolve.
func (c chatModel) wikilinkAtCursor() (string, bool) {
	line := c.lineText(c.curLine)
	for _, m := range wikilinkRe.FindAllStringIndex(line, -1) {
		loRune := utf8.RuneCountInString(line[:m[0]])
		hiRune := utf8.RuneCountInString(line[:m[1]])
		if c.curCol >= loRune && c.curCol < hiRune {
			inner := strings.TrimSpace(line[m[0]+2 : m[1]-2])
			if i := strings.IndexByte(inner, '|'); i >= 0 {
				inner = strings.TrimSpace(inner[:i])
			}
			if inner != "" {
				return inner, true
			}
		}
	}
	return "", false
}

// --- rendering the cursor and current search match ---

// overlayCursor paints the reader cursor cell (and the current "/" match) over
// the rendered viewport rows, mirroring overlaySelection's ANSI-aware cutting.
// Only drawn when the pane is focused.
func (c chatModel) overlayCursor(view string) string {
	if !c.focused || len(c.contentLines) == 0 {
		return view
	}
	rows := strings.Split(view, "\n")
	paint := func(absLine, cellLo, cellHi int, st lipgloss.Style) {
		r := absLine - c.vp.YOffset
		if r < 0 || r >= len(rows) {
			return
		}
		row := rows[r]
		w := ansi.StringWidth(row)
		if cellLo >= w { // cursor past the rendered text: show it as a trailing block
			if cellLo == w {
				rows[r] = row + st.Render(" ")
			}
			return
		}
		if cellHi > w {
			cellHi = w
		}
		if cellLo >= cellHi {
			return
		}
		rows[r] = ansi.Cut(row, 0, cellLo) +
			st.Render(ansi.Strip(ansi.Cut(row, cellLo, cellHi))) +
			ansi.Cut(row, cellHi, w)
	}
	if c.searchMatchCol >= 0 && c.searchMatchLen > 0 {
		runes := c.curRunes()
		lo := runewidth.StringWidth(string(runes[:min(c.searchMatchCol, len(runes))]))
		hi := runewidth.StringWidth(string(runes[:min(c.searchMatchCol+c.searchMatchLen, len(runes))]))
		paint(c.curLine, lo, hi, chatSearchStyle)
	}
	if !c.docVisual {
		_, cell := c.cursorCell()
		runes := c.curRunes()
		wcell := 1
		if c.curCol < len(runes) {
			wcell = max(1, runewidth.RuneWidth(runes[c.curCol]))
		}
		paint(c.curLine, cell, cell+wcell, chatCursorStyle)
	}
	return strings.Join(rows, "\n")
}
