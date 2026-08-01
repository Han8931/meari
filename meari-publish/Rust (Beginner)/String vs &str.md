---
created: "2026-07-08"
id: rust-b-string
source: meari-course
study:
  answer: |
    fn full_name(first: &str, last: &str) -> String {
        format!("{first} {last}")
    }
  kind: code
  lang: rust
  prompt: 'Complete `full_name(first: &str, last: &str) -> String` with `format!("{first} {last}")`. It borrows two string slices and returns a new owned `String`.'
  starter: |
    fn full_name(first: &str, last: &str) -> String {
        String::new()
    }
  tests:
    - assert_eq!(full_name("Ana", "Smith"), "Ana Smith");
    - assert_eq!(full_name("Rust", "Learner"), "Rust Learner");
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

For now, think “owned text” versus “borrowed view of text.” The next lesson,
[[Arrays, Tuples & Slices]], applies the same owner/view relationship to
sequences: `Vec<T>` versus `&[T]`.

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

A literal like `"hello"` is stored with the compiled program, so it can be
borrowed for the program's entire run. Rust writes that type as `&'static str`.
The `'static` part names how long the reference is valid; you do not need to
write it yourself in ordinary literal code.

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

## A string slice is pointer plus length

`&str` is not the string data itself. Conceptually it contains a pointer to the
first byte of valid UTF-8 and a byte length. It may refer to an entire `String`,
a literal stored in the program, or part of another string:

```rust
let owned = String::from("hello world");
let first = &owned[0..5]; // borrows "hello"; no characters are copied
```

The byte indexes must land on UTF-8 character boundaries or slicing panics.
That makes arbitrary user-facing slicing safer with `.chars()` or specialized
Unicode libraries. Also, `&String` means “borrow this particular owned container,”
while `&str` means “borrow text from any source”; this is why parameters normally
use `&str`.

## Syntax checkpoint

Read the ampersand before reading the word `str`:

```rust
fn count(text: &str) -> usize {
    text.len()
}
```

`&str` means “borrowed text”; `String` means “owned text.” The function above
only reads, so it borrows. A caller with a `String` writes `count(&owned)`, while
a string literal can be passed directly with `count("hello")`.

You are ready to continue when you can choose `&str` for read-only input and
`String` for a returned value that must own its text. The next lesson uses the
same owner/view distinction for sequences: an owned array and a borrowed
`&[T]` slice.

## Try it

1. **Classify:** For a literal, a `String::from(...)`, and `&owned`, state which
   value owns text and which only borrows it.
2. **Fill:** Write `greet(name: &str)` and call it with both a literal and a
   borrowed `String`.
3. **Compare ownership:** Combine two strings with `format!` and confirm the
   originals remain usable.
4. **Measure:** Print both `s.len()` and `s.chars().count()` for `"héllo"` and
   explain why the numbers differ.

> **Takeaway:** `String` owns and grows; `&str` borrows and reads. Store owned
> text as `String`, but accept `&str` in your function signatures for maximum
> flexibility — and remember strings are UTF-8, so index by iteration, never by
> integer.
