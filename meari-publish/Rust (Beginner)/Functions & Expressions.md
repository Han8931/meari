---
created: "2026-07-25"
id: rust-b-functions
source: meari-course
study:
  answer: |
    fn rectangle_area(width: i32, height: i32) -> i32 {
        width * height
    }
  kind: code
  lang: rust
  prompt: 'Complete `rectangle_area(width: i32, height: i32) -> i32`. Return `width * height` as a final expression without a semicolon.'
  hint: |
    A function returns its final expression when that line has no semicolon — so the body is just the product of the two parameters, with no `return` keyword.
  starter: |
    fn rectangle_area(width: i32, height: i32) -> i32 {
        0
    }
  tests:
    - assert_eq!(rectangle_area(3, 4), 12);
    - assert_eq!(rectangle_area(0, 9), 0);
    - assert_eq!(rectangle_area(7, 2), 14);
subject: Rust (Beginner)
title: Functions & Expressions
---

You have already seen `fn main()`, but the exercises in this course use small
functions of their own. Before going further, this lesson spells out how to read
and write one. The punctuation is compact, so read the following signature from
left to right:

```rust
fn add(a: i32, b: i32) -> i32 {
    a + b
}
```

> Define a function named `add`. It receives two parameters, `a` and `b`. Each
> parameter is an `i32`. Calling the function produces an `i32`.

## Parameters, arguments, and return types

The names in a function definition are **parameters**. The values supplied by a
caller are **arguments**:

```rust
fn multiply(left: i32, right: i32) -> i32 {
    left * right
}

let area = multiply(3, 4);
//                  └── arguments supplied by this call
```

```
fn multiply(left: i32, right: i32) -> i32
             └────┬────  └────┬────
              parameter     parameter
              name + type   name + type
```

Rust requires every parameter type to be written explicitly. The arrow
`-> Type` states the return type. (`i32` is Rust's usual whole-number type; see
[[Data Types & Type Casting]] if you need a refresher.) When there is no arrow,
the function returns the **unit value** `()`, meaning “no useful result”:

```rust
fn announce() {
    println!("Starting");
}
```

## A block's final expression is its value

Most Rust functions return their result with a final expression:

```rust
fn square(n: i32) -> i32 {
    n * n
}
```

There is deliberately no semicolon after `n * n`. Follow the execution:

```
square(4)
   │
   ├── n is bound to 4
   ├── evaluate n * n
   └── the block's final value is 16 → return 16
```

A semicolon says “finish this expression and discard its value.” This version
does not compile:

```rust
fn broken_square(n: i32) -> i32 {
    n * n;
    //     ^ semicolon discards the i32
}   // the body now produces (), but the signature promised i32
```

When the compiler says **“expected `i32`, found `()`,”** first inspect the final
line for an accidental semicolon or a missing result.

## Expression versus statement

An **expression** produces a value. A **statement** completes an action without
passing a useful value onward.

```rust
let total = 2 + 3;
//  ^^^^^   ^^^^^
// binding  expression whose value is 5
```

`2 + 3` is an expression. The whole `let total = 2 + 3;` line is a statement.
A braced block is also an expression, and its final expression becomes the
block's value:

```rust
let total = {
    let price = 8;
    let quantity = 3;
    price * quantity
};

println!("{total}"); // 24
```

This same rule powers blocks such as `if`, `match`, and `loop` later in the
course.

## Explicit `return`

The `return` keyword exits the whole function immediately:

```rust
fn square_with_return(n: i32) -> i32 {
    return n * n;
}
```

The semicolon is required as part of the `return` statement. Later, `return`
will be useful for exiting early from a branch. For the ordinary result at the
bottom of a function, Rust code normally uses a final expression.

## Calling and tracing a function

```rust
fn twice(n: i32) -> i32 {
    n * 2
}

fn add_one(n: i32) -> i32 {
    n + 1
}

let answer = add_one(twice(5));
```

Trace from the innermost call outward:

```
twice(5)       → 10
add_one(10)    → 11
answer         = 11
```

A function call is an expression because it produces the function's return
value. Rust checks each boundary: the arguments must match the parameter types,
and the body must match the declared return type.

## The optional Python comparison

Python uses indentation where Rust uses braces and usually omits parameter and
return types:

```python
def add(a, b):
    return a + b
```

The Rust difference that matters most is not the braces. It is that the
signature is a checked contract: `fn add(a: i32, b: i32) -> i32` tells both the
reader and compiler exactly what may enter and leave.

## Syntax checkpoint

Read a signature in four pieces instead of as one dense line:

```rust
fn add(a: i32, b: i32) -> i32 {
//  ^^^  ^^^^^^  ^^^^^^    ^^^
//  name input 1 input 2   output
    a + b
}
```

The commas separate inputs, `:` separates each input's name from its type, and
`->` points to the output type. Inside the braces, the final expression supplies
the output. Do not combine this with explicit `return` yet unless you need an
early exit.

You are ready for control flow when you can write one two-parameter function
and explain why its final expression has no semicolon. The next lesson keeps
that same “a block produces its final value” rule and applies it to `if`.

## Try it

1. **Trace:** Without running it, predict `add_one(twice(6))`. Then check.
2. **Repair:** Add a semicolon after `n * 2` in `twice`, run `cargo check`, and
   explain why the error mentions `()`.
3. **Fill:** Write `fn cube(n: i32) -> i32` with a final expression.
4. **Create:** Write a function with two parameters and call it from `main`.

> **Takeaway:** a function signature is a type-checked boundary. Parameters are
> written `name: Type`, `-> Type` declares the result, and the final expression
> without a semicolon supplies that result.
