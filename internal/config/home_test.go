package config

import (
	"os"
	"path/filepath"
	"testing"
)

// $MEARI_HOME wins over everything, with tilde expansion.
func TestBaseDirEnvOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MEARI_HOME", dir)
	got, err := BaseDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != dir {
		t.Fatalf("BaseDir = %q, want %q", got, dir)
	}
}

// A cwd that already looks like a meari home (config.toml or vault/ present)
// keeps state local — the historical, repo/portable behavior.
func TestBaseDirCwdLooksLikeHome(t *testing.T) {
	t.Setenv("MEARI_HOME", "")
	dir := chdir(t, t.TempDir())

	// A bare cwd is NOT a home, so it falls through to the global default.
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	base, err := BaseDir()
	if err != nil {
		t.Fatal(err)
	}
	if base == dir {
		t.Fatal("a bare cwd should not be treated as the home")
	}

	// Drop a config.toml and it becomes the home.
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	base, err = BaseDir()
	if err != nil {
		t.Fatal(err)
	}
	if base != dir {
		t.Fatalf("cwd with config.toml should be the home: got %q, want %q", base, dir)
	}
}

// A vault/ directory also marks the cwd as a home.
func TestBaseDirCwdWithVault(t *testing.T) {
	t.Setenv("MEARI_HOME", "")
	dir := chdir(t, t.TempDir())
	if err := os.Mkdir(filepath.Join(dir, "vault"), 0o755); err != nil {
		t.Fatal(err)
	}
	base, err := BaseDir()
	if err != nil {
		t.Fatal(err)
	}
	if base != dir {
		t.Fatalf("cwd with vault/ should be the home: got %q, want %q", base, dir)
	}
}

// With no override and a plain cwd, BaseDir is <XDG_CONFIG_HOME>/meari, created.
func TestBaseDirGlobalDefault(t *testing.T) {
	t.Setenv("MEARI_HOME", "")
	_ = chdir(t, t.TempDir()) // a blank dir, not a home
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	base, err := BaseDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(xdg, "meari")
	if base != want {
		t.Fatalf("BaseDir = %q, want %q", base, want)
	}
	if !isDir(base) {
		t.Fatal("global home should be created on first run")
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	cases := map[string]string{
		"~":         home,
		"~/notes":   filepath.Join(home, "notes"),
		"/abs/path": "/abs/path",
		"relative":  "relative",
		"~notilde":  "~notilde", // "~" not followed by "/" is left alone
	}
	for in, want := range cases {
		got, err := ExpandHome(in)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("ExpandHome(%q) = %q, want %q", in, got, want)
		}
	}
}

// chdir switches to dir and returns its resolved path (os.Getwd resolves
// symlinks like macOS's /var -> /private/var, which t.TempDir does not).
func chdir(t *testing.T, dir string) string {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	resolved, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}
