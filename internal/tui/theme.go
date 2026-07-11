package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// A theme names every color role the TUI draws from. Values are either
// ANSI-256 codes ("39") or hex ("#7aa2f7") — the render profile degrades hex
// to the nearest 256-color, so both work in any terminal the app supports.
type theme struct {
	border, borderDim string // focused / blurred pane borders
	titleFg, titleBg  string // title bar
	statusFg          string // status bar + sidebar header text
	statusBg          string // status bar
	headerBg          string // editor header, sidebar header row, blurred selection
	hint              string // hints, prompt header, dimmed chrome
	errFg             string // error text
	selBg, selFg      string // selection bar (sidebar cursor, drag-selection)
	done, wip         string // progress glyphs: completed / in-progress
	text              string // chat body text
	badgeFg           string // text on the chat speaker badges
	user, tutor       string // badge backgrounds; tutor also tints busy/ok accents
	lesson, quiz      string // badge backgrounds
	inputBg           string // chat typing-area wash
	rule              string // input rule, blurred prompt
	codeBg            string // code-block wash in the transcript
	system            string // system lines in the transcript
	notice            string // status-bar notices
	fail              string // failed-check text
}

// themes are the palettes :theme can switch between. "meari" is the app's
// original look; the rest are the usual modern terminal suspects.
var themes = map[string]theme{
	"meari": {
		border: "39", borderDim: "240",
		titleFg: "231", titleBg: "25",
		statusFg: "250", statusBg: "236", headerBg: "238",
		hint: "244", errFg: "203",
		selBg: "24", selFg: "231",
		done: "42", wip: "214",
		text: "252", badgeFg: "232",
		user: "81", tutor: "79", lesson: "222", quiz: "215",
		inputBg: "237", rule: "240", codeBg: "236",
		system: "245", notice: "222", fail: "210",
	},
	"tokyo-night": {
		border: "#7aa2f7", borderDim: "#3b4261",
		titleFg: "#c0caf5", titleBg: "#3d59a1",
		statusFg: "#a9b1d6", statusBg: "#24283b", headerBg: "#292e42",
		hint: "#565f89", errFg: "#f7768e",
		selBg: "#3d59a1", selFg: "#c0caf5",
		done: "#9ece6a", wip: "#e0af68",
		text: "#c0caf5", badgeFg: "#1a1b26",
		user: "#7dcfff", tutor: "#73daca", lesson: "#e0af68", quiz: "#ff9e64",
		inputBg: "#292e42", rule: "#3b4261", codeBg: "#24283b",
		system: "#565f89", notice: "#e0af68", fail: "#f7768e",
	},
	"catppuccin": { // mocha
		border: "#89b4fa", borderDim: "#45475a",
		titleFg: "#cdd6f4", titleBg: "#45475a",
		statusFg: "#a6adc8", statusBg: "#313244", headerBg: "#313244",
		hint: "#6c7086", errFg: "#f38ba8",
		selBg: "#45475a", selFg: "#cdd6f4",
		done: "#a6e3a1", wip: "#f9e2af",
		text: "#cdd6f4", badgeFg: "#1e1e2e",
		user: "#89dceb", tutor: "#94e2d5", lesson: "#f9e2af", quiz: "#fab387",
		inputBg: "#313244", rule: "#45475a", codeBg: "#292c3c",
		system: "#6c7086", notice: "#f9e2af", fail: "#f38ba8",
	},
	"nord": {
		border: "#88c0d0", borderDim: "#4c566a",
		titleFg: "#eceff4", titleBg: "#5e81ac",
		statusFg: "#d8dee9", statusBg: "#3b4252", headerBg: "#434c5e",
		hint: "#616e88", errFg: "#bf616a",
		selBg: "#5e81ac", selFg: "#eceff4",
		done: "#a3be8c", wip: "#ebcb8b",
		text: "#d8dee9", badgeFg: "#2e3440",
		user: "#88c0d0", tutor: "#8fbcbb", lesson: "#ebcb8b", quiz: "#d08770",
		inputBg: "#3b4252", rule: "#4c566a", codeBg: "#3b4252",
		system: "#616e88", notice: "#ebcb8b", fail: "#bf616a",
	},
	"dracula": {
		border: "#bd93f9", borderDim: "#44475a",
		titleFg: "#f8f8f2", titleBg: "#44475a",
		statusFg: "#f8f8f2", statusBg: "#343746", headerBg: "#44475a",
		hint: "#6272a4", errFg: "#ff5555",
		selBg: "#44475a", selFg: "#f8f8f2",
		done: "#50fa7b", wip: "#ffb86c",
		text: "#f8f8f2", badgeFg: "#282a36",
		user: "#8be9fd", tutor: "#50fa7b", lesson: "#f1fa8c", quiz: "#ffb86c",
		inputBg: "#343746", rule: "#44475a", codeBg: "#343746",
		system: "#6272a4", notice: "#f1fa8c", fail: "#ff5555",
	},
	"gruvbox": { // dark
		border: "#83a598", borderDim: "#504945",
		titleFg: "#ebdbb2", titleBg: "#504945",
		statusFg: "#d5c4a1", statusBg: "#3c3836", headerBg: "#3c3836",
		hint: "#928374", errFg: "#fb4934",
		selBg: "#504945", selFg: "#ebdbb2",
		done: "#b8bb26", wip: "#fabd2f",
		text: "#ebdbb2", badgeFg: "#282828",
		user: "#83a598", tutor: "#8ec07c", lesson: "#fabd2f", quiz: "#fe8019",
		inputBg: "#3c3836", rule: "#504945", codeBg: "#3c3836",
		system: "#928374", notice: "#fabd2f", fail: "#fb4934",
	},
	"rose-pine": {
		border: "#c4a7e7", borderDim: "#403d52",
		titleFg: "#e0def4", titleBg: "#403d52",
		statusFg: "#908caa", statusBg: "#1f1d2e", headerBg: "#26233a",
		hint: "#6e6a86", errFg: "#eb6f92",
		selBg: "#403d52", selFg: "#e0def4",
		done: "#9ccfd8", wip: "#f6c177",
		text: "#e0def4", badgeFg: "#191724",
		user: "#9ccfd8", tutor: "#ebbcba", lesson: "#f6c177", quiz: "#c4a7e7",
		inputBg: "#26233a", rule: "#403d52", codeBg: "#1f1d2e",
		system: "#6e6a86", notice: "#f6c177", fail: "#eb6f92",
	},
}

const defaultThemeName = "meari"

// currentThemeName is what :theme reports and the theme file persists.
var currentThemeName = defaultThemeName

// themeNames returns the selectable theme names, sorted.
func themeNames() []string {
	names := make([]string, 0, len(themes))
	for n := range themes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// bgSeq returns the raw SGR sequence that opens color as a background, in the
// app's ANSI-256 profile. chat.go re-asserts this after every reset code to
// keep the input wash solid (see inputView).
func bgSeq(color string) string {
	c := termenv.ANSI256.Color(color)
	if c == nil {
		return ""
	}
	return "\x1b[" + c.Sequence(true) + "m"
}

// applyTheme rebuilds every package-level style from the named palette. All
// render paths read these vars each frame, so the switch is instant. Returns
// false for an unknown name.
func applyTheme(name string) bool {
	t, ok := themes[name]
	if !ok {
		return false
	}
	c := func(s string) lipgloss.Color { return lipgloss.Color(s) }

	focusedBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(c(t.border))
	blurredBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(c(t.borderDim))
	titleBar = lipgloss.NewStyle().Bold(true).Foreground(c(t.titleFg)).Background(c(t.titleBg)).Padding(0, 1)
	statusBar = lipgloss.NewStyle().Foreground(c(t.statusFg)).Background(c(t.statusBg)).Padding(0, 1)
	editorHeader = lipgloss.NewStyle().Foreground(c(t.titleFg)).Background(c(t.headerBg)).Padding(0, 1)
	hintStyle = lipgloss.NewStyle().Foreground(c(t.hint))
	errStyle = lipgloss.NewStyle().Foreground(c(t.errFg)).Bold(true)
	checkButton = lipgloss.NewStyle().Bold(true).Foreground(c(t.badgeFg)).Background(c(t.tutor))

	selectedBg = c(t.selBg)
	selectedBlurredBg = c(t.headerBg)
	selectedFg = c(t.selFg)
	doneColor = c(t.done)
	wipColor = c(t.wip)
	selectedRow = lipgloss.NewStyle().Foreground(selectedFg).Background(selectedBg)
	headerRow = lipgloss.NewStyle().Foreground(c(t.statusFg)).Background(c(t.headerBg)).Bold(true)
	doneGlyph = lipgloss.NewStyle().Foreground(doneColor)
	wipGlyph = lipgloss.NewStyle().Foreground(wipColor)
	markedItem = lipgloss.NewStyle().Foreground(wipColor)

	chatBodyStyle = lipgloss.NewStyle().Foreground(c(t.text))
	chatUserBadge = lipgloss.NewStyle().Bold(true).Foreground(c(t.badgeFg)).Background(c(t.user))
	chatTutorBadge = lipgloss.NewStyle().Bold(true).Foreground(c(t.badgeFg)).Background(c(t.tutor))
	chatLessonBadge = lipgloss.NewStyle().Bold(true).Foreground(c(t.badgeFg)).Background(c(t.lesson))
	chatQuizBadge = lipgloss.NewStyle().Bold(true).Foreground(c(t.badgeFg)).Background(c(t.quiz))
	chatBusyStyle = lipgloss.NewStyle().Foreground(c(t.tutor)).Italic(true)

	chatInputBG = c(t.inputBg)
	chatInputBGSeq = bgSeq(t.inputBg)
	chatInputRule = lipgloss.NewStyle().Foreground(c(t.rule))
	chatPromptFocus = lipgloss.NewStyle().Foreground(c(t.user)).Bold(true)
	chatPromptBlur = lipgloss.NewStyle().Foreground(c(t.rule))
	chatPromptNormal = lipgloss.NewStyle().Foreground(c(t.done)).Bold(true)
	chatCodeGutter = lipgloss.NewStyle().Foreground(c(t.user)).Bold(true)
	chatCodeLine = lipgloss.NewStyle().Background(c(t.codeBg))
	chatSystemStyle = lipgloss.NewStyle().Foreground(c(t.system)).Italic(true)
	chatSelStyle = lipgloss.NewStyle().Foreground(selectedFg).Background(selectedBg)

	noticeStyle = lipgloss.NewStyle().Foreground(c(t.notice))
	promptHeaderStyle = lipgloss.NewStyle().Foreground(c(t.hint)).Italic(true)
	backlinkHeaderStyle = lipgloss.NewStyle().Foreground(c(t.tutor)).Bold(true)
	chatOkStyle = lipgloss.NewStyle().Foreground(doneColor).Bold(true)
	chatFailStyle = lipgloss.NewStyle().Foreground(c(t.fail)).Bold(true)

	currentThemeName = name
	return true
}

// themeFile is where the chosen theme persists — in the data dir, next to
// progress.json, never in the user's config.toml.
func themeFile(dataDir string) string { return filepath.Join(dataDir, "theme") }

// loadTheme applies the persisted theme, if any. Unknown or missing names
// keep the default palette (the styles.go initializers).
func loadTheme(dataDir string) {
	if dataDir == "" {
		return
	}
	b, err := os.ReadFile(themeFile(dataDir))
	if err != nil {
		return
	}
	if name := strings.TrimSpace(string(b)); name != "" {
		applyTheme(name)
	}
}

func saveTheme(dataDir, name string) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(themeFile(dataDir), []byte(name+"\n"), 0o644)
}

// themeCommand implements ":theme [<name>]" for both TUIs: no argument lists
// the palettes, a name switches live and persists the choice. Returns the
// status-bar message.
func themeCommand(dataDir, arg string) string {
	if strings.TrimSpace(arg) == "" {
		return "themes: " + strings.Join(themeNames(), " · ") + "  — current: " + currentThemeName
	}
	name := strings.ToLower(strings.TrimSpace(arg))
	if !applyTheme(name) {
		return "unknown theme " + strconv.Quote(name) + " — :theme lists the choices"
	}
	if err := saveTheme(dataDir, name); err != nil {
		return "theme " + name + " (couldn't persist: " + err.Error() + ")"
	}
	return "theme " + name
}

// themeArgCandidates completes ":theme <Tab>" against the palette names.
func themeArgCandidates(input string) []string {
	if !strings.HasPrefix(input, "theme ") {
		return nil
	}
	names := themeNames()
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = "theme " + n
	}
	return out
}
