---
created: "2026-07-08"
id: rust-b-types
source: meari-course
study:
  answer: |
    fn average(a: i32, b: i32) -> f64 {
        (a as f64 + b as f64) / 2.0
    }
  kind: code
  lang: rust
  prompt: 'Write `average(a: i32, b: i32) -> f64` returning the mean of the two values as an f64 (cast with `as f64`).'
  starter: |
    fn average(a: i32, b: i32) -> f64 {
        0.0
    }
  tests:
    - assert_eq!(average(2, 4), 3.0);
    - assert_eq!(average(1, 2), 1.5);
    - assert_eq!(average(-3, 3), 0.0);
subject: Rust (Beginner)
title: Data Types & Type Casting
---

Rust is statically typed: every value has a type known at compile time. Often
the compiler *infers* it, but you can always annotate with `: Type`. This lesson
covers the **scalar** types (single values) and how to convert between them.

## The scalar types

```rust
let count: i32   = -42;      // signed 32-bit integer
let size:  usize = 10;       // pointer-sized unsigned (for indexing)
let ratio: f64   = 3.14;     // 64-bit floating point
let ok:    bool  = true;     // true / false
let grade: char  = 'A';      // a single Unicode scalar, in single quotes
```

### Integers come in many widths

| Signed | Unsigned | Bits | Approx. range (signed)          |
| ------ | -------- | ---- | ------------------------------- |
| `i8`   | `u8`     | 8    | -128 … 127                      |
| `i16`  | `u16`    | 16   | ±32k                            |
| `i32`  | `u32`    | 32   | ±2.1 billion  *(default `int`)* |
| `i64`  | `u64`    | 64   | ±9.2 quintillion                |
| `isize`| `usize`  | arch | pointer-sized (used for indexes)|

If you don't annotate, an integer literal defaults to `i32` and a float to
`f64`. Use `usize` when indexing into collections — that's what the standard
library expects.

## No implicit numeric coercion

This is a common surprise. Rust will **not** silently mix number types for you:

```rust
let a: i32 = 10;
let b: f64 = 2.5;
let c = a + b;          // ❌ error: cannot add f64 to i32
let c = a as f64 + b;   // ✅ 12.5 — you must convert explicitly
```

You convert with the `as` keyword:

```rust
let decimal: f64 = 54.321;
let integer = decimal as u16;   // 54 — the fraction is TRUNCATED, not rounded
let letter  = 65u8 as char;     // 'A'
let byte    = 'A' as u8;        // 65
```

### Watch out when casting *down*

Casting to a smaller type can silently lose information:

```
  300  as  u8      →   44        (300 mod 256, wraps around)
 -1    as  u8      →   255       (bit pattern reinterpreted)
  3.99 as  i32     →   3         (truncates toward zero)
```

```
      i32: [ ........ 300 ........ ]
                       │  cast to u8 keeps only low 8 bits
                       ▼
      u8 :          [ 44 ]        (300 - 256)
```

An `as` conversion always produces a value, even when information is lost. Use
it when that behavior is understood and intentional. When a value might not fit,
`TryFrom` or `try_into` can report failure instead; the checked example below
shows what that difference looks like.

### The same in Python

Python converts with constructor-style functions, and mixes number types
*automatically* — the opposite of Rust:

```python
decimal = 54.321
integer = int(decimal)     # 54 — truncates, like `as u16`
letter  = chr(65)          # 'A'
byte    = ord('A')         # 65
n = 10 + 2.5               # 12.5 — Python silently promotes int → float
```

That last line is the key difference: Python happily evaluates `10 + 2.5`, while
Rust rejects `a + b` across types until you cast explicitly. Python also has
arbitrary-precision integers, so it never overflows the way the next section
describes.

## Overflow behavior

What happens when arithmetic exceeds a type's range depends on the build:

- **Debug build** (`cargo run`): the program **panics** — it crashes loudly so
  you notice the bug.
- **Release build** (`--release`): the value **wraps around** silently.

When you actually *want* a defined behavior, the standard library gives you
explicit methods:

```rust
let x: u8 = 255;
x.wrapping_add(1);   // 0   — wrap on purpose
x.checked_add(1);    // None — returns Option, Some(_) if it fit
x.saturating_add(1); // 255 — clamp at the max
```

That `Option` return connects to [[Option & Result]], where absence-as-a-value
becomes a central theme. Next: putting values to work in [[Control Flow]].

## Literals, annotations, and suffixes

The text `42` is an integer **literal**. Context usually determines its type:

```rust
let a = 42;          // defaults to i32
let b: u64 = 42;     // annotation supplies the type
let c = 42u8;        // a suffix supplies the type
let d = 1_000_000;   // underscores improve readability only
```

When an error says “expected `usize`, found `i32`,” it is describing the two
sides of an operation, not claiming either type is universally wrong.

## Prefer checked conversion for uncertain data

Use `as` when truncation or wrapping is intentional. For uncertain input, a
checked conversion is safer:

```rust
let large: u16 = 300;
let byte = u8::try_from(large);

match byte {
    Ok(n) => println!("converted: {n}"),
    Err(_) => println!("300 does not fit in a u8"),
}
```

You will learn `Result` and `match` later. For now, notice that `as` produces a
value unconditionally, while `try_from` reports failure.


## Try it

1. Create variables of type `i32`, `f64`, `bool`, and `char`, then print them.
2. Try adding an `i32` and an `f64`. Read the error, then fix it with `as`.
3. Cast `300 as u8` and print the result. Does it match your expectation?
