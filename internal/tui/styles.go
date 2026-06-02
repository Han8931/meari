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

	hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	// Sidebar row styles.
	selectedRow = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("24"))
	headerRow   = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238")).Bold(true)
	doneGlyph   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	wipGlyph    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	// Chat transcript role styles.
	chatUserStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true)
	chatTutorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	chatSystemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true)
	chatOkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	chatFailStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
)

// borderStyle returns the focused or blurred border depending on whether the
// pane is the active one.
func borderStyle(active bool) lipgloss.Style {
	if active {
		return focusedBorder
	}
	return blurredBorder
}
