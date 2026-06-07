<div align="center">

# 메아리 · Meari

**Your notes, echoed back as learning.**

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Terminal + Web](https://img.shields.io/badge/runs%20in-terminal%20%26%20browser-blueviolet)
![Local-first](https://img.shields.io/badge/local--first-your%20files-success)
![Offline OK](https://img.shields.io/badge/AI-optional%2C%20any%20provider-lightgrey)

</div>

---

## 🏔️ The name

**메아리** (*meari*) is the Korean word for **echo** — the voice that comes back
when you call out across a mountain valley.

That's the whole idea of this app. You put knowledge *into* your vault — notes you
write, lessons you ask for — and Meari sends it **back at you**: as questions, as
graded essays, as coding challenges, as whole courses built from what you wrote.
Reading fills a vault; the echo is what makes it stick.

## 📖 What it is

Meari is a **self-directed learning vault**: Obsidian-style markdown notes you own,
plus an AI tutor that turns *"I want to learn X"* into a saved lesson note — and then
echoes it back as study sessions until you actually remember it. Any subject:
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

## 💡 Why Meari?

- 📁 **Your notes are files, not a database.** Plain `.md` with frontmatter and
  `[[wikilinks]]`, in a folder you choose — point it at your existing Obsidian vault
  and it just works. Sync with git, edit with anything, leave anytime.
- ✍️ **AI lessons become notes, not chat scroll.** Ask to learn a topic and the tutor
  writes a focused lesson *into your vault*, linked to its prerequisites. Knowledge
  accumulates instead of evaporating.
- 🎓 **Notes become courses.** Open a note, type `:course`, answer two questions — an
  agentic pipeline plans a curriculum *from what you wrote*, writes the missing
  lessons, authors an exercise per topic, and **verifies every coding challenge by
  actually running its tests** before you ever see it. `:revise` polishes a course
  with your feedback.
- 🎁 **Courses are shareable.** `:publish` exports a course as a self-contained
  markdown folder — push it to a git repo, and anyone can drop it into their own
  Meari and study it. Your echo can teach someone else.
- 🧠 **Learning means recall, not re-reading.** Study any note actively: write an essay
  answer and get it graded, with quizzes and spaced-repetition flashcards on the
  roadmap. A built-in coding tutor (write code → hidden tests → feedback) covers the
  programming side today.
- 🔌 **Local-first and provider-agnostic.** Works offline with built-in content; plug in
  OpenAI, a local Ollama model, or any OpenAI-compatible endpoint for the AI parts.
  Nothing leaves your machine except the model calls you configure.

## 🚀 Quick start

```bash
go build -o meari .

./meari -vault        # the vault, in your terminal
./meari serve         # the same vault, in your browser
./meari               # the tutor (launch dashboard: pick a course)
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

| You type | Meari echoes back |
|---|---|
| `:learn the french revolution` | 📝 a lesson note, written into your vault |
| `:essay` → `:grade` | 🧠 an essay prompt on the open note, then a scored critique |
| `:course` | 🎓 a full course built from the open note (`:tutor` → `:topic` to take it) |
| `:revise make it harder` | 🔧 the same course, rebuilt around your feedback |
| `:publish` | 🎁 a shareable copy of the course, ready for `git push` |

Courses are plain markdown in `meari-course/` next to the app — your notes vault
stays untouched.

## ✨ Highlights

- 🤖 **Agentic course building** — `:course` interviews you in the chat pane
  (difficulty, scope, title — or just say "defaults"), then plans, writes, critiques,
  and verifies: code exercises run against the real executor and get repaired or
  demoted before shipping; dead wikilinks are stripped; a completeness critic adds
  what the outline missed. Courses are hand-editable markdown manifests.
- 🎁 **Course sharing** — `:publish` copies a course (manifest *plus* every linked
  topic note, so it's self-contained) into a publish directory meant to be a git
  repo. Recipients drop the folder into their `meari-course/` and study it like
  their own.
- 🌲 **A real file tree** in the sidebar — fold/unfold directories, and manage files
  NERDTree-style: `Space` to mark, `m` then **a**dd / **m**ove / **d**elete.
- ⌨️ **A modal Vim editor** with motions, counts, operators, visual mode, undo/redo, a
  jumplist (`Ctrl-O`/`Ctrl-I`), and markdown syntax highlighting — fenced code blocks
  get real Go/Python highlighting inside your notes.
- 🔍 **Fuzzy finding everywhere** — `,ff` jumps to any note, `,fg` greps every note's
  contents, from any pane.
- 📋 **A chat you can actually copy from** — drag to select transcript text and `Alt-C`
  it to your clipboard; yanks work even over SSH via OSC 52; `Alt-V`/`Cmd-V` paste.
- 💬 **A context-aware tutor chat** in every screen: replies stream live, render
  markdown in color, and always see the open note, lesson, or code you're working on.
- 🔗 **Obsidian-style backlinks**, `[[wikilink]]` navigation, and per-note chat history.
- 🧩 **One core, two faces** — the TUI and web UI stay in feature parity because neither
  contains business logic; both drive the same headless engine.

📘 **[Read the full manual →](docs/MANUAL.md)** — every key, command, and config option.

## 🗺️ Roadmap

- [ ] 🃏 Quiz study kind and spaced-repetition flashcards (SM-2)
- [ ] ⚡ SQLite index: fast search, link graph, SRS store
- [ ] 🕸️ Visual knowledge graph
- [ ] 🖥️ Desktop app (the core + web UI in a native window)

## 🌱 Status

In active development — Meari began as an AI coding tutor and is being generalized
into the subject-agnostic learning vault described above. The coding loop still works
end-to-end today.

If Meari's direction resonates with you, a ⭐ helps others find it — and if you build
a course worth sharing, `:publish` it and pass the echo on.
