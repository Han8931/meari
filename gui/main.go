package main

// Meari — the desktop app: the same core engine as the terminal app in a
// native window (Wails v2; WebKit on macOS, WebKitGTK on Linux). Build with
// `wails build` from this directory. The TUI and web front-ends are untouched;
// all three drive one core.Service over the same markdown vault.

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"meari/internal/config"
	"meari/internal/core"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// version is stamped by the build (-ldflags -X main.version=…); "dev" for a
// plain `wails dev` / `go build`.
var version = "dev"

//go:embed all:frontend/dist
var assets embed.FS

// buildService mirrors main.buildDeps: config rooted at the working directory,
// a vault-backed core.Service with the courses mount and built-in courses
// seeded. Returned separately so tests can construct their own service.
func buildService() (*core.Service, config.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, config.Config{}, err
	}
	cfg, err := config.Load("config.toml", wd)
	if err != nil {
		return nil, config.Config{}, err
	}
	v, err := vault.Open(cfg.VaultDir)
	if err != nil {
		return nil, config.Config{}, err
	}
	svc := core.New(v, tutor.New(cfg.AI))
	svc.SetCourseDir(cfg.CourseDir)
	if err := svc.SeedBuiltinCourses(); err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not seed built-in courses:", err)
	}
	return svc, cfg, nil
}

func main() {
	svc, cfg, err := buildService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "meari:", err)
		os.Exit(1)
	}
	a := newApp(svc, cfg)

	if err := wails.Run(&options.App{
		Title:  "Meari",
		Width:  1120,
		Height: 760,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: a.startup,
		Bind:      []interface{}{a},
	}); err != nil {
		fmt.Fprintln(os.Stderr, "meari:", err)
		os.Exit(1)
	}
}
