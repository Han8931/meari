# Meari desktop app

A native desktop front-end for Meari, built with [Wails v2](https://wails.io)
(Go + WebKit) and a React/TypeScript UI. It drives the **same** `core.Service`
and markdown vault as the terminal app and the web server — no feature logic
lives here, only DTO shaping and stream bookkeeping (`app.go`).

## The editor

The center pane is a [CodeMirror 6](https://codemirror.net) markdown editor with
Vim keybindings via [`@replit/codemirror-vim`](https://github.com/replit/codemirror-vim):
full Normal / Visual / Insert modes, counts, registers, and text objects. Crucially,
**mouse selection integrates with Vim** — drag with the mouse to enter Visual mode
over the dragged span, then `y`/`d`/`c` act on it, exactly like Vim in a terminal.
`:w` / `:wq` save through the host. Toggle Vim off with the checkbox in the title bar
for a plain editor.

## Prerequisites

- Go 1.26+
- Node 18+ and npm
- The Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0`
- Platform toolchain: Xcode CLT (macOS) or WebKitGTK (Linux) — run `wails doctor`.

## Develop

```bash
cd gui
wails dev          # live-reloading app: Go rebuilds + Vite HMR
```

## Build

```bash
cd gui
wails build        # -> gui/build/bin/Meari.app (macOS) / Meari (Linux)
```

## Test

```bash
go test ./gui/            # Go bindings, offline (no Wails runtime needed)
cd gui/frontend && npm run build   # typechecks (tsc) + bundles (vite)
```

## Layout

```
gui/
  main.go            Wails bootstrap; builds the shared core.Service
  app.go             the bound App: methods callable from TypeScript
  app_test.go        binding tests against a temp vault, offline tutor
  wails.json         Wails project config
  frontend/
    src/
      App.tsx        3-pane shell (tree | editor/preview | chat)
      Editor.tsx     CodeMirror + Vim editor (mouse-aware)
      Chat.tsx       streaming tutor conversation
      Tree.tsx       vault file tree
      api.ts         typed bridge to the Go bindings + stream helper
      style.css      theme-aware styling (light/dark)
```

The embedded `frontend/dist/` holds a placeholder `index.html` in git so
`go build ./gui` always compiles; `wails build` regenerates the real bundle.
