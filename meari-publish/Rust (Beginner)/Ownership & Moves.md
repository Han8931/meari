---
created: "2026-07-08"
id: rust-b-ownership
source: meari-course
study:
  answer: |
    fn longer(a: String, b: String) -> String {
        if a.len() >= b.len() {
            a
        } else {
            b
        }
    }
  kind: code
  lang: rust
  prompt: 'Write `longer(a: String, b: String) -> String` that takes ownership of both strings and returns whichever is longer (return `a` on a tie).'
  starter: |
    fn longer(a: String, b: String) -> String {
        a
    }
  tests:
    - assert_eq!(longer(String::from("hi"), String::from("hello")), "hello");
    - assert_eq!(longer(String::from("abc"), String::from("de")), "abc");
    - assert_eq!(longer(String::from("ab"), String::from("cd")), "ab");
subject: Rust (Beginner)
title: Ownership & Moves
---

Ownership is one of Rust's most important ideas, and it may feel unfamiliar at
first. There is no need to understand every consequence immediately. Begin with
one question: **which variable is responsible for this value right now?** Rust
uses the answer to clean up memory safely, without a garbage collector or a
manual `free` call.

## The three rules

```
  1. Every value has exactly ONE owner.
  2. There can only be one owner at a time.
  3. When the owner goes out of scope, the value is DROPPED (freed).
```

That's the entire model. The compiler enforces it, so cleanup is automatic *and*
provably correct.

## Scope and `drop`

When a value's owner leaves its block, Rust automatically frees the value — no
`free()`, no `delete`, no GC:

```rust
fn main() {
    let s = String::from("hello"); // s owns a heap-allocated string
    println!("{s}");
}                                  // s goes out of scope → memory freed here
```

## Stack vs heap (why moves exist)

Simple values (`i32`, `bool`, `char`) sit entirely on the **stack** — cheap to
copy. A `String` or `Vec` is different: a small handle on the stack points to
data on the **heap**.

```
   let s = String::from("hi");

   STACK              HEAP
   ┌───────────┐      ┌───┬───┐
   │ ptr   ●───┼────► │ h │ i │
   │ len   2   │      └───┴───┘
   │ cap   2   │
   └───────────┘
```

## A move transfers ownership

Assigning a heap value to another variable does **not** copy the heap data — it
*moves* the handle, and the old variable becomes invalid:

```rust
let s1 = String::from("hello");
let s2 = s1;            // the handle MOVES from s1 to s2

println!("{s2}");      // ✅ fine
println!("{s1}");      // ❌ error: value borrowed after move
```

```
   before:  s1 ●──► "hello"        after `let s2 = s1;`

   s1 ✗ (invalidated)
   s2 ●──► "hello"       ← only ONE owner, as rule #2 demands
```

Why invalidate `s1`? If both `s1` and `s2` pointed at the same heap buffer, then
when *both* went out of scope Rust would try to free it **twice** — a classic
double-free bug. Moving prevents that by construction.

## The Python contrast — the biggest one in this course

This is where Rust departs most sharply from Python. In Python, assignment just
creates another *name* for the **same** object; both names stay valid, and a
garbage collector frees the object later when nothing references it:

```python
s1 = "hello"
s2 = s1          # s1 and s2 are two names for one object
print(s1, s2)    # ✅ both fine — Python never "moves" or invalidates a name
```

Python can afford this because a garbage collector is always running to decide
when memory is safe to free. Rust has no GC, so it uses ownership instead: the
move *invalidates* `s1` precisely so the heap buffer has exactly one owner and
gets freed exactly once — no GC required, no double-free possible.

## `Copy` types don't move

Types that live entirely on the stack implement the `Copy` trait, so assignment
duplicates them and the original stays valid:

```rust
let x = 5;
let y = x;             // x is COPIED, not moved
println!("{x} {y}");   // ✅ both usable — 5 5
```

| Behavior on assign | Types                                         |
| ------------------ | --------------------------------------------- |
| **Copy** (cheap)   | `i32`, `f64`, `bool`, `char`, tuples of Copy  |
| **Move** (handle)  | `String`, `Vec<T>`, most heap-owning types    |

## `clone()` for an explicit deep copy

When you genuinely want two independent owners of heap data, ask for it
explicitly with `.clone()` — the visible `.clone()` call is Rust telling you
"this costs a heap copy":

```rust
let s1 = String::from("hello");
let s2 = s1.clone();   // deep copy — s1 and s2 each own their own buffer
println!("{s1} {s2}"); // ✅ both valid
```

## Moves happen at function boundaries too

Passing a value to a function moves it in, unless the type is `Copy`:

```rust
fn takes(s: String) { println!("{s}"); }  // s is dropped when this returns

let name = String::from("Ana");
takes(name);
// println!("{name}");   // ❌ name was moved into takes()
```

Moving a value into every function would make many ordinary programs awkward.
**Borrowing** provides another choice: a function can use a value temporarily
without becoming its owner. That is the subject of the next lesson,
[[References & Borrowing]].

## A more precise mental model

“Stack types copy, heap types move” is a useful first approximation, but the
actual rule is: assignment copies only when the type implements `Copy`; otherwise
it moves. A type can contain only stack data and still choose not to be `Copy`.
The compiler follows the trait, not the storage location.

A move also does not physically move the heap bytes. Rust copies the small
`String` handle (pointer, length, and capacity) into the new variable and then
forbids use of the old handle. The buffer stays where it was.

## Trace ownership through a function

```rust
fn add_period(mut text: String) -> String {
    text.push('.'); // the function owns text and may mutate it
    text             // ownership moves back to the caller
}

let sentence = String::from("Hello");
let sentence = add_period(sentence);
println!("{sentence}");
```

The owner is first the outer `sentence`, then the parameter `text`, then the new
outer `sentence`. There is one usable owner at every point. Returning ownership
works, but borrowing is more convenient for temporary access.

### Partial moves

Moving a non-`Copy` field can make only part of a compound value unavailable:

```rust
let pair = (String::from("left"), 7);
let word = pair.0;        // moves the String field
println!("{}", pair.1);   // untouched i32 field is still usable
// println!("{:?}", pair); // error: pair is partially moved
```

You need not use partial moves yet, but recognizing one makes the compiler's
message much less mysterious.

## Try it

1. Move a `String` from `s1` to `s2`, then try to print both. Read the error.
2. Repeat the same experiment with an `i32`. Why does it still work?
3. Use `.clone()` on a `String` and confirm both variables remain usable.

> **Takeaway:** one owner, dropped at end of scope. Assignment *moves* heap
> types (invalidating the source) and *copies* stack types. Reach for `.clone()`
> only when you truly need a second independent owner.
