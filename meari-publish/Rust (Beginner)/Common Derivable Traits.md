---
created: "2026-07-08"
id: rust-b-derive
source: meari-course
study:
  answer: |
    #[derive(Debug, Clone, PartialEq)]
    struct Point {
        x: i32,
        y: i32,
    }
  kind: code
  lang: rust
  prompt: Add the right `#[derive(...)]` so `Point { x, y }` can be compared with `==`, cloned, and printed with `{:?}`.
  hint: |
    Don't implement these by hand — one `#[derive(...)]` line above the struct generates them. You need three traits: one for `==`, one for `{:?}`, and one for `.clone()`.
  starter: |
    struct Point {
        x: i32,
        y: i32,
    }
  tests:
    - 'assert_eq!(Point { x: 1, y: 2 }, Point { x: 1, y: 2 });'
    - 'assert_ne!(Point { x: 1, y: 2 }, Point { x: 3, y: 4 });'
    - 'assert_eq!(format!("{:?}", Point { x: 1, y: 2 }), "Point { x: 1, y: 2 }");'
    - 'let p = Point { x: 5, y: 6 }; assert_eq!(p.clone(), p);'
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
duplication operation via `.clone()`; `Copy` makes assignment duplicate the
value instead of moving it:

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

These enable `<`, `>`, and sorting. A generic function that compares values
would require its type to implement `PartialOrd`:

```rust
#[derive(PartialEq, Eq, PartialOrd, Ord)]
struct Score(u32);

let mut scores = vec![Score(30), Score(10), Score(20)];
scores.sort();          // works because Score implements Ord
```

## Later: `Default` — a sensible zero value

```rust
#[derive(Default)]
struct Config { verbose: bool, level: u32 }

let c = Config::default();   // Config { verbose: false, level: 0 }
```

## The derivable traits at a glance

| Derive        | Gives you                          | Example                    |
| ------------- | ---------------------------------- | -------------------------- |
| `Debug`       | `{:?}` / `{:#?}` printing           | `println!("{p:?}")`        |
| `Clone`       | explicit type-defined duplication   | `let q = p.clone();`       |
| `Copy`        | implicit bitwise copy on assignment | `let q = p;` (p still ok)  |
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

## Later: derive generates ordinary implementations

`#[derive(Debug)]` is an attribute attached to the next item. During compilation
Rust generates an implementation much like one you could write by hand. There
is no reflection or runtime switch involved.

Be careful with the phrase “deep copy.” `Clone` asks every field to clone itself;
what that means follows the field's implementation. For `String` it duplicates
the buffer, while cloning `Rc<T>` creates another shared owner rather than
duplicating `T`. `Copy` is more restrictive: it must be cheap and implicit, and
a type implementing `Drop` cannot also implement `Copy`.

Derive the capabilities your code needs. It is normal for a type to be `Debug`
and `PartialEq` but intentionally not `Clone` or `Copy`.

## Syntax checkpoint

```rust
#[derive(Debug, PartialEq)]
struct Point { x: i32, y: i32 }
```

`#[...]` is an **attribute** attached to the item directly below it. Inside
`derive`, commas separate the traits Rust should implement. Here `Debug` enables
`{:?}` formatting and `PartialEq` enables `==`; deriving does not alter the
fields or create a value.

You are ready to continue when you can select a derive for printing or equality
and understand that every field must support it. The next lesson uses closures
and iterators; their advanced type relationships are traits too, but a first
pass only needs their visible `|x| ...` syntax.

## Try it

1. **Match capability:** Choose the derive needed for `{:?}`, `==`, and
   assignment that leaves the source usable.
2. **Fill:** Add `Debug` and `PartialEq` to a struct, then print and compare it.
3. **Diagnose:** Try deriving `Copy` for a struct containing a `String` and
   explain the compiler error using the field requirement.
4. **Predict:** Explain why cloning an `Rc<T>` need not duplicate `T`.

> **Takeaway:** `#[derive(...)]` auto-implements routine traits from your fields
> — `Debug` for `{:?}` printing, `Clone`/`Copy` for explicit or implicit
> duplication, `PartialEq`/
> `PartialOrd` for comparing and sorting, `Default` for a zero value — as long as
> every field supports the same trait.
