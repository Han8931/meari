// Package editor provides an in-app code editor with a configurable, hand-rolled
// Vim mode built on Bubble Tea + the bubbles textarea component.
//
// Vim is modeled as a small state machine over the textarea:
//
//	Normal  - movement (h/j/k/l), enter Insert (i/a/o), delete char (x),
//	          open the command line (:)
//	Insert  - normal typing; Esc returns to Normal
//	Command - the ":" line; runs submit / w / q (see the design discussion)
//
// Cursor motions in Normal mode are implemented by translating h/j/k/l into the
// arrow-key messages the textarea already understands, so we reuse its cursor
// logic instead of reimplementing it.
//
// When config selects "default" keybindings, the modal layer is bypassed and
// keys are forwarded straight to the textarea.
package editor

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Action is why the editor closed.
type Action string

const (
	ActionSubmit Action = "submit" // ":submit" — check the solution
	ActionSave   Action = "save"   // editor closed after a save (rare; :w stays open)
	ActionQuit   Action = "quit"   // ":q" — leave without checking
)

type mode int

const (
	modeNormal mode = iota
	modeInsert
	modeCommand
	modeVisual
)

// Result is returned to the caller after the editor closes.
type Result struct {
	Action Action
	Code   string
}

// DoneMsg is emitted when the editor reaches a terminal action (submit/quit).
// A standalone editor (Open) turns this into tea.Quit; an embedding parent
// intercepts it to react (e.g. run tests on submit) without killing its program.
type DoneMsg struct{ Action Action }

// OpenConfigMsg is emitted by the ":config" command. The embedding parent
// handles it (e.g. launches the user's $EDITOR on the config file); the editor
// itself stays open in Normal mode.
type OpenConfigMsg struct{}

// RunCommandMsg forwards an ex-command the editor doesn't own (e.g. ":topic go",
// ":progress", ":clear") up to the embedding parent, which dispatches it. Raw is
// the command without the leading colon. This lets the editor's command line
// double as the app's global command line instead of having two of them.
type RunCommandMsg struct{ Raw string }

// SaveFunc persists an in-progress draft. Called on ":w" while editing.
type SaveFunc func(code string) error

// Model is the editor component. It can run on its own via Open or be embedded
// in a parent Bubble Tea model via New + delegated Update/View.
type Model struct {
	vim    bool
	mode   mode
	ta     textarea.Model
	cmd    textinput.Model
	save   SaveFunc
	status string
	lang   string
	width  int
	height int

	// standalone is set by Open; it makes DoneMsg quit the program. When
	// embedded (New), DoneMsg is left for the parent to handle.
	standalone bool

	// pending holds the first key of a two-key Normal-mode command (operator or
	// prefix): 'd' (delete), 'c' (change), 'g' (goto), 'r' (replace),
	// 'y' (yank), '<'/'>' (dedent/indent). 0 = none.
	pending rune

	// register is Vim's unnamed register: the text captured by the latest
	// delete/yank, pasted by p/P. regLinewise marks whole-line content (dd/yy),
	// which p pastes on its own line instead of inline.
	register    string
	regLinewise bool

	// Visual-mode selection: the anchor is the position where v/V was pressed;
	// the selection spans anchor..cursor. visualLine marks linewise mode (V).
	anchorRow, anchorCol int
	visualLine           bool

	// undo/redo snapshots ('u' and Ctrl-R in Normal mode). One snapshot is
	// taken before each mutating command, and once on Insert entry so a whole
	// typing session is a single undo unit.
	undoStack []editState
	redoStack []editState

	// count is the numeric prefix being typed (3w, 2dd, 5x); 0 = none.
	// pendingCount preserves it across a two-key command (2dd arms 'd' with 2).
	count        int
	pendingCount int

	// lastFindOp/Ch remember the latest f/F/t/T so ; and , can repeat it.
	lastFindOp rune
	lastFindCh rune

	// searchMode marks the command line as a / search; lastSearch feeds n/N.
	searchMode bool
	lastSearch string

	action Action
	done   bool
}

// New builds an embeddable editor pre-filled with starter. vim selects Vim
// bindings; save is invoked on ":w". The caller owns the tea.Program and must
// delegate Update/View and drive sizing via SetSize.
func New(starter string, vim bool, save SaveFunc) Model {
	ta := textarea.New()
	ta.SetValue(starter)
	ta.ShowLineNumbers = true
	ta.Prompt = ""
	ta.CharLimit = 0
	ta.Focus()
	// A steady (non-blinking) cursor: when focused it always renders, so it can
	// never "disappear" mid-blink — the per-mode color (set in View) is always
	// visible.
	ta.Cursor.SetMode(cursor.CursorStatic)
	// SetValue leaves the cursor at the end; start at the top, like opening a
	// file. The textarea ignores keys while blurred, so do this after Focus.
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyCtrlHome})

	ci := textinput.New()
	ci.Prompt = ":"

	m := Model{
		vim:  vim,
		mode: modeNormal, // Vim opens in Normal; default mode treats this as "typing"
		ta:   ta,
		cmd:  ci,
		save: save,
		lang: "python",
	}
	if !vim {
		m.mode = modeInsert
	}
	return m
}

// SetSize lays out the editor within w×h, reserving one row for the status line.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.ta.SetWidth(w)
	if h > 1 {
		m.ta.SetHeight(h - 1)
	}
	m.cmd.Width = w
}

// SetValue replaces the editor's contents (e.g. when switching challenges) and
// moves the cursor to the top, as if opening the file fresh. The undo history
// belongs to the previous buffer, so it is dropped — undo must never resurrect
// another note/challenge's text. The textarea ignores keys while blurred, so
// focus is momentarily restored for the move.
func (m *Model) SetValue(s string) {
	m.ta.SetValue(s)
	m.clearHistory()
	wasFocused := m.ta.Focused()
	m.ta.Focus()
	m.send(tea.KeyCtrlHome)
	if !wasFocused {
		m.ta.Blur()
	}
}

// Value returns the current buffer contents.
func (m Model) Value() string { return m.ta.Value() }

// SetLanguage selects the syntax highlighter used when rendering the buffer.
// Empty keeps the historical default: Python.
func (m *Model) SetLanguage(lang string) {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		lang = "python"
	}
	m.lang = lang
}

// Focus gives the editor's textarea keyboard focus; the returned cmd blinks the cursor.
func (m *Model) Focus() tea.Cmd { return m.ta.Focus() }

// Blur removes keyboard focus from the editor's textarea.
func (m *Model) Blur() { m.ta.Blur() }

var (
	hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	// Mode badges in the status line, color-coded so the current mode is
	// obvious at a glance (green = Normal, blue = Insert, amber = plain Edit).
	normalBadge = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("42")).Padding(0, 1)
	insertBadge = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("39")).Padding(0, 1)
	editBadge   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("214")).Padding(0, 1)

	// Cursor colors per mode. The cursor renders as a reverse-video block, so a
	// foreground color here becomes the block's fill: green in Normal, bright
	// magenta in Insert — both always visible (the cursor is static, not blinking).
	normalCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	insertCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))

	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	typeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	stringStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("106"))
	commentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true)
	numberStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("176"))
)

// Open launches the editor as its own full-screen program, pre-filled with
// starter, returning the final Result. vim selects Vim bindings; save is
// invoked on ":w".
func Open(starter string, vim bool, save SaveFunc) (Result, error) {
	m := New(starter, vim, save)
	m.standalone = true

	p := tea.NewProgram(m, tea.WithAltScreen())
	out, err := p.Run()
	if err != nil {
		return Result{}, err
	}
	fm := out.(Model)
	return Result{Action: fm.action, Code: fm.ta.Value()}, nil
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Only honored standalone; an embedding parent drives layout via SetSize.
		if m.standalone {
			m.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	case DoneMsg:
		// Standalone: terminate the program. Embedded: the parent handles it.
		if m.standalone {
			return m, tea.Quit
		}
		return m, nil

	case tea.KeyMsg:
		// Non-Vim mode: forward everything to the textarea; Ctrl-S submits,
		// Ctrl-Q quits. This is the "default" keybindings path.
		if !m.vim {
			switch msg.String() {
			case "ctrl+s":
				return m.finish(ActionSubmit)
			case "ctrl+q":
				return m.finish(ActionQuit)
			case "tab":
				// KeyTab is not a rune key, so the textarea ignores it; indent here.
				m.ta.InsertString(tabIndent)
				return m, nil
			}
			var cmd tea.Cmd
			m.ta, cmd = m.ta.Update(msg)
			return m, cmd
		}
		return m.updateVim(msg)
	}

	var cmd tea.Cmd
	m.ta, cmd = m.ta.Update(msg)
	return m, cmd
}

// finish records the terminal action and emits a DoneMsg for the program (or
// parent) to act on. It does not call tea.Quit directly, so an embedded editor
// never tears down its parent program.
func (m Model) finish(a Action) (tea.Model, tea.Cmd) {
	m.action = a
	return m, func() tea.Msg { return DoneMsg{Action: a} }
}

// updateVim routes a key according to the current Vim mode.
func (m Model) updateVim(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Ctrl-S / Ctrl-Q always work as escape hatches regardless of mode.
	switch msg.String() {
	case "ctrl+s":
		return m.finish(ActionSubmit)
	case "ctrl+q":
		return m.finish(ActionQuit)
	}

	switch m.mode {
	case modeNormal:
		return m.updateNormal(msg)
	case modeVisual:
		return m.updateVisual(msg)
	case modeInsert:
		switch msg.Type {
		case tea.KeyEsc:
			m.mode = modeNormal
			return m, nil
		case tea.KeyTab:
			// KeyTab is not a rune key, so the textarea ignores it; indent here.
			m.ta.InsertString(tabIndent)
			return m, nil
		}
		var cmd tea.Cmd
		m.ta, cmd = m.ta.Update(msg)
		return m, cmd
	case modeCommand:
		return m.updateCommand(msg)
	}
	return m, nil
}

// arrow sends a synthetic arrow key to the textarea so we reuse its cursor math.
func (m *Model) arrow(t tea.KeyType) {
	m.ta, _ = m.ta.Update(tea.KeyMsg{Type: t})
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.status = ""

	// A prefix/operator is pending (d, c, g, r, f, …): this key completes it.
	if m.pending != 0 {
		op := m.pending
		m.pending = 0
		return m.runPending(op, msg)
	}

	key := msg.String()

	// Numeric prefix: digits accumulate a count for the next command (3w, 2dd).
	// "0" is the line-start motion unless a count is already in progress.
	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' || (key == "0" && m.count > 0) {
		m.count = m.count*10 + int(key[0]-'0')
		return m, nil
	}
	if msg.Type == tea.KeyEsc {
		m.count = 0
		return m, nil
	}

	// Motions are shared with Visual mode; a count repeats them.
	if m.applyMotion(key) {
		for n := m.takeCount() - 1; n > 0; n-- {
			m.applyMotion(key)
		}
		return m, nil
	}

	switch key {
	// --- undo / redo ---
	case "u":
		m.undo()
	case "ctrl+r":
		m.redo()

	// --- char find & repeat ---
	case "f", "F", "t", "T":
		m.pendingCount = m.takeCount()
		m.pending = []rune(key)[0]
		return m, nil
	case ";":
		m.repeatFind(false, m.takeCount())
	case ",":
		m.repeatFind(true, m.takeCount())

	// --- line edits ---
	case "J":
		m.pushUndo()
		m.joinLines(m.takeCount())
	case "~":
		m.pushUndo()
		m.toggleCase(m.takeCount())

	// --- search ---
	case "/":
		m.searchMode = true
		m.mode = modeCommand
		m.cmd.Prompt = "/"
		m.cmd.SetValue("")
		m.cmd.Focus()
		return m, textinput.Blink
	case "n":
		if !m.search(m.lastSearch, 1) {
			m.status = "pattern not found: " + m.lastSearch
		}
	case "N":
		if !m.search(m.lastSearch, -1) {
			m.status = "pattern not found: " + m.lastSearch
		}

	// --- enter Visual ---
	case "v":
		m.enterVisual(false)
	case "V":
		m.enterVisual(true)

	// --- enter Insert (one undo unit per session) ---
	case "i":
		m.pushUndo()
		m.mode = modeInsert
	case "I":
		m.pushUndo()
		m.ta.CursorStart()
		m.mode = modeInsert
	case "a":
		m.pushUndo()
		m.arrow(tea.KeyRight)
		m.mode = modeInsert
	case "A":
		m.pushUndo()
		m.ta.CursorEnd()
		m.mode = modeInsert
	case "o":
		// Open below at the current line's indentation, Vim autoindent-style.
		m.pushUndo()
		indent := m.lineIndent()
		m.ta.CursorEnd()
		m.send(tea.KeyEnter)
		if indent != "" {
			m.ta.InsertString(indent)
		}
		m.mode = modeInsert
	case "O":
		// Open above at the current line's indentation.
		m.pushUndo()
		indent := m.lineIndent()
		m.ta.CursorStart()
		m.send(tea.KeyEnter)
		m.arrow(tea.KeyUp)
		if indent != "" {
			m.ta.InsertString(indent)
		}
		m.mode = modeInsert

	// --- edits ---
	case "x":
		m.pushUndo()
		n := m.takeCount()
		m.captureDelete(false, func() {
			for i := 0; i < n; i++ {
				m.send(tea.KeyDelete)
			}
		})
	case "D":
		m.pushUndo()
		m.captureDelete(false, func() { m.send(tea.KeyCtrlK) }) // delete to end of line
	case "C":
		m.pushUndo()
		m.captureDelete(false, func() { m.send(tea.KeyCtrlK) })
		m.mode = modeInsert
	case "p":
		return m.pasteCmd(true) // register first; system clipboard as fallback
	case "P":
		return m.pasteCmd(false)

	// --- prefixes / operators ---
	case "d", "c", "g", "r", "y", "<", ">":
		m.pendingCount = m.takeCount()
		m.pending = []rune(key)[0]
		return m, nil

	// --- command line ---
	case ":":
		m.mode = modeCommand
		m.cmd.SetValue("")
		m.cmd.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

// runPending completes a two-key Normal command begun by op (the pending key)
// using msg as the second key. Unknown combinations are silently ignored, as in
// Vim.
func (m Model) runPending(op rune, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	// Mutating operators snapshot first; undo skips no-op snapshots, so invalid
	// combinations (silently ignored below) leave history untouched.
	switch op {
	case 'd', 'c', 'r', '<', '>':
		m.pushUndo()
	}
	switch op {
	case 'f', 'F', 't', 'T':
		if len(msg.Runes) == 1 {
			n := m.pendingCount
			if n < 1 {
				n = 1
			}
			if !m.findChar(op, msg.Runes[0], n) {
				m.status = "not found: " + string(msg.Runes)
			}
		}
	case 'g':
		if key == "g" {
			m.send(tea.KeyCtrlHome) // gg -> top of buffer
		}
	case 'd':
		switch key {
		case "d":
			n := m.pendingCount
			if n < 1 {
				n = 1
			}
			m.captureDelete(true, func() {
				for i := 0; i < n; i++ {
					m.deleteLine()
				}
			})
		case "w", "e":
			m.captureDelete(false, func() {
				m.send2(tea.KeyMsg{Type: tea.KeyDelete, Alt: true}) // dw -> delete word
			})
		case "$":
			m.captureDelete(false, func() { m.send(tea.KeyCtrlK) })
		case "0", "^":
			m.captureDelete(false, func() { m.send(tea.KeyCtrlU) }) // delete to line start
		}
	case 'c':
		switch key {
		case "c":
			m.captureDelete(true, func() {
				m.ta.CursorStart()
				m.send(tea.KeyCtrlK)
			})
			m.mode = modeInsert
		case "w", "e":
			m.captureDelete(false, func() {
				m.send2(tea.KeyMsg{Type: tea.KeyDelete, Alt: true})
			})
			m.mode = modeInsert
		case "$":
			m.captureDelete(false, func() { m.send(tea.KeyCtrlK) })
			m.mode = modeInsert
		}
	case 'y':
		if key == "y" {
			m.yankLines(max(1, m.pendingCount)) // yy / 3yy -> yank line(s)
		}
	case '<':
		if key == "<" {
			m.indentLines(max(1, m.pendingCount), -1) // << / 3<< -> dedent line(s)
		}
	case '>':
		if key == ">" {
			m.indentLines(max(1, m.pendingCount), 1) // >> / 3>> -> indent line(s)
		}
	case 'r':
		// Replace the character under the cursor with the typed rune.
		if len(msg.Runes) == 1 {
			m.send(tea.KeyDelete)
			m.ta, _ = m.ta.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: msg.Runes})
			m.arrow(tea.KeyLeft)
		}
	}
	return m, nil
}

// deleteLine removes the current line (Vim "dd"): clear its content, then pull
// the following line up.
func (m *Model) deleteLine() {
	m.ta.CursorStart()
	m.send(tea.KeyCtrlK)  // delete to end of line
	m.send(tea.KeyDelete) // remove the trailing newline
}

// send forwards a synthetic key (by type) to the textarea, reusing its editing
// logic instead of reimplementing it.
func (m *Model) send(t tea.KeyType) { m.ta, _ = m.ta.Update(tea.KeyMsg{Type: t}) }

// send2 forwards a fully-specified synthetic key (e.g. with Alt) to the textarea.
func (m *Model) send2(k tea.KeyMsg) { m.ta, _ = m.ta.Update(k) }

func (m Model) updateCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = modeNormal
		m.cmd.Blur()
		m.cmd.Prompt = ":"
		m.searchMode = false
		return m, nil
	case tea.KeyEnter:
		if m.searchMode {
			m.mode = modeNormal
			m.cmd.Blur()
			m.cmd.Prompt = ":"
			m.searchMode = false
			m.lastSearch = m.cmd.Value()
			if !m.search(m.lastSearch, 1) {
				m.status = "pattern not found: " + m.lastSearch
			}
			return m, nil
		}
		return m.runCommand(m.cmd.Value())
	}
	var cmd tea.Cmd
	m.cmd, cmd = m.cmd.Update(msg)
	return m, cmd
}

// runCommand dispatches the typed ex-command. ":submit" checks the solution,
// ":w" saves a draft and stays in the editor, ":q" leaves.
func (m Model) runCommand(raw string) (tea.Model, tea.Cmd) {
	m.mode = modeNormal
	m.cmd.Blur()

	switch raw {
	case "submit":
		return m.finish(ActionSubmit)
	case "w", "write":
		if m.save != nil {
			if err := m.save(m.ta.Value()); err != nil {
				m.status = "save failed: " + err.Error()
				return m, nil
			}
		}
		m.status = "saved ✓  (resume later with the same challenge)"
		return m, nil
	case "wq":
		if m.save != nil {
			_ = m.save(m.ta.Value())
		}
		return m.finish(ActionSubmit)
	case "q", "quit":
		return m.finish(ActionQuit)
	case "config":
		// Hand off to the parent to open the config in an external editor.
		return m, func() tea.Msg { return OpenConfigMsg{} }
	default:
		// Not an editor command — forward it to the parent's global dispatcher
		// (e.g. ":topic go", ":progress", ":clear drafts").
		cmd := raw
		return m, func() tea.Msg { return RunCommandMsg{Raw: cmd} }
	}
}

func (m Model) View() string {
	// Visual mode draws its own buffer view: the textarea cannot render a
	// selection, so the editor takes over while one is active.
	if m.vim && m.mode == modeVisual {
		return m.visualView() + "\n" + m.statusLine()
	}
	// Shape the cursor to the current mode so Normal vs Insert is visible at a
	// glance (block vs underline). m is a copy, so this only affects this frame.
	if m.vim && m.mode == modeNormal {
		m.ta.Cursor.Style = normalCursor
	} else {
		m.ta.Cursor.Style = insertCursor
	}
	return highlightSyntax(m.lang, m.ta.View()) + "\n" + m.statusLine()
}

// Highlight applies lang's syntax colors to plain source text. Exported so
// other panes (e.g. the chat transcript's fenced code blocks) can reuse the
// editor's highlighter. "plain"/"text" (and "physics") return s unchanged.
func Highlight(lang, s string) string { return highlightSyntax(lang, s) }

func highlightSyntax(lang, s string) string {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return highlightCode(s, goSyntax())
	case "physics", "plain", "text":
		return s
	default:
		return highlightCode(s, pythonSyntax())
	}
}

type syntaxRules struct {
	keywords     map[string]bool
	types        map[string]bool
	lineComments []string
	blockStart   string
	blockEnd     string
}

func goSyntax() syntaxRules {
	return syntaxRules{
		keywords: words("break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var"),
		types:    words("any bool byte comparable complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr true false nil iota"),
		lineComments: []string{
			"//",
		},
		blockStart: "/*",
		blockEnd:   "*/",
	}
}

func pythonSyntax() syntaxRules {
	return syntaxRules{
		keywords: words("False None True and as assert async await break class continue def del elif else except finally for from global if import in is lambda nonlocal not or pass raise return try while with yield"),
		types:    words("bool bytes dict float int list object set str tuple"),
		lineComments: []string{
			"#",
		},
	}
}

func words(s string) map[string]bool {
	out := map[string]bool{}
	for _, w := range strings.Fields(s) {
		out[w] = true
	}
	return out
}

func highlightCode(s string, rules syntaxRules) string {
	lines := strings.SplitAfter(s, "\n")
	var b strings.Builder
	inBlock := false
	for _, line := range lines {
		b.WriteString(highlightLine(line, rules, &inBlock))
	}
	return b.String()
}

func highlightLine(line string, rules syntaxRules, inBlock *bool) string {
	var b strings.Builder
	for i := 0; i < len(line); {
		if *inBlock {
			end := strings.Index(line[i:], rules.blockEnd)
			if end < 0 {
				b.WriteString(commentStyle.Render(line[i:]))
				return b.String()
			}
			end += i + len(rules.blockEnd)
			b.WriteString(commentStyle.Render(line[i:end]))
			i = end
			*inBlock = false
			continue
		}

		if raw, word, end, ok := readIdentToken(line, i); ok {
			switch {
			case rules.keywords[word]:
				b.WriteString(renderToken(raw, keywordStyle))
			case rules.types[word]:
				b.WriteString(renderToken(raw, typeStyle))
			default:
				b.WriteString(raw)
			}
			i = end
			continue
		}

		if escEnd := ansiEscapeEnd(line, i); escEnd > i {
			b.WriteString(line[i:escEnd])
			i = escEnd
			continue
		}

		if rules.blockStart != "" && strings.HasPrefix(line[i:], rules.blockStart) {
			end := strings.Index(line[i+len(rules.blockStart):], rules.blockEnd)
			if end < 0 {
				b.WriteString(commentStyle.Render(line[i:]))
				*inBlock = true
				return b.String()
			}
			end += i + len(rules.blockStart) + len(rules.blockEnd)
			b.WriteString(commentStyle.Render(line[i:end]))
			i = end
			continue
		}

		if isLineComment(line[i:], rules.lineComments) {
			b.WriteString(commentStyle.Render(line[i:]))
			return b.String()
		}

		if line[i] == '"' || line[i] == '\'' || line[i] == '`' {
			end := stringEnd(line, i)
			b.WriteString(stringStyle.Render(line[i:end]))
			i = end
			continue
		}

		if isDigit(line[i]) {
			end := numberEnd(line, i)
			b.WriteString(numberStyle.Render(line[i:end]))
			i = end
			continue
		}

		b.WriteByte(line[i])
		i++
	}
	return b.String()
}

func readIdentToken(s string, i int) (raw, word string, end int, ok bool) {
	var rb, wb strings.Builder
	seenIdent := false
	for j := i; j < len(s); {
		if escEnd := ansiEscapeEnd(s, j); escEnd > j {
			rb.WriteString(s[j:escEnd])
			j = escEnd
			continue
		}
		c := s[j]
		if !isIdentPart(c) {
			word := wb.String()
			if !seenIdent || !isIdentStart(word[0]) {
				return "", "", i, false
			}
			return rb.String(), wb.String(), j, true
		}
		if !seenIdent && !isIdentStart(c) {
			return "", "", i, false
		}
		seenIdent = true
		rb.WriteByte(c)
		wb.WriteByte(c)
		j++
	}
	word = wb.String()
	if !seenIdent || !isIdentStart(word[0]) {
		return "", "", i, false
	}
	return rb.String(), wb.String(), len(s), true
}

func renderToken(raw string, style lipgloss.Style) string {
	var b strings.Builder
	for i := 0; i < len(raw); {
		if escEnd := ansiEscapeEnd(raw, i); escEnd > i {
			b.WriteString(raw[i:escEnd])
			i = escEnd
			continue
		}
		j := i + 1
		for j < len(raw) {
			if escEnd := ansiEscapeEnd(raw, j); escEnd > j {
				break
			}
			j++
		}
		b.WriteString(style.Render(raw[i:j]))
		i = j
	}
	return b.String()
}

func ansiEscapeEnd(s string, i int) int {
	if i >= len(s) || s[i] != 0x1b {
		return i
	}
	j := i + 1
	if j < len(s) && s[j] == '[' {
		j++
	}
	for ; j < len(s); j++ {
		if s[j] >= '@' && s[j] <= '~' {
			return j + 1
		}
	}
	return len(s)
}

func isLineComment(s string, markers []string) bool {
	for _, marker := range markers {
		if strings.HasPrefix(s, marker) {
			return true
		}
	}
	return false
}

func stringEnd(s string, i int) int {
	quote := s[i]
	j := i + 1
	for j < len(s) {
		if quote != '`' && s[j] == '\\' {
			j += 2
			continue
		}
		if s[j] == quote {
			return j + 1
		}
		j++
	}
	return j
}

func numberEnd(s string, i int) int {
	j := i + 1
	for j < len(s) {
		c := s[j]
		if !(isDigit(c) || isIdentStart(c) || c == '.' || c == '_') {
			break
		}
		j++
	}
	return j
}

func identEnd(s string, i int) int {
	j := i + 1
	for j < len(s) && isIdentPart(s[j]) {
		j++
	}
	return j
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdentStart(c byte) bool {
	return c == '_' || unicode.IsLetter(rune(c))
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}

func (m Model) statusLine() string {
	if m.mode == modeCommand {
		return m.cmd.View()
	}

	var left string
	if m.vim {
		switch m.mode {
		case modeNormal:
			left = normalBadge.Render("NORMAL")
		case modeInsert:
			left = insertBadge.Render("INSERT")
		case modeVisual:
			if m.visualLine {
				left = visualBadge.Render("V-LINE")
			} else {
				left = visualBadge.Render("VISUAL")
			}
		}
	} else {
		left = editBadge.Render("EDIT")
	}

	var hint string
	switch {
	case m.vim && m.mode == modeVisual:
		hint = "  move to select · d/x delete · y yank · c change · </> indent · o swap · esc"
	case m.vim:
		hint = "  :submit  :w  :q  :topic  :progress"
		if m.pending != 0 {
			hint = "  " + string(m.pending) + "…" // operator pending, awaiting motion
		}
	default:
		hint = "  Ctrl-S submit  Ctrl-Q quit"
	}

	line := left + hintStyle.Render(hint)
	if m.status != "" {
		line += "   " + hintStyle.Render(m.status)
	}
	// Keep the status line to a single row: a wrapped line would push the editor
	// past its allotted height and corrupt the surrounding layout.
	if m.width > 0 {
		line = lipgloss.NewStyle().MaxWidth(m.width).Render(line)
	}
	return line
}
