package tui

import (
	"strings"
	"testing"

	"meari/internal/tutor"
)

// A built-in course topic ships its canonical answer in the challenge. ":answer"
// must reveal it verbatim, wrapped in a fenced block, with no tutor call — so
// answers stay checkable with no AI provider configured.
func TestAnswerRevealsShippedSolutionOffline(t *testing.T) {
	sol := "fn update_score() -> i32 {\n    20\n}"
	m := Model{
		chat: newChat(),
		current: tutor.Challenge{
			ID:       "rust-b-variables",
			Lang:     "rust",
			Solution: sol,
		},
	}
	// A nil Tutor would panic if cmdAnswer fell through to the LLM path; the
	// stored solution must short-circuit before that.
	tm, cmd := m.cmdAnswer()
	if cmd != nil {
		t.Fatal("revealing a shipped answer must not issue a tutor command")
	}
	got := tm.(Model)
	var found bool
	for _, b := range got.chat.blocks {
		if b.role == roleLesson && strings.Contains(b.text, sol) {
			found = true
			if !strings.Contains(b.text, "```rust") {
				t.Errorf("code answer should be fenced with its language, got:\n%s", b.text)
			}
		}
	}
	if !found {
		t.Fatal("the shipped solution should appear in the chat transcript")
	}
}
