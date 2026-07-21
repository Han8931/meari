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

// sseServer fakes an OpenAI-compatible streaming endpoint, capturing the
// request body and replying with the given content split into chunks.
func sseServer(t *testing.T, chunks []string, gotBody *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		*gotBody = string(b)
		w.Header().Set("Content-Type", "text/event-stream")
		for _, c := range chunks {
			payload, _ := json.Marshal(map[string]any{
				"choices": []map[string]any{{"delta": map[string]string{"content": c}}},
			})
			_, _ = w.Write([]byte("data: " + string(payload) + "\n\n"))
		}
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
}

func TestChatStreamDeltasAndContext(t *testing.T) {
	var body string
	srv := sseServer(t, []string{"Hel", "lo ", "world"}, &body)
	defer srv.Close()

	tut := New(config.AIConfig{Provider: "compatible", BaseURL: srv.URL, Model: "m"})
	var deltas []string
	full, err := tut.ChatStream(context.Background(),
		"Current challenge: write SumTo",
		[]ChatTurn{{Role: "user", Content: "help"}},
		func(d string) { deltas = append(deltas, d) })
	if err != nil {
		t.Fatal(err)
	}
	if full != "Hello world" {
		t.Fatalf("full = %q", full)
	}
	if len(deltas) != 3 || deltas[0] != "Hel" {
		t.Fatalf("deltas = %v", deltas)
	}
	// The study context must ride along as a system message.
	if !strings.Contains(body, "write SumTo") {
		t.Fatalf("request body missing the study context:\n%s", body)
	}
	if !strings.Contains(body, `"stream":true`) {
		t.Fatalf("request should ask for streaming:\n%s", body)
	}
}

func TestChatStreamOfflineFallback(t *testing.T) {
	tut := New(config.AIConfig{Provider: "openai"}) // no key -> offline
	var deltas []string
	full, err := tut.ChatStream(context.Background(), "ctx", nil,
		func(d string) { deltas = append(deltas, d) })
	if err != nil {
		t.Fatal(err)
	}
	if full == "" || len(deltas) != 1 || deltas[0] != full {
		t.Fatalf("offline should deliver the canned reply once: %q %v", full, deltas)
	}
}

func TestChatWithoutContextOmitsContextMessage(t *testing.T) {
	var body string
	srv := sseServer(t, []string{"ok"}, &body)
	defer srv.Close()
	tut := New(config.AIConfig{Provider: "compatible", BaseURL: srv.URL, Model: "m"})
	if _, err := tut.ChatStream(context.Background(), "", []ChatTurn{{Role: "user", Content: "hi"}},
		func(string) {}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(body, "Context — what the learner") {
		t.Fatalf("empty context must not inject a context message:\n%s", body)
	}
}

// Regression: ChatStream must send exactly ONE system message, at position 0,
// even with study context — a second system message trips stricter
// OpenAI-compatible servers ("system message must be at the beginning").
func TestChatStreamSingleSystemMessage(t *testing.T) {
	var body string
	srv := sseServer(t, []string{"ok"}, &body)
	defer srv.Close()

	tut := New(config.AIConfig{Provider: "compatible", BaseURL: srv.URL, Model: "m"})
	_, err := tut.ChatStream(context.Background(),
		"STUDY-CONTEXT-MARKER",
		[]ChatTurn{{Role: "user", Content: "help"}, {Role: "assistant", Content: "hi"}, {Role: "user", Content: "more"}},
		func(string) {})
	if err != nil {
		t.Fatal(err)
	}

	var req struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("decode request: %v\n%s", err, body)
	}
	if len(req.Messages) == 0 {
		t.Fatal("no messages sent")
	}
	systems := 0
	for i, m := range req.Messages {
		if m.Role == "system" {
			systems++
			if i != 0 {
				t.Fatalf("system message at index %d, must be first", i)
			}
		}
	}
	if systems != 1 {
		t.Fatalf("sent %d system messages, want exactly 1", systems)
	}
	// The study context must be folded INTO that single system message.
	if !strings.Contains(req.Messages[0].Content, "STUDY-CONTEXT-MARKER") {
		t.Fatalf("study context not folded into the system message:\n%s", req.Messages[0].Content)
	}
	// The last message stays the learner's turn (some servers require this too).
	if last := req.Messages[len(req.Messages)-1]; last.Role != "user" || last.Content != "more" {
		t.Fatalf("last message = %+v, want the final user turn", last)
	}
}
