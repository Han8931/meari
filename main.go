// Command meari is an interactive, AI-powered TUI that teaches programming
// by having the learner write code themselves and checking it against tests.
//
// The screen is split into three panes:
//
//	challenges (left)  -> the learner's challenges & drafts, with progress
//	editor (center)    -> the in-app Vim/default code editor
//	chat (right)       -> lesson, test results, and an interactive tutor chat
//
// All AI calls and test runs happen asynchronously so the UI stays responsive.
package main

import (
	"flag"
	"fmt"
	"os"

	"meari/internal/config"
	"meari/internal/drafts"
	"meari/internal/progress"
	"meari/internal/tutor"
	"meari/internal/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		cfgPath   = flag.String("config", "config.toml", "path to config file")
		topicFlag = flag.String("topic", "", "topic to learn (skips the startup prompt)")
		currFlag  = flag.Bool("curriculum", false, "start in guided curriculum mode")
		vimFlag   = flag.Bool("vim", false, "force Vim keybindings in the editor")
		defFlag   = flag.Bool("default", false, "force default (non-Vim) keybindings")
	)
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(*cfgPath, wd)
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
