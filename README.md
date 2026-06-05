# Meari

A **self-directed learning vault** — Obsidian-style markdown notes you own, with an
AI tutor that turns *"I want to learn X"* into a saved lesson note and then helps you
**practice and remember it** through active-recall study modes. Any subject:
languages, math, history, science — not just programming.

Meari runs as a **terminal app (TUI)** and a **local web GUI** (`meari serve`); both
are thin front-ends over the same vault and tutor, converging on one shared Go core.

> **Status: in active evolution.** Meari began as an AI *coding* tutor (write code →
> run hidden tests → get feedback) and that loop still works end-to-end today. It's
> being generalized into the subject-agnostic learning vault described above. See
> [Roadmap](#roadmap) for what's landed and what's next.

## The idea

- **You own your notes.** Every note is a plain `.md` file with YAML frontmatter in a
  local *vault* — readable, editable, and syncable with anything (git, Obsidian, …).
  Notes link to each other with `[[wikilinks]]`; backlinks and search come from a
  derived index, never from a lock-in database.
- **AI lessons become notes.** Ask to learn a topic and the tutor writes a focused
  lesson as a new note in your vault (with `[[links]]` to prerequisites) — not a
  throwaway chat message.
- **Study modes reinforce any note.** Layer active recall on top of your notes:
  **Essay** (free-text answer, AI-graded), **Quiz**, **Flashcards** with spaced
  repetition, and **Code** (write code checked against tests) as one optional mode.
- **Offline-friendly.** Runs with built-in content when no AI provider is configured;
  configure one for generated lessons, study items, and feedback.

## Run

```bash
go build -o meari .
./meari notes                  # the learning vault, in the terminal
./meari serve                  # the learning vault, in your browser
./meari                        # coding tutor — guided setup wizard
./meari -curriculum            # coding tutor — skip the wizard, built-in path
./meari -topic "spanish subjunctive"  # coding tutor — jump to a topic
./meari -vim / -default        # force Vim / non-Vim editor keybindings
```

The learning vault has two interchangeable front-ends — a **terminal app**
(`meari notes`) and a **browser app** (`meari serve`) — both driven by the same core
over the same `./vault`, so a note created in one shows up in the other.

### Vault in the terminal — `meari notes`

The 3-pane terminal vault: **notes** (left, grouped by subject) │ **editor** (center) │
**chat / study** (right).

- `Ctrl-W` cycle focus · `j`/`k` then `Enter` to open a note · `Ctrl-S` save · `Ctrl-C` quit
- `:learn <topic>` — generate an AI lesson note (e.g. `:learn the french revolution`)
- `:new <title>` — create a blank note
- `:essay` — study the open note: write an answer in the editor, then `:grade` to check
  it; `:answer` reveals a model answer; `:done` ends the study

### Vault in the browser — `meari serve`

```bash
./meari serve                  # http://localhost:8765
./meari serve --addr :9000     # custom port
```

A 3-pane browser app over your vault: **notes** (left) with a "Generate lesson" box, a
**markdown editor + live preview** (center) with `[[wikilink]]` navigation and backlinks,
and a **chat / study** panel (right) with tutor chat and an Essay study mode — write an
answer and **Check answer** grades it; **Show answer** reveals a model answer. Runs offline
with built-in content; configure an AI provider for generated lessons and grading.

### The three panes

```
┌ notes ──────┐┌──────── editor ─────────┐┌──── chat / study ────┐
│ ▸ math/     ││ # Derivatives           ││ lesson  …             │
│   limits    ││ A derivative measures…  ││ tutor   …             │
│ ▸ spanish/  ││ [[Limits]] first.       ││ › ask the tutor…      │
└─────────────┘└─────────────────────────┘└───────────────────────┘
```

- **notes** (left) — your vault / learning path, with progress (`✓` done / `…` started).
- **editor** (center) — the in-app Vim/default editor for the current note or answer.
- **chat / study** (right) — the lesson, study results, tutor feedback, **and** an
  interactive chat where you can ask the tutor questions.

All AI calls and checks run asynchronously, so the UI never freezes.

### Getting started

On launch you're walked through a short **setup wizard** (`↑`/`↓` or `j`/`k` to move,
`Enter` to choose, `Esc` to go back). It drops you into one of:

- **Curriculum mode** — a built-in, ordered, pre-authored learning path (no AI needed,
  works offline), with progress saved so you can **Continue where you left off**. The
  Go track, for example, goes deep — imperative basics & the type system → functions,
  methods, closures, slices/maps → structs, JSON, interfaces, and pointers.
- **Custom mode** — any topic you type; the AI generates the material for it.

The `-curriculum` and `-topic` flags skip the wizard for returning users.

**Global keys** (any pane):

- `Ctrl-W` then `h` / `l` — move focus left / right, Vim window-style
  (`Tab` / `Shift-Tab` also cycle focus)
- `Ctrl-R` — check / submit the current item
- `Ctrl-N` — advance to the next item on the current topic
- `Ctrl-C` — quit (your work and progress are saved)

In the **left pane**, `j`/`k` move and `Enter` opens an item. In the **chat** pane,
type a question and press `Enter` to ask the tutor.

**The chat pane** (both TUIs):

- Speaker **badges** (` you ` / ` tutor ` / ` lesson ` on colored backgrounds) make turns
  easy to tell apart; fenced ``` code blocks in tutor/lesson messages are
  **syntax-highlighted** behind a gutter bar.
- Everything wraps to the pane — long words, URLs, and code lines included (code
  hard-wraps under its gutter rather than being cut off). Need more room? `:compact`
  repeatedly grows the chat pane up to ~60% of the width (`:wide` gives it back to the
  editor) — in both TUIs.
- An animated **progress line** ("⠹ tutor thinking…") shows inside the pane while the
  AI works, and the input area is **three rows tall** so longer questions wrap visibly.
- The transcript is **per-topic**: switching topics/notes gives you a clean pane for the
  new one, and returning to a previous topic restores its chat and study history.
- Replies **stream in live**, and every question carries the **current context** — the
  lesson, the challenge, and your in-progress code (or the open note and your essay
  draft) — so answers relate to what's on screen. Long conversations send only the most
  recent turns to the model.
- **↑/↓ recall your previous questions** (when the input is empty), readline-style.
- **Copy a reply:** with the chat focused, `Option-O` (macOS) / `Alt-O` (Linux) copies the
  tutor's last reply to the clipboard; `:copy code` grabs just its last code block and
  `:copy all` the whole transcript. Copying uses the native clipboard *and* OSC 52, so it
  also works over SSH in supporting terminals. (On macOS, if your terminal sends Option
  as Meta/Esc+, both modes work; `Cmd-O` can't reach a terminal app.)
- **Paste a question:** `:paste` inserts the system clipboard into the chat input and
  focuses it (`Ctrl-V` while typing in the input works too).

**Scrolling the chat** (lessons and replies get long):

- **Left click** — focuses the pane under the cursor.
- **Mouse wheel** — scrolls whatever pane is under the cursor, without changing focus
  (like `ranger`/`lf`).
- With the chat focused: `Ctrl-F`/`Ctrl-B` page, `Ctrl-D`/`Ctrl-U` half-page,
  `Shift-↑`/`Shift-↓` by line, plus `PgUp`/`PgDn`. New messages only jump you to the
  bottom when you were already there, so reading back isn't interrupted.

**Global commands** — type `:` in the left pane (or use the editor's `:` line):

- `:topic <subject>` / `:subject <subject>` — switch subject; no argument opens a
  picker. Keeps your current level.
- `:progress` — progress summary (completion bars + activity).
- `:clear` — clear the chat transcript. `:clear progress` / `:clear drafts` wipe saved
  history / drafts (each confirms first).

## AI providers (OpenAI-compatible)

Every provider is reached through the OpenAI-compatible chat-completions API, so one
code path works for all — only the base URL / model / key differ. Configure in
`config.toml`:

**OpenAI**
```toml
[ai]
provider = "openai"
model = "gpt-4o-mini"
api_key_env = "OPENAI_API_KEY"   # the NAME of the env var holding your key…
# api_key = "sk-..."             # …or paste the key itself (env var wins)
```
```bash
export OPENAI_API_KEY=sk-...     # in the same shell you run meari from
```

**Ollama (local, no key needed)**
```toml
[ai]
provider = "ollama"
model = "llama3.1"
# base_url defaults to http://localhost:11434/v1
```

**Any compatible gateway**
```toml
[ai]
provider = "compatible"
base_url = "https://your-gateway/v1"
model = "your-model"
api_key_env = "YOUR_KEY_ENV"   # optional — no-auth local servers work without a key
# timeout_seconds = 120        # raise for big/slow local models
```

**Diagnose your connection** with `meari check` — it prints the resolved
provider/base URL/model/key status, verifies the model exists upstream, and sends
a real test request:

```
$ meari check
Meari AI connection check
  provider:  ollama
  base url:  http://localhost:11434/v1
  model:     qwen3-coder-next:latest
  api key:   not set (looked in $OPENAI_API_KEY)

✓ provider reachable; model "qwen3-coder-next:latest" is available (7 models total)
✓ chat round-trip OK in 252ms
```

## In-app editor (center pane)

A modal, Vim-style editor (configurable). Set `editor.keybindings` to `"vim"` or
`"default"` in `config.toml`. The current mode is unmistakable: a **green `NORMAL`** /
**blue `INSERT`** badge in the status line and a steady, color-coded cursor.

**Vim mode — Normal**
- Move: `h j k l` · `w` next word · `b` previous word · `e` end-of-word ·
  `0`/`^` line start · `$` line end · `gg`/`G` top/bottom of file
- **Counts:** a numeric prefix repeats motions and edits — `3w`, `5x`, `2dd`, `3yy`, `2>>`, `2J`
- **Char find:** `f`/`F` to a character (forward/back), `t`/`T` till before it; `;`/`,` repeat
- **Search:** `/pattern` then Enter; `n`/`N` next/previous match (wraps)
- `J` joins lines · `~` toggles case
- Enter Insert: `i` `a` · `I`/`A` (line start/end) · `o`/`O` (open line below/above)
- Edit: `x` · `r<char>` · `dd` · `dw` · `D` · `cc`/`cw`/`C` · `<<`/`>>` dedent/indent line
- Register: deletes and `yy` (yank line) fill the unnamed register; `p`/`P` paste
  after/before (falls back to the system clipboard when the register is empty)
- **Undo/redo:** `u` undo · `Ctrl-R` redo (an Insert session is one undo unit; in the
  coding TUI, run tests with `Ctrl-S`/`:submit` while the editor is focused)
- `o`/`O` open a line below/above **at the current line's indentation**
- **Visual mode:** `v` charwise · `V` linewise — motions extend the highlighted
  selection; `d`/`x` delete · `y` yank · `c` change · `<`/`>` indent · `o` swap ends ·
  `Esc` cancels
- `Esc` returns to Normal (and cancels a half-typed operator like `d`)
- **Insert mode:** `Tab` indents (4 spaces)

**Command line (`:`)**
- `:submit` — check the current item (same as `Ctrl-R`)
- `:w` — save a draft and keep editing (resume later)
- `:q` — leave the app (`:wq` saves + submits)
- `:config` — open `config.toml` in your `$EDITOR`; on save, the layout re-applies live

`Ctrl-S` submit / `Ctrl-Q` quit work in any mode.

## Layouts

Set `ui.layout` in `config.toml` (or change it live with `:config`):

- **`vertical`** (default) — three side-by-side columns: notes │ editor │ chat.
- **`horizontal`** — notes on the left, with the **content on top and your input on the
  bottom**. Better for reading- and writing-heavy subjects.

> `:w` (save & resume) is intentionally separate from `:submit` (check), so you can stop
> mid-answer and come back to it.

## Project layout

```
main.go                 entry point: load config, construct deps, launch a front-end
internal/
  core/                 ✅ headless engine: vault+tutor orchestration both         (NEW)
                        front-ends drive (list/open/save/generate/backlinks/chat/essay)
  vault/                ✅ markdown vault: notes + frontmatter + [[wikilinks]]      (NEW)
  web/                  ✅ local web GUI (net/http) + `meari serve`, over core      (NEW)
  tutor/                OpenAI-compatible client; lesson/challenge/feedback/chat,
                        plus subject-agnostic GenerateNote + GradeEssay, + offline
  tui/                  the 3-pane Bubble Tea program (panes, async cmds, layout)
  config/               TOML config + defaults + flag overrides
  curriculum/           built-in ordered learning paths (modules + topics)
  editor/               embeddable Bubble Tea modal Vim editor
  executor/             runs code against tests (timeout-guarded)
  drafts/               save/load/clear in-progress work by id
  progress/             progress.json — attempts + topic status

planned (see Roadmap):
  index/                derived SQLite index: search, link graph, SRS, progress
  study/  srs/          study-mode graders (Quiz/Flashcard/Essay/Code) + scheduling
```

## Roadmap

The pivot from coding tutor → general learning vault, in phases. **Files are the source
of truth; the index is rebuildable. New dependencies are kept pure-Go / cgo-free** so a
desktop build (Wails/Tauri) stays a thin later step.

- **✅ Vault** — own your notes as markdown with frontmatter and `[[wikilinks]]`.
- **✅ Subject-agnostic tutor** — `GenerateNote` (lesson → note) and `GradeEssay` for any subject.
- **✅ Web GUI (first cut)** — `meari serve`: 3-pane browser UI with lesson generation,
  markdown editor + preview, wikilink navigation, backlinks, chat, and Essay study.
- **✅ Core engine** — a headless `core.Service` owns the vault+tutor orchestration and
  returns plain data; both front-ends drive it (no business logic in handlers).
- **✅ TUI vault parity** — `meari notes` browses/edits vault notes, generates lessons,
  and runs Essay study in the terminal, on the same core as the web GUI.
- **▢ Index** — SQLite-backed search, backlinks, and the SRS/progress store (replaces the
  current in-memory backlink scan).
- **▢ Spaced repetition** — flashcards with SM-2 scheduling and a review queue.
- **▢ More study modes** — Quiz, and Code-with-tests restored as one optional mode.
- **▢ Knowledge graph** — backlink panels and a visual link graph.
- **▢ Desktop app** — package the core + web UI as a native window.

## Notes / current limits

- The code-execution path runs Python via `python3` with a timeout. It is **not a
  sandbox** — fine for a single trusted local learner; don't run untrusted code.
- Vim mode is "Vim-*style*" (motions, `d`/`c`/`y` operators, visual mode, registers,
  `r`, the command line), not full Vim — no counts, marks, or undo yet.
- `Tab` and `Ctrl-W` are reserved for pane navigation, so they can't be typed into the
  editor.

## Todo

- Vim keybinding erros (w / p ...)
- LLM connection issue for compatible
- Answer check button
