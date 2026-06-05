package curriculum

// consistency_test.go audits every pre-authored topic for internal agreement:
// the challenge PROMPT (shown in chat), the STARTER CODE (loaded into the
// editor), the hidden SOLUTION, and the TESTS must all talk about the same
// functions. A topic whose prompt asks for one thing while the editor stub
// defines another is unsolvable-by-reading and confusing.

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var (
	pyDefRE = regexp.MustCompile(`(?m)^\s*def\s+([A-Za-z_]\w*)\s*\(`)
	goDefRE = regexp.MustCompile(`(?m)^\s*func\s+(?:\([^)]*\)\s*)?([A-Za-z_]\w*)\s*\(`)
	goTypRE = regexp.MustCompile(`(?m)^\s*type\s+([A-Za-z_]\w*)\s+`)
)

// definedNames extracts top-level function (and Go type) names from source.
func definedNames(lang, src string) []string {
	var names []string
	switch lang {
	case "go":
		for _, m := range goDefRE.FindAllStringSubmatch(src, -1) {
			names = append(names, m[1])
		}
		for _, m := range goTypRE.FindAllStringSubmatch(src, -1) {
			names = append(names, m[1])
		}
	default: // python
		for _, m := range pyDefRE.FindAllStringSubmatch(src, -1) {
			names = append(names, m[1])
		}
	}
	return names
}

func auditCurriculum(t *testing.T, lang string) {
	t.Helper()
	for _, level := range []string{Beginner, Intermediate, Advanced} {
		c, ok := For(lang, level)
		if !ok {
			t.Fatalf("missing curriculum %s/%s", lang, level)
		}
		for _, topic := range c.Topics() {
			topic := topic
			t.Run(fmt.Sprintf("%s/%s/%s", lang, level, topic.ID), func(t *testing.T) {
				ch := topic.Challenge

				starter := definedNames(lang, ch.StarterCode)
				solution := definedNames(lang, ch.Solution)

				// 1. Every name the learner is given a stub for must be named in
				//    the prompt — otherwise the instruction doesn't match the
				//    editor contents. Dunder methods (__init__ …) are plumbing
				//    that starters pre-implement; prompts needn't mention them.
				for _, name := range starter {
					if strings.HasPrefix(name, "__") {
						continue
					}
					if !strings.Contains(ch.Prompt, name) {
						t.Errorf("starter defines %q but the prompt never mentions it\nprompt: %s",
							name, ch.Prompt)
					}
				}

				// 2. The solution must implement exactly the stubs the starter
				//    shows (same names), or the learner solves a different
				//    problem than the editor presents.
				for _, name := range starter {
					if !contains(solution, name) {
						t.Errorf("starter defines %q but the solution does not", name)
					}
				}

				// 3. Every solution-defined name exercised by the tests must
				//    have a stub in the starter, so the editor shows what the
				//    tests will call.
				for _, name := range solution {
					used := false
					for _, tc := range ch.Tests {
						if strings.Contains(tc, name) {
							used = true
							break
						}
					}
					if used && !contains(starter, name) {
						t.Errorf("tests call %q (defined in the solution) but the starter has no stub for it", name)
					}
				}

				// 4. A code topic must give the learner a stub at all.
				if len(starter) == 0 && strings.TrimSpace(ch.StarterCode) == "" {
					t.Errorf("topic has no starter code")
				}
			})
		}
	}
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}

func TestPythonCurriculumConsistency(t *testing.T) { auditCurriculum(t, "python") }
func TestGoCurriculumConsistency(t *testing.T)     { auditCurriculum(t, "go") }
