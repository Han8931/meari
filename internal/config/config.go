// Package config loads and merges user configuration with sensible defaults.
//
// Resolution order (most to least specific): CLI flags -> config file ->
// built-in defaults. CLI flag overrides are applied by the caller via the
// setter helpers after Load returns.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the fully-resolved configuration handed to the rest of the app.
type Config struct {
	AI         AIConfig         `toml:"ai"`
	Editor     EditorConfig     `toml:"editor"`
	Navigation NavigationConfig `toml:"navigation"`
	UI         UIConfig         `toml:"ui"`

	// Paths are derived at load time, not read from the file.
	WorkspaceDir string `toml:"-"`
	DataDir      string `toml:"-"`
	VaultDir     string `toml:"-"` // the learner's markdown note vault
}

// AIConfig selects the model backend. All providers are reached through the
// OpenAI-compatible chat-completions API, so OpenAI, Ollama, and any compatible
// gateway share one code path — only the base URL / model / key differ.
type AIConfig struct {
	// Provider is "openai", "ollama", or "compatible". It only chooses the
	// default BaseURL; an explicit BaseURL always wins.
	Provider string `toml:"provider"`
	// BaseURL overrides the provider default, e.g. http://localhost:11434/v1.
	BaseURL string `toml:"base_url"`
	// Model is the model name, e.g. "gpt-4o-mini" or "llama3.1".
	Model string `toml:"model"`
	// APIKeyEnv names the environment variable holding the API key. Ollama
	// needs no key, so this may be empty for local use.
	APIKeyEnv string `toml:"api_key_env"`
	// TimeoutSeconds bounds each model request. 0 means the default (120s —
	// local models can be slow to load and generate long lessons).
	TimeoutSeconds int `toml:"timeout_seconds"`
}

type EditorConfig struct {
	// Keybindings selects the in-app editor style: "vim" or "default".
	Keybindings string `toml:"keybindings"`
}

type NavigationConfig struct {
	// Keybindings selects menu navigation: "vim" (j/k/q) or "default".
	Keybindings string `toml:"keybindings"`
}

type UIConfig struct {
	Theme string `toml:"theme"`
	// Layout arranges the panes: "vertical" (side-by-side columns: list | editor |
	// chat — good for coding) or "horizontal" (content on top, input on the
	// bottom — good for reading/writing subjects).
	Layout string `toml:"layout"`
}

// Default returns the built-in configuration used when no file is present.
func Default() Config {
	return Config{
		AI: AIConfig{
			Provider:  "openai",
			Model:     "gpt-4o-mini",
			APIKeyEnv: "OPENAI_API_KEY",
		},
		Editor:     EditorConfig{Keybindings: "vim"},
		Navigation: NavigationConfig{Keybindings: "vim"},
		UI:         UIConfig{Theme: "default", Layout: "vertical"},
	}
}

// ResolveBaseURL returns the effective API base URL: an explicit BaseURL if set,
// otherwise the per-provider default.
func (a AIConfig) ResolveBaseURL() string {
	if a.BaseURL != "" {
		return a.BaseURL
	}
	switch a.Provider {
	case "ollama":
		return "http://localhost:11434/v1"
	default: // "openai" and generic "compatible"
		return "https://api.openai.com/v1"
	}
}

// Load reads the config file at path (if it exists), overlaying it on the
// defaults. A missing file is not an error: defaults are used. Derived paths
// are rooted at baseDir.
func Load(path, baseDir string) (Config, error) {
	cfg := Default()

	if path != "" {
		if _, err := os.Stat(path); err == nil {
			if _, err := toml.DecodeFile(path, &cfg); err != nil {
				return cfg, err
			}
		}
	}

	cfg.WorkspaceDir = filepath.Join(baseDir, "workspace")
	cfg.DataDir = filepath.Join(baseDir, "data")
	cfg.VaultDir = filepath.Join(baseDir, "vault")

	// Normalize unknown values back to safe defaults.
	if cfg.Editor.Keybindings != "vim" && cfg.Editor.Keybindings != "default" {
		cfg.Editor.Keybindings = "vim"
	}
	if cfg.Navigation.Keybindings != "vim" && cfg.Navigation.Keybindings != "default" {
		cfg.Navigation.Keybindings = "vim"
	}
	if cfg.UI.Layout != "vertical" && cfg.UI.Layout != "horizontal" {
		cfg.UI.Layout = "vertical"
	}

	return cfg, nil
}

// VimEditor reports whether the in-app editor should use Vim bindings.
func (c Config) VimEditor() bool { return c.Editor.Keybindings == "vim" }

// Horizontal reports whether the stacked (content-over-input) layout is selected.
func (c Config) Horizontal() bool { return c.UI.Layout == "horizontal" }

// defaultConfigTOML is the template written by EnsureFile when no config exists,
// so the :config command always has a well-commented file to open.
const defaultConfigTOML = `# Meari configuration. All fields are optional; defaults are shown.

[ai]
# provider: "openai" | "ollama" | "compatible"
provider = "openai"
# base_url = "http://localhost:11434/v1"   # uncomment for Ollama
model = "gpt-4o-mini"
api_key_env = "OPENAI_API_KEY"

[editor]
# keybindings: "vim" | "default"
keybindings = "vim"

[navigation]
keybindings = "vim"

[ui]
theme = "default"
# layout: "vertical" (list | editor | chat) or "horizontal" (content on top,
# input on the bottom — better for reading/writing subjects)
layout = "vertical"
`

// EnsureFile writes the default config template to path if it does not yet
// exist, so the editor opened by :config always has content to edit.
func EnsureFile(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}
	return os.WriteFile(path, []byte(defaultConfigTOML), 0o644)
}
