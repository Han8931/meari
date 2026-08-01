---
created: "2026-08-01"
id: rust-b-syntax
source: meari-course
study:
  answer: |
    fn next_score(score: i32) -> i32 {
        let bonus = 1;
        score + bonus
    }
  kind: code
  lang: rust
  prompt: 'Complete `next_score`. Make a variable named `bonus` with value `1`, then return `score + bonus` as the final line (no semicolon on that line).'
  starter: |
    fn next_score(score: i32) -> i32 {
        0
    }
  tests:
    - assert_eq!(next_score(0), 1);
    - assert_eq!(next_score(9), 10);
subject: Rust (Beginner)
title: Reading Rust Code, One Piece at a Time
---

The `Hello, world!` program contains a lot of unfamiliar punctuation. That is
normal. You do **not** need to understand it all before writing code. This
lesson gives you just enough syntax to read, change, and run tiny programs.
Later lessons revisit each piece in more depth.

Start with one complete program:

```rust
fn main() {
    let name = "Mina";
    println!("Hello, {name}!");
}
```

Read it from the outside in:

1. `fn main()` says “define the function named `main`.” Rust starts a program
   by calling this particular function.
2. `{` starts its body and `}` ends it. The indented lines are inside `main`.
3. `let name = "Mina";` gives the name `name` the text value `"Mina"`.
4. `println!(...)` prints one line. Text inside double quotes is a **string
   literal**. `{name}` asks `println!` to place the value of `name` there.

For now, think of braces as “inside this part,” parentheses as “the inputs to
this call,” and a semicolon as “this instruction is finished.”

## The smallest useful shapes

You will see these shapes repeatedly:

```rust
let number = 3;             // create a named value
println!("{number}");       // use it

fn double(n: i32) -> i32 {  // define a reusable action
    n * 2                    // its result
}

let result = double(3);     // call it
```

Do not worry about `i32` or `-> i32` yet. Read them as: “`double` takes a whole
number and produces a whole number.” [[Data Types & Type Casting]] and
[[Functions & Expressions]] will make that precise.

The two kinds of line above have one important difference:

```rust
let number = 3; // this instruction ends with `;`
n * 2           // no `;`: this is the value returned by `double`
```

A semicolon usually ends an instruction. In a function, the last line without a
semicolon is often the value the function gives back. You only need to notice
that pattern now; the functions lesson explains why it works.

## Change one thing, then run it

Make a small edit rather than writing a program from memory:

```rust
fn main() {
    let name = "Mina";
    println!("Hello, {name}!");
}
```

Try these in order:

1. Change `"Mina"` to your name and run `cargo run`.
2. Add `let score = 10;` above `println!`, then print `"score: {score}"`.
3. Remove the semicolon after `let score = 10`, run `cargo check`, and read the
   first error. Put it back.

It is fine if the compiler's wording is not clear yet. The useful observation is
that a small punctuation change gives a precise location to inspect.

## A tiny syntax map

| You see | Read it as | Remember for now |
| --- | --- | --- |
| `fn name() { ... }` | define a function | `main` is where a program begins |
| `let name = value;` | make a named value | the semicolon belongs here |
| `"text"` | text written directly in code | use double quotes for text |
| `call(value)` | call something with an input | parentheses hold the input |
| `{ ... }` | a block of code | braces mark its beginning and end |
| `println!(...)` | print a line | copy this shape when experimenting |

Comments begin with `//`. Rust ignores everything after those two slashes on a
line, so they are a safe place to leave yourself a note:

```rust
let score = 10; // the starting score
```

## Before moving on

Make sure you can read—not necessarily reproduce—these two lines:

```rust
let score = 10;          // make the name score refer to 10
println!("{score}");     // print the value named score
```

The next lesson changes `score`. That introduces exactly one new word, `mut`;
every other part of the example above stays familiar.

## Try it

1. Point to the function body, the variable name, and the string literal in the
   first program.
2. Predict what changing `name` changes before you run the program.
3. Add a second `println!` line that prints a short message.
4. In the exercise, use the `let` shape and leave the final calculation without
   a semicolon.

> **Takeaway:** you only need a few reading rules to begin: `fn main` starts a
> program, braces surround a block, `let` names a value, and semicolons finish
> ordinary instructions. The next lessons fill in the meaning gradually.
