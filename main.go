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
	"flag"
	"fmt"
	"os"

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
	// Subcommand dispatch: "meari serve" launches the web UI; anything else is
	// the TUI. We peel the subcommand off os.Args before flag parsing so each
	// mode owns its own flag set.
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		return runServe(os.Args[2:])
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
