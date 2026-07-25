---
created: '2026-07-21'
id: rust-i-iterator-pipelines
source: meari-course
study:
  answer: |
    fn normalized_words(lines: &[String]) -> Vec<String> {
        lines.iter()
            .flat_map(|line| line.split_whitespace())
            .filter(|word| word.len() >= 3)
            .map(str::to_lowercase)
            .collect()
    }
  kind: code
  lang: rust
  prompt: 'Write `normalized_words(lines: &[String]) -> Vec<String>`. Use one iterator pipeline to visit every whitespace-separated word with `flat_map`, discard words shorter than 3 bytes, lowercase the rest into owned `String`s, and collect them in input order.'
  starter: |
    fn normalized_words(lines: &[String]) -> Vec<String> {
        Vec::new()
    }
  tests:
  - 'let lines = vec!["Rust is FAST".to_string(), "and safe".to_string()]; assert_eq!(normalized_words(&lines), vec!["rust", "fast", "and", "safe"]);'
  - 'let lines = vec!["a bb ccc".to_string()]; assert_eq!(normalized_words(&lines), vec!["ccc"]);'
  - 'let lines: Vec<String> = vec![]; assert!(normalized_words(&lines).is_empty());'
subject: Rust (Intermediate)
title: Iterator Pipelines
---

The beginner course used `map`, `filter`, and `collect`. Intermediate pipelines
often have two extra jobs: flatten nested data and deliberately cross from
borrowed items to owned results. The item type at each stage tells you when that
transition happens.

## Reading a pipeline

```rust
let lines = vec![String::from("red green"), String::from("blue")];

let words: Vec<String> = lines
    .iter()                              // &String; `lines` stays usable
    .flat_map(|line| line.split_whitespace()) // &str from every line
    .map(str::to_uppercase)              // allocate owned Strings here
    .collect();
```

`flat_map` turns each line into an iterator of words and flattens those small
iterators into one stream. The `&str` values borrow from `lines`; `to_uppercase`
creates owned `String`s, so the collected result no longer borrows from it.

Write a **type ledger** beside a difficult chain:

| Stage | Item type | Owns text? |
| ----- | --------- | ---------- |
| `lines.iter()` | `&String` | no |
| `split_whitespace()` through `flat_map` | `&str` | no |
| `map(str::to_uppercase)` | `String` | yes |
| `collect::<Vec<String>>()` | collection of owned strings | yes |

If a closure signature surprises you, find the first row where your expected
item type differs from the real one.

## Common adapters

- `flat_map(f)` — turn each item into an iterator, then flatten the results.
- `copied()` / `cloned()` — turn an iterator of references into copied or cloned
  owned items.
- `inspect(f)` — observe borrowed items while debugging without transforming
  them.
- `try_fold(init, f)` — accumulate while allowing an early `Err` or `None`.

Adapters are lazy. A consumer such as `collect` drives the pipeline, normally
without constructing an intermediate collection between adapters.

## A note on the closures

Track the item type one stage at a time. `iter()` over `Vec<T>` yields `&T`, while
`into_iter()` consumes the vector and yields `T`. In addition, `filter` lends
each item to its predicate, so filtering an iterator of `&T` presents `&&T`.
Methods and automatic dereferencing sometimes hide that extra layer; explicit
patterns such as `|&&n|` make it visible.

`filter` is special because it must inspect an item without consuming the item
that may continue down the pipeline. Its predicate therefore receives `&Item`.
If `Item` is already `&T`, the predicate sees `&&T`.

## Your turn

1. **Ledger:** Write the item type after every stage in the example pipeline.
2. **Predict allocation:** Mark the exact adapter where borrowed text becomes
   owned.
3. **Diagnose:** Explain why `filter` may introduce one more `&`.
4. **Implement:** Flatten borrowed lines, filter words, and collect owned
   lowercase strings while leaving the input usable.
