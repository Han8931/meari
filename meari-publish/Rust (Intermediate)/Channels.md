---
created: '2026-07-21'
id: rust-i-channels
source: meari-course
study:
  answer: |
    use std::sync::mpsc;
    use std::thread;

    fn sum_from_workers(groups: Vec<Vec<i64>>) -> i64 {
        let (tx, rx) = mpsc::channel();

        let mut handles = Vec::new();
        for group in groups {
            let worker_tx = tx.clone();
            handles.push(thread::spawn(move || {
                for value in group {
                    worker_tx.send(value).unwrap();
                }
            }));
        }
        drop(tx);

        let sum = rx.into_iter().sum();
        for handle in handles {
            handle.join().unwrap();
        }
        sum
    }
  kind: code
  lang: rust
  prompt: 'Write `sum_from_workers(groups: Vec<Vec<i64>>) -> i64`. Spawn one worker per group, clone the sender for each worker, send every value, drop the original sender before receiving, sum until the channel closes, join every worker, and return the sum.'
  starter: |
    use std::sync::mpsc;
    use std::thread;

    fn sum_from_workers(groups: Vec<Vec<i64>>) -> i64 {
        0
    }
  tests:
  - assert_eq!(sum_from_workers(vec![vec![1, 2], vec![3, 4]]), 10);
  - assert_eq!(sum_from_workers(vec![]), 0);
  - assert_eq!(sum_from_workers(vec![vec![-5, 5], vec![100], vec![]]), 100);
subject: Rust (Intermediate)
title: Channels
---

Locks share state; **channels** share *messages*. A channel is a one-way pipe:
threads send values in one end, and another thread receives them out the other.
A common concurrency guideline is “share memory by communicating” rather than
making every worker mutate the same structure.

## Sender and receiver

`mpsc::channel()` returns a `(tx, rx)` pair — transmitter and receiver. Move `tx`
into a thread to send; loop over `rx` to receive:

```rust
use std::sync::mpsc;
use std::thread;

let (tx, rx) = mpsc::channel();

thread::spawn(move || {
    for i in 1..=3 {
        tx.send(i).unwrap();      // send 1, 2, 3
    }
});                               // tx dropped here -> channel closes

for received in rx {              // yields 1, 2, 3, then stops
    println!("{received}");
}
```

## How the loop knows to stop

Iterating `rx` blocks waiting for the next value and ends automatically when
every sender has been dropped. That's why moving `tx` into the thread matters:
when the thread finishes, `tx` drops, the channel closes, and the `for` loop
falls out cleanly. `mpsc` means *multi-producer, single-consumer* — clone `tx`
to have several senders feed one receiver. Remember that the original `tx` is
also a sender. If the receiver keeps it alive while iterating `rx`, the loop
waits forever because the channel is still open. Call `drop(tx)` after creating
the worker clones.

Trace sender ownership to understand termination:

```
original tx ──clone──► worker tx #1 ──┐
            └clone──► worker tx #2 ──┼── send values to rx
drop original tx                     │
workers finish and drop clones ──────┘
no senders remain → rx iteration ends
```

Dropping the original sender is not a signal value sent through the channel. It
removes one sender handle. The receiver reports disconnection only when **all**
sender handles are gone.

Messages from different producers may interleave nondeterministically. The
exercise sums values, so order does not affect the result. An application that
needs ordering must encode sequence information or coordinate producers.

## Your turn

1. **Trace handles:** Count live sender handles before cloning, after cloning,
   after dropping the original, and after workers finish.
2. **Diagnose a hang:** Explain why keeping the original `tx` alive prevents
   `rx` iteration from ending.
3. **Predict ordering:** State what is and is not guaranteed across producers.
4. **Implement:** Create several producers and sum on the main thread. The empty
   input test checks that the original sender is closed.
