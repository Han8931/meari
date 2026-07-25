---
created: '2026-07-21'
id: rust-i-fn-traits
source: meari-course
study:
  answer: |
    fn repeat<F>(times: usize, mut action: F)
    where
        F: FnMut(),
    {
        for _ in 0..times {
            action();
        }
    }

    fn transform_once<F>(text: String, transform: F) -> usize
    where
        F: FnOnce(String) -> usize,
    {
        transform(text)
    }
  kind: code
  lang: rust
  prompt: 'Write `repeat<F>(times: usize, action: F)` using `F: FnMut()` and call `action` exactly `times` times. Then write `transform_once<F>(text: String, transform: F) -> usize` using `F: FnOnce(String) -> usize` and return `transform(text)`.'
  starter: |
    fn repeat<F>(times: usize, action: F)
    where
        F: FnMut(),
    {
        // call action `times` times
    }

    fn transform_once<F>(text: String, transform: F) -> usize
    where
        F: FnOnce(String) -> usize,
    {
        0
    }
  tests:
  - 'let mut count = 0; repeat(4, || count += 1); assert_eq!(count, 4);'
  - 'let mut values = vec![]; repeat(3, || values.push(values.len())); assert_eq!(values, vec![0, 1, 2]);'
  - 'let mut count = 0; repeat(0, || count += 1); assert_eq!(count, 0);'
  - 'assert_eq!(transform_once(String::from("rust"), |s| s.into_bytes().len()), 4);'
  - 'let captured = String::from("owned"); assert_eq!(transform_once(String::from("x"), move |s| { drop(captured); s.len() }), 1);'
subject: Rust (Intermediate)
title: Fn, FnMut & FnOnce
---

[[rust-b-iterators|Closures & Iterators]] introduced the three closure traits.
Here the goal is to use them when designing an API: the bound you choose tells a
caller what your function may do with their closure.

## The three closure traits

- **`FnOnce`** permits one call that may consume captured state.
- **`FnMut`** permits repeated calls and may mutate captured state.
- **`Fn`** permits repeated calls without mutating or consuming captured state.

The traits nest: every `Fn` closure can be used as `FnMut`, and every `FnMut`
closure can be used as `FnOnce`. Therefore `FnOnce` is the least restrictive
bound—it accepts the widest set of closures. Require only what your implementation
needs.

Picture the capabilities as nested sets:

```
Fn closures       ⊂ FnMut closures       ⊂ FnOnce closures
read repeatedly     may mutate repeatedly  may consume on one call
```

“Least restrictive bound” describes what the **caller may pass**, not how often
your function may call it. A parameter bounded only by `FnOnce` can be called no
more than once because calling it consumes the closure value.

## Taking a closure as a parameter

Use a trait bound, exactly like any other generic. Calling through `FnMut`
requires the parameter binding itself to be mutable:

```rust
fn repeat<F>(times: usize, mut action: F)
where
    F: FnMut(),
{
    for _ in 0..times {
        action();
    }
}

let mut count = 0;
repeat(3, || count += 1);
assert_eq!(count, 3);
```

## Consuming through `FnOnce`

Use `FnOnce` when your function calls a callback once and wants to allow that
callback to consume a value or captured state:

```rust
fn transform<F>(text: String, f: F) -> usize
where
    F: FnOnce(String) -> usize,
{
    f(text)
}

let len = transform(String::from("rust"), |s| s.into_bytes().len());
```

The argument and captured environment are separate:

- `FnOnce(String) -> usize` says the closure receives one owned `String`
  argument.
- A `move` closure may additionally own captured values.
- The body determines whether calling it consumes captured state and therefore
  whether it is only `FnOnce`.

The `move` keyword alone does not force `FnOnce`; a closure can own a capture and
only read it on every call.

## Your turn

1. **Classify:** For closures that read, mutate, and move out of a capture, choose
   the narrowest implemented call trait.
2. **Explain nesting:** State why an `Fn` closure can satisfy an `FnOnce` bound,
   but not vice versa.
3. **Repair:** Try calling an `FnOnce` parameter twice and connect the error to
   ownership of the closure.
4. **Implement:** Build `repeat` for state-mutating callbacks and
   `transform_once` for a callback allowed to consume state.
