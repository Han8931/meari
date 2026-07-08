package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/curriculum"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// courseDeps is testDeps plus a vault service holding one meari-course with a
// code topic and an essay topic.
func courseDeps(t *testing.T) Deps {
	t.Helper()
	d := testDeps(t)
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d.Svc = core.New(v, tutor.New(config.AIConfig{Provider: "openai"}))

	if _, err := v.Write(vault.Note{
		RelPath: "Algo/BST.md", Title: "BST",
		Body: "A binary search tree keeps keys ordered.\n",
		Extra: map[string]any{"study": map[string]any{
			"kind": "code", "lang": "python",
			"prompt":  "Implement insert(root, val).",
			"starter": "def insert(root, val):\n    pass\n",
			"tests":   []any{"assert True"},
			"answer":  "def insert(root, val): ...",
		}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Write(vault.Note{
		RelPath: "Algo/AVL.md", Title: "AVL",
		Body: "A self-balancing BST.\n",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Write(vault.Note{
		RelPath: core.CourseDir + "/Balanced Trees/course.md",
		ID:      "balanced-trees", Title: "Balanced Trees",
		Extra: map[string]any{"level": "intermediate"},
		Body:  "## Trees\n- [[BST]]\n- [[AVL]]\n",
	}); err != nil {
		t.Fatal(err)
	}
	return d
}

// :course offline goes straight to the defaults pipeline and produces a
// loadable course, streaming progress into the chat pane.
func TestVaultCourseCommandOffline(t *testing.T) {
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := core.New(v, tutor.New(config.AIConfig{Provider: "openai"})) // offline
	if _, err := v.Write(vault.Note{
		RelPath: "Git/Branching.md", Title: "Branching",
		Body: "Branches diverge history. See [[Merging]].\n",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Write(vault.Note{
		RelPath: "Git/Merging.md", Title: "Merging", Body: "Merging joins branches.\n",
	}); err != nil {
		t.Fatal(err)
	}

	m := newVaultModel(svc, config.Config{})
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = tm.(VaultModel)

	// :course without an open note refuses.
	tm, _ = m.runEx("course")
	m = tm.(VaultModel)
	if m.courseCh != nil {
		t.Fatal(":course without a note should refuse")
	}

	// Open the seed note, then :course (offline → defaults, no intake).
	opened := vOpenCmd(svc, "Git/Branching.md")().(vOpenedMsg)
	tm, _ = m.Update(opened)
	m = tm.(VaultModel)
	tm, cmd := m.runEx("course")
	m = tm.(VaultModel)
	if m.courseIntake {
		t.Fatal("offline :course must skip the intake")
	}
	if m.courseCh == nil || cmd == nil {
		t.Fatal("offline :course should start the build")
	}

	// Pump the progress channel until completion.
	progress := 0
	for i := 0; i < 100; i++ {
		msg := cmd()
		if done, ok := msg.(vCourseDoneMsg); ok {
			if done.err != nil {
				t.Fatalf("build failed: %v", done.err)
			}
			tm, _ = m.Update(msg)
			m = tm.(VaultModel)
			break
		}
		progress++
		tm, cmd = m.Update(msg)
		m = tm.(VaultModel)
		if cmd == nil {
			t.Fatal("progress handling lost the listener")
		}
	}
	if progress < 2 {
		t.Fatalf("expected progress lines, got %d", progress)
	}
	if m.courseCh != nil {
		t.Fatal("courseCh should be cleared after completion")
	}
	metas, err := svc.ListCourses()
	if err != nil || len(metas) != 1 {
		t.Fatalf("ListCourses = %+v, %v", metas, err)
	}
	if _, err := svc.LoadCourse(metas[0].ID); err != nil {
		t.Fatalf("generated course does not load: %v", err)
	}
}

// Esc abandons a pending intake without touching the normal chat state.
func TestVaultCourseIntakeEscCancels(t *testing.T) {
	m := newTestVaultModel(t)
	m.courseIntake = true
	m.courseHist = []tutor.ChatTurn{{Role: "user", Content: "hi"}}
	m.setFocus(paneChat)

	tm, _ := m.handleKey(tea.KeyMsg{Type: tea.KeyEsc})
	m = tm.(VaultModel)
	if m.courseIntake || m.courseHist != nil {
		t.Fatalf("esc should cancel the intake: intake=%v hist=%v", m.courseIntake, m.courseHist)
	}
}

func TestVaultCourseRunsInTutor(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady

	// :topic balanced-trees enters the vault course.
	tm, _ := m.runEx("topic balanced-trees")
	m = tm.(Model)
	if !m.curriculum || m.lang != "balanced-trees" || m.level != "intermediate" {
		t.Fatalf("course not installed: curriculum=%v lang=%q level=%q", m.curriculum, m.lang, m.level)
	}

	// First topic is the python code challenge: starter loaded, lang python.
	if m.current.Lang != "python" || !strings.Contains(m.editor.Value(), "def insert") {
		t.Fatalf("code topic wrong: lang=%q editor=%q", m.current.Lang, m.editor.Value())
	}
	if !strings.Contains(m.current.Prompt, "insert") {
		t.Fatalf("prompt = %q", m.current.Prompt)
	}
	// The default course row opens the lesson only; the quiz/task prompt lives
	// in its own child row. The lecture renders in the lesson pane (split
	// view), never in the conversation transcript.
	lesson := false
	for _, b := range m.lesson.blocks {
		if b.role == roleLesson && strings.Contains(b.text, "keeps keys ordered") {
			lesson = true
		}
	}
	if !lesson {
		t.Fatalf("lesson row did not show lecture in the lesson pane, got %+v", m.lesson.blocks)
	}
	for _, b := range m.chat.snapshot() {
		if b.role == roleLesson || b.role == roleQuiz || strings.Contains(b.text, "Implement insert") {
			t.Fatalf("lesson row must keep lecture/quiz out of the chat, got %+v", m.chat.snapshot())
		}
	}

	// Advance to the essay topic (no study block on the AVL note): it takes
	// the prose path — "essay" lang, default prompt and starter.
	topics := m.curr.Topics()
	if len(topics) != 2 {
		t.Fatalf("topics = %d, want 2", len(topics))
	}
	tm2, _ := m.startTopic(topics[1]), error(nil)
	_ = tm2
	if m.topicByID[topics[1].ID].Lang != "essay" {
		t.Fatalf("essay topic lang = %q", topics[1].Lang)
	}

	// The picker lists the vault course after the built-ins.
	ids, labels := m.pickerEntries()
	if ids[len(ids)-1] != "balanced-trees" || !strings.Contains(labels[len(labels)-1], "Balanced Trees") {
		t.Fatalf("picker entries = %v / %v", ids, labels)
	}

	// A unique substring finds the course too — by id fragment or title words.
	for _, q := range []string{"topic balanced", "topic Balanced Trees"} {
		m3 := newModel(courseDeps(t))
		m3 = step(t, m3, tea.WindowSizeMsg{Width: 100, Height: 30})
		m3.phase = phaseReady
		tm3, _ := m3.runEx(q)
		m3 = tm3.(Model)
		if m3.lang != "balanced-trees" {
			t.Fatalf("%q: lang = %q, want balanced-trees", q, m3.lang)
		}
	}

	// Tab completion: ":topic bal⇥" completes to the full course id.
	m.cmdLine.SetValue("topic bal")
	got := ""
	if s, ok := m.cmdComp.Next(m.cmdLine.Value(), m.exCandidates(), 1); ok {
		got = s
	}
	if got != "topic balanced-trees" {
		t.Fatalf("arg completion = %q", got)
	}

	// Resume path: loadCurriculum with the course id (what SetLast stored).
	m2 := newModel(courseDeps(t))
	m2 = step(t, m2, tea.WindowSizeMsg{Width: 100, Height: 30})
	m2.phase = phaseReady
	_ = m2.loadCurriculum("balanced-trees", "intermediate", "")
	if !m2.curriculum || m2.lang != "balanced-trees" || m2.current.ID == "" {
		t.Fatalf("resume failed: curriculum=%v lang=%q topic=%q", m2.curriculum, m2.lang, m2.current.ID)
	}
}

func TestExportChat(t *testing.T) {
	c := newChat()
	c.append(roleUser, "what is a BST?")
	c.append(roleTutor, "an ordered binary tree")
	dir := t.TempDir()

	msg := exportChat(&c, dir, "BST Note")
	if !strings.Contains(msg, "exported chat to ") {
		t.Fatalf("export message = %q", msg)
	}
	path := strings.TrimPrefix(msg, "exported chat to ")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "ordered binary tree") {
		t.Fatalf("export content = %q", b)
	}
	if !strings.Contains(filepath.Base(path), "chat-bst-note-") {
		t.Fatalf("export filename = %q", filepath.Base(path))
	}

	empty := newChat()
	if msg := exportChat(&empty, dir, "x"); !strings.Contains(msg, "nothing to export") {
		t.Fatalf("empty export message = %q", msg)
	}
}

// The EDITOR's ":" line completes :topic arguments too, via the parent's
// WithArgCompleter hook (typing ":" in the center pane is the common path).
func TestEditorTopicArgCompletion(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady

	// Focus the editor, open its command line, type "topic bal", press Tab.
	m.setFocus(paneEditor)
	for _, k := range []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune(":")},
		{Type: tea.KeyRunes, Runes: []rune("topic bal")},
		{Type: tea.KeyTab},
	} {
		m = step(t, m, k)
	}
	if got := m.editor.CmdLineValue(); got != "topic balanced-trees" {
		t.Fatalf("editor arg completion = %q, want %q", got, "topic balanced-trees")
	}
}

// Essay topics open chat-centric: no editor pane, prompt in the transcript,
// the chat input answers via :submit, and :view overrides both ways.
func TestChatCentricViewForEssayTopics(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	tm, _ := m.runEx("topic balanced-trees")
	m = tm.(Model)

	// Topic 1 opens on its lecture row: lectures are chat-only.
	if !m.chatCentric() || m.editorW != 0 {
		t.Fatalf("lecture row should use chat view: centric=%v editorW=%d", m.chatCentric(), m.editorW)
	}

	// Move to the essay topic's reflection row: it uses the editor/code view,
	// while the transcript shows only the reflection prompt.
	_ = m.startTopicView(m.curr.Topics()[1], "quiz")
	if m.chatCentric() || m.editorW == 0 {
		t.Fatalf("reflection row should use code view: centric=%v editorW=%d", m.chatCentric(), m.editorW)
	}
	if m.focus != paneEditor {
		t.Fatalf("focus should land on editor, got %v", m.focus)
	}
	task := false
	for _, b := range m.chat.snapshot() {
		if strings.Contains(b.text, "A self-balancing BST") {
			t.Fatalf("reflection row should not include the lesson: %+v", m.chat.snapshot())
		}
		if b.role == roleQuiz &&
			strings.Contains(b.text, "Explain the key ideas") {
			task = true
		}
	}
	if !task {
		t.Fatal("the prompt should be pinned into the transcript")
	}

	// :submit with the untouched starter refuses; edited text is graded from the editor.
	if cmd := m.startRun(); cmd != nil {
		t.Fatal("untouched reflection starter should not run")
	}
	m.editor.SetValue("A BST orders keys so lookup halves the space each step.")
	cmd := m.startRun()
	if cmd == nil {
		t.Fatal("edited answer should run")
	}
	res, ok := cmd().(runResultMsg)
	if !ok || !res.res.Passed { // essay reflection: non-empty passes
		t.Fatalf("essay submit result = %+v ok=%v", res, ok)
	}
	if res.code == "" || !strings.Contains(res.code, "BST orders keys") {
		t.Fatalf("answer not taken from the editor: %q", res.code)
	}

	// :view chat forces the lecture-style chat view; :view auto returns quiz rows to code view.
	tm, _ = m.runEx("view chat")
	m = tm.(Model)
	if !m.chatCentric() || m.editorW != 0 {
		t.Fatalf(":view chat should hide the editor, editorW=%d", m.editorW)
	}
	tm, _ = m.runEx("view auto")
	m = tm.(Model)
	if m.chatCentric() || m.editorW == 0 {
		t.Fatal(":view auto should follow the quiz row back to code")
	}
}

func TestQuizTopicRendersBelowLecture(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	m.curriculum = true
	m.curr = curriculum.Curriculum{Lang: "quiz-course", Level: curriculum.Beginner}
	m.topicByID = map[string]curriculum.Topic{}

	topic := curriculum.Topic{
		ID:     "quiz-one",
		Title:  "Ordering",
		Lesson: "Read this lecture note first.",
		Lang:   "quiz",
		Quiz: []curriculum.QuizQuestion{{
			Q:       "What should come after the lecture?",
			Choices: []string{"A quiz", "Another title"},
			Answer:  0,
		}},
	}
	_ = m.startTopicView(topic, "quiz")
	if m.chatCentric() || m.editorW == 0 || m.focus != paneEditor {
		t.Fatalf("quiz topic should use code view: centric=%v editorW=%d focus=%v", m.chatCentric(), m.editorW, m.focus)
	}
	blocks := m.chat.snapshot()
	if len(blocks) != 1 || blocks[0].role != roleQuiz {
		t.Fatalf("quiz row should show only quiz, got %+v", blocks)
	}
	if strings.Contains(blocks[0].text, "Read this lecture note first.") ||
		!strings.Contains(blocks[0].text, "1. What should come after the lecture?") ||
		!strings.Contains(blocks[0].text, "A. A quiz") ||
		!strings.Contains(blocks[0].text, "B. Another title") {
		t.Fatalf("quiz choices not formatted below lecture: %+v", blocks)
	}
	m.editor.SetValue("1A")
	cmd := m.startRun()
	if cmd == nil {
		t.Fatal("typed quiz answer should run")
	}
	res, ok := cmd().(runResultMsg)
	if !ok || !res.res.Passed {
		t.Fatalf("quiz submit result = %+v ok=%v", res, ok)
	}
}

func TestCourseSidebarShowsLessonAndQuizRows(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	tm, _ := m.runEx("topic balanced-trees")
	m = tm.(Model)

	var lesson, quiz sidebarItem
	for _, it := range m.sidebar.items {
		switch it.id {
		case topicLessonSidebarID(m.currentTopicID):
			lesson = it
		case topicQuizSidebarID(m.currentTopicID):
			quiz = it
		}
	}
	if lesson.id == "" || quiz.id == "" {
		t.Fatalf("sidebar should include lesson and quiz rows, got %+v", m.sidebar.items)
	}
	if lesson.title != "BST" || lesson.depth != 0 {
		t.Fatalf("lesson row wrong: %+v", lesson)
	}
	if quiz.title != "Quiz" || quiz.depth != 1 {
		t.Fatalf("quiz row wrong: %+v", quiz)
	}

	m.sidebar.cursor = 0
	for i, it := range m.sidebar.items {
		if it.id == quiz.id {
			m.sidebar.cursor = i
			break
		}
	}
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(Model)
	if m.currentTopicID != topicIDFromSidebar(quiz.id) {
		t.Fatalf("selecting quiz row should open its topic, got %q from %q", m.currentTopicID, quiz.id)
	}
}

// :export runs through runEx in both TUIs and writes a real file.
func TestExportCommandBothTUIs(t *testing.T) {
	d := courseDeps(t)
	d.Cfg.ExportsDir = filepath.Join(t.TempDir(), "exports")
	m := newModel(d)
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	m.chat.append(roleTutor, "tutor says hello")
	tm, _ := m.runEx("export")
	m = tm.(Model)
	if !strings.Contains(m.notice, "exported chat to ") {
		t.Fatalf("tutor :export notice = %q", m.notice)
	}
	if _, err := os.Stat(strings.TrimPrefix(m.notice, "exported chat to ")); err != nil {
		t.Fatal(err)
	}

	vm := newTestVaultModel(t)
	vm.cfg.ExportsDir = filepath.Join(t.TempDir(), "exports")
	vm.chat.append(roleTutor, "vault tutor reply")
	tvm, _ := vm.runEx("export")
	vm = tvm.(VaultModel)
	if !strings.Contains(vm.notice, "exported chat to ") {
		t.Fatalf("vault :export notice = %q", vm.notice)
	}
}

// Esc puts the chat input into Normal mode: j/k scroll, ":" opens the command
// line, i returns to typing, Enter doesn't send.
func TestChatInputNormalMode(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	m.setFocus(paneChat)

	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if !m.chat.inNormal() {
		t.Fatal("esc should enter Normal mode")
	}
	// Enter must not submit while in Normal mode.
	m.chat.input.SetValue("draft stays")
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.chat.input.Value() != "draft stays" {
		t.Fatalf("enter in Normal mode consumed the draft: %q", m.chat.input.Value())
	}
	// ":" opens the global command line.
	m = step(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	if !m.cmdMode {
		t.Fatal("':' in Normal mode should open the command line")
	}
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc}) // close cmdline
	// i returns to Insert; typing reaches the input again.
	m = step(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")})
	if m.chat.inNormal() {
		t.Fatal("'i' should return to Insert mode")
	}
	m = step(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if got := m.chat.input.Value(); !strings.Contains(got, "x") {
		t.Fatalf("typing after 'i' should reach the input: %q", got)
	}
}

// Normal mode in the chat input edits the DRAFT vim-style.
func TestChatInputNormalModeEditsDraft(t *testing.T) {
	m := newModel(courseDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	m.setFocus(paneChat)

	m.chat.input.SetValue("hello world")
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc}) // Normal
	press := func(s string) {
		m = step(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)})
	}
	press("0") // line start
	press("x") // delete "h"
	if got := m.chat.input.Value(); got != "ello world" {
		t.Fatalf("0x: %q", got)
	}
	press("w") // word forward
	press("D") // delete to end
	if got := m.chat.input.Value(); got != "ello" {
		t.Fatalf("wD: %q", got)
	}
	press("d")
	press("d") // dd clears the line
	if got := m.chat.input.Value(); got != "" {
		t.Fatalf("dd: %q", got)
	}
	press("i") // back to Insert
	if m.chat.inNormal() {
		t.Fatal("i should re-enter Insert")
	}

	// :q from the chat's Normal mode quits via the global command line.
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	m = step(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	if !m.cmdMode {
		t.Fatal("':' should open the command line")
	}
	m.cmdLine.SetValue("q")
	tm, cmd := m.updateCmdLine(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(Model)
	if cmd == nil {
		t.Fatal(":q should produce a quit command")
	}
}
