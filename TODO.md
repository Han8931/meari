# Meari — TODO

Working backlog. See the [Roadmap](README.md#roadmap) for the big-picture phases;
this file tracks concrete, actionable items. Check things off as they land.

## Quick wins

- [ ] `:rename <title>` and `:delete` note commands (vault lifecycle, next to `:new`/`:learn`)
- [ ] Word/line count + reading time in the editor status bar
- [ ] `meari check` — suggest the exact fix when key/model/base-url is wrong
- [ ] URL go
- [ ] Mouse visual block
- [x] Global install — config + progress + drafts live in `~/.config/meari` (honors
      `$XDG_CONFIG_HOME`; a checkout with a `config.toml`/`vault/` stays local; `MEARI_HOME`
      overrides), so `meari` behaves the same from any directory


## Study & learning

- [ ] Spaced repetition / flashcards with SM-2 scheduling + `:review` due-queue
- [ ] Quiz mode — multiple-choice generated from a note, AI-graded
- [ ] Cloze deletions — auto fill-in-the-blank cards from a note's key sentences
- [ ] Daily review streak + "due today" counts on the launch screen (uses progress.json)

## Vault & knowledge graph

- [ ] `[[wikilink]]` autocomplete in the editor
- [ ] Tag support (`#tag` / frontmatter tags) + tag browser in the left pane
- [ ] Link graph view (start with an ASCII/adjacency summary)

## AI tutor

- [ ] Lesson regeneration — "go deeper" / "simplify" the current lesson note
- [ ] Citations / source mode — store references in lesson frontmatter

## Platform

- [ ] Index — SQLite-backed search, backlinks, SRS/progress store
- [ ] Desktop packaging (Wails, cgo-free)
- [ ] Vault git auto-commit (`vault.autocommit`) for free history/sync

## Recently done

- [x] Durability + AI-client robustness — atomic writes (write-temp+rename) for
      progress/drafts/notes, surfaced save errors; per-chunk stream idle-timeout, empty-reply
      validation, and bounded retry on transient (network/429/5xx) AI failures
- [x] Course-completion rewards — celebration card + markdown certificate; `:achievements`
      ledger + 🏅 dashboard medals; fixed `:clear progress` not clearing the ledger
- [x] Fixed goroutine leaks — cancel chat-stream and course-build workers on quit/mode-switch
- [x] AI note editing — `:polish`/`:edit` (whole note or Visual selection) → review in chat
      → `:apply`/`:discard`; `:ask`/`:discuss` a selection with the tutor (excerpt pinned to
      every turn so follow-ups stay grounded)
- [x] Vault sidebar root — a fixed `vault` row anchors the tree (no real path shown);
      new notes default to it; `r` reloads the tree from disk
- [x] Chat drag-to-copy — drag the transcript and release to copy (Alt-C too); works on Linux
- [x] Launch dashboard — one full-screen course list (continue / your courses / topic / vault)
      replacing the step-by-step wizard
- [x] Markdown-only courses — the built-in Go track is seeded as ordinary `:course`-format
      markdown; `:publish` shares a course as a self-contained folder for git
- [x] CJK-locale layout fix — pin ambiguous-width glyphs to one cell so the TUI doesn't
      misalign / show `????` under `LANG=*.UTF-8` CJK locales
- [x] `:vault` / `:tutor` — hop between the coding TUI and the notes vault in one
      process (no relaunch); vault gains an Obsidian-style backlinks panel (`:backlinks`)
- [x] Go curriculum: added Constants & iota, Recursion, Generics, Sorting, Panic &
      recover, and Number parsing topics (drawn from gobyexample.com)
- [x] Clickable "▸ Check answer" button in the coding TUI title bar (runs the tests)
- [x] Grey shaded chat input field with a single `>` prompt
- [x] `api_key` documented in `config.toml` (paste-the-key alternative to `api_key_env`)
- [x] Quit from the launch wizard with `esc` / `q`
