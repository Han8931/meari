// Command meari is an interactive, AI-powered self-learning vault. It runs as a
// terminal app (the default) or a local web app (the "serve" subcommand); both
// drive the same vault and tutor.
//
// The terminal UI splits the screen into three panes:
//
//	notes (left)    -> the learner's vault / learning path, with progress
//	editor (center) -> the in-app Vim/default editor
//	chat (right)    -> lessons, study results, and an interactive tutor chat
//
// All AI calls and checks happen asynchronously so the UI stays responsive.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/drafts"
	"meari/internal/progress"
	"meari/internal/tutor"
	"meari/internal/tui"
	"meari/internal/vault"
	"meari/internal/web"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	// Subcommand dispatch: "serve" launches the web UI, "notes" the vault TUI;
	// anything else is the classic coding TUI. We peel the subcommand off os.Args
	// before flag parsing so each mode owns its own flag set.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve":
			return runServe(os.Args[2:])
		case "notes":
			return runNotes(os.Args[2:])
		case "check":
			return runCheck(os.Args[2:])
		}
	}
	return runTUI()
}

// loadConfig loads config rooted at the working directory.
func loadConfig(cfgPath string) (config.Config, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return config.Config{}, "", err
	}
	cfg, err := config.Load(cfgPath, wd)
	return cfg, wd, err
}

func runTUI() error {
	var (
		cfgPath   = flag.String("config", "config.toml", "path to config file")
		topicFlag = flag.String("topic", "", "topic to learn (skips the startup prompt)")
		currFlag  = flag.Bool("curriculum", false, "start in guided curriculum mode")
		vimFlag   = flag.Bool("vim", false, "force Vim keybindings in the editor")
		defFlag   = flag.Bool("default", false, "force default (non-Vim) keybindings")
	)
	flag.Parse()

	cfg, wd, err := loadConfig(*cfgPath)
	if err != nil {
		return err
	}
	// CLI flags override the config file (most-specific wins).
	if *vimFlag {
		cfg.Editor.Keybindings = "vim"
	}
	if *defFlag {
		cfg.Editor.Keybindings = "default"
	}

	store, err := drafts.New(cfg.WorkspaceDir)
	if err != nil {
		return err
	}
	prog, err := progress.Load(cfg.DataDir)
	if err != nil {
		return err
	}
	tut := tutor.New(cfg.AI)

	return tui.Run(tui.Deps{
		Tutor:      tut,
		Store:      store,
		Progress:   prog,
		Cfg:        cfg,
		Topic:      *topicFlag,
		Curriculum: *currFlag,
		ConfigPath: *cfgPath,
		BaseDir:    wd,
	})
}

// runServe starts the local web UI over the same vault and tutor as the TUI.
func runServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	cfgPath := fs.String("config", "config.toml", "path to config file")
	addr := fs.String("addr", ":8765", "address to listen on")
	_ = fs.Parse(args)

	cfg, _, err := loadConfig(*cfgPath)
	if err != nil {
		return err
	}

	v, err := vault.Open(cfg.VaultDir)
	if err != nil {
		return err
	}
	svc := core.New(v, tutor.New(cfg.AI))

	fmt.Printf("Meari web UI on http://localhost%s  (vault: %s)\n", *addr, cfg.VaultDir)
	if svc.Offline() {
		fmt.Println("(offline — no AI provider configured; set OPENAI_API_KEY or use Ollama for generated lessons)")
	}
	return web.Serve(*addr, svc)
}

// runCheck diagnoses the AI provider connection: resolved settings, whether the
// configured model exists upstream, and a real round-trip request.
func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	cfgPath := fs.String("config", "config.toml", "path to config file")
	_ = fs.Parse(args)

	cfg, _, err := loadConfig(*cfgPath)
	if err != nil {
		return err
	}
	tut := tutor.New(cfg.AI)
	info := tut.Info()

	fmt.Println("Meari AI connection check")
	fmt.Printf("  provider:  %s\n", cfg.AI.Provider)
	fmt.Printf("  base url:  %s\n", info.BaseURL)
	fmt.Printf("  model:     %s\n", info.Model)
	keyState := "not set"
	switch {
	case cfg.AI.APIKeyEnv != "" && os.Getenv(cfg.AI.APIKeyEnv) != "":
		keyState = "set (from $" + cfg.AI.APIKeyEnv + ")"
	case info.KeySet:
		keyState = "set (from api_key in config.toml)"
	case cfg.AI.APIKeyEnv != "":
		keyState = "not set (looked in $" + cfg.AI.APIKeyEnv + " and config api_key)" +
			"\n             note: api_key_env is the NAME of an env var, not the key itself"
	}
	fmt.Printf("  api key:   %s\n", keyState)

	if info.Offline {
		fmt.Println("\n✗ OFFLINE: this endpoint requires an API key and none is set.")
		fmt.Println("  Either `export " + cfg.AI.APIKeyEnv + "=sk-...` in the shell you run meari from,")
		fmt.Println("  or put `api_key = \"sk-...\"` under [ai] in config.toml,")
		fmt.Println("  or point [ai] at a local provider (Ollama).")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Does the configured model exist upstream?
	models, err := tut.Models(ctx)
	switch {
	case err != nil:
		fmt.Printf("\n✗ could not reach the provider: %v\n", err)
		fmt.Println("  Is the server running and the base url correct?")
		return nil
	default:
		found := false
		for _, id := range models {
			if id == info.Model {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("\n✓ provider reachable; model %q is available (%d models total)\n", info.Model, len(models))
		} else {
			fmt.Printf("\n⚠ provider reachable, but model %q was NOT in its model list (%d models).\n", info.Model, len(models))
			max := len(models)
			if max > 8 {
				max = 8
			}
			fmt.Printf("  available include: %s\n", strings.Join(models[:max], ", "))
		}
	}

	// Real round trip through the same code path lessons/chat use.
	fmt.Println("  sending a test request…")
	dur, err := tut.Ping(ctx)
	if err != nil {
		fmt.Printf("✗ chat request failed: %v\n", err)
		return nil
	}
	fmt.Printf("✓ chat round-trip OK in %s\n", dur.Round(time.Millisecond))
	return nil
}

// runNotes starts the vault terminal UI over the same vault and tutor as the web
// UI, sharing the core engine.
func runNotes(args []string) error {
	fs := flag.NewFlagSet("notes", flag.ExitOnError)
	cfgPath := fs.String("config", "config.toml", "path to config file")
	vimFlag := fs.Bool("vim", false, "force Vim keybindings in the editor")
	defFlag := fs.Bool("default", false, "force default (non-Vim) keybindings")
	_ = fs.Parse(args)

	cfg, _, err := loadConfig(*cfgPath)
	if err != nil {
		return err
	}
	if *vimFlag {
		cfg.Editor.Keybindings = "vim"
	}
	if *defFlag {
		cfg.Editor.Keybindings = "default"
	}

	v, err := vault.Open(cfg.VaultDir)
	if err != nil {
		return err
	}
	svc := core.New(v, tutor.New(cfg.AI))
	return tui.RunVault(svc, cfg.VimEditor())
}
