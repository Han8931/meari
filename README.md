# Meari

An interactive, AI-powered **TUI** that teaches programming by having you **write the
code yourself** and checking it against real tests.

The screen is split into three panes:

```
┌ challenges ┐┌──────── editor ────────┐┌──── chat ────┐
│ ✓ list-sum ││ def sum_list(xs):       ││ lesson  …     │
│ … reverse  ││     pass                ││ tutor   …     │
└────────────┘└─────────────────────────┘└── › ask… ────┘
```

- **challenges** (left) — your challenges & draft files, with progress
  (`✓` solved / `…` in-progress).
- **editor** (center) — the in-app Vim/default code editor.
- **chat** (right) — the lesson, test results, tutor feedback, **and** an
  interactive chat where you can ask the tutor questions.

The loop: **pick a topic → lesson + challenge → write code → run the tests →
tutor feedback → next challenge.** All AI calls and test runs happen
asynchronously, so the UI never freezes.

## Run

```bash
go build -o meari .
./meari                      # guided setup wizard (recommended)
./meari -curriculum          # skip the wizard, start the curriculum
./meari -topic "python loops"  # skip the wizard, jump to a topic
./meari -default             # force non-Vim editor keybindings
./meari -vim                 # force Vim editor keybindings
```

### Getting started

On launch you're walked through a short **setup wizard** (use `↑`/`↓` or `j`/`k`,
`Enter` to choose, `Esc` to go back):

1. **Course** — Python, Go, or Physics.
2. **How to learn** — *start from the beginning* (the full curriculum) or *a
   specific topic* (which you then type).
3. **Level** — beginner / intermediate / advanced, which tunes how the tutor
   pitches lessons and challenges.

This gives you one of two modes:

- **Curriculum mode** — a built-in, ordered learning path for **Python**, **Go**,
  or **Physics**, at beginner / intermediate / advanced. Every lesson and challenge
  is pre-authored (no AI needed) and verified, so the path is consistent and works
  offline. The Go track goes deep — imperative basics & the type system
  (floats, the integer types & wraparound, `math/big`, runes/UTF-8, conversions)
  → functions, methods, closures, arrays/slices/`append`/maps → structs, JSON,
  constructors, composition, interfaces, and pointers. The **left pane** shows
  the modules and topics with progress (`✓` done / `…` started); `Enter` on a
  topic loads its lesson and challenge. Your spot is saved — relaunch and pick
  **Continue where you left off**.
- **Custom mode** — any single topic you type. The left pane lists the
  challenges you generate for it (AI-generated; Python).

The `-curriculum` and `-topic` flags skip the wizard for returning users.

**Pane keys** (global, in any pane):

- `Tab` / `Shift-Tab` — cycle focus between panes
- `Ctrl-W` then `h` / `l` — move focus left / right, Vim window-style
- `Ctrl-R` — run the tests for the current challenge
- `Ctrl-N` — generate the next challenge for the current topic
- `Ctrl-C` — quit (your draft and progress are saved)

In the **left pane**, `j`/`k` move and `Enter` opens the topic/challenge. In the
**chat** pane, type a question and press `Enter` to ask the tutor.

**Scrolling the chat** (lessons, history, and tutor replies can get long):

- **Mouse wheel** — scrolls whatever pane is under the cursor (chat history, or
  the left-pane selection), without changing focus — like `ranger`/`lf`.
- With the chat focused: `Ctrl-F`/`Ctrl-B` page down/up, `Ctrl-D`/`Ctrl-U` half
  page, `Shift-↑`/`Shift-↓` a line at a time, plus `PgUp`/`PgDn`. New messages
  only jump you to the bottom when you were already there, so reading back
  through history isn't interrupted.

**Global commands** — type `:` in the left pane (or use the editor's `:` line) to
open a command prompt:

- `:topic <course>` / `:subject <course>` — switch course (e.g. `:topic go`,
  `:subject physics`); with no argument it opens a course picker. Keeps your
  current level.
- `:progress` — show a progress summary (per-course completion bars + test runs).
- `:clear` — clear the chat transcript. `:clear progress` and `:clear drafts`
  wipe your saved learning history / draft code (each asks to confirm first).

Runs **offline** out of the box with built-in content. Configure an AI provider
for generated lessons, challenges, and feedback.

## AI providers (OpenAI-compatible)

Every provider is reached through the OpenAI-compatible chat-completions API, so
the same code path works for all of them — only the base URL / model / key differ.
Configure in `config.toml`:

**OpenAI**
```toml
[ai]
provider = "openai"
model = "gpt-4o-mini"
api_key_env = "OPENAI_API_KEY"
```
```bash
export OPENAI_API_KEY=sk-...
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
api_key_env = "YOUR_KEY_ENV"
```

## In-app editor (center pane)

A modal, Vim-style editor (configurable). Set `editor.keybindings` to `"vim"` or
`"default"` in `config.toml`.

The current mode is unmistakable: a **green `NORMAL`** / **blue `INSERT`** badge in
the editor's status line, and a **steady, color-coded cursor** — a green block in
Normal, a magenta block in Insert (it never blinks, so it can't disappear).

**Vim mode — Normal**
- Move: `h j k l` · `w`/`b` word fwd/back · `e` end-of-word · `0`/`^` line start ·
  `$` line end · `gg`/`G` top/bottom of file
- Enter Insert: `i` `a` (after) · `I`/`A` (line start/end) · `o`/`O` (open line below/above)
- Edit: `x` delete char · `r<char>` replace char · `dd` delete line · `dw` delete word ·
  `D` delete to end of line · `cc`/`cw`/`C` change line/word/to-end · `p` paste
- `Esc` returns to Normal (and cancels a half-typed operator like `d`)

**Command line (`:`)**
- `:submit` — check your solution (same as `Ctrl-R`)
- `:w` — save a draft and keep editing (resume the challenge later)
- `:q` — leave the app (`:wq` saves + submits)
- `:config` — open `config.toml` in your `$EDITOR`; on save, the layout is
  re-applied live (other settings take effect next launch)

`Ctrl-S` submit / `Ctrl-Q` quit work in any mode.

## Layouts

Set `ui.layout` in `config.toml` (or change it live with `:config`):

- **`vertical`** (default) — three side-by-side columns: list │ editor │ chat.
  Best for coding.
- **`horizontal`** — the list on the left, with the **content (lesson/chat) on
  top and your input on the bottom**. Better for reading- and writing-heavy
  subjects.

**Default mode** — ordinary typing; `Ctrl-S` submit, `Ctrl-Q` quit.

> `:w` (save & resume) is intentionally separate from `:submit` (check), so you can
> stop mid-solution and come back to it.

## Layout

```
main.go                 entry point: load config, construct deps, launch the TUI
internal/
  tui/                  the 3-pane Bubble Tea program (panes, async cmds, layout)
  config/               TOML config + defaults + flag overrides
  tutor/                OpenAI-compatible client (lesson/challenge/feedback/chat + offline)
  curriculum/           the built-in ordered learning path (modules + topics)
  editor/               embeddable Bubble Tea modal Vim editor
  executor/             runs the learner's code against tests (timeout-guarded)
  drafts/               save/load/clear in-progress solutions by challenge id
  progress/             progress.json — challenge attempts + curriculum topic status
```

## Notes / current limits

- The executor runs Python via `python3` with a 5s timeout. It is **not a
  sandbox** — fine for a single trusted local learner; don't run untrusted code.
- Vim mode is "Vim-*style*" (a useful core: motions, operators `d`/`c`, `r`, and
  the command line), not full Vim — there's no visual mode, counts, or undo.
- `Tab` and `Ctrl-W` are reserved for pane navigation, so they can't be typed
  into the editor.
- Drafts are preserved and resumed by challenge id; `Ctrl-N` asks the AI for a
  fresh challenge on the same topic, and you can revisit any past challenge from
  the left pane.
```

## Todo
- 
