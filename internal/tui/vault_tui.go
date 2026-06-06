package tui

// vault_tui.go is the terminal front-end for the general learning vault. Like the
// web GUI, it is a thin presentation layer over core.Service: a three-pane
// program (notes | editor | chat/study) where all real work — listing notes,
// opening/saving them, generating a lesson, grading an essay, chatting — is done
// by core and this model only renders the result. It reuses the existing
// sidebar/chat/editor components and styles from this package.

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/editor"
	"meari/internal/tutor"
)

// VaultModel is the root model for `meari notes` — the vault terminal UI.
type VaultModel struct {
	svc *core.Service

	width, height int
	focus         pane
	exit          SwitchTarget // set by :tutor so RunVault can report a mode switch

	sidebar sidebarModel
	editor  editor.Model
	chat    chatModel

	notes        []core.NoteMeta
	current      string // path of the open note ("" = none)
	currentTitle string
	curPath      *string          // shared with the editor save closure
	chatHist     []tutor.ChatTurn // tutor conversation history

	// Per-note chat contexts: each note keeps its own transcript and tutor
	// conversation, restored when the learner reopens it.
	chatByNote map[string][]chatBlock
	histByNote map[string][]tutor.ChatTurn

	// Streaming chat reply state: one reply at a time.
	streaming bool
	streamCh  chan streamChunkMsg

	// Essay study state. While studying, the editor holds the learner's answer
	// (not the note), and autosave to the note is suspended.
	studyMode   bool
	studyPrompt string

	// Backlinks ("notes that link here") for the open note, shown as a footer
	// under the editor (Obsidian-style). showBacklinks toggles the panel.
	backlinks     []core.NoteMeta
	showBacklinks bool

	// global ex-command line (":" from the notes pane)
	cmdMode bool
	cmdLine textinput.Model
	cmdHist editor.CmdHistory

	// Vim-style chords mirroring the coding TUI: pendingWindow is set after
	// Ctrl-W (the next h/j/k/l picks a pane by direction); pendingLeader is set
	// after "," in the editor's Normal mode (",n" folds the sidebar).
	pendingWindow bool
	pendingLeader bool

	// editorBias shifts the editor/chat split (":wide" grows the editor,
	// ":compact" grows the chat), sharing the classic TUI's step/clamp.
	editorBias int

	// sidebarCollapsed hides the notes pane (":fold"); starts from config.
	sidebarCollapsed bool

	// cfg supplies the configured pane ratios and editor keybindings.
	cfg config.Config

	pending  int
	loadKind string
	spin     spinner.Model
	err      error

	// notice is transient command feedback shown in the status bar.
	notice   string
	noticeAt time.Time

	sidebarW, editorW, chatW, contentH int
}

// RunVault constructs and runs the vault terminal UI over svc. The Outcome tells
// the shell loop (main.runShell) whether to quit or hand off to the coding TUI.
func RunVault(svc *core.Service, cfg config.Config) (Outcome, error) {
	m := newVaultModel(svc, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	final, err := p.Run()
	out := Outcome{}
	if fm, ok := final.(VaultModel); ok {
		out = Outcome{Target: fm.exit}
	}
	return out, err
}

func newVaultModel(svc *core.Service, cfg config.Config) VaultModel {
	vim := cfg.VimEditor()
	curPath := new(string)
	m := VaultModel{
		svc:              svc,
		cfg:              cfg,
		curPath:          curPath,
		showBacklinks:    true,
		sidebarCollapsed: cfg.UI.SidebarFolded,
		sidebar:          newSidebar(),
		chat:             newChat(),
		chatByNote:       map[string][]chatBlock{},
		histByNote:       map[string][]tutor.ChatTurn{},
		spin:             spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
	// The editor's save closure persists the open note — but never while the
	// learner is writing an essay answer (curPath is blanked during study).
	save := func(code string) error {
		if *curPath == "" {
			return nil
		}
		_, err := svc.SaveNote(*curPath, code)
		return err
	}
	m.editor = editor.New("", vim, save)
	m.editor.SetLanguage("markdown")
	m.editor.SetShowLineNumbers(cfg.LineNumbers())

	cl := textinput.New()
	cl.Prompt = ":"
	m.cmdLine = cl

	if m.sidebarCollapsed {
		m.focus = paneEditor
	} else {
		m.focus = paneSidebar
		m.sidebar.focused = true
	}
	return m
}

func (m VaultModel) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, vListCmd(m.svc))
}

func (m VaultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		// Mirror in-flight work as an animated progress line inside the chat pane.
		if m.pending > 0 {
			m.chat.setBusy(m.loadKind)
			m.chat.tickBusy()
		} else {
			m.chat.setBusy("")
		}
		return m, cmd

	case vNotesMsg:
		m.notes = msg.notes
		m.rebuildSidebar()
		return m, nil

	case vOpenedMsg:
		m.studyMode = false
		m.switchNoteChat(msg.note)
		m.current = msg.note.Path
		m.currentTitle = msg.note.Title
		*m.curPath = msg.note.Path
		m.editor.SetValue(msg.note.Body)
		m.backlinks = nil // drop the previous note's backlinks until the fetch returns
		m.rebuildSidebar()
		return m, tea.Batch(m.setFocus(paneEditor), vBacklinksCmd(m.svc, m.current))

	case vBacklinksMsg:
		if msg.path == m.current {
			m.backlinks = msg.links
			m.layout() // the footer height may have changed
		}
		return m, nil

	case vGeneratedMsg:
		m.pending--
		m.chat.append(roleLesson, "Created note: "+msg.meta.Title)
		return m, tea.Batch(vListCmd(m.svc), vOpenCmd(m.svc, msg.meta.Path))

	case vSavedMsg:
		// Refresh the list in case the title/subject changed; keep editing.
		return m, vListCmd(m.svc)

	case streamChunkMsg:
		return m.handleStreamChunk(msg)

	case vEssayMsg:
		m.pending--
		pct := int(msg.res.Score*100 + 0.5)
		role := roleOK
		if msg.res.Score < 0.6 {
			role = roleFail
		}
		m.chat.append(role, "Score: "+itoa(pct)+"%")
		m.chat.append(roleTutor, msg.res.Feedback)
		// Join the conversation so follow-up questions can refer to the feedback.
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant",
			Content: "Essay feedback (score " + itoa(pct) + "%): " + msg.res.Feedback})
		return m, nil

	case vAnswerMsg:
		m.pending--
		m.chat.append(roleLesson, "Model answer\n\n"+msg.text)
		m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "assistant", Content: "Model answer:\n" + msg.text})
		return m, nil

	case vErrMsg:
		m.pending--
		m.err = msg.err
		m.chat.append(roleSystem, "⚠ "+msg.kind+" failed: "+msg.err.Error())
		return m, nil

	case editor.DoneMsg:
		switch msg.Action {
		case editor.ActionSubmit:
			return m.submitEditor()
		case editor.ActionQuit:
			return m, tea.Quit
		}
		return m, nil

	case editor.RunCommandMsg:
		return m.runEx(msg.Raw)

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}
	return m.forwardToFocus(msg)
}

func (m VaultModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.cmdMode {
		return m.updateCmdLine(msg)
	}

	// Ctrl-W starts a window command; the next key chooses a pane by direction
	// (Vim window-style), mirroring the coding TUI.
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
		return m, tea.Quit
	case "ctrl+w":
		// Focus moves via the Vim window chord (Ctrl-W then h/l); bare Tab is left
		// for the panes (e.g. indenting in the editor).
		m.pendingWindow = true
		return m, nil
	}

	// A leader chord lives for exactly one keystroke: clear it here so a stray
	// "," can never carry across a focus change or non-editor key.
	leader := m.pendingLeader
	m.pendingLeader = false

	switch m.focus {
	case paneSidebar:
		if msg.String() == ":" {
			m.cmdMode = true
			m.cmdLine.SetValue("")
			m.cmdHist.Open()
			return m, m.cmdLine.Focus()
		}
		var enter bool
		m.sidebar, enter = m.sidebar.Update(msg)
		if enter {
			if it, ok := m.sidebar.selected(); ok {
				return m, vOpenCmd(m.svc, it.id)
			}
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
		// Copy the tutor's last reply: Alt+O (Linux) / Option+O (macOS, which
		// arrives as "ø"/"Ø" unless the terminal sends Option as Meta).
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

// focusDir moves focus left (d<0) or right (d>0) between panes, clamped to the
// visible panes — a folded sidebar isn't a target. Mirrors the coding TUI.
func (m *VaultModel) focusDir(d int) tea.Cmd {
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

// handleMouse routes wheel events to the pane under the cursor (scrolling never
// steals focus) and focuses the pane under the cursor on a left click.
func (m VaultModel) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.cmdMode {
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
			m.chat, cmd = m.chat.Update(msg)
			return m, cmd
		case paneSidebar:
			switch msg.Button {
			case tea.MouseButtonWheelDown:
				m.sidebar.move(1)
			case tea.MouseButtonWheelUp:
				m.sidebar.move(-1)
			}
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

// paneAt maps a terminal cell to the pane drawn there: row 0 is the title bar,
// the last row the status bar, and each box adds 2 border columns.
func (m VaultModel) paneAt(x, y int) (pane, bool) {
	if y < 1 || y > m.height-2 {
		return 0, false
	}
	sidebarSpan := m.sidebarW + 2
	if m.sidebarCollapsed {
		sidebarSpan = 0
	}
	switch {
	case x < sidebarSpan:
		return paneSidebar, true
	case x < sidebarSpan+m.editorW+2:
		return paneEditor, true
	default:
		return paneChat, true
	}
}

func (m VaultModel) updateCmdLine(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		return m, tea.Quit
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

// runEx dispatches a vault ex-command (without the leading colon).
func (m VaultModel) runEx(raw string) (tea.Model, tea.Cmd) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return m, nil
	}
	args := strings.TrimSpace(strings.TrimPrefix(raw, fields[0]))
	switch fields[0] {
	case "learn", "gen", "lesson":
		if args == "" {
			m.flash("usage: :learn <what you want to learn>")
			return m, nil
		}
		m.pending++
		m.loadKind = "generating lesson"
		m.chat.append(roleSystem, "▶ generating a lesson on "+args+"…")
		return m, vGenCmd(m.svc, args)
	case "new":
		if args == "" {
			m.flash("usage: :new <note title>")
			return m, nil
		}
		path := args + ".md"
		return m, vSaveOpenCmd(m.svc, path, "# "+args+"\n\n")
	case "fold", "sidebar":
		return m.cmdFold()
	case "compact":
		return m.cmdResizeEditor(-editorBiasStep)
	case "wide":
		return m.cmdResizeEditor(editorBiasStep)
	case "essay":
		return m.startEssay(args)
	case "grade":
		return m.gradeEssay()
	case "answer":
		return m.revealAnswer()
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
	case "done":
		return m.endEssay()
	case "backlinks", "links":
		m.showBacklinks = !m.showBacklinks
		m.layout()
		if m.showBacklinks {
			m.flash("backlinks panel on")
			if m.current != "" {
				return m, vBacklinksCmd(m.svc, m.current)
			}
		} else {
			m.flash("backlinks panel off")
		}
		return m, nil
	case "tutor", "code":
		m.exit = SwitchToTutor
		return m, tea.Quit // the shell loop opens the coding TUI
	case "q", "quit":
		return m, tea.Quit
	default:
		m.flash("unknown command: :" + raw +
			"  (try :learn · :new · :essay · :grade · :answer · :done · :backlinks · :tutor · :fold · :compact · :wide · :q)")
		return m, nil
	}
}

// cmdFold toggles the notes pane. Folding the focused pane moves focus to the
// editor so keys never vanish into a hidden pane.
func (m VaultModel) cmdFold() (tea.Model, tea.Cmd) {
	m.sidebarCollapsed = !m.sidebarCollapsed
	var cmd tea.Cmd
	if m.sidebarCollapsed && m.focus == paneSidebar {
		cmd = m.setFocus(paneEditor)
	}
	m.layout()
	if m.sidebarCollapsed {
		m.flash("Notes pane folded — :fold to bring it back")
	}
	return m, cmd
}

// cmdResizeEditor nudges the editor/chat split by delta percentage points
// (":compact" grows the chat, ":wide" grows the editor), clamps, and re-flows.
func (m VaultModel) cmdResizeEditor(delta int) (tea.Model, tea.Cmd) {
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

// submitEditor handles Ctrl-S / :submit: in study mode it grades the answer,
// otherwise it saves the open note.
func (m VaultModel) submitEditor() (tea.Model, tea.Cmd) {
	if m.studyMode {
		return m.gradeEssay()
	}
	if m.current == "" {
		return m, nil
	}
	return m, vSaveCmd(m.svc, m.current, m.editor.Value())
}

// startEssay begins an essay study on the current note. The editor is cleared to
// hold the learner's answer; autosave to the note is suspended.
func (m VaultModel) startEssay(prompt string) (tea.Model, tea.Cmd) {
	if m.current == "" {
		m.flash("open a note first, then :essay to study it")
		return m, nil
	}
	// Persist any edits to the note before repurposing the editor.
	_, _ = m.svc.SaveNote(m.current, m.editor.Value())
	if prompt == "" {
		prompt = "Explain " + m.currentTitle + " in your own words."
	}
	m.studyMode = true
	m.studyPrompt = prompt
	*m.curPath = "" // suspend note autosave while answering
	m.editor.SetValue("")
	m.layout() // the pinned prompt header's height depends on the prompt
	m.chat.append(roleSystem, "— essay study started — write your answer under the prompt in the editor, then :grade (or Ctrl-S); :answer reveals a model answer —")
	return m, m.setFocus(paneEditor)
}

func (m VaultModel) gradeEssay() (tea.Model, tea.Cmd) {
	if !m.studyMode {
		m.flash("not studying — :essay to start")
		return m, nil
	}
	answer := strings.TrimSpace(m.editor.Value())
	if answer == "" {
		m.flash("write an answer first")
		return m, nil
	}
	m.pending++
	m.loadKind = "grading"
	return m, vEssayCmd(m.svc, m.studyPrompt, answer)
}

// revealAnswer shows a model answer for the current essay prompt (the learner
// explicitly asked, so revealing is fine — unlike grading feedback).
func (m VaultModel) revealAnswer() (tea.Model, tea.Cmd) {
	if !m.studyMode {
		m.flash("not studying — :essay to start, then :answer to see a model answer")
		return m, nil
	}
	m.pending++
	m.loadKind = "writing model answer"
	return m, vAnswerCmd(m.svc, m.studyPrompt)
}

func (m VaultModel) endEssay() (tea.Model, tea.Cmd) {
	if !m.studyMode {
		return m, nil
	}
	m.studyMode = false
	m.chat.append(roleSystem, "— essay study ended —")
	if m.current != "" {
		return m, vOpenCmd(m.svc, m.current)
	}
	return m, nil
}

// flash shows transient feedback in the status bar for a few seconds.
func (m *VaultModel) flash(s string) {
	if s == "" {
		return
	}
	m.notice = s
	m.noticeAt = time.Now()
}

// switchNoteChat swaps the chat pane to the opened note's own transcript and
// tutor conversation, saving the outgoing note's first. Reopening a note brings
// its past study activity back; a first visit starts a clean pane with a header.
func (m *VaultModel) switchNoteChat(n core.Note) {
	if n.Path == m.current {
		return
	}
	if m.current != "" {
		m.chatByNote[m.current] = m.chat.snapshot()
		m.histByNote[m.current] = m.chatHist
	}
	saved, visited := m.chatByNote[n.Path]
	m.chat.restore(saved)
	m.chatHist = m.histByNote[n.Path]
	if !visited {
		m.chat.append(roleSystem, "— "+n.Title+" —")
	}
}

// chatContext describes what the learner is looking at — the open note (and,
// during essay study, the prompt plus their draft answer) — so chat replies
// stay grounded in the current material.
func (m VaultModel) chatContext() string {
	if m.current == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString("Current note — " + m.currentTitle + "\n")
	if m.studyMode {
		b.WriteString("\nEssay prompt: " + m.studyPrompt + "\n")
		b.WriteString("Learner's draft answer:\n" + m.editor.Value() + "\n")
		if n, err := m.svc.OpenNote(m.current); err == nil {
			b.WriteString("\nNote content:\n" + n.Body)
		}
	} else {
		b.WriteString("\nNote content (as in the learner's editor):\n" + m.editor.Value())
	}
	return strings.TrimSpace(b.String())
}

// submitChat sends the chat input to the tutor, streaming the reply into the
// transcript, grounded in the open note / study state.
func (m VaultModel) submitChat() (tea.Model, tea.Cmd) {
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

	svc := m.svc
	ctxText := m.chatContext()
	hist := append([]tutor.ChatTurn(nil), m.chatHist...) // copy: the goroutine outlives this Update
	ch, cmd := startChatStream(func(onDelta func(string)) (string, error) {
		return svc.ChatStream(context.Background(), ctxText, hist, onDelta)
	})
	m.streamCh = ch
	return m, cmd
}

// handleStreamChunk advances a streaming tutor reply.
func (m VaultModel) handleStreamChunk(msg streamChunkMsg) (tea.Model, tea.Cmd) {
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

func (m VaultModel) forwardToFocus(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.cmdMode {
		var cmd tea.Cmd
		m.cmdLine, cmd = m.cmdLine.Update(msg)
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

func (m *VaultModel) setFocus(p pane) tea.Cmd {
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

// rebuildSidebar groups notes by subject into selectable rows (id = note path).
func (m *VaultModel) rebuildSidebar() {
	bySubject := map[string][]core.NoteMeta{}
	var subjects []string
	for _, n := range m.notes {
		s := n.Subject
		if s == "" {
			s = "—"
		}
		if _, ok := bySubject[s]; !ok {
			subjects = append(subjects, s)
		}
		bySubject[s] = append(bySubject[s], n)
	}
	sort.Strings(subjects)

	var items []sidebarItem
	for _, s := range subjects {
		items = append(items, sidebarItem{title: s, header: true})
		for _, n := range bySubject[s] {
			items = append(items, sidebarItem{
				id:     n.Path,
				title:  n.Title,
				active: n.Path == m.current,
			})
		}
	}
	m.sidebar.setItems(items)
}

// --- layout & view ---

func (m *VaultModel) layout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	m.cmdLine.Width = m.width - 4
	m.contentH = m.height - 4
	if m.contentH < 1 {
		m.contentH = 1
	}
	borders := 6 // three bordered boxes, 2 columns each
	if m.sidebarCollapsed {
		borders = 4
	}
	contentW := m.width - borders
	if contentW < 3 {
		contentW = 3
	}
	// The configured split is the base; :compact / :wide shift it live.
	chatPct := clampRange(m.cfg.ChatPct(32)-m.editorBias, 15, 75)
	if m.sidebarCollapsed {
		m.sidebarW = 0
		m.chatW = clampMin(contentW*chatPct/100, 16)
		m.editorW = clampMin(contentW-m.chatW, 10)
	} else {
		m.sidebarW = clampMin(contentW*m.cfg.SidebarPct(22)/100, 12)
		m.chatW = clampMin(contentW*chatPct/100, 16)
		m.editorW = clampMin(contentW-m.sidebarW-m.chatW, 10)
	}

	m.sidebar.setSize(m.sidebarW, m.contentH)
	reserved := 1 + len(m.essayHeaderLines(m.editorW)) + len(m.backlinkFooterLines(m.editorW))
	m.editor.SetSize(m.editorW, max(1, m.contentH-reserved))
	m.chat.setSize(m.chatW, m.contentH)
}

func (m VaultModel) View() string {
	if m.width == 0 {
		return "starting…"
	}
	if m.width < 60 || m.height < 16 {
		return "Terminal too small — please enlarge to at least 60×16."
	}
	ed := m.box(paneEditor, m.editorW, m.contentH, m.editorPaneView(m.editorW))
	ch := m.box(paneChat, m.chatW, m.contentH, m.chat.view())
	var row string
	if m.sidebarCollapsed {
		row = lipgloss.JoinHorizontal(lipgloss.Top, ed, ch)
	} else {
		sb := m.box(paneSidebar, m.sidebarW, m.contentH, m.sidebar.view())
		row = lipgloss.JoinHorizontal(lipgloss.Top, sb, ed, ch)
	}
	frame := lipgloss.JoinVertical(lipgloss.Left, m.titleView(), row, m.statusView())
	return lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(frame)
}

func (m VaultModel) box(p pane, w, h int, content string) string {
	return borderStyle(m.focus == p).
		Width(w).Height(h).
		MaxWidth(w + 2).MaxHeight(h + 2).
		Render(content)
}

// essayHeaderLines wraps the active essay prompt to the pane width as pinned
// "> " lines shown above the answer buffer.
func (m VaultModel) essayHeaderLines(w int) []string {
	if !m.studyMode || strings.TrimSpace(m.studyPrompt) == "" {
		return nil
	}
	avail := w - 2
	if avail < 8 {
		avail = 8
	}
	lines := wrapWords("Essay: "+m.studyPrompt, avail)
	if len(lines) > maxPromptHeaderLines {
		lines = lines[:maxPromptHeaderLines]
		lines[maxPromptHeaderLines-1] += " …"
	}
	for i := range lines {
		lines[i] = "> " + lines[i]
	}
	return lines
}

func (m VaultModel) editorPaneView(w int) string {
	// Short titles only — the full essay prompt renders as pinned lines below.
	label := "No note open"
	if m.studyMode {
		label = "ESSAY · " + m.currentTitle
	} else if m.current != "" {
		label = "NOTE · " + m.currentTitle
	}
	if w > 0 && lipgloss.Width(label) > w-2 {
		label = truncate(label, max(1, w-2))
	}
	out := editorHeader.Width(w).Render(label)
	for _, ln := range m.essayHeaderLines(w) {
		out += "\n" + promptHeaderStyle.MaxWidth(w).Render(ln)
	}
	out += "\n" + m.editor.View()
	for _, ln := range m.backlinkFooterLines(w) {
		out += "\n" + ln
	}
	return out
}

// backlinkFooterLines renders the "↩ Linked mentions" panel under the editor for
// the open note (Obsidian-style). Empty in study mode, when toggled off via
// :backlinks, or when nothing links here.
func (m VaultModel) backlinkFooterLines(w int) []string {
	if m.studyMode || !m.showBacklinks || len(m.backlinks) == 0 {
		return nil
	}
	avail := max(8, w-2)
	lines := []string{backlinkHeaderStyle.Render(truncate("↩ Linked mentions ("+itoa(len(m.backlinks))+")", avail))}
	const maxShown = 5
	for i, n := range m.backlinks {
		if i >= maxShown {
			lines = append(lines, hintStyle.Render(truncate("  … +"+itoa(len(m.backlinks)-maxShown)+" more", avail)))
			break
		}
		title := n.Title
		if title == "" {
			title = n.Path
		}
		lines = append(lines, truncate("  • "+title, avail))
	}
	return lines
}

func (m VaultModel) titleView() string {
	t := "Meari — vault"
	if m.svc.Offline() {
		t += "  (offline)"
	}
	return titleBar.Width(m.width).Render(t)
}

func (m VaultModel) statusView() string {
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
	hints := "⌃w h·l focus · : cmds · :learn <topic> · enter open · ⌃s save · ⌃c quit"
	switch {
	case m.pendingWindow:
		hints = errStyle.Render("⌃w") + " window: h/l choose pane"
	case m.studyMode:
		hints = ":grade check answer · :answer see a model answer · :done finish"
	case m.focus == paneChat:
		hints = "enter send · ⌥o/:copy copy reply · ⌃f/⌃b scroll"
	}
	return statusBar.Width(m.width).Render(left + "   " + hintStyle.Render(hints))
}

func (m VaultModel) focusName() string {
	switch m.focus {
	case paneSidebar:
		return "notes"
	case paneEditor:
		if m.studyMode {
			return "answer"
		}
		return "editor"
	case paneChat:
		return "chat"
	}
	return ""
}

// vChatTurn builds a one-element tutor history slice (test/readability helper).
func vChatTurn(role, content string) []tutor.ChatTurn {
	return []tutor.ChatTurn{{Role: role, Content: content}}
}

// itoa is a tiny non-negative int formatter (avoids pulling strconv in here).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// --- async commands & messages ---

type (
	vNotesMsg     struct{ notes []core.NoteMeta }
	vOpenedMsg    struct{ note core.Note }
	vBacklinksMsg struct {
		path  string
		links []core.NoteMeta
	}
	vGeneratedMsg struct{ meta core.NoteMeta }
	vSavedMsg     struct{ meta core.NoteMeta }
	vEssayMsg     struct{ res core.EssayResult }
	vAnswerMsg    struct{ text string }
	vErrMsg       struct {
		kind string
		err  error
	}
)

func vListCmd(svc *core.Service) tea.Cmd {
	return func() tea.Msg {
		notes, err := svc.ListNotes()
		if err != nil {
			return vErrMsg{kind: "list", err: err}
		}
		return vNotesMsg{notes: notes}
	}
}

func vOpenCmd(svc *core.Service, path string) tea.Cmd {
	return func() tea.Msg {
		n, err := svc.OpenNote(path)
		if err != nil {
			return vErrMsg{kind: "open", err: err}
		}
		return vOpenedMsg{note: n}
	}
}

// vBacklinksCmd fetches the notes that link to path. Backlinks are advisory, so
// an error just yields an empty panel rather than a visible failure.
func vBacklinksCmd(svc *core.Service, path string) tea.Cmd {
	return func() tea.Msg {
		links, err := svc.Backlinks(path)
		if err != nil {
			return vBacklinksMsg{path: path, links: nil}
		}
		return vBacklinksMsg{path: path, links: links}
	}
}

func vGenCmd(svc *core.Service, request string) tea.Cmd {
	return func() tea.Msg {
		meta, err := svc.GenerateLesson(context.Background(), request)
		if err != nil {
			return vErrMsg{kind: "generate", err: err}
		}
		return vGeneratedMsg{meta: meta}
	}
}

func vSaveCmd(svc *core.Service, path, body string) tea.Cmd {
	return func() tea.Msg {
		meta, err := svc.SaveNote(path, body)
		if err != nil {
			return vErrMsg{kind: "save", err: err}
		}
		return vSavedMsg{meta: meta}
	}
}

// vSaveOpenCmd saves a new note then opens it.
func vSaveOpenCmd(svc *core.Service, path, body string) tea.Cmd {
	return func() tea.Msg {
		if _, err := svc.SaveNote(path, body); err != nil {
			return vErrMsg{kind: "save", err: err}
		}
		n, err := svc.OpenNote(path)
		if err != nil {
			return vErrMsg{kind: "open", err: err}
		}
		return vOpenedMsg{note: n}
	}
}

func vEssayCmd(svc *core.Service, prompt, answer string) tea.Cmd {
	return func() tea.Msg {
		res, err := svc.GradeEssay(context.Background(), prompt, answer)
		if err != nil {
			return vErrMsg{kind: "grade", err: err}
		}
		return vEssayMsg{res: res}
	}
}

func vAnswerCmd(svc *core.Service, prompt string) tea.Cmd {
	return func() tea.Msg {
		text, err := svc.ModelAnswer(context.Background(), prompt)
		if err != nil {
			return vErrMsg{kind: "answer", err: err}
		}
		return vAnswerMsg{text: text}
	}
}
