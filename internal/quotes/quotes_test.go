package quotes

import (
	"testing"
	"time"
)

// Daily is deterministic for a date and moves across days.
func TestDailyDeterministic(t *testing.T) {
	d1 := time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC)
	d1later := time.Date(2026, 7, 25, 23, 59, 0, 0, time.UTC)
	if Daily(d1) != Daily(d1later) {
		t.Fatal("the same date must yield the same quote")
	}
	if Daily(d1).Text == "" {
		t.Fatal("Daily should return a non-empty quote")
	}
	// Over a couple of weeks we should see more than one distinct quote.
	seen := map[string]bool{}
	for i := 0; i < 14; i++ {
		seen[Daily(d1.AddDate(0, 0, i)).Text] = true
	}
	if len(seen) < 2 {
		t.Fatalf("expected the quote to rotate across days, saw %d distinct", len(seen))
	}
}

// Every built-in quote has text and an author (no half-filled entries).
func TestBuiltinWellFormed(t *testing.T) {
	for i, q := range builtin {
		if q.Text == "" || q.Author == "" {
			t.Fatalf("quote %d is missing text or author: %+v", i, q)
		}
	}
}
