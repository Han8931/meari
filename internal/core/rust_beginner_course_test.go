package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"meari/internal/config"
	"meari/internal/executor"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// The Rust courses are hand-authored rather than generated from the in-code Go
// curriculum. Keep their quiz-like study challenges executable: every lesson
// must provide a complete code exercise, and every reference answer must
// compile and satisfy its hidden tests.
func TestRustCourseStudyChallenges(t *testing.T) {
	if _, err := exec.LookPath("rustc"); err != nil {
		t.Skip("rustc not installed")
	}
	if err := exec.Command("rustc", "--version").Run(); err != nil {
		t.Skip("rustc has no usable toolchain")
	}

	courses := []struct {
		name, dir, folder, id string
		topics                int
	}{
		{"beginner", "meari-publish", "Rust (Beginner)", "rust-beginner", 21},
		{"intermediate", "meari-course", "Rust (Intermediate)", "rust-intermediate", 18},
	}
	for _, tc := range courses {
		t.Run(tc.name, func(t *testing.T) {
			root := filepath.Join("..", "..", tc.dir)
			manifest := filepath.Join(root, tc.folder, "course.md")
			if _, err := os.Stat(manifest); os.IsNotExist(err) {
				t.Skipf("local %s course is not present", tc.name)
			} else if err != nil {
				t.Fatal(err)
			}

			v, err := vault.Open(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			svc := New(v, tutor.New(config.AIConfig{Provider: "openai"}))
			svc.SetCourseDir(root)

			course, err := svc.LoadCourse(tc.id)
			if err != nil {
				t.Fatalf("load Rust %s course: %v", tc.name, err)
			}

			var topics int
			for _, module := range course.Modules {
				for _, topic := range module.Topics {
					topics++
					t.Run(topic.Title, func(t *testing.T) {
						if topic.Kind != "code" || topic.Lang != "rust" {
							t.Fatalf("study kind/lang = %q/%q, want code/rust", topic.Kind, topic.Lang)
						}
						if topic.Prompt == "" || topic.Starter == "" ||
							topic.Answer == "" || len(topic.Tests) == 0 {
							t.Fatalf("incomplete study challenge: %+v", topic)
						}

						result, err := executor.Run(topic.Lang, topic.Answer, topic.Tests)
						if err != nil {
							t.Fatal(err)
						}
						if !result.Passed {
							t.Fatalf("reference answer failed:\n%s", result.Output)
						}
					})
				}
			}
			if topics != tc.topics {
				t.Fatalf("loaded %d topics, want %d", topics, tc.topics)
			}
		})
	}
}
