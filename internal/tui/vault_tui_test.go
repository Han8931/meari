package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/tutor"
	"meari/internal/vault"
)

func newTestVaultModel(t *testing.T) VaultModel {
	t.Helper()
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tut := tutor.New(config.AIConfig{Provider: "openai"}) // offline (no key)
	if !tut.Offline() {
		t.Fatal("expected offline tutor")
	}
	m := newVaultModel(core.New(v, tut), config.Config{})
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	return tm.(VaultModel)
}

func TestVaultGenerateAndOpen(t *testing.T) {
	m := newTestVaultModel(t)

	// :learn generates a lesson note (offline placeholder).
	tm, cmd := m.runEx("learn the cold war")
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("expected a generate command")
	}
	gen, ok := cmd().(vGeneratedMsg)
	if !ok {
		t.Fatalf("expected vGeneratedMsg, got %T", cmd())
	}
	if gen.meta.Path == "" {
		t.Fatal("generated note has no path")
	}
	tm, _ = m.Update(gen)
	m = tm.(VaultModel)
	if m.pending != 0 {
		t.Fatalf("pending should be back to 0, got %d", m.pending)
	}

	// Open the generated note and confirm it loads into the editor.
	opened := vOpenCmd(m.svc, gen.meta.Path)().(vOpenedMsg)
	tm, _ = m.Update(opened)
	m = tm.(VaultModel)
	if m.current != gen.meta.Path {
		t.Fatalf("current = %q, want %q", m.current, gen.meta.Path)
	}
	if strings.TrimSpace(m.editor.Value()) == "" {
		t.Fatal("editor should hold the note body")
	}
	if m.focus != paneEditor {
		t.Fatal("opening a note should focus the editor")
	}
}

func TestVaultNewNote(t *testing.T) {
	m := newTestVaultModel(t)
	_, cmd := m.runEx("new My First Note")
	if cmd == nil {
		t.Fatal("expected a command")
	}
	opened, ok := cmd().(vOpenedMsg)
	if !ok {
		t.Fatalf("expected vOpenedMsg, got %T", cmd())
	}
	if opened.note.Title != "My First Note" {
		t.Fatalf("title = %q", opened.note.Title)
	}
}

func TestVaultEssayStudyFlow(t *testing.T) {
	m := newTestVaultModel(t)

	// Create and open a note to study.
	opened := vSaveOpenCmd(m.svc, "history/Cold War.md", "# Cold War\n\nA rivalry.\n")().(vOpenedMsg)
	tm, _ := m.Update(opened)
	m = tm.(VaultModel)

	// Start an essay study; the editor is cleared for the answer and autosave is
	// suspended (curPath blanked) so the answer can't overwrite the note.
	tm, _ = m.startEssay("")
	m = tm.(VaultModel)
	if !m.studyMode {
		t.Fatal("should be in study mode")
	}
	if *m.curPath != "" {
		t.Fatal("note autosave should be suspended during study")
	}

	// Write an answer and grade it.
	m.editor.SetValue("My understanding of the Cold War is...")
	tm, cmd := m.gradeEssay()
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("expected a grade command")
	}
	res, ok := cmd().(vEssayMsg)
	if !ok {
		t.Fatalf("expected vEssayMsg, got %T", cmd())
	}
	if res.res.Score != 1 { // offline: any non-empty answer passes
		t.Fatalf("offline essay score = %v", res.res.Score)
	}
	tm, _ = m.Update(res)
	m = tm.(VaultModel)

	// End study returns to the note.
	tm, cmd = m.endEssay()
	m = tm.(VaultModel)
	if m.studyMode {
		t.Fatal("study mode should be off after :done")
	}
	reopened := cmd().(vOpenedMsg)
	if reopened.note.Title != "Cold War" {
		t.Fatalf("reopened note = %q", reopened.note.Title)
	}
}

func TestVaultGradeWithoutStudyIsSafe(t *testing.T) {
	m := newTestVaultModel(t)
	// :grade outside study mode must not panic or fire a command.
	_, cmd := m.runEx("grade")
	if cmd != nil {
		t.Fatal("grade outside study should be a no-op")
	}
}

func TestVaultAnswerCommand(t *testing.T) {
	m := newTestVaultModel(t)

	// Outside study mode :answer is a guarded no-op.
	_, cmd := m.runEx("answer")
	if cmd != nil {
		t.Fatal(":answer outside study should be a no-op")
	}

	// In study mode it fetches a model answer.
	opened := vSaveOpenCmd(m.svc, "x/N.md", "# N\n\nbody\n")().(vOpenedMsg)
	tm, _ := m.Update(opened)
	m = tm.(VaultModel)
	tm, _ = m.startEssay("")
	m = tm.(VaultModel)

	tm, cmd = m.runEx("answer")
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("expected an answer command in study mode")
	}
	msg, ok := cmd().(vAnswerMsg)
	if !ok {
		t.Fatalf("expected vAnswerMsg, got %T", cmd())
	}
	if msg.text == "" {
		t.Fatal("model answer should not be empty (offline fallback)")
	}
	tm, _ = m.Update(msg)
	m = tm.(VaultModel)
	if m.pending != 0 {
		t.Fatalf("pending should return to 0, got %d", m.pending)
	}
}

func TestVaultViewRenders(t *testing.T) {
	m := newTestVaultModel(t)
	out := m.View()
	if !strings.Contains(out, "Meari") {
		t.Fatalf("view should render the title bar, got:\n%s", out)
	}
	// And after opening a note (study header path exercised).
	opened := vSaveOpenCmd(m.svc, "x/N.md", "# N\n\nbody\n")().(vOpenedMsg)
	tm, _ := m.Update(opened)
	m = tm.(VaultModel)
	if s := m.View(); !strings.Contains(s, "NOTE") {
		t.Fatalf("view should show the note header, got:\n%s", s)
	}
}

func TestVaultRebuildSidebarGroupsBySubject(t *testing.T) {
	m := newTestVaultModel(t)
	m.notes = []core.NoteMeta{
		{Path: "math/A.md", Title: "A", Subject: "math"},
		{Path: "bio/B.md", Title: "B", Subject: "bio"},
		{Path: "math/C.md", Title: "C", Subject: "math"},
	}
	m.rebuildSidebar()
	// Expect: header(bio), B, header(math), A, C  -> 2 headers, 3 notes.
	headers, notes := 0, 0
	for _, it := range m.sidebar.items {
		if it.header {
			headers++
		} else {
			notes++
		}
	}
	if headers != 2 || notes != 3 {
		t.Fatalf("grouping wrong: %d headers, %d notes", headers, notes)
	}
}
