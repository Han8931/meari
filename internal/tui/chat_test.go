package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"meari/internal/editor"
	"meari/internal/tutor"
)

// forceColorTUI makes lipgloss emit ANSI in tests so style assertions work.
func forceColorTUI(t *testing.T) {
	t.Helper()
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.ANSI256)
	t.Cleanup(func() { lipgloss.SetColorProfile(prev) })
}

func TestChatBusyLineShowsProgress(t *testing.T) {
	c := newChat()
	c.setSize(40, 12)
	c.setBusy("tutor thinking")
	c.tickBusy()
	if !strings.Contains(c.view(), "tutor thinking…") {
		t.Fatalf("busy line missing from view:\n%s", c.view())
	}
	// The pane must not grow: busy steals a transcript row instead.
	if got := strings.Count(c.view(), "\n") + 1; got > 12 {
		t.Fatalf("view is %d rows, want <= 12", got)
	}
	c.setBusy("")
	if strings.Contains(c.view(), "thinking") {
		t.Fatal("busy line should disappear when idle")
	}
}

func TestChatInputIsMultiRowAndSubmits(t *testing.T) {
	c := newChat()
	c.setSize(30, 12)
	c.input.SetValue("a rather long question that wraps over the narrow pane width")
	got, ok := c.submit()
	if !ok || !strings.HasPrefix(got, "a rather long") {
		t.Fatalf("submit = %q ok=%v", got, ok)
	}
	if v := c.input.Value(); v != "" {
		t.Fatalf("input should clear after submit, got %q", v)
	}
	// Three input rows are reserved (height permitting).
	if c.input.Height() != chatInputRows {
		t.Fatalf("input height = %d, want %d", c.input.Height(), chatInputRows)
	}
}

func TestChatSpeakerBadges(t *testing.T) {
	forceColorTUI(t)
	c := newChat()
	c.setSize(40, 12)
	c.append(roleUser, "question")
	c.append(roleTutor, "reply")
	content := c.renderBlock(chatBlock{role: roleTutor, text: "reply"})
	if !strings.Contains(content, chatTutorBadge.Render(" tutor ")) {
		t.Fatalf("tutor badge missing:\n%q", content)
	}
	content = c.renderBlock(chatBlock{role: roleUser, text: "question"})
	if !strings.Contains(content, chatUserBadge.Render(" you ")) {
		t.Fatalf("you badge missing:\n%q", content)
	}
}

func TestChatHighlightsFencedCode(t *testing.T) {
	forceColorTUI(t)
	c := newChat()
	c.setSize(60, 12)
	body := "Here is an example:\n```python\ndef add(a, b):\n    return a + b\n```\nTry it."
	got := c.renderRichBody(body)
	if strings.Contains(got, "```") {
		t.Fatalf("fence markers should not be rendered:\n%q", got)
	}
	wantLine := editor.Highlight("python", "def add(a, b):")
	if !strings.Contains(got, wantLine) {
		t.Fatalf("code should be syntax-highlighted\nwant fragment %q\nin:\n%q", wantLine, got)
	}
	if !strings.Contains(got, "Try it.") {
		t.Fatalf("prose after the fence lost:\n%q", got)
	}
}

func TestChatCodeBlocksWrapInsteadOfClip(t *testing.T) {
	c := newChat()
	c.setSize(24, 12)
	long := "```\nresult = compute_something_quite_long(alpha, beta) # TRAILING_MARKER\n```"
	got := c.renderRichBody(long)
	// Content must survive in full — wrapped across rows, never clipped. Strip
	// the wrapping (gutters + newlines) and look for the tail marker.
	joined := strings.NewReplacer("\n", "", "│ ", "").Replace(got)
	if !strings.Contains(joined, "TRAILING_MARKER") {
		t.Fatalf("the end of a long code line must stay visible (wrapped, not clipped):\n%q", got)
	}
	// Every visual row carries the gutter, including wrapped continuations.
	for i, row := range strings.Split(got, "\n") {
		if !strings.Contains(row, "│") {
			t.Fatalf("row %d lost the code gutter: %q", i, row)
		}
	}
}

func TestVaultCompactGrowsChat(t *testing.T) {
	m := newTestVaultModel(t)
	before := m.chatW
	for i := 0; i < 4; i++ { // step to the clamp
		tm, _ := m.runEx("compact")
		m = tm.(VaultModel)
	}
	if m.chatW <= before {
		t.Fatalf(":compact should widen the chat pane: %d -> %d", before, m.chatW)
	}
	if m.chatW <= m.editorW {
		t.Fatalf("at max compact the chat should outsize the editor (chat=%d editor=%d)", m.chatW, m.editorW)
	}
}

func TestChatUnlabeledFenceUsesCodeLang(t *testing.T) {
	forceColorTUI(t)
	c := newChat()
	c.setSize(60, 12)
	c.codeLang = "go"
	got := c.renderRichBody("```\nfunc main() {}\n```")
	if !strings.Contains(got, editor.Highlight("go", "func main() {}")) {
		t.Fatalf("unlabeled fence should highlight as codeLang:\n%q", got)
	}
	// With no codeLang, unlabeled fences stay plain.
	c.codeLang = ""
	got = c.renderRichBody("```\nfunc main() {}\n```")
	if !strings.Contains(got, "func main() {}") {
		t.Fatalf("plain fence content lost:\n%q", got)
	}
}

func TestVaultPerNoteChat(t *testing.T) {
	m := newTestVaultModel(t)

	// Open note A and have some chat activity.
	a := vSaveOpenCmd(m.svc, "x/A.md", "# A\n\nbody\n")().(vOpenedMsg)
	tm, _ := m.Update(a)
	m = tm.(VaultModel)
	m.chat.append(roleUser, "question about A")

	// Open note B: the pane starts fresh (just the topic header).
	b := vSaveOpenCmd(m.svc, "x/B.md", "# B\n\nbody\n")().(vOpenedMsg)
	tm, _ = m.Update(b)
	m = tm.(VaultModel)
	if strings.Contains(m.chat.view(), "question about A") {
		t.Fatal("note B's chat should not show note A's activity")
	}
	if !strings.Contains(m.chat.view(), "— B —") {
		t.Fatalf("fresh note chat should show its header:\n%s", m.chat.view())
	}

	// Back to A: the old transcript is restored.
	tm, _ = m.Update(vOpenCmd(m.svc, "x/A.md")().(vOpenedMsg))
	m = tm.(VaultModel)
	if !strings.Contains(m.chat.view(), "question about A") {
		t.Fatal("returning to note A should restore its chat history")
	}
}

// stubClipboard replaces the real clipboard for a test and returns a pointer
// to the last copied text.
func stubClipboard(t *testing.T) *string {
	t.Helper()
	var got string
	prev := copyToClipboard
	copyToClipboard = func(text string) error {
		got = text
		return nil
	}
	t.Cleanup(func() { copyToClipboard = prev })
	return &got
}

func TestCopyChatLastReply(t *testing.T) {
	got := stubClipboard(t)
	c := newChat()
	c.setSize(60, 12)
	c.append(roleUser, "q1")
	c.append(roleTutor, "first answer")
	c.append(roleUser, "q2")
	c.append(roleTutor, "second answer")
	c.append(roleSystem, "some notice") // must be skipped

	copyChat(&c, "")
	if *got != "second answer" {
		t.Fatalf("copied %q, want the last tutor reply", *got)
	}
	if !strings.Contains(c.view(), "✓ copied last reply") {
		t.Fatalf("feedback missing:\n%s", c.view())
	}
}

func TestCopyChatCode(t *testing.T) {
	got := stubClipboard(t)
	c := newChat()
	c.setSize(60, 12)
	c.append(roleTutor, "Try this:\n```python\nx = 1\n```\nthen this:\n```python\ny = 2\n```")

	copyChat(&c, "code")
	if *got != "y = 2" {
		t.Fatalf("copied %q, want the LAST fenced block", *got)
	}
}

func TestCopyChatAllAndEmpty(t *testing.T) {
	got := stubClipboard(t)
	c := newChat()
	c.setSize(60, 12)

	copyChat(&c, "") // nothing yet
	if *got != "" {
		t.Fatalf("nothing should be copied from an empty chat, got %q", *got)
	}
	if !strings.Contains(c.view(), "nothing to copy") {
		t.Fatalf("empty-chat feedback missing:\n%s", c.view())
	}

	c.append(roleUser, "hello")
	c.append(roleTutor, "hi there")
	copyChat(&c, "all")
	if !strings.Contains(*got, "you: hello") || !strings.Contains(*got, "tutor: hi there") {
		t.Fatalf("transcript copy wrong: %q", *got)
	}
}

func TestPasteIntoChatInput(t *testing.T) {
	prev := pasteFromClipboard
	pasteFromClipboard = func() (string, error) { return "pasted question", nil }
	t.Cleanup(func() { pasteFromClipboard = prev })

	c := newChat()
	c.setSize(60, 12)
	pasteChat(&c)
	if got := c.input.Value(); got != "pasted question" {
		t.Fatalf("input after paste = %q", got)
	}

	// Empty clipboard: friendly notice, input untouched.
	pasteFromClipboard = func() (string, error) { return "  ", nil }
	c2 := newChat()
	c2.setSize(60, 12)
	pasteChat(&c2)
	if c2.input.Value() != "" {
		t.Fatalf("empty clipboard must not modify the input")
	}
	if !strings.Contains(c2.view(), "clipboard is empty") {
		t.Fatalf("empty-clipboard notice missing:\n%s", c2.view())
	}
}

func TestPasteCommandFocusesChat(t *testing.T) {
	prev := pasteFromClipboard
	pasteFromClipboard = func() (string, error) { return "from clipboard", nil }
	t.Cleanup(func() { pasteFromClipboard = prev })

	m := newTestVaultModel(t)
	tm, _ := m.runEx("paste")
	m = tm.(VaultModel)
	if m.focus != paneChat {
		t.Fatalf(":paste should focus the chat pane, got %v", m.focus)
	}
	if got := m.chat.input.Value(); got != "from clipboard" {
		t.Fatalf("chat input = %q", got)
	}
}

func TestChatInputHistory(t *testing.T) {
	c := newChat()
	c.setSize(60, 12)
	up := tea.KeyMsg{Type: tea.KeyUp}
	down := tea.KeyMsg{Type: tea.KeyDown}

	c.input.SetValue("one")
	c.submit()
	c.input.SetValue("two")
	c.submit()

	// Empty input: up recalls the latest, then the one before.
	c.histKey(up)
	if got := c.input.Value(); got != "two" {
		t.Fatalf("after up: %q", got)
	}
	c.histKey(up)
	if got := c.input.Value(); got != "one" {
		t.Fatalf("after up up: %q", got)
	}
	// Down walks back toward the (empty) live draft.
	c.histKey(down)
	if got := c.input.Value(); got != "two" {
		t.Fatalf("after down: %q", got)
	}
	c.histKey(down)
	if got := c.input.Value(); got != "" {
		t.Fatalf("back to live draft: %q", got)
	}

	// A typed draft is preserved across a recall round-trip.
	c.input.SetValue("")
	c.draft = ""
	c.histPos = len(c.inputHist)
	c.input.SetValue("") // empty -> navigable
	c.histKey(up)        // "two"
	c.histKey(down)      // live again
	if got := c.input.Value(); got != "" {
		t.Fatalf("draft restore: %q", got)
	}

	// While composing (non-empty, not a recalled entry), arrows are NOT history.
	c.input.SetValue("half-typed question")
	if c.histKey(up) {
		t.Fatal("up while composing must remain a cursor movement")
	}
	// Submitting dedupes consecutive repeats.
	c.input.SetValue("two")
	c.submit()
	if n := len(c.inputHist); n != 2 {
		t.Fatalf("consecutive duplicate should not be re-added: %v", c.inputHist)
	}
}

func TestLastFence(t *testing.T) {
	if _, ok := lastFence("no code here"); ok {
		t.Fatal("prose has no fence")
	}
	if code, ok := lastFence("```\nunterminated\nfence"); !ok || code != "unterminated\nfence" {
		t.Fatalf("unterminated fence: %q ok=%v", code, ok)
	}
}

func TestAltOCopiesInChatPane(t *testing.T) {
	got := stubClipboard(t)
	m := newTestVaultModel(t)
	tm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")}) // noop key first
	m = tm.(VaultModel)
	m.setFocus(paneChat)
	m.chat.append(roleTutor, "copy me")

	for _, k := range []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune("ø")},            // mac Option+O default
		{Type: tea.KeyRunes, Runes: []rune("o"), Alt: true}, // alt+o
	} {
		*got = ""
		tm, _ = m.Update(k)
		m = tm.(VaultModel)
		if *got != "copy me" {
			t.Fatalf("key %q should copy the last reply, got %q", k.String(), *got)
		}
	}
}

func TestSwitchChatContext(t *testing.T) {
	m := Model{
		chat:      newChat(),
		chatByKey: map[string][]chatBlock{},
		histByKey: map[string][]tutor.ChatTurn{},
	}
	m.chat.setSize(40, 12)

	// A lesson streams in before the first topic is entered (custom-topic flow).
	m.chat.append(roleLesson, "startup lesson")

	// First topic inherits the startup transcript.
	if fresh := m.switchChatContext("topic:a"); !fresh {
		t.Fatal("first visit should be fresh")
	}
	if !strings.Contains(m.chat.view(), "startup lesson") {
		t.Fatal("first topic must inherit the startup transcript")
	}
	m.chat.append(roleUser, "question on A")

	// Switching to topic B starts clean; back to A restores everything.
	if fresh := m.switchChatContext("topic:b"); !fresh {
		t.Fatal("topic B first visit should be fresh")
	}
	if strings.Contains(m.chat.view(), "question on A") {
		t.Fatal("topic B should not show topic A's chat")
	}
	if fresh := m.switchChatContext("topic:a"); fresh {
		t.Fatal("revisiting A should not be fresh")
	}
	if !strings.Contains(m.chat.view(), "question on A") ||
		!strings.Contains(m.chat.view(), "startup lesson") {
		t.Fatalf("topic A's transcript should be fully restored:\n%s", m.chat.view())
	}

	// Same key again is a no-op.
	if m.switchChatContext("topic:a") {
		t.Fatal("switching to the current key must be a no-op")
	}
}

func TestChallengeHeaderShowsShortTitleNotPrompt(t *testing.T) {
	m := Model{
		current: tutor.Challenge{
			ID:     "sum-list",
			Lang:   "python",
			Prompt: "Write a function sum_list(xs) that returns the sum of all numbers in the list xs, handling empty lists by returning 0.",
		},
		topic: "python lists",
	}
	got := m.challengeHeader(40)
	if !strings.Contains(got, "Sum list") {
		t.Fatalf("header should show the prettified id title, got:\n%q", got)
	}
	if strings.Contains(got, "Write a function") {
		t.Fatalf("header must not contain the full prompt:\n%q", got)
	}
	if !strings.Contains(got, "PYTHON") {
		t.Fatalf("header should keep the language tag:\n%q", got)
	}
}

func TestPrettyID(t *testing.T) {
	cases := map[string]string{
		"sum-list":        "Sum list",
		"offline-is-even": "Offline is even",
		"under_score":     "Under score",
		"":                "",
	}
	for in, want := range cases {
		if got := prettyID(in); got != want {
			t.Errorf("prettyID(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestVaultPerNoteTutorHistorySeparate(t *testing.T) {
	m := newTestVaultModel(t)
	a := vSaveOpenCmd(m.svc, "x/A.md", "# A\n")().(vOpenedMsg)
	tm, _ := m.Update(a)
	m = tm.(VaultModel)
	m.chatHist = append(m.chatHist, vChatTurn("user", "about A")...)

	b := vSaveOpenCmd(m.svc, "x/B.md", "# B\n")().(vOpenedMsg)
	tm, _ = m.Update(b)
	m = tm.(VaultModel)
	if len(m.chatHist) != 0 {
		t.Fatalf("note B should start with an empty tutor conversation, got %d turns", len(m.chatHist))
	}

	tm, _ = m.Update(vOpenCmd(m.svc, "x/A.md")().(vOpenedMsg))
	m = tm.(VaultModel)
	if len(m.chatHist) != 1 || m.chatHist[0].Content != "about A" {
		t.Fatalf("note A's tutor conversation should be restored, got %+v", m.chatHist)
	}
}
