---
created: "2026-07-08"
id: rust-b-control
source: meari-course
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
afterward — this is your first brush with [[Ownership & Moves]]:

```rust
let names = vec!["Ana", "Bo", "Cy"];

for n in &names {          // BORROW each item — names is still usable after
    println!("{n}");
}
println!("{}", names.len()); // ✅ still fine

for n in names {           // CONSUME names — it's moved into the loop
    println!("{n}");
}
// println!("{}", names.len()); // ❌ names was moved away
```

A quick mental picture of the difference:

```
  for n in &names   →   loop borrows a view;   names survives
  for n in names    →   loop takes ownership;  names is gone
```

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

## Try it

1. Write an `if` expression that stores either `"small"` or `"big"` in a variable.
2. Use a `for` loop to print the numbers 1 through 5.
3. Create a `Vec` and loop over it once with `&v` and once by consuming `v`.

> **Takeaway:** lean on Rust's expression orientation — `let x = if …` and
> `let x = loop { … break v }` replace clumsy mutable temporaries. And decide
> deliberately whether a `for` loop should **borrow** (`&`) or **consume** its
> collection.
