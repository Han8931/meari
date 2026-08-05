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
  hint: |
    Chain off `xs.iter()`: map each element through a closure that squares it, then collapse the results into one total. Let the return type drive inference — no manual loop or mutable accumulator needed.
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

### Later: how closures capture

A closure borrows or takes what it uses, following the same
[[References & Borrowing|borrowing]] rules. The trait it implements, though, is
decided by what its **body does** to the captured values — not by how they were
captured:

| Trait    | The body…                     | Callable       |
| -------- | ----------------------------- | -------------- |
| `Fn`     | only *reads* captures         | many times     |
| `FnMut`  | *mutates* captured state      | many times     |
| `FnOnce` | *consumes*/moves captures out | at most once* |

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

Do not memorize the trait names on a first pass. Start with the behavior:

```
body only reads captured state    → can call repeatedly
body changes captured state       → closure binding may need `mut`
body moves captured state away    → that closure call consumes it
```

The trait names become useful when a function accepts a closure as a parameter.

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

## Later: `iter`, `iter_mut`, and `into_iter`

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

## Later: transforming `Option` and `Result`

The same closure syntax works with the wrappers from [[Option & Result]].
These methods run the closure only when the appropriate inner value exists:

```rust
let len: Option<usize> =
    Some("hello").map(|text| text.len());       // Some(5)

let doubled: Option<i32> =
    Some(4).filter(|n| *n > 0).map(|n| n * 2); // Some(8)

let absent: Option<i32> =
    None.map(|n: i32| n * 2);                  // None; closure never runs
```

| Method | Mental expansion |
| ------ | ---------------- |
| `option.map(f)` | `Some(x) → Some(f(x))`; `None → None` |
| `result.map(f)` | `Ok(x) → Ok(f(x))`; `Err(e) → Err(e)` |
| `and_then(f)` | run a step that itself returns `Option` or `Result` |
| `unwrap_or(x)` | use the inner success value or a fallback |

Start with `match` when the flow is unclear. Reach for these methods after you
can expand them back into both variants in your head.

## Syntax checkpoint

Read this pipeline from left to right:

```rust
values.iter().map(|x| x * 2).sum()
```

1. `.iter()` borrows the values one at a time.
2. `|x| x * 2` is a small unnamed function; the bars surround its input.
3. `.map(...)` applies that function to each item lazily.
4. `.sum()` drives the iterator and combines its outputs.

Start with `map` and one consumer. Add `filter`, capture rules, and closure trait
names only after this shape is comfortable. You are ready to continue when you
can translate a simple `for` loop into `.iter().map(...).collect()` or `.sum()`.
The final core lesson organizes the functions and types you have learned into
modules and introduces external crates.

## Try it

1. **Trace:** Predict how many times a closure runs in
   `(1..=5).filter(...).take(1).collect()`.
2. **Fill:** Write a closure `|x| x * 2` and call it with a number.
3. **Transform:** Use `.iter().map(...).collect()` to double a vector.
4. **Compare ownership:** Run equivalent pipelines starting with `.iter()` and
   `.into_iter()`, then check whether the original vector survives.
5. **Expand:** Rewrite one `Option::map` call as a `match`.

> **Takeaway:** closures are anonymous functions whose bodies determine whether
> they read, mutate, or consume captured state; `move` changes captures to owned
> captures. Iterators are lazy pipelines where adapters build a recipe and a
> consumer runs it.
