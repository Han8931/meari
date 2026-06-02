package drafts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClearAllRemovesDraftsOnly(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Save("py-b-vars", "x = 1"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := s.Save("go-b-types", "var x int"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// A non-draft file that must survive the clear.
	keep := filepath.Join(dir, "drafts", "notes.txt")
	if err := os.WriteFile(keep, []byte("keep me"), 0o644); err != nil {
		t.Fatalf("write keep: %v", err)
	}

	if err := s.ClearAll(); err != nil {
		t.Fatalf("ClearAll: %v", err)
	}
	if s.Has("py-b-vars") || s.Has("go-b-types") {
		t.Fatal("ClearAll left drafts behind")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("ClearAll removed a non-draft file: %v", err)
	}
}

func TestClearAllMissingDirIsOK(t *testing.T) {
	s := &Store{dir: filepath.Join(t.TempDir(), "does-not-exist")}
	if err := s.ClearAll(); err != nil {
		t.Fatalf("ClearAll on missing dir: %v", err)
	}
}
