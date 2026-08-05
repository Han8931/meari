---
created: "2026-07-08"
id: rust-b-ownership
source: meari-course
study:
  answer: |
    fn return_to_caller(text: String) -> String {
        text
    }
  kind: code
  lang: rust
  prompt: 'Complete `return_to_caller(text: String) -> String` by returning `text`. The `String` moves into the function and then moves back to the caller.'
  hint: |
    `text` was moved in through the parameter, so you already own it — just name it on the last line to move it back out. No `String::new()` and no `.clone()`.
  starter: |
    fn return_to_caller(text: String) -> String {
        String::new()
    }
  tests:
    - assert_eq!(return_to_caller(String::from("hello")), "hello");
    - assert_eq!(return_to_caller(String::new()), "");
subject: Rust (Beginner)
title: Ownership & Moves
---

Ownership is one of Rust's most important ideas, and it may feel unfamiliar at
first. There is no need to understand every consequence immediately. Begin with
one question: **which variable is responsible for this value right now?** Rust
uses the answer to clean up memory safely, without a garbage collector or a
manual `free` call.

## The three rules for ordinary owned values

```
  1. Every owned value has a variable or field responsible for it.
  2. That responsibility moves; it is not silently duplicated.
  3. When the owner goes out of scope, the value is DROPPED.
```

Call that responsible variable or field the **owner**. Most of this course uses
one clear owner at a time. Optional smart pointers such as `Rc` can coordinate
shared ownership explicitly; they do not invalidate this starting model.

## Scope and `drop`

A **scope** is a region of code, usually surrounded by `{` and `}`. A variable
can normally be used from its `let` statement until the end of that scope.
When its owner leaves the scope, Rust **drops** the value: Rust runs whatever
cleanup that value needs.

```rust
fn main() {
    let s = String::from("hello"); // s owns the String
    println!("{s}");
}                                  // s leaves scope; Rust drops its String
```

For a number, cleanup has almost nothing to do. For a `String`, cleanup includes
returning the memory used for its text. Rust inserts that cleanup automatically;
you do not call `free()` or wait for a garbage collector.

## A small memory detour: stack and heap

Why does a `String` have memory to return while an `i32` does not? Answering that
requires two new terms: **stack** and **heap**. They are simply two places a
program can keep data. You do not need to memorize their implementation to
follow ownership.

- The **stack** holds fixed-size data used by function calls. Rust knows how
  much space such data needs before the program runs and manages that space
  automatically.
- The **heap** holds allocations requested while the program is running. It is
  useful for data whose size can vary or grow, but each allocation eventually
  has to be returned so that its memory can be reused.

An `i32` is always four bytes, so its entire value can be stored directly with
the local variable. The text in a `String` can have any length and can grow, so
a `String` uses both places:

1. The text bytes are stored in a buffer allocated on the heap.
2. The local variable stores a small, fixed-size description of that buffer.
   This description contains a **pointer** (the buffer's address), its current
   **length**, and its **capacity** (how much room the buffer currently has).
3. The `String` owns that heap buffer, so dropping the `String` returns the
   buffer to the allocator.

For `let s = String::from("hi");`, the picture is:

```
   local variable s                 allocated buffer
   (usually on the stack)           (on the heap)

   ┌──────────────┐                 ┌───┬───┐
   │ pointer  ●───┼────────────────►│ h │ i │
   │ length   2   │                 └───┴───┘
   │ capacity 2   │
   └──────────────┘
```

The boxes are one logical value, not two separate Rust values: `s` is a
`String`, and that `String` owns its buffer. This layout helps explain what a
`String` move costs, but it is **not** the rule for deciding whether a type
moves. Assignment copies a type only if that type implements `Copy`; otherwise,
it moves. A user-defined type can therefore move even if all of its data is
stored directly in its local variable.

## A move transfers ownership

Assigning a non-`Copy` value to another variable transfers responsibility. Rust
calls this a **move**. For a `String`, the pointer, length, and capacity are
transferred to the new variable; the text bytes do not have to move to a new
place on the heap. The old variable becomes unusable:

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

Why invalidate `s1`? Imagine that Rust duplicated only the pointer, length, and
capacity. Both variables would point to the same buffer, and both would believe
that they must return it when dropped. The first cleanup would free the buffer;
the second would try to free the same buffer again. This is a **double free** and
can corrupt a program's memory. Rust prevents it by leaving only `s2` usable and
transferring the one cleanup responsibility to `s2`.

Trace the state after each line:

| Line | `s1` | `s2` | Who will free the buffer? |
| ---- | ---- | ---- | ------------------------- |
| `let s1 = String::from("hello");` | usable | absent | `s1` |
| `let s2 = s1;` | moved, unusable | usable | `s2` |
| end of scope | no longer usable | dropped | cleanup has run |

“Unusable” does not mean the name vanished or the bytes were erased. It means
the compiler will reject any path that tries to use that name as though it still
owned the value.

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

Some types implement the `Copy` trait. For them, assignment duplicates the
value and both variables remain usable:

```rust
let x = 5;
let y = x;             // x is COPIED, not moved
println!("{x} {y}");   // ✅ both usable — 5 5
```

| Behavior on assign | Decided by                         | Common examples |
| ------------------ | ---------------------------------- | --------------- |
| **Copy**           | the type implements `Copy`         | numbers, `bool`, `char`, references, tuples of `Copy` values |
| **Move**           | the type does not implement `Copy` | `String`, `Vec<T>`, many structs |

## `clone()` for explicit duplication

When you genuinely want another value, ask for it explicitly with `.clone()`.
For `String`, cloning duplicates the heap buffer:

```rust
let s1 = String::from("hello");
let s2 = s1.clone();   // deep copy — s1 and s2 each own their own buffer
println!("{s1} {s2}"); // ✅ both valid
```

`clone()` means “run this type's explicit duplication behavior,” not universally
“deep-copy every reachable object.” A later example is `Rc::clone`, which
creates another owner of the same allocation. The important promise is that
cloning is visible in your source and may do meaningful work.

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

The rule to retain is: **assignment copies only when the type implements
`Copy`; otherwise it moves.** A type can contain only direct, fixed-size data and
still choose not to be `Copy`. The compiler follows the trait, not an informal
stack-versus-heap classification.

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

### Later: partial moves

Moving a non-`Copy` field can make only part of a compound value unavailable:

```rust
let pair = (String::from("left"), 7);
let word = pair.0;        // moves the String field
println!("{}", pair.1);   // untouched i32 field is still usable
// println!("{:?}", pair); // error: pair is partially moved
```

You need not use partial moves yet, but recognizing one makes the compiler's
message much less mysterious.

## Syntax checkpoint: trace names, not memory diagrams

For a first pass, write “usable” or “moved” beside each name:

```rust
let first = String::from("hello"); // first: usable
let second = first;                // first: moved; second: usable
```

`String::from(...)` constructs owned text. The `::` selects the `from` function
associated with the `String` type. The important ownership event is the second
line: because `String` is not `Copy`, assignment transfers the value.

You are ready to continue when you can predict which name is usable after an
assignment and after passing a `String` into a function. The next lesson adds
`&`, which lets a function use that value without taking it.

## Try it

1. **Trace:** For `let a = String::from("x"); let b = a;`, write the owner after
   each statement before compiling it.
2. **Repair:** Move a `String` from `s1` to `s2`, try to print `s1`, then repair
   the program once by using only `s2` and once by cloning.
3. **Compare:** Repeat the assignment with an `i32`. Explain the result using
   the words “implements `Copy`.”
4. **Create:** Pass a `String` into a function that returns it to the caller.

> **Takeaway:** one clear owner is responsible for cleanup. Assignment copies a
> `Copy` type; otherwise it moves ownership and invalidates the source. A move
> usually transfers responsibility without moving the heap allocation itself.
