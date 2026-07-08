---
created: "2026-07-08"
id: rust-b-derive
source: meari-course
subject: Rust (Beginner)
title: Common Derivable Traits
---

Implementing [[Traits]] by hand is fine for behavior unique to your type — but
several traits are so routine (printing, copying, comparing) that writing them
out would be pure boilerplate. The `#[derive(...)]` attribute tells the compiler
to generate those implementations for you, straight from your struct's fields.

## `#[derive(...)]` in one line

```rust
#[derive(Debug, Clone, PartialEq)]
struct Point {
    x: i32,
    y: i32,
}
```

That single attribute just gave `Point` three capabilities. Let's meet the most
common derivable traits.

## `Debug` — developer-facing printing

You've used `{}` to print scalars, but it won't print a struct. `Debug` enables
the `{:?}` and pretty `{:#?}` formats, meant for debugging:

```rust
let p = Point { x: 1, y: 2 };
println!("{p:?}");     // Point { x: 1, y: 2 }
println!("{p:#?}");    // pretty, multi-line
```

> Note: `{}` (the `Display` trait) is *not* derivable — you write it by hand
> when you want polished user-facing output. `{:?}` (`Debug`) is the one you
> derive, and the one you'll reach for constantly while developing.

## `Clone` and `Copy` — duplicating values

These tie straight back to [[Ownership & Moves]]. `Clone` gives you an explicit
deep copy via `.clone()`; `Copy` makes assignment *duplicate* the value instead
of moving it:

```rust
#[derive(Clone, Copy)]
struct Coord { x: i32, y: i32 }

let a = Coord { x: 1, y: 2 };
let b = a;             // COPIED, not moved (because Coord is Copy)
println!("{}", a.x);   // ✅ a is still valid — recall the move rules
```

Two rules worth remembering:

- `Copy` requires `Clone` (it's the cheap, implicit subset), so you derive them
  together.
- You can only derive `Copy` if **every field is itself `Copy`** — so a struct
  containing a `String` or `Vec` can be `Clone` but never `Copy`.

## `PartialEq` / `Eq` — equality

Derive `PartialEq` to compare with `==` and `!=`:

```rust
#[derive(PartialEq)]
struct Version(u32, u32);

Version(1, 0) == Version(1, 0);   // true
Version(1, 0) == Version(2, 0);   // false
```

## `PartialOrd` / `Ord` — ordering

These enable `<`, `>`, and sorting. This is what the `largest<T: PartialOrd>`
function in [[Generics]] required of its type:

```rust
#[derive(PartialEq, Eq, PartialOrd, Ord)]
struct Score(u32);

let mut scores = vec![Score(30), Score(10), Score(20)];
scores.sort();          // works because Score implements Ord
```

## `Default` — a sensible zero value

```rust
#[derive(Default)]
struct Config { verbose: bool, level: u32 }

let c = Config::default();   // Config { verbose: false, level: 0 }
```

## The derivable traits at a glance

| Derive        | Gives you                          | Example                    |
| ------------- | ---------------------------------- | -------------------------- |
| `Debug`       | `{:?}` / `{:#?}` printing           | `println!("{p:?}")`        |
| `Clone`       | explicit `.clone()` deep copy       | `let q = p.clone();`       |
| `Copy`        | implicit copy on assign (stack only)| `let q = p;` (p still ok)  |
| `PartialEq`   | `==` and `!=`                       | `a == b`                   |
| `PartialOrd`/`Ord` | `<`, `>`, `.sort()`              | `v.sort()`                 |
| `Default`     | `Type::default()`                   | `Config::default()`        |

The one requirement: derive only works if **every field also implements that
trait**. Derive `PartialEq` on a struct whose fields are all comparable and it
just works; include a field that isn't, and the compiler tells you.

## The same in Python

Python's `@dataclass` is a strikingly close parallel — it auto-generates
`__init__`, `__repr__`, and `__eq__` from your field list, just as derive
generates `Debug`, `Clone`, and `PartialEq`:

```python
from dataclasses import dataclass

@dataclass                      # ~ #[derive(Debug, Clone, PartialEq)]
class Point:
    x: int
    y: int

p = Point(1, 2)
print(p)                        # Point(x=1, y=2)   ~ derive(Debug)
Point(1, 2) == Point(1, 2)      # True              ~ derive(PartialEq)
```

The difference in flavor: Python's dataclass bundles a common set on by default,
while Rust makes each capability an explicit opt-in — and, for `Copy`, ties it
directly to the ownership model you learned earlier.

## Try it

1. Add `#[derive(Debug)]` to a struct and print it with `{:?}`.
2. Add `PartialEq` and compare two values with `==`.
3. Try deriving `Copy` for a struct containing a `String`. Read the compiler error.

> **Takeaway:** `#[derive(...)]` auto-implements routine traits from your fields
> — `Debug` for `{:?}` printing, `Clone`/`Copy` for duplication, `PartialEq`/
> `PartialOrd` for comparing and sorting, `Default` for a zero value — as long as
> every field supports the same trait.
