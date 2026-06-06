package tui

import (
	"strings"
	"testing"

	"meari/internal/config"
	"meari/internal/drafts"
	"meari/internal/progress"
	"meari/internal/tutor"
)

// :vault in the coding tutor sets the exit target so the shell loop opens the vault.
func TestTutorVaultSwitchCommand(t *testing.T) {
	store, err := drafts.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	prog, err := progress.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := Deps{
		Tutor:    tutor.New(config.AIConfig{Provider: "openai"}), // offline
		Store:    store,
		Progress: prog,
		Cfg:      config.Config{},
		Topic:    "spanish basics", // skip the wizard -> phaseReady
	}
	m := newModel(d)

	tm, _ := m.runEx("vault")
	if got := tm.(Model).exit; got != SwitchToVault {
		t.Fatalf("exit = %v, want SwitchToVault", got)
	}
}

// :tutor in the vault sets the exit target so the shell loop opens the coding TUI.
func TestVaultTutorSwitchCommand(t *testing.T) {
	m := newTestVaultModel(t)
	tm, _ := m.runEx("tutor")
	if got := tm.(VaultModel).exit; got != SwitchToTutor {
		t.Fatalf("exit = %v, want SwitchToTutor", got)
	}
}

// Opening a note that another note links to populates the backlinks panel.
func TestVaultBacklinksPanel(t *testing.T) {
	m := newTestVaultModel(t)

	// A is the target; B links to A via [[A]].
	_ = vSaveOpenCmd(m.svc, "A.md", "# A\n\nalpha\n")()
	_ = vSaveOpenCmd(m.svc, "B.md", "# B\n\nsee [[A]] for context\n")()

	links := vBacklinksCmd(m.svc, "A.md")().(vBacklinksMsg)
	if links.path != "A.md" {
		t.Fatalf("path = %q, want A.md", links.path)
	}
	if len(links.links) != 1 || links.links[0].Path != "B.md" {
		t.Fatalf("backlinks = %+v, want [B.md]", links.links)
	}

	// Open A, then feed the backlinks message; the panel should render.
	opened := vOpenCmd(m.svc, "A.md")().(vOpenedMsg)
	tm, _ := m.Update(opened)
	m = tm.(VaultModel)
	tm, _ = m.Update(links)
	m = tm.(VaultModel)
	if len(m.backlinks) != 1 {
		t.Fatalf("model backlinks = %d, want 1", len(m.backlinks))
	}
	if !strings.Contains(m.View(), "Linked mentions") {
		t.Fatalf("view should show the backlinks panel, got:\n%s", m.View())
	}

	// :backlinks toggles the panel off.
	tm, _ = m.runEx("backlinks")
	m = tm.(VaultModel)
	if m.showBacklinks {
		t.Fatal(":backlinks should toggle the panel off")
	}
	if strings.Contains(m.View(), "Linked mentions") {
		t.Fatal("panel should be hidden after toggling off")
	}
}
