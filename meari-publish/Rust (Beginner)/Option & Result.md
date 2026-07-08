---
created: "2026-07-08"
id: rust-b-option-result
source: meari-course
subject: Rust (Beginner)
title: Option & Result
---

Most languages represent "no value" with `null` and "it failed" with exceptions.
Rust has **neither**. Instead it uses two ordinary enums — `Option` and `Result`
— so absence and failure become values the type system forces you to handle.
This is [[Enums & Pattern Matching]] put to work.

## `Option<T>`: a value that might be absent

```rust
enum Option<T> {      // built into the standard library
    Some(T),          // there IS a value, here it is
    None,             // there is no value
}
```

Because there's no `null`, a function that might not return something returns an
`Option`:

```rust
fn first_char(s: &str) -> Option<char> {
    s.chars().next()          // Some(c), or None if the string is empty
}

match first_char("hi") {
    Some(c) => println!("first is {c}"),
    None    => println!("empty string"),
}
```

The compiler won't let you use the inner value without first dealing with the
`None` case — that's how Rust abolishes the null-pointer crash.

### The Python contrast

Python usually signals "no value" with `None` returned directly, with nothing at
the type level forcing you to check for it:

```python
def first_char(s):
    return s[0] if s else None    # returns None, but the caller may forget it

c = first_char("hi")              # 'h'
print(c.upper())                  # crashes later if `c` was None
```

Rust's `Option<char>` turns that latent crash into a compile-time requirement:
you *must* deal with the `None` case before touching the value. The famous
"`NoneType` has no attribute…" error simply can't happen.

## `Result<T, E>`: a value or an error

```rust
enum Result<T, E> {
    Ok(T),            // success, carrying the value
    Err(E),           // failure, carrying an error
}
```

Anything that can fail — file I/O, parsing, network calls — returns a `Result`:

```rust
fn parse_age(s: &str) -> Result<u32, std::num::ParseIntError> {
    s.parse::<u32>()          // Ok(number) or Err(parse error)
}

match parse_age("30") {
    Ok(n)  => println!("age is {n}"),
    Err(e) => println!("bad input: {e}"),
}
```

## Option vs Result — which to use

```
   "might there be nothing here?"     →   Option<T>   (Some / None)
   "might this operation FAIL, and
    if so, why?"                      →   Result<T,E> (Ok / Err)
```

| Aspect      | `Option<T>`         | `Result<T, E>`             |
| ----------- | ------------------- | -------------------------- |
| Variants    | `Some(T)` / `None`  | `Ok(T)` / `Err(E)`         |
| Models      | presence / absence  | success / failure + reason |
| Carries why | no                  | yes — the `E` error value  |

## Getting the value out

**`match`** is the fully explicit way, but the standard library gives you
concise helpers for common cases:

```rust
let maybe: Option<i32> = Some(5);

maybe.unwrap();            // 5    — but PANICS if None
maybe.expect("need a value"); // 5 — panic with your message if None
maybe.unwrap_or(0);       // 5, or 0 if it were None (a safe default)
maybe.unwrap_or_else(|| compute_default()); // default computed lazily
```

> ⚠️ `unwrap` and `expect` **crash** on `None`/`Err`. They're fine in quick
> experiments, tests, and cases you can prove are impossible — but in real code
> prefer handling the empty case or propagating it (next lesson).

## Combinators: transform without unwrapping

You can operate on the value *inside* an `Option`/`Result` without tearing it
open, which keeps error handling flat and readable:

```rust
let len: Option<usize> =
    Some("hello").map(|s| s.len());        // Some(5)

let doubled: Option<i32> =
    Some(4).filter(|&n| n > 0).map(|n| n * 2); // Some(8)

// chain fallible steps; the first None/Err short-circuits
let n: Option<i32> = Some("42").and_then(|s| s.parse().ok()); // Some(42)
```

| Method        | Does                                              |
| ------------- | ------------------------------------------------- |
| `map`         | transform the inner value if present              |
| `and_then`    | chain another Option/Result-returning step        |
| `unwrap_or`   | supply a fallback value                           |
| `ok_or`       | turn an `Option` into a `Result` with an error    |

Handling every `Option`/`Result` by hand gets verbose when errors need to travel
up through many function calls. The `?` operator streamlines exactly that — see
[[Error Propagation & Panics]].

## Try it

1. Write a function that returns `Option<char>` for the first character of a string.
2. Match on that `Option` and handle both `Some` and `None`.
3. Parse a string into a number with `"42".parse::<i32>()` and handle the `Result`.

> **Takeaway:** Rust replaces `null` with `Option` and exceptions with `Result`,
> making "might be absent" and "might fail" explicit in every signature. The
> compiler then guarantees you can't forget to handle them.
