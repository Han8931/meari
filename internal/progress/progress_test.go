package progress

import "testing"

func TestResetClearsHistoryAndPersists(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	s.RecordAttempt("py-b-vars", true)
	s.MarkTopicDone("phys-b-motion")
	s.RecordCompletion(Completion{CourseID: "go-beginner", Title: "Go (Beginner)", Date: "2026-06-01", Topics: 9})
	s.SetLast("python", "beginner", "py-b-vars", "Variables")
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := s.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if len(s.Challenges) != 0 || len(s.Topics) != 0 || len(s.Completions) != 0 || s.Last != nil {
		t.Fatalf("Reset left state behind: %+v", s)
	}

	// Reset must persist: a fresh Load sees the cleared state.
	reloaded, err := Load(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(reloaded.Challenges) != 0 || len(reloaded.Topics) != 0 || len(reloaded.Completions) != 0 || reloaded.Last != nil {
		t.Fatalf("Reset did not persist: %+v", reloaded)
	}
}

func TestCompletionLedger(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	s.RecordCompletion(Completion{CourseID: "go-beginner", Title: "Go (Beginner)", Level: "beginner", Date: "2026-06-01", Topics: 9, FirstTry: 9, Flawless: true})
	s.RecordCompletion(Completion{CourseID: "nc-150", Title: "NeetCode 150", Level: "intermediate", Date: "2026-06-10", Topics: 150, FirstTry: 140})

	// CompletionOf finds a recorded course.
	if c, ok := s.CompletionOf("go-beginner"); !ok || !c.Flawless || c.Topics != 9 {
		t.Fatalf("CompletionOf(go-beginner) = %+v ok=%v", c, ok)
	}
	if _, ok := s.CompletionOf("missing"); ok {
		t.Fatal("CompletionOf should miss an unrecorded course")
	}

	// Re-completing preserves the FIRST date but updates the stats.
	s.RecordCompletion(Completion{CourseID: "nc-150", Title: "NeetCode 150", Date: "2026-07-01", Topics: 150, FirstTry: 150, Flawless: true})
	c, _ := s.CompletionOf("nc-150")
	if c.Date != "2026-06-10" {
		t.Fatalf("re-completion should keep the first date, got %q", c.Date)
	}
	if !c.Flawless || c.FirstTry != 150 {
		t.Fatalf("re-completion should refresh stats, got %+v", c)
	}

	// CompletedCourses comes back newest first.
	list := s.CompletedCourses()
	if len(list) != 2 || list[0].CourseID != "nc-150" || list[1].CourseID != "go-beginner" {
		t.Fatalf("CompletedCourses order = %+v", list)
	}

	// The ledger persists across a reload.
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
	reloaded, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := reloaded.CompletionOf("nc-150"); !ok || c.Date != "2026-06-10" {
		t.Fatalf("ledger did not persist: %+v ok=%v", c, ok)
	}
}
