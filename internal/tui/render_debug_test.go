package tui

import (
	"fmt"
	"regexp"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/tutor"
)

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*[A-Za-z]")

func TestRenderFrameInspection(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 26})
	m.phase = phaseReady
	_ = m.loadChallenge(tutor.Challenge{
		ID:   "x",
		Lang: "go",
		Prompt: "Write SumTo(n int) int returning 1 + 2 + ... + n computed with a for " +
			"loop (return 0 when n < 1). This is a deliberately long prompt to force wrapping.",
		StarterCode: "func SumTo(n int) int {\n\treturn 0\n}\n",
	})
	view := ansiRE.ReplaceAllString(m.View(), "")
	fmt.Println("=== FRAME 100x26 ===")
	fmt.Println(view)
	fmt.Println("=== line count:", len(splitLines(view)), "===")
}

func splitLines(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == '\n' {
			out = append(out, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	return append(out, cur)
}

func TestRenderCurriculumAndHorizontal(t *testing.T) {
	// Curriculum topic at a narrow 80x24.
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m.phase = phaseReady
	_ = m.loadCurriculum("go", "beginner", "")
	view := ansiRE.ReplaceAllString(m.View(), "")
	fmt.Println("=== CURRICULUM 80x24 ===")
	fmt.Println(view)
	fmt.Println("=== line count:", len(splitLines(view)), "===")

	// Horizontal layout.
	m2 := newModel(testDeps(t))
	m2.horizontal = true
	m2 = step(t, m2, tea.WindowSizeMsg{Width: 80, Height: 24})
	m2.phase = phaseReady
	_ = m2.loadCurriculum("go", "beginner", "")
	view2 := ansiRE.ReplaceAllString(m2.View(), "")
	fmt.Println("=== HORIZONTAL 80x24 ===")
	fmt.Println(view2)
	fmt.Println("=== line count:", len(splitLines(view2)), "===")
}

func TestRenderWithoutLineNumbers(t *testing.T) {
	d := testDeps(t)
	off := false
	d.Cfg.Editor.LineNumbers = &off
	m := newModel(d)
	m = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 18})
	m.phase = phaseReady
	_ = m.loadCurriculum("go", "beginner", "")
	view := ansiRE.ReplaceAllString(m.View(), "")
	fmt.Println("=== NO LINE NUMBERS 80x18 ===")
	fmt.Println(view)
}
