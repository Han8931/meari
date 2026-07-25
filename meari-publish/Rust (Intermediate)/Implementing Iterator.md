---
created: '2026-07-21'
id: rust-i-custom-iterator
source: meari-course
study:
  answer: |
    struct Counter {
        count: u32,
        max: u32,
    }

    impl Counter {
        fn new(max: u32) -> Self {
            Counter { count: 0, max }
        }
    }

    impl Iterator for Counter {
        type Item = u32;

        fn next(&mut self) -> Option<u32> {
            if self.count < self.max {
                self.count += 1;
                Some(self.count)
            } else {
                None
            }
        }
    }
  kind: code
  lang: rust
  prompt: 'Define `struct Counter { count: u32, max: u32 }` with `Counter::new(max: u32) -> Self` (starting count 0) and implement `Iterator` (`Item = u32`) so it yields `1, 2, ..., max`.'
  starter: |
    struct Counter {
        count: u32,
        max: u32,
    }

    impl Counter {
        fn new(max: u32) -> Self {
            Counter { count: 0, max }
        }
    }

    // impl Iterator for Counter  (Item = u32)
  tests:
  - assert_eq!(Counter::new(3).collect::<Vec<u32>>(), vec![1, 2, 3]);
  - assert_eq!(Counter::new(0).count(), 0);
  - assert_eq!(Counter::new(5).sum::<u32>(), 15);
  - 'let mut counter = Counter::new(1); assert_eq!(counter.next(), Some(1)); assert_eq!(counter.next(), None); assert_eq!(counter.next(), None);'
subject: Rust (Intermediate)
title: Implementing Iterator
---

You've *used* iterators; now you'll *build* one. Implement the `Iterator` trait
for your own type and it instantly gains the whole toolbox — `map`, `filter`,
`collect`, `sum`, `for` loops — all for free.

## One method: `next`

`Iterator` needs an associated `Item` type and a single required method, `next`,
which returns `Some(value)` for each element and `None` when finished:

```rust
struct Countdown {
    n: u32,
}

impl Iterator for Countdown {
    type Item = u32;

    fn next(&mut self) -> Option<u32> {
        if self.n == 0 {
            None
        } else {
            self.n -= 1;
            Some(self.n + 1)     // 3, 2, 1
        }
    }
}
```

`next` takes `&mut self` because an iterator is a state machine: every call
advances its position. Returning `None` says the stream is exhausted. Most
iterators keep returning `None` after that, and consumers are allowed to stop
calling as soon as they see it.

Trace the state rather than only the output:

| Before call | Operation | Return | State afterward |
| ----------- | --------- | ------ | --------------- |
| `n = 3` | decrement | `Some(3)` | `n = 2` |
| `n = 2` | decrement | `Some(2)` | `n = 1` |
| `n = 1` | decrement | `Some(1)` | `n = 0` |
| `n = 0` | no change | `None` | `n = 0` |

## Everything else comes free

Once `next` exists, the default methods light up:

```rust
let taken: Vec<u32> = Countdown { n: 3 }.collect();   // [3, 2, 1]
let sum: u32 = Countdown { n: 3 }.sum();               // 6
```

That's the payoff of implementing a standard trait: write the one method that's
truly yours, inherit dozens of others.

Iterator consumers differ in how they use ownership. `collect` and `sum` take
the iterator by value, so the iterator is consumed. Calling `next` manually only
borrows it mutably for that call:

```rust
let mut countdown = Countdown { n: 3 };
assert_eq!(countdown.next(), Some(3));
assert_eq!(countdown.next(), Some(2));
```

You also get `for`-loop support because the standard library provides
`IntoIterator` for every iterator. Do not add a second `IntoIterator`
implementation for `Counter`; the blanket implementation already covers it.

Returning `None` forever is conventional and the exercise tests it, but the
base `Iterator` trait does not require every iterator to be permanently
exhausted after its first `None`. The marker trait `FusedIterator` expresses
that stronger promise when generic code needs it.

## Your turn

1. **Trace:** Fill a state table for `Counter::new(3)`.
2. **Explain the receiver:** State why `next` needs `&mut self` but does not take
   ownership of the iterator.
3. **Predict ownership:** Compare one manual `next()` call with `.collect()`.
4. **Implement:** Build `Counter`, yielding `1, 2, ..., max` and remaining
   exhausted after `None`.
