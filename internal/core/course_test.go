package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meari/internal/config"
	"meari/internal/tutor"
	"meari/internal/vault"
)

func newCourseVault(t *testing.T) *Service {
	t.Helper()
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(v, tutor.New(config.AIConfig{Provider: "openai"})) // offline

	write := func(path, body string) {
		t.Helper()
		if _, err := svc.SaveNote(path, body); err != nil {
			t.Fatal(err)
		}
	}
	// Topic notes: one with a code study block, one with an essay block, one
	// with no study block at all (defaults to an essay).
	if _, err := v.Write(vault.Note{
		RelPath: "Algo/BST.md",
		Title:   "BST",
		Body:    "# BST\n\nA binary search tree keeps keys ordered.\n",
		Extra: map[string]any{"study": map[string]any{
			"kind":    "code",
			"lang":    "python",
			"prompt":  "Implement insert(root, val).",
			"starter": "def insert(root, val):\n    pass\n",
			"tests":   []any{"assert insert(None, 1) is not None"},
			"answer":  "def insert(root, val): ...",
		}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Write(vault.Note{
		RelPath: "Algo/Tree Terminology.md",
		Title:   "Tree Terminology",
		Body:    "Roots, leaves, depth.\n",
		Extra: map[string]any{"study": map[string]any{
			"kind":   "essay",
			"prompt": "Define root, leaf, and depth.",
			"answer": "The root is…",
		}},
	}); err != nil {
		t.Fatal(err)
	}
	write("Algo/AVL.md", "# AVL\n\nSelf-balancing BST.\n")

	if _, err := v.Write(vault.Note{
		RelPath: CourseDir + "/Balanced Trees/course.md",
		ID:      "balanced-trees",
		Title:   "Balanced Trees",
		Extra:   map[string]any{"level": "intermediate"},
		Body: "## Foundations\n" +
			"- [[Tree Terminology]]\n" +
			"- [[BST]]\n" +
			"\n## Balancing\n" +
			"- [[AVL]]\n" +
			"- [[No Such Note]]\n", // unresolvable: skipped, not fatal
	}); err != nil {
		t.Fatal(err)
	}
	return svc
}

// Difficulty variants of one course can nest under a shared directory; a
// hand-authored manifest without an explicit id derives a unique one from its
// full folder path, keeping the manifest's level for display.
func TestNestedCourseLevels(t *testing.T) {
	svc := newCourseVault(t)
	courseDir := t.TempDir()
	svc.SetCourseDir(courseDir)
	if _, err := svc.SaveNote("Notes/Pointers.md", "# Pointers\n\nAddresses.\n"); err != nil {
		t.Fatal(err)
	}
	// Hand-authored on disk, no explicit ids: the folder-path fallback applies.
	write := func(rel, body string) {
		t.Helper()
		abs := filepath.Join(courseDir, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("Rust/Beginner/course.md", "---\ntitle: Rust (Beginner)\n---\n## Basics\n- [[Pointers]]\n")
	write("Rust/Advanced/course.md", "---\ntitle: Rust (Advanced)\nlevel: advanced\n---\n## Deep\n- [[Pointers]]\n")

	metas, err := svc.ListCourses()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]CourseMeta{}
	for _, m := range metas {
		byID[m.ID] = m
	}
	if _, ok := byID["rust-beginner"]; !ok {
		t.Fatalf("nested id-less manifest should get id rust-beginner, got %+v", metas)
	}
	if _, ok := byID["rust-advanced"]; !ok {
		t.Fatalf("nested id-less manifest should get id rust-advanced, got %+v", metas)
	}
	if _, err := svc.LoadCourse("rust-advanced"); err != nil {
		t.Fatalf("LoadCourse(rust-advanced): %v", err)
	}
	if byID["rust-advanced"].Level != "advanced" {
		t.Fatalf("nested manifest level lost: %+v", byID["rust-advanced"])
	}
}

func TestListAndLoadCourse(t *testing.T) {
	svc := newCourseVault(t)

	metas, err := svc.ListCourses()
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 1 || metas[0].ID != "balanced-trees" || metas[0].Level != "intermediate" {
		t.Fatalf("ListCourses = %+v", metas)
	}

	// Loadable by id, title, or path.
	for _, key := range []string{"balanced-trees", "Balanced Trees", metas[0].Path} {
		if _, err := svc.LoadCourse(key); err != nil {
			t.Fatalf("LoadCourse(%q): %v", key, err)
		}
	}

	c, err := svc.LoadCourse("balanced-trees")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Modules) != 2 || c.Modules[0].Name != "Foundations" || c.Modules[1].Name != "Balancing" {
		t.Fatalf("modules = %+v", c.Modules)
	}
	if got := len(c.Modules[0].Topics); got != 2 {
		t.Fatalf("Foundations has %d topics, want 2", got)
	}
	if got := len(c.Modules[1].Topics); got != 1 { // the dead link was skipped
		t.Fatalf("Balancing has %d topics, want 1", got)
	}

	bst := c.Modules[0].Topics[1]
	if bst.Kind != "code" || bst.Lang != "python" || len(bst.Tests) != 1 ||
		!strings.Contains(bst.Prompt, "insert") {
		t.Fatalf("BST topic mis-parsed: %+v", bst)
	}
	if !strings.Contains(bst.Lesson, "keeps keys ordered") {
		t.Fatalf("lesson should be the note body, got %q", bst.Lesson)
	}

	// The bare note defaulted to an essay with a usable prompt.
	avl := c.Modules[1].Topics[0]
	if avl.Kind != "essay" || avl.Prompt == "" || avl.Starter == "" {
		t.Fatalf("default essay topic wrong: %+v", avl)
	}
}

func TestCourseToCurriculum(t *testing.T) {
	svc := newCourseVault(t)
	c, err := svc.LoadCourse("balanced-trees")
	if err != nil {
		t.Fatal(err)
	}
	cur := c.Curriculum()
	if cur.Lang != "balanced-trees" || cur.Level != "intermediate" {
		t.Fatalf("curriculum header = %q/%q", cur.Lang, cur.Level)
	}
	topics := cur.Topics()
	if len(topics) != 3 {
		t.Fatalf("topics = %d, want 3", len(topics))
	}
	seen := map[string]bool{}
	for _, tp := range topics {
		if !strings.HasPrefix(tp.ID, "course-balanced-trees-") {
			t.Fatalf("topic id %q not namespaced", tp.ID)
		}
		if seen[tp.ID] {
			t.Fatalf("duplicate topic id %q", tp.ID)
		}
		seen[tp.ID] = true
	}
	// Essay topics run the prose path; the code topic keeps python.
	if topics[0].Lang != "essay" || topics[1].Lang != "python" || topics[2].Lang != "essay" {
		t.Fatalf("langs = %q %q %q", topics[0].Lang, topics[1].Lang, topics[2].Lang)
	}
	if topics[1].Challenge.Tests == nil || topics[1].Challenge.Solution == "" {
		t.Fatalf("code challenge incomplete: %+v", topics[1].Challenge)
	}
}
