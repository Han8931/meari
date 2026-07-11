---
created: "2026-07-08"
id: rust-b-generics
source: meari-course
study:
  answer: |
    fn largest<T: PartialOrd + Copy>(xs: &[T]) -> T {
        let mut max = xs[0];
        for &x in xs {
            if x > max {
                max = x;
            }
        }
        max
    }
  kind: code
  lang: rust
  prompt: 'Write a generic `largest<T: PartialOrd + Copy>(xs: &[T]) -> T` returning the maximum element. Assume `xs` is non-empty.'
  starter: |
    fn largest<T: PartialOrd + Copy>(xs: &[T]) -> T {
        xs[0]
    }
  tests:
    - assert_eq!(largest(&[1, 5, 3]), 5);
    - assert_eq!(largest(&[1.5, 2.5, 0.5]), 2.5);
    - assert_eq!(largest(&['a', 'c', 'b']), 'c');
subject: Rust (Beginner)
title: Generics
---

You'll see generics everywhere in Rust: `Vec<T>`, `Option<T>`, `Result<T, E>`
all have that `<T>`. **Generics** let you write one piece of code that works
over many types, instead of copying it once per type — with no loss of speed or
type safety.

## The problem generics solve

Say you want the largest element of a slice. Without generics you'd write it
once for integers, again for characters, again for… everything:

```rust
fn largest_i32(list: &[i32]) -> &i32 { /* ... */ }
fn largest_char(list: &[char]) -> &char { /* ... identical logic ... */ }
```

The bodies are identical; only the type differs. That's exactly the duplication
a generic erases.

## Generic functions

Introduce a **type parameter** `T` in angle brackets after the function name,
then use it like any type:

```rust
fn largest<T: PartialOrd>(list: &[T]) -> &T {
    let mut biggest = &list[0];
    for item in list {
        if item > biggest {     // needs T to be comparable — see the bound below
            biggest = item;
        }
    }
    biggest
}

let nums = vec![3, 7, 2, 9, 4];
let best = largest(&nums);           // T = i32
let letters = vec!['q', 'a', 'z'];
let top = largest(&letters);         // T = char
```

The `T: PartialOrd` part is a **trait bound** — it says "`T` can be any type,
*as long as* it can be compared with `>`." Without it, `item > biggest` wouldn't
compile, because not every type is orderable. Traits are the whole next lesson,
[[Traits]]; for now, read `<T: PartialOrd>` as "some comparable type `T`."

## Generic structs and enums

Types can be generic too. A `Point` that works for any coordinate type:

```rust
struct Point<T> {
    x: T,
    y: T,
}

let ints = Point { x: 1, y: 2 };        // Point<i32>
let floats = Point { x: 1.5, y: 2.5 };  // Point<f64>
```

Use several parameters when fields may differ:

```rust
struct Pair<T, U> {
    first: T,
    second: U,
}

let mixed = Pair { first: "age", second: 30 };  // Pair<&str, i32>
```

And you've been using generic **enums** all along — this is literally how the
standard library defines them:

```rust
enum Option<T> { Some(T), None }
enum Result<T, E> { Ok(T), Err(E) }
```

You'll study these two enums in detail in [[Option & Result]].

## Generic methods

Add the parameter to the `impl` block, then to the methods that use it:

```rust
impl<T> Point<T> {
    fn x(&self) -> &T {          // returns a reference to the x field
        &self.x
    }
}
```

## Zero-cost: monomorphization

Here's the payoff. Generics cost **nothing** at runtime. At compile time Rust
performs *monomorphization*: it stamps out a concrete copy of the generic code
for each type you actually use, exactly as if you'd hand-written them:

```
   generic:   fn largest<T>(list: &[T]) -> &T
                        │  you called it with i32 and char
          ┌─────────────┴─────────────┐
          ▼                           ▼
   fn largest_i32(&[i32])       fn largest_char(&[char])
        (concrete, inlined, as fast as bespoke code)
```

So generics give you the flexibility of "write once" with the performance of
"write by hand" — a recurring Rust theme.

## The same in Python

Python is dynamically typed, so it gets generics "for free" — a function just
works on whatever you pass:

```python
def largest(items):            # no type parameter needed
    biggest = items[0]
    for item in items:
        if item > biggest:
            biggest = item
    return biggest
```

But that freedom is unchecked: pass a list mixing incompatible types and it
crashes at *runtime* with a `TypeError`. Rust's `<T: PartialOrd>` proves, at
compile time, that whatever `T` you pick can actually be compared. Python's
`typing.TypeVar` adds generic *hints* for documentation and tooling, but nothing
enforces them the way Rust does.

## Read a generic signature aloud

Read this from left to right:

```rust
fn first<T>(items: &[T]) -> Option<&T>
```

“For any type `T`, `first` borrows a slice of `T` values and may return a
reference to one `T`.” Every occurrence of `T` must mean the same concrete type
for one call. Calling it with `&[i32]` makes all three `T`s mean `i32`.

A bound narrows “any type” to types with a capability:

```rust
fn print_twice<T: std::fmt::Display>(value: T) {
    println!("{value} {value}");
}
```

Without `Display`, the function body is not allowed to format an unknown `T`.
Bounds are promises available to the generic implementation, not merely
restrictions placed on callers.

## Try it

1. Write a generic `first<T>(items: &[T]) -> &T` function that returns `&items[0]`.
2. Call it with a `Vec<i32>` and a `Vec<&str>`.
3. Explain in plain English what the `T` stands for.

> **Takeaway:** generics (`<T>`) let one function, struct, or enum serve many
> types. Trait bounds like `<T: PartialOrd>` say what a type must be able to
> *do*, and monomorphization compiles it all down to concrete, zero-cost code.
