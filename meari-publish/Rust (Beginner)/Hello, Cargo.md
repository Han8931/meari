---
created: "2026-07-08"
id: rust-b-cargo
source: meari-course
study:
  answer: |
    fn fastest_check_command() -> &'static str {
        "cargo check"
    }
  kind: code
  lang: rust
  prompt: Fill in the missing string returned by `fastest_check_command`. It should return the Cargo command that checks your code without producing the final executable. Only replace `"TODO"`; the function syntax is provided for you.
  starter: |
    fn fastest_check_command() -> &'static str {
        "TODO" // replace this string
    }
  tests:
    - assert_eq!(fastest_check_command(), "cargo check");
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
- `println!` prints a line. The `!` means it is a **macro**, a Rust feature
  that looks like a function call. For now, treat `println!` as the standard
  way to print; the next lesson explains the punctuation in this line one
  piece at a time. You'll meet `vec!` and `format!` later.

## The everyday Cargo commands

| Command        | What it does                                            |
| -------------- | ------------------------------------------------------ |
| `cargo check`  | Type-checks fast **without** producing a binary.       |
| `cargo build`  | Compiles a debug binary into `target/debug/`.          |
| `cargo run`    | Builds (if needed) and runs.                           |
| `cargo test`   | Compiles and runs your tests.                          |
| `cargo fmt`    | Formats your code in the standard Rust style.          |
| `cargo build --release` | Optimized build into `target/release/` (slower to compile, faster to run). |

`cargo check` is the quickest way to ask Rust whether the code is valid. Run it
often while learning, especially after a small change. This keeps error messages
focused on the idea you are currently practicing. Save `--release` for
benchmarking and shipping.

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

Read the message from the top and begin with the first reported error. Later
errors can be consequences of that first problem. The explanation after the
caret describes what Rust expected at the marked location, and the suggestion
is a possible fix rather than a command you must follow. Next up:
[[Reading Rust Code, One Piece at a Time]].

## What happens when you run `cargo run`?

It helps to separate the steps that Cargo normally bundles together:

1. Cargo reads `Cargo.toml` to learn the package name and dependencies.
2. It asks `rustc` to compile `src/main.rs` and any dependencies.
3. If compilation succeeds, it writes an executable under `target/debug/`.
4. It starts that executable. Only now does your `main` function run.

This explains an important distinction: a **compiler error** means no program
was produced, while a **runtime error** happens after compilation, while the
program is executing. Rust catches many mistakes in the first category.

```rust
fn main() {                     // define main; it takes no input
    println!("answer: {}", 42); // ! calls a macro; ; ends the statement
}
```

Braces delimit the function body and parentheses hold arguments. Rust ignores
most whitespace, but punctuation is meaningful. You need not memorize it all
now—use `cargo fmt`, and let repeated examples make it familiar.


## Syntax checkpoint

At this point, you only need to recognize this shape—you do not need to write it
from memory yet:

```rust
fn main() {
    println!("Hello");
}
```

Read it as: “define `main`; inside it, print one line.” If that sentence makes
sense and you can run the program with `cargo run`, you are ready to continue.
The next lesson pauses on every punctuation mark in this example before adding
variables.

## Try it

1. **Order:** Put “compile,” “run `main`,” and “read `Cargo.toml`” in the order
   used by `cargo run`.
2. **Build:** Create a project with `cargo new hello_rust`, then run it.
3. **Repair:** Introduce an undefined variable, run `cargo check`, and identify
   the error code, source location, and explanation.
4. **Practice:** Change the message, run `cargo fmt`, then run the program again.
