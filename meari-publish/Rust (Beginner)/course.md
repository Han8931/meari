---
created: "2026-07-08"
id: rust-beginner
level: beginner
source: meari-course
title: Rust (Beginner)
---

This course assumes no Rust experience. Follow the lessons in order through
**Project Structure**: later chapters deliberately reuse the vocabulary and
mental models introduced earlier. **Optional: Deeper Memory** is enrichment,
not a requirement for being productive in beginner Rust.

The lessons are meant to be read slowly. New terms are introduced before they
are used, and examples explain not only what Rust accepts, but why it behaves
that way. You are not expected to infer missing steps or remember a rule after
seeing it once. Important ideas return in later lessons from a different angle.

For each lesson, use this learning loop:

1. Type the examples into a small Cargo project instead of only reading them.
2. Predict whether each example compiles and what it prints.
3. Run `cargo check`; read the first error fully before changing code.
4. Complete the exercises, then explain the takeaway in your own words.

Rust often feels slow at first because the compiler makes ownership, failure,
and types visible. Treat compiler errors as evidence about your current mental
model. Change one thing at a time and check again.

## Getting Started
- [[rust-b-why|Why Rust]]
- [[rust-b-cargo|Hello, Cargo]]
- [[rust-b-variables|Variables & Mutability]]
- [[rust-b-types|Data Types & Type Casting]]
- [[rust-b-control|Control Flow]]

## Ownership & Borrowing
- [[rust-b-ownership|Ownership & Moves]]
- [[rust-b-borrowing|References & Borrowing]]

## Compound Data
- [[rust-b-compound|Arrays, Tuples & Slices]]
- [[rust-b-structs|Structs & Methods]]
- [[rust-b-enums|Enums & Pattern Matching]]
- [[rust-b-string|String vs &str]]

## Errors & Collections
- [[rust-b-option-result|Option & Result]]
- [[rust-b-question|Error Propagation & Panics]]
- [[rust-b-collections|Vec & HashMap]]
- [[rust-b-iterators|Closures & Iterators]]

## Generics & Traits
- [[rust-b-generics|Generics]]
- [[rust-b-traits|Traits]]
- [[rust-b-derive|Common Derivable Traits]]

## Project Structure
- [[rust-b-modules|Modules & Cargo Dependencies]]

## Optional: Deeper Memory
- [[rust-b-smart-pointers|Box, Rc & RefCell]]
