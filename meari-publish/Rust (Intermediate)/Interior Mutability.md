---
created: '2026-07-21'
id: rust-i-interior-mutability
source: meari-course
study:
  answer: |
    use std::cell::RefCell;

    struct Logger {
        entries: RefCell<Vec<String>>,
    }

    impl Logger {
        fn new() -> Self {
            Logger { entries: RefCell::new(Vec::new()) }
        }

        fn record(&self, message: &str) {
            self.entries.borrow_mut().push(message.to_string());
        }

        fn snapshot(&self) -> Vec<String> {
            self.entries.borrow().clone()
        }
    }
  kind: code
  lang: rust
  prompt: 'Define `Logger { entries: RefCell<Vec<String>> }`. Implement `new()`, `record(&self, message: &str)` to append an owned message through the `RefCell`, and `snapshot(&self) -> Vec<String>` to return a cloned snapshot. The public mutation method must take `&self`, not `&mut self`.'
  starter: |
    use std::cell::RefCell;

    struct Logger {
        entries: RefCell<Vec<String>>,
    }

    impl Logger {
        fn new() -> Self {
            todo!()
        }

        fn record(&self, message: &str) {
            // mutate entries without changing this signature to &mut self
        }

        fn snapshot(&self) -> Vec<String> {
            Vec::new()
        }
    }
  tests:
  - 'let log = Logger::new(); log.record("start"); log.record("done"); assert_eq!(log.snapshot(), vec!["start", "done"]);'
  - 'let log = Logger::new(); let first = log.snapshot(); log.record("later"); assert!(first.is_empty()); assert_eq!(log.snapshot(), vec!["later"]);'
subject: Rust (Intermediate)
title: Interior Mutability
---

The optional [[rust-b-smart-pointers|Box, Rc & RefCell]] lesson introduced
`RefCell`. Here we use it as an API-design tool. **Interior mutability** means a
type can change its internal state through a shared `&self` reference while
keeping that mutation controlled behind its methods.

## Why put a `RefCell` inside a type?

Some values are logically shared even though they record internal bookkeeping:
a test spy records calls, a cache remembers computed values, or a logger stores
events. Requiring `&mut self` would make the entire outer value exclusively
borrowed. A `RefCell<T>` narrows that flexibility to one field:

Put the `RefCell` inside the outer type, and its methods can mutate that field
through `&self` while the rest of the type remains ordinarily borrowed:

```rust
use std::cell::RefCell;

struct CallCounter {
    calls: RefCell<usize>,
}

impl CallCounter {
    fn record(&self) {
        *self.calls.borrow_mut() += 1; // mutation through &self
    }
}
```

## Borrow guards have a scope

`borrow()` and `borrow_mut()` return guard values. The dynamic borrow lasts until
its guard is dropped, which is usually the end of its scope:

```rust
let cell = RefCell::new(vec![1]);
{
    let mut values = cell.borrow_mut();
    values.push(2);
}                                      // mutable guard drops here
let length = cell.borrow().len();      // now a shared borrow is fine
```

Overlapping an incompatible borrow panics. For code that must handle contention,
`try_borrow()` and `try_borrow_mut()` return a `Result` instead. Keep guards
short-lived and never call unknown code while holding one.

The ordinary borrow rule still exists; only its enforcement time changes:

| Access | Ordinary reference | `RefCell` guard |
| ------ | ------------------ | --------------- |
| many readers | checked by compiler | checked at runtime |
| one writer | checked by compiler | checked at runtime |
| conflict | compile error | panic or `try_borrow*` error |

`RefCell` is for single-threaded interior mutability. It is not a lock and does
not make shared state safe across threads. [[rust-i-arc-mutex|Arc & Mutex]]
applies a related guard pattern with thread synchronization later.

`Rc<RefCell<T>>` adds multiple ownership when that is genuinely needed. Do not
add `Rc` merely to obtain interior mutability; a field inside one owning struct
can use `RefCell<T>` directly.

## Your turn

1. **Compare timing:** Explain how an overlapping ordinary borrow and
   `RefCell` borrow fail differently.
2. **Trace guards:** Mark where a mutable guard begins and drops.
3. **Choose a wrapper:** Explain why adding `Rc` is unnecessary when one
   `Logger` already owns the `RefCell`.
4. **Implement:** Build a logger whose `record` takes `&self` and whose owned
   snapshot does not expose a borrow guard.
