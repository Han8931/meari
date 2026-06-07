package core

// course_seed.go materializes the built-in Go track as ordinary markdown
// courses on first run, so EVERY course in the app — shipped or learner-built —
// is the same thing: a plain, editable, revisable, publishable markdown folder
// under the courses mount. The in-code curriculum (internal/curriculum) is only
// the seed's source of truth, where its solutions stay executor-verified by
// tests; at runtime all courses load through the vault course loader.
//
// Seeding runs once: a marker file in the courses directory records that it
// happened, so a learner who deletes (or prunes) a seeded course never has it
// resurrected on the next launch.

import (
	"os"
	"path/filepath"
	"strings"

	"meari/internal/curriculum"
	"meari/internal/vault"
)

// seedMarkerName is created in the courses directory after the one-time seed.
const seedMarkerName = ".seeded-builtin"

// builtinCourses names the seeded courses: one per level of the Go track, all
// nested under one Go/ directory so the difficulty levels of a family of
// courses live side by side (folder = "Go/<Level>").
var builtinCourses = []struct {
	id, title, folder, level string
}{
	{"go-beginner", "Go (Beginner)", "Go/Beginner", curriculum.Beginner},
	{"go-intermediate", "Go (Intermediate)", "Go/Intermediate", curriculum.Intermediate},
	{"go-advanced", "Go (Advanced)", "Go/Advanced", curriculum.Advanced},
}

// SeedBuiltinCourses writes the built-in Go track as markdown courses — the
// same format :course produces — unless it already ran once.
func (s *Service) SeedBuiltinCourses() error {
	marker := filepath.Join(s.coursesDir, seedMarkerName)
	if _, err := os.Stat(marker); err == nil {
		return nil
	}
	for _, b := range builtinCourses {
		cur, ok := curriculum.For("go", b.level)
		if !ok {
			continue
		}
		if err := s.writeSeedCourse(b.id, b.title, b.folder, b.level, cur); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(s.coursesDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(marker, []byte("built-in courses seeded; delete freely — they will not be re-created\n"), 0o644)
}

// writeSeedCourse writes one curriculum as a course folder: a topic note per
// challenge (lesson body + study: frontmatter) and a course.md manifest whose
// outline links topics by their stable ids ([[id|Title]]), so a learner note
// with the same title can never shadow a seeded topic.
func (s *Service) writeSeedCourse(id, title, dir, level string, cur curriculum.Curriculum) error {
	folder := CourseDir + "/" + dir
	var outline strings.Builder
	for _, mod := range cur.Modules {
		outline.WriteString("## " + mod.Name + "\n")
		for _, t := range mod.Topics {
			study := map[string]any{
				"kind": "code",
				"lang": "go",
			}
			if t.Challenge.Prompt != "" {
				study["prompt"] = t.Challenge.Prompt
			}
			if t.Challenge.StarterCode != "" {
				study["starter"] = t.Challenge.StarterCode
			}
			if len(t.Challenge.Tests) > 0 {
				tests := make([]any, 0, len(t.Challenge.Tests))
				for _, tc := range t.Challenge.Tests {
					tests = append(tests, tc)
				}
				study["tests"] = tests
			}
			if t.Challenge.Solution != "" {
				study["answer"] = t.Challenge.Solution
			}
			if len(t.Quiz) > 0 {
				study["kind"] = "quiz"
				qs := make([]any, 0, len(t.Quiz))
				for _, q := range t.Quiz {
					choices := make([]any, 0, len(q.Choices))
					for _, c := range q.Choices {
						choices = append(choices, c)
					}
					qs = append(qs, map[string]any{
						"q": q.Q, "choices": choices, "answer": q.Answer, "why": q.Why,
					})
				}
				study["questions"] = qs
			}
			note := vault.Note{
				RelPath: folder + "/" + vault.CleanFilename(t.Title) + ".md",
				ID:      t.ID,
				Title:   t.Title,
				Source:  "imported:builtin-go",
				Body:    t.Lesson,
				Extra:   map[string]any{"study": study},
			}
			if _, err := s.writeNote(note); err != nil {
				return err
			}
			outline.WriteString("- [[" + t.ID + "|" + t.Title + "]]\n")
		}
		outline.WriteString("\n")
	}
	manifest := vault.Note{
		RelPath: folder + "/course.md",
		ID:      id,
		Title:   title,
		Source:  "imported:builtin-go",
		Body:    outline.String(),
		Extra:   map[string]any{"level": level},
	}
	_, err := s.writeNote(manifest)
	return err
}
