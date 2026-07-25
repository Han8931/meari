---
created: '2026-07-21'
id: rust-i-match-advanced
source: meari-course
study:
  answer: |
    fn parse_port(text: &str) -> Option<u16> {
        let Ok(port) = text.parse::<u16>() else {
            return None;
        };
        Some(port)
    }

    fn classify_status(code: u16) -> String {
        match code {
            n @ 200..=299 => format!("success:{n}"),
            n if n >= 500 => format!("server-error:{n}"),
            n => format!("other:{n}"),
        }
    }
  kind: code
  lang: rust
  prompt: 'Write `parse_port(text: &str) -> Option<u16>` using `let Ok(port) = ... else { return None; }`. Then write `classify_status(code: u16) -> String` using `n @ 200..=299` for `"success:{n}"`, a guard `n if n >= 500` for `"server-error:{n}"`, and a final binding for `"other:{n}"`.'
  starter: |
    fn parse_port(text: &str) -> Option<u16> {
        None
    }

    fn classify_status(code: u16) -> String {
        String::new()
    }
  tests:
  - assert_eq!(parse_port("8080"), Some(8080));
  - assert_eq!(parse_port("eight"), None);
  - assert_eq!(parse_port("70000"), None);
  - assert_eq!(classify_status(204), "success:204");
  - assert_eq!(classify_status(500), "server-error:500");
  - assert_eq!(classify_status(503), "server-error:503");
  - assert_eq!(classify_status(404), "other:404");
subject: Rust (Intermediate)
title: Advanced Pattern Matching
---

[[rust-b-enums|Enums & Pattern Matching]] covered the basics of `match`. Rust's
patterns go further: ranges, guards, and bindings let one `match` express rules
that would otherwise need a tangle of `if`s.

## Ranges and guards

A pattern can be a **range**, and a `match` arm can carry a **guard** — an extra
`if` that must also hold:

```rust
match n {
    0 => "zero",
    x if x < 0 => "negative",     // guard: runs only when x < 0
    1..=9 => "small",             // inclusive range 1 through 9
    _ => "big",
}
```

Arms are tried top to bottom, so order matters: put the specific cases first.

A guard is evaluated only after its pattern matches. Guards also do not make a
non-exhaustive set of patterns exhaustive, because the compiler does not reason
that arbitrary guard expressions cover every remaining value. Keep an
unguarded fallback when other values are possible.

## Binding with `@`

`name @ pattern` matches the pattern *and* captures the value in `name`:

```rust
match code {
    n @ 200..=299 => println!("success ({n})"),
    n => println!("other ({n})"),
}
```

## `let ... else`

When you only care about the success case, `let else` binds it or bails out (with
`return`, `break`, ...) in the failing case — no nesting:

```rust
let Ok(port) = "8080".parse::<u16>() else {
    return;                       // not a number: give up here
};
// `port` is in scope from here on
```

## Your turn

1. **Order arms:** Predict which arm wins when a value satisfies both a broad
   pattern and a later guarded arm.
2. **Separate jobs:** In `n @ 200..=299`, identify the test and the new binding.
3. **Explain divergence:** State why the `else` block of `let ... else` must
   leave the surrounding flow with `return`, `break`, `continue`, or panic.
4. **Implement:** Parse a port with `let ... else`, then classify status codes
   with a range binding, guard, and fallback.
