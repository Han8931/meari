// Package fsutil holds small filesystem helpers shared across the app.
package fsutil

import (
	"os"
	"path/filepath"
)

// WriteFile atomically replaces the file at path with data. It writes to a
// temporary file in the SAME directory, flushes it to disk, then renames it
// over the target — and rename is atomic on POSIX, so a crash or a full disk
// mid-write leaves the previous file intact rather than a truncated or empty
// one. This is the durable replacement for os.WriteFile when the file holds
// state we can't afford to corrupt (progress, drafts, notes).
func WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	// Remove the temp file on every error path; after a successful Rename it no
	// longer exists, so this no-ops.
	defer func() { _ = os.Remove(tmp) }()

	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil { // get the bytes to disk before the rename
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
