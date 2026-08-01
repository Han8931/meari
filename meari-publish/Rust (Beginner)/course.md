---
created: "2026-07-08"
id: rust-beginner
level: beginner
source: meari-course
title: Rust (Beginner)
---

This course assumes no Rust experience and no experience with any particular
other language. Follow the lessons in order through **Project Structure**.
The first six lessons deliberately teach the punctuation and shapes of Rust
code before asking you to use them in larger examples. Python comparisons are
optional side notes; you never need Python to understand the Rust explanation.
**Optional: Deeper Memory** is enrichment, not a requirement for productive
beginner Rust.

Read slowly and type the examples. A new syntax form is first named, then
shown in a tiny program, then reused in an exercise. Do not try to memorize all
of Rust's punctuation at once: initially, focus only on the marked pieces in
the current lesson. Important ideas return later from a different angle.

For each lesson, use this learning loop:

1. Read the opening example and the **Takeaway** first.
2. Type that one example into a small Cargo project and change one small thing.
3. Run `cargo check`; read the first error fully before changing code.
4. Complete the exercise, then explain the takeaway in your own words.

### What “finished” means on a first pass

You do not need to retain every section of every lesson before moving on. The
first pass should leave you with one usable idea per lesson:

| Lessons | Leave able to… | It is fine to skim for now… |
| --- | --- | --- |
| Cargo → Syntax → Variables | run a project; read `fn`, `{}`, `let`, strings, and `;`; choose `mut` when changing a value | macros, shadowing, and `const` details |
| Types → Functions → Control Flow | read basic types and `as`; write a typed function; use `if` and a `for` range | integer widths, overflow policies, labeled loops |
| Why Rust → Ownership → Borrowing → Strings | explain why a `String` moves; choose `&T` or `&mut T`; choose `String` or `&str` | heap diagrams, partial moves, inferred lifetimes |
| Compound Data → Structs → Enums | make an array, struct, or enum; read a `match` | tuple/unit structs, struct update syntax, `let else` |
| Option → Errors → Collections | handle `Some`/`None` and `Ok`/`Err`; use `?`; use a `Vec` or `HashMap` | `unwrap` alternatives, mixed error types, capacity and map variants |
| Generics → Traits → Derive → Iterators | recognize `<T>`; implement a small trait; derive `Debug`; use `.iter().map(...).sum()` | monomorphization, trait objects, closure trait names |
| Modules | split a small program into modules and add a dependency | visibility edge cases and larger package layouts |

Sections explicitly marked **Later** or **Optional** are reference material.
They are there when you meet the syntax in real code, not a gate you must pass
before the next lesson. Rust often feels slow at first because the compiler
makes ownership, failure, and types visible. Treat compiler errors as evidence
about your current mental model. Change one thing at a time and check again.

## Getting Started
- [[rust-b-cargo|Hello, Cargo]]
- [[rust-b-syntax|Reading Rust Code, One Piece at a Time]]
- [[rust-b-variables|Variables & Mutability]]
- [[rust-b-types|Data Types & Type Casting]]
- [[rust-b-functions|Functions & Expressions]]
- [[rust-b-control|Control Flow]]
- [[rust-b-why|Why Rust]]

## Ownership & Borrowing
- [[rust-b-ownership|Ownership & Moves]]
- [[rust-b-borrowing|References & Borrowing]]
- [[rust-b-string|String vs &str]]

## Compound Data
- [[rust-b-compound|Arrays, Tuples & Slices]]
- [[rust-b-structs|Structs & Methods]]
- [[rust-b-enums|Enums & Pattern Matching]]

## Errors & Collections
- [[rust-b-option-result|Option & Result]]
- [[rust-b-question|Error Propagation & Panics]]
- [[rust-b-collections|Vec & HashMap]]

## Generics & Traits
- [[rust-b-generics|Generics]]
- [[rust-b-traits|Traits]]
- [[rust-b-derive|Common Derivable Traits]]

## Functional Programming
- [[rust-b-iterators|Closures & Iterators]]

## Project Structure
- [[rust-b-modules|Modules & Cargo Dependencies]]

## Optional: Deeper Memory
- [[rust-b-smart-pointers|Box, Rc & RefCell]]
