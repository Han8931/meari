# Meari — TODO

Working backlog. See the [Roadmap](README.md#roadmap) for the big-picture phases;
this file tracks concrete, actionable items. Check things off as they land.

## Quick wins

- [ ] Reading progress (80%, 90%...)
- [x] Lecture update feature — `:capture` / `:capture all` save chat Q&A to a linked
      companion note (`My Notes/<Lecture>.md`), leaving the lecture untouched;
      `:weave [instruction]` reorganizes that note into coherent subject-grouped
      prose as a reviewable proposal (`:apply` / `:discard`).
- [ ] `:rename <title>` and `:delete` note commands (vault lifecycle, next to `:new`/`:learn`)
- [ ] Word/line count + reading time in the editor status bar
- [ ] `meari check` — suggest the exact fix when key/model/base-url is wrong
- [ ] URL go
- [ ] Mouse visual block
- [x] Global install: config/`workspace`/`data`/`exports`/courses root at a global home
      (`~/.config/meari`, or `$MEARI_HOME`), so `meari` behaves the same from any directory;
      a cwd holding `config.toml`/`vault/` stays local. `:config` edits the global file.


## Study & learning

- [x] Resume fixing the Rust Intermediate course: review the current files under
      `meari-course/Rust (Intermediate)`, finish the remaining revisions, compile every
      reference solution against its tests, and verify the final lesson order. Do not
      modify unrelated changes.
- [ ] Offline course checkpoints — make conceptual answers checkable without an LLM:
      keep code exercises compiler/test-backed; add deterministic multiple-choice,
      classify/order, and “does this compile / what happens?” questions after each
      module; use `rustc` as the oracle for Rust compile/output checks (match stable
      outcomes or error codes, not full diagnostics); give every distractor targeted
      feedback; and make genuinely free-form explanations self-checked with a model
      answer, required-idea rubric, common misconceptions, reveal/retry controls, and
      learner confidence rather than brittle keyword grading. Start with separate
      `kind: quiz` checkpoint topics using the existing schema; consider multiple study
      activities per lesson later so code, quiz, and reflection can coexist.
- [ ] Spaced repetition / flashcards with SM-2 scheduling + `:review` due-queue
- [ ] Quiz mode — multiple-choice generated from a note, AI-graded
- [ ] Cloze deletions — auto fill-in-the-blank cards from a note's key sentences
- [ ] AI Q&A cards — generate question/answer pairs from a note (same pipeline as essay/challenge)
- [ ] Guided course tutoring — `:study` a course topic-by-topic (teach → quiz → next) instead of
      only essay/challenge study, so a built course closes the recall loop
- [ ] Daily review streak + "due today" counts on the launch screen (uses progress.json)

## Vault & knowledge graph

- [ ] Follow links — jump to a linked note from the backlinks ("↩ Linked mentions") panel
      (focus + `enter` to open) and from inline `[[wikilinks]]` under the cursor in the editor
      (e.g. `gf`/`⌃]`); currently both are display-only, not navigable
- [ ] `[[wikilink]]` autocomplete in the editor
- [ ] Full-text search across the vault (in-memory inverted index now; SQLite-backed later)
- [ ] Tag support (`#tag` / frontmatter tags) + tag browser in the left pane
- [ ] Link graph view (start with an ASCII/adjacency summary)
- [ ] Note templates / daily notes — `:today` and `:template` (trivial over `vSaveOpenCmd`)

## Desktop app (`gui/`)

The Wails desktop app exposes a subset of the TUI — bring the AI/vault features across:

- [ ] Wire `Explain selection` into the chat pane (the Go stream exists; the UI logs the id)
- [ ] AI note editing — `:polish`/`:edit`/`:ask` on a selection
- [ ] `:course` / `:revise` / `:publish` from the app
- [ ] Fuzzy find and a backlinks panel UI (`Backlinks` binding already exists)
- [ ] Stamp the build with `main.version` from a build script, like the CLI's `meari version`

## AI tutor

- [ ] Lesson regeneration — "go deeper" / "simplify" the current lesson note
- [ ] Citations / source mode — store references in lesson frontmatter

## Platform

- [ ] Index — SQLite-backed search, backlinks, SRS/progress store
- [x] Desktop app (Wails) — `gui/`, native window over the shared core with a Vim editor
- [ ] Vault git auto-commit (`vault.autocommit`) for free history/sync

## Correctness & safety

- [x] **Path traversal** — `vault.Read`/`Write` now route through `safeAbs()`; covered by
      `internal/vault/traversal_test.go` and verified against the live `/api/note` endpoint
- [x] Atomic writes (temp-file + rename) for notes, `progress.json`, drafts, and the chat
      store — `internal/fsutil.WriteFileAtomic`

## Recently done

- [x] More Vim in the editor — text objects (`iw`/`aw`, quotes, brackets) with `d`/`c`/`y`
      (`ciw`, `di"`, `da{`…), `y` with motions (`yw`/`y$`/`yiw`), `X`/`s`/`S`, backward
      search `?`, word search `*`/`#`, direction-aware `n`/`N`; lecture reader gained
      counts (`5j`) and `f`/`F`/`t`/`T` + `;`/`,`. (Still missing: `.` repeat, marks,
      macros, blockwise Visual, `:s` — see the audit.)
- [x] Launch screen — daily learning-quote epigraph (page-sage style, `internal/quotes`)
      + a "Just chat with meari" blank-tutor option (`dashScratch`)
- [x] Lecture UX — Check-answer button hidden while reading a lecture (quiz rows only);
      mouse-wheel scrolling carries the reader cursor so keys don't snap the view back
- [x] Tutor lecture pane is a Vim reader/editor — keyboard cursor + Visual selection
      (`hjkl`/`w`/`b`/`gg`/`G`, `v`+`y`/`⌥c` copy), `/` search with `n`/`N`, `gd` follows the
      `[[wikilink]]` under the cursor to that lecture, `⌃o`/`⌃i` jump back/forward across
      lectures + positions, `,ff` fuzzy-jumps to any lecture, and `i`/`:edit` edits the
      lecture source (file-backed courses only; `⌃s`/`:w` saves, preserving frontmatter)
- [x] Global home hardening — `meari` never treats `$HOME` as its root, so a stray
      `config.toml`/`vault/` in the home dir can't scatter `data/ workspace/ meari-course/ …`
      there (`config.BaseDir` excludes `$HOME`; `$MEARI_HOME` still opts in explicitly)
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
