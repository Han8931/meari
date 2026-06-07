package core

// course.go is the vault-native course format: a learner-owned curriculum
// built FROM vault notes, stored AS vault notes, runnable in the tutor TUI.
//
// A course lives under meari-course/<Title>/course.md. The manifest's
// frontmatter carries id/title/level; its body is the ordered outline —
// "## Module" headings with one [[wikilink]] list item per topic, each link
// resolving to an ordinary vault note (by title, id, or filename stem, the
// same rules backlinks use):
//
//	---
//	id: balanced-search-trees
//	title: Balanced Search Trees
//	level: intermediate
//	---
//	## Foundations
//	- [[BST]]
//	- [[Tree Terminology]]
//	## Balancing
//	- [[AVL]]
//
// A topic note's body is its lesson. Its study material rides in the note's
// frontmatter under "study:" (preserved by the vault's Extra passthrough):
//
//	study:
//	  kind: code            # "code" | "essay"
//	  lang: python          # code only
//	  prompt: Implement insert() for a BST…
//	  starter: |            # code only
//	    def insert(root, val): ...
//	  tests:                # code only
//	    - assert insert(...) ...
//	  answer: |             # the model answer / reference solution
//	    ...
//
// A topic with no study block defaults to an essay on the note ("explain it
// in your own words"), so a hand-written manifest over existing notes is a
// complete, runnable course with zero extra authoring.

import (
	"fmt"
	"sort"
	"strings"

	"meari/internal/curriculum"
	"meari/internal/vault"
)

// CourseDir is the vault folder that holds generated/authored courses.
const CourseDir = "meari-course"

// CourseTopic is one unit of study: a lesson (the note body) plus one study
// item (a code challenge, an essay prompt, or a quiz step's questions).
type CourseTopic struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	NotePath string `json:"notePath"`
	Lesson   string `json:"lesson"`

	Kind      string                    `json:"kind"` // "code" | "essay" | "quiz"
	Lang      string                    `json:"lang"` // code topics: "python" | "go"
	Prompt    string                    `json:"prompt"`
	Starter   string                    `json:"starter"`
	Tests     []string                  `json:"tests"`
	Answer    string                    `json:"answer"`
	Questions []curriculum.QuizQuestion `json:"questions"` // quiz steps
}

// CourseModule groups ordered topics under a heading.
type CourseModule struct {
	Name   string        `json:"name"`
	Topics []CourseTopic `json:"topics"`
}

// Course is a full vault-backed curriculum.
type Course struct {
	ID      string         `json:"id"`
	Title   string         `json:"title"`
	Level   string         `json:"level"`
	Path    string         `json:"path"` // the manifest's vault-relative path
	Modules []CourseModule `json:"modules"`
}

// CourseMeta is the list-view of a course.
type CourseMeta struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Level string `json:"level"`
	Path  string `json:"path"`
}

// ListCourses returns every course manifest under meari-course/, sorted by title.
func (s *Service) ListCourses() ([]CourseMeta, error) {
	notes, err := s.allNotes()
	if err != nil {
		return nil, err
	}
	var out []CourseMeta
	for _, n := range notes {
		if !isCourseManifest(n.RelPath) {
			continue
		}
		out = append(out, CourseMeta{
			ID:    courseID(n),
			Title: n.Title,
			Level: stringExtra(n, "level"),
			Path:  n.RelPath,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out, nil
}

// LoadCourse loads a course by id, title, or manifest path (case-insensitive)
// and resolves its topic links against the vault.
func (s *Service) LoadCourse(key string) (Course, error) {
	notes, err := s.allNotes()
	if err != nil {
		return Course{}, err
	}
	var manifest *vault.Note
	for i, n := range notes {
		if !isCourseManifest(n.RelPath) {
			continue
		}
		if strings.EqualFold(key, courseID(notes[i])) ||
			strings.EqualFold(key, n.Title) || strings.EqualFold(key, n.RelPath) {
			manifest = &notes[i]
			break
		}
	}
	if manifest == nil {
		return Course{}, fmt.Errorf("no course %q under %s/", key, CourseDir)
	}

	c := Course{
		ID:    courseID(*manifest),
		Title: manifest.Title,
		Level: stringExtra(*manifest, "level"),
		Path:  manifest.RelPath,
	}
	if c.Level == "" {
		c.Level = curriculum.Beginner
	}

	// Walk the outline: "## Heading" opens a module, each list-item wikilink
	// adds the linked note as a topic. Unresolvable links are skipped — a
	// course must stay loadable after a note is renamed away.
	cur := CourseModule{Name: "Topics"} // links before any heading land here
	for _, line := range strings.Split(manifest.Body, "\n") {
		t := strings.TrimSpace(line)
		if name, ok := strings.CutPrefix(t, "## "); ok {
			if len(cur.Topics) > 0 {
				c.Modules = append(c.Modules, cur)
			}
			cur = CourseModule{Name: strings.TrimSpace(name)}
			continue
		}
		if !strings.HasPrefix(t, "- ") && !strings.HasPrefix(t, "* ") {
			continue
		}
		for _, l := range vault.ParseLinks(t) {
			if n, ok := resolveLink(l.Target, notes); ok {
				cur.Topics = append(cur.Topics, courseTopic(c.ID, n))
			}
		}
	}
	if len(cur.Topics) > 0 {
		c.Modules = append(c.Modules, cur)
	}
	if len(c.Modules) == 0 {
		return Course{}, fmt.Errorf("course %q has no resolvable topics", c.Title)
	}
	return c, nil
}

// Curriculum converts the course into the tutor's runnable curriculum form.
// Essay topics use the "essay" reflection language (graded prose, like the
// built-in physics track); code topics carry their own language.
func (c Course) Curriculum() curriculum.Curriculum {
	out := curriculum.Curriculum{Lang: c.ID, Level: c.Level}
	for _, m := range c.Modules {
		mod := curriculum.Module{Name: m.Name}
		for _, t := range m.Topics {
			lang := "essay"
			switch {
			case t.Kind == "quiz" && len(t.Questions) > 0:
				lang = "quiz"
			case t.Kind == "code":
				lang = t.Lang
				if lang == "" {
					lang = "python"
				}
			}
			mod.Topics = append(mod.Topics, curriculum.Topic{
				ID:     t.ID,
				Title:  t.Title,
				Lesson: t.Lesson,
				Lang:   lang,
				Quiz:   t.Questions,
				Challenge: curriculum.Challenge{
					Prompt:      t.Prompt,
					StarterCode: t.Starter,
					Tests:       t.Tests,
					Solution:    t.Answer,
				},
			})
		}
		out.Modules = append(out.Modules, mod)
	}
	return out
}

// --- helpers ---

func isCourseManifest(relPath string) bool {
	return strings.HasPrefix(relPath, CourseDir+"/") &&
		strings.HasSuffix(strings.ToLower(relPath), "/course.md")
}

// courseID prefers the manifest's explicit id, falling back to its folder name.
func courseID(n vault.Note) string {
	if n.ID != "" {
		return n.ID
	}
	parts := strings.Split(n.RelPath, "/")
	if len(parts) >= 2 {
		return vault.Slug(parts[len(parts)-2])
	}
	return vault.Slug(n.Title)
}

// resolveLink finds the note a wikilink target refers to.
func resolveLink(target string, notes []vault.Note) (vault.Note, bool) {
	for _, n := range notes {
		if linkMatches(target, n2target(n)) {
			return n, true
		}
	}
	return vault.Note{}, false
}

// courseTopic builds a CourseTopic from a resolved note, reading the study
// block from its frontmatter and defaulting to an essay on the note.
func courseTopic(courseID string, n vault.Note) CourseTopic {
	t := CourseTopic{
		ID:       "course-" + courseID + "-" + vault.Slug(n.Title),
		Title:    n.Title,
		NotePath: n.RelPath,
		Lesson:   n.Body,
		Kind:     "essay",
	}
	study, _ := n.Extra["study"].(map[string]any)
	get := func(k string) string {
		if v, ok := study[k]; ok {
			return strings.TrimSpace(fmt.Sprint(v))
		}
		return ""
	}
	if k := get("kind"); k != "" {
		t.Kind = k
	}
	t.Lang = get("lang")
	t.Prompt = get("prompt")
	t.Starter = get("starter")
	t.Answer = get("answer")
	if tests, ok := study["tests"].([]any); ok {
		for _, v := range tests {
			if s := strings.TrimSpace(fmt.Sprint(v)); s != "" {
				t.Tests = append(t.Tests, s)
			}
		}
	}
	if qs, ok := study["questions"].([]any); ok {
		t.Questions = parseQuizQuestions(qs)
	}
	if t.Kind == "quiz" && len(t.Questions) == 0 {
		t.Kind = "essay" // a quiz step without questions degrades to an essay
	}
	if t.Prompt == "" {
		t.Prompt = "Explain the key ideas of \"" + n.Title +
			"\" in your own words, with concrete examples."
	}
	if t.Kind == "essay" && t.Starter == "" {
		t.Starter = "Write your answer here:\n\n"
	}
	return t
}

// parseQuizQuestions decodes the study.questions frontmatter list, dropping
// malformed entries (a question needs text, 2+ choices, and a valid answer
// index) — a broken question must never reach the learner.
func parseQuizQuestions(raw []any) []curriculum.QuizQuestion {
	var out []curriculum.QuizQuestion
	for _, v := range raw {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		q := curriculum.QuizQuestion{}
		if s, ok := m["q"]; ok {
			q.Q = strings.TrimSpace(fmt.Sprint(s))
		}
		if s, ok := m["why"]; ok {
			q.Why = strings.TrimSpace(fmt.Sprint(s))
		}
		if cs, ok := m["choices"].([]any); ok {
			for _, c := range cs {
				q.Choices = append(q.Choices, strings.TrimSpace(fmt.Sprint(c)))
			}
		}
		if a, ok := m["answer"]; ok {
			fmt.Sscanf(fmt.Sprint(a), "%d", &q.Answer)
		}
		if q.Q != "" && len(q.Choices) >= 2 && q.Answer >= 0 && q.Answer < len(q.Choices) {
			out = append(out, q)
		}
	}
	return out
}

func stringExtra(n vault.Note, key string) string {
	if v, ok := n.Extra[key]; ok {
		return strings.TrimSpace(fmt.Sprint(v))
	}
	return ""
}
