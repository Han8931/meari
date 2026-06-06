package tui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// noticeTTL is how long a status-bar notice stays visible.
const noticeTTL = 5 * time.Second

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

	// Chat transcript styles. Each speaker turn opens with a colored BADGE on
	// its own line (" you " / " tutor " / " lesson " on a filled background) so
	// who is talking is unmistakable even when skimming; the message body stays
	// a calm, high-contrast neutral so long passages read comfortably.
	chatBodyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	chatUserBadge   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("81"))  // learner — cyan-blue
	chatTutorBadge  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("79"))  // tutor — teal
	chatLessonBadge = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("222")) // lesson — warm gold

	// chatBusyStyle renders the in-pane "⠹ tutor thinking…" progress line.
	chatBusyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("79")).Italic(true)

	// Input-area styles. A dim rule (chatInputRule) separates the transcript
	// from the typing area, and the "┃ " prompt prefix is painted bright cyan
	// when the pane is focused (matching the "you" badge) and dim when it isn't,
	// so it's always obvious where typing happens.
	chatInputRule   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	chatPromptFocus = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	chatPromptBlur  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// chatCodeGutter is the "│ " bar marking syntax-highlighted code blocks in
	// tutor/lesson messages.
	chatCodeGutter = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	chatSystemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)

	// noticeStyle renders transient command feedback in the status bar (copy
	// confirmations, resize/fold notices, unknown commands…).
	noticeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("222"))

	// promptHeaderStyle renders the pinned challenge/essay statement above the
	// editor, styled like a code comment.
	promptHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true)
	chatOkStyle       = lipgloss.NewStyle().Foreground(doneColor).Bold(true)             // match the sidebar "done" green
	chatFailStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("210")).Bold(true) // soft red, easier on the eyes than 203
)

// borderStyle returns the focused or blurred border depending on whether the
// pane is the active one.
func borderStyle(active bool) lipgloss.Style {
	if active {
		return focusedBorder
	}
	return blurredBorder
}
