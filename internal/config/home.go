package config

// home.go resolves the app home: the one directory that roots every piece of
// per-user state — config.toml, data/ (progress), workspace/ (drafts),
// exports/, and generated courses — so an installed meari behaves the same no
// matter which directory you launch it from.

import (
	"os"
	"path/filepath"
	"strings"
)

// BaseDir picks the directory that roots all per-user state. Order:
//
//  1. $MEARI_HOME, if set (an explicit override; "~/" expands).
//  2. The current directory IF it already looks like a meari home/checkout —
//     a config.toml or a vault/ is present — so running from the repo or a
//     portable folder keeps everything local (the historical behavior). The
//     home directory itself is deliberately excluded: a stray config.toml or an
//     unrelated vault/ (e.g. an Obsidian vault) in $HOME must never turn $HOME
//     into the root and scatter data/ workspace/ meari-course/ … across it.
//     Set $MEARI_HOME to opt into a home-rooted layout explicitly.
//  3. $XDG_CONFIG_HOME/meari (default ~/.config/meari) — the global default,
//     created on first run.
func BaseDir() (string, error) {
	if h := strings.TrimSpace(os.Getenv("MEARI_HOME")); h != "" {
		return ExpandHome(h)
	}
	if wd, err := os.Getwd(); err == nil && !isHomeDir(wd) {
		if isFile(filepath.Join(wd, "config.toml")) || isDir(filepath.Join(wd, "vault")) {
			return wd, nil
		}
	}
	cfgRoot := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME"))
	if !filepath.IsAbs(cfgRoot) { // unset, or a relative value the XDG spec says to ignore
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cfgRoot = filepath.Join(home, ".config")
	}
	base := filepath.Join(cfgRoot, "meari")
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	return base, nil
}

// ExpandHome resolves a leading "~"/"~/" to the user's home directory. "~" on
// its own and "~/sub" both work; other values pass through unchanged. Shared
// with the vault-dir resolver so "~/" means the same everywhere.
func ExpandHome(p string) (string, error) {
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, strings.TrimPrefix(p[1:], "/")), nil
	}
	return p, nil
}

// StrayHomeFiles returns the meari-owned names sitting directly in $HOME that
// the app is NOT using — a footprint left by running an older build from the
// home directory before BaseDir excluded it. It's empty when base already IS
// $HOME, or when nothing meari-shaped is there. Callers should treat a lone
// generic hit (e.g. an unrelated ~/data) with a grain of salt; a ~/config.toml
// or two+ matches is a strong signal.
func StrayHomeFiles(base string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	if filepath.Clean(base) == filepath.Clean(home) {
		return nil
	}
	var found []string
	for _, name := range []string{"config.toml", "vault", "data", "workspace", "meari-course", "meari-publish", "exports"} {
		if _, err := os.Stat(filepath.Join(home, name)); err == nil {
			found = append(found, name)
		}
	}
	return found
}

// isHomeDir reports whether p is the user's home directory (so the local-root
// detection can skip it). A failure to resolve $HOME is treated as "not home",
// leaving the normal detection in place.
func isHomeDir(p string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	ap, err1 := filepath.Abs(p)
	ah, err2 := filepath.Abs(home)
	if err1 != nil || err2 != nil {
		return false
	}
	return filepath.Clean(ap) == filepath.Clean(ah)
}

func isFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}
