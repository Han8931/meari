package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppBaseDirMeariHomeWins(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MEARI_HOME", dir)
	got, err := appBaseDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != dir {
		t.Fatalf("MEARI_HOME base = %q, want %q", got, dir)
	}
}

// A directory that already looks like a meari home (has config.toml or vault/)
// keeps everything local — the portable / repo workflow.
func TestAppBaseDirPortable(t *testing.T) {
	t.Setenv("MEARI_HOME", "")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte("[ai]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	got, err := appBaseDir()
	if err != nil {
		t.Fatal(err)
	}
	wd, _ := os.Getwd()
	if got != wd {
		t.Fatalf("portable base = %q, want cwd %q", got, wd)
	}
}

// With no override and a plain working directory, state lands in
// ~/.config/meari (the XDG config home).
func TestAppBaseDirGlobalDefault(t *testing.T) {
	t.Setenv("MEARI_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	t.Chdir(t.TempDir()) // a plain dir with no config.toml / vault

	got, err := appBaseDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(fakeHome, ".config", "meari")
	if got != want {
		t.Fatalf("global base = %q, want %q", got, want)
	}
	if !isDir(want) {
		t.Fatal("the global app home should be created on first run")
	}
}

// $XDG_CONFIG_HOME relocates the global home.
func TestAppBaseDirXDGConfigHome(t *testing.T) {
	t.Setenv("MEARI_HOME", "")
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Chdir(t.TempDir())

	got, err := appBaseDir()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(xdg, "meari"); got != want {
		t.Fatalf("XDG base = %q, want %q", got, want)
	}
}

func TestExpandTilde(t *testing.T) {
	home, _ := os.UserHomeDir()
	for in, want := range map[string]string{
		"~/notes": filepath.Join(home, "notes"),
		"~":       home, // bare "~" now expands too (shared with config.ExpandHome)
	} {
		got, err := expandTilde(in)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("expandTilde(%q) = %q, want %q", in, got, want)
		}
	}
}
