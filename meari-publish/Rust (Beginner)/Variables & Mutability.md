---
created: "2026-07-08"
id: rust-b-variables
source: meari-course
study:
  answer: |
    fn update_score() -> i32 {
        let mut score = 10;
        score = 15;
        let score = score + 5;
        score
    }
  kind: code
  lang: rust
  prompt: Complete `update_score()`: declare mutable `score` as 10, change it to 15, then shadow it with `score + 5` and return the result.
  starter: |
    fn update_score() -> i32 {
        0
    }
  tests:
    - assert_eq!(update_score(), 20);
subject: Rust (Beginner)
title: Variables & Mutability
---

In most languages a variable is a box you can freely overwrite. In Rust the
default is the opposite: **variables are immutable unless you say otherwise.**
This isn't a nuisance — it's a design choice that makes code easier to reason
about and unlocks safe concurrency later.

## `let` and the `mut` opt-in

```rust
let x = 5;
x = 6;        // ❌ error: cannot assign twice to immutable variable `x`

let mut y = 5;
y = 6;        // ✅ fine — y was declared mutable
```

If you try to mutate an immutable binding, the compiler stops you *and* suggests
adding `mut`. The rule of thumb: reach for `mut` only when you genuinely need to
change a value in place.

## Shadowing is not mutation

You can declare a new variable with the same name as an old one. This
**shadows** the previous binding — it's a brand-new variable that happens to
reuse the name, not a mutation of the old one:

```rust
let spaces = "   ";        // spaces is a &str
let spaces = spaces.len(); // spaces is now a usize (3)
```

Notice the *type changed* — that's allowed with shadowing but impossible with
`mut` (a `mut` variable keeps its type forever). Shadowing shines when you
transform a value through a pipeline and want to keep one clear name.

```
  mut  :  same binding, value changes,  type FIXED
          let mut n = 5;  n = 6;        // still i32

shadow :  new binding each time,        type may CHANGE
          let n = "5";  let n = 5;      // &str → i32
```

## `const` — compile-time constants

```rust
const MAX_POINTS: u32 = 100_000;
```

A `const` differs from an immutable `let` in three ways:

| Feature            | `let` (immutable)        | `const`                       |
| ------------------ | ------------------------ | ----------------------------- |
| Value known at     | Runtime is fine          | **Must** be compile-time      |
| Type annotation    | Optional (inferred)      | **Required**                  |
| Scope              | Block-local              | Any scope, incl. global       |
| Naming convention  | `snake_case`             | `SCREAMING_SNAKE_CASE`        |

Use `const` for fixed facts about your program — a maximum, a conversion factor,
a version number.

## The Python contrast

Python has no immutable-by-default binding — every variable is freely
reassignable, so Rust's `mut` distinction is brand new:

```python
x = 5
x = 6          # always fine in Python; no `mut` needed or possible
```

The nearest thing to a `const` is a `SCREAMING_SNAKE_CASE` name — but that's a
convention only, and nothing stops you reassigning it. Rust makes immutability a
rule the compiler enforces.

## Scope: variables live in blocks

A binding is valid from where it's declared until the end of its enclosing
`{ }` block, then it's gone:

```rust
fn main() {
    let outer = 1;
    {
        let inner = 2;
        println!("{outer} {inner}"); // both visible
    }
    // inner is gone here
    println!("{outer}");             // only outer
}
```

This scope-based lifetime is the seed of Rust's whole memory model — a value is
cleaned up when its owner goes out of scope. You'll see that idea again, much
more powerfully, in [[Ownership & Moves]].

## Binding, value, and type are different ideas

In `let mut score: i32 = 10`, `score` is the **binding** (the name), `10` is the
**value**, and `i32` is the **type**. `mut` applies to the binding: it allows that
name to receive another `i32` value. It does not make Rust dynamically typed.

```rust
let mut score: i32 = 10;
score = 11;          // same binding, same type, new value
// score = "high";  // error: expected i32, found &str
```

Type inference only means Rust can fill in an omitted annotation. The inferred
type still cannot change. Shadowing can also happen inside a smaller scope:

```rust
let message = "outer";
{
    let message = "inner";
    println!("{message}"); // inner: nearest declaration wins
}
println!("{message}");     // outer is visible again
```

## Try it

1. Write `let x = 5; x = 6;` and read the compiler error. Then fix it with `mut`.
2. Shadow a variable so it changes from a string to a number: `"42"` → `42`.
3. Create a `const` for a maximum score or limit.

> **Takeaway:** immutable-by-default plus shadowing lets you write transformation
> pipelines that are both safe and readable. Add `mut` deliberately, not by
> reflex.
