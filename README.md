# Meari

> **메아리** (*meari*) is Korean for **"echo"** — what you learn should come back to you.

Meari is a **self-directed learning vault**: Obsidian-style markdown notes you own,
plus an AI tutor that turns *"I want to learn X"* into a saved lesson note — and then
**echoes it back** as study sessions until you actually remember it. Any subject:
languages, math, history, science, code.

```
┌ notes ──────┐┌──────── editor ─────────┐┌──── chat / study ────┐
│ ▾ math      ││ # Derivatives           ││ lesson  …            │
│    limits   ││ A derivative measures…  ││ tutor   …            │
│ ▸ spanish   ││ [[Limits]] first.       ││ > ask the tutor…     │
└─────────────┘└─────────────────────────┘└──────────────────────┘
```

It runs as a fast **terminal app** and a **local web app** — two thin front-ends over
one shared Go core, working on the same plain-markdown vault.

## Why Meari?

- **Your notes are files, not a database.** Plain `.md` with frontmatter and
  `[[wikilinks]]`, in a folder you choose — point it at your existing Obsidian vault
  and it just works. Sync with git, edit with anything, leave anytime.
- **AI lessons become notes, not chat scroll.** Ask to learn a topic and the tutor
  writes a focused lesson *into your vault*, linked to its prerequisites. Knowledge
  accumulates instead of evaporating.
- **Notes become courses.** Open a note, type `:course`, answer two questions — an
  agentic pipeline plans a curriculum *from what you wrote*, writes the missing
  lessons, authors an exercise per topic, and **verifies every coding challenge by
  actually running its tests** before you ever see it. `:revise` polishes a course
  with your feedback. Then study it in the tutor: essays graded by AI, code checked
  against hidden tests.
- **Learning means recall, not re-reading.** Study any note actively: write an essay
  answer and get it graded, with quizzes and spaced-repetition flashcards on the
  roadmap. A built-in coding tutor (write code → hidden tests → feedback) covers the
  programming side today.
- **Local-first and provider-agnostic.** Works offline with built-in content; plug in
  OpenAI, a local Ollama model, or any OpenAI-compatible endpoint for the AI parts.
  Nothing leaves your machine except the model calls you configure.

## Quick start

```bash
go build -o meari .

./meari -vault        # the vault, in your terminal
./meari serve         # the same vault, in your browser
./meari               # the coding tutor (guided setup)
```

Point it at notes you already have, and optionally wire up an AI:

```toml
# config.toml
[vault]
dir = "~/Documents/my-notes"   # default: ./vault

[ai]
provider = "ollama"            # or "openai" / any compatible endpoint
model = "llama3.1"
```

`meari check` verifies your AI setup end-to-end. Then, inside the vault:
`:learn the french revolution` writes a lesson note; `:essay` quizzes you on the open
note and `:grade` scores your answer; **`:course` turns the open note into a full
course** (`:tutor` → `:topic <name>` to take it). Courses are plain markdown in
`meari-course/` next to the app — your notes vault stays untouched.

## Highlights

- **Agentic course building** — `:course` interviews you in the chat pane
  (difficulty, scope, title — or just say "defaults"), then plans, writes, critiques,
  and verifies: code exercises run against the real executor and get repaired or
  demoted before shipping; dead wikilinks are stripped; a completeness critic adds
  what the outline missed. Courses are hand-editable markdown manifests.
- **A real file tree** in the sidebar — fold/unfold directories, and manage files
  NERDTree-style: `Space` to mark, `m` then **a**dd / **m**ove / **d**elete.
- **A modal Vim editor** with motions, counts, operators, visual mode, undo/redo, a
  jumplist (`Ctrl-O`/`Ctrl-I`), and markdown syntax highlighting — fenced code blocks
  get real Go/Python highlighting inside your notes.
- **System clipboard both ways** — yanks land in your clipboard (even over SSH, via
  OSC 52); `Alt-V`/`Cmd-V` paste from anywhere.
- **A context-aware tutor chat** in every screen: replies stream live and always see
  the open note, lesson, or code you're working on.
- **Obsidian-style backlinks**, `[[wikilink]]` navigation, and per-note chat history.
- **One core, two faces** — the TUI and web UI stay in feature parity because neither
  contains business logic; both drive the same headless engine.

**[Read the full manual →](docs/MANUAL.md)** — every key, command, and config option.

## Roadmap

- [x] Markdown vault with frontmatter + wikilinks (bring your own Obsidian vault)
- [x] AI lessons as notes · essay study with grading · tutor chat
- [x] **Agentic courses from your notes** — `:course` / `:revise`, executor-verified
      exercises, runnable in the tutor
- [x] Terminal + web front-ends over one headless core
- [x] Built-in offline coding curricula (Go, Python) with hidden-test challenges
- [ ] Quiz study kind and spaced-repetition flashcards (SM-2)
- [ ] SQLite index: fast search, link graph, SRS store
- [ ] Visual knowledge graph
- [ ] Desktop app (the core + web UI in a native window)

## Status

In active development — Meari began as an AI coding tutor and is being generalized
into the subject-agnostic learning vault described above. The coding loop still works
end-to-end today.

If Meari's direction resonates with you, a ⭐ helps others find it.
