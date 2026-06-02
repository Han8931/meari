package curriculum

import (
	"testing"

	"meari/internal/executor"
)

// TestChallengesAreSolvable runs every curriculum challenge's hidden reference
// Solution against its Tests and requires a pass. This guarantees the authored
// content is correct: the challenge is solvable and its tests actually check it.
func TestChallengesAreSolvable(t *testing.T) {
	for _, lang := range Languages() {
		for _, level := range []string{Beginner, Intermediate, Advanced} {
			c, ok := For(lang, level)
			if !ok {
				t.Fatalf("missing curriculum for %s/%s", lang, level)
			}
			for _, topic := range c.Topics() {
				topic := topic
				t.Run(lang+"/"+level+"/"+topic.ID, func(t *testing.T) {
					ch := topic.Challenge
					res, err := executor.Run(lang, ch.Solution, ch.Tests)
					if err != nil {
						t.Fatalf("executor error: %v", err)
					}
					if !res.Passed {
						t.Fatalf("reference solution failed its own tests:\n%s", res.Output)
					}
				})
			}
		}
	}
}

// TestTopicIDsUnique guards against duplicate topic IDs across the whole
// registry (they're used as progress keys and must be globally unique).
func TestTopicIDsUnique(t *testing.T) {
	seen := map[string]string{}
	for _, lang := range Languages() {
		for _, level := range []string{Beginner, Intermediate, Advanced} {
			c, _ := For(lang, level)
			for _, topic := range c.Topics() {
				if where, dup := seen[topic.ID]; dup {
					t.Errorf("duplicate topic ID %q (in %s and %s/%s)", topic.ID, where, lang, level)
				}
				seen[topic.ID] = lang + "/" + level
			}
		}
	}
}
