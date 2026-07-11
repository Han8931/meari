---
created: "2026-07-08"
id: rust-b-smart-pointers
source: meari-course
study:
  answer: |
    fn boxed_value(value: i32) -> Box<i32> {
        Box::new(value)
    }
  kind: code
  lang: rust
  prompt: 'Write `boxed_value(value: i32) -> Box<i32>` that stores the given value in a `Box` and returns it.'
  starter: |
    fn boxed_value(value: i32) -> Box<i32> {
        Box::new(0)
    }
  tests:
    - assert_eq!(*boxed_value(7), 7);
    - assert_eq!(*boxed_value(-3), -3);
subject: Rust (Beginner)
title: Box, Rc & RefCell
---

So far ownership has been strict: one owner, known at compile time. Sometimes
that's too rigid — you need heap allocation for a value of unknown size, or
several parts of a program need to share the *same* data. **Smart pointers** are
types that own heap data and add extra powers on top. Three carry you a long way.

> You've already *used* `Box` — `Box<dyn Error>` in [[Error Propagation & Panics]]
> and `Box<dyn Summary>` in [[Traits]] — so treat the `Box` section below as core;
> it just names what you've been using. `Rc`, and especially `Rc<RefCell<T>>`, are
> the genuinely advanced, optional material: grasp what each is *for* on a first
> pass and let the finer details settle later. Everyday beginner Rust rarely needs
> `Rc`/`RefCell`.

## Recap: stack vs heap

```
   STACK: fast, fixed-size, auto-managed (function-local values)
   HEAP:  flexible size, lives as long as an owner keeps it
          — you reach it through a pointer on the stack
```

## `Box<T>`: put a value on the heap

A `Box` is the simplest smart pointer: it owns a single value stored on the
heap. Two classic uses:

```rust
// 1. Explicitly heap-allocate a value
let boxed: Box<i32> = Box::new(5);
println!("{}", *boxed);       // deref with * to reach the value

// 2. Enable a RECURSIVE type (which would otherwise be infinitely sized)
enum List {
    Cons(i32, Box<List>),     // Box gives the variant a known, pointer size
    Nil,
}
```

Without the `Box`, `List` would contain a `List` directly — infinite size. The
`Box` is a fixed-size pointer, breaking the cycle.

```
   Box<List>:   stack ┌─────┐      heap ┌──────────────┐
                      │ ●───┼────────►  │ Cons(1, ●───) │──► …
                      └─────┘           └──────────────┘
```

## `Rc<T>`: multiple owners (single-threaded)

Ownership says *one* owner — but sometimes several parts genuinely need to own
the same data (think a node referenced by many others). `Rc<T>` — "reference
counted" — allows it by counting owners and freeing the data only when the last
one goes away:

```rust
use std::rc::Rc;

let a = Rc::new(String::from("shared"));
let b = Rc::clone(&a);        // NOT a deep copy — just bumps the count to 2
let c = Rc::clone(&a);        // count is now 3

println!("owners: {}", Rc::strong_count(&a));  // 3
// data freed only when a, b, AND c are all dropped
```

```
        a ●─┐
        b ●─┼──►  "shared"   (refcount = 3)
        c ●─┘
```

`Rc::clone` is cheap — it copies a pointer and increments a counter, unlike
`.clone()` on the inner data. `Rc` is **single-threaded only**; its thread-safe
sibling is `Arc`.

## The Python contrast

In Python this machinery is entirely invisible: *every* object is already
heap-allocated and reference-counted, and any shared object can be mutated
freely:

```python
a = ["shared"]
b = a                 # both names share one object — like Rc, but automatic
b.append("more")
print(a)              # ['shared', 'more'] — the mutation shows through both
```

Python's runtime does the ref-counting (that's `Rc`) and always permits shared
mutation (that's `RefCell`), so you never think about it. Rust has no such
runtime, so it makes you *opt in* to each capability with `Box`, `Rc`, and
`RefCell` — the trade-off being that the cost is always visible in the code.

## `RefCell<T>`: interior mutability

`Rc` gives shared ownership, but it only hands out *immutable* access — and the
borrow rules forbid mutating shared data. `RefCell<T>` bends this by moving the
borrow check from **compile time to runtime**: you can mutate through a shared
reference, and if you break the "one writer" rule, it **panics** at runtime
instead of refusing to compile.

```rust
use std::cell::RefCell;

let cell = RefCell::new(5);
*cell.borrow_mut() += 10;         // mutate through a shared reference
println!("{}", cell.borrow());    // 15

// breaking the rule panics at RUNTIME:
let a = cell.borrow_mut();
let b = cell.borrow_mut();        // 💥 panic: already mutably borrowed
```

## Combining shared ownership and mutation: `Rc<RefCell<T>>`

Put them together and you get data with **multiple owners that can also be
mutated** — the standard beginner pattern for shared, editable structures like
graph or tree nodes:

```rust
use std::rc::Rc;
use std::cell::RefCell;

let shared = Rc::new(RefCell::new(vec![1, 2, 3]));
let clone  = Rc::clone(&shared);

clone.borrow_mut().push(4);       // mutate via one owner…
println!("{:?}", shared.borrow()); // …visible through the other: [1, 2, 3, 4]
```

## Which one when?

| Type          | Gives you                    | Borrow check | Threads |
| ------------- | ---------------------------- | ------------ | ------- |
| `Box<T>`      | heap allocation, one owner   | compile time | any     |
| `Rc<T>`       | many owners, read-only       | compile time | single  |
| `RefCell<T>`  | mutate through shared `&`    | **runtime**  | single  |
| `Rc<RefCell<T>>` | shared **and** mutable     | runtime      | single  |

One caveat: `Rc` cycles (A owns B owns A) leak memory, because the count never
reaches zero. The fix is `Weak<T>`, a non-owning reference — worth knowing exists,
though the details are beyond the beginner track.

## Why this lesson is optional

These wrappers are not upgrades that should replace ordinary ownership. Most
Rust programs should begin with owned values and normal `&`/`&mut` borrows:

```
one clear owner?                 use T
temporary shared access?        use &T / &mut T
must allocate behind a pointer? use Box<T>
truly needs several owners?     use Rc<T>
shared owners must mutate?      consider Rc<RefCell<T>> carefully
```

Each step adds flexibility and cost. `Rc::clone` is cheap because it increments
a count, but shared ownership can make cleanup harder to reason about. `RefCell`
moves an error from compile time to a possible runtime panic. If the compiler
rejects an ordinary borrow, first reconsider the data design; do not automatically
wrap everything in `Rc<RefCell<_>>`.

## Try it

1. Put an integer in a `Box` with `Box::new(5)` and print it with `*boxed`.
2. Create an `Rc<String>`, clone it with `Rc::clone`, and print `Rc::strong_count`.
3. Read the `RefCell` example, but do not worry if it feels advanced on the first pass.

> **Takeaway:** reach for `Box` to heap-allocate, `Rc` to share ownership, and
> `RefCell` to mutate shared data (accepting a runtime borrow check). `Rc<RefCell<T>>`
> combines the last two — the go-to for shared mutable structures.
