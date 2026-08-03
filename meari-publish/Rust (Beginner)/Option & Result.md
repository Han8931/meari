---
created: "2026-07-08"
id: rust-b-option-result
source: meari-course
study:
  answer: |
    fn first_char(text: &str) -> Option<char> {
        text.chars().next()
    }
  kind: code
  lang: rust
  prompt: 'Complete `first_char(text: &str) -> Option<char>` with `text.chars().next()`. `next()` returns `Some(character)` for non-empty text and `None` for empty text.'
  hint: |
    Iterate the text as characters with `.chars()`, then ask that iterator for its next item — it already returns `Some(c)` or `None`, exactly the return type, so you can hand it straight out.
  starter: |
    fn first_char(text: &str) -> Option<char> {
        None
    }
  tests:
    - assert_eq!(first_char("rust"), Some('r'));
    - assert_eq!(first_char(""), None);
    - assert_eq!(first_char("éclair"), Some('é'));
subject: Rust (Beginner)
title: Option & Result
---

Many languages commonly represent “no value” with `null` and failure with
exceptions. Safe, idiomatic Rust instead uses two ordinary enums—`Option` and
`Result`—so absence and failure become values the type system forces you to
handle. Raw pointers can be null and Rust can panic, but neither is the ordinary
way to model an optional value or a recoverable error. This lesson is
[[Enums & Pattern Matching]] put to work.

## `Option<T>`: a value that might be absent

```rust
enum Option<T> {      // built into the standard library
    Some(T),          // there IS a value, here it is
    None,             // there is no value
}
```

Here `T` means “whatever type may be inside.” If `T` is `char`, the complete
type is `Option<char>` and the present variant is `Some(char)`. You can read
angle brackets this way for now; [[Generics]] later explains how one definition
works for many choices of `T`.

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

## Handling every possibility

Do not picture an `Option<i32>` as “an `i32` that might be broken.” It is a
complete value with one of two shapes:

```
Option<i32>
   ├── Some(42)  → contains an i32
   └── None      → contains no i32
```

A `match` converts both shapes into the result your program needs:

```rust
let maybe_number: Option<i32> = Some(4);

let doubled: i32 = match maybe_number {
    Some(number) => number * 2,
    None => 0,
};
```

Trace it as:

1. Inspect which variant `maybe_number` contains.
2. If it is `Some`, bind the inner value to `number`.
3. If it is `None`, use the fallback.
4. Both arms produce an `i32`, so the whole `match` produces an `i32`.

The wrapper has not been “ignored”; handling both variants is what safely turns
an `Option<i32>` into a definite `i32`.

## Later: convenience methods

**`match`** is the fully explicit way, but the standard library gives you
concise helpers for common cases:

```rust
let maybe: Option<i32> = Some(5);

maybe.unwrap();            // 5    — but PANICS if None
maybe.expect("need a value"); // 5 — panic with your message if None
maybe.unwrap_or(0);       // 5, or 0 if it were None (a safe default)
```

> `unwrap` and `expect` end the program with a panic on `None` or `Err`. They can
> be convenient in a small experiment or a test. In application code, first ask
> what the program should do when the value is absent or the operation fails;
> handle that case or propagate it to the caller (next lesson).

Handling every `Option`/`Result` by hand gets verbose when errors need to travel
up through many function calls. The `?` operator streamlines exactly that — see
[[Error Propagation & Panics]]. Later, once closures have been introduced,
[[Closures & Iterators]] shows `map`, `filter`, and `and_then` as compact ways
to transform wrapped values.

## Follow the type parameter

In `Option<i32>`, `i32` is the type inside `Some`; `None` carries no number. In
`Result<u32, ParseIntError>`, `u32` is inside `Ok` and the error type is inside
`Err`. The wrapper is part of the type—you cannot use an `Option<i32>` directly
where an `i32` is required because the value might be absent.

Ask “what should happen in every variant?” before reaching for `unwrap`.

## Syntax checkpoint

Read the outer type first, then what it may contain:

```rust
Option<i32>       // Option containing an i32 when it is Some
Result<i32, E>    // Result containing an i32 on success or E on failure
```

The angle brackets provide type information; they do not access a value. At
runtime you handle the variants you already know:

```rust
match maybe {
    Some(number) => number,
    None => 0,
}
```

You are ready to continue when you can choose `Option` for absence, `Result` for
failure, and match both cases. The next lesson does not replace `match`; it adds
`?` as shorthand for the repeated “return the error, otherwise continue” case.

## Try it

1. **Classify:** Decide whether “find a matching item” and “parse user input”
   should return `Option` or `Result`, and explain why.
2. **Trace:** Follow both `Some(4)` and `None` through the `doubled` match.
3. **Fill:** Write a function returning `Option<char>` for the first character
   of a string.
4. **Handle failure:** Parse `"42"` and `"nope"` with `.parse::<i32>()`, matching
   both `Ok` and `Err`.

> **Takeaway:** Rust replaces `null` with `Option` and exceptions with `Result`,
> making "might be absent" and "might fail" explicit in every signature. The
> compiler then guarantees you can't forget to handle them.
