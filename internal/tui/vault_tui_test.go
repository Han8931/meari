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

func TestVaultSidebarTree(t *testing.T) {
	m := newTestVaultModel(t)
	m.tree = []core.TreeEntry{
		{Path: "bio", Name: "bio", Dir: true},
		{Path: "bio/B.md", Name: "B"},
		{Path: "math", Name: "math", Dir: true},
		{Path: "math/A.md", Name: "A"},
		{Path: "math/calc", Name: "calc", Dir: true},
		{Path: "math/calc/C.md", Name: "C"},
		{Path: "root.md", Name: "root"},
	}

	// Collapsed by default: only the top level is visible, directories first.
	m.rebuildSidebar()
	got := make([]string, 0, len(m.sidebar.items))
	for _, it := range m.sidebar.items {
		got = append(got, it.id)
	}
	want := []string{"bio", "math", "root.md"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("collapsed tree = %v, want %v", got, want)
	}

	// Expanding math reveals its children (subdirectory before file), at depth 1.
	m.expanded["math"] = true
	m.rebuildSidebar()
	got = got[:0]
	for _, it := range m.sidebar.items {
		got = append(got, it.id)
	}
	want = []string{"bio", "math", "math/calc", "math/A.md", "root.md"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("expanded tree = %v, want %v", got, want)
	}
	for _, it := range m.sidebar.items {
		if it.id == "math/A.md" && it.depth != 1 {
			t.Fatalf("math/A.md depth = %d, want 1", it.depth)
		}
		if it.id == "math" && (!it.dir || !it.expanded) {
			t.Fatalf("math should render as an expanded dir: %+v", it)
		}
	}
}

func TestVaultSidebarMarkAndDirToggle(t *testing.T) {
	m := newTestVaultModel(t)
	// A real vault: one dir with a note, one root note.
	vSaveCmd(m.svc, "math/A.md", "# A\n")()
	vSaveCmd(m.svc, "root.md", "# Root\n")()
	tm, _ := m.Update(vListCmd(m.svc)())
	m = tm.(VaultModel)
	m.setFocus(paneSidebar)

	// Enter on the collapsed "math" dir unfolds it.
	if it, _ := m.sidebar.selected(); it.id != "math" {
		t.Fatalf("cursor should start on the math dir, got %q", it.id)
	}
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(VaultModel)
	if !m.expanded["math"] {
		t.Fatal("enter should unfold the directory")
	}
	if len(m.sidebar.items) != 3 {
		t.Fatalf("after unfold the tree should show 3 rows, got %d", len(m.sidebar.items))
	}

	// Space marks the row under the cursor and steps down.
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = tm.(VaultModel)
	if !m.marked["math"] {
		t.Fatal("space should mark the dir under the cursor")
	}
	if it, _ := m.sidebar.selected(); it.id == "math" {
		t.Fatal("space should step the cursor down after marking")
	}

	// m then d arms deletion of the marked rows; "n" cancels.
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = tm.(VaultModel)
	if len(m.confirmDel) != 1 || m.confirmDel[0] != "math" {
		t.Fatalf("confirmDel = %v, want [math]", m.confirmDel)
	}
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = tm.(VaultModel)
	if len(m.confirmDel) != 0 {
		t.Fatal("any non-y key should cancel the delete")
	}

	// m then d then y deletes the marked dir and its note.
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = tm.(VaultModel)
	tm, cmd := m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("y should issue the delete command")
	}
	del, ok := cmd().(vDeletedMsg)
	if !ok {
		t.Fatalf("expected vDeletedMsg, got %T", cmd())
	}
	tm, cmd = m.Update(del)
	m = tm.(VaultModel)
	tm, _ = m.Update(cmd().(vNotesMsg)) // refresh the tree
	m = tm.(VaultModel)
	for _, it := range m.sidebar.items {
		if it.id == "math" || it.id == "math/A.md" {
			t.Fatalf("deleted %q still in the tree", it.id)
		}
	}
}
