---
created: "2026-07-08"
id: rust-b-borrowing
source: meari-course
study:
  answer: |
    fn add_exclamation(text: &mut String) {
        text.push('!');
    }
  kind: code
  lang: rust
  prompt: 'Write `add_exclamation(text: &mut String)` so it mutably borrows the caller''s string and appends one `!` without taking ownership.'
  starter: |
    fn add_exclamation(text: &mut String) {
    }
  tests:
    - 'let mut s = String::from("hello"); add_exclamation(&mut s); assert_eq!(s, "hello!");'
    - 'let mut s = String::new(); add_exclamation(&mut s); assert_eq!(s, "!");'
subject: Rust (Beginner)
title: References & Borrowing
---

[[Ownership & Moves]] showed that passing a value hands over ownership — which
would make sharing data painful. **Borrowing** is the fix: a reference lets you
*use* a value without *owning* it, like lending a book instead of giving it away.

## Creating a reference with `&`

```rust
fn length(s: &str) -> usize {      // borrows text, doesn't take ownership
    s.len()
}

let name = String::from("Ana");
let n = length(&name);             // lend name as a &str
println!("{name} has {n} chars");  // ✅ name is STILL OWNED here
```

```
   name ●──────► "Ana"      (name owns the data)
                  ▲
   &name ─────────┘         (a reference borrows it, owns nothing)
```

When the reference goes out of scope, nothing is freed — you only borrowed it.

> `&String` would also borrow the string, but `&str` is the more flexible
> parameter type for read-only text: it accepts both string literals and borrowed
> `String`s. [[String vs &str]] explains the distinction in detail.

## Shared vs mutable references

There are two kinds of borrow:

| Reference | Written  | Grants          | How many at once      |
| --------- | -------- | --------------- | --------------------- |
| Shared    | `&T`     | read-only       | any number            |
| Mutable   | `&mut T` | read **&** write | exactly one           |

```rust
fn append(s: &mut String) {        // mutable borrow
    s.push_str(" Smith");
}

let mut name = String::from("Ana");
append(&mut name);                 // must be `mut name` and `&mut`
println!("{name}");                // "Ana Smith"
```

## The borrow rules

At any given time, for a piece of data you may have **either**:

```
   ┌─────────────────────────────────────────────┐
   │   MANY shared readers   &T   &T   &T   …     │
   │              — OR —                          │
   │   ONE exclusive writer  &mut T               │
   │                                              │
   │   never both at the same time.               │
   └─────────────────────────────────────────────┘
```

This "shared XOR mutable" rule is what prevents data races *and* subtle
aliasing bugs. Violate it and the compiler stops you:

```rust
let mut v = String::from("hi");
let r1 = &v;
let r2 = &mut v;        // ❌ cannot borrow `v` as mutable:
                        //    it's already borrowed as immutable
println!("{r1}");
```

## The Python contrast

Python already passes objects "by reference," but with **no borrow checker** —
any number of parts of your program can hold and mutate the same list at the
same time:

```python
def append(s):        # receives a reference to the caller's list
    s.append("Smith")

name = ["Ana"]
append(name)          # mutates the caller's data; no &mut, no rules
print(name)           # ['Ana', 'Smith']
```

Convenient — but that unrestricted aliasing is exactly what allows subtle bugs
and, in threaded code, data races. Rust's "many readers **or** one writer" rule
forbids the dangerous combination at compile time, which is the whole basis of
its fearless concurrency.

## No dangling references

Rust also refuses to let a reference outlive the data it points to — a dangling
pointer is a compile error, not a runtime crash:

```rust
fn dangle() -> &String {           // ❌ returns a reference to…
    let s = String::from("oops");
    &s                             // …s, which is dropped when the fn ends
}
```

The fix is to return the owned value instead (`-> String`, return `s`), handing
ownership to the caller. This "a borrow can't outlive its owner" guarantee is
enforced by **lifetimes**; as a beginner the compiler usually infers them for
you, and the error messages tell you when it can't.

## Choosing a signature

A practical rule of thumb for function parameters:

| You want to…                        | Take…       |
| ----------------------------------- | ----------- |
| read the value                      | `&T` / `&str` for text |
| modify the caller's value in place  | `&mut T`    |
| take/store/consume the value        | `T` (owned) |

Prefer borrowing (`&T`) by default — it's the least restrictive on your caller.

## A borrow has a usable region

A borrow usually lasts until its **last use**, not necessarily until the closing
brace. This is why the following is accepted:

```rust
let mut text = String::from("hi");
let view = &text;
println!("{view}"); // last use of the shared borrow
text.push('!');     // mutable access is now safe
```

If you add another `println!("{view}")` after `push`, the program is rejected.
When debugging a borrow error, look for the reference's last use, not only where
it was created.

## References do not own or clone

`&value` creates temporary access to the same value. It neither transfers
ownership nor duplicates data. Dereferencing with `*` means “access the value
behind this reference”:

```rust
fn increment(n: &mut i32) {
    *n += 1;
}

let mut count = 4;
increment(&mut count);
println!("{count}"); // 5
```

For method calls such as `s.len()`, Rust often inserts the needed borrowing or
dereferencing automatically. That convenience is why `*` appears less often
than you might expect.

## Try it

1. Write a function that takes `&str` and returns its length.
2. Write a function that takes `&mut String` and appends an exclamation mark.
3. Try to create one shared borrow and one mutable borrow of the same value at the same time. Read the error.

> **Takeaway:** borrowing lets many parts of a program read shared data, or one
> part mutate it, but never both at once — and never past the data's lifetime.
> These rules, checked at compile time, are the whole basis of Rust's
> "fearless" safety.
