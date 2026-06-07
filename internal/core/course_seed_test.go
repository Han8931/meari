package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"meari/internal/curriculum"
)

func TestSeedBuiltinCoursesMatchesCurriculum(t *testing.T) {
	svc := newTestService(t)
	svc.SetCourseDir(t.TempDir())
	if err := svc.SeedBuiltinCourses(); err != nil {
		t.Fatal(err)
	}

	metas, err := svc.ListCourses()
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 3 {
		t.Fatalf("seeded %d courses, want 3: %+v", len(metas), metas)
	}
	// The levels of one course family nest under a single Go/ directory.
	paths := map[string]bool{}
	for _, m := range metas {
		paths[m.Path] = true
	}
	for _, want := range []string{
		"meari-course/Go/Beginner/course.md",
		"meari-course/Go/Intermediate/course.md",
		"meari-course/Go/Advanced/course.md",
	} {
		if !paths[want] {
			t.Errorf("seeded manifest missing: %s (got %+v)", want, metas)
		}
	}

	for _, b := range []struct{ id, level string }{
		{"go-beginner", curriculum.Beginner},
		{"go-intermediate", curriculum.Intermediate},
		{"go-advanced", curriculum.Advanced},
	} {
		c, err := svc.LoadCourse(b.id)
		if err != nil {
			t.Fatalf("LoadCourse(%s): %v", b.id, err)
		}
		src, _ := curriculum.For("go", b.level)
		want := len(src.Topics())
		got := 0
		for _, m := range c.Modules {
			got += len(m.Topics)
		}
		if got != want {
			t.Errorf("%s: %d topics, want %d (every outline link must resolve)", b.id, got, want)
		}
		if c.Level != b.level {
			t.Errorf("%s: level = %q, want %q", b.id, c.Level, b.level)
		}
	}

	// The markdown round-trip preserves the study material: the first
	// beginner topic keeps its prompt, starter, tests, and solution. (The
	// course loader TrimSpaces study values and the vault normalizes the
	// body's final newline, for seeded and generated courses alike — so
	// compare modulo edge whitespace.)
	c, _ := svc.LoadCourse("go-beginner")
	first := c.Modules[0].Topics[0]
	src, _ := curriculum.For("go", curriculum.Beginner)
	want := src.Modules[0].Topics[0]
	if first.Kind != "code" || first.Lang != "go" {
		t.Errorf("first topic kind/lang = %q/%q, want code/go", first.Kind, first.Lang)
	}
	eq := func(field, got, want string) {
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			t.Errorf("%s drifted:\n got %q\nwant %q", field, got, want)
		}
	}
	eq("prompt", first.Prompt, want.Challenge.Prompt)
	eq("starter", first.Starter, want.Challenge.StarterCode)
	eq("answer", first.Answer, want.Challenge.Solution)
	eq("lesson", first.Lesson, want.Lesson)
	if len(first.Tests) != len(want.Challenge.Tests) {
		t.Errorf("tests drifted: %d, want %d", len(first.Tests), len(want.Challenge.Tests))
	}
}

func TestSeedBuiltinCoursesRunsOnce(t *testing.T) {
	svc := newTestService(t)
	dir := t.TempDir()
	svc.SetCourseDir(dir)
	if err := svc.SeedBuiltinCourses(); err != nil {
		t.Fatal(err)
	}
	// The learner deletes a seeded course; a second seed must not revive it.
	if err := svc.Delete(CourseDir + "/Go/Advanced"); err != nil {
		t.Fatal(err)
	}
	if err := svc.SeedBuiltinCourses(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "Go", "Advanced")); !os.IsNotExist(err) {
		t.Fatal("a deleted seeded course must stay deleted")
	}
	metas, _ := svc.ListCourses()
	if len(metas) != 2 {
		t.Fatalf("got %d courses after delete + re-seed, want 2", len(metas))
	}
}
