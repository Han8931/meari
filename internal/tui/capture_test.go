package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/tutor"
)

func TestLastExchange(t *testing.T) {
	hist := []tutor.ChatTurn{
		{Role: "user", Content: "first?"},
		{Role: "assistant", Content: "first answer"},
		{Role: "user", Content: "second?"},
		{Role: "assistant", Content: "second answer"},
	}
	ex, ok := lastExchange(hist)
	if !ok {
		t.Fatal("expected an exchange")
	}
	if ex.Q != "second?" || ex.A != "second answer" {
		t.Fatalf("got %+v, want the LAST pair", ex)
	}

	// A question still awaiting its answer is not capturable.
	if _, ok := lastExchange([]tutor.ChatTurn{{Role: "user", Content: "pending?"}}); ok {
		t.Fatal("an unanswered question should not be capturable")
	}
	if _, ok := lastExchange(nil); ok {
		t.Fatal("empty history should not be capturable")
	}
}

func TestAllExchanges(t *testing.T) {
	hist := []tutor.ChatTurn{
		{Role: "user", Content: "q1"},
		{Role: "assistant", Content: "a1"},
		{Role: "user", Content: "q2"},
		{Role: "assistant", Content: "a2"},
		{Role: "user", Content: "unanswered"},
	}
	got := allExchanges(hist)
	if len(got) != 2 {
		t.Fatalf("got %d exchanges, want 2 (the unanswered one is skipped)", len(got))
	}
	if got[0].Q != "q1" || got[1].A != "a2" {
		t.Fatalf("exchanges out of order: %+v", got)
	}
}

// Capturing twice must reuse the section rather than repeat the heading.
func TestAppendCaptureIsIdempotentOnHeading(t *testing.T) {
	body := newCompanionBody("Ownership & Moves")
	if !strings.Contains(body, "[[Ownership & Moves]]") {
		t.Fatalf("companion note should link back to the lecture:\n%s", body)
	}

	body = appendCapture(body, formatExchanges([]qaExchange{{Q: "why move?", A: "ownership transfers"}}))
	body = appendCapture(body, formatExchanges([]qaExchange{{Q: "and clone?", A: "it copies"}}))

	if n := strings.Count(body, notesSection); n != 1 {
		t.Fatalf("heading appears %d times, want 1:\n%s", n, body)
	}
	for _, want := range []string{"why move?", "ownership transfers", "and clone?", "it copies"} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in:\n%s", want, body)
		}
	}
	// Order preserved: the first capture comes before the second.
	if strings.Index(body, "why move?") > strings.Index(body, "and clone?") {
		t.Fatal("captures should append in order")
	}
}

// A note with no section yet (e.g. a hand-written companion) gains one.
func TestAppendCaptureAddsMissingSection(t *testing.T) {
	got := appendCapture("Some existing prose.", "**Q:** x\n\ny")
	if !strings.Contains(got, notesSection) {
		t.Fatalf("section not added:\n%s", got)
	}
	if !strings.HasPrefix(got, "Some existing prose.") {
		t.Fatalf("existing content should be preserved:\n%s", got)
	}
	// An empty block leaves the note untouched.
	if out := appendCapture("body", "   "); out != "body" {
		t.Fatalf("empty capture changed the note: %q", out)
	}
}

func TestCompanionPath(t *testing.T) {
	if got := companionPath("Ownership & Moves"); got != "My Notes/Ownership & Moves.md" {
		t.Fatalf("companionPath = %q", got)
	}
	// Illegal filename characters are stripped, not passed through.
	if got := companionPath("A/B:C"); strings.ContainsAny(got[len("My Notes/"):], `/\:*?"<>|`) {
		t.Fatalf("unsafe characters survived: %q", got)
	}
}

// End-to-end: :capture writes a companion note that links back to the lecture,
// leaves the lecture untouched, and appends on a second capture.
func TestVaultCaptureWritesCompanionNote(t *testing.T) {
	m := newTestVaultModel(t)
	const lecture = "Ownership.md"
	m = openNote(t, m, lecture, "# Ownership\n\nThe original lecture.\n")

	// Nothing asked yet.
	tm, _ := m.runEx("capture")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "nothing to capture") {
		t.Fatalf("notice = %q", m.notice)
	}

	m.chatHist = []tutor.ChatTurn{
		{Role: "user", Content: "why does a move happen?"},
		{Role: "assistant", Content: "Because ownership transfers."},
	}
	tm, _ = m.runEx("capture")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "captured") {
		t.Fatalf("notice = %q", m.notice)
	}

	companion, err := m.svc.OpenNote(companionPath("Ownership"))
	if err != nil {
		t.Fatalf("companion note not created: %v", err)
	}
	for _, want := range []string{"[[Ownership]]", "why does a move happen?", "Because ownership transfers."} {
		if !strings.Contains(companion.Body, want) {
			t.Fatalf("companion missing %q:\n%s", want, companion.Body)
		}
	}

	// The lecture itself must be untouched — that's the whole point.
	lec, err := m.svc.OpenNote(lecture)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(lec.Body, "why does a move happen?") {
		t.Fatalf("capture leaked into the lecture:\n%s", lec.Body)
	}

	// A second capture appends under the same heading.
	m.chatHist = append(m.chatHist,
		tutor.ChatTurn{Role: "user", Content: "what about clone?"},
		tutor.ChatTurn{Role: "assistant", Content: "It copies instead."})
	tm, _ = m.runEx("capture")
	m = tm.(VaultModel)

	companion, err = m.svc.OpenNote(companionPath("Ownership"))
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(companion.Body, notesSection); n != 1 {
		t.Fatalf("heading duplicated (%d):\n%s", n, companion.Body)
	}
	if !strings.Contains(companion.Body, "It copies instead.") {
		t.Fatalf("second capture missing:\n%s", companion.Body)
	}
}

// ":capture all" saves the whole conversation; a bad argument is rejected.
func TestVaultCaptureAllAndUsage(t *testing.T) {
	m := newTestVaultModel(t)
	m = openNote(t, m, "Traits.md", "# Traits\n")
	m.chatHist = []tutor.ChatTurn{
		{Role: "user", Content: "q1"},
		{Role: "assistant", Content: "a1"},
		{Role: "user", Content: "q2"},
		{Role: "assistant", Content: "a2"},
	}

	tm, _ := m.runEx("capture bogus")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "usage:") {
		t.Fatalf("bad argument should show usage, got %q", m.notice)
	}

	tm, _ = m.runEx("capture all")
	m = tm.(VaultModel)
	companion, err := m.svc.OpenNote(companionPath("Traits"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"q1", "a1", "q2", "a2"} {
		if !strings.Contains(companion.Body, want) {
			t.Fatalf("capture all missing %q:\n%s", want, companion.Body)
		}
	}
}

// Typing ":capture" into the command line (not calling runEx directly) must
// reach the command — this covers the key-handler wiring the other tests skip.
func TestVaultCaptureViaCommandLine(t *testing.T) {
	m := newTestVaultModel(t)
	m = openNote(t, m, "Borrowing.md", "# Borrowing\n")
	m.chatHist = []tutor.ChatTurn{
		{Role: "user", Content: "why &mut?"},
		{Role: "assistant", Content: "Exclusive access."},
	}
	m.setFocus(paneSidebar) // ":" opens the global command line from the notes pane

	send := func(msg tea.Msg) {
		tm, _ := m.Update(msg)
		m = tm.(VaultModel)
	}
	send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	for _, r := range "capture" {
		send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	send(tea.KeyMsg{Type: tea.KeyEnter})

	companion, err := m.svc.OpenNote(companionPath("Borrowing"))
	if err != nil {
		t.Fatalf("typing :capture did not create the companion note: %v", err)
	}
	if !strings.Contains(companion.Body, "Exclusive access.") {
		t.Fatalf("companion missing the answer:\n%s", companion.Body)
	}
}

// :weave needs an AI provider and something already captured.
func TestVaultCaptureWeaveGuards(t *testing.T) {
	m := newTestVaultModel(t) // offline tutor
	m = openNote(t, m, "Lifetimes.md", "# Lifetimes\n")

	tm, _ := m.runEx("weave")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "needs an AI provider") {
		t.Fatalf("offline notice = %q", m.notice)
	}
	if m.pendingEdit != "" || m.streaming {
		t.Fatal("offline weave must not start a stream")
	}
}

// A completed weave stream parks the rewrite as a DETACHED proposal for the
// companion note; :apply writes that note without it being open, and the
// lecture on screen is untouched.
func TestVaultCaptureWeaveApplyWritesCompanion(t *testing.T) {
	m := newTestVaultModel(t)
	const lecture = "Traits.md"
	m = openNote(t, m, lecture, "# Traits\n\nThe lecture body.\n")

	// Seed a companion note as :capture would have.
	companion := companionPath("Traits")
	if _, err := m.svc.SaveNote(companion, newCompanionBody("Traits")+"\n**Q:** what is dyn?\n\nDynamic dispatch.\n"); err != nil {
		t.Fatal(err)
	}

	// Simulate the weave stream in flight, then completing (no network).
	m.streaming, m.polishing, m.pending = true, true, 1
	m.pendingEditPath = companion
	m.pendingDetached = true
	m.chat.beginStream()
	woven := "Notes captured while studying [[Traits]].\n\n## Dynamic dispatch\n\n`dyn` selects the method at runtime.\n"
	tm, _ := m.handleStreamChunk(streamChunkMsg{done: true, full: woven})
	m = tm.(VaultModel)

	if m.pendingEdit != woven {
		t.Fatalf("pendingEdit = %q", m.pendingEdit)
	}
	// Not written until :apply — and the open lecture is unaffected.
	got, _ := m.svc.OpenNote(companion)
	if strings.Contains(got.Body, "selects the method at runtime") {
		t.Fatal("weave wrote the note before :apply")
	}
	if m.current != lecture {
		t.Fatalf("weave should not switch away from the lecture, got %q", m.current)
	}

	tm, _ = m.runEx("apply")
	m = tm.(VaultModel)

	got, err := m.svc.OpenNote(companion)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.Body, "selects the method at runtime") {
		t.Fatalf("apply did not write the companion note:\n%s", got.Body)
	}
	if strings.Contains(got.Body, "**Q:** what is dyn?") {
		t.Fatalf("weave should have reorganized the raw Q&A away:\n%s", got.Body)
	}
	// The lecture itself is still untouched, and the editor still shows it.
	lec, _ := m.svc.OpenNote(lecture)
	if !strings.Contains(lec.Body, "The lecture body.") || strings.Contains(lec.Body, "Dynamic dispatch") {
		t.Fatalf("the lecture was modified:\n%s", lec.Body)
	}
	if m.pendingEdit != "" || m.pendingDetached {
		t.Fatal("apply should clear the proposal state")
	}
}

// A detached proposal survives opening another note (unlike a :polish edit,
// which is discarded because it targets the open editor).
func TestDetachedProposalSurvivesNoteSwitch(t *testing.T) {
	m := newTestVaultModel(t)
	m = openNote(t, m, "A.md", "# A\n")
	m.pendingEdit = "woven"
	m.pendingEditPath = companionPath("A")
	m.pendingDetached = true

	m = openNote(t, m, "B.md", "# B\n")
	if m.pendingEdit == "" {
		t.Fatal("a detached proposal should survive switching notes")
	}

	// A non-detached (:polish) proposal still gets dropped on a switch.
	m.pendingDetached = false
	m.pendingEditPath = "A.md"
	m.pendingEdit = "polished"
	m = openNote(t, m, "C.md", "# C\n")
	if m.pendingEdit != "" {
		t.Fatal("a :polish proposal for another note should be discarded")
	}
}

// :weave is a top-level command (not a :capture subcommand), it is offered in
// tab-completion, and it accepts a steering instruction.
func TestWeaveIsTopLevelCommand(t *testing.T) {
	found := false
	for _, c := range vaultExCmds {
		if c == "weave" {
			found = true
		}
	}
	if !found {
		t.Fatal("weave missing from vaultExCmds (tab completion)")
	}

	m := newTestVaultModel(t) // offline: the guard fires, proving dispatch reached cmdWeave
	m = openNote(t, m, "Enums.md", "# Enums\n")
	tm, _ := m.runEx("weave keep it terse")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "needs an AI provider") {
		t.Fatalf(":weave with an instruction did not reach cmdWeave: %q", m.notice)
	}

	// The old spelling is no longer a weave; it's an invalid :capture argument.
	tm, _ = m.runEx("capture weave")
	m = tm.(VaultModel)
	if !strings.Contains(m.notice, "usage:") {
		t.Fatalf("`:capture weave` should now be a usage error, got %q", m.notice)
	}
}
