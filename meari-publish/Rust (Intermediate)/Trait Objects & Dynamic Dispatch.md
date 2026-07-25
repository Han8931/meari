---
created: '2026-07-21'
id: rust-i-trait-objects
source: meari-course
study:
  answer: |
    trait Animal {
        fn legs(&self) -> u32;
    }

    struct Dog;
    struct Spider;
    struct Snake;

    impl Animal for Dog {
        fn legs(&self) -> u32 {
            4
        }
    }

    impl Animal for Spider {
        fn legs(&self) -> u32 {
            8
        }
    }

    impl Animal for Snake {
        fn legs(&self) -> u32 {
            0
        }
    }

    fn total_legs(animals: &[&dyn Animal]) -> u32 {
        animals.iter().map(|a| a.legs()).sum()
    }
  kind: code
  lang: rust
  prompt: 'Define a trait `Animal` with `fn legs(&self) -> u32`. Implement it for `Dog` (4), `Spider` (8), and `Snake` (0). Then write `total_legs(animals: &[&dyn Animal]) -> u32`, using borrowed trait objects rather than taking ownership.'
  starter: |
    trait Animal {
        fn legs(&self) -> u32;
    }

    struct Dog;
    struct Spider;
    struct Snake;

    // impl Animal for Dog, Spider, Snake

    fn total_legs(animals: &[&dyn Animal]) -> u32 {
        0
    }
  tests:
  - 'let dog = Dog; let spider = Spider; let snake = Snake; let zoo: Vec<&dyn Animal> = vec![&dog, &spider, &snake]; assert_eq!(total_legs(&zoo), 12);'
  - 'let a = Dog; let b = Dog; let pack: Vec<&dyn Animal> = vec![&a, &b]; assert_eq!(total_legs(&pack), 8);'
  - 'let empty: Vec<&dyn Animal> = vec![]; assert_eq!(total_legs(&empty), 0);'
subject: Rust (Intermediate)
title: Trait Objects & Dynamic Dispatch
---

[[rust-b-traits|Traits]] gave you shared behavior, and trait *bounds* (`fn f<T:
Trait>`) let one function work over many types. That's **static dispatch**: the
compiler stamps out a specialized copy per type, and the exact method is known
at compile time.

Sometimes you want the opposite — a single collection holding *different* types
that share a trait. For that you use a **trait object**: `dyn Trait`, usually
behind a pointer like `Box<dyn Trait>`.

## The problem trait objects solve

A `Vec<T>` holds one type. But "a list of shapes" or "a list of animals" mixes
types that only share a trait. `Box<dyn Animal>` says "some value on the heap
that implements `Animal`; I don't care which concrete type":

```rust
trait Animal {
    fn legs(&self) -> u32;
}

struct Dog;
struct Bird;
impl Animal for Dog  { fn legs(&self) -> u32 { 4 } }
impl Animal for Bird { fn legs(&self) -> u32 { 2 } }

let zoo: Vec<Box<dyn Animal>> = vec![Box::new(Dog), Box::new(Bird)];
```

## Dynamic dispatch

Calling `a.legs()` on a `Box<dyn Animal>` looks up the right method *at runtime*
through a small table of function pointers (a "vtable"). That's **dynamic
dispatch** — slightly slower than static, but it's what lets one loop walk a
mixed collection:

```rust
let total: u32 = zoo.iter().map(|a| a.legs()).sum();
```

Rule of thumb: reach for `dyn` when you truly need to mix types in one place;
prefer plain trait bounds otherwise.

The representation is usefully pictured as two pieces:

```
&dyn Animal
├── data pointer  ──► one concrete Dog or Bird value
└── vtable pointer ─► Animal methods for that concrete type
```

This is a mental model, not a layout guarantee to depend on in application
code. The important result is that the concrete type is erased from the
variable's static type while calls still reach its implementation.

## Why `dyn Trait` needs a pointer

Different implementors have different sizes, so a bare `dyn Animal` has no size
known at compile time. Put it behind a pointer whose size is known:

- `&dyn Animal` temporarily borrows an existing value.
- `Box<dyn Animal>` owns the value and is useful in an owning collection.
- `Arc<dyn Animal>` adds shared ownership when the trait's thread-safety bounds
  allow it. [[rust-i-arc-mutex|Shared State with Arc & Mutex]] introduces `Arc`
  later.

The pointer carries both a pointer to the value and the information needed to
find its trait methods.

## Not every trait is dyn-compatible

A trait used as `dyn Trait` must let Rust build one usable method table. As a
practical first rule, methods callable through the trait object cannot be generic
and cannot return a bare `Self`:

```rust
trait Factory {
    fn make<T>(&self) -> T;       // generic method
    fn duplicate(&self) -> Self;  // result size is unknown through `dyn Factory`
}
```

Such methods can sometimes be restricted with `where Self: Sized`, leaving the
rest of the trait usable as an object. When the compiler says a trait “is not dyn
compatible,” inspect its method signatures first.

Separate two questions that are often blurred:

1. **Does a concrete type implement the trait?**
2. **Can the trait be used behind `dyn`?**

A trait can be useful as a generic bound even when its method shapes prevent a
trait object.

## Your turn

1. **Choose dispatch:** Decide whether a homogeneous generic algorithm and a
   mixed animal collection need static or dynamic dispatch.
2. **Trace:** For `&dyn Animal`, identify the data pointer, method table, owner,
   and borrower.
3. **Diagnose:** Explain why a method returning bare `Self` is not callable
   through a value whose concrete size is erased.
4. **Implement:** Give `Dog`, `Spider`, and `Snake` leg counts, then total
   borrowed `&dyn Animal` values without allocating or taking ownership.
