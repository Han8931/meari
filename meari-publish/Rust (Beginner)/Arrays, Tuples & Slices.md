---
created: "2026-07-08"
id: rust-b-compound
source: meari-course
study:
  answer: |
    fn first_and_last(xs: &[i32]) -> (i32, i32) {
        (xs[0], xs[xs.len() - 1])
    }
  kind: code
  lang: rust
  prompt: 'Complete `first_and_last(xs: &[i32]) -> (i32, i32)`. For this non-empty slice, use `xs[0]` for the first value and `xs[xs.len() - 1]` for the last, then return them as a tuple.'
  hint: |
    A slice knows its own length via `.len()`, but indexing starts at 0 — so the last element sits one before that count. Wrap the two values in parentheses to return a tuple.
  starter: |
    fn first_and_last(xs: &[i32]) -> (i32, i32) {
        (0, 0)
    }
  tests:
    - assert_eq!(first_and_last(&[3, 1, 2]), (3, 2));
    - assert_eq!(first_and_last(&[5]), (5, 5));
    - assert_eq!(first_and_last(&[-2, 4, 0]), (-2, 0));
subject: Rust (Beginner)
title: Arrays, Tuples & Slices
---

Scalars hold one value; **compound types** group several. This lesson starts
with arrays and tuples, which have a fixed shape, then introduces slices, which
are borrowed views into a sequence. (The growable cousins — `Vec` and `HashMap`
— get their own lesson in [[Vec & HashMap]].)

## Arrays: same type, fixed length

An array's type is written `[T; N]` — element type `T`, exactly `N` of them. The
length is part of the type and fixed at compile time. That means `[i32; 3]` and
`[i32; 4]` are different types; a function that requires one cannot receive the
other. Use a slice when the function should accept either length.

```rust
let days: [i32; 3] = [1, 2, 3];
let zeros = [0u8; 4];          // [0, 0, 0, 0] — shorthand for repeats

println!("{}", days[0]);       // 1 — indexing
println!("{}", days.len());    // 3
```

For now, treat an array as one local value with all of its elements together.
Every index is **bounds-checked** at runtime —
reading `days[9]` panics rather than reading random memory (that's the safety
guarantee in action):

```
  stack:  days = [ 1 | 2 | 3 ]
                    ↑
                  days[0]        days[9] → panic! index out of bounds
```

## Tuples: mixed types, fixed length

A tuple groups a fixed number of values that can each be a **different type**:

```rust
let person: (&str, i32, bool) = ("Ana", 30, true);

// access by position with .0, .1, .2
println!("{}", person.0);      // "Ana"

// or destructure into names
let (name, age, active) = person;
println!("{name} is {age}");
```

The empty tuple `()` is called the **unit type**. It means "no meaningful
value" and is what expressions like a `println!` statement evaluate to — Rust's
equivalent of `void`.

## Slices: a borrowed window

A slice `&[T]` is a **reference to a contiguous run** of an array or vector — a
view, not a copy. You don't own the data; you borrow a range of it. Conceptually,
a slice carries two facts: where the first element is and how many elements are
in the view. The elements remain in the original collection.

```rust
let nums = [10, 20, 30, 40, 50];
let mid = &nums[1..4];         // &[20, 30, 40]

println!("{}", mid.len());     // 3
println!("{}", mid[0]);        // 20
```

```
  nums:  [ 10 | 20 | 30 | 40 | 50 ]
                └─────────┘
   &nums[1..4]  points here (len 3), owns nothing
```

The range `1..4` includes index 1 but excludes index 4, so it selects indexes 1,
2, and 3. This “start inclusive, end exclusive” rule makes the slice length
`end - start`.

Because a slice is a borrow, it cannot outlive the array or vector it points
into, and the ordinary rules from [[References & Borrowing]] still apply. A
**mutable** slice `&mut [T]` lets you edit the underlying elements in place:

```rust
let mut data = [3, 1, 2];
let s = &mut data[..];         // whole-array mutable slice
s.sort();                       // data is now [1, 2, 3]
```

## The same in Python

Python's slicing syntax looks almost identical, but there's a crucial semantic
difference — a Python slice **copies** the elements into a new list, while a Rust
slice only *borrows* a view of the existing data:

```python
nums = [10, 20, 30, 40, 50]
mid = nums[1:4]            # [20, 30, 40] — a NEW list (a copy, owns its data)
```

Python also doesn't split its sequences by type and mutability the way Rust
splits arrays, tuples, and `Vec`: a Python `list` is closest to a Rust `Vec`,
and a Python `tuple` to a Rust tuple.

## Array vs Tuple vs Slice at a glance

| Type       | Written   | Elements       | Length       | Owns data? |
| ---------- | --------- | -------------- | ------------ | ---------- |
| Array      | `[T; N]`  | all same `T`   | fixed, known | yes        |
| Tuple      | `(A, B…)` | may differ     | fixed, known | yes        |
| Slice      | `&[T]`    | all same `T`   | runtime len  | **no** (borrows) |

Why doesn't a slice know its length at compile time? Because it can point into a
region of *any* size — that flexibility is exactly why slices are the standard
way to pass "some sequence of `T`" into a function without caring whether the
caller had an array or a [[Vec & HashMap|Vec]].

## Reading the type notation

Read `[i32; 4]` as “four owned `i32` values” and `&[i32]` as “a borrowed view of
some number of `i32` values.” The semicolon in an array type separates element
type from length; a slice omits the length because it is stored at runtime.

```rust
fn sum(values: &[i32]) -> i32 {
    let mut total = 0;
    for &value in values {
        total += value;
    }
    total
}

let array = [1, 2, 3];
let vector = vec![4, 5, 6];
println!("{} {}", sum(&array), sum(&vector));
```

One function accepts both collections without copying either. Rust automatically
turns `&array` and `&vector` into the `&[i32]` view the parameter needs. No
elements are copied.

Indexing and slicing are checked at runtime. `xs[0]` panics when `xs` is empty,
and `&xs[start..end]` panics unless `start <= end <= xs.len()`. This is why the
study exercise explicitly promises a non-empty slice: that promise makes its
first and last indexes valid. Later, `xs.first()` provides a non-panicking way
to handle possibly empty input by returning an `Option`.

## Syntax checkpoint

These brackets have related but different jobs:

```rust
let values: [i32; 3] = [10, 20, 30];
let view: &[i32] = &values[0..2];
```

Read `[i32; 3]` as “an array of exactly three `i32` values.” Read `&[i32]` as
“a borrowed slice containing some number of `i32` values.” In `[0..2]`, the
ending index 2 is excluded, so the view contains positions 0 and 1.

You are ready to continue when you can index an array, destructure a tuple, and
recognize that a slice borrows rather than copies. The next lesson gives names
to groups of fields so code can say `user.name` instead of `person.0`.

## Try it

1. **Classify:** For `[i32; 4]`, `(i32, bool)`, and `&[i32]`, state whether the
   type owns data and whether its length is part of the type.
2. **Build:** Make an array of five numbers and print its first element.
3. **Destructure:** Create `("Ana", 30)` and bind its fields to two names.
4. **Borrow:** Take `&arr[1..3]`, then explain the indexes and resulting length.

> **Takeaway:** arrays and tuples are fixed-shape, owned bundles; slices are
> cheap borrowed windows into sequences. Prefer accepting a slice `&[T]`
> in function signatures — it's the most general, allocation-free choice.
