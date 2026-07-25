---
created: '2026-07-21'
id: rust-i-associated-types
source: meari-course
study:
  answer: |
    trait Container {
        type Item;
        fn get(&self, i: usize) -> Option<&Self::Item>;

        fn first(&self) -> Option<&Self::Item> {
            self.get(0)
        }
    }

    struct Stack {
        items: Vec<i32>,
    }

    impl Container for Stack {
        type Item = i32;

        fn get(&self, i: usize) -> Option<&i32> {
            self.items.get(i)
        }
    }
  kind: code
  lang: rust
  prompt: 'Define a trait `Container` with an associated type `Item`, a required `fn get(&self, i: usize) -> Option<&Self::Item>`, and a default `fn first(&self) -> Option<&Self::Item>` that returns element 0. Implement it for `struct Stack { items: Vec<i32> }`.'
  starter: |
    trait Container {
        type Item;
        fn get(&self, i: usize) -> Option<&Self::Item>;

        fn first(&self) -> Option<&Self::Item> {
            self.get(0)
        }
    }

    struct Stack {
        items: Vec<i32>,
    }

    // impl Container for Stack
  tests:
  - 'let s = Stack { items: vec![10, 20, 30] }; assert_eq!(s.get(1), Some(&20));'
  - 'let s = Stack { items: vec![10, 20, 30] }; assert_eq!(s.first(), Some(&10));'
  - 'let e = Stack { items: vec![] }; assert_eq!(e.first(), None);'
  - 'let s = Stack { items: vec![10] }; let item: Option<&i32> = s.get(0); assert_eq!(item, Some(&10));'
subject: Rust (Intermediate)
title: Associated Types
---

A trait can leave a *type* blank for each implementor to fill in, called an
**associated type**. You've already used one: `Iterator` has `type Item`, and
each iterator says what it yields. Now you'll define your own.

## A placeholder type inside a trait

Write `type Item;` in the trait, then refer to it as `Self::Item`:

```rust
trait Container {
    type Item;                                   // filled in per implementor
    fn get(&self, i: usize) -> Option<&Self::Item>;
}
```

Each implementor picks the concrete type once:

```rust
struct Bag { items: Vec<String> }

impl Container for Bag {
    type Item = String;                          // Bag's Item is String
    fn get(&self, i: usize) -> Option<&String> {
        self.items.get(i)
    }
}
```

`Self` means the type currently implementing the trait (`Bag` here), so
`Self::Item` means `Bag`'s chosen `Item` type. Callers do not have to repeat
`String` in the trait name or method call; the implementation fixes that
relationship.

When reading `C::Item`, substitute from the implementation:

```
C = Bag
Bag implements Container with Item = String
therefore C::Item = String
```

`Self` inside a trait or implementation means the implementing type. `C::Item`
in generic code means the associated `Item` selected by `C`'s implementation.

## Why not a generic `<T>`?

A generic trait (`Container<T>`) could be implemented many times for one type,
with different `T`s. An associated type says "there is exactly **one** Item type
per implementor" — cleaner when only one answer makes sense, which is why
`Iterator` uses it.

That distinction affects inference. With `Container<T>`, code may need help
choosing which implementation it means. With `Container`, Rust can start from
the implementor and determine its item type:

```rust
fn inspect<C: Container>(container: &C, index: usize) -> Option<&C::Item> {
    container.get(index)
}
```

## Default methods can use it

Because the trait knows `Self::Item` exists, a default method can build on the
required one:

```rust
fn first(&self) -> Option<&Self::Item> {
    self.get(0)
}
```

Implementors inherit `first` unless they override it. Default methods are a
useful way to keep a trait's core contract small: require the primitive
operation, then define convenient operations in terms of it.

## Your turn

1. **Substitute:** If `Bag::Item = String`, read
   `Option<&Bag::Item>` as a concrete type.
2. **Compare designs:** Explain when `Container<T>` permits a choice that
   `Container { type Item; }` deliberately removes.
3. **Fill:** Define the default `first` method in terms of required `get`.
4. **Implement:** Give `Stack` the one associated item type `i32`.
