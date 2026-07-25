---
created: '2026-07-21'
id: rust-i-operator
source: meari-course
study:
  answer: |
    use std::ops::Add;

    #[derive(Debug, Clone, Copy, PartialEq)]
    struct Point {
        x: i32,
        y: i32,
    }

    impl Add for Point {
        type Output = Point;

        fn add(self, other: Point) -> Point {
            Point {
                x: self.x + other.x,
                y: self.y + other.y,
            }
        }
    }
  kind: code
  lang: rust
  prompt: 'Define `struct Point { x: i32, y: i32 }` deriving `Debug, Clone, Copy, PartialEq`, and implement `std::ops::Add` so `p1 + p2` adds `x` and `y` componentwise.'
  starter: |
    use std::ops::Add;

    #[derive(Debug, Clone, Copy, PartialEq)]
    struct Point {
        x: i32,
        y: i32,
    }

    // impl Add for Point
  tests:
  - 'assert_eq!(Point { x: 1, y: 2 } + Point { x: 3, y: 4 }, Point { x: 4, y: 6 });'
  - 'assert_eq!(Point { x: 0, y: 0 } + Point { x: -1, y: 5 }, Point { x: -1, y: 5 });'
  - 'assert_eq!(Point { x: 10, y: -2 } + Point { x: -10, y: 2 }, Point { x: 0, y: 0 });'
subject: Rust (Intermediate)
title: Operator Overloading
---

Operators like `+` use traits. Roughly, `a + b` becomes
`std::ops::Add::add(a, b)`. Implement that trait and your own type gets `+` too,
with the operand and result types made explicit by the implementation.

## Implementing `Add`

`Add` has an associated type `Output` (the result type) and one method, `add`:

```rust
use std::ops::Add;

#[derive(Debug, Clone, Copy, PartialEq)]
struct Money(i32);        // cents

impl Add for Money {
    type Output = Money;
    fn add(self, other: Money) -> Money {
        Money(self.0 + other.0)
    }
}

let total = Money(150) + Money(99);   // Money(249)
```

Note `add` takes `self` by value, so deriving `Copy` (as above) keeps the
operands usable afterward instead of moving them away.

Trace the types in `Money(150) + Money(99)`:

```
Self = Money
Rhs  = Money       (the default)
Output = Money     (the associated type)
add(self, other) consumes two Money values and returns one Money
```

The full trait is `Add<Rhs = Self>`. `Rhs` has a default, so `impl Add for Money`
means the right-hand side is also `Money`. You can deliberately choose a
different input or output when the operation calls for it:

```rust
// Duration + u32 seconds -> Duration
// impl Add<u32> for Duration { type Output = Duration; ... }
```

For non-`Copy` values, decide whether consuming the operands is appropriate. You
can instead implement `Add` for references, but that is a separate implementation
with its own ownership contract.

```rust
// This implementation would allow `&left + &right` without consuming either:
// impl Add<&Money> for &Money {
//     type Output = Money;
//     ...
// }
```

The operator syntax does not hide ownership. The chosen trait implementation
determines whether the left and right operands are values or references.

## Overload responsibly

Only implement an operator when its meaning is obvious. `+` on two points or two
money amounts reads naturally; a surprising overload makes code harder to follow,
not easier. The other operators follow the same shape: `Sub`, `Mul`, `Index`,
and so on, each its own trait in `std::ops`.

## Your turn

1. **Substitute:** Name `Self`, `Rhs`, and `Output` in the `Money` implementation.
2. **Predict ownership:** Explain why non-`Copy` operands would be unavailable
   after value-based `+`.
3. **Choose semantics:** Give an example where defining `+` would be surprising
   and should be avoided.
4. **Implement:** Give a 2-D `Point` componentwise addition.
