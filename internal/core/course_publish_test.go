package core

import (
	"os"
	"path/filepath"
	"testing"

	"meari/internal/config"
	"meari/internal/tutor"
	"meari/internal/vault"
)

func TestPublishCourseSelfContained(t *testing.T) {
	svc := newCourseVault(t)
	dest := t.TempDir()

	res, err := svc.PublishCourse("balanced-trees", dest)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dest, "Balanced Trees"); res.Dir != want {
		t.Fatalf("published dir = %q, want %q", res.Dir, want)
	}
	// The manifest plus the three resolvable topics — all of which live
	// OUTSIDE the course folder in this fixture, so they must be copied in.
	if res.Notes != 4 {
		t.Fatalf("published %d notes, want 4", res.Notes)
	}
	for _, f := range []string{"course.md", "BST.md", "Tree Terminology.md", "AVL.md"} {
		if _, err := os.Stat(filepath.Join(res.Dir, f)); err != nil {
			t.Errorf("published file missing: %s (%v)", f, err)
		}
	}

	// Self-contained: a fresh service over an EMPTY vault, with the courses
	// mount pointed at the published directory, loads the course fully —
	// exactly what a recipient of the shared git repo would do.
	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	recipient := New(v, tutor.New(config.AIConfig{Provider: "openai"}))
	recipient.SetCourseDir(dest)
	c, err := recipient.LoadCourse("balanced-trees")
	if err != nil {
		t.Fatalf("published course should load from the publish dir alone: %v", err)
	}
	topics := 0
	for _, m := range c.Modules {
		topics += len(m.Topics)
	}
	if topics != 3 {
		t.Fatalf("recipient sees %d topics, want 3", topics)
	}
	// Frontmatter study blocks must survive the copy.
	var bst *CourseTopic
	for i := range c.Modules {
		for j := range c.Modules[i].Topics {
			if c.Modules[i].Topics[j].Title == "BST" {
				bst = &c.Modules[i].Topics[j]
			}
		}
	}
	if bst == nil {
		t.Fatal("BST topic missing from the published course")
	}
	if bst.Kind != "code" || bst.Lang != "python" || len(bst.Tests) != 1 {
		t.Fatalf("BST study block did not survive publishing: %+v", bst)
	}

	// Publishing again refreshes in place rather than failing.
	if _, err := svc.PublishCourse("balanced-trees", dest); err != nil {
		t.Fatalf("re-publish: %v", err)
	}
}

func TestPublishCourseErrors(t *testing.T) {
	svc := newCourseVault(t)
	if _, err := svc.PublishCourse("no-such-course", t.TempDir()); err == nil {
		t.Fatal("publishing an unknown course should fail")
	}
	if _, err := svc.PublishCourse("balanced-trees", "  "); err == nil {
		t.Fatal("publishing with no destination should fail")
	}
}
