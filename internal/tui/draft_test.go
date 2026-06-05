package tui

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/drafts"
	"meari/internal/tutor"
	"meari/internal/vault"
)

func TestDraftMatches(t *testing.T) {
	fizz := tutor.Challenge{
		Lang:        "go",
		StarterCode: "func FizzBuzz(n int) []string {\n\treturn nil\n}\n",
	}
	// The user-reported case: a SumTo draft saved when this topic was a
	// different challenge must NOT shadow the FizzBuzz starter.
	if draftMatches(fizz, "func SumTo(n int) int {\n    return 0\n}") {
		t.Fatal("an unrelated draft must be detected as stale")
	}
	// A genuine in-progress FizzBuzz draft is kept.
	if !draftMatches(fizz, "func FizzBuzz(n int) []string {\n\tout := []string{}\n\treturn out\n}") {
		t.Fatal("a matching draft must be kept")
	}

	py := tutor.Challenge{StarterCode: "def sum_to(n):\n    pass\n"}
	if draftMatches(py, "def evens(nums):\n    pass\n") {
		t.Fatal("stale python draft must be detected")
	}
	if !draftMatches(py, "def sum_to(n):\n    return n\n") {
		t.Fatal("matching python draft must be kept")
	}

	// Prose challenges (physics) extract no names: any draft is acceptable.
	prose := tutor.Challenge{Lang: "physics", StarterCode: "Write your reasoning here:\n"}
	if !draftMatches(prose, "My reasoning about velocity is...") {
		t.Fatal("prose drafts must always be accepted")
	}

	// Go type-based starters match on the type name too.
	typed := tutor.Challenge{Lang: "go", StarterCode: "type Celsius int\n\nfunc (c Celsius) String() string { return \"\" }\n"}
	if !draftMatches(typed, "type Celsius int\nfunc (c Celsius) String() string { return \"x\" }") {
		t.Fatal("type-name match must be kept")
	}
}

func TestVaultClickFocusesPane(t *testing.T) {
	m := newTestVaultModel(t) // sized 100x40 by the helper
	if m.focus != paneSidebar {
		t.Fatalf("initial focus = %v", m.focus)
	}
	click := func(x, y int) {
		tm, _ := m.Update(tea.MouseMsg{
			X: x, Y: y,
			Action: tea.MouseActionPress,
			Button: tea.MouseButtonLeft,
		})
		m = tm.(VaultModel)
	}
	// Click in the middle of the editor column.
	click(m.sidebarW+2+m.editorW/2, 10)
	if m.focus != paneEditor {
		t.Fatalf("click in editor: focus = %v", m.focus)
	}
	// Click far right -> chat.
	click(m.sidebarW+2+m.editorW+2+m.chatW/2, 10)
	if m.focus != paneChat {
		t.Fatalf("click in chat: focus = %v", m.focus)
	}
	// Click far left -> sidebar.
	click(1, 10)
	if m.focus != paneSidebar {
		t.Fatalf("click in sidebar: focus = %v", m.focus)
	}
	// Clicking the title bar (row 0) changes nothing.
	click(1, 0)
	if m.focus != paneSidebar {
		t.Fatalf("click on title bar should not change focus, got %v", m.focus)
	}
}

func TestClassicClickFocusesPane(t *testing.T) {
	m := newModel(testDeps(t))
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady

	click := func(x, y int) {
		m = step(t, m, tea.MouseMsg{
			X: x, Y: y,
			Action: tea.MouseActionPress,
			Button: tea.MouseButtonLeft,
		})
	}
	click(m.sidebarW+2+m.editorW/2, 10) // editor column
	if m.focus != paneEditor {
		t.Fatalf("click in editor: focus = %v", m.focus)
	}
	click(m.sidebarW+2+m.editorW+2+m.chatW/2, 10) // chat column
	if m.focus != paneChat {
		t.Fatalf("click in chat: focus = %v", m.focus)
	}
	click(1, 10) // sidebar
	if m.focus != paneSidebar {
		t.Fatalf("click in sidebar: focus = %v", m.focus)
	}
	// Wheel under the editor must NOT steal focus (ranger-style).
	m = step(t, m, tea.MouseMsg{
		X: m.sidebarW + 4, Y: 10,
		Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown,
	})
	if m.focus != paneSidebar {
		t.Fatalf("wheel should not change focus, got %v", m.focus)
	}
}

func TestConfigurablePaneRatios(t *testing.T) {
	m := newModel(testDeps(t))
	m.deps.Cfg.UI.ChatPercent = 50
	m.deps.Cfg.UI.SidebarPercent = 12
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady

	contentW := 120 - 6 // three bordered boxes
	if want := clampMin(contentW*50/100, 16); m.chatW != want {
		t.Fatalf("chatW = %d, want %d (50%%)", m.chatW, want)
	}
	if want := clampMin(contentW*12/100, 12); m.sidebarW != want {
		t.Fatalf("sidebarW = %d, want %d (12%%)", m.sidebarW, want)
	}

	// Unset config keeps the built-in defaults.
	m2 := newModel(testDeps(t))
	m2 = step(t, m2, tea.WindowSizeMsg{Width: 120, Height: 40})
	if want := clampMin(contentW*30/100, 16); m2.chatW != want {
		t.Fatalf("default chatW = %d, want %d", m2.chatW, want)
	}
}

func TestCmdLineHistory(t *testing.T) {
	m := newTestVaultModel(t)

	type tm = tea.KeyMsg
	typeStr := func(s string) {
		for _, r := range s {
			mm, _ := m.Update(tm{Type: tea.KeyRunes, Runes: []rune{r}})
			m = mm.(VaultModel)
		}
	}
	press := func(k tea.KeyType) {
		mm, _ := m.Update(tm{Type: k})
		m = mm.(VaultModel)
	}

	// Run two commands through the global ":" line.
	typeStr(":")
	typeStr("copy")
	press(tea.KeyEnter)
	typeStr(":")
	typeStr("wide")
	press(tea.KeyEnter)

	// Reopen and recall: ↑ gives the latest, ↑ again the older, ↓ back down.
	typeStr(":")
	press(tea.KeyUp)
	if got := m.cmdLine.Value(); got != "wide" {
		t.Fatalf("↑ = %q, want \"wide\"", got)
	}
	press(tea.KeyUp)
	if got := m.cmdLine.Value(); got != "copy" {
		t.Fatalf("↑↑ = %q, want \"copy\"", got)
	}
	press(tea.KeyDown)
	if got := m.cmdLine.Value(); got != "wide" {
		t.Fatalf("↓ = %q, want \"wide\"", got)
	}
	press(tea.KeyDown)
	if got := m.cmdLine.Value(); got != "" {
		t.Fatalf("↓ past newest should restore the empty draft, got %q", got)
	}
}

func TestConfigReloadRebuildsTutor(t *testing.T) {
	dir := t.TempDir()
	cfgPath := dir + "/config.toml"
	mustWrite := func(s string) {
		if err := os.WriteFile(cfgPath, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Start pointed at OpenAI with no key: offline.
	mustWrite("[ai]\nprovider = \"openai\"\nmodel = \"m\"\n")
	cfg, err := config.Load(cfgPath, dir)
	if err != nil {
		t.Fatal(err)
	}
	m := Model{
		deps: Deps{Tutor: tutor.New(cfg.AI), Cfg: cfg, ConfigPath: cfgPath, BaseDir: dir},
		chat: newChat(),
	}
	m.chat.setSize(60, 12)
	if !m.deps.Tutor.Offline() {
		t.Fatal("setup: tutor should start offline")
	}

	// Edit the config to a local compatible provider and reload: the tutor must
	// be rebuilt and come online without restarting the app.
	mustWrite("[ai]\nprovider = \"compatible\"\nbase_url = \"http://localhost:9999/v1\"\nmodel = \"m\"\n")
	m.applyConfigReload(configReloadMsg{})
	if m.deps.Tutor.Offline() {
		t.Fatal("reload must rebuild the tutor from the new [ai] section")
	}
	if got := m.deps.Tutor.Info().BaseURL; got != "http://localhost:9999/v1" {
		t.Fatalf("tutor base URL after reload = %q", got)
	}
}

func TestLoadStarterOrDraftStale(t *testing.T) {
	store, err := drafts.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	m := Model{deps: Deps{Store: store}}

	ch := tutor.Challenge{
		ID:          "go-b-loops",
		Lang:        "go",
		StarterCode: "func FizzBuzz(n int) []string {\n\treturn nil\n}\n",
	}

	// No draft: the starter, topped with the prompt as a comment.
	code, stale := m.loadStarterOrDraft(ch)
	if stale || code != promptComment(ch)+ch.StarterCode {
		t.Fatalf("no-draft case wrong: stale=%v code=%q", stale, code)
	}

	// Stale draft (the reported bug): starter wins and it's flagged.
	if err := store.Save(ch.ID, "func SumTo(n int) int {\n    return 0\n}"); err != nil {
		t.Fatal(err)
	}
	code, stale = m.loadStarterOrDraft(ch)
	if !stale {
		t.Fatal("a mismatched draft must be flagged stale")
	}
	if code != promptComment(ch)+ch.StarterCode {
		t.Fatalf("stale draft must not shadow the starter, got %q", code)
	}

	// Matching draft: kept, not stale.
	want := "func FizzBuzz(n int) []string {\n\tvar out []string\n\treturn out\n}"
	if err := store.Save(ch.ID, want); err != nil {
		t.Fatal(err)
	}
	code, stale = m.loadStarterOrDraft(ch)
	if stale || code != want {
		t.Fatalf("matching draft case wrong: stale=%v code=%q", stale, code)
	}
}

func TestSidebarFoldedFromConfig(t *testing.T) {
	// Classic TUI starts folded when configured.
	d := testDeps(t)
	d.Cfg.UI.SidebarFolded = true
	m := newModel(d)
	m = step(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.phase = phaseReady
	if !m.sidebarCollapsed {
		t.Fatal("classic TUI should start folded")
	}
	if m.sidebarW != 0 {
		t.Fatalf("folded sidebar width = %d, want 0", m.sidebarW)
	}

	// Vault TUI: starts folded, focus on the editor, :fold brings it back.
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tut := tutor.New(config.AIConfig{Provider: "openai"})
	vm := newVaultModel(core.New(v, tut), config.Config{UI: config.UIConfig{SidebarFolded: true}})
	tm, _ := vm.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	vm = tm.(VaultModel)
	if !vm.sidebarCollapsed || vm.sidebarW != 0 {
		t.Fatalf("vault should start folded (collapsed=%v w=%d)", vm.sidebarCollapsed, vm.sidebarW)
	}
	if vm.focus != paneEditor {
		t.Fatalf("folded start should focus the editor, got %v", vm.focus)
	}
	tm, _ = vm.runEx("fold")
	vm = tm.(VaultModel)
	if vm.sidebarCollapsed || vm.sidebarW == 0 {
		t.Fatal(":fold should unfold the vault sidebar")
	}
}

func TestPromptCommentLanguages(t *testing.T) {
	goCh := tutor.Challenge{Lang: "go", Prompt: "Write SumTo(n int) int returning the sum."}
	if got := promptComment(goCh); !strings.HasPrefix(got, "// Write SumTo") {
		t.Fatalf("go comment = %q", got)
	}
	pyCh := tutor.Challenge{Prompt: "Write sum_to(n) returning the sum."}
	if got := promptComment(pyCh); !strings.HasPrefix(got, "# Write sum_to") {
		t.Fatalf("python comment = %q", got)
	}
	prose := tutor.Challenge{Lang: "physics", Prompt: "Explain the energy transfer."}
	if got := promptComment(prose); !strings.HasPrefix(got, "Explain the energy") {
		t.Fatalf("prose header = %q", got)
	}
	// Long prompts wrap into multiple comment lines.
	long := tutor.Challenge{Lang: "go", Prompt: strings.Repeat("word ", 40)}
	if lines := strings.Count(promptComment(long), "// "); lines < 2 {
		t.Fatalf("long prompt should wrap, got %d comment lines", lines)
	}
}

func TestEssayHeaderSeededAndStripped(t *testing.T) {
	m := newTestVaultModel(t)
	opened := vSaveOpenCmd(m.svc, "x/N.md", "# N\n\nbody\n")().(vOpenedMsg)
	tm, _ := m.Update(opened)
	m = tm.(VaultModel)
	tm, _ = m.startEssay("Explain N.")
	m = tm.(VaultModel)

	if !strings.HasPrefix(m.editor.Value(), "> Essay: Explain N.") {
		t.Fatalf("answer buffer should open with the prompt header: %q", m.editor.Value())
	}
	// Header alone is not an answer.
	tm, cmd := m.gradeEssay()
	m = tm.(VaultModel)
	if cmd != nil {
		t.Fatal("grading with only the header should be refused")
	}
	// The header is stripped from what gets graded.
	if got := stripEssayHeader("> Essay: Explain N.\n\nmy actual answer"); strings.TrimSpace(got) != "my actual answer" {
		t.Fatalf("stripEssayHeader = %q", got)
	}
}
