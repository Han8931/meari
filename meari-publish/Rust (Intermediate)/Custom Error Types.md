---
created: '2026-07-21'
id: rust-i-custom-errors
source: meari-course
study:
  answer: |
    use std::error::Error;

    #[derive(Debug, PartialEq)]
    enum ParseError {
        Empty,
        NotANumber(String),
    }

    impl std::fmt::Display for ParseError {
        fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
            match self {
                ParseError::Empty => write!(f, "input was empty"),
                ParseError::NotANumber(s) => write!(f, "'{s}' is not a number"),
            }
        }
    }

    impl Error for ParseError {}

    fn parse_u32(s: &str) -> Result<u32, ParseError> {
        let s = s.trim();
        if s.is_empty() {
            return Err(ParseError::Empty);
        }
        s.parse::<u32>().map_err(|_| ParseError::NotANumber(s.to_string()))
    }
  kind: code
  lang: rust
  prompt: 'Define `enum ParseError { Empty, NotANumber(String) }` (derive `Debug, PartialEq`), implement `std::fmt::Display` and `std::error::Error` for it, and write `parse_u32(s: &str) -> Result<u32, ParseError>` that trims `s`, errors `Empty` on blank input, and `NotANumber` when it isn''t a `u32`.'
  starter: |
    use std::error::Error;

    #[derive(Debug, PartialEq)]
    enum ParseError {
        Empty,
        NotANumber(String),
    }

    // impl std::fmt::Display and Error for ParseError, then parse_u32

    fn parse_u32(s: &str) -> Result<u32, ParseError> {
        Ok(0)
    }
  tests:
  - assert_eq!(parse_u32("42"), Ok(42));
  - assert_eq!(parse_u32("  7 "), Ok(7));
  - assert_eq!(parse_u32(""), Err(ParseError::Empty));
  - assert_eq!(parse_u32("abc"), Err(ParseError::NotANumber("abc".to_string())));
  - 'assert_eq!(format!("{}", ParseError::Empty), "input was empty");'
  - 'fn assert_error<E: std::error::Error>() {} assert_error::<ParseError>();'
subject: Rust (Intermediate)
title: Custom Error Types
---

[[rust-b-question|Error Propagation & Panics]] used `Result` with ready-made
error types. Real programs define their *own* error type — usually an `enum` of
everything that can go wrong — so callers can tell the cases apart.

## An error enum

```rust
#[derive(Debug, PartialEq)]
enum ParseError {
    Empty,
    NotANumber(String),
}
```

## Make it a proper error

Two traits turn an enum into a first-class error. `Display` gives a human
message; `std::error::Error` lets it slot into the wider error ecosystem and
supports features such as exposing an underlying source error. `Error` requires
`Display` + `Debug`, which is why we add both:

```rust
impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            ParseError::Empty => write!(f, "input was empty"),
            ParseError::NotANumber(s) => write!(f, "'{s}' is not a number"),
        }
    }
}

impl std::error::Error for ParseError {}     // default methods are enough
```

Now `ParseError` prints nicely and works anywhere an `Error` is expected.

Keep the roles separate:

| Piece | Audience | Purpose |
| ----- | -------- | ------- |
| enum variants | program | structured cases callers can match |
| `Debug` | developers | diagnostic representation |
| `Display` | people | readable message |
| `Error` | generic error-handling code | common error contract and optional source chain |

Returning only strings throws away the structured distinction between empty and
invalid input. An enum preserves it while still providing a message.

## Your turn

1. **Classify:** Identify the structured case and human message for each variant.
2. **Fill:** Implement `Display`, then verify `{}` and `{:?}` serve different
   purposes.
3. **Satisfy the contract:** Implement `Error` and explain why `Debug` and
   `Display` are required first.
4. **Implement:** Write `parse_u32`, preserving empty and non-numeric failures as
   different variants.
