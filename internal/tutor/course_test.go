package tutor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"meari/internal/config"
)

func TestCourseLessonPromptEncouragesUsefulVisuals(t *testing.T) {
	tut, body, close := promptCaptureTutor(t)
	defer close()

	_, err := tut.CourseLesson(context.Background(), "Trees",
		OutlineTopic{Title: "BST search", Summary: "search path and ordering"},
		"source note", []string{"BST"})
	if err != nil {
		t.Fatal(err)
	}
	assertVisualPrompt(t, *body, "CourseLesson")
}

func TestReviseCourseLessonPromptEncouragesUsefulVisuals(t *testing.T) {
	tut, body, close := promptCaptureTutor(t)
	defer close()

	_, err := tut.ReviseCourseLesson(context.Background(), "Trees",
		OutlineTopic{Title: "BST search", Summary: "search path and ordering"},
		"Existing lesson", "source note", []string{"BST"})
	if err != nil {
		t.Fatal(err)
	}
	assertVisualPrompt(t, *body, "ReviseCourseLesson")
	if !strings.Contains(*body, "Existing lesson") {
		t.Fatalf("revision prompt should include the existing lesson:\n%s", *body)
	}
}

func promptCaptureTutor(t *testing.T) (*Tutor, *string, func()) {
	t.Helper()
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{
				"message": map[string]string{
					"role":    "assistant",
					"content": "Lesson body",
				},
			}},
		})
	}))
	tut := New(config.AIConfig{Provider: "compatible", BaseURL: srv.URL, Model: "m"})
	return tut, &body, srv.Close
}

func assertVisualPrompt(t *testing.T, body, label string) {
	t.Helper()
	flat := strings.Join(strings.Fields(strings.ReplaceAll(body, `\n`, " ")), " ")
	for _, want := range []string{"markdown table", "Mermaid diagram", "ASCII diagram", "Do not force a visual"} {
		if !strings.Contains(flat, want) {
			t.Fatalf("%s prompt missing %q:\n%s", label, want, body)
		}
	}
}
