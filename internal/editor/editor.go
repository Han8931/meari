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
	width  int
	height int

	// standalone is set by Open; it makes DoneMsg quit the program. When
	// embedded (New), DoneMsg is left for the parent to handle.
	standalone bool

	// pending holds the first key of a two-key Normal-mode command (operator or
	// prefix): 'd' (delete), 'c' (change), 'g' (goto), 'r' (replace). 0 = none.
	pending rune

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
// moves the cursor to the top, as if opening the file fresh. The textarea
// ignores keys while blurred, so focus is momentarily restored for the move.
func (m *Model) SetValue(s string) {
	m.ta.SetValue(s)
	wasFocused := m.ta.Focused()
	m.ta.Focus()
	m.send(tea.KeyCtrlHome)
	if !wasFocused {
		m.ta.Blur()
	}
}

// Value returns the current buffer contents.
func (m Model) Value() string { return m.ta.Value() }

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
	case modeInsert:
		if msg.Type == tea.KeyEsc {
			m.mode = modeNormal
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

	// A prefix/operator is pending (d, c, g, r): this key completes it.
	if m.pending != 0 {
		op := m.pending
		m.pending = 0
		return m.runPending(op, msg)
	}

	switch msg.String() {
	// --- motions ---
	case "h", "left":
		m.arrow(tea.KeyLeft)
	case "l", "right":
		m.arrow(tea.KeyRight)
	case "j", "down":
		m.arrow(tea.KeyDown)
	case "k", "up":
		m.arrow(tea.KeyUp)
	case "w", "W", "e":
		m.wordForward()
	case "b", "B":
		m.wordBackward()
	case "0":
		m.ta.CursorStart()
	case "^", "_":
		m.ta.CursorStart()
	case "$":
		m.ta.CursorEnd()
	case "G":
		m.send(tea.KeyCtrlEnd) // jump to end of buffer

	// --- enter Insert ---
	case "i":
		m.mode = modeInsert
	case "I":
		m.ta.CursorStart()
		m.mode = modeInsert
	case "a":
		m.arrow(tea.KeyRight)
		m.mode = modeInsert
	case "A":
		m.ta.CursorEnd()
		m.mode = modeInsert
	case "o":
		m.ta.CursorEnd()
		m.send(tea.KeyEnter)
		m.mode = modeInsert
	case "O":
		m.ta.CursorStart()
		m.send(tea.KeyEnter)
		m.arrow(tea.KeyUp)
		m.mode = modeInsert

	// --- edits ---
	case "x":
		m.send(tea.KeyDelete)
	case "D":
		m.send(tea.KeyCtrlK) // delete to end of line
	case "C":
		m.send(tea.KeyCtrlK)
		m.mode = modeInsert
	case "p":
		m.send(tea.KeyCtrlV) // paste from the system clipboard

	// --- prefixes / operators ---
	case "d", "c", "g", "r":
		m.pending = []rune(msg.String())[0]
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
	switch op {
	case 'g':
		if key == "g" {
			m.send(tea.KeyCtrlHome) // gg -> top of buffer
		}
	case 'd':
		switch key {
		case "d":
			m.deleteLine()
		case "w", "e":
			m.send2(tea.KeyMsg{Type: tea.KeyDelete, Alt: true}) // dw -> delete word
		case "$":
			m.send(tea.KeyCtrlK)
		case "0", "^":
			m.send(tea.KeyCtrlU) // delete to line start
		}
	case 'c':
		switch key {
		case "c":
			m.ta.CursorStart()
			m.send(tea.KeyCtrlK)
			m.mode = modeInsert
		case "w", "e":
			m.send2(tea.KeyMsg{Type: tea.KeyDelete, Alt: true})
			m.mode = modeInsert
		case "$":
			m.send(tea.KeyCtrlK)
			m.mode = modeInsert
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

func (m *Model) wordForward()  { m.send2(tea.KeyMsg{Type: tea.KeyRight, Alt: true}) }
func (m *Model) wordBackward() { m.send2(tea.KeyMsg{Type: tea.KeyLeft, Alt: true}) }

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
		return m, nil
	case tea.KeyEnter:
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
		m.status = "unknown command: :" + raw + "  (try :submit, :w, :q, :config)"
		return m, nil
	}
}

func (m Model) View() string {
	// Shape the cursor to the current mode so Normal vs Insert is visible at a
	// glance (block vs underline). m is a copy, so this only affects this frame.
	if m.vim && m.mode == modeNormal {
		m.ta.Cursor.Style = normalCursor
	} else {
		m.ta.Cursor.Style = insertCursor
	}
	return m.ta.View() + "\n" + m.statusLine()
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
		}
	} else {
		left = editBadge.Render("EDIT")
	}

	var hint string
	if m.vim {
		hint = "  :submit  :w  :q  :config"
		if m.pending != 0 {
			hint = "  " + string(m.pending) + "…" // operator pending, awaiting motion
		}
	} else {
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
