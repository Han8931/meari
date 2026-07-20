package vault

import (
	"os"
	"path/filepath"
	"testing"
)

// Read and Write must refuse paths that escape the vault root — the web
// server passes ?path= straight through, so this is the security boundary.
func TestReadWriteRejectEscapingPaths(t *testing.T) {
	dir := t.TempDir()
	v, err := Open(filepath.Join(dir, "vault"))
	if err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(secret, []byte("s3cret"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, p := range []string{
		"../secret.txt",
		"../../etc/passwd",
		"a/../../secret.txt",
		"/etc/passwd",
		"..",
	} {
		if _, err := v.Read(p); err == nil {
			t.Errorf("Read(%q) escaped the vault", p)
		}
		if _, err := v.Write(Note{RelPath: p, Title: "x", Body: "x"}); err == nil {
			t.Errorf("Write(%q) escaped the vault", p)
		}
	}

	// An empty RelPath is legal for Write (it derives a path from the title)
	// but must not read anything.
	if _, err := v.Read(""); err == nil {
		t.Error(`Read("") should fail`)
	}

	// Normal notes keep working, including in subfolders.
	if _, err := v.Write(Note{RelPath: "math/limits.md", Title: "Limits", Body: "hi"}); err != nil {
		t.Fatalf("legit Write failed: %v", err)
	}
	if _, err := v.Read("math/limits.md"); err != nil {
		t.Fatalf("legit Read failed: %v", err)
	}
}
