package tui

import "github.com/charmbracelet/lipgloss"

// Pane border styles. The focused pane gets a bright border; others are dim.
// We switch only BorderForeground (not the border runes) so box widths stay
// stable across focus changes and the layout never jumps.
var (
	focusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39"))

	blurredBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	titleBar = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	editorHeader = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("238")).
			Padding(0, 1)

	hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	// Sidebar row styles and colors. selectedBg paints the cursor bar when the
	// pane is focused; selectedBlurredBg is the dimmed bar shown when the pane has
	// lost focus (the way ranger/lf fade an inactive pane's selection). doneColor
	// and wipColor are reused both for the standalone glyph styles and for the
	// glyphs painted onto the selection bar.
	selectedBg        = lipgloss.Color("24")
	selectedBlurredBg = lipgloss.Color("238")
	selectedFg        = lipgloss.Color("231")
	doneColor         = lipgloss.Color("42")
	wipColor          = lipgloss.Color("214")

	selectedRow = lipgloss.NewStyle().Foreground(selectedFg).Background(selectedBg)
	headerRow   = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238")).Bold(true)
	doneGlyph   = lipgloss.NewStyle().Foreground(doneColor)
	wipGlyph    = lipgloss.NewStyle().Foreground(wipColor)

	// Chat transcript styles. Speaker labels are short and saturated so the eye
	// can find who's talking; the message body stays a calm, high-contrast
	// neutral so long passages read comfortably (the old scheme tinted whole
	// paragraphs in a low-contrast role color, which was hard to read).
	chatBodyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	chatUserLabel   = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)  // learner — cyan-blue
	chatTutorLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("79")).Bold(true)  // tutor — teal
	chatLessonLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("222")).Bold(true) // lesson — warm gold

	chatSystemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
	chatOkStyle     = lipgloss.NewStyle().Foreground(doneColor).Bold(true)             // match the sidebar "done" green
	chatFailStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("210")).Bold(true) // soft red, easier on the eyes than 203
)

// borderStyle returns the focused or blurred border depending on whether the
// pane is the active one.
func borderStyle(active bool) lipgloss.Style {
	if active {
		return focusedBorder
	}
	return blurredBorder
}
