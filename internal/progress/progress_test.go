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
	s.SetLast("python", "beginner", "py-b-vars", "Variables")
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := s.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if len(s.Challenges) != 0 || len(s.Topics) != 0 || s.Last != nil {
		t.Fatalf("Reset left state behind: %+v", s)
	}

	// Reset must persist: a fresh Load sees the cleared state.
	reloaded, err := Load(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(reloaded.Challenges) != 0 || len(reloaded.Topics) != 0 || reloaded.Last != nil {
		t.Fatalf("Reset did not persist: %+v", reloaded)
	}
}
