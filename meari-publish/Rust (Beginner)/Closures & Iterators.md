---
created: "2026-07-08"
id: rust-b-iterators
source: meari-course
study:
  answer: |
    fn sum_of_squares(xs: &[i32]) -> i32 {
        xs.iter().map(|x| x * x).sum()
    }
  kind: code
  lang: rust
  prompt: 'Write `sum_of_squares(xs: &[i32]) -> i32` using iterator adapters (`.iter().map(...).sum()`).'
  starter: |
    fn sum_of_squares(xs: &[i32]) -> i32 {
        0
    }
  tests:
    - assert_eq!(sum_of_squares(&[1, 2, 3]), 14);
    - assert_eq!(sum_of_squares(&[]), 0);
    - assert_eq!(sum_of_squares(&[-2]), 4);
subject: Rust (Beginner)
title: Closures & Iterators
---

Rust has a rich functional side. **Closures** are anonymous functions that can
capture their surroundings; **iterators** are lazy sequences you transform with
composable adapters. Together they let you replace many index loops with clear,
declarative pipelines — and thanks to zero-cost abstractions, they compile down
to code as fast as the hand-written loop.

## Closures

A closure is written with `|params| body` and can **capture** variables from the
scope where it's defined:

```rust
let add = |a: i32, b: i32| a + b;
println!("{}", add(2, 3));         // 5

let factor = 10;
let scale = |x: i32| x * factor;   // captures `factor` from the environment
println!("{}", scale(5));          // 50
```

Types are usually inferred, so closures are far terser than named functions.

### How closures capture

A closure borrows or takes what it uses, following the same
[[References & Borrowing|borrowing]] rules. The trait it implements, though, is
decided by what its **body does** to the captured values — not by how they were
captured:

| Trait    | The body…                     | Callable       |
| -------- | ----------------------------- | -------------- |
| `Fn`     | only *reads* captures         | many times     |
| `FnMut`  | *mutates* captured state      | many times     |
| `FnOnce` | *consumes*/moves captures out | at least once* |

<sub>*A `FnOnce` may be callable only once. The traits nest: every `Fn` is also
`FnMut`, and every `FnMut` is also `FnOnce`.</sub>

Separately, the `move` keyword forces a closure to *take ownership* of what it
captures — essential when the closure outlives the current scope (e.g. a new
thread):

```rust
let name = String::from("Ana");
let greet = move || println!("Hi, {name}"); // name is MOVED into the closure
greet();
greet();                                     // still fine — callable repeatedly
```

Note `move` only changes *how* captures are taken (by ownership); it does **not**
force `FnOnce`. `greet` here still implements `Fn`, because it only *reads* the
`name` it now owns.

(`Fn`, `FnMut`, and `FnOnce` are *traits* — shared-behavior contracts formally
introduced in [[Traits]].)

## Iterators

An iterator produces a sequence one item at a time. You get one from a collection
with `.iter()` (borrowing) or `.into_iter()` (consuming):

```rust
let v = vec![1, 2, 3];
let mut it = v.iter();
it.next();     // Some(&1)
it.next();     // Some(&2)
```

You rarely call `.next()` by hand. The power is in chaining **adapters**.

### Adapters are lazy

Adapters like `map` and `filter` build a *recipe* and do nothing until a
**consumer** drives them. This laziness means no intermediate collections are
built:

```rust
let v = vec![1, 2, 3, 4, 5, 6];

let result: Vec<i32> = v.iter()
    .filter(|&&x| x % 2 == 0)   // keep evens:      2, 4, 6
    .map(|&x| x * 10)           // transform:      20, 40, 60
    .collect();                 // consumer: run it, gather into a Vec

// result == [20, 40, 60]
```

```
   [1,2,3,4,5,6]
        │  .filter(even)
        ▼
     [2,4,6]              ← nothing has run yet; this is a "recipe"
        │  .map(×10)
        ▼
   [20,40,60]             ← .collect() finally DRIVES the whole pipeline
```

### Common adapters and consumers

| Adapter (lazy) | Produces                              |
| -------------- | ------------------------------------- |
| `map(f)`       | each item transformed by `f`          |
| `filter(pred)` | only items where `pred` is true       |
| `take(n)`      | the first `n` items                   |
| `zip(other)`   | pairs from two iterators              |
| `enumerate()`  | `(index, item)` pairs                 |

| Consumer (drives it) | Produces                          |
| -------------------- | --------------------------------- |
| `collect()`          | a collection (Vec, HashMap, …)    |
| `sum()` / `product()`| a single folded number            |
| `count()`            | how many items                    |
| `fold(init, f)`      | a custom accumulation             |
| `for_each(f)`        | runs `f` for its side effects     |

```rust
let total: i32 = (1..=100).sum();               // 5050
let words: Vec<&str> = "a b c".split(' ').collect();
for (i, c) in "rust".chars().enumerate() {
    println!("{i}: {c}");                       // 0: r, 1: u, …
}
```

## The same in Python

Python expresses the very same pipeline with a comprehension (or `map`/`filter`),
and `lambda` is its closure:

```python
v = [1, 2, 3, 4, 5, 6]
result = [x * 10 for x in v if x % 2 == 0]   # [20, 40, 60]

add = lambda a, b: a + b                      # ~ let add = |a, b| a + b;
total = sum(range(1, 101))                    # 5050
```

Python generators are lazy, just like Rust iterators. Two differences stand out:
Python has no `Fn`/`FnMut`/`FnOnce` distinction (closures capture by reference to
the enclosing scope), and Rust's iterator chains are **zero-cost** — they compile
down to the same machine code as the hand-written loop, with no per-item
overhead.

## Loop vs iterator

```rust
// imperative
let mut sum = 0;
for &x in &v { if x % 2 == 0 { sum += x; } }

// declarative — same result, often clearer, and just as fast
let sum: i32 = v.iter().filter(|&&x| x % 2 == 0).sum();
```

Iterator chains express *what* you want, not the bookkeeping of *how*. Prefer
them for transformations; reach for an explicit loop when the logic is genuinely
imperative or the borrow interplay gets awkward.

## `iter`, `iter_mut`, and `into_iter`

These three starting points control ownership:

| Call | Items yielded | Effect on collection |
| --- | --- | --- |
| `values.iter()` | `&T` | shared borrow; collection survives |
| `values.iter_mut()` | `&mut T` | mutable borrow; edit in place |
| `values.into_iter()` | `T` | consumes the collection |

This is the ownership lesson in iterator form. If a closure receives extra `&`
characters, first ask which iterator you created. For beginners, splitting a
long chain into named intermediate values and adding a type annotation often
makes the compiler message easier to understand.

`collect()` also needs to know its destination. Either annotate the variable
(`let result: Vec<_> = ...`) or use `collect::<Vec<_>>()`; `_` asks Rust to infer
the element type while you specify the container.

## Try it

1. Write a closure `|x| x * 2` and call it with a number.
2. Use `.iter().map(...).collect()` to double every number in a vector.
3. Use `.filter(...)` to keep only even numbers.

> **Takeaway:** closures are capturing anonymous functions (`Fn`/`FnMut`/`FnOnce`
> by how they capture; `move` to take ownership); iterators are lazy pipelines
> where adapters build a recipe and a consumer runs it — expressive *and*
> zero-cost.
