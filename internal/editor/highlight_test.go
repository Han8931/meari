package editor

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// withANSI forces a color profile so styles emit real escape sequences (tests
// run without a TTY, where lipgloss otherwise renders plain text).
func withANSI(t *testing.T) {
	t.Helper()
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.ANSI256)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })
}

// Markdown (and any non-code language) must pass through the highlighter
// untouched: code rules over prose bold ordinary words ("for", "and"…) and
// their resets punch holes in the textarea's cursor-line background, which
// reads as leftover text following the cursor.
func TestHighlightProsePassesThrough(t *testing.T) {
	withANSI(t)
	for _, lang := range []string{"markdown", "md", "plain", "text", "physics"} {
		s := "for tips and notes import these words"
		if got := highlightSyntax(lang, s); got != s {
			t.Errorf("highlightSyntax(%q) altered prose: %q", lang, got)
		}
	}
}

// On the focused cursor row the textarea opens a background (SGR 40) that must
// survive every highlighted token: each token's trailing reset has to be
// followed by a re-assertion of the ambient state.
func TestHighlightReassertsAmbientSGR(t *testing.T) {
	withANSI(t)
	m := New("", true, nil)
	m.SetLanguage("python")
	m.SetSize(60, 8)
	m.SetValue("for x in items:\n    pass\n")
	m.Focus()
	row := strings.Split(m.View(), "\n")[0]

	// The keyword "in" sits mid-row, clear of the cursor: it must be styled,
	// and the cursor-line background must resume right after its reset — the
	// regression was "\x1b[0m x items:" leaving the rest of the bar unpainted.
	if !strings.Contains(row, "\x1b[1;38;5;81min\x1b[0m\x1b[40m") {
		t.Fatalf("keyword styling or bg re-assert missing in row: %q", row)
	}
	// The non-keyword text between tokens keeps the background too.
	if !strings.Contains(row, "\x1b[40m x ") && !strings.Contains(row, "\x1b[40m items:") {
		t.Fatalf("plain text between tokens lost the cursor-line bg: %q", row)
	}
}
