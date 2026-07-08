---
created: "2026-07-08"
id: rust-b-string
source: meari-course
subject: Rust (Beginner)
title: String vs &str
---

Newcomers trip over Rust having *two* string types. It's not arbitrary — it's
[[Ownership & Moves|ownership]] and [[References & Borrowing|borrowing]] applied
to text. Once you see that, `String` vs `&str` clicks.

## The two types

| Type   | What it is                    | Owns the data? | Growable? | Lives where       |
| ------ | ----------------------------- | -------------- | --------- | ----------------- |
| `String` | an owned, heap-allocated string | **yes**      | yes       | heap + stack handle |
| `&str`   | a borrowed *view* into a string | no           | no        | points elsewhere  |

Think of it exactly like `Vec<T>` vs `&[T]` from
[[Arrays, Tuples & Slices]]: `String` is the owner, `&str` is a slice of one.

```
   let owned = String::from("hello");

   STACK              HEAP
   ┌──────────┐       ┌───┬───┬───┬───┬───┐
   │ ptr  ●───┼─────► │ h │ e │ l │ l │ o │   ← String owns this buffer
   │ len  5   │       └───┴───┴───┴───┴───┘
   │ cap  5   │              ▲
   └──────────┘              │
   let view: &str = &owned;  │   ← &str just borrows a window, owns nothing
```

## Where each comes from

```rust
let literal: &str = "hello";           // string literals are &'static str
let owned:   String = String::from("hello"); // or "hello".to_string()

let view: &str = &owned;               // borrow a String as a &str
let view2: &str = &owned[0..2];        // a sub-slice: "he"
```

A literal like `"hello"` is baked into your binary, so it's a `&'static str` — a
borrow that lives for the entire program. You never *own* a literal.

## The rule of thumb

```
   Need to build, grow, or own text?   →   String
   Only need to read/pass text?        →   &str  (take &str in parameters!)
```

Accepting `&str` in a function is more flexible than `String`, because a
`String` can be borrowed *as* a `&str` for free, but not vice versa:

```rust
fn greet(name: &str) {                 // accepts BOTH a literal and a &String
    println!("Hi, {name}");
}

greet("Ana");                          // &str literal
greet(&String::from("Bo"));            // String borrowed as &str
```

## Building and combining strings

```rust
let mut s = String::from("Hello");
s.push_str(", world");                 // append a &str
s.push('!');                           // append one char

let a = String::from("foo");
let b = String::from("bar");
let c = a + &b;                        // "foobar" — `a` is MOVED, `b` is borrowed
// println!("{a}");                    // ❌ a was consumed by `+`, no longer valid
```

The `+` operator reuses the left operand's buffer, so it *moves* `a` — which is
why `a` is unusable afterward. When you need all your inputs to survive, reach
for `format!`, which only *borrows* its arguments:

```rust
let first = String::from("Ana");
let last  = String::from("Smith");
let full  = format!("{first} {last}");  // "Ana Smith"
println!("{first} is still usable");    // ✅ format! borrowed — nothing moved
```

## The same in Python

Python has just **one** string type, so the `String` vs `&str` split is
Rust-specific. Python strings are also immutable, so "modifying" one actually
builds a new string:

```python
s = "Hello"
s += ", world"        # creates a NEW string; the old one is discarded
```

Loosely, a Rust `String` plays the role of the owned, growable buffer and `&str`
the role of a borrowed view — two jobs Python's `str` blurs together behind its
garbage collector. Like Rust, though, Python strings are Unicode, so iterating
by character (rather than raw bytes) is the safe habit in both languages.

## UTF-8: no integer indexing

Rust strings are UTF-8, and a character may be several bytes. So Rust
deliberately forbids `s[0]` — it would be ambiguous (a byte? a character?) and
could split a multi-byte character. Iterate instead:

```rust
let s = "héllo";
// let c = s[0];          // ❌ not allowed
for ch in s.chars() {      // iterate by Unicode character
    print!("{ch} ");       // h é l l o
}
println!("{}", s.len());   // 6 — BYTES, not characters (é is 2 bytes)
println!("{}", s.chars().count()); // 5 — actual character count
```

## Try it

1. Write a `greet(name: &str)` function and call it with both a string literal and a `String`.
2. Combine two strings with `format!` and confirm the originals are still usable.
3. Print both `s.len()` and `s.chars().count()` for `"héllo"`.

> **Takeaway:** `String` owns and grows; `&str` borrows and reads. Store owned
> text as `String`, but accept `&str` in your function signatures for maximum
> flexibility — and remember strings are UTF-8, so index by iteration, never by
> integer.
