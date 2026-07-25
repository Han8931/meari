package core

import (
	"os/exec"
	"path/filepath"
	"testing"

	"meari/internal/config"
	"meari/internal/executor"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// The published Rust course is hand-authored rather than generated from the
// in-code Go curriculum. Keep its quiz-like study challenges executable: every
// lesson must provide a complete code exercise, and every reference answer must
// compile and satisfy its hidden tests.
func TestPublishedRustBeginnerStudyChallenges(t *testing.T) {
	if _, err := exec.LookPath("rustc"); err != nil {
		t.Skip("rustc not installed")
	}
	if err := exec.Command("rustc", "--version").Run(); err != nil {
		t.Skip("rustc has no usable toolchain")
	}

	v, err := vault.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(v, tutor.New(config.AIConfig{Provider: "openai"}))
	svc.SetCourseDir(filepath.Join("..", "..", "meari-publish"))

	course, err := svc.LoadCourse("rust-beginner")
	if err != nil {
		t.Fatalf("load published Rust beginner course: %v", err)
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
	if topics != 21 {
		t.Fatalf("loaded %d topics, want 21", topics)
	}
}
