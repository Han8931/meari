---
created: "2026-07-08"
id: rust-b-traits
source: meari-course
subject: Rust (Beginner)
title: Traits
---

A **trait** defines shared behavior — a set of methods a type promises to
provide. If [[Generics]] answered "code over many types," traits answer "…that
all share some capability." They're Rust's version of interfaces, and they're
everywhere: `PartialOrd`, `Fn`, and `std::error::Error` are all traits you'll
see throughout Rust.

## Defining a trait

A trait is a named contract of method signatures:

```rust
trait Summary {
    fn summarize(&self) -> String;         // required — implementors must supply it

    fn preview(&self) -> String {          // default — implementors MAY override it
        String::from("(read more)")
    }
}
```

## Implementing a trait

Use `impl Trait for Type` to fulfil the contract for a specific type:

```rust
struct Article { title: String, body: String }
struct Tweet   { user: String, text: String }

impl Summary for Article {
    fn summarize(&self) -> String {
        format!("{}: {}", self.title, self.body)
    }
}

impl Summary for Tweet {
    fn summarize(&self) -> String {
        format!("@{}: {}", self.user, self.text)
    }
    fn preview(&self) -> String {          // override the default
        format!("@{}...", self.user)
    }
}

let a = Article { title: "Rust".into(), body: "is fast".into() };
println!("{}", a.summarize());   // "Rust: is fast"
println!("{}", a.preview());     // "(read more)" — uses the default
```

Two different types, one shared vocabulary. Any code that works with the
`Summary` trait can now handle both.

## Traits as bounds

This is where traits and generics meet. A trait bound restricts a generic to
types that implement the trait — so inside the function you may call the trait's
methods:

```rust
fn notify<T: Summary>(item: &T) {
    println!("Breaking! {}", item.summarize());
}

// identical, using the `impl Trait` shorthand:
fn notify(item: &impl Summary) {
    println!("Breaking! {}", item.summarize());
}
```

That `<T: PartialOrd>` from the generics lesson was exactly this pattern —
`PartialOrd` is just a trait for "can be ordered."

## Static vs dynamic dispatch

There are two ways to be generic over a trait, and the difference matters:

```rust
// STATIC dispatch: one type per call site, resolved at compile time (fast)
fn print_it(item: &impl Summary) { println!("{}", item.summarize()); }

// DYNAMIC dispatch: a trait OBJECT — a mixed bag of types behind one pointer
let feed: Vec<Box<dyn Summary>> = vec![
    Box::new(a),
    Box::new(Tweet { user: "bo".into(), text: "hi".into() }),
];
for item in &feed {
    println!("{}", item.summarize());   // looked up at runtime via a vtable
}
```

A trait object (`dyn Summary`) carries a hidden pointer to a **vtable** — a
table of the type's method implementations, consulted at runtime:

```
   Box<dyn Summary> ──► ┌────────┐
                        │ data   │  the Article/Tweet value
                        │ vptr ──┼──► vtable: summarize(), preview()
                        └────────┘
```

| Approach                       | Dispatch            | Cost              | Use when                          |
| ------------------------------ | ------------------- | ----------------- | --------------------------------- |
| `<T: Trait>` / `impl Trait`    | static (compile)    | zero-cost, inlined| the type is known at compile time |
| `dyn Trait` (trait object)     | dynamic (runtime)   | small indirection | you need a *mix* of types together|

Reach for generics by default; reach for `dyn Trait` when you genuinely need a
heterogeneous collection like that `Vec<Box<dyn Summary>>`.

## You've already met traits

- `Box<dyn std::error::Error>` in [[Error Propagation & Panics]] is a trait
  object — any error type behind one pointer.
- `Fn`, `FnMut`, `FnOnce` in [[Closures & Iterators]] are traits closures
  implement.
- `PartialOrd` in [[Generics]] is the "can be compared" trait.

One rule to know: you can `impl` a trait for a type only if **you define the
trait, or you define the type** (the "orphan rule"). It stops two crates from
adding conflicting implementations for someone else's types.

## The same in Python

Python approximates traits two ways. **Duck typing** — if it has the method, it
works, checked at runtime:

```python
def notify(item):                 # no declared contract
    print(item.summarize())       # works on anything with .summarize()
```

Or, more explicitly, an **abstract base class** (or `typing.Protocol`), which is
the closest analog to a trait:

```python
from abc import ABC, abstractmethod

class Summary(ABC):               # ~ trait Summary
    @abstractmethod
    def summarize(self) -> str: ...

    def preview(self) -> str:     # ~ a default method
        return "(read more)"
```

The difference is enforcement and timing. Python checks "does it have this
method?" when the call runs; Rust checks that the type implements the trait when
the program *compiles*. Rust's `dyn Trait` is the nearest thing to Python's
"pass any object that quacks right."

## Try it

1. Define a trait `Named` with a method `name(&self) -> &str`.
2. Implement it for a simple `User` struct.
3. Write a function that takes `&impl Named` and prints the name.

> **Takeaway:** a trait is a contract of shared behavior; `impl Trait for Type`
> fulfils it, and default methods cut boilerplate. Use trait *bounds* to make
> generics callable, static dispatch (`impl Trait`) for speed, and `dyn Trait`
> when you need a mixed collection of types.
