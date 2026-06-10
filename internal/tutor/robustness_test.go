package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"meari/internal/config"
)

func compatTutor(url string) *Tutor {
	return New(config.AIConfig{Provider: "compatible", BaseURL: url, Model: "m"})
}

func jsonReply(content string) string {
	b, _ := json.Marshal(map[string]any{
		"choices": []map[string]any{{"message": map[string]string{"content": content}}},
	})
	return string(b)
}

// A transient 503 then 500 then a 200 should succeed via retry.
func TestChatRetriesTransientFailures(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch atomic.AddInt32(&calls, 1) {
		case 1:
			w.WriteHeader(http.StatusServiceUnavailable)
		case 2:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			fmt.Fprint(w, jsonReply("recovered"))
		}
	}))
	defer srv.Close()

	got, err := compatTutor(srv.URL).chat(context.Background(), "sys", "hi")
	if err != nil {
		t.Fatal(err)
	}
	if got != "recovered" {
		t.Fatalf("content = %q, want recovered", got)
	}
	if n := atomic.LoadInt32(&calls); n != 3 {
		t.Fatalf("made %d calls, want 3 (two retries)", n)
	}
}

// A client error (401) must NOT be retried.
func TestChatDoesNotRetryClientError(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	if _, err := compatTutor(srv.URL).chat(context.Background(), "sys", "hi"); err == nil {
		t.Fatal("a 401 should be an error")
	}
	if n := atomic.LoadInt32(&calls); n != 1 {
		t.Fatalf("made %d calls, want 1 (no retry on a client error)", n)
	}
}

// An empty reply is reported clearly, not passed through as "".
func TestChatRejectsEmptyReply(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonReply("   ")) // whitespace-only content
	}))
	defer srv.Close()

	_, err := compatTutor(srv.URL).chat(context.Background(), "sys", "hi")
	if err == nil || !strings.Contains(err.Error(), "empty reply") {
		t.Fatalf("want an 'empty reply' error, got %v", err)
	}
}

// A stream that goes quiet (the server holds the connection open but never
// sends data) is aborted by the IDLE WATCHDOG on its own — no caller
// cancellation — rather than hanging forever.
func TestChatStreamAbortsOnStall(t *testing.T) {
	prev := streamIdleTimeout
	streamIdleTimeout = 200 * time.Millisecond
	t.Cleanup(func() { streamIdleTimeout = prev })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		<-r.Context().Done() // never send any data; wait until the client gives up
	}))
	defer srv.Close()

	tut := compatTutor(srv.URL)

	done := make(chan error, 1)
	go func() {
		// A plain background ctx: only the watchdog can end this.
		_, err := tut.chatStreamRaw(context.Background(),
			[]chatMessage{{Role: "user", Content: "hi"}}, func(string) {})
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "stalled") {
			t.Fatalf("want a 'stalled' error from the watchdog, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("stalled stream hung — the idle watchdog did not fire")
	}
}
