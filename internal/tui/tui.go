// Package tui is the three-pane terminal UI for Meari: a challenge/draft
// list (left), the code editor (center), and an interactive tutor chat (right).
// It replaces the old linear stdin loop — one Bubble Tea program owns all three
// panes, and every AI call and test run is dispatched as an async tea.Cmd so the
// UI never blocks.
package tui

import (
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"meari/internal/config"
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

	// Curriculum mode: the left pane is the guided learning path instead of the
	// generated-challenge list. Lessons and challenges are pre-authored (no LLM).
	curriculum     bool
	curr           curriculum.Curriculum
	topicByID      map[string]curriculum.Topic
	currentTopicID string

	// horizontal selects the stacked layout (content on top, input on the
	// bottom) instead of the side-by-side columns. Toggled live via :config.
	horizontal bool

	// pendingWindow is set after Ctrl-W; the next key (h/j/k/l) picks a pane,
	// Vim window-command style.
	pendingWindow bool

	pending  int    // in-flight async ops; spinner shows while > 0
	loadKind string // label for the most recent op
	spin     spinner.Model
	err      error

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

	sp := spinner.New(spinner.WithSpinner(spinner.Dot))

	m := Model{
		deps:       d,
		horizontal: d.Cfg.Horizontal(),
		sidebar:    newSidebar(),
		editor:     editor.New("", d.Cfg.VimEditor(), save),
		chat:       newChat(),
		topicInput: ti,
		curID:      curID,
		challenges: map[string]tutor.Challenge{},
		spin:       sp,
	}
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

	case chatReplyMsg:
		m.pending--
		m.chat.append(roleTutor, msg.text)
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant", Content: msg.text})
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
		m.pendingWindow = true
		return m, nil
	case "tab":
		return m, m.setFocus(m.shiftPane(1))
	case "shift+tab":
		return m, m.setFocus(m.shiftPane(-1))
	case "ctrl+r":
		return m, m.startRun()
	case "ctrl+n":
		return m, m.nextChallenge()
	}

	switch m.focus {
	case paneSidebar:
		var enter bool
		m.sidebar, enter = m.sidebar.Update(msg)
		if enter {
			return m, m.openSelected()
		}
		return m, nil

	case paneEditor:
		tm, cmd := m.editor.Update(msg)
		m.editor = tm.(editor.Model)
		return m, cmd

	case paneChat:
		if msg.Type == tea.KeyEnter {
			return m.submitChat()
		}
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(msg)
		return m, cmd
	}
	return m, nil
}

// handleMouse routes wheel events to the pane under the cursor, like ranger/lf:
// scrolling never changes focus, it just moves the hovered pane. Non-wheel mouse
// events (clicks, motion) are ignored so the terminal keeps its own selection.
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.phase == phaseSetup || !tea.MouseEvent(msg).IsWheel() {
		return m, nil
	}
	p, ok := m.paneAt(msg.X, msg.Y)
	if !ok {
		return m, nil
	}
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

// paneAt maps a terminal cell to the pane drawn there, accounting for the title
// bar (top row), status bar (bottom row), and the one-cell border around each
// box. ok is false for cells outside any pane (the bars or border gutters).
func (m Model) paneAt(x, y int) (pane, bool) {
	if y < 1 || y > m.height-2 { // row 0 is the title bar, the last row the status bar
		return 0, false
	}
	if m.horizontal {
		// sidebar on the left; chat stacked over the editor on the right.
		if x < m.sidebarW+2 {
			return paneSidebar, true
		}
		if y < 1+m.chatH+2 {
			return paneChat, true
		}
		return paneEditor, true
	}
	// vertical: sidebar | editor | chat, each box +2 columns for its border.
	switch {
	case x < m.sidebarW+2:
		return paneSidebar, true
	case x < m.sidebarW+2+m.editorW+2:
		return paneEditor, true
	default:
		return paneChat, true
	}
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

// submitChat sends the chat input as a question to the tutor.
func (m Model) submitChat() (tea.Model, tea.Cmd) {
	text, ok := m.chat.submit()
	if !ok {
		return m, nil
	}
	m.chat.append(roleUser, text)
	m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "user", Content: text})
	m.pending++
	m.loadKind = "tutor thinking"
	return m, chatCmd(m.deps.Tutor, m.chatHist)
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
	m.chat.append(roleSystem, "No challenge data for \""+it.id+"\" this session — press Ctrl-N to generate a new one.")
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

	code := ch.StarterCode
	if d, ok := m.deps.Store.Load(ch.ID); ok {
		code = d
	}
	m.editor.SetValue(code)
	m.rebuildSidebar()
	m.chat.append(roleLesson, t.Title+"\n\n"+t.Lesson)
	m.chat.append(roleSystem, "📝 Challenge: "+ch.Prompt)
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

	code := ch.StarterCode
	if d, ok := m.deps.Store.Load(ch.ID); ok {
		code = d
	}
	m.editor.SetValue(code)
	m.rebuildSidebar()
	m.chat.append(roleSystem, "📝 Challenge: "+ch.Prompt)
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

// applyConfigReload re-reads the config after editing and live-applies the
// settings that can change mid-session (currently the layout).
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
	m.chat.append(roleSystem, "✓ Config reloaded — layout is now "+cfg.UI.Layout+
		". (Editor keybindings and AI settings apply on next launch.)")
}

// --- focus & ordering helpers ---

func (m *Model) shiftPane(delta int) pane {
	return pane((int(m.focus) + delta + 3) % 3)
}

// focusDir moves focus one pane in direction d (-1 left, +1 right), clamping at
// the edges (Vim windows don't wrap).
func (m *Model) focusDir(d int) tea.Cmd {
	n := int(m.focus) + d
	if n < 0 {
		n = 0
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
		// sidebar | (content on top, input on the bottom)
		contentW := m.width - 4 // two columns -> 4 border columns
		if contentW < 3 {
			contentW = 3
		}
		m.sidebarW = clampMin(contentW*22/100, 14)
		m.rightW = clampMin(contentW-m.sidebarW, 20)
		// The two stacked boxes cost 4 border rows vs the sidebar's 2, so their
		// combined content height is contentH-2.
		rightContent := m.contentH - 2
		if rightContent < 2 {
			rightContent = 2
		}
		m.chatH = clampMin(rightContent*55/100, 3)
		m.editorH = clampMin(rightContent-m.chatH, 3)

		m.sidebar.setSize(m.sidebarW, m.contentH)
		m.chat.setSize(m.rightW, m.chatH)
		m.editor.SetSize(m.rightW, max(1, m.editorH-1))
		return
	}

	// vertical: sidebar | editor | chat (three borders eat 6 columns)
	contentW := m.width - 6
	if contentW < 3 {
		contentW = 3
	}
	m.sidebarW = clampMin(contentW*22/100, 12)
	m.chatW = clampMin(contentW*30/100, 16)
	m.editorW = clampMin(contentW-m.sidebarW-m.chatW, 10)

	m.sidebar.setSize(m.sidebarW, m.contentH)
	m.editor.SetSize(m.editorW, max(1, m.contentH-1))
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

	var row string
	if m.horizontal {
		// sidebar on the left; content (chat) above the input (editor) on the right.
		sb := m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view())
		ch := m.box(paneChat, m.rightW, m.chatH, m.chat.view())
		ed := m.box(paneEditor, m.rightW, m.editorH, m.editorPaneView(m.rightW))
		right := lipgloss.JoinVertical(lipgloss.Left, ch, ed)
		row = lipgloss.JoinHorizontal(lipgloss.Top, sb, right)
	} else {
		sb := m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view())
		ed := m.box(paneEditor, m.editorW, m.contentH, m.editorPaneView(m.editorW))
		ch := m.box(paneChat, m.chatW, m.contentH, m.chat.view())
		row = lipgloss.JoinHorizontal(lipgloss.Top, sb, ed, ch)
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
	return m.challengeHeader(w) + "\n" + m.editor.View()
}

func (m Model) challengeHeader(w int) string {
	label := "No challenge selected"
	if m.current.ID != "" {
		prompt := firstLine(m.current.Prompt)
		if prompt == "" {
			prompt = m.topic
		}
		label = strings.ToUpper(m.current.Lang) + " · " + prompt
	}
	if w > 0 && lipgloss.Width(label) > w-2 {
		label = truncate(label, max(1, w-2))
	}
	return editorHeader.Width(w).Render(label)
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
	left := "[" + m.focusName() + "]"
	if m.pending > 0 {
		left += " " + m.spin.View() + " " + m.loadKind
	} else if m.err != nil {
		left += " " + errStyle.Render("error: "+m.err.Error())
	}
	hints := "tab / ⌃w h·l focus · ⌃r run · ⌃n next · ⌃c quit"
	switch {
	case m.pendingWindow:
		hints = errStyle.Render("⌃w") + " window: h/l choose pane"
	case m.focus == paneChat:
		hints = "⌃f/⌃b page · ⌃d/⌃u half · ⇧↑/↓ line · wheel scrolls · enter send"
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

// --- small utilities ---

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
