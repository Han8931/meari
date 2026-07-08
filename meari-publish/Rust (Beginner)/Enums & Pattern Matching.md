---
created: "2026-07-08"
id: rust-b-enums
source: meari-course
subject: Rust (Beginner)
title: Enums & Pattern Matching
---

If a struct is "**A and** B **and** C" (all fields at once), an **enum** is
"**A or** B **or** C" (exactly one variant at a time). Enums plus `match` are one
of Rust's most powerful and distinctive features — they replace null checks,
type-code fields, and much of what class hierarchies do elsewhere.

## An enum is a set of variants

```rust
enum Direction {
    North,
    South,
    East,
    West,
}

let heading = Direction::North;
```

The real power appears when variants **carry data** — and each variant can carry
a *different shape* of data:

```rust
enum Message {
    Quit,                        // no data
    Move { x: i32, y: i32 },     // named fields, like a struct
    Write(String),               // a single value
    ChangeColor(u8, u8, u8),     // a tuple of values
}
```

```
   Message
   ├── Quit                      (nothing)
   ├── Move { x, y }             (struct-like)
   ├── Write(String)             (one payload)
   └── ChangeColor(u8, u8, u8)   (tuple payload)
        exactly ONE of these at a time
```

## `match`: exhaustive pattern matching

`match` compares a value against patterns and runs the first arm that fits. Its
superpower is **exhaustiveness** — the compiler forces you to handle *every*
variant, so you can't forget a case:

```rust
fn describe(msg: Message) -> String {
    match msg {
        Message::Quit => "quitting".to_string(),
        Message::Move { x, y } => format!("moving to ({x}, {y})"),
        Message::Write(text) => format!("writing: {text}"),
        Message::ChangeColor(r, g, b) => format!("color #{r:02x}{g:02x}{b:02x}"),
    }
}
```

Each arm can **bind** the data inside the variant (`x`, `y`, `text`, `r/g/b`)
straight into local variables. If you forget a variant, the program won't
compile — add a new `Message` variant later and the compiler points you at every
`match` that needs updating.

### The `_` catch-all

When you don't want to spell out every case, `_` matches "anything else":

```rust
let n = 3;
match n {
    1 => println!("one"),
    2 => println!("two"),
    _ => println!("something else"),   // required to be exhaustive
}
```

## The same in Python

Modern Python (3.10+) has structural pattern matching that reads much like
`match`:

```python
match msg:
    case Quit():           print("quitting")
    case Move(x, y):       print(f"moving to ({x}, {y})")
    case Write(text):      print(f"writing: {text}")
    case _:                print("something else")
```

The decisive difference is **exhaustiveness**. Rust checks at compile time that
every variant is handled — forget one and the program won't build. Python's
`match` has no such check, so a missing case simply falls through at runtime.

## `if let` and `let else`: matching one case

When you care about a *single* pattern, a full `match` is noisy. `if let` is the
concise form:

```rust
let config: Option<i32> = Some(5);

// verbose
match config {
    Some(v) => println!("got {v}"),
    None => {}
}

// concise — same thing
if let Some(v) = config {
    println!("got {v}");
}
```

`let else` handles the "bind it or bail out" pattern cleanly:

```rust
let Some(v) = config else {
    println!("no config; giving up");
    return;
};
// v is available for the rest of the function
```

## Why this matters

Enums-plus-match are how Rust encodes "this value is one of a fixed set of
shapes" *in the type system*. Two of the most important types in the whole
language — `Option` and `Result` — are just enums, which is exactly why a later
lesson, [[Option & Result]], builds directly on everything here.

## Try it

1. Define an enum `Light` with `Red`, `Yellow`, and `Green`.
2. Use `match` to print what a driver should do for each light.
3. Rewrite one single-case `match` as `if let`.

> **Takeaway:** model mutually exclusive alternatives as enum variants (which can
> carry data), then handle them with `match`. Exhaustiveness turns "I forgot a
> case" from a runtime bug into a compile error.
