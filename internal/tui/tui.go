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
	Svc        *core.Service // the vault engine, for vault-built courses
	Cfg        config.Config
	Topic      string // from --topic; empty means prompt for it at startup
	Curriculum bool   // from -tutor; start in guided curriculum mode
	ConfigPath string // path to config.toml, for the :config command
	BaseDir    string // working dir, for reloading config after :config
}

type pane int

const (
	paneSidebar pane = iota
	paneEditor
	paneChat
	// paneLesson is the read-only lecture pane shown on lesson rows between
	// the sidebar and the chat (the editor's geometric slot). Focus order is
	// visual, via visiblePanes — the constant's value carries no position.
	paneLesson
)

type phase int

const (
	phaseSetup phase = iota // the guided startup wizard
	phaseReady              // the three panes
)

// setupStep is the current screen of the startup flow.
type setupStep int

const (
	stepDashboard setupStep = iota // the launch dashboard: pick a course or destination
	stepTopic                      // type the specific topic
	stepLevel                      // beginner / intermediate / advanced
)

// dashKind says what selecting a launch-dashboard row does.
type dashKind int

const (
	dashContinue dashKind = iota // resume the saved session
	dashCourse                   // study a course (seeded or learner-built)
	dashTopic                    // ask the AI for a custom topic
	dashVault                    // switch to the notes vault
)

// dashEntry is one selectable row of the launch dashboard.
type dashEntry struct {
	kind    dashKind
	id      string // the course id / language, for dashCourse and dashLang
	section string // group header, drawn above the section's first entry
	title   string
	meta    string // dim right-aligned detail (level, progress, what happens)
}

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
	exit   SwitchTarget // set by :vault so Run can report a mode switch to the shell

	sidebar    sidebarModel
	editor     editor.Model
	chat       chatModel
	lesson     chatModel // read-only lecture pane (noInput); shown on lesson rows
	topicInput textinput.Model

	// Startup flow state (phaseSetup).
	setupStep    setupStep
	setupCursor  int
	setupHistory []setupStep // for Esc = back
	dash         []dashEntry // the launch dashboard's rows
	dashCursor   int         // dashboard position, restored when Esc returns to it
	lang         string      // the running course's id (curriculum.Lang)
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
	streaming      bool
	streamStopping bool
	streamCancel   context.CancelFunc
	streamCh       chan streamChunkMsg

	// Curriculum mode: the left pane is the guided learning path instead of the
	// generated-challenge list. Lessons and challenges are pre-authored (no LLM).
	curriculum       bool
	curr             curriculum.Curriculum
	topicByID        map[string]curriculum.Topic
	currentTopicID   string
	currentTopicView string // "lesson" or "quiz"; mirrors the selected course tree row

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

	// view selects the screen (":view"): "auto" follows the topic kind — essay
	// topics hide the editor so the lesson, conversation, and answer all live
	// in the chat pane; "chat"/"code" force it either way.
	view string

	// chatCollapsed hides the chat pane (":chat"); the freed width goes to the
	// editor (or the lesson pane on lecture rows). Ignored wherever the chat
	// is the only content surface — see chatHidden.
	chatCollapsed bool

	// pendingWindow is set after Ctrl-W; the next key (h/j/k/l) picks a pane,
	// Vim window-command style.
	pendingWindow bool

	// Chat drag-selection state: a left press on the chat anchors a selection;
	// motion with the button held sweeps it out (Alt-C copies).
	dragChat bool
	// dragEditor mirrors dragChat for the editor pane, so a drag over the editor
	// sweeps out (and on release copies) its text.
	dragEditor bool
	// dragLesson mirrors dragChat for the lesson pane.
	dragLesson bool

	// pendingLeader is set after "," in the editor's Normal mode, in the
	// sidebar, or in the chat's Vim Normal mode; it starts a leader chord.
	// "n" folds the sidebar; from the editor any other key replays the
	// swallowed "," so its repeat-find binding still works.
	pendingLeader bool

	// Global ex-command line (":topic", ":clear", ":progress"). cmdMode shows the
	// cmdLine input in the status row; it's opened with ":" from the sidebar and
	// also driven by RunCommandMsg forwarded from the editor's own command line.
	cmdMode bool
	cmdLine textinput.Model
	cmdHist editor.CmdHistory
	cmdComp editor.CmdCompleter // Tab completion over tutorExCmds

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
	lessonW  int // lesson pane width; 0 whenever the lesson split is off
	chatW    int
	contentH int
	// horizontal-layout extras: the right column width and the stacked heights.
	rightW  int
	chatH   int
	editorH int
}

// Run constructs the model and runs the full-screen program.
// SwitchTarget is how a TUI tells the shell loop (main.runShell) what to do when
// it exits: quit the process, or hand off to the other TUI.
type SwitchTarget int

const (
	StayQuit SwitchTarget = iota
	SwitchToVault
	SwitchToTutor
)

// Outcome is returned by Run/RunVault: where to go next, plus enough of the
// tutor's session (Topic/Curriculum) to resume it without the setup wizard.
type Outcome struct {
	Target     SwitchTarget
	Topic      string
	Curriculum bool
}

func Run(d Deps) (Outcome, error) {
	enableTUIColor()
	loadTheme(d.Cfg.DataDir)
	m := newModel(d)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	final, err := p.Run()
	out := Outcome{}
	if fm, ok := final.(Model); ok {
		// A curriculum session resumes via Curriculum (which restores the exact
		// topic from saved progress); only a custom-topic session resumes via
		// Topic. In curriculum mode fm.topic holds the current challenge title,
		// so it must NOT be used as a resume topic.
		out = Outcome{Target: fm.exit, Curriculum: fm.curriculum}
		if !fm.curriculum {
			out.Topic = fm.topic
		}
	}
	return out, err
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
		view:             d.Cfg.UI.View,
		sidebar:          newSidebar(),
		editor: editor.New("", d.Cfg.VimEditor(), save).
			WithGlobalCmds(tutorExCmds).
			WithArgCompleter(func(input string) []string {
				return topicArgCandidates(d.Svc, input)
			}),
		chat:       newChat(),
		lesson:     newLessonPane(),
		topicInput: ti,
		cmdLine:    cl,
		curID:      curID,
		challenges: map[string]tutor.Challenge{},
		chatByKey:  map[string][]chatBlock{},
		histByKey:  map[string][]tutor.ChatTurn{},
		spin:       sp,
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
		// Curriculum via flag: resume the saved session if any, else the
		// seeded Go beginner course.
		key, level, startID := "go-beginner", curriculum.Beginner, ""
		if last := d.Progress.Last; last != nil {
			key, level, startID = last.Lang, last.Level, last.TopicID
		}
		m.loadCurriculum(key, level, startID)
		m.phase = phaseReady
	default:
		// The launch dashboard: every way into the app on one screen.
		m.setupStep = stepDashboard
		m.dash = m.dashboardEntries()
	}
	return m
}

// dashboardEntries assembles the launch dashboard: the saved session, every
// course (the seeded Go track and learner-built ones alike), and the two
// destinations that aren't a course (a custom AI topic, the vault itself).
func (m Model) dashboardEntries() []dashEntry {
	var out []dashEntry
	var metas []core.CourseMeta
	if m.deps.Svc != nil {
		metas, _ = m.deps.Svc.ListCourses()
	}
	if last := m.deps.Progress.Last; last != nil && m.resumable(last.Lang, metas) {
		out = append(out, dashEntry{
			kind:    dashContinue,
			section: "Continue",
			title:   last.Title,
			meta:    "pick up where you left off",
		})
	}
	for _, c := range metas {
		level := c.Level
		if level == "" {
			level = curriculum.Beginner
		}
		meta := level
		if done := m.courseTopicsDone(c.ID); done > 0 {
			meta = level + " · " + itoa(done) + " done"
		}
		out = append(out, dashEntry{
			kind:    dashCourse,
			id:      c.ID,
			section: "Courses",
			title:   c.Title,
			meta:    meta,
		})
	}
	return append(out,
		dashEntry{kind: dashTopic, section: "Or", title: "A topic of my own",
			meta: "the AI writes a lesson + challenge"},
		dashEntry{kind: dashVault, section: "Or", title: "Open the vault",
			meta: "notes, lessons, build courses"},
	)
}

// resumable reports whether the saved session's course still exists, so the
// dashboard never leads with a dead Continue row (e.g. after a course was
// deleted, or for sessions from removed built-ins). With no vault service the
// check is skipped.
func (m Model) resumable(key string, metas []core.CourseMeta) bool {
	if m.deps.Svc == nil {
		return true
	}
	for _, c := range metas {
		if strings.EqualFold(c.ID, key) || strings.EqualFold(c.Title, key) {
			return true
		}
	}
	return false
}

// courseTopicsDone counts a course's completed topics from saved progress
// (topic ids are "course-<id>-…"), without loading the course itself.
func (m Model) courseTopicsDone(courseID string) int {
	done := 0
	for tid, status := range m.deps.Progress.Topics {
		if status == "done" && strings.HasPrefix(tid, "course-"+courseID+"-") {
			done++
		}
	}
	return done
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
		if m.lessonSplit() {
			// Async-generated lectures land in the lesson pane like any other.
			m.lesson.setLesson(msg.text)
		} else {
			m.chat.append(roleLesson, msg.text)
			// append() pins to the tail, so a long lesson would otherwise open at
			// its very end — start the reader at the beginning.
			m.chat.gotoTop()
		}
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
		// Re-issue focus through setFocus: a layout switch may have removed
		// the pane that held it (e.g. the lesson pane under horizontal).
		return m, m.setFocus(m.focus)

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
	if msg.Type == tea.KeyEsc && m.streaming {
		return m.stopStream()
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
		// The ",n" fold chord works from the course list too — same keys as
		// the editor's Normal mode (mirrors the vault TUI's sidebar).
		if leader {
			if msg.String() == "n" {
				return m.cmdFold()
			}
			return m, nil
		}
		switch msg.String() {
		case ":":
			return m.openCmdLine()
		case ",":
			m.pendingLeader = true
			return m, nil
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

	case paneLesson:
		// The lecture pane is never a typing surface, so bare keys are safe —
		// same leader/command treatment as the sidebar.
		if leader {
			if msg.String() == "n" {
				return m.cmdFold()
			}
			return m, nil
		}
		switch msg.String() {
		case ":":
			return m.openCmdLine()
		case ",":
			m.pendingLeader = true
			return m, nil
		case "alt+c", "ç", "Ç":
			m.flash(copySelection(&m.lesson))
			return m, nil
		case "alt+o", "ø", "Ø":
			m.flash(copyChat(&m.lesson, "all")) // the whole lecture document
			return m, nil
		}
		var cmd tea.Cmd
		m.lesson, cmd = m.lesson.Update(msg)
		return m, cmd

	case paneChat:
		switch msg.String() {
		// Copy the tutor's last reply: Alt+O on Linux; on macOS, Option+O —
		// which arrives as "ø"/"Ø" unless the terminal sends Option as Meta.
		// (Cmd+O never reaches a terminal app; the emulator consumes it.)
		case "alt+o", "ø", "Ø":
			m.flash(copyChat(&m.chat, ""))
			return m, nil
		// Copy the mouse drag-selection: Alt+C / Option+C (macOS sends "ç").
		case "alt+c", "ç", "Ç":
			m.flash(copySelection(&m.chat))
			return m, nil
		// Paste the system clipboard into the chat input: Alt+V / Option+V
		// (macOS sends "√" for Option+V). Cmd+V also works — the terminal
		// delivers it as a bracketed paste straight into the input.
		case "alt+v", "√":
			m.flash(pasteChat(&m.chat))
			return m, nil
		}
		// The ",n" fold chord works from the chat's Vim Normal mode too (never
		// from Insert, where "," must type a comma).
		if leader {
			if msg.String() == "n" {
				return m.cmdFold()
			}
			return m, nil
		}
		if m.chat.inNormal() && msg.String() == "," {
			m.pendingLeader = true
			return m, nil
		}
		// The input's Vim Normal mode (Esc): ":" opens the command line right
		// from the chat; Enter doesn't send while in it.
		if m.chat.inNormal() && msg.String() == ":" {
			return m.openCmdLine()
		}
		if msg.Type == tea.KeyEnter && !m.chat.inNormal() {
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

	// Clicking the title-bar "Check answer" button runs the tests (same as Ctrl-S
	// / :submit). Hit-test it before the pane routing, since it lives on row 0.
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		if x0, x1, ok := m.checkButtonBounds(); ok && msg.Y == 0 && msg.X >= x0 && msg.X < x1 {
			return m, m.startRun()
		}
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
		case paneLesson:
			var cmd tea.Cmd
			m.lesson, cmd = m.lesson.Update(msg)
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

	// Left click: focus the pane under the cursor; on the chat it also anchors
	// a text SELECTION, so dragging sweeps out transcript text. RELEASING the
	// drag copies the selection to the clipboard automatically (Alt-C also
	// copies, but many Linux terminals eat Alt-<key> as a menu mnemonic, so
	// release-to-copy is the reliable path). Scrolling stays on the wheel and
	// Ctrl-F/B; the terminal's native bypass still works too — Option+drag on
	// macOS, Shift+drag on Linux skip mouse reporting entirely.
	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button == tea.MouseButtonLeft {
			if p == paneSidebar && !m.sidebarCollapsed {
				// Click a curriculum row to select and load it, mirroring Enter.
				// Content rows begin two cells down: the title bar (row 0) and the
				// box's top border (row 1).
				m.chat.clearSelect()
				m.setFocus(paneSidebar)
				if m.sidebar.clickRow(msg.Y - 2) {
					return m, m.openSelected()
				}
				return m, nil
			}
			m.chat.clearSelect()
			m.lesson.clearSelect()
			m.dragChat, m.dragEditor, m.dragLesson = false, false, false
			switch p {
			case paneChat:
				m.dragChat = true
				lx, ly := m.chatLocal(msg.X, msg.Y)
				m.chat.startSelect(lx, ly)
			case paneLesson:
				m.dragLesson = true
				lx, ly := m.lessonLocal(msg.X, msg.Y)
				m.lesson.startSelect(lx, ly)
			case paneEditor:
				lx, ly := m.editorLocal(msg.X, msg.Y)
				m.dragEditor = m.editor.MouseSelectStart(lx, ly)
			}
			return m, m.setFocus(p)
		}
	case tea.MouseActionMotion:
		if msg.Button == tea.MouseButtonLeft {
			if m.dragChat {
				lx, ly := m.chatLocal(msg.X, msg.Y)
				m.chat.dragSelect(lx, ly)
				return m, nil
			}
			if m.dragLesson {
				lx, ly := m.lessonLocal(msg.X, msg.Y)
				m.lesson.dragSelect(lx, ly)
				return m, nil
			}
			if m.dragEditor {
				lx, ly := m.editorLocal(msg.X, msg.Y)
				m.editor.MouseSelectTo(lx, ly)
				return m, nil
			}
		}
	case tea.MouseActionRelease:
		if m.dragChat && m.chat.sel.active { // a real drag, not a bare click
			m.flash(copySelection(&m.chat))
		}
		if m.dragLesson && m.lesson.sel.active {
			m.flash(copySelection(&m.lesson))
		}
		if m.dragEditor {
			if text := m.editor.MouseSelectEnd(); text != "" {
				m.flash(fmt.Sprintf("✓ copied selection (%d chars)", len([]rune(text))))
			}
		}
		m.dragChat, m.dragEditor, m.dragLesson = false, false, false
	}
	return m, nil
}

// chatLocal converts a terminal cell to chat-viewport-local coordinates: past
// the title row, each box's border cell, and the panes left of the chat. In
// the horizontal layout the chat is the top box of the right column.
func (m Model) chatLocal(x, y int) (int, int) {
	sidebarSpan := m.sidebarW + 2
	if m.sidebarCollapsed {
		sidebarSpan = 0
	}
	if m.horizontal {
		return x - (sidebarSpan + 1), y - 2
	}
	editorSpan := m.editorW + 2
	if m.editorW == 0 {
		editorSpan = 0
	}
	lessonSpan := m.lessonW + 2
	if m.lessonW == 0 {
		lessonSpan = 0
	}
	return x - (sidebarSpan + lessonSpan + editorSpan + 1), y - 2
}

// lessonLocal converts a terminal cell to lesson-viewport-local coordinates:
// the lesson pane is the first box right of the sidebar (vertical layout only
// — the split never renders under the horizontal layout).
func (m Model) lessonLocal(x, y int) (int, int) {
	sidebarSpan := m.sidebarW + 2
	if m.sidebarCollapsed {
		sidebarSpan = 0
	}
	return x - (sidebarSpan + 1), y - 2
}

// editorLocal converts a terminal cell to editor-textarea-local coordinates: past
// the title row, the box border, the panes above/left of the editor, and the
// editor pane's own header (the language label plus any pinned prompt lines). x
// lands inside the editor box; the line-number gutter is subtracted by the editor.
func (m Model) editorLocal(x, y int) (int, int) {
	sidebarSpan := m.sidebarW + 2
	if m.sidebarCollapsed {
		sidebarSpan = 0
	}
	if m.horizontal {
		headerRows := 1 + len(m.promptHeaderLines(m.rightW))
		return x - (sidebarSpan + 1), y - (m.chatH + 4 + headerRows)
	}
	headerRows := 1 + len(m.promptHeaderLines(m.editorW))
	return x - (sidebarSpan + 1), y - (2 + headerRows)
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
		// sidebar on the left; chat stacked over the editor on the right (the
		// editor row is absent in the chat-centric view, the chat row under
		// ":chat").
		if x < sidebarSpan {
			return paneSidebar, true
		}
		if m.chatH == 0 {
			return paneEditor, true
		}
		if m.editorH == 0 || y < 1+m.chatH+2 {
			return paneChat, true
		}
		return paneEditor, true
	}
	// vertical: sidebar | lesson-or-editor | chat, each box +2 columns for its
	// border; a hidden pane occupies no columns at all (the lesson and editor
	// spans are never both nonzero — one fills the middle slot).
	editorSpan := m.editorW + 2
	if m.editorW == 0 {
		editorSpan = 0
	}
	lessonSpan := m.lessonW + 2
	if m.lessonW == 0 {
		lessonSpan = 0
	}
	switch {
	case x < sidebarSpan:
		return paneSidebar, true
	case x < sidebarSpan+lessonSpan:
		return paneLesson, true
	case x < sidebarSpan+lessonSpan+editorSpan || m.chatW == 0:
		// A ":chat"-hidden pane occupies no columns; the middle surface owns
		// the rest of the row.
		if m.editorW == 0 && m.lessonW > 0 {
			return paneLesson, true
		}
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
// still quits, Tab/Shift-Tab cycle command completions, anything else edits
// the text.
func (m Model) updateCmdLine(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type != tea.KeyTab && msg.Type != tea.KeyShiftTab {
		m.cmdComp.Reset() // any other key ends the completion cycle
	}
	switch msg.Type {
	case tea.KeyTab, tea.KeyShiftTab:
		dir := 1
		if msg.Type == tea.KeyShiftTab {
			dir = -1
		}
		if s, ok := m.cmdComp.Next(m.cmdLine.Value(), m.exCandidates(), dir); ok {
			m.cmdLine.SetValue(s)
			m.cmdLine.CursorEnd()
		}
		return m, nil
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

// tutorExCmds lists every command runEx accepts (aliases included), sorted,
// for Tab completion in the command prompt.
var tutorExCmds = []string{
	"answer", "chat", "clear", "compact", "config", "copy", "course", "export", "fold",
	"help", "notes", "paste", "progress", "q", "quit", "run", "sidebar",
	"subject", "submit", "theme", "topic", "vault", "view", "wide", "yank",
}

// exCandidates returns the command line's Tab-completion candidates: command
// names — or, once a course verb is typed, its argument completed against the
// built-in languages and vault courses (":topic nos⇥" → the full course id).
func (m Model) exCandidates() []string {
	if c := topicArgCandidates(m.deps.Svc, m.cmdLine.Value()); c != nil {
		return c
	}
	return tutorExCmds
}

// topicArgCandidates completes a course verb's argument: for "topic …" (and
// aliases) it returns "topic <id>" for every built-in language and vault
// course, nil for any other input. Shared by the global prompt and the
// editor's ":" line (via WithArgCompleter).
func topicArgCandidates(svc *core.Service, input string) []string {
	if strings.HasPrefix(input, "view ") {
		return []string{"view auto", "view chat", "view code"}
	}
	if c := themeArgCandidates(input); c != nil {
		return c
	}
	for _, verb := range []string{"topic ", "course ", "subject "} {
		if !strings.HasPrefix(input, verb) {
			continue
		}
		var ids []string
		if svc != nil {
			if metas, err := svc.ListCourses(); err == nil {
				for _, c := range metas {
					ids = append(ids, c.ID)
				}
			}
		}
		out := make([]string, len(ids))
		for i, id := range ids {
			out[i] = verb + id
		}
		return out
	}
	return nil
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
	case "chat":
		return m.cmdChat()
	case "theme":
		m.flash(themeCommand(m.deps.Cfg.DataDir, strings.Join(fields[1:], " ")))
		return m, nil
	case "view":
		return m.cmdView(fields[1:])
	case "compact":
		return m.cmdResizeEditor(-editorBiasStep)
	case "wide":
		return m.cmdResizeEditor(editorBiasStep)
	case "answer":
		return m.cmdAnswer()
	case "submit", "run":
		// Submission without the Ctrl-S chord: works from any pane (the
		// editor's own ":submit" already handles the editing case).
		return m, m.startRun()
	case "copy", "yank":
		what := ""
		if len(fields) > 1 {
			what = fields[1]
		}
		m.flash(copyChat(&m.chat, what))
		return m, nil
	case "export":
		m.flash(exportChat(&m.chat, m.deps.Cfg.ExportsDir, m.topic))
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
	case "q", "quit":
		return m, m.quit() // saves drafts/progress on the way out
	case "vault", "notes":
		m.exit = SwitchToVault
		return m, m.quit() // saves the draft, then the shell loop opens the vault
	case "learn", "essay", "grade":
		m.flash(":" + fields[0] + " lives in the vault — type :vault to switch to it")
		return m, nil
	default:
		m.flash("unknown command: :" + raw + "  (try :help)")
		return m, nil
	}
}

// cmdTopic switches to another course. With no argument it opens the picker;
// multi-word names are taken whole (":topic nosql databases").
func (m Model) cmdTopic(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		m.overlay = overlayPicker
		m.pickerCursor = 0
		for i, id := range m.courseNames() {
			if id == m.lang {
				m.pickerCursor = i
				break
			}
		}
		return m, nil
	}
	return m.switchCourse(strings.ToLower(strings.Join(args, " ")))
}

// switchCourse enters a course (by id or title; all courses — the seeded Go
// track included — are markdown meari-courses); unknown names list what exists.
func (m Model) switchCourse(course string) (tea.Model, tea.Cmd) {
	if vc, err := m.loadVaultCourse(course); err == nil {
		m.chat.append(roleSystem, "— now studying "+vc.Lang+" ("+vc.Level+") —")
		return m, m.installCurriculum(vc, "")
	}
	// Generated course ids are long ("nosql-databases-a-beginner-s-…"): a
	// unique case-insensitive substring of the id or title is enough.
	if m.deps.Svc != nil {
		if metas, err := m.deps.Svc.ListCourses(); err == nil {
			q := strings.ToLower(course)
			var hits []core.CourseMeta
			for _, c := range metas {
				if strings.Contains(strings.ToLower(c.ID), q) ||
					strings.Contains(strings.ToLower(c.Title), q) {
					hits = append(hits, c)
				}
			}
			if len(hits) == 1 {
				if vc, err := m.loadVaultCourse(hits[0].ID); err == nil {
					m.chat.append(roleSystem, "— now studying "+hits[0].Title+" ("+vc.Level+") —")
					return m, m.installCurriculum(vc, "")
				}
			}
			if len(hits) > 1 {
				var names []string
				for _, c := range hits {
					names = append(names, c.ID)
				}
				m.flash("\"" + course + "\" matches several courses: " + strings.Join(names, ", "))
				return m, nil
			}
		}
	}
	m.flash("no course \"" + course + "\" — try: " + strings.Join(m.courseNames(), ", "))
	return m, nil
}

// courseNames lists everything :topic accepts: the vault's meari-courses.
func (m Model) courseNames() []string {
	ids, _ := m.pickerEntries()
	return ids
}

// cmdView switches the screen: "chat" hides the editor (the lesson,
// conversation, and your answer all live in the chat pane), "code" forces the
// three-pane screen, "auto" (default) follows the topic kind.
func (m Model) cmdView(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		eff := "code"
		if m.chatCentric() {
			eff = "chat"
		}
		m.flash("view: " + m.view + " (" + eff + " here) — :view auto|chat|code")
		return m, nil
	}
	v := strings.ToLower(args[0])
	if v != "auto" && v != "chat" && v != "code" {
		m.flash("usage: :view auto|chat|code")
		return m, nil
	}
	m.view = v
	m.syncLessonPresentation() // the lecture moves between pane and transcript
	m.layout()
	m.flash("view: " + v)
	// Re-issue focus through setFocus: its retarget rules land focus on a
	// visible pane whichever surface (editor, lesson pane) just vanished.
	return m, m.setFocus(m.focus)
}

// cmdChat toggles the chat pane (":chat"); the editor — or the lesson pane on
// lecture rows — absorbs the freed width. Refused where the chat is the only
// content surface, since hiding it would blank the screen. Hiding the focused
// chat moves focus left so keys never vanish into a hidden pane; anything
// that lands focus in the chat later unfolds it again via setFocus.
func (m Model) cmdChat() (tea.Model, tea.Cmd) {
	if m.chatCentric() && !m.lessonSplit() {
		m.flash("the chat is the whole screen here — :view code first")
		return m, nil
	}
	m.chatCollapsed = !m.chatCollapsed
	var cmd tea.Cmd
	if m.chatCollapsed && m.focus == paneChat {
		cmd = m.setFocus(paneEditor) // retargets to the lesson pane on lecture rows
	}
	m.layout()
	if m.chatCollapsed {
		m.flash("Chat pane hidden — :chat to bring it back")
	}
	return m, cmd
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

// cmdFold toggles the left tree pane. Folding a focused tree moves focus to the
// editor; unfolding focuses the newly visible tree so navigation can begin
// immediately. The layout is re-flowed as the pane appears or disappears.
func (m Model) cmdFold() (tea.Model, tea.Cmd) {
	m.sidebarCollapsed = !m.sidebarCollapsed
	var cmd tea.Cmd
	if !m.sidebarCollapsed {
		cmd = m.setFocus(paneSidebar)
	} else if m.focus == paneSidebar {
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
		ids, _ := m.pickerEntries()
		switch msg.String() {
		case "esc", "q":
			m.overlay = overlayNone
		case "j", "down":
			if m.pickerCursor < len(ids)-1 {
				m.pickerCursor++
			}
		case "k", "up":
			if m.pickerCursor > 0 {
				m.pickerCursor--
			}
		case "enter":
			m.overlay = overlayNone
			return m.switchCourse(ids[m.pickerCursor])
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

// pickerEntries returns everything the course picker offers — the vault's
// meari-courses (the seeded Go track included). ids are what :topic accepts;
// labels are for display.
func (m Model) pickerEntries() (ids, labels []string) {
	if m.deps.Svc != nil {
		if metas, err := m.deps.Svc.ListCourses(); err == nil {
			for _, c := range metas {
				ids = append(ids, c.ID)
				labels = append(labels, c.Title)
			}
		}
	}
	return ids, labels
}

func (m Model) pickerView() string {
	_, labels := m.pickerEntries()
	var b strings.Builder
	for i, l := range labels {
		if i == m.pickerCursor {
			b.WriteString(selectedRow.Render("▸ " + l))
		} else {
			b.WriteString("  " + l)
		}
		if i < len(labels)-1 {
			b.WriteString("\n")
		}
	}
	return modalCard("Switch course", b.String(), "↑/↓ or j/k · enter to switch · esc to cancel")
}

func (m Model) progressView() string {
	var b strings.Builder
	ids, labels := m.pickerEntries()
	if len(ids) == 0 {
		b.WriteString(hintStyle.Render("no courses yet — :vault then :course builds one"))
		b.WriteString("\n")
	}
	for i, id := range ids {
		done, total := m.courseCompletion(id)
		b.WriteString(courseProgressLine(labels[i], done, total))
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
		"  :topic <course>    switch course (any meari-course, by id or title)",
		"  :subject <course>  alias for :topic",
		"  :fold              fold/unfold the left tree pane (also ,n)",
		"  :view auto|chat|code  lesson|chat split · single full-width chat · editor",
		"  :chat              hide/show the chat pane (editor takes the width)",
		"  :compact / :wide   shrink/grow the editor (frees chat space)",
		"  :theme [<name>]    switch color theme (no name lists them)",
		"  :answer            reveal a model solution for the open challenge",
		"  :vault             switch to the notes vault (Obsidian-style)",
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
		"  click ▸Check answer run the tests (title bar; same as ⌃s / :submit)",
		"  ⌃r / ⌃n            run tests / next challenge (in the editor, ⌃r is Vim redo;",
		"                     run with ⌃s or :submit there)",
		"  chat ⌃f ⌃b ⌃d ⌃u   page / half-page scroll",
		"  lesson { / }       previous / next paragraph (Vim; also chat Normal mode)",
		"  mouse wheel        scroll the pane under the cursor",
		"  ⌃c                 quit",
	}, "\n")
	return modalCard("Meari — help", cmds, "esc / q to close")
}

func bold(s string) string { return lipgloss.NewStyle().Bold(true).Render(s) }

// courseCompletion counts done vs total topics of one course (by id/title),
// loading it through the same course loader the tutor runs it with.
func (m Model) courseCompletion(course string) (done, total int) {
	c, err := m.loadVaultCourse(course)
	if err != nil {
		return 0, 0
	}
	for _, t := range c.Topics() {
		total++
		if m.deps.Progress.TopicStatus(t.ID) == "done" {
			done++
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
	return fmt.Sprintf("%-18s %s  %d/%d", truncate(course, 18), bar, done, total)
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

// setupOptions returns the selectable options for the current selection step
// (empty for the free-text topic step; the dashboard has its own entries).
func (m Model) setupOptions() []string {
	if m.setupStep == stepLevel {
		return []string{"Beginner", "Intermediate", "Advanced"}
	}
	return nil
}

func (m *Model) gotoStep(s setupStep) {
	m.setupStep = s
	if s == stepDashboard {
		m.setupCursor = m.dashCursor // Esc-back lands on the row that was chosen
	} else {
		m.setupCursor = 0
	}
	if s == stepTopic {
		m.topicInput.SetValue("")
		m.topicInput.Placeholder = "e.g. spanish past tense · binary search"
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
		// At the first step there's nothing to go back to, so Esc leaves the
		// app — letting users bail straight out of the launch screen.
		if len(m.setupHistory) == 0 {
			return m, tea.Quit
		}
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

	// Selection steps: the dashboard scrolls its own entries, the level step
	// its option list.
	n := len(m.setupOptions())
	if m.setupStep == stepDashboard {
		n = len(m.dash)
	}
	switch msg.String() {
	case "q":
		// "q" quits from any selection step (it's never typed here).
		return m, tea.Quit
	case "g":
		m.setupCursor = 0
	case "G":
		m.setupCursor = n - 1
	case "j", "down":
		if m.setupCursor < n-1 {
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

// setupSelect advances the startup flow after an option is chosen.
func (m Model) setupSelect() (tea.Model, tea.Cmd) {
	switch m.setupStep {
	case stepDashboard:
		if len(m.dash) == 0 {
			return m, nil
		}
		m.dashCursor = m.setupCursor
		e := m.dash[m.setupCursor]
		switch e.kind {
		case dashContinue:
			return m.finishResume()
		case dashCourse:
			// Courses carry their own level; no further questions.
			m.phase = phaseReady
			m.topicInput.Blur()
			return m.switchCourse(e.id)
		case dashTopic:
			m.curriculum = false
			m.advance(stepTopic)
		case dashVault:
			m.exit = SwitchToVault
			return m, tea.Quit // the shell loop opens the vault
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
	m.topic = v
	m.advance(stepLevel)
	return m, nil
}

// finishSetup leaves the launch flow into an AI-generated custom topic (the
// only dashboard path that still asks follow-up questions).
func (m Model) finishSetup() (tea.Model, tea.Cmd) {
	m.deps.Tutor.SetLevel(m.level)
	m.phase = phaseReady
	m.topicInput.Blur()

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
	case paneLesson:
		var cmd tea.Cmd
		m.lesson, cmd = m.lesson.Update(msg)
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
	ctx, cancel := context.WithCancel(context.Background())
	m.streamCancel = cancel
	m.streamStopping = false
	ch, cmd := startChatStream(ctx, func(ctx context.Context, onDelta func(string)) (string, error) {
		return tut.ChatStream(ctx, ctxText, hist, onDelta)
	})
	m.streamCh = ch
	return m, cmd
}

func (m Model) stopStream() (tea.Model, tea.Cmd) {
	if !m.streaming {
		return m, nil
	}
	if m.streamCancel != nil {
		m.streamCancel()
	}
	m.streamStopping = true
	m.loadKind = "stopping tutor"
	m.chat.append(roleSystem, "— stopping tutor reply —")
	return m, nil
}

// handleStreamChunk advances a streaming tutor reply: grow the transcript on
// each delta, finalize the conversation history when done.
func (m Model) handleStreamChunk(msg streamChunkMsg) (tea.Model, tea.Cmd) {
	if m.streamStopping {
		if msg.done || msg.err != nil {
			m.pending--
			m.streaming = false
			m.streamStopping = false
			m.streamCancel = nil
			m.chat.append(roleSystem, "— tutor reply stopped —")
			return m, nil
		}
		return m, listenStream(m.streamCh)
	}
	if msg.err != nil {
		m.pending--
		m.streaming = false
		m.streamCancel = nil
		m.chat.failStream("⚠ chat failed: " + msg.err.Error())
		return m, nil
	}
	if msg.done {
		m.pending--
		m.streaming = false
		m.streamCancel = nil
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
		part := topicViewFromSidebar(it.id)
		if t, ok := m.topicByID[topicIDFromSidebar(it.id)]; ok {
			return m.startTopicView(t, part)
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

// loadCurriculum enters curriculum mode for the course named by key (a
// meari-course id or title — the seeded Go track included), indexes its
// topics, and starts at startID (or the first unfinished topic when empty).
// level rides along only for the warning when the course no longer exists.
func (m *Model) loadCurriculum(key, level, startID string) tea.Cmd {
	c, err := m.loadVaultCourse(key)
	if err != nil {
		m.chat.append(roleSystem, "⚠ No course "+key+" ("+level+") — :topic lists what exists.")
		return nil
	}
	return m.installCurriculum(c, startID)
}

// loadVaultCourse loads a meari-course by id/title and converts it to the
// runnable curriculum form.
func (m *Model) loadVaultCourse(key string) (curriculum.Curriculum, error) {
	if m.deps.Svc == nil {
		return curriculum.Curriculum{}, fmt.Errorf("no vault service")
	}
	c, err := m.deps.Svc.LoadCourse(key)
	if err != nil {
		return curriculum.Curriculum{}, err
	}
	return c.Curriculum(), nil
}

// installCurriculum makes c the active curriculum and starts a topic.
func (m *Model) installCurriculum(c curriculum.Curriculum, startID string) tea.Cmd {
	m.curriculum = true
	m.lang, m.level = c.Lang, c.Level
	m.curr = c
	m.topicByID = map[string]curriculum.Topic{}
	for _, t := range c.Topics() {
		m.topicByID[t.ID] = t
	}
	m.deps.Tutor.SetLevel(c.Level)

	start := m.firstIncompleteTopic()
	if startID != "" {
		if t, ok := m.topicByID[startID]; ok {
			start = t
		}
	}
	return m.startTopicView(start, "lesson")
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
	return m.startTopicView(t, "lesson")
}

func (m *Model) startTopicView(t curriculum.Topic, view string) tea.Cmd {
	if m.current.ID != "" {
		_ = m.deps.Store.Save(m.current.ID, m.editor.Value())
	}
	if view != "quiz" {
		view = "lesson"
	}
	m.currentTopicID = t.ID
	m.currentTopicView = view
	m.topic = t.Title
	m.deps.Progress.MarkTopicInProgress(t.ID)
	m.deps.Progress.SetLast(m.curr.Lang, m.curr.Level, t.ID, t.Title)
	_ = m.deps.Progress.Save()

	lang := m.curr.Lang
	if t.Lang != "" {
		lang = t.Lang // vault courses mix kinds: code topics vs essay topics
	}
	ch := tutor.Challenge{
		ID:          t.ID,
		Prompt:      t.Challenge.Prompt,
		StarterCode: t.Challenge.StarterCode,
		Tests:       t.Challenge.Tests,
		Lang:        lang,
	}
	m.current = ch
	*m.curID = ch.ID
	m.editor.SetLanguage(ch.Lang)
	m.chat.setCodeLang(challengeLang(ch))

	code, stale := m.loadStarterOrDraft(ch)
	m.editor.SetValue(code)
	m.layout() // the pinned prompt header's height depends on the new prompt
	m.rebuildSidebar()
	// Each lesson/quiz row keeps its own chat context. The lesson row shows only
	// lecture content; the quiz/reflection/challenge row shows only the prompt.
	fresh := m.switchChatContext("topic:" + t.ID + "#" + view)
	if fresh && view != "lesson" {
		if quiz := m.topicQuizText(t, ch); quiz != "" {
			m.chat.append(roleQuiz, quiz)
		}
	}
	if view == "lesson" {
		// The lecture lives in the lesson pane when split, in the transcript
		// otherwise — syncLessonPresentation is idempotent across fresh,
		// restored, and cross-mode-restored contexts.
		m.lesson.setCodeLang(challengeLang(ch))
		m.syncLessonPresentation()
	}
	// A lesson lecture reads top-to-bottom, so in the legacy single-pane view
	// it must always open at the first line — whether freshly shown or
	// restored (restore() jumps to the transcript tail). In split mode the
	// lesson pane owns the top-scroll (setLesson) and the chat keeps its
	// conversation tail. Interactive views reset to the top only when fresh.
	if (view == "lesson" && !m.lessonSplit()) || fresh {
		m.chat.gotoTop()
	}
	if stale {
		m.chat.append(roleSystem, staleDraftNotice)
	}
	if view == "lesson" {
		if m.lessonSplit() {
			return m.setFocus(paneLesson) // the reading surface: bare j/k scroll it
		}
		return m.setFocus(paneChat)
	}
	return m.setFocus(paneEditor)
}

// lessonBlockText is the exact transcript form of a topic's lecture; the
// materialize/strip pair below match on it verbatim.
func lessonBlockText(t curriculum.Topic) string {
	return strings.TrimRight(t.Title+"\n\n"+t.Lesson, "\n")
}

// materializeLessonInChat prepends the lecture to the conversation transcript
// (the legacy single-pane view) unless an identical block is already there.
// Prepending never displaces the tail block a stream may be writing to.
func (m *Model) materializeLessonInChat(text string) {
	for _, b := range m.chat.blocks {
		if b.role == roleLesson && b.text == text {
			return
		}
	}
	m.chat.clearSelect()
	m.chat.blocks = append([]chatBlock{{role: roleLesson, text: text}}, m.chat.blocks...)
	m.chat.reflow()
}

// stripLessonFromChat removes the lecture block that split mode shows in its
// own pane. Exact-text matching keeps every other roleLesson block (e.g. an
// ":answer" model solution) intact.
func (m *Model) stripLessonFromChat(text string) {
	for i, b := range m.chat.blocks {
		if b.role == roleLesson && b.text == text {
			m.chat.clearSelect()
			m.chat.blocks = append(m.chat.blocks[:i], m.chat.blocks[i+1:]...)
			m.chat.reflow()
			return
		}
	}
}

// syncLessonPresentation places the current lesson row's lecture on the right
// surface for the active view mode — the lesson pane when split, the chat
// transcript otherwise. Idempotent, so it's safe to call on every topic,
// view, or config change.
func (m *Model) syncLessonPresentation() {
	if !m.curriculum || m.currentTopicView != "lesson" {
		return
	}
	t, ok := m.topicByID[m.currentTopicID]
	if !ok {
		return
	}
	text := lessonBlockText(t)
	if m.lessonSplit() {
		m.lesson.setLesson(text)
		m.stripLessonFromChat(text)
		return
	}
	m.materializeLessonInChat(text)
	// A fresh lecture reads from the top; an ongoing conversation isn't yanked.
	for _, b := range m.chat.blocks {
		if b.role == roleUser || b.role == roleTutor {
			return
		}
	}
	m.chat.gotoTop()
}

func (m Model) topicQuizText(t curriculum.Topic, ch tutor.Challenge) string {
	if len(t.Quiz) > 0 {
		var b strings.Builder
		for i, q := range t.Quiz {
			if i > 0 {
				b.WriteString("\n\n")
			}
			fmt.Fprintf(&b, "%d. %s", i+1, strings.TrimSpace(q.Q))
			for j, choice := range q.Choices {
				if strings.TrimSpace(choice) == "" {
					continue
				}
				fmt.Fprintf(&b, "\n   %s. %s", quizChoiceLabel(j), strings.TrimSpace(choice))
			}
		}
		return b.String()
	}
	prompt := strings.TrimSpace(ch.Prompt)
	if prompt == "" {
		return ""
	}
	return prompt
}

func quizChoiceLabel(i int) string {
	if i >= 0 && i < 26 {
		return string(rune('A' + i))
	}
	return fmt.Sprint(i + 1)
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
	if m.curriculum && m.currentTopicView != "quiz" {
		m.flash("open the quiz row, then submit your answer")
		return nil
	}
	code := m.editor.Value()
	if m.chatCentric() {
		// Chat-centric view: the answer is whatever sits in the chat input
		// (Enter still chats; :submit grades). It joins the transcript as the
		// learner's turn, like any spoken answer.
		text, ok := m.chat.submit()
		if !ok {
			m.flash("type your answer in the chat input, then :submit")
			return nil
		}
		m.chat.append(roleUser, text)
		code = text
	}
	if isReflectionLang(m.current.Lang) && strings.TrimSpace(code) == strings.TrimSpace(m.current.StarterCode) {
		m.flash("write your answer in the editor, then :submit")
		return nil
	}
	_ = m.deps.Store.Save(m.current.ID, code)
	m.chat.append(roleSystem, "▶ checking your answer…")
	m.pending++
	m.loadKind = "checking"
	return runCmd(m.current.Lang, code, m.current)
}

func isReflectionLang(lang string) bool {
	switch strings.ToLower(lang) {
	case "physics", "essay", "quiz":
		return true
	default:
		return false
	}
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
	// A live layout switch can turn the lesson split on or off (it is
	// vertical-only), so the lecture must move between its pane and the
	// transcript before the panes re-flow.
	m.syncLessonPresentation()
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
// visiblePanes lists the focusable panes in visual order, left to right. The
// lesson pane and the editor are mutually exclusive (the lesson fills the
// editor's slot on lesson rows).
func (m Model) visiblePanes() []pane {
	ps := make([]pane, 0, 3)
	if !m.sidebarCollapsed {
		ps = append(ps, paneSidebar)
	}
	if m.lessonSplit() {
		ps = append(ps, paneLesson)
	} else if !m.chatCentric() {
		ps = append(ps, paneEditor)
	}
	if !m.chatHidden() {
		ps = append(ps, paneChat)
	}
	return ps
}

// focusDir moves focus left (d<0) or right (d>0) through the visible panes,
// clamped at the edges — hidden panes are simply not in the list.
func (m *Model) focusDir(d int) tea.Cmd {
	ps := m.visiblePanes()
	i := 0 // focus on a since-hidden pane restarts at the left edge
	for j, p := range ps {
		if p == m.focus {
			i = j
			break
		}
	}
	i += d
	if i < 0 {
		i = 0
	}
	if i >= len(ps) {
		i = len(ps) - 1
	}
	return m.setFocus(ps[i])
}

// setFocus moves keyboard focus, blurring the old input and focusing the new so
// exactly one component captures keys / blinks a cursor at a time. Focus aimed
// at a hidden pane lands on its stand-in: the hidden editor's slot belongs to
// the lesson pane on lesson rows and to the chat otherwise; a hidden lesson
// pane hands off to the chat.
func (m *Model) setFocus(p pane) tea.Cmd {
	if p == paneEditor && m.chatCentric() {
		if m.lessonSplit() {
			p = paneLesson
		} else {
			p = paneChat
		}
	}
	if p == paneLesson && !m.lessonSplit() {
		p = paneChat
	}
	// Focusing the ":chat"-hidden pane (a reply landing, :copy…) unfolds it —
	// keys must never land in a pane that isn't on screen.
	if p == paneChat && m.chatHidden() {
		m.chatCollapsed = false
		m.layout()
	}
	m.editor.Blur()
	m.chat.blur()
	m.lesson.blur()
	m.sidebar.focused = false
	m.focus = p
	switch p {
	case paneEditor:
		return m.editor.Focus()
	case paneChat:
		return m.chat.focus()
	case paneLesson:
		return m.lesson.focus()
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
				status := m.deps.Progress.TopicStatus(t.ID)
				items = append(items, sidebarItem{
					id:     topicLessonSidebarID(t.ID),
					title:  t.Title,
					active: t.ID == m.currentTopicID && m.currentTopicView == "lesson",
				})
				items = append(items, sidebarItem{
					id:     topicQuizSidebarID(t.ID),
					title:  topicQuizSidebarTitle(t),
					status: status,
					active: t.ID == m.currentTopicID && m.currentTopicView == "quiz",
					depth:  1,
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

func topicLessonSidebarID(topicID string) string { return topicID + "#lesson" }
func topicQuizSidebarID(topicID string) string   { return topicID + "#quiz" }

func topicIDFromSidebar(id string) string {
	if base, _, ok := strings.Cut(id, "#"); ok {
		return base
	}
	return id
}

func topicViewFromSidebar(id string) string {
	_, part, ok := strings.Cut(id, "#")
	if ok && part == "quiz" {
		return "quiz"
	}
	return "lesson"
}

func topicQuizSidebarTitle(t curriculum.Topic) string {
	if len(t.Quiz) > 0 {
		if len(t.Quiz) == 1 {
			return "Quiz"
		}
		return fmt.Sprintf("Quiz (%d)", len(t.Quiz))
	}
	return "Quiz"
}

// --- layout ---

// chatCentric reports whether the screen should drop the editor: an explicit
// :view override wins, otherwise curriculum lesson rows read as chat-only
// lectures and quiz/reflection/challenge rows use the editor/code view.
func (m Model) chatCentric() bool {
	switch m.view {
	case "chat":
		return true
	case "code":
		return false
	}
	return m.curriculum && m.currentTopicView == "lesson"
}

// chatHidden reports whether the chat pane is actually dropped this frame:
// the ":chat" toggle holds, except where the chat is the only content
// surface (":view chat", lecture rows without the lesson split) — hiding it
// there would blank the screen, so the pane stays until the view changes.
func (m Model) chatHidden() bool {
	return m.chatCollapsed && !(m.chatCentric() && !m.lessonSplit())
}

// lessonSplit reports whether lesson rows render as sidebar | lesson | chat:
// the lecture as its own read-only pane, the chat holding just the
// conversation. Auto mode only — ":view chat"/":view code" are the legacy
// single-surface modes (m.view's zero value "" IS auto, mirroring
// chatCentric). Invariant: lessonSplit() implies chatCentric().
// TODO(horizontal): stacking a third box is deferred; the gate keeps the
// horizontal layout byte-identical to today.
func (m Model) lessonSplit() bool {
	return !m.horizontal && m.view != "chat" && m.view != "code" &&
		m.curriculum && m.currentTopicView == "lesson"
}

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

	chatOnly := m.chatCentric() // the editor pane is hidden entirely
	chatOff := m.chatHidden()   // the chat pane is dropped (":chat")

	if m.horizontal {
		m.lessonW = 0 // the lesson split is vertical-only (see lessonSplit)
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
		if chatOnly {
			// One right box instead of two: the chat takes the whole column.
			m.chatH = m.contentH
			m.editorH = 0
		} else if chatOff {
			// ":chat" dropped the top box: the editor takes the whole column.
			m.chatH = 0
			m.editorH = m.contentH
		} else {
			// The two stacked boxes cost 4 border rows vs the sidebar's 2, so
			// their combined content height is contentH-2.
			rightContent := m.contentH - 2
			if rightContent < 2 {
				rightContent = 2
			}
			// In this stacked layout the editor sits below the chat, so :compact /
			// :wide trade height between them rather than width.
			chatPct := clampRange(m.deps.Cfg.ChatPct(55)-m.editorBias, 30, 85)
			m.chatH = clampMin(rightContent*chatPct/100, 3)
			m.editorH = clampMin(rightContent-m.chatH, 3)
		}

		m.sidebar.setSize(m.sidebarW, m.contentH)
		if m.chatH > 0 {
			m.chat.setSize(m.rightW, m.chatH)
		}
		if m.editorH > 0 {
			m.editor.SetSize(m.rightW, max(1, m.editorH-1-len(m.promptHeaderLines(m.rightW))))
		}
		return
	}

	// vertical: (sidebar) | editor-or-lesson | chat. Each visible box eats 2
	// border columns: 6 with the sidebar, 4 when it's folded away — minus 2
	// again in the chat-centric view WITHOUT the lesson split, which drops
	// the middle box entirely (the split brings a third box back).
	borders := 6
	if m.sidebarCollapsed {
		borders = 4
	}
	if chatOnly && !m.lessonSplit() {
		borders -= 2
	}
	if chatOff {
		borders -= 2
	}
	contentW := m.width - borders
	if contentW < 3 {
		contentW = 3
	}
	if m.sidebarCollapsed {
		m.sidebarW = 0
	} else {
		m.sidebarW = clampMin(contentW*m.deps.Cfg.SidebarPct(22)/100, 12)
	}
	switch {
	case chatOnly && m.lessonSplit() && chatOff:
		// ":chat" on a lecture row: the lesson pane reads full width.
		m.chatW, m.editorW = 0, 0
		m.lessonW = clampMin(contentW-m.sidebarW, 10)
	case chatOnly && m.lessonSplit():
		// Lesson rows: the lecture pane takes the editor's slot, so the same
		// :compact / :wide bias trades width between the lesson and the chat.
		chatPct := clampRange(m.deps.Cfg.ChatPct(30)-m.editorBias, 15, 75)
		m.chatW = clampMin(contentW*chatPct/100, 16)
		m.editorW = 0
		m.lessonW = clampMin(contentW-m.sidebarW-m.chatW, 10)
	case chatOnly:
		m.editorW, m.lessonW = 0, 0
		m.chatW = clampMin(contentW-m.sidebarW, 16)
	case chatOff:
		// ":chat" on an editor row: the editor absorbs the chat's width.
		m.chatW, m.lessonW = 0, 0
		m.editorW = clampMin(contentW-m.sidebarW, 10)
	default:
		// The configured split is the base; :compact / :wide shift it live.
		m.lessonW = 0
		chatPct := clampRange(m.deps.Cfg.ChatPct(30)-m.editorBias, 15, 75)
		m.chatW = clampMin(contentW*chatPct/100, 16)
		m.editorW = clampMin(contentW-m.sidebarW-m.chatW, 10)
	}

	m.sidebar.setSize(m.sidebarW, m.contentH)
	if m.editorW > 0 {
		m.editor.SetSize(m.editorW, max(1, m.contentH-1-len(m.promptHeaderLines(m.editorW))))
	}
	if m.lessonW > 0 {
		m.lesson.setSize(m.lessonW, m.contentH)
	}
	if m.chatW > 0 {
		m.chat.setSize(m.chatW, m.contentH)
	}
}

func (m Model) View() string {
	if m.width == 0 {
		return "starting…"
	}
	if m.width < 60 || m.height < 16 {
		return "Terminal too small — please enlarge to at least 60×16."
	}

	if m.phase == phaseSetup {
		if m.setupStep == stepDashboard {
			return m.dashboardView()
		}
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.setupView())
	}

	if m.overlay != overlayNone {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.overlayView())
	}

	var row string
	if m.horizontal {
		// sidebar on the left; content (chat) above the input (editor) on the
		// right — the editor box disappears in the chat-centric view, the chat
		// box under ":chat".
		var right string
		if m.chatH > 0 {
			right = m.box(paneChat, m.rightW, m.chatH, m.chat.view())
		}
		if m.editorH > 0 {
			ed := m.box(paneEditor, m.rightW, m.editorH, m.editorPaneView(m.rightW))
			if right == "" {
				right = ed
			} else {
				right = lipgloss.JoinVertical(lipgloss.Left, right, ed)
			}
		}
		if m.sidebarCollapsed {
			row = right
		} else {
			sb := m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view())
			row = lipgloss.JoinHorizontal(lipgloss.Top, sb, right)
		}
	} else {
		var panes []string
		if !m.sidebarCollapsed {
			panes = append(panes, m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view()))
		}
		if m.lessonW > 0 { // lesson rows: the lecture pane fills the editor's slot
			panes = append(panes, m.box(paneLesson, m.lessonW, m.contentH, m.lesson.view()))
		}
		if m.editorW > 0 {
			panes = append(panes, m.box(paneEditor, m.editorW, m.contentH, m.editorPaneView(m.editorW)))
		}
		if m.chatW > 0 {
			panes = append(panes, m.box(paneChat, m.chatW, m.contentH, m.chat.view()))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, panes...)
	}

	frame := lipgloss.JoinVertical(lipgloss.Left, m.titleView(), row, m.statusView())
	// Final guard: never emit more than the screen holds, so a too-tall pane
	// can't scroll the alt-screen and leave residue behind.
	return lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(frame)
}

// setupView renders the follow-up steps after a dashboard choice (the custom
// topic input and the level question) as a centered card.
func (m Model) setupView() string {
	var title, body string

	switch m.setupStep {
	case stepTopic:
		title = "What do you want to learn?"
		body = m.topicInput.View()
	case stepLevel:
		title = "What's your experience level?"
		body = m.setupMenu()
	}

	hint := "↑/↓ or j/k to move · enter to choose · esc to go back · q to quit"
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

// dashboardView renders the launch screen: a full-screen dashboard of every
// way into the app — resume, your vault courses, the built-in tracks, a
// custom topic, the vault — grouped into sections with the cursor row
// highlighted, and scrolled to keep the cursor visible on short terminals.
func (m Model) dashboardView() string {
	colW := m.width - 8
	if colW > 76 {
		colW = 76
	}

	// Build display rows; the cursor's entry maps to one row.
	var rows []string
	cursorRow := 0
	section := ""
	for i, e := range m.dash {
		if e.section != section {
			section = e.section
			if len(rows) > 0 {
				rows = append(rows, "")
			}
			head := "── " + section + " "
			rows = append(rows, hintStyle.Render(head+strings.Repeat("─", max(0, colW-lipgloss.Width(head)))))
		}
		meta := e.meta
		title := truncate(e.title, max(8, colW-2-lipgloss.Width(meta)-2))
		pad := max(2, colW-2-lipgloss.Width(title)-lipgloss.Width(meta))
		if i == m.setupCursor {
			cursorRow = len(rows)
			rows = append(rows, selectedRow.Render("▸ "+title+strings.Repeat(" ", pad)+meta))
		} else {
			rows = append(rows, "  "+title+strings.Repeat(" ", pad)+hintStyle.Render(meta))
		}
	}

	// Window the rows around the cursor when the terminal is short.
	maxRows := m.height - 7 // title bar, greeting, paddings, status bar
	if maxRows < 3 {
		maxRows = 3
	}
	if len(rows) > maxRows {
		start := cursorRow - maxRows/2
		if start < 0 {
			start = 0
		}
		if start > len(rows)-maxRows {
			start = len(rows) - maxRows
		}
		rows = rows[start : start+maxRows]
	}

	greeting := lipgloss.NewStyle().Bold(true).Render("What will you learn today?")
	body := greeting + "\n\n" + strings.Join(rows, "\n")

	head := titleBar.Width(m.width).Render("Meari — 메아리")
	hints := "↑/↓ or j/k move · enter choose · q / esc quit"
	foot := statusBar.Width(m.width).Render(hintStyle.Render(hints))
	mid := lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, body)
	return head + "\n" + mid + "\n" + foot
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

	// Right-align the clickable "Check answer" button when one is shown, filling
	// the gap so the title bar background spans the full width.
	if x0, _, ok := m.checkButtonBounds(); ok {
		left := fitWidth(t, x0-1) // inner cols 1..x0-1 (titleBar pads 1 on the left)
		return titleBar.Width(m.width).Render(left + checkButton.Render(checkButtonText))
	}
	return titleBar.Width(m.width).Render(t)
}

// checkButtonBounds returns the half-open screen-column range [x0, x1) on row 0
// occupied by the "Check answer" button, and ok=false when it isn't drawn (no
// active challenge, a modal/command line is up, or the window is too narrow).
func (m Model) checkButtonBounds() (x0, x1 int, ok bool) {
	if m.phase != phaseReady || m.cmdMode || m.overlay != overlayNone {
		return 0, 0, false
	}
	if m.current.ID == "" {
		return 0, 0, false
	}
	w := lipgloss.Width(checkButtonText)
	x0 = m.width - 1 - w // titleBar pads 1 on the right, so the button ends at width-2
	x1 = m.width - 1
	if x0 < 8 { // keep room for at least a little title text
		return 0, 0, false
	}
	return x0, x1, true
}

// fitWidth pads s with spaces (or truncates it with an ellipsis) so it occupies
// exactly w display cells. Used to place a right-aligned element on a line.
func fitWidth(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if cur := lipgloss.Width(s); cur <= w {
		return s + strings.Repeat(" ", w-cur)
	}
	var b strings.Builder
	used := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if used+rw > w-1 { // leave one cell for the ellipsis
			break
		}
		b.WriteRune(r)
		used += rw
	}
	out := b.String() + "…"
	if pad := w - lipgloss.Width(out); pad > 0 {
		out += strings.Repeat(" ", pad)
	}
	return out
}

func (m Model) statusView() string {
	if m.cmdMode {
		line := m.cmdLine.View()
		if h := m.cmdComp.Hint(); h != "" {
			line += "   " + hintStyle.Render(h)
		}
		// MaxWidth keeps a long wildmenu to the single status row.
		return statusBar.Width(m.width).MaxWidth(m.width).Render(line)
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
	case m.focus == paneLesson:
		hints = "j/k ⌃d/⌃u scroll · drag/⌥c copy · ⌃w l chat · :view chat single pane"
	case m.focus == paneChat && m.chatCentric():
		hints = "enter chat · :submit grade your answer · ⌥o/:copy copy · ⌃f/⌃b page · :view code"
	case m.focus == paneChat:
		hints = "enter send · drag to copy · ⌥o/:copy copy reply · ⌃f/⌃b page"
	case m.focus == paneSidebar:
		hints = "j/k move · enter/click open · : cmds (:help) · ⌃r run · ⌃c quit"
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
	case paneLesson:
		return "lesson"
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
	if m.curriculum && m.currentTopicView != "quiz" {
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
