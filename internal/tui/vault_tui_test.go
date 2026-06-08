package tui

import (
	"os"
	"path/filepath"
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

func TestVaultEscStopsStreamingTutorReply(t *testing.T) {
	m := newTestVaultModel(t)
	m.setFocus(paneChat)
	m.streaming = true
	m.pending = 1
	m.streamCh = make(chan streamChunkMsg, 1)
	cancelled := false
	m.streamCancel = func() { cancelled = true }
	m.chat.beginStream()

	tm, cmd := m.handleKey(tea.KeyMsg{Type: tea.KeyEsc})
	m = tm.(VaultModel)
	if cmd != nil {
		t.Fatal("esc while streaming should not quit or dispatch another command")
	}
	if !cancelled || !m.streamStopping || !m.streaming {
		t.Fatalf("stream not marked stopping: cancelled=%v stopping=%v streaming=%v", cancelled, m.streamStopping, m.streaming)
	}

	tm, _ = m.Update(streamChunkMsg{done: true, full: "ignored"})
	m = tm.(VaultModel)
	if m.streaming || m.streamStopping || m.pending != 0 {
		t.Fatalf("done after stop should clear stream state: streaming=%v stopping=%v pending=%d", m.streaming, m.streamStopping, m.pending)
	}
	if len(m.chatHist) != 0 {
		t.Fatalf("stopped reply should not be recorded as assistant history: %+v", m.chatHist)
	}
}

func TestVaultFileFinderOpensSelectedNote(t *testing.T) {
	m := newTestVaultModel(t)
	vSaveCmd(m.svc, "math/Algebra.md", "# Algebra\n\nSymbols.\n")()
	vSaveCmd(m.svc, "physics/Waves.md", "# Waves\n\nOscillation.\n")()
	tm, _ := m.Update(vListCmd(m.svc)())
	m = tm.(VaultModel)

	tm, _ = m.openFinder("file")
	m = tm.(VaultModel)
	m.finderInput.SetValue("alg")
	m.refreshFinderResults()
	if len(m.finderResults) == 0 {
		t.Fatal("expected file finder matches")
	}
	if got := m.finderResults[0].path; got != "math/Algebra.md" {
		t.Fatalf("top file result = %q, want math/Algebra.md", got)
	}

	tm, cmd := m.updateFinder(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("enter should open the selected note")
	}
	opened, ok := cmd().(vOpenedMsg)
	if !ok {
		t.Fatalf("expected vOpenedMsg, got %T", cmd())
	}
	tm, _ = m.Update(opened)
	m = tm.(VaultModel)
	if m.current != "math/Algebra.md" {
		t.Fatalf("current = %q, want math/Algebra.md", m.current)
	}
	if !m.expanded["math"] {
		t.Fatal("opening from finder should unfold the note's directory")
	}
}

func TestVaultContentFinderShowsSnippetAndOpensNote(t *testing.T) {
	m := newTestVaultModel(t)
	vSaveCmd(m.svc, "math/Linear Algebra.md", "# Linear Algebra\n\nEigenvalues measure a linear map.\n")()
	vSaveCmd(m.svc, "history/Rome.md", "# Rome\n\nRepublic and empire.\n")()
	tm, _ := m.Update(vListCmd(m.svc)())
	m = tm.(VaultModel)

	tm, _ = m.openFinder("grep")
	m = tm.(VaultModel)
	m.finderInput.SetValue("eigen")
	m.refreshFinderResults()
	if len(m.finderResults) == 0 {
		t.Fatal("expected content finder matches")
	}
	first := m.finderResults[0]
	if first.path != "math/Linear Algebra.md" {
		t.Fatalf("top grep result = %q, want math/Linear Algebra.md", first.path)
	}
	if !strings.Contains(first.context, "Eigenvalues") {
		t.Fatalf("grep context should include matching line, got %q", first.context)
	}

	tm, cmd := m.updateFinder(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("enter should open the selected content match")
	}
	opened := cmd().(vOpenedMsg)
	tm, _ = m.Update(opened)
	m = tm.(VaultModel)
	if m.current != "math/Linear Algebra.md" {
		t.Fatalf("current = %q, want math/Linear Algebra.md", m.current)
	}
}

func TestVaultFinderLeaderChords(t *testing.T) {
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tut := tutor.New(config.AIConfig{Provider: "openai"})
	m := newVaultModel(core.New(v, tut), config.Config{Editor: config.EditorConfig{Keybindings: "vim"}})
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = tm.(VaultModel)
	m.setFocus(paneEditor)

	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{','}})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	m = tm.(VaultModel)
	if m.finderMode != "file" {
		t.Fatalf(",ff finder mode = %q, want file", m.finderMode)
	}

	tm, _ = m.updateFinder(tea.KeyMsg{Type: tea.KeyEsc})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{','}})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = tm.(VaultModel)
	if m.finderMode != "grep" {
		t.Fatalf(",fg finder mode = %q, want grep", m.finderMode)
	}
}

func TestVaultPublishCommand(t *testing.T) {
	m := newTestVaultModel(t)
	dest := t.TempDir()
	m.cfg.PublishDir = dest

	// Without a course open, :publish only flashes guidance.
	tm, cmd := m.runEx("publish")
	m = tm.(VaultModel)
	if cmd != nil || !strings.Contains(m.notice, "open a course first") {
		t.Fatalf(":publish without a course: notice = %q", m.notice)
	}

	// Author a minimal course: a topic note plus a manifest linking it.
	vSaveCmd(m.svc, "Algo/BST.md", "# BST\n\nOrdered keys.\n")()
	vSaveCmd(m.svc, core.CourseDir+"/Trees/course.md", "## Basics\n- [[BST]]\n")()
	m.current = core.CourseDir + "/Trees/course.md"

	tm, cmd = m.runEx("publish")
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal(":publish with a course open should run")
	}
	msg, ok := cmd().(vPublishedMsg)
	if !ok {
		t.Fatalf("expected vPublishedMsg, got %#v", cmd())
	}
	if msg.res.Dir != filepath.Join(dest, "Trees") || msg.res.Notes != 2 {
		t.Fatalf("publish result = %+v", msg.res)
	}
	tm, _ = m.Update(msg)
	m = tm.(VaultModel)
	if !strings.Contains(m.chat.view(), "published 2 notes") {
		t.Fatalf("publish confirmation missing from chat:\n%s", m.chat.view())
	}
	for _, f := range []string{"course.md", "BST.md"} {
		if _, err := os.Stat(filepath.Join(dest, "Trees", f)); err != nil {
			t.Errorf("published file missing: %s (%v)", f, err)
		}
	}
}

// On a fresh vault no note is open and focus starts on the sidebar — the
// finder chords must still work (regression: they were editor-pane-only).
func TestVaultFinderOpensFromSidebarOnFreshVault(t *testing.T) {
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tut := tutor.New(config.AIConfig{Provider: "openai"})
	m := newVaultModel(core.New(v, tut), config.Config{Editor: config.EditorConfig{Keybindings: "vim"}})
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = tm.(VaultModel)
	if m.focus != paneSidebar {
		t.Fatalf("fresh vault should focus the sidebar, got %v", m.focus)
	}

	for _, r := range ",ff" {
		tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = tm.(VaultModel)
	}
	if m.finderMode != "file" {
		t.Fatalf(",ff from the sidebar: finder mode = %q, want file", m.finderMode)
	}

	tm, _ = m.updateFinder(tea.KeyMsg{Type: tea.KeyEsc})
	m = tm.(VaultModel)
	for _, r := range ",fg" {
		tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = tm.(VaultModel)
	}
	if m.finderMode != "grep" {
		t.Fatalf(",fg from the sidebar: finder mode = %q, want grep", m.finderMode)
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

	// The vault root is row 0 (a directory); top-level entries nest under it at
	// depth 1, directories first.
	m.rebuildSidebar()
	got := make([]string, 0, len(m.sidebar.items))
	for _, it := range m.sidebar.items {
		got = append(got, it.id)
	}
	want := []string{vaultRootID, "bio", "math", "root.md"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("collapsed tree = %v, want %v", got, want)
	}
	if it := m.sidebar.items[0]; !it.root || !it.dir || !it.expanded || it.depth != 0 {
		t.Fatalf("row 0 should be the open vault root: %+v", it)
	}

	// Expanding math reveals its children (subdirectory before file), at depth 2.
	m.expanded["math"] = true
	m.rebuildSidebar()
	got = got[:0]
	for _, it := range m.sidebar.items {
		got = append(got, it.id)
	}
	want = []string{vaultRootID, "bio", "math", "math/calc", "math/A.md", "root.md"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("expanded tree = %v, want %v", got, want)
	}
	for _, it := range m.sidebar.items {
		if it.id == "math/A.md" && it.depth != 2 {
			t.Fatalf("math/A.md depth = %d, want 2", it.depth)
		}
		if it.id == "math" && (!it.dir || !it.expanded || it.depth != 1) {
			t.Fatalf("math should render as an expanded dir at depth 1: %+v", it)
		}
	}

	// Folding the vault root hides everything beneath it.
	delete(m.expanded, vaultRootID)
	m.rebuildSidebar()
	if len(m.sidebar.items) != 1 || !m.sidebar.items[0].root {
		t.Fatalf("a folded root should leave only the root row, got %d rows", len(m.sidebar.items))
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

	// The cursor starts on the vault root; step down to the "math" dir.
	if it, _ := m.sidebar.selected(); !it.root {
		t.Fatalf("cursor should start on the vault root, got %q", it.id)
	}
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = tm.(VaultModel)
	if it, _ := m.sidebar.selected(); it.id != "math" {
		t.Fatalf("j should land on the math dir, got %q", it.id)
	}
	// Enter on the collapsed "math" dir unfolds it.
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	m = tm.(VaultModel)
	if !m.expanded["math"] {
		t.Fatal("enter should unfold the directory")
	}
	if len(m.sidebar.items) != 4 { // root, math, math/A.md, root.md
		t.Fatalf("after unfold the tree should show 4 rows, got %d", len(m.sidebar.items))
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

// On a course-only vault (the common fresh-install case), the vault root is
// still shown so a new note can be created there — "m a" on the root opens an
// add prompt rooted at the vault, not inside meari-course/.
func TestVaultRootAddCreatesAtRoot(t *testing.T) {
	m := newTestVaultModel(t)
	m.svc.SetCourseDir(t.TempDir())
	if _, err := m.svc.SaveNote("math/A.md", "# A\n"); err != nil { // seed a course
		t.Fatal(err)
	}
	if _, err := m.svc.SaveNote(core.CourseDir+"/Trees/course.md", "## Basics\n- [[A]]\n"); err != nil {
		t.Fatal(err)
	}
	// Simulate a vault whose only user content is the course mount.
	m.tree = []core.TreeEntry{
		{Path: core.CourseDir, Name: core.CourseDir, Dir: true},
		{Path: core.CourseDir + "/Trees", Name: "Trees", Dir: true},
	}
	m.rebuildSidebar()
	m.setFocus(paneSidebar)

	// Row 0 is the vault root and the cursor sits on it.
	if it, ok := m.sidebar.selected(); !ok || !it.root {
		t.Fatalf("cursor should start on the vault root, got %+v", it)
	}

	// "m" then "a" opens the add prompt with an empty (root-level) prefill.
	tm, _ := m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = tm.(VaultModel)
	if !m.cmdMode || m.promptMode != "add" {
		t.Fatalf("m a should open the add prompt, cmdMode=%v mode=%q", m.cmdMode, m.promptMode)
	}
	if v := m.cmdLine.Value(); v != "" {
		t.Fatalf("add prompt on the root should prefill empty (root-level), got %q", v)
	}

	// Submitting a bare name creates the note at the vault root.
	tm, cmd := m.runNodePrompt("add", "", "Ideas")
	m = tm.(VaultModel)
	if cmd == nil {
		t.Fatal("add should issue a save+open command")
	}
	opened, ok := cmd().(vOpenedMsg)
	if !ok {
		t.Fatalf("expected vOpenedMsg, got %T", cmd())
	}
	if opened.note.Path != "Ideas.md" {
		t.Fatalf("new note path = %q, want Ideas.md (vault root)", opened.note.Path)
	}
}

// The vault root is structural — it can't be moved, deleted, or marked.
func TestVaultRootIsProtected(t *testing.T) {
	m := newTestVaultModel(t)
	m.rebuildSidebar()
	m.setFocus(paneSidebar)
	if it, ok := m.sidebar.selected(); !ok || !it.root {
		t.Fatalf("cursor should start on the vault root, got %+v", it)
	}

	// m d on the root refuses with a flash, arms no deletion.
	tm, _ := m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	m = tm.(VaultModel)
	if len(m.confirmDel) != 0 {
		t.Fatalf("deleting the vault root should be refused, confirmDel=%v", m.confirmDel)
	}

	// m m on the root refuses too (no move prompt).
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	m = tm.(VaultModel)
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	m = tm.(VaultModel)
	if m.cmdMode || m.promptMode == "move" {
		t.Fatal("moving the vault root should be refused")
	}

	// Space doesn't mark the root.
	tm, _ = m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = tm.(VaultModel)
	if m.marked[vaultRootID] {
		t.Fatal("the vault root must not be markable")
	}
}
