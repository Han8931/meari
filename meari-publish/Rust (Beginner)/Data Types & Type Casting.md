---
created: "2026-07-08"
id: rust-b-types
source: meari-course
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

So `as` is a blunt tool — it always succeeds and never warns. For fallible,
checked conversions you'll later use `TryFrom`/`try_into`, but `as` is the
beginner's workhorse.

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


## Try it

1. Create variables of type `i32`, `f64`, `bool`, and `char`, then print them.
2. Try adding an `i32` and an `f64`. Read the error, then fix it with `as`.
3. Cast `300 as u8` and print the result. Does it match your expectation?
