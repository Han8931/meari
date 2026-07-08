---
created: "2026-07-08"
id: rust-b-modules
source: meari-course
subject: Rust (Beginner)
title: Modules & Cargo Dependencies
---

As a program grows past a single file, you need a way to organize code into
namespaces and to pull in libraries other people wrote. Rust's **module system**
handles the first; **Cargo** and crates.io handle the second. This final lesson
ties the toolchain from [[Hello, Cargo]] back together.

## The vocabulary

```
   PACKAGE   what `cargo new` makes; has a Cargo.toml
     └── CRATE   a compilation unit; a binary (has main) or a library
           └── MODULE   a namespace inside a crate (mod ...)
```

- **Package** — a bundle Cargo builds; described by `Cargo.toml`.
- **Crate** — the tree of code compiled together. A *binary* crate has a
  `main`; a *library* crate is meant to be used by others.
- **Module** — an in-crate namespace you create with `mod`.

## Defining modules

Modules group related items and, crucially, control **visibility**. Everything
is **private by default** — you opt into exposure with `pub`:

```rust
mod math {
    pub fn add(a: i32, b: i32) -> i32 {   // pub → visible outside `math`
        a + b
    }

    fn secret() -> i32 { 42 }             // private to `math`

    pub mod trig {                        // nested module
        pub fn sine(x: f64) -> f64 { x.sin() }
    }
}

fn main() {
    let s = math::add(2, 3);              // path with ::
    let y = math::trig::sine(0.0);
    // math::secret();                    // ❌ private — won't compile
}
```

| Visibility     | Reachable from                          |
| -------------- | --------------------------------------- |
| *(default)*    | the current module and its children     |
| `pub`          | anywhere the parent is reachable        |
| `pub(crate)`   | anywhere in this crate, but not outside |

## Modules across files

Once a module grows, move it to its own file. `mod foo;` (note the semicolon,
not a block) tells Rust to load `foo` from a file:

```
   src/
   ├── main.rs        contains `mod math;`
   ├── math.rs        the contents of module `math`
   └── shapes/        a module with submodules
       └── mod.rs     (or, newer style: shapes.rs beside a shapes/ folder)
```

```rust
// main.rs
mod math;                 // load src/math.rs as module `math`

fn main() {
    println!("{}", math::add(2, 3));
}
```

## `use`: shortening paths

Typing full paths everywhere is tedious. `use` brings a name into scope:

```rust
use std::collections::HashMap;   // now write HashMap, not the full path
use math::trig::sine;

// path keywords you'll see:
use crate::math::add;   // crate  = this crate's root
use self::helpers;      // self   = the current module
use super::config;      // super  = the parent module
```

## Adding external crates

Rust's real leverage comes from **crates.io**, the package registry. To use a
library, add it as a dependency. The easy way:

```bash
cargo add rand          # edits Cargo.toml for you
```

…which adds a line to `Cargo.toml`:

```toml
[dependencies]
rand = "0.8"
```

Then `use` it in your code just like a standard-library module:

```rust
use rand::Rng;

fn main() {
    let n = rand::thread_rng().gen_range(1..=6);   // roll a die
    println!("rolled a {n}");
}
```

`cargo build` downloads the crate (and its dependencies), records exact versions
in `Cargo.lock`, and compiles everything together.

## The same in Python

Python's module system is the closest analog you already know, so this table is
a handy Rosetta stone:

| Rust                        | Python                                     |
| --------------------------- | ------------------------------------------ |
| `mod math;` (loads a file)  | `import math` (a `math.py` module)         |
| `use math::add;`            | `from math import add`                     |
| `pub fn` vs private default | public `def` vs `_name` "private"          |
| `crate` / `super` / `self`  | package / relative imports (`from . import`) |
| `cargo add rand`            | `pip install ...`                          |

The sharpest contrast is visibility: Rust items are **private by default** and
the compiler enforces it, whereas Python's privacy (a leading `_`) is purely a
naming convention that nothing prevents you from ignoring.

## Where to go next

You've covered the beginner core: the toolchain, the type system, the ownership
model that defines Rust, generics and traits, error handling, the key
collections, and how to organize and extend a project. One optional lesson
remains — [[Box, Rc & RefCell]] — for when you need heap allocation or shared
ownership. From here the natural next steps are **lifetimes** in depth, **writing
tests** with `#[test]`, and **concurrency** — where the borrow rules you learned
pay off as fearless parallelism.

## Try it

1. Create a `math` module with a public `add` function.
2. Call it from `main` using `math::add(2, 3)`.
3. Run `cargo fmt` after editing your code.

> **Takeaway:** modules (`mod`, `pub`, `use`) organize code into namespaces with
> private-by-default visibility, while `cargo add` plus `use` pulls in the vast
> ecosystem on crates.io. Package → crate → module is the hierarchy that holds
> it all together.
