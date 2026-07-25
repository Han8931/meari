---
created: '2026-07-21'
id: rust-i-bounds
source: meari-course
study:
  answer: |
    fn largest<T>(items: &[T]) -> Option<T>
    where
        T: PartialOrd + Copy,
    {
        let mut iter = items.iter();
        let mut max = *iter.next()?;
        for &x in iter {
            if x > max {
                max = x;
            }
        }
        Some(max)
    }
  kind: code
  lang: rust
  prompt: 'Write `largest<T>(items: &[T]) -> Option<T>` that returns the maximum element, or `None` for an empty slice. Put its `PartialOrd + Copy` bounds in a `where` clause.'
  starter: |
    fn largest<T>(items: &[T]) -> Option<T>
    where
        T: PartialOrd + Copy,
    {
        None
    }
  tests:
  - assert_eq!(largest(&[3, 7, 2, 9, 4]), Some(9));
  - assert_eq!(largest(&['a', 'z', 'm']), Some('z'));
  - assert_eq!(largest::<i32>(&[]), None);
subject: Rust (Intermediate)
title: Trait Bounds & where Clauses
---

A trait bound restricts a generic to types that implement a trait, so inside the
function you may use that trait's abilities. You met the short form in
[[rust-b-generics|Generics]]; here we combine several bounds and meet the
`where` clause that keeps them readable.

## Combining bounds with `+`

Require more than one trait by joining them with `+`:

```rust
fn later<T: PartialOrd + std::fmt::Debug>(a: &T, b: &T) {
    if a > b {
        println!("{a:?} comes later");
    }
}
```

Common building blocks: `PartialOrd` for `<`/`>`, `Copy` when a value may be
copied out through a reference, `Clone` for `.clone()`, and `Debug` for `{:?}`.
Dereferencing `&T` does not itself promise a copy; the `T: Copy` bound is what
makes `let value = *reference` legal without moving out of borrowed data.

## The `where` clause

When bounds pile up, move them below the signature with `where` — same meaning,
easier to read:

```rust
fn describe_pair<T, U>(left: &T, right: &U) -> String
where
    T: std::fmt::Debug,
    U: std::fmt::Display,
{
    format!("{left:?}: {right}")
}
```

## Why bounds are required

Without a bound, the compiler knows nothing about `T`, so it won't let you
compare or copy it. The bound is your promise — and the reason the body type-
checks for *every* `T` a caller might pick.

Read a bounded function from two perspectives:

```
implementation receives: every operation promised by the bounds
caller must provide:      one concrete type implementing every bound
```

The compiler checks the generic body once against those promises. It does not
wait for a convenient caller to make an otherwise invalid operation safe.

One sharp edge: `PartialOrd` permits values that are not comparable, notably
floating-point `NaN`. The exercise returns the greatest item according to the
comparisons it actually observes; it does not define special `NaN` behavior.
Use `Ord` when your algorithm requires a total ordering.

## Your turn

1. **Read aloud:** Explain what the implementation receives and caller promises
   in `T: PartialOrd + Copy`.
2. **Repair:** Remove each bound in turn and connect the compiler error to the
   operation that needed it.
3. **Rewrite:** Express one short inline bound as an equivalent `where` clause.
4. **Implement:** Write `largest`, returning `None` for an empty slice, with both
   requirements in its `where` clause.
