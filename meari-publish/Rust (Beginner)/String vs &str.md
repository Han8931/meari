---
created: "2026-07-08"
id: rust-b-string
source: meari-course
study:
  answer: |
    fn full_name(first: &str, last: &str) -> String {
        format!("{first} {last}")
    }
  kind: code
  lang: rust
  prompt: 'Complete `full_name(first: &str, last: &str) -> String` with `format!("{first} {last}")`. It borrows two string slices and returns a new owned `String`.'
  starter: |
    fn full_name(first: &str, last: &str) -> String {
        String::new()
    }
  tests:
    - assert_eq!(full_name("Ana", "Smith"), "Ana Smith");
    - assert_eq!(full_name("Rust", "Learner"), "Rust Learner");
subject: Rust (Beginner)
title: String vs &str
---

Rust has two common text types:

- `String` is text that **owns its data**.
- `&str` is a **borrowed view** of text owned somewhere else.

That is the main idea. Most of the time, your choice follows one simple rule:

> Use `String` when a value must own or change its text. Use `&str` when code
> only needs to read text for a while.

This is [[Ownership & Moves|ownership]] and
[[References & Borrowing|borrowing]] applied to text.

## Why does Rust need both?

A **string literal** and the `String` type solve different problems.

```rust
let fixed = "hello";                    // string literal: &'static str
let mut changing = String::from("hello"); // owned String
changing.push('!');
```

The literal `"hello"` is text written directly in the source code. Its contents
are known when the program is compiled, so Rust can place those bytes in the
compiled program. They stay available for the program's entire run and cannot
be resized. The variable `fixed` only refers to that existing text, which is why
its type is `&str`—more precisely, `&'static str`.

A `String` is needed for text whose contents or size are decided while the
program runs: user input, a file's contents, a formatted message, or text built
in a loop. It owns a growable buffer, so that text can outlive the operation
that created it and can be changed when the `String` is mutable.

```rust
let name = String::from("Ana");
let message: String = format!("Hello, {name}!"); // created at runtime
```

Rust keeps these ideas separate so ownership and allocation are visible:

| Text | Where it comes from | Can grow? | Common type |
| --- | --- | --- | --- |
| `"hello"` | fixed text in the compiled program | no | `&'static str` |
| `String::from("hello")` | an owned runtime buffer | yes, if mutable | `String` |
| `&owned[0..5]` | a borrowed view into existing text | no | `&str` |

A string literal is an expression that denotes particular fixed text; `String`
is a type that can hold many owned text values. Also, `&str` does **not** mean only
“string literal.” It can borrow literal text, all of a `String`, or part of one.
Without `String`, Rust would have no general growable, owned text value. Without
`&str`, code would often have to copy or take ownership of text merely to read
it.

## Start with one example

```rust
let name: String = String::from("Ana");
let view: &str = &name;

println!("{name}"); // the String owns "Ana"
println!("{view}"); // view borrows the same text
```

No text is copied when `view` is created. It temporarily refers to the text
inside `name`.

Read `&str` as **“a borrowed string slice.”** The `&` is the same borrow symbol
used for other references.

```text
name: String  ── owns ──>  "Ana"
view: &str    ── borrows ──┘
```

Because `view` borrows from `name`, Rust will not let the view remain valid
after `name` is gone.

## How to create each type

A quoted string literal has type `&'static str`, which can be used wherever an
`&str` is expected. `'static` means the literal's text remains valid for the
entire program:

```rust
let a: &'static str = "hello";
let b: &str = "world"; // the shorter spelling is usually enough
```

Use `String::from` or `.to_string()` when you need owned text:

```rust
let b: String = String::from("hello");
let c: String = "hello".to_string();
```

Borrow a `String` to get an `&str`:

```rust
let owned = String::from("hello");
let borrowed: &str = &owned;
```

So the common conversions are:

```text
"hello"                 has type &str
String::from("hello")   creates a String
&owned                  borrows a String as &str
```

## Function parameters: usually take `&str`

If a function only reads text, give it an `&str` parameter:

```rust
fn greet(name: &str) {
    println!("Hello, {name}!");
}

let owned = String::from("Ana");

greet("Bo");     // works: a literal is an &str
greet(&owned);   // works: borrow the String
println!("{owned}"); // still works: greet did not take ownership
```

This is flexible: callers can pass either a string literal or a borrowed
`String`.

Compare that with taking a `String`:

```rust
fn greet_owned(name: String) {
    println!("Hello, {name}!");
}

let owned = String::from("Ana");
greet_owned(owned);     // moves the String into the function
// println!("{owned}"); // error: owned was moved
```

Taking a `String` is not wrong. Do it when the function needs to keep or take
ownership of the text. Do not require ownership when borrowing is enough.

## Return `String` for newly created text

A function cannot return a borrowed view of a temporary local string. When a
function builds new text, it usually returns a `String` that owns the result:

```rust
fn full_name(first: &str, last: &str) -> String {
    format!("{first} {last}")
}

let full = full_name("Ana", "Smith");
println!("{full}");
```

Notice the useful pattern:

```text
borrow inputs as &str  →  build new text  →  return String
```

`format!` creates a new `String` and only borrows the values inserted into it.

## Changing text requires a mutable `String`

`String` owns a growable buffer, so a mutable `String` can be changed:

```rust
let mut message = String::from("Hello");
message.push_str(", world"); // append text
message.push('!');           // append one character

assert_eq!(message, "Hello, world!");
```

An `&str` is a read-only view, so it cannot be grown with `push_str`.

One way to combine strings is `+`, but it moves its left side:

```rust
let a = String::from("hello");
let b = String::from(" world");
let c = a + &b;

// a is no longer usable
println!("{b}"); // b was only borrowed, so it is still usable
```

For beginner code, `format!` is often clearer because it does not consume its
arguments:

```rust
let a = String::from("hello");
let b = String::from("world");
let c = format!("{a} {b}");

println!("{a}, {b}, {c}"); // all three are usable
```

## What does “slice” mean?

An `&str` may view all of a string or only part of it:

```rust
let message = String::from("hello world");
let hello: &str = &message[0..5];

println!("{hello}");
```

The slice stores where the borrowed text starts and how many bytes it contains.
It does not own or copy those bytes.

You do not need to slice strings often as a beginner. Rust strings use UTF-8,
so slice positions are byte positions and must be on valid character
boundaries. For example, arbitrary indexes can be unsafe for text containing
characters such as `é` or `日`.

## UTF-8 means no `text[0]`

Rust does not allow integer indexing into a string:

```rust
let text = "héllo";
// let first = text[0]; // error
```

A character can use more than one byte, so `0` could mean “first byte” or
“first character.” Iterate explicitly instead:

```rust
for ch in "héllo".chars() {
    println!("{ch}");
}

println!("{}", "héllo".len());           // 6 bytes
println!("{}", "héllo".chars().count()); // 5 characters
```

This UTF-8 rule applies to both `String` and `&str`.

## Quick decision guide

Ask what the code needs to do:

| Need | Choose |
| --- | --- |
| Store text in a struct or collection | `String` |
| Build, append to, or modify text | `String` (usually `mut`) |
| Return newly created text | `String` |
| Read text in a function | `&str` |
| Refer to a string literal | `&str` |
| Refer to part of existing text | `&str` |

A useful function-signature pattern is:

```rust
fn transform(input: &str) -> String {
    // borrow existing text, return newly owned text
    input.to_uppercase()
}
```

## Try it

1. Identify the type of each value:

   ```rust
   let a = "hello";
   let b = String::from("hello");
   let c: &str = &b;
   ```

2. Write `greet(name: &str)` and call it with both `"Ana"` and a borrowed
   `String`.
3. Write `full_name(first: &str, last: &str) -> String` using `format!`.
4. Print both `.len()` and `.chars().count()` for `"héllo"` and explain the
   difference.

> **Takeaway:** `String` owns text; `&str` borrows text. For read-only function
> inputs, prefer `&str`. When building or storing new text, use `String`.
