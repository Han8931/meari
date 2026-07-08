---
created: "2026-07-08"
id: rust-b-question
source: meari-course
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
step. (One new piece of syntax shows up in the return type below:
`Box<dyn std::error::Error>` means "some error value, boxed on the heap" — a
catch-all that lets one function surface several different kinds of error. Read
it as "any error" for now; [[Box, Rc & RefCell]] covers `Box` in full.)

```rust
fn read_number(path: &str) -> Result<i32, Box<dyn std::error::Error>> {
    let contents = match std::fs::read_to_string(path) {
        Ok(c) => c,
        Err(e) => return Err(Box::new(e)),   // forward the I/O error by hand
    };
    let n = match contents.trim().parse::<i32>() {
        Ok(n) => n,
        Err(e) => return Err(Box::new(e)),   // ...and again for the parse error
    };
    Ok(n)
}
```

That `match … return Err(e)` boilerplate repeats for every fallible call.

## The `?` operator

`?` collapses that pattern into one character. On a `Result`, it means:
**if `Ok`, unwrap the value and keep going; if `Err`, return that error from the
whole function immediately.**

```rust
fn read_number(path: &str) -> Result<i32, Box<dyn std::error::Error>> {
    let contents = std::fs::read_to_string(path)?;  // Err → early return
    let n: i32 = contents.trim().parse()?;          // Err → early return
    Ok(n)
}
```

```
   read_to_string(path)?
        │
        ├── Ok(text)  ──►  bind text, continue
        └── Err(e)    ──►  return Err(e) from read_number  (short-circuit)
```

Two conditions to use `?`:

1. The enclosing function must itself return a `Result` (or `Option`) — `?`
   needs somewhere to return the error *to*.
2. The error types must be compatible. `Box<dyn Error>` is a common beginner
   choice because most errors convert into it automatically.

`?` also works on `Option`, returning `None` early.

## `?` in `main`

`main` can return a `Result`, letting you use `?` at the top level:

```rust
fn main() -> Result<(), Box<dyn std::error::Error>> {
    let n = read_number("count.txt")?;
    println!("read {n}");
    Ok(())
}
```

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

## Try it

1. Write a function that reads a file with `std::fs::read_to_string` and returns a `Result`.
2. Use `?` instead of a manual `match` to propagate the error.
3. Replace an `unwrap()` with `expect("explain what went wrong")`, then with proper `match` handling.

> **Takeaway:** `?` propagates errors up the call stack with one character,
> turning verbose `match` chains into linear code. Use `Result` for expected,
> recoverable failures and `panic!` only for unrecoverable programmer errors.
