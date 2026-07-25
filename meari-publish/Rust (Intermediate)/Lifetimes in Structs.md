---
created: '2026-07-21'
id: rust-i-lifetime-structs
source: meari-course
study:
  answer: |
    struct Parser<'a> {
        input: &'a str,
    }

    impl<'a> Parser<'a> {
        fn new(input: &'a str) -> Self {
            Parser { input }
        }

        fn first_word(&self) -> &'a str {
            self.input.split_whitespace().next().unwrap_or("")
        }
    }
  kind: code
  lang: rust
  prompt: 'Define `struct Parser<''a>` holding `input: &''a str`. Give it `Parser::new(input: &''a str) -> Self` and `first_word(&self) -> &''a str` returning the first whitespace-separated word (or `""` if there is none).'
  starter: |
    struct Parser<'a> {
        input: &'a str,
    }

    impl<'a> Parser<'a> {
        fn new(input: &'a str) -> Self {
            Parser { input }
        }

        fn first_word(&self) -> &'a str {
            ""
        }
    }
  tests:
  - assert_eq!(Parser::new("hello world").first_word(), "hello");
  - assert_eq!(Parser::new("  spaced  out").first_word(), "spaced");
  - assert_eq!(Parser::new("").first_word(), "");
  - 'let owned = String::from("borrowed input"); let parser = Parser::new(&owned); assert_eq!(parser.first_word(), "borrowed"); assert_eq!(owned.len(), 14);'
subject: Rust (Intermediate)
title: Lifetimes in Structs
---

A struct can hold a reference instead of owning its data — useful when you want a
lightweight view over something that lives elsewhere. But a struct holding a
reference must promise not to outlive what it points to, and you state that
promise with a lifetime, just like on a function.

## Declaring the lifetime

Put `<'a>` after the struct name and tag the borrowed field with it:

```rust
struct Excerpt<'a> {
    text: &'a str,          // this struct may not outlive `text`
}
```

Now `Excerpt` is only valid for as long as the string it borrows. The compiler
enforces that for you — a dangling `Excerpt` simply won't compile.

The lifetime parameter is part of the struct's type, but it stores no timestamp
or counter at runtime. It is compile-time evidence connecting the view to its
source. Because `Excerpt` borrows, constructing one does not move or copy the
underlying string data.

Trace construction as two values with different jobs:

```
String `source` owns UTF-8 bytes
        ▲
        │ borrowed as &'a str
Excerpt<'a> stores the reference, not the bytes
```

Dropping the `Excerpt` never drops the source text. Dropping the source while
the `Excerpt` might still be used is rejected.

## Methods on a borrowing struct

The `impl` block repeats the lifetime, and methods can hand back references that
share it:

```rust
impl<'a> Excerpt<'a> {
    fn new(text: &'a str) -> Self {
        Excerpt { text }
    }

    fn shout(&self) -> &'a str {
        self.text            // returns a reference living as long as the original
    }
}
```

There are two lifetimes involved in `shout`: the usually short borrow of
`&self`, and `'a`, the lifetime of the text stored inside the struct. Writing
`-> &'a str` says the returned slice comes from that stored text, not from the
temporary borrow of the `Excerpt` handle.

This can let a returned reference remain usable after the struct value itself
is dropped, provided the original text is still alive:

```rust
let text = String::from("an important sentence");
let first;
{
    let excerpt = Excerpt { text: &text };
    first = excerpt.shout();
}
println!("{first}");             // valid: `text` is still alive
```

## Borrowing struct or owning struct?

A borrowing struct avoids allocation and is ideal for a temporary view, parser,
or index into caller-owned data. The tradeoff is that the caller must keep the
source alive:

| Shape | Field | Benefit | Cost |
| ----- | ----- | ------- | ---- |
| borrowing | `text: &'a str` | no copy; can return slices | struct is tied to source lifetime |
| owning | `text: String` | self-contained; easier to store or return | allocation or ownership transfer |

Do not add a lifetime because it looks more advanced. If the value should keep
the text independently, own a `String`. Use a lifetime when “this is a view into
someone else's data” is genuinely part of the type's meaning.

Also notice the two distinct relationships in
`fn first_word(&self) -> &'a str`: `&self` is a short borrow needed to call the
method, while `'a` describes the stored input. The returned slice comes from
`self.input`, so it may remain usable after that short method call ends.

## Your turn

1. **Classify:** Decide whether a long-lived user profile and a temporary parser
   should own or borrow their text.
2. **Trace:** Identify the short `&self` borrow and stored `'a` relationship in
   `first_word`.
3. **Diagnose:** Try moving a `Parser` outside the scope of its owned input.
4. **Implement:** Build `Parser<'a>` and return its first word as a slice into
   the original text, with no copying.
