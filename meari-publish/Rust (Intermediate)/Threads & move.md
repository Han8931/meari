---
created: '2026-07-21'
id: rust-i-threads
source: meari-course
study:
  answer: |
    use std::thread;

    fn parallel_squares(nums: Vec<i32>) -> Vec<i32> {
        let mut handles = Vec::new();
        for n in nums {
            handles.push(thread::spawn(move || n * n));
        }

        let mut results = Vec::new();
        for handle in handles {
            results.push(handle.join().unwrap());
        }
        results
    }
  kind: code
  lang: rust
  prompt: 'Write `parallel_squares(nums: Vec<i32>) -> Vec<i32>` that spawns one thread per number to square it (with a `move` closure), joins them in order, and returns the squares in the original order.'
  starter: |
    use std::thread;

    fn parallel_squares(nums: Vec<i32>) -> Vec<i32> {
        Vec::new()
    }
  tests:
  - assert_eq!(parallel_squares(vec![1, 2, 3, 4]), vec![1, 4, 9, 16]);
  - assert_eq!(parallel_squares(vec![]), vec![]);
  - assert_eq!(parallel_squares(vec![-3]), vec![9]);
subject: Rust (Intermediate)
title: Threads & move
---

Rust runs code in parallel with **threads**, and its ownership rules make that
far safer than usual: the compiler simply won't let a thread reference data that
might disappear out from under it.

## Spawning and joining

`thread::spawn` starts a thread with a closure; the returned handle's `join`
waits for it to finish and hands back whatever the closure returned:

```rust
use std::thread;

let handle = thread::spawn(|| {
    1 + 1                         // this runs on another thread
});

let result = handle.join().unwrap();   // wait, then collect -> 2
```

## Why `move`

A spawned thread may outlive the function that started it, so it must **own**
whatever it uses. The `move` keyword forces the closure to take ownership of its
captures instead of borrowing them:

```rust
let name = String::from("Rust");
let handle = thread::spawn(move || {
    format!("hello from {name}")    // name is owned by the thread now
});
```

Without `move`, the closure would borrow `name`, and the compiler would refuse —
it can't prove `name` lives long enough. `move` resolves that by transferring
ownership.

More precisely, `thread::spawn` requires a closure that is `Send + 'static`:
`Send` means its captured values may cross to another thread, and `'static`
means it contains no short-lived borrowed references. Owned values satisfy the
lifetime requirement when their types are safe to send. This is why `move` is
the usual solution, though references to truly static data can also work.

Do not read `'static` here as “the thread runs forever.” It constrains borrowed
data inside the closure, because the thread *could* outlive the current stack
frame. An owned `String` qualifies even though it is dropped normally when the
thread ends.

## `Send`, joining, and failure

`Send` is a marker trait meaning ownership of a value may move to another
thread. Most ordinary owned types are `Send`; `Rc<T>` is a notable exception,
which is why shared threaded ownership uses `Arc<T>` in the next lesson.

`join()` returns a `Result` because the child may panic:

```
child returns value  → join() = Ok(value)
child panics         → join() = Err(panic payload)
```

The examples use `.unwrap()` to keep the exercise focused, which deliberately
propagates a child panic as a panic in the joining thread. Production code
should decide whether to propagate, report, or recover.

When threads are guaranteed to finish before a scope ends, `thread::scope` can
permit them to borrow local data. Ordinary `thread::spawn` uses the simpler
owned/`'static` model taught here.

One thread per integer, as used by the study challenge, is a learning exercise
for captures and join handles—not a performance recommendation. Real parallel
work normally uses a bounded worker pool or a parallel-iterator library so
thread creation does not cost more than the computation.

## Your turn

1. **Trace ownership:** Mark when each number moves into a closure and when its
   result moves back through `join`.
2. **Interpret bounds:** Explain why owned data can satisfy `'static` without
   living forever.
3. **Handle failure:** State what `join().unwrap()` does when a child panics.
4. **Implement:** Square each number on its own thread and join handles in input
   order, while recognizing this as an ownership exercise rather than an
   efficient squaring algorithm.
