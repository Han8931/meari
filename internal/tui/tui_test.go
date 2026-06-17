package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/curriculum"
	"meari/internal/drafts"
	"meari/internal/editor"
	"meari/internal/executor"
	"meari/internal/progress"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// testDeps builds an offline TUI over temp dirs so tests need no network/python.
func testDeps(t *testing.T) Deps {
	t.Helper()
	dir := t.TempDir()
	store, err := drafts.New(dir)
	if err != nil {
		t.Fatalf("drafts.New: %v", err)
	}
	prog, err := progress.Load(dir)
	if err != nil {
		t.Fatalf("progress.Load: %v", err)
	}
	// No API key + non-ollama provider => offline tutor.
	tut := tutor.New(config.AIConfig{Provider: "openai", APIKeyEnv: "CODETUTOR_TEST_NO_KEY"})
	cfg := config.Default()
	cfg.WorkspaceDir, cfg.DataDir = dir, dir
	return Deps{Tutor: tut, Store: store, Progress: prog, Cfg: cfg}
}

// step applies a message and returns the concrete model, failing on panic.
func step(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	tm, _ := m.Update(msg)
	return tm.(Model)
}

func keyRunes(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

// dashboardSelect moves the dashboard cursor to the entry with the given kind
// and presses enter.
func dashboardSelect(t *testing.T, m Model, kind dashKind) Model {
	t.Helper()
	at := -1
	for i, e := range m.dash {
		if e.kind == kind {
			at = i
			break
		}
	}
	if at < 0 {
		t.Fatalf("no dashboard entry of kind %v in %+v", kind, m.dash)
	}
	m = step(t, m, keyRunes("g")) // cursor to the top
	for i := 0; i < at; i++ {
		m = step(t, m, keyRunes("j"))
	}
	return step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
}

func TestSetupDashboardCustomTopicFlow(t *testing.T) {
	m := newModel(testDeps(t))
	if m.phase != phaseSetup || m.setupStep != stepDashboard {
		t.Fatalf("expected the launch dashboard, got phase=%v step=%v", m.phase, m.setupStep)
	}

	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	// Pick "A topic of my own".
	m = dashboardSelect(t, m, dashTopic)
	if m.setupStep != stepTopic {
		t.Fatalf("step = %v, want stepTopic", m.setupStep)
	}
	// Topic: type and submit.
	m = step(t, m, keyRunes("loops"))
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.setupStep != stepLevel {
		t.Fatalf("step = %v, want stepLevel", m.setupStep)
	}
	// Level: Intermediate (second option).
	m = step(t, m, keyRunes("j"))
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if m.phase != phaseReady {
		t.Fatalf("expected phaseReady after the dashboard, got %v", m.phase)
	}
	if m.topic != "loops" {
		t.Fatalf("topic = %q, want \"loops\"", m.topic)
	}
	if m.level != "intermediate" {
		t.Fatalf("level = %q, want intermediate", m.level)
	}
	if m.curriculum {
		t.Fatal("custom-topic path should not be curriculum mode")
	}

	// Deliver async results that the dispatched commands would produce.
	m = step(t, m, lessonMsg{text: "Loops repeat work."})
	ch := tutor.Challenge{ID: "sum-list", Prompt: "Write sum_list(xs).", StarterCode: "def sum_list(xs):\n    pass", Tests: []string{"assert sum_list([1,2])==3"}}
	m = step(t, m, challengeMsg{ch: ch})

	if m.current.ID != "sum-list" {
		t.Fatalf("current challenge = %q, want sum-list", m.current.ID)
	}
	if got := m.editor.Value(); !strings.Contains(got, "def sum_list") {
		t.Fatalf("editor not loaded with starter code, got %q", got)
	}
	if it, ok := m.sidebar.selected(); !ok || it.id != "sum-list" {
		t.Fatalf("sidebar missing the new challenge, got %+v ok=%v", it, ok)
	}

	view := m.View()
	// Title bar, sidebar title, editor buffer, and chat transcript should all
	// be visible — one assertion per pane plus the chrome.
	for _, want := range []string{
		"Meari",          // title bar
		"loops",          // topic in title bar
		"Write sum_list", // sidebar item title (first line of prompt)
		"def sum_list",   // editor buffer (starter code)
		"Loops repeat",   // chat transcript (lesson)
		"editor",         // status bar focused-pane name
	} {
		if !strings.Contains(view, want) {
			t.Errorf("view missing %q", want)
		}
	}
}

func TestRunResultRecordsProgressAndChainsFeedback(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.topic = "loops"

	ch := tutor.Challenge{ID: "pass-id", Prompt: "do it", Tests: []string{"assert True"}}
	m = step(t, m, challengeMsg{ch: ch})

	// Simulate a passing run.
	tm, cmd := m.Update(runResultMsg{res: executor.Result{Passed: true, Output: ""}, ch: ch, code: "x=1"})
	m = tm.(Model)
	if cmd == nil {
		t.Fatal("expected a feedback command to be chained after a run")
	}
	if !m.deps.Progress.Done("pass-id") {
		t.Fatal("passing run should mark the challenge done in progress")
	}

	// The chained feedback command resolves to a feedbackMsg (offline canned text).
	msg := cmd()
	fb, ok := msg.(feedbackMsg)
	if !ok {
		t.Fatalf("chained cmd produced %T, want feedbackMsg", msg)
	}
	m = step(t, m, fb)
	if !strings.Contains(m.chat.view(), "tutor") {
		t.Errorf("feedback not shown in chat transcript:\n%s", m.chat.view())
	}
}

func TestRunFailureShowsStructuredChatSummary(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	ch := tutor.Challenge{ID: "fail-id", Prompt: "do it", Tests: []string{"assert is_even(4) == True"}}
	m = step(t, m, challengeMsg{ch: ch})

	out := "Traceback (most recent call last):\n  File \"solution.py\", line 8, in <module>\n    assert is_even(4) == True\nAssertionError"
	tm, cmd := m.Update(runResultMsg{res: executor.Result{Output: out}, ch: ch, code: "def is_even(n): return False"})
	m = tm.(Model)
	if cmd == nil {
		t.Fatal("expected a feedback command to be chained after a failing run")
	}

	view := m.chat.view()
	for _, want := range []string{"✗ Tests failed", "Failed:", "assert is_even(4) == True", "Output:"} {
		if !strings.Contains(view, want) {
			t.Fatalf("failure summary missing %q:\n%s", want, view)
		}
	}
	if got := m.deps.Progress.Challenges["fail-id"].Status; got != "in_progress" {
		t.Fatalf("failing run status = %q, want in_progress", got)
	}
}

func TestFailureSummaryTimeoutIncludesNextStep(t *testing.T) {
	got := failureSummary(executor.Result{TimedOut: true})
	for _, want := range []string{"Reason:", "execution timed out", "Try:", "loops"} {
		if !strings.Contains(got, want) {
			t.Fatalf("timeout summary missing %q:\n%s", want, got)
		}
	}
}

func TestTabDoesNotSwitchFocus(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.setFocus(paneEditor)
	// Tab is no longer a focus switch — it belongs to the focused pane (e.g.
	// indenting in the editor), so focus must stay put.
	m = step(t, m, tea.KeyMsg{Type: tea.KeyTab})
	m = step(t, m, tea.KeyMsg{Type: tea.KeyShiftTab})
	if m.focus != paneEditor {
		t.Fatalf("Tab should not change focus, got %v", m.focus)
	}
}

func TestCtrlWWindowChordSwitchesPanes(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	if m.focus != paneSidebar {
		t.Fatalf("expected initial focus on sidebar, got %v", m.focus)
	}

	// Ctrl-W then 'l' moves focus right (sidebar -> editor).
	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlW})
	if !m.pendingWindow {
		t.Fatal("Ctrl-W should arm a pending window command")
	}
	m = step(t, m, keyRunes("l"))
	if m.focus != paneEditor {
		t.Fatalf("after ⌃w l focus = %v, want paneEditor", m.focus)
	}
	if m.pendingWindow {
		t.Fatal("pending window command should be cleared after the second key")
	}

	// Ctrl-W l again -> chat; once more should clamp (no wrap).
	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlW})
	m = step(t, m, keyRunes("l"))
	if m.focus != paneChat {
		t.Fatalf("focus = %v, want paneChat", m.focus)
	}
	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlW})
	m = step(t, m, keyRunes("l"))
	if m.focus != paneChat {
		t.Fatalf("⌃w l at the right edge should clamp to paneChat, got %v", m.focus)
	}

	// Ctrl-W h moves back left.
	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlW})
	m = step(t, m, keyRunes("h"))
	if m.focus != paneEditor {
		t.Fatalf("after ⌃w h focus = %v, want paneEditor", m.focus)
	}
}

// testDepsSeeded is testDeps plus a vault service with the built-in Go track
// seeded as markdown courses — what a real first launch produces.
func testDepsSeeded(t *testing.T) Deps {
	t.Helper()
	d := testDeps(t)
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d.Svc = core.New(v, d.Tutor)
	d.Svc.SetCourseDir(t.TempDir())
	if err := d.Svc.SeedBuiltinCourses(); err != nil {
		t.Fatal(err)
	}
	return d
}

// Entering a seeded Go course from the dashboard starts it directly — courses
// carry their own level, so no follow-up question.
func TestSetupDashboardEntersSeededCourse(t *testing.T) {
	m := newModel(testDepsSeeded(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	m = dashboardSelect(t, m, dashCourse)
	if m.phase != phaseReady {
		t.Fatalf("phase = %v, want phaseReady (no level question for courses)", m.phase)
	}
	if !m.curriculum || len(m.curr.Modules) == 0 {
		t.Fatalf("seeded course should install a curriculum, got %+v", m.curr)
	}
	if m.current.Lang != "go" {
		t.Fatalf("first challenge lang = %q, want go", m.current.Lang)
	}
}

func TestSetupDashboardBackWithEsc(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m = dashboardSelect(t, m, dashTopic)         // -> stepTopic
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc}) // back -> dashboard
	if m.setupStep != stepDashboard {
		t.Fatalf("Esc should return to the dashboard, got %v", m.setupStep)
	}
	// The cursor lands back on the row that was chosen, not the top.
	if m.dash[m.setupCursor].kind != dashTopic {
		t.Fatalf("cursor should restore to the chosen row, got %+v", m.dash[m.setupCursor])
	}
}

func TestResumeOffersAndContinuesSavedSession(t *testing.T) {
	d := testDepsSeeded(t)
	// Simulate a prior session saved to disk: the 2nd topic of go-beginner.
	course, err := d.Svc.LoadCourse("go-beginner")
	if err != nil {
		t.Fatal(err)
	}
	second := course.Curriculum().Topics()[1]
	d.Progress.SetLast("go-beginner", "beginner", second.ID, second.Title)

	m := newModel(d)
	if m.setupStep != stepDashboard || len(m.dash) == 0 || m.dash[0].kind != dashContinue {
		t.Fatalf("with a saved session, the dashboard should lead with Continue, got %+v", m.dash)
	}
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	// Choose "Continue" (first row).
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.phase != phaseReady {
		t.Fatalf("phase = %v, want phaseReady", m.phase)
	}
	if !m.curriculum || m.curr.Lang != "go-beginner" {
		t.Fatalf("resumed into %q, want go-beginner", m.curr.Lang)
	}
	if m.currentTopicID != second.ID {
		t.Fatalf("resumed topic = %q, want %q", m.currentTopicID, second.ID)
	}
}

// A saved session whose course no longer exists must not produce a dead
// Continue row.
func TestStaleSavedSessionHidesContinue(t *testing.T) {
	d := testDepsSeeded(t)
	d.Progress.SetLast("python", "beginner", "py-b-vars", "Variables")
	m := newModel(d)
	for _, e := range m.dash {
		if e.kind == dashContinue {
			t.Fatalf("stale session should hide Continue, got %+v", e)
		}
	}
}

// The dashboard lists vault-built courses and entering one starts it directly
// (no level question — the course carries its own).
func TestSetupDashboardListsVaultCourses(t *testing.T) {
	d := testDeps(t)
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d.Svc = core.New(v, d.Tutor)
	if _, err := d.Svc.SaveNote("Algo/BST.md", "# BST\n\nOrdered keys.\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Svc.SaveNote(core.CourseDir+"/Trees/course.md", "## Basics\n- [[BST]]\n"); err != nil {
		t.Fatal(err)
	}

	m := newModel(d)
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	found := false
	for _, e := range m.dash {
		if e.kind == dashCourse {
			found = true
			if e.title != "Trees" && e.id == "" {
				t.Fatalf("course entry incomplete: %+v", e)
			}
		}
	}
	if !found {
		t.Fatalf("dashboard should list the vault course, got %+v", m.dash)
	}

	m = dashboardSelect(t, m, dashCourse)
	if m.phase != phaseReady {
		t.Fatalf("phase = %v, want phaseReady after entering a vault course", m.phase)
	}
	if !m.curriculum || len(m.curr.Modules) == 0 {
		t.Fatalf("vault course should install a curriculum, got %+v", m.curr)
	}
}

func TestDashboardViewRenders(t *testing.T) {
	d := testDepsSeeded(t)
	d.Progress.SetLast("go-beginner", "beginner", "course-go-beginner-hello-go", "Hello, Go")
	m := newModel(d)
	m = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})

	out := m.View()
	for _, want := range []string{
		"What will you learn today?",
		"Continue", "Hello, Go",
		"Courses", "Go (Beginner)", "Go (Advanced)",
		"A topic of my own", "Open the vault",
	} {
		if !strings.Contains(ansi.Strip(out), want) {
			t.Errorf("dashboard missing %q:\n%s", want, out)
		}
	}
	// Full-screen: as tall as the terminal, not a small card.
	if rows := strings.Count(out, "\n") + 1; rows < 24 {
		t.Errorf("dashboard should fill the screen, got %d rows", rows)
	}
}

// Choosing "Open the vault" leaves the tutor with the switch target set, so
// the shell loop opens the vault TUI.
func TestSetupDashboardOpensVault(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m = dashboardSelect(t, m, dashVault)
	if m.exit != SwitchToVault {
		t.Fatalf("exit = %v, want SwitchToVault", m.exit)
	}
}

func TestStartingTopicPersistsResumePoint(t *testing.T) {
	m := readyModel(t)
	if m.currentTopicID == "" {
		t.Fatal("setup: no current topic")
	}
	if m.deps.Progress.Last == nil || m.deps.Progress.Last.TopicID != m.currentTopicID {
		t.Fatalf("starting a topic should save it as the resume point, got %+v", m.deps.Progress.Last)
	}
}

func TestCurriculumModeStartsTracksAndSwitches(t *testing.T) {
	d := testDepsSeeded(t)
	d.Curriculum = true
	m := newModel(d)
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	if !m.curriculum {
		t.Fatal("expected curriculum mode")
	}
	if m.topic == "" || m.currentTopicID == "" {
		t.Fatalf("a starting topic should be selected (topic=%q id=%q)", m.topic, m.currentTopicID)
	}
	if m.deps.Progress.TopicStatus(m.currentTopicID) != "in_progress" {
		t.Fatalf("starting topic should be in_progress, got %q", m.deps.Progress.TopicStatus(m.currentTopicID))
	}

	// Sidebar should list module headers and selectable topics.
	var headers, topics int
	for _, it := range m.sidebar.items {
		if it.header {
			headers++
		} else {
			topics++
		}
	}
	if headers == 0 || topics == 0 {
		t.Fatalf("sidebar should have module headers and topics, got %d/%d", headers, topics)
	}

	// Passing a challenge completes the current topic.
	startID := m.currentTopicID
	ch := tutor.Challenge{ID: "ch-1", Tests: []string{"assert True"}}
	tm, _ := m.Update(runResultMsg{res: executor.Result{Passed: true}, ch: ch, code: "x"})
	m = tm.(Model)
	if m.deps.Progress.TopicStatus(startID) != "done" {
		t.Fatalf("passing a challenge should complete the topic, got %q", m.deps.Progress.TopicStatus(startID))
	}

	// Selecting another topic from the sidebar switches the active topic.
	m.setFocus(paneSidebar)
	m = step(t, m, keyRunes("G")) // jump to the last selectable topic
	sel, ok := m.sidebar.selected()
	if !ok || sel.header {
		t.Fatalf("expected a selectable topic under the cursor, got %+v", sel)
	}
	tm2, _ := m.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm2.(Model)
	wantTopicID := topicIDFromSidebar(sel.id)
	if m.currentTopicID != wantTopicID {
		t.Fatalf("active topic = %q, want %q", m.currentTopicID, wantTopicID)
	}
	// Baked content loads synchronously: the editor should hold the new topic's
	// starter code and the tutor should have its tests.
	if m.current.ID != wantTopicID {
		t.Fatalf("current challenge = %q, want %q", m.current.ID, wantTopicID)
	}
	if len(m.current.Tests) == 0 {
		t.Fatal("selected topic should have baked-in tests")
	}
}

// TestViewNeverExceedsScreen guards against the leftover-render bug: the frame
// must always fit exactly within the terminal, even when the sidebar (here, the
// full curriculum) has far more rows than fit.
func TestViewNeverExceedsScreen(t *testing.T) {
	d := testDeps(t)
	d.Curriculum = true
	m := newModel(d)

	for _, sz := range []struct{ w, h int }{{80, 24}, {120, 40}, {64, 16}} {
		mm := step(t, m, tea.WindowSizeMsg{Width: sz.w, Height: sz.h})
		mm = step(t, mm, challengeMsg{ch: tutor.Challenge{ID: "x", Prompt: "p", StarterCode: "def f():\n    pass"}})
		// Add a long chat line and a long editor status to stress wrapping.
		mm.chat.append(roleSystem, strings.Repeat("verylongword ", 40))

		out := mm.View()
		lines := strings.Split(out, "\n")
		if len(lines) > sz.h {
			t.Errorf("%dx%d: View produced %d lines, exceeds height %d", sz.w, sz.h, len(lines), sz.h)
		}
		for i, ln := range lines {
			if w := lipgloss.Width(ln); w > sz.w {
				t.Errorf("%dx%d: line %d width %d exceeds %d", sz.w, sz.h, i, w, sz.w)
			}
		}
	}
}

func TestHorizontalLayoutFitsScreen(t *testing.T) {
	d := testDeps(t)
	d.Cfg.UI.Layout = "horizontal"
	m := newModel(d)
	if !m.horizontal {
		t.Fatal("layout=horizontal should put the model in horizontal mode")
	}
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.phase = phaseReady
	m = step(t, m, challengeMsg{ch: tutor.Challenge{ID: "x", Prompt: "p", StarterCode: "def f():\n    pass"}})
	m.chat.append(roleSystem, strings.Repeat("long ", 60))

	out := m.View()
	lines := strings.Split(out, "\n")
	if len(lines) > 30 {
		t.Errorf("horizontal View produced %d lines, exceeds height 30", len(lines))
	}
	for i, ln := range lines {
		if w := lipgloss.Width(ln); w > 100 {
			t.Errorf("horizontal line %d width %d exceeds 100", i, w)
		}
	}
}

func TestConfigReloadTogglesLayout(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte("[ui]\nlayout = \"horizontal\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	d := testDeps(t)
	d.ConfigPath = cfgPath
	d.BaseDir = dir

	m := newModel(d) // starts vertical (testDeps default cfg)
	if m.horizontal {
		t.Fatal("expected the session to start vertical")
	}
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m.applyConfigReload(configReloadMsg{})
	if !m.horizontal {
		t.Fatal("reloading a config with layout=horizontal should switch the layout live")
	}
}

func TestChatQuestionRoundTrip(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.setFocus(paneChat)

	m.chat.input.SetValue("why a loop?")
	tm, cmd := m.submitChat()
	m = tm.(Model)
	if cmd == nil {
		t.Fatal("submitChat should dispatch a chat command")
	}
	if len(m.chatHist) != 1 || m.chatHist[0].Role != "user" {
		t.Fatalf("chat history not recording the user turn: %+v", m.chatHist)
	}
	if !m.streaming {
		t.Fatal("submitChat should mark a stream in flight")
	}

	// Drive the stream to completion: deltas grow the transcript, done
	// finalizes the conversation history. (Offline: one delta, then done.)
	for i := 0; i < 10; i++ {
		msg := cmd()
		chunk, ok := msg.(streamChunkMsg)
		if !ok {
			t.Fatalf("stream produced %T, want streamChunkMsg", msg)
		}
		var tm2 tea.Model
		tm2, cmd = m.Update(chunk)
		m = tm2.(Model)
		if chunk.done {
			break
		}
	}
	if m.streaming {
		t.Fatal("stream should be finished")
	}
	if len(m.chatHist) != 2 || m.chatHist[1].Role != "assistant" || m.chatHist[1].Content == "" {
		t.Fatalf("assistant turn not recorded: %+v", m.chatHist)
	}
	if !strings.Contains(m.chat.view(), "offline") {
		t.Fatalf("streamed reply should be visible in the transcript:\n%s", m.chat.view())
	}
}

func TestEscStopsStreamingTutorReply(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.setFocus(paneChat)
	m.streaming = true
	m.pending = 1
	m.streamCh = make(chan streamChunkMsg, 1)
	cancelled := false
	m.streamCancel = func() { cancelled = true }
	m.chat.beginStream()

	tm, cmd := m.handleKey(tea.KeyMsg{Type: tea.KeyEsc})
	m = tm.(Model)
	if cmd != nil {
		t.Fatal("esc while streaming should not quit or dispatch another command")
	}
	if !cancelled || !m.streamStopping || !m.streaming {
		t.Fatalf("stream not marked stopping: cancelled=%v stopping=%v streaming=%v", cancelled, m.streamStopping, m.streaming)
	}

	tm, cmd = m.Update(streamChunkMsg{delta: "late text"})
	m = tm.(Model)
	if !m.streamStopping || !m.streaming || cmd == nil {
		t.Fatalf("late delta should be ignored while listener stays armed: stopping=%v streaming=%v cmd=%v", m.streamStopping, m.streaming, cmd)
	}
	if strings.Contains(m.chat.view(), "late text") {
		t.Fatalf("late delta should not appear after stop:\n%s", m.chat.view())
	}

	tm, _ = m.Update(streamChunkMsg{done: true, full: "ignored"})
	m = tm.(Model)
	if m.streaming || m.streamStopping || m.pending != 0 {
		t.Fatalf("done after stop should clear stream state: streaming=%v stopping=%v pending=%d", m.streaming, m.streamStopping, m.pending)
	}
	if len(m.chatHist) != 0 {
		t.Fatalf("stopped reply should not be recorded as assistant history: %+v", m.chatHist)
	}
}

// fillChat appends enough blocks that the transcript overflows its viewport, so
// scrolling has somewhere to go.
func fillChat(m *Model) {
	for i := 0; i < 60; i++ {
		m.chat.append(roleTutor, "transcript line for scrolling")
	}
}

func TestChatTypingDoesNotScroll(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.setFocus(paneChat)
	fillChat(&m)
	if !m.chat.vp.AtBottom() {
		t.Fatal("transcript should start pinned to the bottom")
	}

	// "b" is bound to page-up in the viewport's default keymap; typed in the chat
	// box it must insert a character, not scroll the history.
	m = step(t, m, keyRunes("b"))
	if got := m.chat.input.Value(); got != "b" {
		t.Fatalf("typed key not inserted into input: got %q", got)
	}
	if !m.chat.vp.AtBottom() {
		t.Fatal("typing scrolled the transcript; it should stay at the bottom")
	}
}

func TestChatVimKeysScroll(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.setFocus(paneChat)
	fillChat(&m)

	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlU}) // half page up
	if m.chat.vp.AtBottom() {
		t.Fatal("ctrl+u did not scroll the transcript up")
	}
	if m.chat.input.Value() != "" {
		t.Fatalf("ctrl+u leaked into the input: %q", m.chat.input.Value())
	}

	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlD}) // half page down
	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlD})
	if !m.chat.vp.AtBottom() {
		t.Fatal("ctrl+d did not return to the bottom")
	}
}

func TestNewMessageKeepsScrollPosition(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	fillChat(&m)

	m.chat.vp.GotoTop() // reader scrolled up into the history
	m.chat.append(roleTutor, "a reply arrives while reading")
	if m.chat.vp.AtBottom() {
		t.Fatal("a new message yanked the reader to the bottom")
	}

	m.chat.vp.GotoBottom()
	m.chat.append(roleTutor, "another reply")
	if !m.chat.vp.AtBottom() {
		t.Fatal("a new message should follow the tail when already at the bottom")
	}
}

func TestWheelScrollsPaneUnderCursor(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.setFocus(paneEditor) // focus stays elsewhere; the wheel shouldn't move it
	fillChat(&m)

	m.sidebar.setItems([]sidebarItem{{id: "a", title: "A"}, {id: "b", title: "B"}, {id: "c", title: "C"}})

	// Wheel up over the chat column (rightmost) scrolls the transcript.
	chatX := m.sidebarW + m.editorW + 5
	m = step(t, m, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp, X: chatX, Y: 10})
	if m.chat.vp.AtBottom() {
		t.Fatal("wheel over the chat pane did not scroll the transcript")
	}
	if m.focus != paneEditor {
		t.Fatal("wheel scrolling should not steal focus")
	}

	// Wheel down over the sidebar column moves its selection (ranger-style).
	m = step(t, m, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown, X: 2, Y: 5})
	if m.sidebar.cursor != 1 {
		t.Fatalf("wheel over the sidebar should move the cursor, got %d", m.sidebar.cursor)
	}
}

func TestStartingTopicShowsLessonFromTop(t *testing.T) {
	m := readyModel(t)

	// A lesson far taller than the chat pane, so opening at the top vs. the
	// bottom is observable.
	lesson := strings.Repeat("A line of lesson content that the reader should see first.\n", 120)
	topic := curriculum.Topic{ID: "synthetic-long-lesson", Title: "Long Lesson", Lesson: lesson}

	_ = m.startTopic(topic)

	if m.chat.vp.TotalLineCount() <= m.chat.vp.Height {
		t.Fatalf("setup: lesson not tall enough to scroll (lines=%d height=%d)",
			m.chat.vp.TotalLineCount(), m.chat.vp.Height)
	}
	if !m.chat.vp.AtTop() || m.chat.vp.AtBottom() {
		t.Fatalf("a freshly started lesson should open at its top, got AtTop=%v AtBottom=%v",
			m.chat.vp.AtTop(), m.chat.vp.AtBottom())
	}
}

func TestLoadedLessonShowsFromTop(t *testing.T) {
	m := readyModel(t)

	// The chat starts pinned to the tail (as it is while the "loading lesson"
	// spinner runs), so a naive append would open the lesson at its end.
	m.chat.vp.GotoBottom()

	lesson := strings.Repeat("A line of lesson content that the reader should see first.\n", 120)
	m = step(t, m, lessonMsg{text: lesson})

	if m.chat.vp.TotalLineCount() <= m.chat.vp.Height {
		t.Fatalf("setup: lesson not tall enough to scroll (lines=%d height=%d)",
			m.chat.vp.TotalLineCount(), m.chat.vp.Height)
	}
	if !m.chat.vp.AtTop() || m.chat.vp.AtBottom() {
		t.Fatalf("a freshly loaded lesson should open at its top, got AtTop=%v AtBottom=%v",
			m.chat.vp.AtTop(), m.chat.vp.AtBottom())
	}
}

func TestTutorClickOpensTopic(t *testing.T) {
	m := readyModel(t)
	startID := m.currentTopicID

	// Find a selectable row that loads a DIFFERENT topic than the current one.
	target, wantID := -1, ""
	for i, it := range m.sidebar.items {
		if it.header {
			continue
		}
		if id := topicIDFromSidebar(it.id); id != startID {
			target, wantID = i, id
			break
		}
	}
	if target < 0 {
		t.Fatalf("no row for a different topic in sidebar: %+v", m.sidebar.items)
	}

	lo, _ := m.sidebar.window()
	y := target - lo + 2 // title bar (row 0) + box top border (row 1)
	tm, _ := m.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 2, Y: y})
	m = tm.(Model)

	// The click selects the clicked row and loads its topic — just like Enter
	// (opening then moves focus to chat/editor, so focus is not the signal here).
	if m.sidebar.cursor != target {
		t.Fatalf("click should select row %d, got %d", target, m.sidebar.cursor)
	}
	if m.currentTopicID != wantID {
		t.Fatalf("click should open topic %q, current is %q", wantID, m.currentTopicID)
	}
}

// readyModel returns a sized, post-setup model already in a Python beginner
// curriculum, so command tests can act on a live session.
func readyModel(t *testing.T) Model {
	t.Helper()
	m := newModel(testDepsSeeded(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m.loadCurriculum("go-beginner", curriculum.Beginner, "")
	return m
}

func ex(t *testing.T, m Model, raw string) Model {
	t.Helper()
	tm, _ := m.runEx(raw)
	return tm.(Model)
}

func TestCommandSwitchesCourse(t *testing.T) {
	m := readyModel(t)
	if m.lang != "go-beginner" {
		t.Fatalf("setup: lang=%q", m.lang)
	}
	m = ex(t, m, "topic go-intermediate")
	if !m.curriculum || m.lang != "go-intermediate" || m.current.Lang != "go" {
		t.Fatalf("topic go-intermediate did not switch course: curriculum=%v lang=%q chLang=%q", m.curriculum, m.lang, m.current.Lang)
	}

	// :subject is an alias; an unknown course is rejected without switching.
	m = ex(t, m, "subject go-advanced")
	if m.lang != "go-advanced" {
		t.Fatalf("subject alias failed: %q", m.lang)
	}
	m = ex(t, m, "topic ruby")
	if m.lang != "go-advanced" {
		t.Fatalf("unknown course should not switch, lang=%q", m.lang)
	}
}

func TestBareTopicOpensPicker(t *testing.T) {
	m := readyModel(t)
	m = ex(t, m, "topic")
	if m.overlay != overlayPicker {
		t.Fatalf("bare :topic should open the picker, overlay=%d", m.overlay)
	}
	ids, _ := m.pickerEntries()
	if ids[m.pickerCursor] != "go-beginner" {
		t.Fatalf("picker cursor should start on the current course, got %q", ids[m.pickerCursor])
	}
	// Move down and choose: it should switch to that course and close the modal.
	m = step(t, m, keyRunes("j"))
	chosen := ids[m.pickerCursor]
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.overlay != overlayNone || m.lang != chosen {
		t.Fatalf("picker enter should switch to %q and close, got lang=%q overlay=%d", chosen, m.lang, m.overlay)
	}
}

func TestClearChatTranscript(t *testing.T) {
	m := readyModel(t)
	if len(m.chat.blocks) == 0 {
		t.Fatal("expected lesson/challenge blocks after setup")
	}
	m.chatHist = append(m.chatHist, tutor.ChatTurn{Role: "user", Content: "hi"})
	m = ex(t, m, "clear")
	if len(m.chat.blocks) != 0 || len(m.chatHist) != 0 {
		t.Fatalf("clear did not empty the transcript: %d blocks, %d hist", len(m.chat.blocks), len(m.chatHist))
	}
}

func TestClearProgressRequiresConfirmation(t *testing.T) {
	m := readyModel(t)
	id := m.curr.Topics()[0].ID
	m.deps.Progress.MarkTopicDone(id)

	m = ex(t, m, "clear progress")
	if m.overlay != overlayConfirm || m.confirmKind != confirmClearProgress {
		t.Fatalf("clear progress should open a confirm modal, overlay=%d", m.overlay)
	}
	if m.deps.Progress.TopicStatus(id) != "done" {
		t.Fatal("progress wiped before confirmation")
	}

	// 'n' cancels and keeps progress.
	m = step(t, m, keyRunes("n"))
	if m.overlay != overlayNone || m.deps.Progress.TopicStatus(id) != "done" {
		t.Fatal("cancel should close the modal and keep progress")
	}

	// Re-open and confirm with 'y'.
	m = ex(t, m, "clear progress")
	m = step(t, m, keyRunes("y"))
	if m.overlay != overlayNone {
		t.Fatal("y should close the modal")
	}
	if m.deps.Progress.TopicStatus(id) == "done" {
		t.Fatal("confirm should clear progress")
	}
}

func TestClearDraftsConfirmRemovesFiles(t *testing.T) {
	m := readyModel(t)
	if err := m.deps.Store.Save("py-b-vars", "x = 1"); err != nil {
		t.Fatalf("seed draft: %v", err)
	}
	m = ex(t, m, "clear drafts")
	if m.confirmKind != confirmClearDrafts {
		t.Fatalf("clear drafts should arm the drafts confirm, got %d", m.confirmKind)
	}
	m = step(t, m, keyRunes("y"))
	if m.deps.Store.Has("py-b-vars") {
		t.Fatal("confirm should delete saved drafts")
	}
}

func TestProgressSummaryListsCourses(t *testing.T) {
	m := readyModel(t)
	m = ex(t, m, "progress")
	if m.overlay != overlayProgress {
		t.Fatalf("progress overlay not open, overlay=%d", m.overlay)
	}
	view := m.progressView()
	for _, want := range []string{"Go (Beginner)", "Go (Intermediate)", "Go (Advanced)"} {
		if !strings.Contains(view, want) {
			t.Fatalf("progress view missing %q:\n%s", want, view)
		}
	}
	// Esc closes it.
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if m.overlay != overlayNone {
		t.Fatal("esc should close the progress overlay")
	}
}

func TestColonOpensCommandLineFromSidebar(t *testing.T) {
	m := readyModel(t)
	m.setFocus(paneSidebar)
	m = step(t, m, keyRunes(":"))
	if !m.cmdMode {
		t.Fatal(": from the sidebar should open the command line")
	}
	for _, r := range "progress" {
		m = step(t, m, keyRunes(string(r)))
	}
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.cmdMode {
		t.Fatal("enter should close the command line")
	}
	if m.overlay != overlayProgress {
		t.Fatalf("typed :progress should open the overlay, overlay=%d", m.overlay)
	}
}

func TestHelpCommandOpensModal(t *testing.T) {
	m := readyModel(t)
	m = ex(t, m, "help")
	if m.overlay != overlayHelp {
		t.Fatalf(":help should open the help overlay, overlay=%d", m.overlay)
	}
	view := m.overlayView()
	for _, want := range []string{":topic", ":clear", ":progress", "focus"} {
		if !strings.Contains(view, want) {
			t.Fatalf("help view missing %q:\n%s", want, view)
		}
	}
	m = step(t, m, keyRunes("q"))
	if m.overlay != overlayNone {
		t.Fatal("q should close the help overlay")
	}
}

func TestFoldHidesSidebarAndReclaimsWidth(t *testing.T) {
	m := readyModel(t)
	if tpc, ok := m.topicByID[m.currentTopicID]; ok {
		_ = m.startTopicView(tpc, "quiz")
	}
	m.setFocus(paneSidebar)
	wantTitle := firstLine(m.curr.Modules[0].Name)
	if !strings.Contains(m.View(), wantTitle) {
		t.Fatalf("sidebar module %q should be visible before folding", wantTitle)
	}
	sidebarW, editorW := m.sidebarW, m.editorW

	// :fold collapses the pane, moves focus off it, and widens the editor.
	m = ex(t, m, "fold")
	if !m.sidebarCollapsed {
		t.Fatal(":fold should collapse the sidebar")
	}
	if m.focus == paneSidebar {
		t.Fatal("folding the focused sidebar should move focus away")
	}
	if m.sidebarW != 0 {
		t.Fatalf("folded sidebar width = %d, want 0", m.sidebarW)
	}
	if m.editorW <= editorW {
		t.Fatalf("editor should reclaim width when folded: %d -> %d", editorW, m.editorW)
	}
	if strings.Contains(m.View(), wantTitle) {
		t.Fatal("folded sidebar should not render its module headers")
	}

	// Ctrl-W h must not return focus to the hidden pane (clamps at the editor).
	m.setFocus(paneEditor)
	m = step(t, m, tea.KeyMsg{Type: tea.KeyCtrlW})
	m = step(t, m, keyRunes("h"))
	if m.focus == paneSidebar {
		t.Fatal("⌃w h should not focus a folded sidebar")
	}

	// :fold again restores the original geometry.
	m = ex(t, m, "fold")
	if m.sidebarCollapsed {
		t.Fatal("second :fold should unfold the sidebar")
	}
	if m.sidebarW != sidebarW || m.editorW != editorW {
		t.Fatalf("unfold should restore widths: sidebar %d->%d, editor %d->%d", sidebarW, m.sidebarW, editorW, m.editorW)
	}
	if !strings.Contains(m.View(), wantTitle) {
		t.Fatal("unfolded sidebar should render again")
	}
}

func TestCompactAndWideResizeEditor(t *testing.T) {
	m := readyModel(t)
	if tpc, ok := m.topicByID[m.currentTopicID]; ok {
		_ = m.startTopicView(tpc, "quiz")
	}
	baseEditor, baseChat := m.editorW, m.chatW

	// :compact shrinks the editor and widens the chat (the whole point).
	m = ex(t, m, "compact")
	if m.editorBias >= 0 {
		t.Fatalf(":compact should bias toward the chat, got %d", m.editorBias)
	}
	if m.editorW >= baseEditor {
		t.Fatalf(":compact should narrow the editor: %d -> %d", baseEditor, m.editorW)
	}
	if m.chatW <= baseChat {
		t.Fatalf(":compact should widen the chat: %d -> %d", baseChat, m.chatW)
	}
	// Total width must still fit (sidebar + editor + chat + 6 borders <= width).
	if tot := m.sidebarW + m.editorW + m.chatW + 6; tot > m.width {
		t.Fatalf("panes overflow after :compact: %d > %d", tot, m.width)
	}

	// :wide reverses it; two :wide from one :compact lands net wider than default.
	m = ex(t, m, "wide")
	if m.editorBias != 0 {
		t.Fatalf("one :compact then one :wide should return to default, got %d", m.editorBias)
	}
	if m.editorW != baseEditor || m.chatW != baseChat {
		t.Fatalf("split not restored: editor %d->%d chat %d->%d", baseEditor, m.editorW, baseChat, m.chatW)
	}

	// Repeated :wide clamps instead of letting the chat vanish.
	for i := 0; i < 10; i++ {
		m = ex(t, m, "wide")
	}
	if m.editorBias != editorBiasMax {
		t.Fatalf("editorBias should clamp at %d, got %d", editorBiasMax, m.editorBias)
	}
	if m.chatW < 16 {
		t.Fatalf("chat pane fell below its floor: %d", m.chatW)
	}
}

func TestEditorForwardsGlobalCommand(t *testing.T) {
	m := readyModel(t)
	m = step(t, m, editor.RunCommandMsg{Raw: "topic go-advanced"})
	if m.lang != "go-advanced" {
		t.Fatalf("editor-forwarded command did not switch course: %q", m.lang)
	}
}
