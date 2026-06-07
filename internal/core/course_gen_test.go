package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meari/internal/config"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// The offline pipeline must still produce a loadable course: one essay topic
// per source note, manifest written under meari-course/.
func TestGenerateCourseOffline(t *testing.T) {
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(v, tutor.New(config.AIConfig{Provider: "openai"})) // offline

	mustWrite := func(n vault.Note) {
		t.Helper()
		if _, err := v.Write(n); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite(vault.Note{
		RelPath: "Git/Branching.md", Title: "Branching",
		Body: "Branches diverge history. See [[Merging]] and [[Rebasing]].\n",
	})
	mustWrite(vault.Note{
		RelPath: "Git/Merging.md", Title: "Merging",
		Body: "Merging joins branches.\n",
	})
	// A backlink neighbor, with a pre-existing study block that must survive.
	mustWrite(vault.Note{
		RelPath: "Git/Rebasing.md", Title: "Rebasing",
		Body: "Rebasing replays commits onto [[Branching]].\n",
		Extra: map[string]any{"study": map[string]any{
			"kind": "essay", "prompt": "KEEP-ME", "answer": "kept",
		}},
	})

	var steps []string
	meta, err := svc.GenerateCourse(context.Background(), CourseRequest{
		NotePath:      "Git/Branching.md",
		IncludeLinked: true,
	}, func(s string) { steps = append(steps, s) })
	if err != nil {
		t.Fatal(err)
	}
	if meta.ID == "" || !strings.HasPrefix(meta.Path, CourseDir+"/") {
		t.Fatalf("meta = %+v", meta)
	}
	if len(steps) < 3 {
		t.Fatalf("expected progress reports, got %v", steps)
	}

	// The generated course loads and covers seed + both neighbors.
	c, err := svc.LoadCourse(meta.ID)
	if err != nil {
		t.Fatal(err)
	}
	var titles []string
	for _, m := range c.Modules {
		for _, tp := range m.Topics {
			titles = append(titles, tp.Title)
		}
	}
	joined := strings.Join(titles, " ")
	for _, want := range []string{"Branching", "Merging", "Rebasing"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("course misses %q: %v", want, titles)
		}
	}
	// And it converts to a runnable curriculum.
	if got := len(c.Curriculum().Topics()); got != len(titles) {
		t.Fatalf("curriculum topics = %d, want %d", got, len(titles))
	}

	// The pre-existing study block on Rebasing was preserved, not clobbered.
	n, err := v.Read("Git/Rebasing.md")
	if err != nil {
		t.Fatal(err)
	}
	study, _ := n.Extra["study"].(map[string]any)
	if study == nil || study["prompt"] != "KEEP-ME" {
		t.Fatalf("existing study block clobbered: %+v", n.Extra)
	}

	// Cancelation propagates.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := svc.GenerateCourse(ctx, CourseRequest{NotePath: "Git/Branching.md"}, nil); err == nil {
		t.Fatal("canceled context should abort generation")
	}
}

func TestParseCourseRequest(t *testing.T) {
	reply := "Great — I'll build an intermediate course including the linked notes.\n" +
		`{"course_request": {"title": "Balanced Trees", "level": "intermediate", "includeLinked": true, "extra": "focus on rotations"}}`
	req, ok := ParseCourseRequest(reply)
	if !ok || req.Title != "Balanced Trees" || req.Level != "intermediate" ||
		!req.IncludeLinked || req.Extra != "focus on rotations" {
		t.Fatalf("ParseCourseRequest = %+v, %v", req, ok)
	}

	// Mid-conversation replies (no JSON yet) must not end the intake.
	for _, s := range []string{
		"1. What difficulty? 2. Include linked notes?",
		"the course_request will come later",
		`{"course_request": broken`,
	} {
		if _, ok := ParseCourseRequest(s); ok {
			t.Fatalf("ParseCourseRequest(%q) should not match", s)
		}
	}
}

// The code verifier runs generated exercises against the REAL executor: a
// passing item survives untouched, a failing one (offline = no repair
// available) demotes to an essay so a broken challenge can never ship.
func TestVerifyCodeItemAgainstExecutor(t *testing.T) {
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(v, tutor.New(config.AIConfig{Provider: "openai"})) // offline

	good := tutor.StudyItem{
		Kind: "code", Lang: "python",
		Prompt: "Implement add.",
		Answer: "def add(a, b):\n    return a + b\n",
		Tests:  []string{"assert add(1, 2) == 3"},
	}
	var msgs []string
	report := func(f string, args ...any) { msgs = append(msgs, fmt.Sprintf(f, args...)) }

	if got := svc.verifyCodeItem(context.Background(), "Add", good, report); got.Kind != "code" {
		t.Fatalf("passing item demoted: %+v (%v)", got, msgs)
	}

	bad := good
	bad.Tests = []string{"assert add(1, 2) == 4"} // reference can never pass
	got := svc.verifyCodeItem(context.Background(), "Add", bad, report)
	if got.Kind != "essay" {
		t.Fatalf("failing item not demoted: %+v", got)
	}
	if got.Prompt != bad.Prompt {
		t.Fatalf("demotion lost the prompt: %+v", got)
	}

	// No tests at all → demoted without touching the executor.
	if got := svc.verifyCodeItem(context.Background(), "X",
		tutor.StudyItem{Kind: "code", Lang: "python", Prompt: "p"}, report); got.Kind != "essay" {
		t.Fatalf("testless item not demoted: %+v", got)
	}
}

func TestStripDeadLinks(t *testing.T) {
	ok := func(target string) bool { return target == "Real Note" || target == "Topic A" }
	body := "See [[Real Note]] and [[Ghost]] plus [[Phantom|the phantom]] and [[Topic A|this]]."
	got, dead := stripDeadLinks(body, ok)
	want := "See [[Real Note]] and Ghost plus the phantom and [[Topic A|this]]."
	if got != want || dead != 2 {
		t.Fatalf("stripDeadLinks = %q (%d dead), want %q (2 dead)", got, dead, want)
	}
}

// ReviseCourse in maintenance mode: dead links in generated lessons are
// stripped, broken code exercises demote, the manifest is rewritten, and the
// course keeps its id and folder.
func TestReviseCourseOffline(t *testing.T) {
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(v, tutor.New(config.AIConfig{Provider: "openai"})) // offline

	mustWrite := func(n vault.Note) {
		t.Helper()
		if _, err := v.Write(n); err != nil {
			t.Fatal(err)
		}
	}
	// A generated-style lesson with a dead link, and a broken code exercise.
	mustWrite(vault.Note{
		RelPath: CourseDir + "/Trees/Rotations.md", Title: "Rotations",
		Source: "meari-course",
		Body:   "Rotations rebalance. See [[Ghost Note]] and [[Heights]].\n",
		Extra: map[string]any{"study": map[string]any{
			"kind": "code", "lang": "python",
			"prompt": "Implement rot().",
			"tests":  []any{"assert rot() == 1"},
			"answer": "def rot():\n    return 2\n", // fails its own test
		}},
	})
	mustWrite(vault.Note{
		RelPath: CourseDir + "/Trees/Heights.md", Title: "Heights",
		Source: "meari-course", Body: "Height of a tree.\n",
		Extra: map[string]any{"study": map[string]any{
			"kind": "essay", "prompt": "Define height.", "answer": "…",
		}},
	})
	mustWrite(vault.Note{
		RelPath: CourseDir + "/Trees/course.md", ID: "trees", Title: "Trees",
		Source: "meari-course",
		Extra:  map[string]any{"level": "beginner"},
		Body:   "## All\n- [[Rotations]]\n- [[Heights]]\n",
	})

	var steps []string
	meta, err := svc.ReviseCourse(context.Background(), "trees", "",
		func(s string) { steps = append(steps, s) })
	if err != nil {
		t.Fatal(err)
	}
	if meta.ID != "trees" || meta.Path != CourseDir+"/Trees/course.md" {
		t.Fatalf("revision changed the course identity: %+v", meta)
	}

	// Dead link stripped, valid link kept.
	n, err := v.Read(CourseDir + "/Trees/Rotations.md")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(n.Body, "[[Ghost Note]]") || !strings.Contains(n.Body, "Ghost Note") {
		t.Fatalf("dead link not unlinked: %q", n.Body)
	}
	if !strings.Contains(n.Body, "[[Heights]]") {
		t.Fatalf("valid link removed: %q", n.Body)
	}
	// The failing code exercise was demoted to an essay (offline = no repair).
	study, _ := n.Extra["study"].(map[string]any)
	if study == nil || study["kind"] != "essay" {
		t.Fatalf("broken exercise not demoted: %+v", study)
	}
	// Still loadable afterwards.
	if _, err := svc.LoadCourse("trees"); err != nil {
		t.Fatal(err)
	}
	if len(steps) < 3 {
		t.Fatalf("expected progress, got %v", steps)
	}
}

// With SetCourseDir, course material lives OUTSIDE the notes vault (the app
// project dir) yet keeps its virtual meari-course/ paths: generation writes
// there, courses load, the tree mounts it, and the vault stays clean.
func TestCourseDirOutsideVault(t *testing.T) {
	vaultDir, appDir := t.TempDir(), t.TempDir()
	v, err := vault.Open(vaultDir)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(v, tutor.New(config.AIConfig{Provider: "openai"})) // offline
	courseDir := filepath.Join(appDir, CourseDir)
	svc.SetCourseDir(courseDir)

	if _, err := v.Write(vault.Note{
		RelPath: "Git/Branching.md", Title: "Branching", Body: "Branches diverge.\n",
	}); err != nil {
		t.Fatal(err)
	}

	meta, err := svc.GenerateCourse(context.Background(),
		CourseRequest{NotePath: "Git/Branching.md"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Physically in the app dir, NOT in the vault.
	if _, err := os.Stat(filepath.Join(courseDir, "Branching Course", "course.md")); err != nil {
		t.Fatalf("manifest not in the app dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(vaultDir, CourseDir)); !os.IsNotExist(err) {
		t.Fatal("the notes vault must stay free of course material")
	}

	// Virtually addressable as before: load, open, tree.
	if _, err := svc.LoadCourse(meta.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.OpenNote(meta.Path); err != nil {
		t.Fatalf("OpenNote(%s): %v", meta.Path, err)
	}
	tree, err := svc.Tree()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range tree {
		if e.Path == CourseDir && e.Dir {
			found = true
		}
	}
	if !found {
		t.Fatalf("tree does not mount %s/: %+v", CourseDir, tree)
	}

	// NERDTree ops route through the mount too.
	if err := svc.Rename(meta.Path, CourseDir+"/Branching Course/renamed.md"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Delete(CourseDir + "/Branching Course"); err != nil {
		t.Fatal(err)
	}
	if metas, _ := svc.ListCourses(); len(metas) != 0 {
		t.Fatalf("course should be gone, got %+v", metas)
	}
}
