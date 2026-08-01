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

When a value's owner leaves its block, Rust runs that value's cleanup — no
manual `free()` or `delete`, and no garbage collector deciding when to act:

```rust
fn main() {
    let s = String::from("hello"); // s owns a heap-allocated string
    println!("{s}");
}                                  // s goes out of scope → memory freed here
```

## Stack vs heap (why moves exist)

The stack and heap explain why some values need cleanup, but they do **not**
decide whether assignment copies or moves.

- An `i32` contains its whole value directly and needs no heap allocation.
- A `String` value is a small handle—pointer, length, and capacity—that owns a
  separate byte buffer on the heap.
- A user-defined type may contain only direct fields and still choose move
  behavior. The actual rule is whether the type implements `Copy`.

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

Assigning a non-`Copy` value to another variable transfers responsibility. Rust
calls this a **move**. For a `String`, the small handle is copied internally, the
heap bytes stay in place, and the old variable becomes unusable:

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

Why invalidate `s1`? If both handles independently believed they owned the same
buffer, both would try to clean it up—a double free. Rust instead transfers the
one cleanup responsibility to `s2`.

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
