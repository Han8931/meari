---
created: '2026-07-21'
id: rust-i-from-into
source: meari-course
study:
  answer: |
    struct Celsius(f64);
    struct Fahrenheit(f64);

    impl From<Celsius> for Fahrenheit {
        fn from(c: Celsius) -> Self {
            Fahrenheit(c.0 * 9.0 / 5.0 + 32.0)
        }
    }
  kind: code
  lang: rust
  prompt: Given `struct Celsius(f64)` and `struct Fahrenheit(f64)`, implement `From<Celsius> for Fahrenheit` using `F = C * 9/5 + 32`. (You then get `.into()` for free.)
  starter: |
    struct Celsius(f64);
    struct Fahrenheit(f64);

    // impl From<Celsius> for Fahrenheit
  tests:
  - assert_eq!(Fahrenheit::from(Celsius(100.0)).0, 212.0);
  - assert_eq!(Fahrenheit::from(Celsius(0.0)).0, 32.0);
  - 'let f: Fahrenheit = Celsius(-40.0).into(); assert_eq!(f.0, -40.0);'
subject: Rust (Intermediate)
title: From & Into
---

Rust has a standard way to say "this type can be built from that one": the
`From` trait. Implement `From`, and you get `Into` for free — two views of the
same conversion.

## Implement `From`, get `Into` free

```rust
struct Celsius(f64);
struct Kelvin(f64);

impl From<Celsius> for Kelvin {
    fn from(c: Celsius) -> Self {
        Kelvin(c.0 + 273.15)
    }
}

let k = Kelvin::from(Celsius(0.0));    // 273.15  — using From
let k: Kelvin = Celsius(0.0).into();   // 273.15  — same thing, using Into
```

You wrote one impl; a standard-library **blanket implementation** supplies the
matching `Into` implementation. Nothing is derived and no code-generation
attribute is involved. Prefer implementing `From` and let callers choose
whichever direction reads better.

## The `From` contract

Use `From` only when conversion is:

- **infallible** — every source value has a destination;
- **value-preserving enough for the domain** — it does not silently reject or
  clamp an invalid source;
- **obvious** — callers should not have to guess which interpretation was used.

If validation can fail, use [[rust-i-tryfrom|TryFrom & Fallible Conversion]].
If several equally plausible conversions exist, a named constructor may
communicate intent better.

## Why it matters

`From`/`Into` are everywhere: `String::from("hi")`, `Vec::from(...)`, and the `?`
operator uses `From` to convert error types automatically. After `TryFrom` and
custom errors, the final lesson in this module puts that connection to work.

## Your turn

1. **Classify:** Decide whether integer-to-percentage validation belongs in
   `From` or `TryFrom`.
2. **Expand:** Rewrite an `.into()` conversion as the equivalent `From::from`.
3. **Explain:** State where the `Into` implementation came from.
4. **Implement:** Convert temperatures with
   `Fahrenheit = Celsius * 9/5 + 32`.
