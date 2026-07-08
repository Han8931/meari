---
created: "2026-07-08"
id: rust-b-ownership
source: meari-course
subject: Rust (Beginner)
title: Ownership & Moves
---

This is *the* lesson. Ownership is the idea that makes Rust Rust — it's how the
language guarantees memory safety with **no garbage collector and no manual
`free`**. Everything else you've learned was warm-up for this.

## The three rules

```
  1. Every value has exactly ONE owner.
  2. There can only be one owner at a time.
  3. When the owner goes out of scope, the value is DROPPED (freed).
```

That's the entire model. The compiler enforces it, so cleanup is automatic *and*
provably correct.

## Scope and `drop`

When a value's owner leaves its block, Rust automatically frees the value — no
`free()`, no `delete`, no GC:

```rust
fn main() {
    let s = String::from("hello"); // s owns a heap-allocated string
    println!("{s}");
}                                  // s goes out of scope → memory freed here
```

## Stack vs heap (why moves exist)

Simple values (`i32`, `bool`, `char`) sit entirely on the **stack** — cheap to
copy. A `String` or `Vec` is different: a small handle on the stack points to
data on the **heap**.

```
   let s = String::from("hi");

   STACK              HEAP
   ┌───────────┐      ┌───┬───┐
   │ ptr   ●───┼────► │ h │ i │
   │ len   2   │      └───┴───┘
   │ cap   2   │
   └───────────┘
```

## A move transfers ownership

Assigning a heap value to another variable does **not** copy the heap data — it
*moves* the handle, and the old variable becomes invalid:

```rust
let s1 = String::from("hello");
let s2 = s1;            // the handle MOVES from s1 to s2

println!("{s2}");      // ✅ fine
println!("{s1}");      // ❌ error: value borrowed after move
```

```
   before:  s1 ●──► "hello"        after `let s2 = s1;`

   s1 ✗ (invalidated)
   s2 ●──► "hello"       ← only ONE owner, as rule #2 demands
```

Why invalidate `s1`? If both `s1` and `s2` pointed at the same heap buffer, then
when *both* went out of scope Rust would try to free it **twice** — a classic
double-free bug. Moving prevents that by construction.

## The Python contrast — the biggest one in this course

This is where Rust departs most sharply from Python. In Python, assignment just
creates another *name* for the **same** object; both names stay valid, and a
garbage collector frees the object later when nothing references it:

```python
s1 = "hello"
s2 = s1          # s1 and s2 are two names for one object
print(s1, s2)    # ✅ both fine — Python never "moves" or invalidates a name
```

Python can afford this because a garbage collector is always running to decide
when memory is safe to free. Rust has no GC, so it uses ownership instead: the
move *invalidates* `s1` precisely so the heap buffer has exactly one owner and
gets freed exactly once — no GC required, no double-free possible.

## `Copy` types don't move

Types that live entirely on the stack implement the `Copy` trait, so assignment
duplicates them and the original stays valid:

```rust
let x = 5;
let y = x;             // x is COPIED, not moved
println!("{x} {y}");   // ✅ both usable — 5 5
```

| Behavior on assign | Types                                         |
| ------------------ | --------------------------------------------- |
| **Copy** (cheap)   | `i32`, `f64`, `bool`, `char`, tuples of Copy  |
| **Move** (handle)  | `String`, `Vec<T>`, most heap-owning types    |

## `clone()` for an explicit deep copy

When you genuinely want two independent owners of heap data, ask for it
explicitly with `.clone()` — the visible `.clone()` call is Rust telling you
"this costs a heap copy":

```rust
let s1 = String::from("hello");
let s2 = s1.clone();   // deep copy — s1 and s2 each own their own buffer
println!("{s1} {s2}"); // ✅ both valid
```

## Moves happen at function boundaries too

Passing a value to a function moves it in, unless the type is `Copy`:

```rust
fn takes(s: String) { println!("{s}"); }  // s is dropped when this returns

let name = String::from("Ana");
takes(name);
// println!("{name}");   // ❌ name was moved into takes()
```

Constantly moving values into functions and moving them back out would be
miserable. That's exactly the problem **borrowing** solves — lending a value
without giving up ownership — which is the very next lesson,
[[References & Borrowing]].

## Try it

1. Move a `String` from `s1` to `s2`, then try to print both. Read the error.
2. Repeat the same experiment with an `i32`. Why does it still work?
3. Use `.clone()` on a `String` and confirm both variables remain usable.

> **Takeaway:** one owner, dropped at end of scope. Assignment *moves* heap
> types (invalidating the source) and *copies* stack types. Reach for `.clone()`
> only when you truly need a second independent owner.
