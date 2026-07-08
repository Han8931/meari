---
created: "2026-07-08"
id: rust-b-cargo
source: meari-course
subject: Rust (Beginner)
title: Hello, Cargo
---

Before writing Rust you need the toolchain, and you'll drive almost all of it
through **Cargo**, Rust's build tool and package manager. Cargo compiles your
code, runs it, manages dependencies, runs tests, and builds documentation.

## The pieces of the toolchain

| Tool     | Job                                                            |
| -------- | ------------------------------------------------------------- |
| `rustup` | Installs and updates Rust itself; manages toolchain versions. |
| `rustc`  | The actual compiler. You rarely call it directly.            |
| `cargo`  | The tool you use daily — build, run, test, add dependencies. |

Install everything with one command from <https://rustup.rs>. After that,
`cargo` is your front door.

## Your first project

```bash
cargo new hello
cd hello
cargo run
```

`cargo new` scaffolds a project. Here's what it creates and how the parts fit:

```
hello/
├── Cargo.toml        ← project manifest: name, version, dependencies
├── Cargo.lock        ← exact resolved versions (auto-managed)
└── src/
    └── main.rs       ← your code; main() is the entry point
   (target/  appears after the first build — compiled output)
```

The generated `src/main.rs`:

```rust
fn main() {
    println!("Hello, world!");
}
```

- `fn main()` is the entry point — execution starts here.
- `println!` ends in `!`, which marks it as a **macro**, not a function. Macros
  are expanded at compile time; `println!` uses this to type-check your format
  string and arguments. You'll meet more macros like `vec!` and `format!` later.

## The everyday Cargo commands

| Command        | What it does                                            |
| -------------- | ------------------------------------------------------ |
| `cargo check`  | Type-checks fast **without** producing a binary.       |
| `cargo build`  | Compiles a debug binary into `target/debug/`.          |
| `cargo run`    | Builds (if needed) and runs.                           |
| `cargo test`   | Compiles and runs your tests.                          |
| `cargo fmt`    | Formats your code in the standard Rust style.          |
| `cargo build --release` | Optimized build into `target/release/` (slower to compile, faster to run). |

`cargo check` is your fastest feedback loop — use it constantly while wrestling
with the compiler. Save `--release` for benchmarking and shipping.

### If you're coming from Python

The everyday workflow maps almost one-to-one — Cargo simply bundles jobs that
Python spreads across several tools:

| Task              | Python                  | Rust (Cargo)        |
| ----------------- | ----------------------- | ------------------- |
| Run a program     | `python main.py`        | `cargo run`         |
| Start a project   | *(just make a `.py`)*   | `cargo new hello`   |
| Add a library     | `pip install requests`  | `cargo add reqwest` |
| Dependency file   | `requirements.txt`      | `Cargo.toml`        |
| Lock exact versions | `requirements.lock`   | `Cargo.lock`        |

## Reading a compiler error

Rust's compiler is unusually helpful. When something's wrong you'll see the
line, a caret pointing at it, and often a suggested fix:

```
error[E0425]: cannot find value `x` in this scope
 --> src/main.rs:2:20
  |
2 |     println!("{}", x);
  |                    ^ not found in this scope
```

Get used to reading these carefully — in Rust the compiler is a teacher, not
just a gatekeeper. Next up: [[Variables & Mutability]].


## Try it

1. Create a new project with `cargo new hello_rust`, then run it with `cargo run`.
2. Change the message in `println!` and run it again.
3. Run `cargo check` for fast feedback, then `cargo fmt` to format the code.
