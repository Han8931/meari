package tui

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/config"
	"meari/internal/drafts"
	"meari/internal/tutor"
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

	// No draft: starter.
	code, stale := m.loadStarterOrDraft(ch)
	if stale || code != ch.StarterCode {
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
	if code != ch.StarterCode {
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
