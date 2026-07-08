---
created: "2026-07-08"
id: rust-b-structs
source: meari-course
subject: Rust (Beginner)
title: Structs & Methods
---

A **struct** lets you name a bundle of related data as a single custom type —
the equivalent of a class's fields, minus the inheritance. Combined with `impl`
blocks for behavior, structs are how you model the "nouns" of your program.

## Defining and creating a struct

```rust
struct User {
    name: String,
    age: u32,
    active: bool,
}

let u = User {
    name: String::from("Ana"),
    age: 30,
    active: true,
};

println!("{} is {}", u.name, u.age);   // dot access
```

```
   User {
     name:   "Ana"     ← fields, each with its own type
     age:    30
     active: true
   }
```

To mutate a field, the whole binding must be `mut` (you can't mark individual
fields mutable):

```rust
let mut u = u;
u.age += 1;             // needs `mut u`
```

### Two ergonomic shortcuts

**Field init shorthand** — when a variable has the same name as the field:

```rust
fn make_user(name: String, age: u32) -> User {
    User { name, age, active: true }   // name:name, age:age implied
}
```

**Struct update syntax** — build one struct from another, changing only some
fields:

```rust
let u2 = User { age: 31, ..u };        // take the rest of the fields from `u`
```

Because `User` owns a `String`, `..u` **moves** those non-`Copy` fields out of
`u` — so after this line `u` can no longer be used as a whole value (its `name`
was moved into `u2`). This is the same move rule from [[Ownership & Moves]]
showing up in struct syntax.

## Other struct shapes

| Kind         | Definition              | Use for                                   |
| ------------ | ----------------------- | ----------------------------------------- |
| Named-field  | `struct P { x: i32 }`   | the usual case — self-documenting fields  |
| Tuple struct | `struct Point(i32,i32)` | lightweight wrapper, fields by position   |
| Unit struct  | `struct Marker;`        | a type with no data (markers, traits)     |

```rust
struct Point(i32, i32);       // tuple struct
let p = Point(3, 4);
println!("{}", p.0);          // access by position → 3

struct AlwaysEqual;           // unit struct — carries no data
```

## Methods with `impl`

Behavior lives in an `impl` block. The first parameter, `self`, is the instance
the method is called on — and *how* you take `self` matters a lot:

```rust
struct Rectangle { width: u32, height: u32 }

impl Rectangle {
    // associated function (no self) — the idiomatic constructor
    fn new(width: u32, height: u32) -> Rectangle {
        Rectangle { width, height }
    }

    fn area(&self) -> u32 {          // &self: read-only borrow
        self.width * self.height
    }

    fn scale(&mut self, factor: u32) { // &mut self: mutable borrow
        self.width *= factor;
        self.height *= factor;
    }

    fn consume(self) -> u32 {        // self: takes ownership, ends the instance
        self.width + self.height
    }
}

let mut r = Rectangle::new(3, 4);    // :: calls an associated function
println!("{}", r.area());            // 12  — . calls a method
r.scale(2);                          // r is now 6×8
```

The three receiver forms, and when to reach for each:

| Receiver    | Meaning                        | Reach for it when…                     |
| ----------- | ------------------------------ | -------------------------------------- |
| `&self`     | borrow immutably               | you only need to read (most methods)   |
| `&mut self` | borrow mutably                 | you modify the instance in place       |
| `self`      | take ownership (consumes it)   | you transform/finalize and it's done   |

Note the two call syntaxes: `Rectangle::new(..)` uses `::` because `new` is an
**associated function** (no `self`, like a static/class method). `r.area()` uses
`.` because `area` is a **method** on an instance. `new` is just a convention —
Rust has no built-in constructors.

Those `&self` / `&mut self` distinctions *are* the borrowing rules from
[[References & Borrowing]] applied to methods. Next, we model choices between
alternatives with [[Enums & Pattern Matching]].

## The same in Python

A Rust struct plus its `impl` block is a Python **class**: fields become
instance attributes, methods stay methods, and the `new` associated function
becomes `__init__`:

```python
class Rectangle:
    def __init__(self, width, height):   # ~ Rectangle::new
        self.width = width
        self.height = height

    def area(self):                       # ~ fn area(&self)
        return self.width * self.height

    def scale(self, factor):              # ~ fn scale(&mut self, ...)
        self.width *= factor
        self.height *= factor

r = Rectangle(3, 4)
print(r.area())                           # 12
```

The distinction Rust adds is on `self`: `&self` / `&mut self` / `self` declare
whether a method *reads*, *mutates*, or *consumes* the instance. Python's `self`
makes no such promise — any method can mutate the object.

## Try it

1. Define a `Book` struct with `title: String` and `pages: u32`.
2. Add a method `is_long(&self) -> bool` that returns true if `pages > 300`.
3. Create a `Book`, call the method, and print the result.

> **Takeaway:** structs name your data; `impl` blocks attach behavior. Default to
> `&self`, upgrade to `&mut self` when you mutate, and reserve `self` for methods
> that consume the value.
