package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileCreatesAndOverwrites(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")

	if err := WriteFile(p, []byte("first"), 0o644); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(p); string(b) != "first" {
		t.Fatalf("got %q, want first", b)
	}

	// Overwriting replaces the contents atomically (via rename).
	if err := WriteFile(p, []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(p); string(b) != "second" {
		t.Fatalf("got %q, want second", b)
	}

	// No temp files are left behind in the directory.
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("expected only state.json, got %v", names)
	}

	// The permission is applied.
	fi, _ := os.Stat(p)
	if fi.Mode().Perm() != 0o644 {
		t.Fatalf("perm = %v, want 0644", fi.Mode().Perm())
	}
}

// A failed write (unwritable directory) leaves any existing file untouched —
// the whole point of writing to a temp file and renaming.
func TestWriteFilePreservesOnFailure(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "keep.txt")
	if err := WriteFile(p, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Make the directory read-only so CreateTemp/Rename fails.
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Skipf("cannot chmod dir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	if err := WriteFile(p, []byte("new"), 0o644); err == nil {
		t.Skip("directory still writable (running as root?) — skipping")
	}
	_ = os.Chmod(dir, 0o755)
	if b, _ := os.ReadFile(p); string(b) != "original" {
		t.Fatalf("a failed write corrupted the file: got %q, want original", b)
	}
}
