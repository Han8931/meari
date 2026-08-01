---
created: "2026-07-08"
id: rust-b-question
source: meari-course
study:
  answer: |
    fn parse_sum(a: &str, b: &str) -> Result<i32, std::num::ParseIntError> {
        let x: i32 = a.parse()?;
        let y: i32 = b.parse()?;
        Ok(x + y)
    }
  kind: code
  lang: rust
  prompt: 'Write `parse_sum(a: &str, b: &str) -> Result<i32, std::num::ParseIntError>` that parses both strings to `i32` with the `?` operator and returns their sum.'
  starter: |
    fn parse_sum(a: &str, b: &str) -> Result<i32, std::num::ParseIntError> {
        Ok(0)
    }
  tests:
    - assert_eq!(parse_sum("2", "3"), Ok(5));
    - assert_eq!(parse_sum("10", "-4"), Ok(6));
    - assert!(parse_sum("x", "3").is_err());
subject: Rust (Beginner)
title: Error Propagation & Panics
---

[[Option & Result]] gave us failure-as-a-value. But real programs call functions
that call functions — and an error often needs to travel *up* several layers to
whoever can actually deal with it. This lesson is about moving errors around
cleanly, and about the escape hatch when an error truly can't be recovered:
`panic!`.

## The problem: manual propagation is noisy

Without help, passing an error upward means unwrapping and re-returning at every
step. Start with two operations that have the same error type: both parse text
as an `i32`.

```rust
fn parse_sum(a: &str, b: &str) -> Result<i32, std::num::ParseIntError> {
    let left = match a.parse::<i32>() {
        Ok(number) => number,
        Err(error) => return Err(error),
    };

    let right = match b.parse::<i32>() {
        Ok(number) => number,
        Err(error) => return Err(error),
    };

    Ok(left + right)
}
```

Read `Result<i32, ParseIntError>` as “either an `i32` result, or a parsing
error.” Each `match` does the same plumbing: extract the success value or return
the error unchanged.

## The `?` operator

`?` collapses that pattern into one character. On a `Result`, it means:
**if `Ok`, unwrap the value and keep going; if `Err`, return that error from the
whole function immediately.**

```rust
fn parse_sum(a: &str, b: &str) -> Result<i32, std::num::ParseIntError> {
    let left: i32 = a.parse()?;   // Err → return it; Ok → bind the number
    let right: i32 = b.parse()?;  // same flow again
    Ok(left + right)
}
```

```
   a.parse()?
        │
        ├── Ok(number) ──► bind number to left, continue
        └── Err(error) ──► return Err(error) from parse_sum
```

Two conditions to use `?`:

1. The enclosing function must itself return a `Result` (or `Option`) — `?`
   needs somewhere to return the error *to*.
2. The error returned by the operation must fit the error type declared by the
   enclosing function. Here both are `ParseIntError`, so no conversion is
   needed.

`?` also works on `Option`, returning `None` early.

## `?` in `main`

`main` can return a `Result`, letting you use `?` at the top level. This example
performs only file I/O, so one concrete error type is enough:

```rust
fn main() -> std::io::Result<()> {
    let text = std::fs::read_to_string("message.txt")?;
    println!("{text}");
    Ok(())
}
```

Read `Result<(), std::io::Error>` as “success carries no interesting value;
failure carries an I/O error.” `Ok(())` constructs that successful empty result.

## The Python contrast

Python propagates errors with **exceptions**, which travel up the call stack
automatically until a `try`/`except` catches them:

```python
def read_number(path):
    with open(path) as f:          # raises on failure; bubbles up on its own
        return int(f.read().strip())

try:
    n = read_number("count.txt")
except (OSError, ValueError) as e:
    print(f"failed: {e}")
```

The philosophies differ: in Python the error path is *invisible* in the
signature — any call might raise. Rust makes it explicit — the return type spells
out that the function can fail, and `?` is the visible marker that says "bubble
this up." Nothing fails silently, and nothing fails invisibly.

## `panic!`: for the unrecoverable

`Result` is for errors you expect and can handle. `panic!` is for **bugs and
impossible states** — situations where continuing makes no sense. By default, a
panic unwinds the stack and runs cleanup; if it reaches the top of a thread,
that thread stops. Rust can also be configured to abort immediately:

```rust
fn withdraw(balance: u32, amount: u32) -> u32 {
    if amount > balance {
        panic!("withdrew {amount} from a balance of {balance}"); // a bug!
    }
    balance - amount
}
```

`unwrap()` and `expect()` from the previous lesson are just panics in disguise —
they panic on `None`/`Err`.

## Choosing: recover or panic?

```
   Can the caller reasonably do something about it?
        │
        ├── YES → return Result / Option   (let them decide)
        │
        └── NO  → panic!  (programmer error, invariant broken, "can't happen")
```

| Situation                              | Use          |
| -------------------------------------- | ------------ |
| File missing, bad user input, network  | `Result`     |
| Index proven in-bounds, logic invariant| `panic!`/`unwrap` |
| Prototype / test where a crash is fine | `unwrap`/`expect` |

The guiding principle: **make the caller's life easy.** Return `Result` from
library-ish code so callers choose how to react; reserve `panic!` for genuine
bugs. Next we put `Option` and borrowing to work with growable collections in
[[Vec & HashMap]].

## Expand `?` in your head

For a `Result`, `operation()?` means approximately:

```rust
let value = match operation() {
    Ok(value) => value,
    Err(error) => return Err(error.into()),
};
```

On success, `?` unwraps the `Ok` value and execution continues. On failure, it
returns early from the **whole enclosing function**, not merely the current
block. Therefore that function must return a compatible `Result` (or `Option`
when `?` is used on an `Option`). The `.into()` allows certain error types to be
converted into the function's declared error type.

Use `match` when this function can recover locally; use `?` when the caller is
better placed to decide.

## Later: operations with different error types

A function may read a file and then parse its contents, producing either an
I/O error or a parsing error. Real applications commonly define an error enum
with one variant for each case. Libraries can then match those variants and
respond precisely.

You will also encounter `Box<dyn std::error::Error>`, which roughly means “an
owned value of some type implementing the `Error` trait.” It is convenient for
small applications that only need to report several possible errors. Its parts
are introduced later:

- [[Traits]] explains `dyn Error`.
- [[Box, Rc & RefCell]] explains `Box`.

You do not need that catch-all syntax to understand `?`: first master
propagating one concrete error type, as `parse_sum` does above.

## Syntax checkpoint

Read `?` as a control-flow instruction attached to the result immediately
before it:

```rust
let number: i32 = text.parse()?;
```

This means: parse the text; if it is `Ok(number)`, place the number in the
binding; if it is `Err(error)`, return that error from the current function.
Therefore the surrounding function must advertise failure with a compatible
`Result<..., ...>` return type.

You are ready to continue when you can expand one `?` into a two-arm `match`.
The next lesson shows why collections return `Option` from lookups and how
borrowing controls iteration.

## Try it

1. **Expand:** Rewrite one `parse()?` as an explicit `match`.
2. **Trace:** Follow `parse_sum("2", "x")` and identify the exact point where the
   function returns.
3. **Fill:** Write a two-step parsing function using `?`.
4. **Choose:** For bad user input and an impossible internal invariant, explain
   which should return `Result` and which may justify a panic.
5. **Repair:** Replace an `unwrap()` first with `expect`, then with proper
   `match` handling.

> **Takeaway:** `?` propagates errors up the call stack with one character,
> turning verbose `match` chains into linear code. Use `Result` for expected,
> recoverable failures and `panic!` only for unrecoverable programmer errors.
