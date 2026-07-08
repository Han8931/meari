package editor

// table.go renders GFM tables ("| a | b |" rows under a "| --- | --- |"
// separator) as Unicode box-drawing grids for the chat pane, where the prose
// path's word-wrap would otherwise destroy column alignment. The renderer is
// width-aware: columns get their natural width when it fits, and shrink —
// wrapping cell text within the column — when it doesn't. Cell interiors go
// through mdInlineSeg, so `code`, **bold**, and [[wikilinks]] stay styled
// inside cells.
//
// The note editor never re-boxes tables (its buffer must stay byte-identical
// for cursor mapping); it only tints pipes — see mdTableEditorRow in
// markdown.go, which shares this file's separator predicate.

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// TableAlign is a column's alignment, from the separator row's colons.
type TableAlign int

const (
	AlignDefault TableAlign = iota // no colons: left, like GFM renderers
	AlignLeft                      // :---
	AlignCenter                    // :---:
	AlignRight                     // ---:
)

// mdTableBorder tints the grid's box-drawing characters (and the editor's
// pipe characters) the same dimness as horizontal rules.
var mdTableBorder = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

// IsTableSeparator reports whether the trimmed line is a GFM header separator
// row: cells of dashes with optional alignment colons (---, :---, :---:,
// ---:) split by pipes. At least one pipe is required, so a bare "---"
// stays a thematic break.
func IsTableSeparator(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, "|") {
		return false
	}
	cells := splitTableRow(line)
	if len(cells) == 0 {
		return false
	}
	for _, c := range cells {
		if _, ok := parseAlign(c); !ok {
			return false
		}
	}
	return true
}

// parseAlign classifies one separator cell, rejecting anything that isn't
// dashes with optional flanking colons.
func parseAlign(cell string) (TableAlign, bool) {
	cell = strings.TrimSpace(cell)
	left := strings.HasPrefix(cell, ":")
	right := strings.HasSuffix(cell, ":") && len(cell) > 1
	dashes := strings.TrimSuffix(strings.TrimPrefix(cell, ":"), ":")
	if dashes == "" || strings.Trim(dashes, "-") != "" {
		return AlignDefault, false
	}
	switch {
	case left && right:
		return AlignCenter, true
	case right:
		return AlignRight, true
	case left:
		return AlignLeft, true
	}
	return AlignDefault, true
}

// splitTableRow splits a table line into trimmed cell texts: optional outer
// pipes drop, unescaped pipes split, and "\|" unescapes to a literal pipe.
// Per the GFM spec, pipes inside code spans DO split cells — only backslash
// escapes them — so a plain byte scan is spec-correct.
func splitTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	var cells []string
	var cur strings.Builder
	for i := 0; i < len(line); i++ {
		switch {
		case line[i] == '\\' && i+1 < len(line) && line[i+1] == '|':
			cur.WriteByte('|')
			i++
		case line[i] == '|':
			cells = append(cells, strings.TrimSpace(cur.String()))
			cur.Reset()
		default:
			cur.WriteByte(line[i])
		}
	}
	// Text after the last pipe is a final cell only if non-empty (a trailing
	// "|" is the row's closing border, not an empty cell).
	if tail := strings.TrimSpace(cur.String()); tail != "" {
		cells = append(cells, tail)
	}
	return cells
}

// RenderTable renders GFM table lines (header, separator, body rows) as a
// box-drawing grid fitted to width — one string per visual row, every row's
// display width <= width. ok is false when lines don't parse as a table, so
// the caller can fall back to the prose path.
func RenderTable(lines []string, width int) (rows []string, ok bool) {
	if len(lines) < 2 || !IsTableSeparator(lines[1]) {
		return nil, false
	}
	header := splitTableRow(lines[0])
	if len(header) == 0 {
		return nil, false
	}
	aligns := make([]TableAlign, len(header))
	for i, c := range splitTableRow(strings.TrimSpace(lines[1])) {
		if i >= len(aligns) {
			break // excess separator cells: ragged, drop
		}
		aligns[i], _ = parseAlign(c)
	}

	// Normalize every row to the header's column count: missing cells are
	// empty, excess cells drop (GFM's ragged-row rule).
	ncols := len(header)
	cellRows := [][]string{header}
	for _, ln := range lines[2:] {
		cellRows = append(cellRows, splitTableRow(ln))
	}
	for r := range cellRows {
		for len(cellRows[r]) < ncols {
			cellRows[r] = append(cellRows[r], "")
		}
		cellRows[r] = cellRows[r][:ncols]
	}

	// Natural width per column; floor 1 so an all-empty column still draws.
	natural := make([]int, ncols)
	for _, row := range cellRows {
		for i, cell := range row {
			if w := lipgloss.Width(cell); w > natural[i] {
				natural[i] = w
			}
		}
	}
	for i := range natural {
		if natural[i] < 1 {
			natural[i] = 1
		}
	}
	widths := fitColumns(natural, width-(3*ncols+1)) // borders + "│ c │" padding

	var out []string
	emit := func(row string) { out = append(out, ansi.Hardwrap(row, width, true)) }
	emit(tableBorderRow("┌", "┬", "┐", widths))
	for r, row := range cellRows {
		for _, visual := range tableTextRows(row, widths, aligns, r == 0) {
			emit(visual)
		}
		if r == 0 {
			emit(tableBorderRow("├", "┼", "┤", widths))
		}
	}
	emit(tableBorderRow("└", "┴", "┘", widths))
	return out, true
}

// fitColumns picks display widths: natural when they fit, otherwise a
// proportional shrink with a floor so no column vanishes. avail excludes all
// border/padding chrome.
func fitColumns(natural []int, avail int) []int {
	total := 0
	widest := 0
	for i, n := range natural {
		total += n
		if n > natural[widest] {
			widest = i
		}
	}
	if total <= avail {
		return natural
	}
	// A small overflow comes out of the widest column alone (kept at or above
	// its fair share), so a 3-cell debt never wraps every header. Larger
	// overflows shrink all columns proportionally below.
	if fair := avail / len(natural); total-avail <= natural[widest]-fair {
		widths := make([]int, len(natural))
		copy(widths, natural)
		widths[widest] -= total - avail
		return widths
	}
	minW := 3
	if avail < minW*len(natural) {
		minW = 1 // very narrow pane: keep every column alive at width 1
	}
	widths := make([]int, len(natural))
	sum := 0
	for i, n := range natural {
		w := n * avail / total
		if floor := min(n, minW); w < floor {
			w = floor
		}
		widths[i] = w
		sum += w
	}
	// Rounding overflow: take from the widest columns still above the floor.
	for sum > avail {
		widest := -1
		for i, w := range widths {
			if w > minW && (widest < 0 || w > widths[widest]) {
				widest = i
			}
		}
		if widest < 0 {
			break // all at floor; Hardwrap backstops the impossible cases
		}
		widths[widest]--
		sum--
	}
	// Rounding shortfall: give back to the columns that lost the most.
	for sum < avail {
		neediest := -1
		for i, w := range widths {
			if w < natural[i] && (neediest < 0 || natural[i]-w > natural[neediest]-widths[neediest]) {
				neediest = i
			}
		}
		if neediest < 0 {
			break
		}
		widths[neediest]++
		sum++
	}
	return widths
}

// tableBorderRow draws one horizontal border in a single styled run.
func tableBorderRow(left, mid, right string, widths []int) string {
	var b strings.Builder
	b.WriteString(left)
	for i, w := range widths {
		if i > 0 {
			b.WriteString(mid)
		}
		b.WriteString(strings.Repeat("─", w+2))
	}
	b.WriteString(right)
	return mdTableBorder.Render(b.String())
}

// tableTextRows renders one logical table row as one or more visual rows:
// each cell word-wraps within its column, the row is as tall as its tallest
// cell, and shorter cells pad with blanks. Header cells render bold via a
// seeded ambient SGR, which mdInlineSeg's token re-assertion keeps alive
// across inline styles' trailing resets.
func tableTextRows(cells []string, widths []int, aligns []TableAlign, header bool) []string {
	wrapped := make([][]string, len(cells))
	height := 1
	for i, cell := range cells {
		wrapped[i] = strings.Split(ansi.Wrap(cell, widths[i], ""), "\n")
		if len(wrapped[i]) > height {
			height = len(wrapped[i])
		}
	}
	bar := mdTableBorder.Render("│")
	var out []string
	for line := 0; line < height; line++ {
		var b strings.Builder
		for i := range cells {
			b.WriteString(bar)
			b.WriteString(" ")
			text := ""
			if line < len(wrapped[i]) {
				text = wrapped[i][line]
			}
			b.WriteString(tableCell(text, widths[i], aligns[i], header))
			b.WriteString(" ")
		}
		b.WriteString(bar)
		out = append(out, b.String())
	}
	return out
}

// tableCell styles one wrapped cell line and pads it to w by alignment.
// Padding is computed from the raw text (styling never changes width).
func tableCell(text string, w int, align TableAlign, header bool) string {
	pad := w - lipgloss.Width(text)
	if pad < 0 {
		pad = 0
	}
	styled := ""
	if header {
		amb := "\x1b[1m"
		styled = amb + mdInlineSeg(text, &amb) + "\x1b[0m"
	} else {
		amb := ""
		styled = mdInlineSeg(text, &amb)
	}
	switch align {
	case AlignRight:
		return strings.Repeat(" ", pad) + styled
	case AlignCenter:
		left := pad / 2
		return strings.Repeat(" ", left) + styled + strings.Repeat(" ", pad-left)
	default:
		return styled + strings.Repeat(" ", pad)
	}
}
