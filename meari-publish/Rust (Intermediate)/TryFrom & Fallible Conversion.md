---
created: '2026-07-21'
id: rust-i-tryfrom
source: meari-course
study:
  answer: |
    #[derive(Debug, PartialEq)]
    struct Percent(u8);

    impl TryFrom<i32> for Percent {
        type Error = String;

        fn try_from(value: i32) -> Result<Self, Self::Error> {
            if (0..=100).contains(&value) {
                Ok(Percent(value as u8))
            } else {
                Err(format!("{value} is out of range"))
            }
        }
    }
  kind: code
  lang: rust
  prompt: 'Define `struct Percent(u8)` (derive `Debug, PartialEq`) and implement `TryFrom<i32>` with `Error = String`: return `Ok(Percent(..))` when the value is in `0..=100`, otherwise an `Err`.'
  starter: |
    #[derive(Debug, PartialEq)]
    struct Percent(u8);

    // impl TryFrom<i32> for Percent  (type Error = String)
  tests:
  - assert_eq!(Percent::try_from(50), Ok(Percent(50)));
  - assert_eq!(Percent::try_from(0), Ok(Percent(0)));
  - assert_eq!(Percent::try_from(100), Ok(Percent(100)));
  - assert_eq!(Percent::try_from(150), Err("150 is out of range".to_string()));
  - assert_eq!(Percent::try_from(-1), Err("-1 is out of range".to_string()));
subject: Rust (Intermediate)
title: TryFrom & Fallible Conversion
---

Some conversions can fail: not every `i32` is a valid percentage, not every
string is a number. `From` is for conversions that *always* succeed. When a
conversion might fail, use `TryFrom`, which returns a [[rust-b-option-result|
Result]].

## `TryFrom` returns a Result

Like `From`, but `try_from` returns `Result<Self, Self::Error>`, and you pick the
error type:

```rust
#[derive(Debug, PartialEq)]
struct Even(i32);

impl TryFrom<i32> for Even {
    type Error = String;
    fn try_from(v: i32) -> Result<Self, Self::Error> {
        if v % 2 == 0 {
            Ok(Even(v))
        } else {
            Err(format!("{v} is odd"))
        }
    }
}

let ok  = Even::try_from(4);    // Ok(Even(4))
let bad = Even::try_from(5);    // Err("5 is odd")
```

Implementing `TryFrom` also gives you `try_into()` for free, mirroring the
`From`/`Into` pair. In edition 2021 both traits are in the prelude, so there's
nothing to import.

```rust
let even = Even::try_from(4);       // target type is explicit
let even: Result<Even, _> = 4.try_into(); // annotation selects the target
```

The annotation on `try_into` matters when context does not otherwise determine
the destination type. Many types could conceivably be built from one `i32`, so
Rust needs to know which conversion you want.

## Validate before narrowing

The exercise stores a valid percentage as `u8`, but its input is `i32`. Check
the range *before* writing `value as u8`: an `as` cast would wrap or truncate an
invalid integer rather than report failure. After `(0..=100).contains(&value)`
succeeds, the cast is known to preserve the value.

This pattern makes an invalid state hard to construct: code outside the module
can work with `Percent` after validation instead of repeatedly checking a raw
integer. In a real module, keeping the tuple field private preserves that
invariant.

Trace a successful conversion:

```
i32 value 50
  └── validate 0..=100
        └── narrow to u8 (now proven safe)
              └── construct Percent(50)
```

An invalid input exits before the cast. The wrapper type is useful only if code
cannot bypass its checked constructor, which is why field visibility matters.

## Your turn

1. **Predict:** State what `150 as u8` does and why it must not happen before
   validation.
2. **Trace:** Follow `50` and `-1` through the conversion branches.
3. **Infer:** Explain why `let p: Result<Percent, _> = 50.try_into()` supplies
   information the method call needs.
4. **Implement:** Build `Percent`, accepting only `0..=100`.
