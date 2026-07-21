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
//     portable folder keeps everything local (the historical behavior).
//  3. $XDG_CONFIG_HOME/meari (default ~/.config/meari) — the global default,
//     created on first run.
func BaseDir() (string, error) {
	if h := strings.TrimSpace(os.Getenv("MEARI_HOME")); h != "" {
		return ExpandHome(h)
	}
	if wd, err := os.Getwd(); err == nil {
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

func isFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}
