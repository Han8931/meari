package editor

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func TestIsTableSeparator(t *testing.T) {
	for _, valid := range []string{
		"| --- | --- |", "|---|---|", "| :-- | :-: | --: |", "---|---", "| - |",
	} {
		if !IsTableSeparator(valid) {
			t.Errorf("IsTableSeparator(%q) = false, want true", valid)
		}
	}
	for _, invalid := range []string{
		"---", "| a | b |", "| -- | text |", "", "| : |", "plain prose",
	} {
		if IsTableSeparator(invalid) {
			t.Errorf("IsTableSeparator(%q) = true, want false", invalid)
		}
	}
}

func TestRenderTableBasicGrid(t *testing.T) {
	withANSI(t)
	rows, ok := RenderTable([]string{
		"| Type | Owns? |",
		"| --- | --- |",
		"| String | yes |",
		"| &str | no |",
	}, 60)
	if !ok {
		t.Fatal("table did not parse")
	}
	joined := strings.Join(rows, "\n")
	plain := ansi.Strip(joined)
	for _, want := range []string{"┌", "┬", "┐", "├", "┼", "┤", "└", "┴", "┘", "│"} {
		if !strings.Contains(plain, want) {
			t.Errorf("grid missing %q:\n%s", want, plain)
		}
	}
	// Header is bold; borders carry the dim border color.
	if !strings.Contains(joined, "\x1b[1mType") {
		t.Errorf("header not bold:\n%q", rows[1])
	}
	if !strings.Contains(joined, "\x1b[38;5;240m") {
		t.Errorf("borders not tinted:\n%q", rows[0])
	}
	// Body cells present, raw pipes and separator dashes gone.
	if !strings.Contains(plain, "String") || !strings.Contains(plain, "&str") {
		t.Errorf("body cells missing:\n%s", plain)
	}
	if strings.Contains(plain, "---") {
		t.Errorf("separator row leaked into output:\n%s", plain)
	}
	// Every │ column aligns across rows.
	var pipeCols []int
	for i, r := range strings.Split(plain, "\n") {
		if !strings.Contains(r, "│") {
			continue
		}
		var cols []int
		for j, ch := range r {
			if ch == '│' {
				cols = append(cols, j)
			}
		}
		if pipeCols == nil {
			pipeCols = cols
		} else if len(cols) != len(pipeCols) {
			t.Errorf("row %d has %d pipes, want %d: %q", i, len(cols), len(pipeCols), r)
		}
	}
}

func TestRenderTableAlignment(t *testing.T) {
	rows, ok := RenderTable([]string{
		"| L | C | R |",
		"| :--- | :---: | ---: |",
		"| a | b | c |",
	}, 40)
	if !ok {
		t.Fatal("table did not parse")
	}
	plain := ansi.Strip(strings.Join(rows, "\n"))
	body := strings.Split(plain, "\n")[3] // top, header, sep, body
	cells := strings.Split(strings.Trim(body, "│"), "│")
	if len(cells) != 3 {
		t.Fatalf("want 3 cells, got %d: %q", len(cells), body)
	}
	if !strings.HasPrefix(cells[0], " a") {
		t.Errorf("left cell not left-aligned: %q", cells[0])
	}
	if !strings.HasSuffix(cells[2], "c ") {
		t.Errorf("right cell not right-aligned: %q", cells[2])
	}
	l, r := strings.Index(cells[1], "b"), len(cells[1])-strings.Index(cells[1], "b")-1
	if l < 1 || r < 1 || l-r > 1 || r-l > 1 {
		t.Errorf("center cell not centered (%d left, %d right): %q", l, r, cells[1])
	}
}

// A table whose natural width exceeds the pane must shrink columns and wrap
// cell text within them — never exceed the width, never drop a column.
func TestRenderTableShrinksAndWraps(t *testing.T) {
	rows, ok := RenderTable([]string{
		"| Approach | Dispatch | Use when |",
		"| --- | --- | --- |",
		"| generics and impl Trait everywhere | static at compile time, zero cost | the concrete type is known at compile time |",
	}, 40)
	if !ok {
		t.Fatal("table did not parse")
	}
	for i, r := range rows {
		if w := lipgloss.Width(r); w > 40 {
			t.Errorf("row %d width %d > 40: %q", i, w, ansi.Strip(r))
		}
	}
	plain := ansi.Strip(strings.Join(rows, "\n"))
	// The long body must wrap across multiple visual rows, all content kept.
	if len(rows) < 7 {
		t.Errorf("expected wrapped body rows, got %d rows:\n%s", len(rows), plain)
	}
	for _, word := range []string{"generics", "static", "concrete"} {
		if !strings.Contains(plain, word) {
			t.Errorf("cell content %q lost in wrap:\n%s", word, plain)
		}
	}
	// Pipes still align on every text row.
	var pipeCount int
	for _, r := range strings.Split(plain, "\n") {
		if strings.HasPrefix(r, "│") {
			if n := strings.Count(r, "│"); pipeCount == 0 {
				pipeCount = n
			} else if n != pipeCount {
				t.Errorf("ragged pipes (%d vs %d): %q", n, pipeCount, r)
			}
		}
	}
}

func TestRenderTableEdgeCases(t *testing.T) {
	// Escaped pipe becomes a literal; inline code with a pipe splits per GFM.
	rows, ok := RenderTable([]string{
		"| Symbol | Meaning |",
		"| --- | --- |",
		`| \| | or |`,
		"| a | b |",
	}, 40)
	if !ok {
		t.Fatal("escaped-pipe table did not parse")
	}
	if plain := ansi.Strip(strings.Join(rows, "\n")); !strings.Contains(plain, "│ |") {
		t.Errorf("escaped pipe not rendered as literal:\n%s", plain)
	}

	// Single column.
	if rows, ok = RenderTable([]string{"| Only |", "| --- |", "| one |"}, 30); !ok {
		t.Fatal("single-column table did not parse")
	} else if plain := ansi.Strip(strings.Join(rows, "\n")); !strings.Contains(plain, "one") {
		t.Errorf("single-column cell missing:\n%s", plain)
	}

	// Ragged rows: missing cells go empty, excess cells drop.
	rows, ok = RenderTable([]string{
		"| A | B | C |",
		"| --- | --- | --- |",
		"| 1 |",
		"| 1 | 2 | 3 | 4 |",
	}, 40)
	if !ok {
		t.Fatal("ragged table did not parse")
	}
	if plain := ansi.Strip(strings.Join(rows, "\n")); strings.Contains(plain, "4") {
		t.Errorf("excess cell not dropped:\n%s", plain)
	}

	// Empty cells and a tiny pane must not panic; header-only table renders.
	if _, ok := RenderTable([]string{"| a |  | c |", "| --- | --- | --- |"}, 4); !ok {
		t.Error("header-only/empty-cell table at width 4 did not render")
	}

	// Non-tables refuse.
	if _, ok := RenderTable([]string{"| a | b |", "| not | sep |"}, 40); ok {
		t.Error("non-separator second line accepted as table")
	}
	if _, ok := RenderTable([]string{"| only one line |"}, 40); ok {
		t.Error("single line accepted as table")
	}
}

// A slight overflow shrinks only the widest column, so narrow headers like
// "Growable?" never wrap over a few cells of debt.
func TestRenderTableSmallOverflowHitsWidestColumn(t *testing.T) {
	rows, ok := RenderTable([]string{
		"| Type   | What it is                      | Growable? |",
		"| ------ | ------------------------------- | --------- |",
		"| String | an owned, heap-allocated string | yes       |",
	}, 50) // a hair narrower than natural
	if !ok {
		t.Fatal("table did not parse")
	}
	plain := ansi.Strip(strings.Join(rows, "\n"))
	if !strings.Contains(plain, "│ Growable? │") {
		t.Errorf("narrow header wrapped despite small overflow:\n%s", plain)
	}
}

// Inline markdown inside cells keeps its styling.
func TestRenderTableInlineStyling(t *testing.T) {
	withANSI(t)
	rows, ok := RenderTable([]string{
		"| Feature | Note |",
		"| --- | --- |",
		"| `code` | **bold** and [[Link]] |",
	}, 60)
	if !ok {
		t.Fatal("table did not parse")
	}
	joined := strings.Join(rows, "\n")
	if !strings.Contains(joined, "\x1b[38;5;222m`code`") {
		t.Errorf("inline code not styled in cell:\n%q", joined)
	}
	if !strings.Contains(joined, "\x1b[1;38;5;222m**bold**") {
		t.Errorf("bold not styled in cell:\n%q", joined)
	}
	if !strings.Contains(joined, "\x1b[38;5;79m[[Link]]") {
		t.Errorf("wikilink not styled in cell:\n%q", joined)
	}
}
