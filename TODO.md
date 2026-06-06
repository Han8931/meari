# Meari — TODO

Working backlog. See the [Roadmap](README.md#roadmap) for the big-picture phases;
this file tracks concrete, actionable items. Check things off as they land.

## Quick wins

- [ ] Fuzzy note finder — `Ctrl-P` / `:find` palette that jumps to a note by title
- [ ] `:rename <title>` and `:delete` note commands (vault lifecycle, next to `:new`/`:learn`)
- [ ] Word/line count + reading time in the editor status bar
- [x] Answer check button — clickable "▸ Check answer" in the coding TUI title bar
- [ ] `meari check` — suggest the exact fix when key/model/base-url is wrong

## Study & learning

- [ ] Spaced repetition / flashcards with SM-2 scheduling + `:review` due-queue
- [ ] Quiz mode — multiple-choice generated from a note, AI-graded
- [ ] Cloze deletions — auto fill-in-the-blank cards from a note's key sentences
- [ ] Daily review streak + "due today" counts on the launch screen (uses progress.json)

## Vault & knowledge graph

- [ ] `[[wikilink]]` autocomplete in the editor
- [x] Backlinks panel (vault `:backlinks` — "↩ Linked mentions" under the editor)
- [ ] Tag support (`#tag` / frontmatter tags) + tag browser in the left pane
- [ ] Link graph view (start with an ASCII/adjacency summary)

## AI tutor

- [ ] "Explain this selection" — explain/expand selected editor text inline
- [ ] Lesson regeneration — "go deeper" / "simplify" the current lesson note
- [ ] Citations / source mode — store references in lesson frontmatter

## Platform

- [ ] Index — SQLite-backed search, backlinks, SRS/progress store
- [ ] Desktop packaging (Wails, cgo-free)
- [ ] Vault git auto-commit (`vault.autocommit`) for free history/sync

## Recently done

- [x] `:vault` / `:tutor` — hop between the coding TUI and the notes vault in one
      process (no relaunch); vault gains an Obsidian-style backlinks panel (`:backlinks`)
- [x] Go curriculum: added Constants & iota, Recursion, Generics, Sorting, Panic &
      recover, and Number parsing topics (drawn from gobyexample.com)
- [x] Clickable "▸ Check answer" button in the coding TUI title bar (runs the tests)
- [x] Grey shaded chat input field with a single `>` prompt
- [x] `api_key` documented in `config.toml` (paste-the-key alternative to `api_key_env`)
- [x] Quit from the launch wizard with `esc` / `q`
