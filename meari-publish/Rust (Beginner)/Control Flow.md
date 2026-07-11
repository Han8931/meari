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

## Iterating a collection: borrow vs consume

How you write the `for` loop decides whether you can still use the collection
afterward. The important detail is that a `for` loop must first obtain an
**iterator** from the expression after `in`. It can obtain that iterator by
borrowing the collection or by taking ownership of it.

### Borrowing the collection

```rust
let names = vec!["Ana", "Bo", "Cy"];

for n in &names {          // BORROW each item — names is still usable after
    println!("{n}");
}
println!("{}", names.len()); // ✅ still fine
```

Here the loop expression is `&names`, a shared reference to the vector. The
iterator therefore borrows `names` and yields a reference to each element. It
never owns the vector. When the loop ends, the borrow ends, so the original
owner—`names`—is still valid.

Trace the ownership:

```text
names owns the Vec
        │
        └── loop temporarily borrows &names
                └── n borrows one element at a time

loop ends → temporary borrows end → names still owns the Vec
```

### Consuming the collection

```rust
let names = vec!["Ana", "Bo", "Cy"];

for n in names {           // CONSUME names — it's moved into the loop
    println!("{n}");
}
// println!("{}", names.len()); // ❌ names was moved away
```

This time the expression after `in` is `names` itself, not `&names`. `Vec` does
not implement `Copy`, so giving it to the loop **moves ownership** into the
vector's iterator. The iterator then yields owned elements one at a time. After
the loop, the iterator and remaining vector storage are dropped. The binding
`names` still exists as a name, but it no longer owns a value, so Rust refuses
to let you call `.len()` on it.

Conceptually, the loop behaves roughly like this:

```rust
let mut iterator = names.into_iter(); // names is moved here
while let Some(n) = iterator.next() {
    println!("{n}");
}
// iterator is dropped; names cannot be used again
```

The compiler error often says **“borrow of moved value: `names`.”** The word
“borrow” refers to what `.len()` tries to do: method calls borrow their receiver
temporarily. That borrow is impossible because `names` lost ownership earlier
at `for n in names`.

The difference is therefore not caused by `println!` or by the loop body. It is
caused by the expression given to the loop:

```
  for n in &names   → iterator borrows the Vec → names survives
  for n in names    → iterator owns the Vec    → names is moved
```

Use `&names` when you only need to read and want to keep the collection. Use
`names` when you are finished with the collection and want the loop to take its
elements. A third form, `&mut names`, temporarily borrows the vector mutably and
lets the loop edit each element in place.

## Labeled breaks for nested loops

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

## Expressions, statements, and `()`

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

## Try it

1. Write an `if` expression that stores either `"small"` or `"big"` in a variable.
2. Use a `for` loop to print the numbers 1 through 5.
3. Create a `Vec` and loop over it once with `&v` and once by consuming `v`.

> **Takeaway:** lean on Rust's expression orientation — `let x = if …` and
> `let x = loop { … break v }` replace clumsy mutable temporaries. And decide
> deliberately whether a `for` loop should **borrow** (`&`) or **consume** its
> collection.
