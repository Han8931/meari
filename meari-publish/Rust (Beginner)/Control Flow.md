---
created: "2026-07-08"
id: rust-b-control
source: meari-course
study:
  answer: |
    fn describe(n: i32) -> i32 {
        if n < 0 {
            -1
        } else if n == 0 {
            0
        } else {
            1
        }
    }
  kind: code
  lang: rust
  prompt: 'Write `describe(n: i32) -> i32` as an `if` expression: return -1 for a negative number, 0 for zero, and 1 for a positive number.'
  hint: |
    Make the whole `if / else if / else` the return value: each branch ends in a bare number with no semicolon, so the matching branch's value falls out as the result.
  starter: |
    fn describe(n: i32) -> i32 {
        0
    }
  tests:
    - assert_eq!(describe(-8), -1);
    - assert_eq!(describe(0), 0);
    - assert_eq!(describe(12), 1);
subject: Rust (Beginner)
title: Control Flow
---

Rust has the control-flow tools you'd expect — `if`, loops, `for` — but with one
twist that shapes how idiomatic Rust reads: **almost everything is an
expression that produces a value**, not just a statement that does something.

## `if` is an expression

You can use `if`/`else` the ordinary way, but you can also let it *return* a
value directly into a `let`:

```rust
let n = 7;

// classic branching
if n % 2 == 0 {
    println!("even");
} else {
    println!("odd");
}

// if AS an expression — no ternary operator needed
let label = if n % 2 == 0 { "even" } else { "odd" };
```

Two rules when using `if` as an expression:

1. Every branch must produce the **same type** (`"even"` and `"odd"` are both
   `&str`).
2. The branch value is the block's **last expression with no semicolon** — a
   trailing `;` turns it into a statement that yields `()`.

```
  { "even" }   →  block evaluates to "even"   (expression)
  { "even"; }  →  block evaluates to ()        (statement — usually a bug here)
```

## Three kinds of loop

```rust
// 1. loop — infinite until you break; can RETURN a value
let mut n = 0;
let doubled = loop {
    n += 1;
    if n == 10 { break n * 2; }   // break carries a value out
};                                 // doubled == 20

// 2. while — loops while a condition holds
let mut count = 3;
while count > 0 {
    println!("{count}...");
    count -= 1;
}

// 3. for — iterate over a collection or range (the one you'll use most)
for i in 1..=5 {           // 1, 2, 3, 4, 5
    println!("{i}");
}
```

### Ranges

| Syntax  | Meaning                | Example expands to |
| ------- | ---------------------- | ------------------ |
| `0..n`  | exclusive end          | `0..3` → 0, 1, 2   |
| `0..=n` | inclusive end          | `0..=3` → 0,1,2,3  |

## The same in Python

Python's conditional expression is its version of `if`-as-a-value, and `range`
mirrors Rust's exclusive `0..n` range:

```python
label = "even" if n % 2 == 0 else "odd"   # ~ let label = if … { } else { };

for i in range(1, 6):                      # ~ for i in 1..6  (end-exclusive)
    print(i)                               # 1, 2, 3, 4, 5
```

Python has no `loop … break value` construct — you'd use `while True:` with a
`break`, but it can't *return* a value out of the loop the way Rust's `loop`
does.

For now, use `for` with numeric ranges. Later, after ownership and collections
have names, [[Vec & HashMap]] spells out why `for item in values`,
`for item in &values`, and `for item in &mut values` behave differently.

## Later: labeled breaks for nested loops

When loops nest, a plain `break` only exits the innermost one. Label a loop with
`'name:` to break out of an outer loop directly:

```rust
'outer: for i in 0..5 {
    for j in 0..5 {
        if i * j > 6 { break 'outer; }   // jumps all the way out
        println!("{i}·{j}");
    }
}
```

## Expressions, statements, and `()`: a reminder

An **expression** evaluates to a value; a **statement** performs an action and
does not pass a useful value on.

```rust
let a = 2 + 3; // `2 + 3` is an expression with value 5
let b = {
    let x = 10;
    x * 2      // no semicolon: the block's value is 20
};
```

`let x = 10;` is a statement. Adding a semicolon to `x * 2` changes the block's
result to the unit value `()`, roughly “no meaningful value.” Thus “expected
`i32`, found `()`” often means a semicolon discarded a number you meant to
return.

For `for i in 1..=3`, Rust obtains `1`, `2`, and `3` one at a time, binds each to
`i`, and runs the body. The loop creates `i`; you do not declare it beforehand.
Tracing these values on paper quickly exposes most off-by-one errors.

## Syntax checkpoint

Start with only the most common forms:

```rust
if temperature < 0 {
    println!("freezing");
}

for n in 1..4 {
    println!("{n}");
}
```

Read the first as “if this `bool` condition is true, run this block.” Read the
second as “for each number from 1 up to, but not including, 4, run this block.”
You can postpone `loop` values and labels until ordinary `if` and `for` feel
comfortable.

You are ready to continue when you can predict which branch runs and which
numbers a range produces. The following lesson is a conceptual pause: it
explains why Rust performs so many checks before the ownership section begins.

## Try it

1. Write an `if` expression that stores either `"small"` or `"big"` in a variable.
2. Use a `for` loop to print the numbers 1 through 5.
3. Use a `while` loop to count down from 3.
4. Remove the final expression from an `if` branch and read the resulting type
   error.

> **Takeaway:** lean on Rust's expression orientation — `let x = if …` and
> `let x = loop { … break v }` replace clumsy mutable temporaries. And decide
> use a final expression rather than a mutable temporary when a branch or loop
> is calculating a value.
