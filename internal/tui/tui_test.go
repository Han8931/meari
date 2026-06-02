package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"meari/internal/config"
	"meari/internal/drafts"
	"meari/internal/executor"
	"meari/internal/progress"
	"meari/internal/tutor"
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

func TestSetupWizardCustomTopicFlow(t *testing.T) {
	m := newModel(testDeps(t))
	if m.phase != phaseSetup || m.setupStep != stepLanguage {
		t.Fatalf("expected setup wizard at language step, got phase=%v step=%v", m.phase, m.setupStep)
	}

	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	// Language: Python (enter on first option).
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.setupStep != stepPath {
		t.Fatalf("after choosing Python, step = %v, want stepPath", m.setupStep)
	}
	// Path: move to "specific topic" and choose it.
	m = step(t, m, keyRunes("j"))
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
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
		t.Fatalf("expected phaseReady after the wizard, got %v", m.phase)
	}
	if m.topic != "python loops" {
		t.Fatalf("topic = %q, want \"python loops\"", m.topic)
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
		"Meari",       // title bar
		"loops",           // topic in title bar
		"Write sum_list",  // sidebar item title (first line of prompt)
		"def sum_list",    // editor buffer (starter code)
		"Loops repeat",    // chat transcript (lesson)
		"editor",          // status bar focused-pane name
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

func TestFocusCyclesAcrossPanes(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	m = step(t, m, tea.KeyMsg{Type: tea.KeyTab})
	m = step(t, m, tea.KeyMsg{Type: tea.KeyTab})
	// From the default (sidebar=0), two tabs land on chat (2).
	if m.focus != paneChat {
		t.Fatalf("after two tabs focus = %v, want paneChat", m.focus)
	}
	if !m.chat.focused {
		t.Error("chat pane should report focused after tabbing to it")
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

func TestSetupWizardCurriculumPath(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter}) // Language: Python -> stepPath
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter}) // Path: full curriculum -> stepLevel
	if m.setupStep != stepLevel {
		t.Fatalf("step = %v, want stepLevel", m.setupStep)
	}
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter}) // Level: Beginner -> ready

	if m.phase != phaseReady {
		t.Fatalf("phase = %v, want phaseReady", m.phase)
	}
	if !m.curriculum {
		t.Fatal("choosing 'full curriculum' should enable curriculum mode")
	}
	if m.level != "beginner" || m.topic == "" {
		t.Fatalf("level=%q topic=%q after wizard", m.level, m.topic)
	}
}

func TestSetupWizardBackWithEsc(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter}) // -> stepPath
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEsc})   // back -> stepLanguage
	if m.setupStep != stepLanguage {
		t.Fatalf("Esc should return to the language step, got %v", m.setupStep)
	}
}

func TestResumeOffersAndContinuesSavedSession(t *testing.T) {
	d := testDeps(t)
	// Simulate a prior session saved to disk.
	d.Progress.SetLast("go", "beginner", "go-b-loops", "Loops")

	m := newModel(d)
	if m.setupStep != stepResume {
		t.Fatalf("with a saved session, wizard should start at stepResume, got %v", m.setupStep)
	}
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	// Choose "Continue" (first option).
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.phase != phaseReady {
		t.Fatalf("phase = %v, want phaseReady", m.phase)
	}
	if !m.curriculum || m.curr.Lang != "go" || m.curr.Level != "beginner" {
		t.Fatalf("resumed into %s/%s, want go/beginner", m.curr.Lang, m.curr.Level)
	}
	if m.currentTopicID != "go-b-loops" {
		t.Fatalf("resumed topic = %q, want go-b-loops", m.currentTopicID)
	}
}

func TestStartingTopicPersistsResumePoint(t *testing.T) {
	d := testDeps(t)
	m := newModel(d)
	m.loadCurriculum("python", "beginner", "py-b-loops")
	if m.deps.Progress.Last == nil || m.deps.Progress.Last.TopicID != "py-b-loops" {
		t.Fatalf("starting a topic should save it as the resume point, got %+v", m.deps.Progress.Last)
	}
}

func TestCurriculumModeStartsTracksAndSwitches(t *testing.T) {
	d := testDeps(t)
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
	if m.currentTopicID != sel.id {
		t.Fatalf("active topic = %q, want %q", m.currentTopicID, sel.id)
	}
	// Baked content loads synchronously: the editor should hold the new topic's
	// starter code and the tutor should have its tests.
	if m.current.ID != sel.id {
		t.Fatalf("current challenge = %q, want %q", m.current.ID, sel.id)
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
	reply := cmd()
	if _, ok := reply.(chatReplyMsg); !ok {
		t.Fatalf("chat cmd produced %T, want chatReplyMsg", reply)
	}
}
