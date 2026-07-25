---
created: '2026-07-21'
id: rust-intermediate
level: intermediate
source: meari-course
title: Rust (Intermediate)
---

This course picks up where **Rust (Beginner)** ends. It assumes you are
comfortable with ownership, borrowing, structs, enums, `Option`/`Result`,
collections, closures, generics, and basic traits. If any of those feel shaky,
revisit the beginner course first — every topic here builds on them.

The same learning loop applies:

1. Type each example into a small Cargo project instead of only reading it.
2. Before compiling, trace the important type, owner, borrow, or thread at each
   line.
3. Predict whether it compiles and what it produces, then run `cargo check`.
4. Read the first diagnostic completely and change one thing at a time.
5. Do the study challenge, then explain the idea back in one sentence.

Intermediate Rust is mostly about the type system working *for* you: lifetimes
that make borrowing precise, traits that abstract behavior, and the ownership
rules extending cleanly to threads. Read the lessons in order — later chapters
lean on the vocabulary earlier ones introduce.

Each lesson ends with a progression rather than a single leap: first trace or
classify existing code, then repair or fill a small part, and finally implement
the embedded study challenge. Sections marked optional explain a sharp edge
without making it a prerequisite for the following lesson.

## Lifetimes
- [[rust-i-lifetimes|Lifetime Annotations]]
- [[rust-i-lifetime-structs|Lifetimes in Structs]]

## Traits in Depth
- [[rust-i-bounds|Trait Bounds & where Clauses]]
- [[rust-i-associated-types|Associated Types]]
- [[rust-i-operator|Operator Overloading]]
- [[rust-i-trait-objects|Trait Objects & Dynamic Dispatch]]

## Conversions & Errors
- [[rust-i-from-into|From & Into]]
- [[rust-i-tryfrom|TryFrom & Fallible Conversion]]
- [[rust-i-custom-errors|Custom Error Types]]
- [[rust-i-error-propagation|Error Conversion with ?]]

## Closures & Iterators
- [[rust-i-fn-traits|Fn, FnMut & FnOnce]]
- [[rust-i-custom-iterator|Implementing Iterator]]
- [[rust-i-iterator-pipelines|Iterator Pipelines]]

## Shared State
- [[rust-i-interior-mutability|Interior Mutability]]

## Concurrency
- [[rust-i-threads|Threads & move]]
- [[rust-i-arc-mutex|Shared State with Arc & Mutex]]
- [[rust-i-channels|Channels]]

## Advanced Patterns
- [[rust-i-match-advanced|Advanced Pattern Matching]]
