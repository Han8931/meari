package tui

import (
	"context"
	"runtime"
	"testing"
	"time"
)

// startChatStream's worker must exit when its context is cancelled even if the
// reader (the Update loop) has stopped draining the channel — otherwise a
// mode-switch or quit mid-stream leaks the goroutine and its open request.
func TestStreamWorkerExitsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	started := make(chan struct{})
	// fn floods deltas without ever returning on its own; only ctx cancel can
	// stop it. The channel (buffer 64) fills quickly since nobody reads it.
	fn := func(ctx context.Context, onDelta func(string)) (string, error) {
		close(started)
		for {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				onDelta("x") // blocks once the buffer fills (no reader)
			}
		}
	}

	before := runtime.NumGoroutine()
	_, _ = startChatStream(ctx, fn) // ignore the channel — simulate an abandoned reader
	<-started

	cancel() // a quit / mode-switch

	// The worker should unwind promptly; the goroutine count returns to baseline.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before {
			return // worker exited — no leak
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("stream worker did not exit after cancel: goroutines %d > baseline %d",
		runtime.NumGoroutine(), before)
}
