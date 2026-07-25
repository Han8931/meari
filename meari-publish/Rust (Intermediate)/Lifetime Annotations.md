---
created: '2026-07-21'
id: rust-i-lifetimes
source: meari-course
study:
  answer: |
    fn before_comma<'a>(text: &'a str, fallback: &'a str) -> &'a str {
        match text.split_once(',') {
            Some((before, _)) => before,
            None => fallback,
        }
    }
  kind: code
  lang: rust
  prompt: 'Write `before_comma<''a>(text: &''a str, fallback: &''a str) -> &''a str`. Return the part of `text` before its first comma, or `fallback` when there is no comma.'
  starter: |
    fn before_comma<'a>(text: &'a str, fallback: &'a str) -> &'a str {
        fallback
    }
  tests:
  - assert_eq!(before_comma("red,green", "none"), "red");
  - assert_eq!(before_comma("red", "none"), "none");
  - assert_eq!(before_comma(",green", "none"), "");
subject: Rust (Intermediate)
title: Lifetime Annotations
---

In the beginner course, [[rust-b-borrowing|References & Borrowing]] let you look
at data without owning it. Most of the time the compiler tracks how long each
reference stays valid on its own. A **lifetime annotation** is how you help it
when a function's output borrows from its input and Rust can't tell which one.

A lifetime is *not* a duration you compute. It is a relationship, written with a
label such as `'a` ("tick a"), that tells the compiler which input references an
output may borrow from. It changes nothing at runtime.

Using the same label does **not** say the inputs live equally long. At each call,
Rust chooses an `'a` that fits inside the usable lifetime of every reference
marked `'a`. With two inputs, that is usually their overlapping, shorter scope.

## Why a function sometimes needs one

This function returns one of its two arguments:

```rust
fn longest(a: &str, b: &str) -> &str {   // error: which input does the output borrow?
    if a.len() >= b.len() { a } else { b }
}
```

Rust rejects it. The returned `&str` borrows from `a` *or* `b`, and the compiler
must know the result can't outlive whichever one it came from. You spell that
out by giving both inputs and the output the **same** lifetime `'a`:

```rust
fn longest<'a>(a: &'a str, b: &'a str) -> &'a str {
    if a.len() >= b.len() { a } else { b }
}
```

Read it as: "for some lifetime `'a`, both inputs live at least `'a`, and the
result lives no longer than that." The result is valid only while *both* inputs
are — exactly the guarantee that makes returning either one safe.

More precisely, the signature describes the **overlap** the call can safely use:

```
`a` input:    ├──────────────────────────┤
`b` input:             ├────────────┤
safe `'a`:             ├────────────┤
returned ref:          ├───────┤       must stay inside the overlap
```

The annotation does not force either input to be dropped sooner. It limits how
long the returned reference may be used because the function body is allowed to
return either input.

```rust
let outer = String::from("a long string");
let chosen;
{
    let inner = String::from("short");
    chosen = longest(&outer, &inner);
    println!("{chosen}");             // valid: both strings are alive here
}
// println!("{chosen}");             // error: it might refer to `inner`
```

## An annotation connects references; it does not extend them

Lifetime syntax cannot make data live longer:

```rust
fn bad<'a>() -> &'a str {
    let local = String::from("temporary");
    &local
} // error: local is dropped here
```

The caller is allowed to choose `'a`, but the local string cannot satisfy that
promise. The usual fixes change ownership—return `String`—or borrow from an
input whose lifetime the signature can connect to the output.

When reading a signature, ask two questions:

1. Which inputs can the result point into?
2. What scope is common to all references carrying the same label?

Do not begin by asking “how many seconds is `'a`?” A lifetime is a compile-time
relationship between usable regions, not a runtime duration.

## Elision and `'static`

Rust omits obvious annotations using **lifetime elision** rules. For example,
`fn first(s: &str) -> &str` needs no label because there is only one borrowed
input. An annotation is needed when those rules cannot determine the relationship.

The special lifetime `'static` means a reference can remain valid for the rest
of the program. String literals have this type because they are stored in the
program binary:

```rust
let name: &'static str = "Rust";
```

That does not mean ordinary owned data lives forever, and adding `'static` is
not a general fix for lifetime errors.

## Your turn

1. **Trace:** In the nested-scope `longest` example, draw the usable region of
   `outer`, `inner`, and `chosen`.
2. **Compare signatures:** Explain why `fn first(s: &str) -> &str` can use
   elision but `longest` cannot.
3. **Diagnose:** Try returning a reference to a local `String`, then explain why
   adding `'a` cannot repair ownership.
4. **Implement:** Write `before_comma`. Its answer can borrow from either input,
   so its signature must connect both inputs to the returned reference.
