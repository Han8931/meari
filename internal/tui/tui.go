// Package tui is the three-pane terminal UI for Meari: a challenge/draft
// list (left), the code editor (center), and an interactive tutor chat (right).
// It replaces the old linear stdin loop — one Bubble Tea program owns all three
// panes, and every AI call and test run is dispatched as an async tea.Cmd so the
// UI never blocks.
package tui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/curriculum"
	"meari/internal/drafts"
	"meari/internal/editor"
	"meari/internal/executor"
	"meari/internal/progress"
	"meari/internal/tutor"
)

// Deps are the collaborators the TUI drives, constructed in main.
type Deps struct {
	Tutor      *tutor.Tutor
	Store      *drafts.Store
	Progress   *progress.State
	Cfg        config.Config
	Topic      string // from --topic; empty means prompt for it at startup
	Curriculum bool   // from --curriculum; start in guided curriculum mode
	ConfigPath string // path to config.toml, for the :config command
	BaseDir    string // working dir, for reloading config after :config
}

type pane int

const (
	paneSidebar pane = iota
	paneEditor
	paneChat
)

type phase int

const (
	phaseSetup phase = iota // the guided startup wizard
	phaseReady              // the three panes
)

// setupStep is the current question in the startup wizard.
type setupStep int

const (
	stepResume   setupStep = iota // offer to continue a saved session (if any)
	stepLanguage                  // pick a language
	stepPath                      // full curriculum vs a specific topic
	stepTopic                     // type the specific topic
	stepLevel                     // beginner / intermediate / advanced
)

// overlayKind selects which full-screen modal is active, if any.
type overlayKind int

const (
	overlayNone     overlayKind = iota
	overlayPicker               // course picker (":topic" with no argument)
	overlayProgress             // progress summary (":progress")
	overlayConfirm              // destructive-action confirmation (":clear progress|drafts")
	overlayHelp                 // command & key reference (":help")
)

// confirmKind is the action a confirmation modal will run if accepted.
type confirmKind int

const (
	confirmNone confirmKind = iota
	confirmClearProgress
	confirmClearDrafts
)

// Model is the root Bubble Tea model coordinating the three panes.
type Model struct {
	deps   Deps
	width  int
	height int
	focus  pane
	phase  phase

	sidebar    sidebarModel
	editor     editor.Model
	chat       chatModel
	topicInput textinput.Model

	// Startup wizard state (phaseSetup).
	setupStep    setupStep
	setupCursor  int
	setupHistory []setupStep // for Esc = back
	lang         string      // chosen language ("python" | "go")
	level        string      // chosen experience level

	topic      string
	current    tutor.Challenge
	curID      *string                    // shared with the editor's save closure
	challenges map[string]tutor.Challenge // id -> challenge, for reload on select
	order      []string                   // sidebar order (stable across rebuilds)
	chatHist   []tutor.ChatTurn           // turns sent to tutor.Chat

	// Per-topic chat contexts: each topic/challenge keeps its own transcript and
	// tutor conversation, restored when the learner returns to it, so the chat
	// pane always shows only the current topic's activity.
	chatKey   string                      // key of the active chat context
	chatByKey map[string][]chatBlock      // saved transcripts
	histByKey map[string][]tutor.ChatTurn // saved tutor conversations

	// Streaming chat reply state: one reply at a time.
	streaming bool
	streamCh  chan streamChunkMsg

	// Curriculum mode: the left pane is the guided learning path instead of the
	// generated-challenge list. Lessons and challenges are pre-authored (no LLM).
	curriculum     bool
	curr           curriculum.Curriculum
	topicByID      map[string]curriculum.Topic
	currentTopicID string

	// horizontal selects the stacked layout (content on top, input on the
	// bottom) instead of the side-by-side columns. Toggled live via :config.
	horizontal bool

	// sidebarCollapsed hides the left tree pane so the editor and chat get the
	// full width. Toggled live via ":fold" (alias ":sidebar").
	sidebarCollapsed bool

	// editorBias shifts the editor/chat split, in percentage points, so the chat
	// pane can be widened when it's cramped. Positive favors the editor (":wide"),
	// negative favors the chat (":compact"). Applied in layout(); 0 is the default.
	editorBias int

	// pendingWindow is set after Ctrl-W; the next key (h/j/k/l) picks a pane,
	// Vim window-command style.
	pendingWindow bool

	// pendingLeader is set after "," in the editor's Normal mode; it starts a
	// leader chord. "n" folds the sidebar; any other key replays the swallowed
	// "," to the editor so its repeat-find binding still works.
	pendingLeader bool

	// Global ex-command line (":topic", ":clear", ":progress"). cmdMode shows the
	// cmdLine input in the status row; it's opened with ":" from the sidebar and
	// also driven by RunCommandMsg forwarded from the editor's own command line.
	cmdMode bool
	cmdLine textinput.Model
	cmdHist editor.CmdHistory

	// Modal overlays drawn full-screen over the panes (like the setup wizard).
	overlay       overlayKind
	pickerCursor  int // selection in the course picker
	confirmKind   confirmKind
	confirmPrompt string // the question shown in the confirm modal

	pending  int    // in-flight async ops; spinner shows while > 0
	loadKind string // label for the most recent op
	spin     spinner.Model
	err      error

	// notice is transient command feedback shown in the status bar (instead of
	// cluttering the chat transcript); it fades after noticeTTL.
	notice   string
	noticeAt time.Time

	// cached layout dims (content sizes inside borders)
	sidebarW int
	editorW  int
	chatW    int
	contentH int
	// horizontal-layout extras: the right column width and the stacked heights.
	rightW  int
	chatH   int
	editorH int
}

// Run constructs the model and runs the full-screen program.
func Run(d Deps) error {
	m := newModel(d)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

func newModel(d Deps) Model {
	curID := new(string)
	save := func(code string) error {
		if *curID == "" {
			return nil
		}
		return d.Store.Save(*curID, code)
	}

	ti := textinput.New()
	ti.Placeholder = "python functions"
	ti.Focus()

	cl := textinput.New()
	cl.Prompt = ":"

	sp := spinner.New(spinner.WithSpinner(spinner.Dot))

	m := Model{
		deps:             d,
		horizontal:       d.Cfg.Horizontal(),
		sidebarCollapsed: d.Cfg.UI.SidebarFolded,
		sidebar:          newSidebar(),
		editor:           editor.New("", d.Cfg.VimEditor(), save),
		chat:             newChat(),
		topicInput:       ti,
		cmdLine:          cl,
		curID:            curID,
		challenges:       map[string]tutor.Challenge{},
		chatByKey:        map[string][]chatBlock{},
		histByKey:        map[string][]tutor.ChatTurn{},
		spin:             sp,
	}
	m.editor.SetShowLineNumbers(d.Cfg.LineNumbers())
	m.seedOrder()

	// Decide the starting screen.
	switch {
	case d.Topic != "":
		// Custom topic via flag: jump straight in; Init fetches lesson+challenge.
		m.topic = d.Topic
		m.phase = phaseReady
		m.pending = 2
		m.setFocus(paneEditor)
	case d.Curriculum:
		// Curriculum via flag: resume the saved session if any, else Python beginner.
		lang, level, startID := "python", curriculum.Beginner, ""
		if last := d.Progress.Last; last != nil {
			lang, level, startID = last.Lang, last.Level, last.TopicID
		}
		m.loadCurriculum(lang, level, startID)
		m.phase = phaseReady
	default:
		// Guided wizard. Start on the resume step only if there's a saved session.
		if d.Progress.Last != nil {
			m.setupStep = stepResume
		} else {
			m.setupStep = stepLanguage
		}
	}
	return m
}

// seedOrder pre-populates the sidebar order from drafts and progress on disk so
// past work shows up across restarts (stable, sorted).
func (m *Model) seedOrder() {
	seen := map[string]bool{}
	addAll := func(ids []string) {
		sort.Strings(ids)
		for _, id := range ids {
			if !seen[id] {
				seen[id] = true
				m.order = append(m.order, id)
			}
		}
	}
	if ids, err := m.deps.Store.IDs(); err == nil {
		addAll(ids)
	}
	var pids []string
	for id := range m.deps.Progress.Challenges {
		pids = append(pids, id)
	}
	addAll(pids)
	m.rebuildSidebar()
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spin.Tick}
	switch {
	case m.phase == phaseSetup:
		cmds = append(cmds, textinput.Blink)
	case m.curriculum:
		// Curriculum content is pre-authored and already loaded; nothing to fetch.
		cmds = append(cmds, m.editor.Focus())
	default:
		// --topic was supplied: focus the editor and kick off lesson + challenge.
		cmds = append(cmds,
			m.editor.Focus(),
			lessonCmd(m.deps.Tutor, m.topic),
			challengeCmd(m.deps.Tutor, m.topic),
		)
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		// Mirror in-flight work as an animated progress line inside the chat
		// pane, where the eye already is — the status bar alone is easy to miss.
		if m.pending > 0 {
			m.chat.setBusy(m.loadKind)
			m.chat.tickBusy()
		} else {
			m.chat.setBusy("")
		}
		return m, cmd

	case lessonMsg:
		m.pending--
		m.chat.append(roleLesson, msg.text)
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant", Content: msg.text})
		return m, nil

	case challengeMsg:
		m.pending--
		return m, m.loadChallenge(msg.ch)

	case feedbackMsg:
		m.pending--
		m.chat.append(roleTutor, msg.text)
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant", Content: msg.text})
		return m, nil

	case streamChunkMsg:
		return m.handleStreamChunk(msg)

	case answerMsg:
		m.pending--
		m.chat.append(roleLesson, "Model answer\n\n"+msg.text)
		// Join the conversation so follow-up questions can refer to it.
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant", Content: "Model answer:\n" + msg.text})
		return m, nil

	case runResultMsg:
		m.pending--
		return m, m.handleRunResult(msg)

	case errMsg:
		m.pending--
		m.err = msg.err
		m.chat.append(roleSystem, "⚠ "+msg.kind+" failed: "+msg.err.Error())
		return m, nil

	case editor.DoneMsg:
		switch msg.Action {
		case editor.ActionSubmit:
			return m, m.startRun()
		case editor.ActionQuit:
			return m, m.quit()
		}
		return m, nil

	case editor.OpenConfigMsg:
		return m, m.openConfig()

	case editor.RunCommandMsg:
		return m.runEx(msg.Raw)

	case configReloadMsg:
		m.applyConfigReload(msg)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}

	return m.forwardToFocus(msg)
}

// handleKey dispatches a key: topic-entry first, then global bindings, then the
// focused pane.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.phase == phaseSetup {
		return m.updateSetup(msg)
	}

	// Modals and the command line capture all keys while they're open.
	if m.overlay != overlayNone {
		return m.updateOverlay(msg)
	}
	if m.cmdMode {
		return m.updateCmdLine(msg)
	}

	// Ctrl-W starts a window command; the next key chooses a pane by direction.
	if m.pendingWindow {
		m.pendingWindow = false
		switch msg.String() {
		case "h", "left", "k", "up", "shift+tab":
			return m, m.focusDir(-1)
		case "l", "right", "j", "down", "tab", "ctrl+w":
			return m, m.focusDir(1)
		}
		return m, nil // unknown window command: ignore, like Vim
	}

	switch msg.String() {
	case "ctrl+c":
		return m, m.quit()
	case "ctrl+w":
		// Focus moves only via the Vim window chord (Ctrl-W then h/l); bare Tab is
		// left for the panes (e.g. indenting in the editor).
		m.pendingWindow = true
		return m, nil
	case "ctrl+r":
		// In the editor, Ctrl-R is Vim redo; from any other pane it runs the
		// tests (running is always available via Ctrl-S / :submit while editing).
		if m.focus != paneEditor {
			return m, m.startRun()
		}
	case "ctrl+n":
		return m, m.nextChallenge()
	}

	// A leader chord lives for exactly one keystroke: clear it here so a stray
	// "," can never carry across a focus change or non-editor key.
	leader := m.pendingLeader
	m.pendingLeader = false

	switch m.focus {
	case paneSidebar:
		if msg.String() == ":" {
			return m.openCmdLine()
		}
		var enter bool
		m.sidebar, enter = m.sidebar.Update(msg)
		if enter {
			return m, m.openSelected()
		}
		return m, nil

	case paneEditor:
		// Leader chord ",n" folds the sidebar — but only in Vim Normal mode, so
		// it never disturbs typing or a pending multi-key Vim command.
		if leader {
			if msg.String() == "n" {
				return m.cmdFold()
			}
			// Not the fold chord: replay the swallowed "," (its Normal-mode
			// repeat-find), then deliver the key that followed it.
			comma := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{','}}
			tm, _ := m.editor.Update(comma)
			m.editor = tm.(editor.Model)
			tm, cmd := m.editor.Update(msg)
			m.editor = tm.(editor.Model)
			return m, cmd
		}
		if msg.String() == "," && m.editor.NormalMode() {
			m.pendingLeader = true
			return m, nil
		}
		tm, cmd := m.editor.Update(msg)
		m.editor = tm.(editor.Model)
		return m, cmd

	case paneChat:
		switch msg.String() {
		// Copy the tutor's last reply: Alt+O on Linux; on macOS, Option+O —
		// which arrives as "ø"/"Ø" unless the terminal sends Option as Meta.
		// (Cmd+O never reaches a terminal app; the emulator consumes it.)
		case "alt+o", "ø", "Ø":
			m.flash(copyChat(&m.chat, ""))
			return m, nil
		}
		if msg.Type == tea.KeyEnter {
			return m.submitChat()
		}
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(msg)
		return m, cmd
	}
	return m, nil
}

// flash shows transient feedback in the status bar for a few seconds.
func (m *Model) flash(s string) {
	if s == "" {
		return
	}
	m.notice = s
	m.noticeAt = time.Now()
}

// handleMouse routes wheel events to the pane under the cursor, like ranger/lf:
// scrolling never changes focus, it just moves the hovered pane. A left click
// focuses the pane under the cursor. Other mouse events (motion, other buttons)
// are ignored.
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.phase == phaseSetup || m.overlay != overlayNone || m.cmdMode {
		return m, nil
	}
	p, ok := m.paneAt(msg.X, msg.Y)
	if !ok {
		return m, nil
	}

	if tea.MouseEvent(msg).IsWheel() {
		switch p {
		case paneChat:
			var cmd tea.Cmd
			m.chat, cmd = m.chat.Update(msg) // the viewport scrolls on wheel events
			return m, cmd
		case paneSidebar:
			// The sidebar has no viewport; move the selection instead (ranger-style).
			switch msg.Button {
			case tea.MouseButtonWheelDown:
				m.sidebar.move(1)
			case tea.MouseButtonWheelUp:
				m.sidebar.move(-1)
			}
			return m, nil
		case paneEditor:
			tm, cmd := m.editor.Update(msg)
			m.editor = tm.(editor.Model)
			return m, cmd
		}
		return m, nil
	}

	// Left click: focus the pane under the cursor.
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		return m, m.setFocus(p)
	}
	return m, nil
}

// paneAt maps a terminal cell to the pane drawn there, accounting for the title
// bar (top row), status bar (bottom row), and the one-cell border around each
// box. ok is false for cells outside any pane (the bars or border gutters).
func (m Model) paneAt(x, y int) (pane, bool) {
	if y < 1 || y > m.height-2 { // row 0 is the title bar, the last row the status bar
		return 0, false
	}
	// When folded, the sidebar occupies no columns; sidebarW is 0 and its border
	// is gone, so the +2 offsets below collapse to the editor's left edge.
	sidebarSpan := m.sidebarW + 2
	if m.sidebarCollapsed {
		sidebarSpan = 0
	}
	if m.horizontal {
		// sidebar on the left; chat stacked over the editor on the right.
		if x < sidebarSpan {
			return paneSidebar, true
		}
		if y < 1+m.chatH+2 {
			return paneChat, true
		}
		return paneEditor, true
	}
	// vertical: sidebar | editor | chat, each box +2 columns for its border.
	switch {
	case x < sidebarSpan:
		return paneSidebar, true
	case x < sidebarSpan+m.editorW+2:
		return paneEditor, true
	default:
		return paneChat, true
	}
}

// --- global ex-command line ---

// openCmdLine starts the ":" command prompt, shown in the status row.
func (m Model) openCmdLine() (tea.Model, tea.Cmd) {
	m.cmdMode = true
	m.cmdLine.SetValue("")
	m.cmdHist.Open()
	return m, m.cmdLine.Focus()
}

// updateCmdLine drives the command prompt: Enter runs it, Esc cancels, Ctrl-C
// still quits, anything else edits the text.
func (m Model) updateCmdLine(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if s, ok := m.cmdHist.Prev(m.cmdLine.Value()); ok {
			m.cmdLine.SetValue(s)
			m.cmdLine.CursorEnd()
		}
		return m, nil
	case tea.KeyDown:
		if s, ok := m.cmdHist.Next(); ok {
			m.cmdLine.SetValue(s)
			m.cmdLine.CursorEnd()
		}
		return m, nil
	case tea.KeyCtrlC:
		return m, m.quit()
	case tea.KeyEnter:
		raw := strings.TrimSpace(m.cmdLine.Value())
		m.cmdMode = false
		m.cmdLine.Blur()
		if raw == "" {
			return m, nil
		}
		m.cmdHist.Record(raw)
		return m.runEx(raw)
	case tea.KeyEsc:
		m.cmdMode = false
		m.cmdLine.Blur()
		return m, nil
	}
	var cmd tea.Cmd
	m.cmdLine, cmd = m.cmdLine.Update(msg)
	return m, cmd
}

// runEx dispatches an ex-command (without the leading colon), from either the
// global prompt or the editor's forwarded command line.
func (m Model) runEx(raw string) (tea.Model, tea.Cmd) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return m, nil
	}
	switch fields[0] {
	case "topic", "subject", "course":
		return m.cmdTopic(fields[1:])
	case "clear":
		return m.cmdClear(fields[1:])
	case "fold", "sidebar":
		return m.cmdFold()
	case "compact":
		return m.cmdResizeEditor(-editorBiasStep)
	case "wide":
		return m.cmdResizeEditor(editorBiasStep)
	case "answer":
		return m.cmdAnswer()
	case "copy", "yank":
		what := ""
		if len(fields) > 1 {
			what = fields[1]
		}
		m.flash(copyChat(&m.chat, what))
		return m, nil
	case "paste":
		m.flash(pasteChat(&m.chat))
		return m, m.setFocus(paneChat) // land where the pasted text is
	case "progress":
		m.overlay = overlayProgress
		return m, nil
	case "help":
		m.overlay = overlayHelp
		return m, nil
	case "config":
		return m, m.openConfig()
	case "learn", "essay", "grade":
		m.flash(":" + fields[0] + " lives in the learning vault — quit and run `meari notes` (or `meari serve`)")
		return m, nil
	default:
		m.flash("unknown command: :" + raw + "  (try :help)")
		return m, nil
	}
}

// cmdTopic switches to another course. With no argument it opens the picker.
func (m Model) cmdTopic(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		m.overlay = overlayPicker
		m.pickerCursor = courseIndex(m.lang)
		return m, nil
	}
	return m.switchCourse(strings.ToLower(args[0]))
}

// switchCourse enters the curriculum for course at the learner's current level
// (defaulting to beginner), or reports an error if there's no such course.
func (m Model) switchCourse(course string) (tea.Model, tea.Cmd) {
	if !curriculum.HasCurriculum(course) {
		m.flash("no course \"" + course + "\" — try: " + strings.Join(curriculum.Languages(), ", "))
		return m, nil
	}
	level := m.level
	if level == "" {
		level = curriculum.Beginner
	}
	m.chat.append(roleSystem, "— now learning "+titleCase(course)+" ("+level+") —")
	return m, m.loadCurriculum(course, level, "")
}

// cmdAnswer reveals a model solution for the current challenge. The learner
// explicitly asked, so revealing is fine — unlike run feedback, which never does.
func (m Model) cmdAnswer() (tea.Model, tea.Cmd) {
	if m.current.ID == "" {
		m.flash("no challenge open — pick one on the left first")
		return m, nil
	}
	m.pending++
	m.loadKind = "writing model answer"
	return m, answerCmd(m.deps.Tutor, m.current)
}

// cmdClear handles ":clear" (chat transcript), ":clear progress", and
// ":clear drafts". The destructive variants route through a confirm modal.
func (m Model) cmdClear(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		m.chat.blocks = nil
		m.chatHist = nil
		m.chat.reflow()
		return m, nil
	}
	switch args[0] {
	case "progress":
		m.overlay = overlayConfirm
		m.confirmKind = confirmClearProgress
		m.confirmPrompt = "Erase ALL learning progress — completed topics, attempts, and your resume point?"
	case "drafts":
		m.overlay = overlayConfirm
		m.confirmKind = confirmClearDrafts
		m.confirmPrompt = "Delete ALL saved draft code?"
	default:
		m.chat.append(roleSystem, "clear what? try :clear · :clear progress · :clear drafts")
	}
	return m, nil
}

// cmdFold toggles the left tree pane. When folding away the pane that currently
// has focus, focus jumps to the editor so keys never vanish into a hidden pane;
// the layout is re-flowed so the editor and chat reclaim the freed width.
func (m Model) cmdFold() (tea.Model, tea.Cmd) {
	m.sidebarCollapsed = !m.sidebarCollapsed
	var cmd tea.Cmd
	if m.sidebarCollapsed && m.focus == paneSidebar {
		cmd = m.setFocus(paneEditor)
	}
	m.layout()
	if m.sidebarCollapsed {
		m.flash("Sidebar folded — :fold to bring it back")
	}
	return m, cmd
}

// editorBias bounds: one step per :compact / :wide, clamped so neither the
// editor nor the chat can be squeezed past a usable share. The range is wide
// enough that repeated :compact lets the chat take well over half the width
// for reading long lessons and feedback.
const (
	editorBiasStep = 8
	editorBiasMax  = 32
)

// cmdResizeEditor nudges the editor/chat split by delta percentage points
// (":compact" shrinks the editor so the chat gets more room, ":wide" grows it),
// clamps it, re-flows, and reports the new emphasis.
func (m Model) cmdResizeEditor(delta int) (tea.Model, tea.Cmd) {
	prev := m.editorBias
	m.editorBias = clampRange(m.editorBias+delta, -editorBiasMax, editorBiasMax)
	if m.editorBias == prev {
		edge := "widest"
		if delta < 0 {
			edge = "narrowest"
		}
		m.flash("Editor already at its " + edge + " — chat can't go further")
		return m, nil
	}
	m.layout()
	switch {
	case m.editorBias < 0:
		m.flash("Editor narrowed — more room for chat (:wide to grow it back)")
	case m.editorBias > 0:
		m.flash("Editor widened (:compact to give chat more room)")
	default:
		m.flash("Editor/chat split reset to default")
	}
	return m, nil
}

// runConfirm performs the pending destructive action after the learner accepts.
func (m Model) runConfirm() (tea.Model, tea.Cmd) {
	switch m.confirmKind {
	case confirmClearProgress:
		if err := m.deps.Progress.Reset(); err != nil {
			m.flash("⚠ couldn't clear progress: " + err.Error())
		} else {
			m.flash("✓ learning progress cleared")
			m.rebuildSidebar()
		}
	case confirmClearDrafts:
		if err := m.deps.Store.ClearAll(); err != nil {
			m.flash("⚠ couldn't clear drafts: " + err.Error())
		} else {
			m.flash("✓ saved drafts cleared")
		}
	}
	m.overlay = overlayNone
	m.confirmKind = confirmNone
	return m, nil
}

// --- modal overlays ---

// updateOverlay handles keys while a full-screen modal is open.
func (m Model) updateOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, m.quit()
	}
	switch m.overlay {
	case overlayPicker:
		courses := curriculum.Languages()
		switch msg.String() {
		case "esc", "q":
			m.overlay = overlayNone
		case "j", "down":
			if m.pickerCursor < len(courses)-1 {
				m.pickerCursor++
			}
		case "k", "up":
			if m.pickerCursor > 0 {
				m.pickerCursor--
			}
		case "enter":
			m.overlay = overlayNone
			return m.switchCourse(courses[m.pickerCursor])
		}
	case overlayProgress, overlayHelp:
		switch msg.String() {
		case "esc", "q", "enter":
			m.overlay = overlayNone
		}
	case overlayConfirm:
		switch msg.String() {
		case "y", "Y":
			return m.runConfirm()
		case "n", "N", "esc", "q":
			m.overlay = overlayNone
			m.confirmKind = confirmNone
		}
	}
	return m, nil
}

// overlayView renders the active modal as a centered card.
func (m Model) overlayView() string {
	switch m.overlay {
	case overlayPicker:
		return m.pickerView()
	case overlayProgress:
		return m.progressView()
	case overlayConfirm:
		return modalCard("Are you sure?", errStyle.Render("This cannot be undone.")+"\n\n"+m.confirmPrompt, "y to confirm · n / esc to cancel")
	case overlayHelp:
		return helpView()
	}
	return ""
}

func (m Model) pickerView() string {
	courses := curriculum.Languages()
	var b strings.Builder
	for i, c := range courses {
		if i == m.pickerCursor {
			b.WriteString(selectedRow.Render("▸ " + titleCase(c)))
		} else {
			b.WriteString("  " + titleCase(c))
		}
		if i < len(courses)-1 {
			b.WriteString("\n")
		}
	}
	return modalCard("Switch course", b.String(), "↑/↓ or j/k · enter to switch · esc to cancel")
}

func (m Model) progressView() string {
	var b strings.Builder
	for _, course := range curriculum.Languages() {
		done, total := m.courseCompletion(course)
		b.WriteString(courseProgressLine(course, done, total))
		b.WriteString("\n")
	}
	attempts, passes := 0, 0
	for _, e := range m.deps.Progress.Challenges {
		attempts += e.Attempts
		passes += e.Passes
	}
	b.WriteString("\n")
	b.WriteString(hintStyle.Render(fmt.Sprintf("Challenge runs: %d attempt(s), %d passed", attempts, passes)))
	return modalCard("Your progress", b.String(), "esc / q to close")
}

// helpView lists the global commands and key bindings (":help").
func helpView() string {
	cmds := strings.Join([]string{
		bold("Commands"),
		"  :help              this screen",
		"  :topic <course>    switch course (python/go/physics)",
		"  :subject <course>  alias for :topic",
		"  :fold              fold/unfold the left tree pane",
		"  :compact / :wide   shrink/grow the editor (frees chat space)",
		"  :answer            reveal a model solution for the open challenge",
		"  :progress          progress summary",
		"  :copy [code|all]   copy the tutor's last reply / its code / everything",
		"  :paste             paste the clipboard into the chat input",
		"  :clear             clear the chat transcript",
		"  :clear progress    erase saved progress (asks first)",
		"  :clear drafts      delete saved drafts (asks first)",
		"  :config            edit config.toml",
		"",
		bold("Keys"),
		"  : (left pane)      open this command prompt",
		"  ⌃w then h/l        move focus between panes",
		"  mouse click        focus the pane under the cursor",
		"  ⌃r / ⌃n            run tests / next challenge (in the editor, ⌃r is Vim redo;",
		"                     run with ⌃s or :submit there)",
		"  chat ⌃f ⌃b ⌃d ⌃u   page / half-page scroll",
		"  mouse wheel        scroll the pane under the cursor",
		"  ⌃c                 quit",
	}, "\n")
	return modalCard("Meari — help", cmds, "esc / q to close")
}

func bold(s string) string { return lipgloss.NewStyle().Bold(true).Render(s) }

// courseCompletion counts done vs total topics across all levels of a course.
func (m Model) courseCompletion(course string) (done, total int) {
	for _, level := range []string{curriculum.Beginner, curriculum.Intermediate, curriculum.Advanced} {
		c, ok := curriculum.For(course, level)
		if !ok {
			continue
		}
		for _, t := range c.Topics() {
			total++
			if m.deps.Progress.TopicStatus(t.ID) == "done" {
				done++
			}
		}
	}
	return done, total
}

// courseProgressLine renders one course's completion as a labeled bar.
func courseProgressLine(course string, done, total int) string {
	const width = 14
	filled := 0
	if total > 0 {
		filled = done * width / total
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%-9s %s  %d/%d", titleCase(course), bar, done, total)
}

// modalCard wraps a title/body/hint in the same bordered card the setup wizard
// uses, so modals look consistent with the rest of the app.
func modalCard(title, body, hint string) string {
	inner := lipgloss.NewStyle().Bold(true).Render(title) + "\n\n" + body
	if hint != "" {
		inner += "\n\n" + hintStyle.Render(hint)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 3).
		Width(56).
		Render(inner)
}

// courseIndex returns the position of course in the supported list, or 0.
func courseIndex(course string) int {
	for i, c := range curriculum.Languages() {
		if c == course {
			return i
		}
	}
	return 0
}

// setupOptions returns the selectable options for the current selection step
// (empty for the free-text topic step).
func (m Model) setupOptions() []string {
	switch m.setupStep {
	case stepResume:
		cont := "Continue where you left off"
		if last := m.deps.Progress.Last; last != nil {
			cont = "Continue: " + titleCase(last.Lang) + " · " + last.Level + " — " + last.Title
		}
		return []string{cont, "Start something new"}
	case stepLanguage:
		langs := curriculum.Languages()
		opts := make([]string, 0, len(langs))
		for _, lang := range langs {
			opts = append(opts, titleCase(lang))
		}
		return opts
	case stepPath:
		return []string{"Start from the beginning (full curriculum)", "Learn a specific topic (AI-generated)"}
	case stepLevel:
		return []string{"Beginner", "Intermediate", "Advanced"}
	}
	return nil
}

func (m *Model) gotoStep(s setupStep) {
	m.setupStep = s
	m.setupCursor = 0
	if s == stepTopic {
		m.topicInput.SetValue("")
		m.topicInput.Placeholder = "e.g. list comprehensions"
		m.topicInput.Focus()
	}
}

func (m *Model) advance(to setupStep) {
	m.setupHistory = append(m.setupHistory, m.setupStep)
	m.gotoStep(to)
}

func (m *Model) back() {
	if len(m.setupHistory) == 0 {
		return
	}
	prev := m.setupHistory[len(m.setupHistory)-1]
	m.setupHistory = m.setupHistory[:len(m.setupHistory)-1]
	m.gotoStep(prev)
}

// updateSetup drives the startup wizard.
func (m Model) updateSetup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if msg.String() == "esc" {
		m.back()
		return m, nil
	}

	// The only free-text step is typing a specific topic.
	if m.setupStep == stepTopic {
		if msg.String() == "enter" {
			return m.setupSubmitInput()
		}
		var cmd tea.Cmd
		m.topicInput, cmd = m.topicInput.Update(msg)
		return m, cmd
	}

	// Selection steps.
	opts := m.setupOptions()
	switch msg.String() {
	case "j", "down":
		if m.setupCursor < len(opts)-1 {
			m.setupCursor++
		}
	case "k", "up":
		if m.setupCursor > 0 {
			m.setupCursor--
		}
	case "enter":
		return m.setupSelect()
	}
	return m, nil
}

// setupSelect advances the wizard after an option is chosen on a selection step.
func (m Model) setupSelect() (tea.Model, tea.Cmd) {
	switch m.setupStep {
	case stepResume:
		if m.setupCursor == 0 { // continue the saved session
			return m.finishResume()
		}
		m.advance(stepLanguage)
	case stepLanguage:
		langs := curriculum.Languages()
		if m.setupCursor < 0 || m.setupCursor >= len(langs) {
			m.setupCursor = 0
		}
		m.lang = langs[m.setupCursor]
		if m.lang == "python" {
			m.advance(stepPath) // Python offers curriculum or a specific topic
		} else {
			m.curriculum = true // non-Python languages are curriculum-only
			m.advance(stepLevel)
		}
	case stepPath:
		if m.setupCursor == 0 {
			m.curriculum = true
			m.advance(stepLevel)
		} else {
			m.curriculum = false
			m.advance(stepTopic)
		}
	case stepLevel:
		m.level = []string{"beginner", "intermediate", "advanced"}[m.setupCursor]
		return m.finishSetup()
	}
	return m, nil
}

// setupSubmitInput advances after typing a specific (custom) topic.
func (m Model) setupSubmitInput() (tea.Model, tea.Cmd) {
	v := strings.TrimSpace(m.topicInput.Value())
	if v == "" {
		v = "basics"
	}
	m.topic = strings.TrimSpace(m.lang + " " + v)
	m.advance(stepLevel)
	return m, nil
}

// finishSetup leaves the wizard and starts the session: a pre-authored
// curriculum (no LLM) or an AI-generated custom topic.
func (m Model) finishSetup() (tea.Model, tea.Cmd) {
	m.deps.Tutor.SetLevel(m.level)
	m.phase = phaseReady
	m.topicInput.Blur()

	if m.curriculum {
		return m, m.loadCurriculum(m.lang, m.level, "")
	}

	// Custom topic via the LLM.
	m.pending += 2
	m.loadKind = "loading lesson"
	return m, tea.Batch(
		m.setFocus(paneEditor),
		lessonCmd(m.deps.Tutor, m.topic),
		challengeCmd(m.deps.Tutor, m.topic),
	)
}

// finishResume continues the saved curriculum session.
func (m Model) finishResume() (tea.Model, tea.Cmd) {
	last := m.deps.Progress.Last
	m.phase = phaseReady
	m.topicInput.Blur()
	return m, m.loadCurriculum(last.Lang, last.Level, last.TopicID)
}

// forwardToFocus routes non-key messages (e.g. cursor blinks) to whatever
// component currently has focus, so its cursor keeps animating.
func (m Model) forwardToFocus(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.cmdMode {
		var cmd tea.Cmd
		m.cmdLine, cmd = m.cmdLine.Update(msg)
		return m, cmd
	}
	if m.phase == phaseSetup {
		var cmd tea.Cmd
		m.topicInput, cmd = m.topicInput.Update(msg)
		return m, cmd
	}
	switch m.focus {
	case paneEditor:
		tm, cmd := m.editor.Update(msg)
		m.editor = tm.(editor.Model)
		return m, cmd
	case paneChat:
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(msg)
		return m, cmd
	}
	return m, nil
}

// chatContext describes what the learner is currently looking at — the lesson,
// the challenge, and their in-progress code — so chat replies stay grounded in
// the material instead of answering in a vacuum.
func (m Model) chatContext() string {
	var b strings.Builder
	if m.curriculum {
		if t, ok := m.topicByID[m.currentTopicID]; ok {
			b.WriteString("Current lesson — " + t.Title + ":\n" + t.Lesson + "\n\n")
		}
	} else if m.topic != "" {
		b.WriteString("Current topic: " + m.topic + "\n\n")
	}
	if m.current.ID != "" {
		b.WriteString("Current challenge (" + challengeLang(m.current) + "): " + m.current.Prompt + "\n\n")
		b.WriteString("Learner's current code:\n" + m.editor.Value() + "\n")
	}
	return strings.TrimSpace(b.String())
}

// submitChat sends the chat input as a question to the tutor, streaming the
// reply into the transcript. The call carries the current lesson/challenge/code
// as context so answers relate to what's on screen.
func (m Model) submitChat() (tea.Model, tea.Cmd) {
	if m.streaming {
		m.flash("the tutor is still replying — one question at a time")
		return m, nil
	}
	text, ok := m.chat.submit()
	if !ok {
		return m, nil
	}
	m.chat.append(roleUser, text)
	m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "user", Content: text})
	m.pending++
	m.loadKind = "tutor thinking"
	m.streaming = true
	m.chat.beginStream()

	tut := m.deps.Tutor
	ctxText := core.ClampContext(m.chatContext())
	hist := core.TrimTurns(append([]tutor.ChatTurn(nil), m.chatHist...)) // copy: the goroutine outlives this Update
	ch, cmd := startChatStream(func(onDelta func(string)) (string, error) {
		return tut.ChatStream(context.Background(), ctxText, hist, onDelta)
	})
	m.streamCh = ch
	return m, cmd
}

// handleStreamChunk advances a streaming tutor reply: grow the transcript on
// each delta, finalize the conversation history when done.
func (m Model) handleStreamChunk(msg streamChunkMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.pending--
		m.streaming = false
		m.chat.failStream("⚠ chat failed: " + msg.err.Error())
		return m, nil
	}
	if msg.done {
		m.pending--
		m.streaming = false
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant", Content: msg.full})
		return m, nil
	}
	m.chat.appendStream(msg.delta)
	return m, listenStream(m.streamCh)
}

// openSelected loads the highlighted sidebar item's challenge into the editor.
func (m *Model) openSelected() tea.Cmd {
	it, ok := m.sidebar.selected()
	if !ok {
		return nil
	}
	if m.curriculum {
		if t, ok := m.topicByID[it.id]; ok {
			return m.startTopic(t)
		}
		return nil
	}
	if ch, ok := m.challenges[it.id]; ok {
		return m.loadChallenge(ch)
	}
	// A draft/progress id from a previous session — we don't have its tests this
	// session, so it can't be run. Tell the learner.
	m.flash("No challenge data for \"" + it.id + "\" this session — Ctrl-N generates a new one")
	return nil
}

// loadCurriculum enters curriculum mode for (lang, level), indexes its topics,
// and starts at startID (or the first unfinished topic when startID is empty).
func (m *Model) loadCurriculum(lang, level, startID string) tea.Cmd {
	c, ok := curriculum.For(lang, level)
	if !ok {
		m.chat.append(roleSystem, "⚠ No curriculum for "+lang+" / "+level+".")
		return nil
	}
	m.curriculum = true
	m.lang, m.level = lang, level
	m.curr = c
	m.topicByID = map[string]curriculum.Topic{}
	for _, t := range c.Topics() {
		m.topicByID[t.ID] = t
	}
	m.deps.Tutor.SetLevel(level)

	start := m.firstIncompleteTopic()
	if startID != "" {
		if t, ok := m.topicByID[startID]; ok {
			start = t
		}
	}
	return m.startTopic(start)
}

// firstIncompleteTopic returns the first not-yet-done topic, or the first topic
// if the whole level is complete.
func (m *Model) firstIncompleteTopic() curriculum.Topic {
	topics := m.curr.Topics()
	for _, t := range topics {
		if m.deps.Progress.TopicStatus(t.ID) != "done" {
			return t
		}
	}
	if len(topics) > 0 {
		return topics[0]
	}
	return curriculum.Topic{}
}

// switchChatContext swaps the chat pane to key's own transcript and tutor
// conversation, saving the outgoing topic's first so it can be restored when
// the learner comes back. It reports whether the new context is fresh (so the
// caller knows to show the intro content once, not on every revisit).
func (m *Model) switchChatContext(key string) (fresh bool) {
	if key == m.chatKey {
		return false
	}
	prev := m.chatKey
	if prev != "" {
		m.chatByKey[prev] = m.chat.snapshot()
		m.histByKey[prev] = m.chatHist
	}
	m.chatKey = key
	saved, visited := m.chatByKey[key]
	if !visited && prev == "" {
		// First topic of the session: inherit the startup transcript as-is (a
		// custom topic's lesson may already have streamed in before its first
		// challenge arrived).
		return true
	}
	m.chat.restore(saved)
	m.chatHist = m.histByKey[key]
	return !visited
}

// startTopic switches the active curriculum topic, showing its pre-authored
// lesson in chat and loading its challenge into the editor — no LLM call.
func (m *Model) startTopic(t curriculum.Topic) tea.Cmd {
	if m.current.ID != "" {
		_ = m.deps.Store.Save(m.current.ID, m.editor.Value())
	}
	m.currentTopicID = t.ID
	m.topic = t.Title
	m.deps.Progress.MarkTopicInProgress(t.ID)
	m.deps.Progress.SetLast(m.curr.Lang, m.curr.Level, t.ID, t.Title)
	_ = m.deps.Progress.Save()

	ch := tutor.Challenge{
		ID:          t.ID,
		Prompt:      t.Challenge.Prompt,
		StarterCode: t.Challenge.StarterCode,
		Tests:       t.Challenge.Tests,
		Lang:        m.curr.Lang,
	}
	m.current = ch
	*m.curID = ch.ID
	m.editor.SetLanguage(ch.Lang)
	m.chat.setCodeLang(challengeLang(ch))

	code, stale := m.loadStarterOrDraft(ch)
	m.editor.SetValue(code)
	m.layout() // the pinned prompt header's height depends on the new prompt
	m.rebuildSidebar()
	// Each topic keeps its own chat: show the lesson only on the first visit —
	// revisits restore the prior transcript (which already contains it). The
	// challenge statement is NOT echoed here: it lives at the top of the editor
	// as a comment, where it stays visible regardless of chat length.
	if m.switchChatContext("topic:" + t.ID) {
		m.chat.append(roleLesson, t.Title+"\n\n"+t.Lesson)
	}
	if stale {
		m.chat.append(roleSystem, staleDraftNotice)
	}
	return m.setFocus(paneEditor)
}

// loadChallenge saves the current buffer, then loads ch (its saved draft or
// starter code) into the editor and focuses it.
func (m *Model) loadChallenge(ch tutor.Challenge) tea.Cmd {
	if m.current.ID != "" {
		_ = m.deps.Store.Save(m.current.ID, m.editor.Value())
	}
	m.challenges[ch.ID] = ch
	m.addOrder(ch.ID)
	m.current = ch
	*m.curID = ch.ID
	m.editor.SetLanguage(ch.Lang)
	m.chat.setCodeLang(challengeLang(ch))

	code, stale := m.loadStarterOrDraft(ch)
	m.editor.SetValue(code)
	m.layout() // the pinned prompt header's height depends on the new prompt
	m.rebuildSidebar()
	m.switchChatContext("challenge:" + ch.ID)
	if stale {
		m.chat.append(roleSystem, staleDraftNotice)
	}
	return m.setFocus(paneEditor)
}

// startRun saves the buffer and runs the tests for the current challenge.
func (m *Model) startRun() tea.Cmd {
	if m.current.ID == "" {
		return nil
	}
	code := m.editor.Value()
	_ = m.deps.Store.Save(m.current.ID, code)
	m.chat.append(roleSystem, "▶ running tests…")
	m.pending++
	m.loadKind = "running tests"
	return runCmd(m.current.Lang, code, m.current)
}

// handleRunResult records progress, reports pass/fail, and chains tutor feedback.
func (m *Model) handleRunResult(msg runResultMsg) tea.Cmd {
	if msg.res.Passed {
		m.chat.append(roleOK, "✓ All tests passed!")
		_ = m.deps.Store.Clear(msg.ch.ID)
		if m.curriculum && m.currentTopicID != "" {
			m.deps.Progress.MarkTopicDone(m.currentTopicID)
			m.chat.append(roleSystem, "🎓 Topic complete! Press Ctrl-N for the next topic, or pick any topic on the left.")
		}
	} else {
		m.chat.append(roleFail, "✗ Tests failed")
		m.chat.append(roleSystem, failureSummary(msg.res))
		m.deps.Progress.MarkInProgress(msg.ch.ID)
	}
	m.deps.Progress.RecordAttempt(msg.ch.ID, msg.res.Passed)
	_ = m.deps.Progress.Save()
	m.rebuildSidebar()

	m.pending++
	m.loadKind = "tutor feedback"
	return feedbackCmd(m.deps.Tutor, msg.ch, msg.code, msg.res.Output, msg.res.Passed)
}

func failureSummary(res executor.Result) string {
	out := strings.TrimSpace(res.Output)
	if res.TimedOut {
		if out == "" {
			out = "execution timed out (possible infinite loop)"
		}
		return "Reason:\n" + out + "\n\nTry:\nCheck loops and recursion for a condition that always makes progress."
	}
	if out == "" {
		return "No output was captured. Check the challenge tests and rerun."
	}

	var b strings.Builder
	if line := likelyFailureLine(out); line != "" {
		b.WriteString("Failed:\n")
		b.WriteString(line)
		b.WriteString("\n\n")
	}
	b.WriteString("Output:\n")
	b.WriteString(trimFailureOutput(out, 14))
	return b.String()
}

func likelyFailureLine(out string) string {
	lines := strings.Split(out, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "assert ") {
			return line
		}
	}
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "--- FAIL:"),
			strings.Contains(line, ".go:"),
			strings.Contains(line, "AssertionError"),
			strings.Contains(line, "panic:"),
			strings.Contains(line, "Error:"):
			return line
		}
	}
	return ""
}

func trimFailureOutput(out string, maxLines int) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) <= maxLines {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[:maxLines], "\n") + "\n... (" + strconv.Itoa(len(lines)-maxLines) + " more lines)"
}

// nextChallenge advances to the next topic in curriculum mode, or asks the LLM
// for a fresh challenge on the current topic in custom mode.
func (m *Model) nextChallenge() tea.Cmd {
	if m.curriculum {
		topics := m.curr.Topics()
		for i, t := range topics {
			if t.ID == m.currentTopicID && i+1 < len(topics) {
				return m.startTopic(topics[i+1])
			}
		}
		m.chat.append(roleSystem, "🎉 That's the last topic at this level — revisit any topic on the left, or restart for a higher level.")
		return nil
	}
	if m.topic == "" {
		return nil
	}
	m.pending++
	m.loadKind = "new challenge"
	return challengeCmd(m.deps.Tutor, m.topic)
}

func (m *Model) quit() tea.Cmd {
	if m.current.ID != "" {
		_ = m.deps.Store.Save(m.current.ID, m.editor.Value())
	}
	_ = m.deps.Progress.Save()
	return tea.Quit
}

// configReloadMsg is delivered after the external config editor exits.
type configReloadMsg struct{ err error }

// openConfig suspends the TUI and opens the config file in the user's editor
// ($EDITOR / $VISUAL, falling back to vi), creating it first if needed.
func (m *Model) openConfig() tea.Cmd {
	path := m.deps.ConfigPath
	if path == "" {
		m.chat.append(roleSystem, "⚠ No config path is set.")
		return nil
	}
	if err := config.EnsureFile(path); err != nil {
		m.chat.append(roleSystem, "⚠ Could not create config: "+err.Error())
		return nil
	}
	ed := os.Getenv("EDITOR")
	if ed == "" {
		ed = os.Getenv("VISUAL")
	}
	if ed == "" {
		ed = "vi"
	}
	return tea.ExecProcess(exec.Command(ed, path), func(err error) tea.Msg {
		return configReloadMsg{err: err}
	})
}

// applyConfigReload re-reads the config after editing and live-applies what can
// change mid-session: the layout AND the AI provider — the tutor is rebuilt so
// a new key/model/base URL works immediately, without restarting the app.
func (m *Model) applyConfigReload(msg configReloadMsg) {
	if msg.err != nil {
		m.chat.append(roleSystem, "⚠ Config editor exited with an error: "+msg.err.Error())
		return
	}
	cfg, err := config.Load(m.deps.ConfigPath, m.deps.BaseDir)
	if err != nil {
		m.chat.append(roleSystem, "⚠ Config has an error and was not applied: "+err.Error())
		return
	}
	m.deps.Cfg = cfg
	m.horizontal = cfg.Horizontal()
	m.layout() // re-flow the panes for the (possibly new) layout

	// Rebuild the AI client from the new [ai] section, keeping the learner's level.
	m.deps.Tutor = tutor.New(cfg.AI)
	if m.level != "" {
		m.deps.Tutor.SetLevel(m.level)
	}
	ai := "AI: " + cfg.AI.Model + " @ " + cfg.AI.ResolveBaseURL()
	if m.deps.Tutor.Offline() {
		ai = "AI: OFFLINE — key required but not found (try `meari check`)"
	}
	m.flash("✓ Config reloaded — layout " + cfg.UI.Layout + "; " + ai)
}

// --- focus & ordering helpers ---

// focusDir moves focus one pane in direction d (-1 left, +1 right), clamping at
// the edges (Vim windows don't wrap).
func (m *Model) focusDir(d int) tea.Cmd {
	// A folded sidebar isn't a focus target, so the left edge becomes the editor.
	lo := int(paneSidebar)
	if m.sidebarCollapsed {
		lo = int(paneEditor)
	}
	n := int(m.focus) + d
	if n < lo {
		n = lo
	}
	if n > int(paneChat) {
		n = int(paneChat)
	}
	return m.setFocus(pane(n))
}

// setFocus moves keyboard focus, blurring the old input and focusing the new so
// exactly one component captures keys / blinks a cursor at a time.
func (m *Model) setFocus(p pane) tea.Cmd {
	m.editor.Blur()
	m.chat.blur()
	m.sidebar.focused = false
	m.focus = p
	switch p {
	case paneEditor:
		return m.editor.Focus()
	case paneChat:
		return m.chat.focus()
	case paneSidebar:
		m.sidebar.focused = true
	}
	return nil
}

func (m *Model) addOrder(id string) {
	for _, x := range m.order {
		if x == id {
			return
		}
	}
	m.order = append(m.order, id)
}

func (m *Model) rebuildSidebar() {
	if m.curriculum {
		var items []sidebarItem
		for _, mod := range m.curr.Modules {
			items = append(items, sidebarItem{title: mod.Name, header: true})
			for _, t := range mod.Topics {
				items = append(items, sidebarItem{
					id:     t.ID,
					title:  t.Title,
					status: m.deps.Progress.TopicStatus(t.ID),
					active: t.ID == m.currentTopicID,
				})
			}
		}
		m.sidebar.setItems(items)
		return
	}

	items := make([]sidebarItem, 0, len(m.order))
	for _, id := range m.order {
		title := id
		if ch, ok := m.challenges[id]; ok {
			title = firstLine(ch.Prompt)
		}
		status := ""
		if e, ok := m.deps.Progress.Challenges[id]; ok {
			status = e.Status
		}
		items = append(items, sidebarItem{id: id, title: title, status: status})
		items[len(items)-1].active = id == m.current.ID
	}
	m.sidebar.setItems(items)
}

// --- layout ---

func (m *Model) layout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	m.topicInput.Width = min(50, m.width-6)

	// Two bars (title + status) and a pane's top/bottom border (2 rows) eat height.
	m.contentH = m.height - 4
	if m.contentH < 1 {
		m.contentH = 1
	}

	if m.horizontal {
		// (sidebar) | (content on top, input on the bottom). Each visible column
		// costs 2 border columns: 4 with the sidebar, 2 when it's folded away.
		borders := 4
		if m.sidebarCollapsed {
			borders = 2
		}
		contentW := m.width - borders
		if contentW < 3 {
			contentW = 3
		}
		if m.sidebarCollapsed {
			m.sidebarW = 0
			m.rightW = clampMin(contentW, 20)
		} else {
			m.sidebarW = clampMin(contentW*m.deps.Cfg.SidebarPct(22)/100, 14)
			m.rightW = clampMin(contentW-m.sidebarW, 20)
		}
		// The two stacked boxes cost 4 border rows vs the sidebar's 2, so their
		// combined content height is contentH-2.
		rightContent := m.contentH - 2
		if rightContent < 2 {
			rightContent = 2
		}
		// In this stacked layout the editor sits below the chat, so :compact /
		// :wide trade height between them rather than width.
		chatPct := clampRange(m.deps.Cfg.ChatPct(55)-m.editorBias, 30, 85)
		m.chatH = clampMin(rightContent*chatPct/100, 3)
		m.editorH = clampMin(rightContent-m.chatH, 3)

		m.sidebar.setSize(m.sidebarW, m.contentH)
		m.chat.setSize(m.rightW, m.chatH)
		m.editor.SetSize(m.rightW, max(1, m.editorH-1-len(m.promptHeaderLines(m.rightW))))
		return
	}

	// vertical: (sidebar) | editor | chat. Each visible box eats 2 border
	// columns: 6 with the sidebar, 4 when it's folded away.
	borders := 6
	if m.sidebarCollapsed {
		borders = 4
	}
	contentW := m.width - borders
	if contentW < 3 {
		contentW = 3
	}
	// The configured split is the base; :compact / :wide shift it live.
	chatPct := clampRange(m.deps.Cfg.ChatPct(30)-m.editorBias, 15, 75)
	if m.sidebarCollapsed {
		m.sidebarW = 0
		m.chatW = clampMin(contentW*chatPct/100, 16)
		m.editorW = clampMin(contentW-m.chatW, 10)
	} else {
		m.sidebarW = clampMin(contentW*m.deps.Cfg.SidebarPct(22)/100, 12)
		m.chatW = clampMin(contentW*chatPct/100, 16)
		m.editorW = clampMin(contentW-m.sidebarW-m.chatW, 10)
	}

	m.sidebar.setSize(m.sidebarW, m.contentH)
	m.editor.SetSize(m.editorW, max(1, m.contentH-1-len(m.promptHeaderLines(m.editorW))))
	m.chat.setSize(m.chatW, m.contentH)
}

func (m Model) View() string {
	if m.width == 0 {
		return "starting…"
	}
	if m.width < 60 || m.height < 16 {
		return "Terminal too small — please enlarge to at least 60×16."
	}

	if m.phase == phaseSetup {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.setupView())
	}

	if m.overlay != overlayNone {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.overlayView())
	}

	var row string
	if m.horizontal {
		// sidebar on the left; content (chat) above the input (editor) on the right.
		ch := m.box(paneChat, m.rightW, m.chatH, m.chat.view())
		ed := m.box(paneEditor, m.rightW, m.editorH, m.editorPaneView(m.rightW))
		right := lipgloss.JoinVertical(lipgloss.Left, ch, ed)
		if m.sidebarCollapsed {
			row = right
		} else {
			sb := m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view())
			row = lipgloss.JoinHorizontal(lipgloss.Top, sb, right)
		}
	} else {
		ed := m.box(paneEditor, m.editorW, m.contentH, m.editorPaneView(m.editorW))
		ch := m.box(paneChat, m.chatW, m.contentH, m.chat.view())
		if m.sidebarCollapsed {
			row = lipgloss.JoinHorizontal(lipgloss.Top, ed, ch)
		} else {
			sb := m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view())
			row = lipgloss.JoinHorizontal(lipgloss.Top, sb, ed, ch)
		}
	}

	frame := lipgloss.JoinVertical(lipgloss.Left, m.titleView(), row, m.statusView())
	// Final guard: never emit more than the screen holds, so a too-tall pane
	// can't scroll the alt-screen and leave residue behind.
	return lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(frame)
}

// setupView renders the current wizard step as a centered card.
func (m Model) setupView() string {
	var title, body string

	switch m.setupStep {
	case stepResume:
		title = "Welcome back!"
		body = m.setupMenu()
	case stepLanguage:
		title = "Which language do you want to learn?"
		body = m.setupMenu()
	case stepPath:
		title = "How would you like to learn " + titleCase(m.lang) + "?"
		body = m.setupMenu()
	case stepTopic:
		title = "What topic do you want to learn?"
		body = m.topicInput.View()
	case stepLevel:
		title = "What's your experience level?"
		body = m.setupMenu()
	}

	hint := "↑/↓ or j/k to move · enter to choose · esc to go back · ctrl+c to quit"
	if m.setupStep == stepTopic {
		hint = "type, then enter · esc to go back · ctrl+c to quit"
	}

	card := titleBar.Render(" Meari ") + "\n\n" +
		lipgloss.NewStyle().Bold(true).Render(title) + "\n\n" +
		body + "\n\n" +
		hintStyle.Render(hint)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 3).
		Width(56).
		Render(card)
}

// setupMenu renders the current selection step's options with the cursor row
// highlighted.
func (m Model) setupMenu() string {
	opts := m.setupOptions()
	var b strings.Builder
	for i, opt := range opts {
		if i == m.setupCursor {
			b.WriteString(selectedRow.Render("▸ " + opt))
		} else {
			b.WriteString("  " + opt)
		}
		if i < len(opts)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// titleCase upper-cases the first letter for display (e.g. "python" -> "Python").
func titleCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return strings.ToUpper(string(r[0])) + string(r[1:])
}

// box renders one pane's content inside a focus-aware border, hard-clamped to
// exactly w×h (plus the 1-cell border) so it can never overflow the layout
// regardless of how many lines the sub-view produced.
func (m Model) box(p pane, w, h int, content string) string {
	return borderStyle(m.focus == p).
		Width(w).Height(h).
		MaxWidth(w + 2).MaxHeight(h + 2).
		Render(content)
}

func (m Model) editorPaneView(w int) string {
	out := m.challengeHeader(w)
	for _, ln := range m.promptHeaderLines(w) {
		out += "\n" + promptHeaderStyle.MaxWidth(w).Render(ln)
	}
	return out + "\n" + m.editor.View()
}

// challengeHeader labels the editor pane with the language and a SHORT title.
// The full problem statement lives at the top of the editor as a comment — a
// long prompt squeezed onto one truncated line here told the learner nothing.
func (m Model) challengeHeader(w int) string {
	label := "No challenge selected"
	if m.current.ID != "" {
		label = strings.ToUpper(challengeLang(m.current)) + " · " + m.shortTitle()
	}
	if w > 0 && lipgloss.Width(label) > w-2 {
		label = truncate(label, max(1, w-2))
	}
	return editorHeader.Width(w).Render(label)
}

// shortTitle is a compact name for the current challenge: the curriculum
// topic's title, or the challenge id prettified ("sum-list" -> "Sum list").
func (m Model) shortTitle() string {
	if m.curriculum {
		if t, ok := m.topicByID[m.currentTopicID]; ok && t.Title != "" {
			return t.Title
		}
	}
	if title := prettyID(m.current.ID); title != "" {
		return title
	}
	return m.topic
}

// prettyID turns a kebab/snake challenge id into a readable title.
func prettyID(id string) string {
	s := strings.NewReplacer("-", " ", "_", " ").Replace(id)
	return titleCase(strings.TrimSpace(s))
}

func (m Model) titleView() string {
	t := "Meari"
	if m.curriculum {
		t += " · " + m.curr.Lang + "/" + m.curr.Level
	}
	if m.topic != "" {
		t += " — " + m.topic
	}
	if m.deps.Tutor.Offline() {
		t += "  (offline)"
	}
	return titleBar.Width(m.width).Render(t)
}

func (m Model) statusView() string {
	if m.cmdMode {
		return statusBar.Width(m.width).Render(m.cmdLine.View())
	}
	left := "[" + m.focusName() + "]"
	if m.pending > 0 {
		left += " " + m.spin.View() + " " + m.loadKind
	} else if m.err != nil {
		left += " " + errStyle.Render("error: "+m.err.Error())
	}
	if m.notice != "" && time.Since(m.noticeAt) < noticeTTL {
		return statusBar.Width(m.width).Render(left + "   " + noticeStyle.Render(m.notice))
	}
	hints := "⌃w h·l focus · : cmds (:help) · ⌃s/:submit run · u/⌃r undo/redo · ⌃c quit"
	switch {
	case m.pendingWindow:
		hints = errStyle.Render("⌃w") + " window: h/l choose pane"
	case m.focus == paneChat:
		hints = "enter send · ⌥o/:copy copy reply · ⌃f/⌃b page · ⌃d/⌃u half · wheel scrolls"
	case m.focus == paneSidebar:
		hints = "j/k move · enter open · : cmds (:help) · ⌃r run · ⌃c quit"
	}
	line := left + "   " + hintStyle.Render(hints)
	return statusBar.Width(m.width).Render(line)
}

func (m Model) focusName() string {
	switch m.focus {
	case paneSidebar:
		if m.curriculum {
			return "curriculum"
		}
		return "challenges"
	case paneEditor:
		return "editor"
	case paneChat:
		return "chat"
	}
	return ""
}

// --- draft staleness ---

var (
	pyDefNameRE = regexp.MustCompile(`(?m)^\s*def\s+([A-Za-z_]\w*)\s*\(`)
	goDefNameRE = regexp.MustCompile(`(?m)^\s*(?:func\s+(?:\([^)]*\)\s*)?|type\s+)([A-Za-z_]\w*)`)
)

// starterNames extracts the function/type names a challenge's starter defines.
func starterNames(ch tutor.Challenge) []string {
	re := pyDefNameRE
	if strings.EqualFold(ch.Lang, "go") || strings.EqualFold(ch.Lang, "golang") {
		re = goDefNameRE
	}
	var names []string
	for _, m := range re.FindAllStringSubmatch(ch.StarterCode, -1) {
		names = append(names, m[1])
	}
	return names
}

// draftMatches reports whether a saved draft plausibly belongs to ch. Challenge
// content can change between versions of the app while drafts (keyed only by
// id) survive on disk; a draft that mentions none of the starter's function
// names would silently shadow the real challenge with an unrelated problem.
// Prose challenges (no extractable names) accept any draft.
func draftMatches(ch tutor.Challenge, draft string) bool {
	names := starterNames(ch)
	if len(names) == 0 {
		return true
	}
	for _, n := range names {
		if strings.Contains(draft, n) {
			return true
		}
	}
	return false
}

// loadStarterOrDraft picks the editor contents for ch: its saved draft when one
// exists and still matches the challenge, else the fresh starter. The problem
// statement is NOT embedded in the buffer — it renders as a pinned header above
// the editor (see promptHeaderLines), wrapped to the live pane width, so it can
// never scramble the textarea's own wrapping or line numbers. stale is true
// when a draft existed but belonged to an older version of the challenge.
func (m *Model) loadStarterOrDraft(ch tutor.Challenge) (code string, stale bool) {
	code = ch.StarterCode
	d, ok := m.deps.Store.Load(ch.ID)
	if !ok {
		return code, false
	}
	if !draftMatches(ch, d) {
		return code, true
	}
	return d, false
}

// maxPromptHeaderLines caps how much pane height the pinned statement may take.
const maxPromptHeaderLines = 6

// promptHeaderLines wraps the current challenge's statement to the pane width
// as a pinned description block above the numbered editor buffer.
func (m Model) promptHeaderLines(w int) []string {
	if m.current.ID == "" || strings.TrimSpace(m.current.Prompt) == "" {
		return nil
	}
	marker := "Challenge: "
	avail := w - len(marker)
	if avail < 8 {
		avail = 8
	}
	lines := wrapWords(m.current.Prompt, avail)
	if len(lines) > maxPromptHeaderLines {
		lines = lines[:maxPromptHeaderLines]
		lines[maxPromptHeaderLines-1] += " …"
	}
	for i := range lines {
		if i == 0 {
			lines[i] = marker + lines[i]
			continue
		}
		lines[i] = strings.Repeat(" ", len(marker)) + lines[i]
	}
	return lines
}

// wrapWords greedily wraps s into lines of at most width characters.
func wrapWords(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		if len(cur)+1+len(w) > width {
			lines = append(lines, cur)
			cur = w
		} else {
			cur += " " + w
		}
	}
	return append(lines, cur)
}

const staleDraftNotice = "⚠ Your saved draft was for an older version of this challenge, " +
	"so the editor shows the fresh starter. (The old draft is replaced when you save.)"

// --- small utilities ---

// challengeLang is the language a challenge's unlabeled code fences should be
// highlighted as ("" historically means Python).
func challengeLang(ch tutor.Challenge) string {
	if ch.Lang == "" {
		return "python"
	}
	return ch.Lang
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return s
}

func clampMin(v, lo int) int {
	if v < lo {
		return lo
	}
	return v
}

func clampRange(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
