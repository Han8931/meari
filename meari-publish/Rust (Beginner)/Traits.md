---
created: "2026-07-08"
id: rust-b-traits
source: meari-course
study:
  answer: |
    trait Animal {
        fn sound(&self) -> String;
    }

    struct Dog;
    struct Cat;

    impl Animal for Dog {
        fn sound(&self) -> String {
            "woof".to_string()
        }
    }

    impl Animal for Cat {
        fn sound(&self) -> String {
            "meow".to_string()
        }
    }
  kind: code
  lang: rust
  prompt: The trait, structs, and `impl` blocks are already provided. Replace each `String::new()` so `Dog.sound()` returns `"woof"` and `Cat.sound()` returns `"meow"`.
  starter: |
    trait Animal {
        fn sound(&self) -> String;
    }

    struct Dog;
    struct Cat;

    impl Animal for Dog {
        fn sound(&self) -> String {
            String::new()
        }
    }

    impl Animal for Cat {
        fn sound(&self) -> String {
            String::new()
        }
    }
  tests:
    - assert_eq!(Dog.sound(), "woof");
    - assert_eq!(Cat.sound(), "meow");
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

## Static dispatch: the beginner default

With `<T: Summary>` or `&impl Summary`, Rust knows each call's concrete type at
compile time and can call the implementation directly. This is **static
dispatch**. It is the ordinary starting point; use it unless you have a reason
to store different concrete types together.

## Optional preview: dynamic dispatch for mixed collections

You do not need this section to define or use ordinary traits. Read it only
when you need one collection to hold several different concrete types. There
are two ways to be generic over a trait:

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

Do not try to memorize the representation. The decision is:

```
one concrete type per call, known at compile time → &impl Trait / <T: Trait>
different concrete types in one collection       → Box<dyn Trait>
```

Reach for generics by default. Treat `dyn Trait` as a second-stage tool for a
genuine heterogeneous collection such as `Vec<Box<dyn Summary>>`.

## Traits across this course

- `Box<dyn std::error::Error>` in [[Error Propagation & Panics]] is a trait
  object — any error type behind one pointer.
- `Fn`, `FnMut`, and `FnOnce` in the next lesson,
  [[Closures & Iterators]], are traits implemented by closures.
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

## Trait, implementation, and bound

These three roles are easy to blur:

1. `trait Named { ... }` defines a capability.
2. `impl Named for User { ... }` teaches one concrete type that capability.
3. `fn show(x: &impl Named)` accepts any value whose type has that capability.

The trait does not store data and implementing it does not create a new object.
It is a compile-time relationship between a behavior and a type.

Start with `impl Trait` or `<T: Trait>`; both use static dispatch and are usually
the simplest choice. Treat `dyn Trait` as a separate, later tool for cases where
the concrete types genuinely must differ at runtime.

## Syntax checkpoint

Keep the three roles on separate lines:

```rust
trait Named { fn name(&self) -> &str; } // define a capability
impl Named for User { /* method */ }    // give User that capability
fn show(x: &impl Named) { /* use it */ }// accept any capable type
```

Read `impl Named for User` as “implement the `Named` contract for `User`.” Read
`&impl Named` as “borrow a value of some type that implements `Named`.” You do
not need `dyn`, boxes, or vtables to use ordinary trait bounds.

You are ready to continue when you can define one trait, implement it for one
struct, and call the method. The next lesson shows a special case where the
compiler can generate common implementations from a type's fields.

## Try it

1. **Classify:** In the `Summary` example, identify the trait, two implementing
   types, a required method, and a default method.
2. **Fill:** Define `Named` with `name(&self) -> &str` and implement it for
   `User`.
3. **Use a bound:** Write a function taking `&impl Named` and print the name.
4. **Explain:** State why implementing a trait does not create a new object.
5. **Optional:** Explain why a mixed `Article`/`Tweet` vector needs `dyn Summary`
   while a single call to `notify` does not.

> **Takeaway:** a trait is a contract of shared behavior; `impl Trait for Type`
> fulfils it, and default methods cut boilerplate. Use trait *bounds* to make
> generics callable, static dispatch (`impl Trait`) for speed, and `dyn Trait`
> when you need a mixed collection of types.
