---
created: "2026-07-08"
id: rust-b-compound
source: meari-course
subject: Rust (Beginner)
title: Arrays, Tuples & Slices
---

Scalars hold one value; **compound types** group several. Rust's three
fixed-shape building blocks are arrays, tuples, and slices. (The growable
cousins — `Vec` and `HashMap` — get their own lesson in [[Vec & HashMap]].)

## Arrays: same type, fixed length

An array's type is written `[T; N]` — element type `T`, exactly `N` of them. The
length is part of the type and fixed at compile time.

```rust
let days: [i32; 3] = [1, 2, 3];
let zeros = [0u8; 4];          // [0, 0, 0, 0] — shorthand for repeats

println!("{}", days[0]);       // 1 — indexing
println!("{}", days.len());    // 3
```

A local array lives **on the stack** (like any value, though, an array lives
wherever its owner does — inside a `Box`, `Vec`, or heap struct it rides on the
heap). Every index is **bounds-checked** at runtime —
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
view, not a copy. You don't own the data; you borrow a range of it:

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

Because a slice is a borrow, it follows the borrowing rules from
[[References & Borrowing]]. A **mutable** slice `&mut [T]` lets you edit the
underlying elements in place:

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

## Try it

1. Make an array of five numbers and print its first element.
2. Create a tuple like `("Ana", 30)` and destructure it into two variables.
3. Take a slice of an array with `&arr[1..3]` and print its length.

> **Takeaway:** arrays and tuples are fixed, owned, stack-friendly bundles;
> slices are cheap borrowed windows into them. Prefer accepting a slice `&[T]`
> in function signatures — it's the most general, allocation-free choice.
