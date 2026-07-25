---
created: '2026-07-21'
id: rust-i-error-propagation
source: meari-course
study:
  answer: |
    #[derive(Debug, PartialEq)]
    enum AppError {
        Parse(String),
    }

    impl From<std::num::ParseIntError> for AppError {
        fn from(e: std::num::ParseIntError) -> Self {
            AppError::Parse(e.to_string())
        }
    }

    fn total(text: &str) -> Result<i32, AppError> {
        let mut sum = 0;
        for token in text.split_whitespace() {
            sum += token.parse::<i32>()?;
        }
        Ok(sum)
    }
  kind: code
  lang: rust
  prompt: 'Define `enum AppError { Parse(String) }` (derive `Debug, PartialEq`), implement `From<std::num::ParseIntError>` for it, then write `total(text: &str) -> Result<i32, AppError>` that sums the whitespace-separated integers, using `?` so the parse error converts automatically.'
  starter: |
    #[derive(Debug, PartialEq)]
    enum AppError {
        Parse(String),
    }

    // impl From<std::num::ParseIntError> for AppError

    fn total(text: &str) -> Result<i32, AppError> {
        Ok(0)
    }
  tests:
  - assert_eq!(total("1 2 3"), Ok(6));
  - assert_eq!(total("10  -4 2"), Ok(8));
  - assert_eq!(total(""), Ok(0));
  - 'match total("1 x 3") { Err(AppError::Parse(message)) => assert!(!message.is_empty()), other => panic!("unexpected result: {other:?}"), }'
subject: Rust (Intermediate)
title: Error Conversion with ?
---

The `?` operator does more than early-return an error — it *converts* it. When
`?` returns an error whose type differs from the function's, it calls `From` to
turn one into the other. Implement `From`, and `?` stitches libraries with
different error types into your own.

## `?` calls `From` for you

```rust
#[derive(Debug, PartialEq)]
enum AppError {
    Parse(String),
}

impl From<std::num::ParseIntError> for AppError {
    fn from(e: std::num::ParseIntError) -> Self {
        AppError::Parse(e.to_string())
    }
}

fn double(s: &str) -> Result<i32, AppError> {
    let n: i32 = s.parse()?;    // parse() yields ParseIntError;
    Ok(n * 2)                   // ? converts it to AppError via From
}
```

`s.parse()` fails with a `ParseIntError`, but `double` returns `AppError`. The
`?` sees the mismatch, calls your `From` impl, and returns the converted error —
no manual `match` or `map_err`.

For `Result<T, E>`, you can read `expression?` approximately as:

```rust
match expression {
    Ok(value) => value,
    Err(error) => return Err(AppError::from(error)),
}
```

The real definition is more general, but this expansion is the right mental
model here. Notice that `?` unwraps only the success value; the surrounding
function still returns a `Result`.

Trace the two paths:

```
"12".parse()?  → Ok(12)  → bind 12 and continue
"x".parse()?   → Err(ParseIntError)
                    └── AppError::from(error)
                          └── return Err(AppError::Parse(...))
```

The conversion happens only on the error path. The success type still has to
fit the expression around `?`.

## Why this is the good path

Each library reports its own error type. Your `From` impls are the one place you
translate them into *your* vocabulary; everywhere else, a plain `?` just works.
That's how idiomatic Rust keeps error handling terse without hiding it.

One tradeoff: converting an error into a string loses its structured source.
Production error enums often store the original `ParseIntError` in their
variant and expose it through `Error::source`. This exercise uses a string so it
can focus narrowly on the `From`/`?` connection.

## Your turn

1. **Expand:** Rewrite `s.parse()?` as a `match` containing the explicit
   conversion.
2. **Trace:** Follow `"1 x 3"` and identify which later token is never visited.
3. **Compare errors:** Explain what is lost by storing only `e.to_string()`.
4. **Implement:** Sum whitespace-separated integers and surface a bad token as
   `AppError` through `?`.
