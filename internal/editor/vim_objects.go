package editor

// vim_objects.go adds range-based Vim operators: text objects (iw/aw, quotes,
// brackets), operator+motion for y (yw, y$, yiw…) on par with d/c, and the
// extra motions db/dj/dk. Everything is computed as a half-open flat-rune span
// [start, cut) and applied uniformly by operateCharRange / operateLineRange, so
// d, c, and y share one implementation.

import (
	"strings"
)

// lineBounds returns the [start, end) flat-rune span of the line containing idx,
// excluding the trailing newline.
func lineBounds(runes []rune, idx int) (start, end int) {
	if idx > len(runes) {
		idx = len(runes)
	}
	if idx < 0 {
		idx = 0
	}
	start = idx
	for start > 0 && runes[start-1] != '\n' {
		start--
	}
	end = idx
	for end < len(runes) && runes[end] != '\n' {
		end++
	}
	return start, end
}

// operateCharRange applies a charwise operator over the half-open flat span
// [start, cut): y yanks, d deletes, c deletes and enters Insert. The register is
// set for d/c (and the yank register for y).
func (m *Model) operateCharRange(op rune, start, cut int) {
	runes := []rune(m.ta.Value())
	if start < 0 {
		start = 0
	}
	if cut > len(runes) {
		cut = len(runes)
	}
	if start >= cut {
		if op == 'c' {
			m.mode = modeInsert
		}
		return
	}
	text := string(runes[start:cut])
	switch op {
	case 'y':
		m.setYank(text, false)
		r, c := rowColOf(runes, start)
		m.moveTo(r, c)
	case 'd', 'c':
		m.pushUndo() // a duplicate pre-edit snapshot is skipped by undo
		m.register, m.regLinewise = text, false
		rest := append(append([]rune{}, runes[:start]...), runes[cut:]...)
		m.ta.SetValue(string(rest))
		r, c := rowColOf(rest, start)
		m.moveTo(r, c)
		if op == 'c' {
			m.mode = modeInsert
		}
	}
}

// operateLineRange applies a linewise operator over rows [r1, r2] (inclusive).
func (m *Model) operateLineRange(op rune, r1, r2 int) {
	lines := strings.Split(m.ta.Value(), "\n")
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	if r1 < 0 {
		r1 = 0
	}
	if r2 >= len(lines) {
		r2 = len(lines) - 1
	}
	text := strings.Join(lines[r1:r2+1], "\n")
	switch op {
	case 'y':
		m.setYank(text, true)
		m.moveTo(r1, 0)
	case 'd':
		m.pushUndo()
		m.register, m.regLinewise = text, true
		rest := append(append([]string{}, lines[:r1]...), lines[r2+1:]...)
		if len(rest) == 0 {
			rest = []string{""}
		}
		m.ta.SetValue(strings.Join(rest, "\n"))
		nr := r1
		if nr >= len(rest) {
			nr = len(rest) - 1
		}
		m.moveTo(nr, 0)
	case 'c':
		m.pushUndo()
		m.register, m.regLinewise = text, true
		out := append([]string{}, lines[:r1]...)
		out = append(out, "")
		out = append(out, lines[r2+1:]...)
		m.ta.SetValue(strings.Join(out, "\n"))
		m.moveTo(r1, 0)
		m.mode = modeInsert
	}
}

// operatorMotionKey applies operator op over the motion named by key (the extra
// motions beyond the ones editor.go handles directly): b/B, h/l, j/k (linewise),
// and G (linewise to end).
func (m *Model) operatorMotionKey(op rune, key string) {
	runes := []rune(m.ta.Value())
	row, col := m.cursorPos()
	cur := flatIndex(runes, row, col)
	n := m.pendingCount
	if n < 1 {
		n = 1
	}
	switch key {
	case "j", "down":
		m.operateLineRange(op, row, row+n)
	case "k", "up":
		m.operateLineRange(op, row-n, row)
	case "G":
		m.operateLineRange(op, row, len(strings.Split(m.ta.Value(), "\n"))-1)
	default:
		if s, c, ok := charMotionRange(runes, cur, key, n); ok {
			m.operateCharRange(op, s, c)
		}
	}
}

// charMotionRange returns the half-open charwise span a motion covers from cur,
// clamped to the current line for the within-line motions.
func charMotionRange(runes []rune, cur int, key string, n int) (int, int, bool) {
	if n < 1 {
		n = 1
	}
	ls, le := lineBounds(runes, cur)
	switch key {
	case "w", "W":
		e := cur
		for i := 0; i < n; i++ {
			e = nextWordStart(runes, e)
		}
		if e <= cur {
			return 0, 0, false
		}
		return cur, e, true
	case "e", "E":
		e := cur
		for i := 0; i < n; i++ {
			e = nextWordEnd(runes, e)
		}
		if e+1 <= cur {
			return 0, 0, false
		}
		return cur, min(e+1, len(runes)), true
	case "b", "B":
		s := cur
		for i := 0; i < n; i++ {
			s = prevWordStart(runes, s)
		}
		if s >= cur {
			return 0, 0, false
		}
		return s, cur, true
	case "$":
		return cur, le, true
	case "0":
		return ls, cur, true
	case "^":
		fnb := ls
		for fnb < le && isWordSpace(runes[fnb]) {
			fnb++
		}
		if fnb <= cur {
			return fnb, cur, true
		}
		return cur, fnb, true
	case "h", "left":
		s := cur - n
		if s < ls {
			s = ls
		}
		return s, cur, true
	case "l", "right", " ", "space":
		e := cur + n
		if e > le {
			e = le
		}
		return cur, e, true
	}
	return 0, 0, false
}

// --- text objects ---

// applyTextObject runs operator op over the text object named by objKey (w, ",
// ', `, (, ), b, {, }, B, [, ], <, >). around selects a<obj> (include the
// delimiters / trailing space), otherwise i<obj>.
func (m *Model) applyTextObject(op rune, around bool, objKey string) {
	if op == 0 {
		return
	}
	runes := []rune(m.ta.Value())
	row, col := m.cursorPos()
	cur := flatIndex(runes, row, col)
	var start, cut int
	var ok bool
	switch objKey {
	case "w", "W":
		start, cut, ok = textObjectWord(runes, cur, around)
	case "\"", "'", "`":
		start, cut, ok = textObjectQuote(runes, cur, []rune(objKey)[0], around)
	case "(", ")", "b":
		start, cut, ok = pairObject(runes, cur, '(', ')', around)
	case "{", "}", "B":
		start, cut, ok = pairObject(runes, cur, '{', '}', around)
	case "[", "]":
		start, cut, ok = pairObject(runes, cur, '[', ']', around)
	case "<", ">":
		start, cut, ok = pairObject(runes, cur, '<', '>', around)
	}
	if ok {
		m.operateCharRange(op, start, cut)
	}
}

// textObjectWord returns the span of the word (or whitespace run) under the
// cursor. aw extends over the trailing whitespace, or the leading run if there
// is no trailing one — Vim's behavior. Runs never cross line breaks.
func textObjectWord(runes []rune, cur int, around bool) (int, int, bool) {
	n := len(runes)
	if n == 0 {
		return 0, 0, false
	}
	if cur >= n {
		cur = n - 1
	}
	if runes[cur] == '\n' {
		return 0, 0, false
	}
	space := isWordSpace(runes[cur])
	start, end := cur, cur
	for start > 0 && runes[start-1] != '\n' && isWordSpace(runes[start-1]) == space {
		start--
	}
	for end+1 < n && runes[end+1] != '\n' && isWordSpace(runes[end+1]) == space {
		end++
	}
	cut := end + 1
	if around && !space {
		grew := false
		for cut < n && runes[cut] != '\n' && isWordSpace(runes[cut]) {
			cut++
			grew = true
		}
		if !grew {
			for start > 0 && runes[start-1] != '\n' && isWordSpace(runes[start-1]) {
				start--
			}
		}
	}
	return start, cut, true
}

// textObjectQuote returns the span between (i) or including (a) the pair of q
// characters on the cursor's line that encloses — or first follows — the cursor.
func textObjectQuote(runes []rune, cur int, q rune, around bool) (int, int, bool) {
	ls, le := lineBounds(runes, cur)
	var pos []int
	for i := ls; i < le; i++ {
		if runes[i] == q {
			pos = append(pos, i)
		}
	}
	for k := 0; k+1 < len(pos); k += 2 {
		a, b := pos[k], pos[k+1]
		if cur <= b { // enclosing, or the first pair at/after the cursor
			if around {
				return a, b + 1, true
			}
			return a + 1, b, true
		}
	}
	return 0, 0, false
}

// pairObject returns the span inside (i) or including (a) the bracket pair that
// encloses the cursor. Brackets may span lines; nesting is respected.
func pairObject(runes []rune, cur int, open, close rune, around bool) (int, int, bool) {
	oi, ci, ok := enclosingPair(runes, cur, open, close)
	if !ok {
		return 0, 0, false
	}
	if around {
		return oi, ci + 1, true
	}
	return oi + 1, ci, true
}

// enclosingPair finds the innermost open/close bracket pair surrounding cur.
func enclosingPair(runes []rune, cur int, open, close rune) (int, int, bool) {
	n := len(runes)
	if n == 0 {
		return 0, 0, false
	}
	if cur >= n {
		cur = n - 1
	}
	depth, oi := 0, -1
	for i := cur; i >= 0; i-- {
		switch runes[i] {
		case close:
			if i != cur {
				depth++
			}
		case open:
			if depth > 0 {
				depth--
			} else {
				oi = i
			}
		}
		if oi >= 0 {
			break
		}
	}
	if oi < 0 {
		return 0, 0, false
	}
	depth, ci := 0, -1
	for i := oi + 1; i < n; i++ {
		switch runes[i] {
		case open:
			depth++
		case close:
			if depth > 0 {
				depth--
			} else {
				ci = i
			}
		}
		if ci >= 0 {
			break
		}
	}
	if ci < 0 {
		return 0, 0, false
	}
	return oi, ci, true
}

// deleteCharsBefore implements X: delete up to n characters before the cursor,
// not crossing the line start.
func (m *Model) deleteCharsBefore(n int) {
	if n < 1 {
		n = 1
	}
	runes := []rune(m.ta.Value())
	row, col := m.cursorPos()
	cur := flatIndex(runes, row, col)
	ls, _ := lineBounds(runes, cur)
	start := cur - n
	if start < ls {
		start = ls
	}
	if start >= cur {
		return
	}
	m.pushUndo()
	m.register, m.regLinewise = string(runes[start:cur]), false
	rest := append(append([]rune{}, runes[:start]...), runes[cur:]...)
	m.ta.SetValue(string(rest))
	r, c := rowColOf(rest, start)
	m.moveTo(r, c)
}

// searchRepeat implements n (sign +1) and N (sign -1), honoring the direction
// of the last search (?, # reverse it).
func (m *Model) searchRepeat(sign int) {
	if m.lastSearch == "" {
		return
	}
	dir := sign
	if m.searchReverse {
		dir = -sign
	}
	if from := m.curJumpPos(); m.search(m.lastSearch, dir) {
		m.recordJump(from)
	} else {
		m.status = "pattern not found: " + m.lastSearch
	}
}

// searchWordUnderCursor implements * and #: search for the word under the cursor
// forward (backward=false) or backward.
func (m *Model) searchWordUnderCursor(backward bool) {
	runes := []rune(m.ta.Value())
	row, col := m.cursorPos()
	cur := flatIndex(runes, row, col)
	if cur >= len(runes) || isWordSpace(runes[cur]) {
		return
	}
	s, e := cur, cur
	for s > 0 && runes[s-1] != '\n' && !isWordSpace(runes[s-1]) {
		s--
	}
	for e+1 < len(runes) && runes[e+1] != '\n' && !isWordSpace(runes[e+1]) {
		e++
	}
	word := string(runes[s : e+1])
	if strings.TrimSpace(word) == "" {
		return
	}
	m.lastSearch = word
	m.searchReverse = backward
	dir := 1
	if backward {
		dir = -1
	}
	from := m.curJumpPos()
	if m.search(word, dir) {
		m.recordJump(from)
	} else {
		m.status = "pattern not found: " + word
	}
}
