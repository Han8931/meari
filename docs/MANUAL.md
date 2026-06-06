# Meari Manual

The complete guide: every command, key, and config option. For the concept and a
quick start, see the [README](../README.md).

- [Running Meari](#running-meari)
- [Configuration](#configuration)
  - [Vault location](#vault-location)
  - [AI providers](#ai-providers-openai-compatible)
  - [Layouts & pane sizes](#layouts--pane-sizes)
- [The vault TUI (`meari -vault`)](#the-vault-tui-meari--vault)
  - [The notes tree](#the-notes-tree)
  - [Commands & study](#commands--study)
  - [Markdown highlighting](#markdown-highlighting)
- [The coding tutor (`meari -tutor`)](#the-coding-tutor-meari--tutor)
- [The chat pane (both TUIs)](#the-chat-pane-both-tuis)
- [The editor (center pane)](#the-editor-center-pane)
- [The command line (`:`)](#the-command-line-)
- [The web UI (`meari serve`)](#the-web-ui-meari-serve)
- [Project layout](#project-layout)
- [Notes & current limits](#notes--current-limits)

## Running Meari

```bash
go build -o meari .
./meari -vault   (or -v)       # the learning vault, in the terminal
./meari serve                  # the learning vault, in your browser
./meari                        # coding tutor — guided setup wizard
./meari -tutor   (or -t)       # coding tutor — skip the wizard into the curriculum
./meari -topic "spanish subjunctive"  # coding tutor — jump to a topic
./meari -vim / -default        # force Vim / non-Vim editor keybindings
./meari check                  # diagnose the AI provider connection
```

The two vault front-ends — terminal (`-vault`) and browser (`serve`) — are driven by
the same core over the same vault directory, so a note created in one shows up in the
other. `:vault` / `:tutor` switch between the two TUIs without quitting; the process
stays up and your session resumes.

## Configuration

All configuration lives in `config.toml` next to meari (or `-config <path>`). The
`:config` command opens it in your `$EDITOR` from inside the app; on save, the layout
re-applies live.

### Vault location

By default the vault lives in `./vault` next to meari. Point it anywhere — e.g. an
existing Obsidian vault (notes with hand-written or unparseable frontmatter still
load; the header just stays in the body):

```toml
[vault]
dir = "~/Documents/my-notes"   # "~/" expands; relative paths are rooted at meari
```

### AI providers (OpenAI-compatible)

Every provider is reached through the OpenAI-compatible chat-completions API, so one
code path works for all — only the base URL / model / key differ.

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

### Layouts & pane sizes

Set `ui.layout` in `config.toml` (or change it live with `:config`):

- **`vertical`** (default) — three side-by-side columns: notes │ editor │ chat.
- **`horizontal`** — notes on the left, with the **content on top and your input on the
  bottom**. Better for reading- and writing-heavy subjects.

Set your **default pane split** with `sidebar_percent` / `chat_percent` under `[ui]`
(percent of the width; the editor takes the rest — e.g. `chat_percent = 45` for a
chat-heavy layout). `:compact` / `:wide` still adjust live from that base, and
`sidebar_folded = true` starts with the left pane folded away (`:fold` toggles it
live, in both TUIs).

Transient command feedback (copy confirmations, resize/fold notices, unknown
commands…) appears briefly in the **bottom status bar**, keeping the chat transcript
for actual conversation. The editor pins the **challenge statement above the code** as
a labeled description block, wrapped to the pane width — so the problem stays readable
however long the chat gets, without polluting your buffer (essay study pins the prompt
as a `>` header the same way).

## The vault TUI (`meari -vault`)

The 3-pane terminal vault: **notes** (left, a file tree mirroring your real folder
structure) │ **editor** (center) │ **chat / study** (right).

### The notes tree

The tree shows your vault as it is on disk, NERDTree/Obsidian-style — directories
first (`▸` folded / `▾` unfolded), files indented beneath them, the open note in bold:

- `j`/`k` move · `Enter` opens a note / folds-unfolds a directory
- `Space` marks files/folders (amber) for a batch operation
- `m` opens the node menu: **(a)dd** — type a path, end with `/` for a folder;
  **(m)ove/rename** — edit the prefilled path; **(d)elete** — the marked rows (or the
  cursor row) after a `y/n` confirm. Deletes clear the editor if the open note went
  with them; renames follow the open note.

### Commands & study

- `Ctrl-W` cycle focus · `Ctrl-S` save · `Ctrl-C` quit
- `:learn <topic>` — generate an AI lesson note (e.g. `:learn the french revolution`)
- `:new <title>` — create a blank note
- `:essay` — study the open note: write an answer in the editor, then `:grade` to check
  it; `:answer` reveals a model answer; `:done` ends the study
- `:backlinks` — toggle the "↩ Linked mentions" panel under the editor, listing the
  notes whose `[[wikilinks]]` point at the open note (Obsidian-style backlinks)
- `:tutor` — hand off to the coding tutor without quitting (the coding TUI's `:vault`
  comes back); the process stays up and your curriculum session resumes

### Markdown highlighting

The editor highlights `# headings`, ` ```fenced code``` ` (with real Go/Python
highlighting when the fence names the language), `` `inline code` ``, `[[wikilinks]]`,
`> blockquotes` (spanning wrapped lines), `-`/`*`/`+` list markers, and `*italic*` /
`**bold**` / `***both***`. Highlighting is color-only — it never shifts a character —
and is stable while the cursor moves or the view scrolls.

## The coding tutor (`meari -tutor`)

On a bare `meari` launch you're walked through a short **setup wizard** (`↑`/`↓` or
`j`/`k` to move, `Enter` to choose, `Esc` to go back). It drops you into one of:

- **Curriculum mode** — a built-in, ordered, pre-authored learning path (no AI needed,
  works offline), with progress saved so you can **Continue where you left off**. The
  Go track, for example, goes deep — imperative basics & the type system → functions,
  methods, closures, slices/maps → structs, JSON, interfaces, and pointers.
- **Custom mode** — any topic you type; the AI generates the material for it.

The `-tutor` (`-t`) and `-topic` flags skip the wizard for returning users.

**Global keys** (any pane):

- `Ctrl-W` then `h` / `l` — move focus left / right, Vim window-style
  (`Tab` / `Shift-Tab` also cycle focus)
- `Ctrl-R` — check / submit the current item
- `Ctrl-N` — advance to the next item on the current topic
- `Ctrl-C` — quit (your work and progress are saved)

In the **left pane**, `j`/`k` move and `Enter` opens an item. In the **chat** pane,
type a question and press `Enter` to ask the tutor.

**Global commands** — type `:` in the left pane (or use the editor's `:` line):

- `:topic <subject>` / `:subject <subject>` — switch subject; no argument opens a
  picker. Keeps your current level.
- `:vault` — switch to the notes vault (the `meari -vault` UI) without quitting; `:tutor`
  from there switches back, and your coding session resumes where you left off.
- `:progress` — progress summary (completion bars + activity).
- `:clear` — clear the chat transcript. `:clear progress` / `:clear drafts` wipe saved
  history / drafts (each confirms first).

## The chat pane (both TUIs)

- Speaker **badges** (` you ` / ` tutor ` / ` lesson ` on colored backgrounds) make turns
  easy to tell apart; fenced ``` code blocks in tutor/lesson messages are
  **syntax-highlighted** behind a gutter bar.
- Everything wraps to the pane — long words, URLs, and code lines included (code
  hard-wraps under its gutter rather than being cut off). Need more room? `:compact`
  repeatedly grows the chat pane up to ~60% of the width (`:wide` gives it back to the
  editor) — in both TUIs.
- An animated **progress line** ("⠹ tutor thinking…") shows inside the pane while the
  AI works. The input area sits in a **shaded grey field** under a `>` prompt and is
  **three rows tall** so longer questions wrap visibly.
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
- **Paste a question:** `Option-V` (macOS) / `Alt-V` (Linux) inserts the system
  clipboard into the chat input; `:paste` does the same and focuses the input
  (`Ctrl-V` while typing, or the terminal's own `Cmd-V` paste, work too).

**Scrolling the chat** (lessons and replies get long):

- **Left click** — focuses the pane under the cursor.
- **Mouse wheel** — scrolls whatever pane is under the cursor, without changing focus
  (like `ranger`/`lf`).
- With the chat focused: `Ctrl-F`/`Ctrl-B` page, `Ctrl-D`/`Ctrl-U` half-page,
  `Shift-↑`/`Shift-↓` by line, plus `PgUp`/`PgDn`. New messages only jump you to the
  bottom when you were already there, so reading back isn't interrupted.

## The editor (center pane)

A modal, Vim-style editor (configurable). Set `editor.keybindings` to `"vim"` or
`"default"` in `config.toml`. The current mode is unmistakable: a **green `NORMAL`** /
**blue `INSERT`** badge in the status line and a steady, color-coded cursor.

**Vim mode — Normal**
- Move: `h j k l` · `w` next word · `b` previous word · `e` end-of-word ·
  `0`/`^` line start · `$` line end · `gg`/`G` top/bottom of file ·
  `{`/`}` (Shift-[ / Shift-]) previous/next paragraph · `Ctrl-E`/`Ctrl-Y` scroll by line
- **Jumplist:** `G`, `{`/`}`, and searches record where you came from — `Ctrl-O` walks
  back, `Ctrl-I` (`Tab`) forward, Vim-style; the view always follows the cursor
- **Counts:** a numeric prefix repeats motions and edits — `3w`, `5x`, `2dd`, `3yy`, `2>>`, `2J`
- **Char find:** `f`/`F` to a character (forward/back), `t`/`T` till before it; `;`/`,` repeat
- **Search:** `/pattern` then Enter; `n`/`N` next/previous match (wraps)
- `J` joins lines · `~` toggles case
- Enter Insert: `i` `a` · `I`/`A` (line start/end) · `o`/`O` (open line below/above,
  **at the current line's indentation**)
- Edit: `x` · `r<char>` · `dd` · `dw` · `D` · `cc`/`cw`/`C` · `<<`/`>>` dedent/indent line
- Register: deletes and `yy` (yank line) fill the unnamed register; `p`/`P` paste
  after/before (falls back to the system clipboard when the register is empty).
  **Yanks (`yy`, visual `y`) also copy to the system clipboard** (native + OSC 52, so
  it works over SSH) — paste them into any other app. Deletes stay register-only.
- **Paste from other apps:** `Alt-V`/`Option-V` inserts the system clipboard at the
  cursor in any mode; the terminal's `Cmd-V` paste lands literally too (it can never
  fire as Vim commands)
- **Undo/redo:** `u` undo · `Ctrl-R` redo (an Insert session is one undo unit; the
  restored change is centered on screen; in the coding TUI, run tests with
  `Ctrl-S`/`:submit` while the editor is focused)
- **Visual mode:** `v` charwise · `V` linewise — motions extend the highlighted
  selection; `d`/`x` delete · `y` yank · `c` change · `<`/`>` indent · `o` swap ends ·
  `Esc` cancels
- `Esc` returns to Normal (and cancels a half-typed operator like `d`)
- **Insert mode:** `Tab` indents (4 spaces); **Enter auto-indents** the new line — one
  level deeper after `{`, `(`, `[` or `:` — and typing `}` on a blank-indented line
  **dedents it electrically**; `o`/`O` follow the same rules

`Ctrl-S` submit / `Ctrl-Q` quit work in any mode.

## The command line (`:`)

- `:submit` — check the current item (same as `Ctrl-R`)
- `:w` — save a draft and keep editing (resume later)
- `:q` — leave the app (`:wq` saves + submits)
- `:config` — open `config.toml` in your `$EDITOR`; on save, the layout re-applies live

The `:` command line (and the editor's `:` / `/` prompts) recall **previous commands
with ↑/↓**, with separate histories for commands and searches — and **Tab-complete
command names** (Shift-Tab cycles backward), showing a Vim-wildmenu-style candidate
list in the status bar: `:co⇥` → `[compact]  config  copy  course`.

> `:w` (save & resume) is intentionally separate from `:submit` (check), so you can stop
> mid-answer and come back to it.

## The web UI (`meari serve`)

```bash
./meari serve                  # http://localhost:8765
./meari serve --addr :9000     # custom port
```

A 3-pane browser app over your vault: **notes** (left) with a "Generate lesson" box, a
**markdown editor + live preview** (center) with `[[wikilink]]` navigation and backlinks,
and a **chat / study** panel (right) with tutor chat and an Essay study mode — write an
answer and **Check answer** grades it; **Show answer** reveals a model answer. Runs offline
with built-in content; configure an AI provider for generated lessons and grading.

## Project layout

```
main.go                 entry point: load config, construct deps, launch a front-end
internal/
  core/                 headless engine: vault+tutor orchestration both front-ends
                        drive (list/open/save/generate/backlinks/chat/essay)
  vault/                markdown vault: notes + frontmatter + [[wikilinks]] + file ops
  web/                  local web GUI (net/http) + `meari serve`, over core
  tutor/                OpenAI-compatible client; lesson/challenge/feedback/chat,
                        plus subject-agnostic GenerateNote + GradeEssay, + offline
  tui/                  the 3-pane Bubble Tea program (panes, async cmds, layout)
  config/               TOML config + defaults + flag overrides
  curriculum/           built-in ordered learning paths (modules + topics)
  editor/               embeddable Bubble Tea modal Vim editor + highlighters
  executor/             runs code against tests (timeout-guarded)
  drafts/               save/load/clear in-progress work by id
  progress/             progress.json — attempts + topic status

planned (see the README roadmap):
  index/                derived SQLite index: search, link graph, SRS, progress
  study/  srs/          study-mode graders (Quiz/Flashcard/Essay/Code) + scheduling
```

## Notes & current limits

- The code-execution path runs Python via `python3` with a timeout. It is **not a
  sandbox** — fine for a single trusted local learner; don't run untrusted code.
- Vim mode is "Vim-*style*" (motions, counts, `d`/`c`/`y` operators, visual mode,
  registers, undo/redo, the jumplist, the command line), not full Vim — no marks,
  macros, or `:s` substitution yet.
- `Ctrl-W` is reserved for pane navigation. In the editor, `Tab` indents in Insert
  mode and walks the jumplist forward in Normal mode.
- `Ctrl-E`/`Ctrl-Y` move the cursor line-by-line (the view scrolls with it at the
  window edges) — the textarea's viewport can't scroll independently of the cursor.
