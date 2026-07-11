---
created: "2026-07-08"
id: rust-b-why
source: meari-course
study:
  answer: |
    fn rust_checks_memory_safety() -> &'static str {
        "compile time"
    }
  kind: code
  lang: rust
  prompt: The function is already written. Replace the `"TODO"` placeholder with `"compile time"` to show when Rust checks many memory-safety rules. You only need to fill in one line.
  starter: |
    fn rust_checks_memory_safety() -> &'static str {
        "TODO" // replace this string
    }
  tests:
    - assert_eq!(rust_checks_memory_safety(), "compile time");
subject: Rust (Beginner)
title: Why Rust
---

Rust is a systems programming language built around one bold promise: **memory
safety without a garbage collector**. It aims to give you the raw speed of C and
C++ while eliminating whole categories of bugs — dangling pointers, data races,
buffer overflows — *before your program ever runs*, at compile time.

## The three pitches

1. **Safety** — the compiler's *borrow checker* proves your memory access is
   valid. In *safe* Rust the compiler prevents null-pointer dereferences,
   use-after-free, double-free, and data races. (Rust still has `unsafe`, FFI,
   deliberate panics, and ordinary logic bugs — safety here means memory safety,
   not the absence of all bugs.)
2. **Speed** — Rust compiles to native machine code with no runtime and no
   garbage collector pausing your program. "Zero-cost abstractions": high-level
   code compiles down to what you'd have written by hand.
3. **Concurrency** — the same rules that keep memory safe also make data races a
   *compile error*. This is Rust's famous "fearless concurrency."

## Rust vs. a language you already know

If you're coming from Python, the mental model is quite different:

| Dimension        | Python                    | Rust                          |
| ---------------- | ------------------------- | ----------------------------- |
| Execution        | Interpreted at runtime    | Compiled to a native binary   |
| Typing           | Dynamic (checked at run)  | Static (checked at compile)   |
| Memory managed by| Garbage collector         | Ownership & lifetimes         |
| Errors surface   | Often at runtime          | Mostly at compile time        |
| Speed            | Slower, GC pauses         | C-level, predictable          |

The trade is real: you do more work up front to satisfy the compiler. In return,
"if it compiles, it usually works" becomes a genuine experience rather than a
slogan.

Even "Hello, world" hints at the difference in philosophy:

```rust
fn main() {
    println!("Hello, world!");   // typed, compiled, no runtime
}
```

```python
print("Hello, world!")           # dynamic, interpreted at runtime
```

Python wins on brevity here — but Rust's extra ceremony is what buys the
compile-time guarantees in the table above. Throughout this course each Rust
idea is paired with its Python equivalent, so you can lean on what you already
know while learning what's genuinely new.

## How the safety model fits together

Rust's guarantees come from a small set of rules that build on each other. This
whole beginner course walks up that ladder:

```
        Ownership          "each value has exactly one owner"
            │
            ▼
        Borrowing          "you can lend a value out with & references"
            │
            ▼
        Lifetimes          "a borrow can't outlive what it points to"
            │
            ▼
   Memory safety + no data races, checked at COMPILE TIME
```

Don't worry about the details yet — [[Ownership & Moves]] and
[[References & Borrowing]] cover them. The point is that this is the *core idea*
of the language, and everything else orbits it.

## Where Rust is used

- **Systems & CLI tools** — ripgrep, fd, the `uv` Python packager.
- **WebAssembly** — compiling to fast, safe code that runs in the browser.
- **Embedded & OS work** — parts of the Linux kernel, Android, Windows.
- **Backend services** — where predictable latency (no GC pauses) matters.

## When Rust is *not* the answer

Rust asks for patience. For a quick script, a one-off data-munging task, or a
prototype you'll throw away tomorrow, a garbage-collected language is often the
faster path. Rust pays off when correctness, performance, or long-term
maintenance justify the up-front rigor.

## Try it

1. In one sentence, explain what Rust means by “memory safety without a garbage collector.”
2. Name one task where Python is probably the better choice, and one where Rust might pay off.
3. Look at the Python/Rust comparison table and pick the difference you think will affect you most.

> **Takeaway:** Rust moves an entire class of bugs from *runtime crashes* to
> *compiler errors*, at the cost of a steeper learning curve driven by the
> ownership system.
