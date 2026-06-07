// Package curriculum defines built-in, pre-authored learning paths: ordered
// modules of topics for a given language and level. Unlike LLM-generated topics,
// every lesson and challenge here is written and verified ahead of time, so the
// app teaches consistently and works fully offline.
//
// Each Topic carries its own Lesson (shown in the chat pane) and a Challenge
// with a hidden reference Solution that the test suite uses to prove the
// challenge is solvable and its tests are correct (see curriculum_test.go).
package curriculum

import "strings"

// Levels.
const (
	Beginner     = "beginner"
	Intermediate = "intermediate"
	Advanced     = "advanced"
)

// Challenge is a pre-authored exercise.
type Challenge struct {
	Prompt      string   // what to implement (shown in chat)
	StarterCode string   // the stub the learner edits
	Solution    string   // hidden reference answer (used only to self-verify tests)
	Tests       []string // assertions: Python asserts, or Go test-body statements
}

// QuizQuestion is one multiple-choice item of a quiz step.
type QuizQuestion struct {
	Q       string   `json:"q"`
	Choices []string `json:"choices"`
	Answer  int      `json:"answer"` // index into Choices
	Why     string   `json:"why"`    // shown after answering
}

// Topic is one unit of study with a baked lesson and challenge.
type Topic struct {
	ID     string // globally unique, e.g. "py-b-variables"
	Title  string // short sidebar label
	Lesson string // explanation + worked example, shown in the chat pane
	// Lang overrides the curriculum's language for this topic. Vault-built
	// courses mix kinds — code topics carry "python"/"go", essay topics
	// "essay", quiz steps "quiz" — while the built-in curricula leave it
	// empty (inherit).
	Lang      string
	Challenge Challenge
	// Quiz holds a quiz step's questions (Lang "quiz"); graded locally.
	Quiz []QuizQuestion
}

// Module groups related topics under a heading.
type Module struct {
	Name   string
	Topics []Topic
}

// Curriculum is the full path for one (language, level).
type Curriculum struct {
	Lang    string
	Level   string
	Modules []Module
}

// Topics flattens the curriculum into its topics in order.
func (c Curriculum) Topics() []Topic {
	var out []Topic
	for _, m := range c.Modules {
		out = append(out, m.Topics...)
	}
	return out
}

// For returns the curriculum for a language and level, or ok=false if none.
func For(lang, level string) (Curriculum, bool) {
	lang = strings.ToLower(lang)
	level = strings.ToLower(level)
	builders, ok := registry[lang]
	if !ok {
		return Curriculum{}, false
	}
	build, ok := builders[level]
	if !ok {
		return Curriculum{}, false
	}
	return Curriculum{Lang: lang, Level: level, Modules: build()}, true
}

// HasCurriculum reports whether any curriculum exists for a language.
func HasCurriculum(lang string) bool {
	_, ok := registry[strings.ToLower(lang)]
	return ok
}

// registry maps language -> level -> module builder. The built-in Go track is
// the SOURCE for the markdown course core.SeedBuiltinCourses materializes on
// first run; at runtime every course — seeded or learner-built — is loaded
// from markdown through the same vault course loader.
var registry = map[string]map[string]func() []Module{
	"go": {
		Beginner:     goBeginner,
		Intermediate: goIntermediate,
		Advanced:     goAdvanced,
	},
}

// Languages returns the built-in curriculum languages in display order.
func Languages() []string { return []string{"go"} }
