---
created: "2026-07-08"
id: rust-b-generics
source: meari-course
study:
  answer: |
    fn choose_first<T>(first: T, _second: T) -> T {
        first
    }
  kind: code
  lang: rust
  prompt: 'Complete `choose_first<T>(first: T, _second: T) -> T` so it returns the first value. Both arguments must have the same type, but the function must work for any type.'
  starter: |
    fn choose_first<T>(first: T, _second: T) -> T {
        todo!()
    }
  tests:
    - assert_eq!(choose_first(1, 2), 1);
    - assert_eq!(choose_first("left", "right"), "left");
    - assert_eq!(choose_first(String::from("a"), String::from("b")), "a");
subject: Rust (Beginner)
title: Generics
---

You'll see generics everywhere in Rust: `Vec<T>`, `Option<T>`, `Result<T, E>`
all have that `<T>`. **Generics** let you write one piece of code that works
over many types, instead of copying it once per type — with no loss of speed or
type safety.

## The problem generics solve

Say you want a function that returns the first item in a slice. Without generics
you would write it once for integers, again for characters, and again for every
other type:

```rust
fn first_i32(items: &[i32]) -> Option<&i32> { items.first() }
fn first_char(items: &[char]) -> Option<&char> { items.first() }
```

The bodies are identical; only the type differs. That's exactly the duplication
a generic erases.

## Generic functions

Introduce a **type parameter** `T` in angle brackets after the function name,
then use it like any type:

```rust
fn first<T>(items: &[T]) -> Option<&T> {
    items.first()
}

let nums = vec![3, 7, 2];
let number = first(&nums);       // T = i32; result is Option<&i32>

let letters = vec!['q', 'a'];
let letter = first(&letters);    // T = char; result is Option<&char>
```

`T` is a placeholder for one concrete type chosen at each call. In the first
call every `T` becomes `i32`; in the second every `T` becomes `char`. The
function body does not need to know which one because `.first()` works without
performing any type-specific operation on the element.

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
   generic:   fn first<T>(items: &[T]) -> Option<&T>
                        │  you called it with i32 and char
          ┌─────────────┴─────────────┐
          ▼                           ▼
   first for &[i32]             first for &[char]
        (concrete, inlined, as fast as bespoke code)
```

So generics give you the flexibility of "write once" with the performance of
"write by hand" — a recurring Rust theme.

## The same in Python

Python is dynamically typed, so it gets generics "for free" — a function just
works on whatever you pass:

```python
def first(items):              # no type parameter needed
    return items[0] if items else None
```

Python's freedom is checked dynamically. Rust checks at compile time that one
call uses one consistent `T`. Python's `typing.TypeVar` can express a similar
relationship for documentation and external type checkers, but the Python
runtime does not enforce it.

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
restrictions placed on callers. `Display` is a **trait** and this requirement is
a **trait bound**. [[Traits]] explains how to define capabilities and use bounds
in detail.

## Try it

1. **Read aloud:** Explain each occurrence of `T` in
   `fn first<T>(items: &[T]) -> Option<&T>`.
2. **Fill:** Write that `first` function using `.first()`.
3. **Compare calls:** Call it with a `Vec<i32>` and a `Vec<&str>`, naming the
   concrete `T` in each call.
4. **Diagnose:** Try `choose_first(1, "two")` and explain why one call cannot
   choose two different meanings for `T`.

> **Takeaway:** generics (`<T>`) let one function, struct, or enum serve many
> types. Every `T` in one call represents the same concrete type. Trait bounds
> add capabilities that the generic body may use, and monomorphization compiles
> generic code into concrete implementations.
