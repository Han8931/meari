---
created: '2026-07-21'
id: rust-i-arc-mutex
source: meari-course
study:
  answer: |
    use std::sync::{Arc, Mutex};
    use std::thread;

    fn concurrent_increment(threads: usize, per_thread: usize) -> usize {
        let counter = Arc::new(Mutex::new(0usize));
        let mut handles = Vec::new();

        for _ in 0..threads {
            let counter = Arc::clone(&counter);
            handles.push(thread::spawn(move || {
                let mut n = counter.lock().unwrap();
                *n += per_thread;
            }));
        }

        for handle in handles {
            handle.join().unwrap();
        }

        let total = *counter.lock().unwrap();
        total
    }
  kind: code
  lang: rust
  prompt: 'Write `concurrent_increment(threads: usize, per_thread: usize) -> usize` that spawns `threads` threads. Give each worker an `Arc` clone, lock once, add `per_thread` to the shared `Mutex<usize>`, join every worker, and return the final total.'
  starter: |
    use std::sync::{Arc, Mutex};
    use std::thread;

    fn concurrent_increment(threads: usize, per_thread: usize) -> usize {
        0
    }
  tests:
  - assert_eq!(concurrent_increment(4, 1000), 4000);
  - assert_eq!(concurrent_increment(1, 50), 50);
  - assert_eq!(concurrent_increment(0, 100), 0);
subject: Rust (Intermediate)
title: Shared State with Arc & Mutex
---

Last lesson each thread owned its own data. To let several threads touch *one*
shared value, you need two things: shared ownership that's safe across threads,
and locked access so only one thread mutates at a time. That's `Arc<Mutex<T>>`.

## Arc: shared ownership across threads

`Arc` is an atomically reference-counted counterpart to `Rc`. Clone it to give
each thread a handle to the same value. `Arc<T>` does not make an unsafe inner
type safe: it can cross threads only when `T` satisfies the required `Send` and
`Sync` bounds. `Mutex<T>` supplies those synchronization guarantees for many
shared-mutation designs.

## Mutex: one writer at a time

`Mutex<T>` guards the value. `lock()` waits until the lock is free, then returns
a guard you deref to read or write; the lock releases when the guard drops:

```rust
use std::sync::{Arc, Mutex};
use std::thread;

let counter = Arc::new(Mutex::new(0));
let mut handles = vec![];

for _ in 0..3 {
    let counter = Arc::clone(&counter);
    handles.push(thread::spawn(move || {
        *counter.lock().unwrap() += 1;    // lock, mutate, unlock on drop
    }));
}
for h in handles { h.join().unwrap(); }

assert_eq!(*counter.lock().unwrap(), 3);
```

The `Mutex` makes the increments take turns, so no update is ever lost — the data
race is impossible, not just unlikely.

Trace the wrapper layers:

```
Arc clone  → owns one shared handle
Mutex      → controls access to the inner usize
lock()     → returns a guard
*guard     → accesses the usize
drop guard → unlocks the Mutex
```

Keep the locked region as short as correctness permits. The study challenge has
each worker add its whole contribution under one lock rather than lock once per
integer; synchronization itself has a cost.

`lock()` returns a `Result` because another thread might panic while holding the
lock, leaving it **poisoned**. `.unwrap()` turns poisoning into another panic.
That is acceptable for this focused exercise, but applications may inspect the
error and decide whether the protected data can be recovered.

Two more distinctions:

- `Arc` solves ownership, not mutation.
- `Mutex` synchronizes access, but careless lock ordering can still deadlock.

Rust prevents data races; it cannot prove that every locking strategy makes
progress.

## Your turn

1. **Unwrap the layers:** For `Arc<Mutex<usize>>`, name the job of each wrapper.
2. **Trace counts and guards:** Mark every `Arc::clone`, lock acquisition, guard
   drop, and final owner.
3. **Compare critical sections:** Explain why adding a batch under one lock is
   preferable to locking for every unit.
4. **Implement:** Have several threads add fixed contributions to one shared
   counter, then join and return the result.
