---
created: "2026-07-08"
id: rust-b-collections
source: meari-course
study:
  answer: |
    fn word_count(text: &str) -> HashMap<String, u32> {
        let mut counts = HashMap::new();
        for word in text.split_whitespace() {
            *counts.entry(word.to_string()).or_insert(0) += 1;
        }
        counts
    }
  kind: code
  lang: rust
  prompt: 'Write `word_count(text: &str) -> HashMap<String, u32>` counting how many times each whitespace-separated word appears. (`HashMap` is already in scope.)'
  starter: |
    fn word_count(text: &str) -> HashMap<String, u32> {
        HashMap::new()
    }
  tests:
    - assert_eq!(word_count("a b a").get("a"), Some(&2));
    - assert_eq!(word_count("a b a").get("b"), Some(&1));
    - assert_eq!(word_count("a b a").get("c"), None);
subject: Rust (Beginner)
title: Vec & HashMap
---

Arrays and tuples from [[Arrays, Tuples & Slices]] are fixed-size. Real programs
need collections that **grow and shrink at runtime**. The two workhorses live on
the heap: `Vec<T>`, a growable list, and `HashMap<K, V>`, a key-value store.

## `Vec<T>`: a growable list

```rust
let mut nums: Vec<i32> = Vec::new();
nums.push(10);                 // add to the end
nums.push(20);
nums.push(30);

let more = vec![1, 2, 3];      // vec! macro — build from literals

println!("{}", nums.len());    // 3
let last = nums.pop();         // Some(30) — removes & returns the end
```

### Two ways to read an element

```rust
let v = vec![10, 20, 30];

let a = v[1];                  // 20 — direct index, PANICS if out of bounds
let b = v.get(1);              // Some(20) — returns Option, safe
let c = v.get(99);             // None — no panic
```

Indexing (`v[i]`) is concise but panics on a bad index; `v.get(i)` returns an
[[Option & Result|Option]] so you handle the miss safely. Choose based on
whether an out-of-range index is a bug (`[i]`) or an expected possibility
(`.get`).

### Iterating

```rust
let mut v = vec![1, 2, 3];             // `mut` so we can borrow it mutably below

for x in &v      { print!("{x} "); }   // borrow — v survives
for x in &mut v  { *x *= 2; }          // mutable borrow — doubles in place
for x in v       { print!("{x} "); }   // consumes v
```

Same borrow-vs-consume choice as in [[Control Flow]].

### How a Vec grows

A `Vec` keeps spare **capacity**; when it fills, it allocates a bigger buffer and
moves the elements. Length ≤ capacity always:

```
   len = 3, cap = 4     [ 1 | 2 | 3 | _ ]     one slot free
   push(4)              [ 1 | 2 | 3 | 4 ]     now full
   push(5) → reallocate [ 1 | 2 | 3 | 4 | 5 | _ | _ | _ ]  cap grows (→8)
```

## `HashMap<K, V>`: key–value pairs

```rust
use std::collections::HashMap;

let mut scores: HashMap<String, i32> = HashMap::new();
scores.insert(String::from("Ana"), 10);
scores.insert(String::from("Bo"), 7);

// lookup returns an Option — the key might not be present
match scores.get("Ana") {
    Some(s) => println!("Ana: {s}"),
    None    => println!("no score"),
}
```

### The `entry` API

A frequent need is "update if present, insert a default if not." The `entry` API
does it in one clean expression — perfect for counting:

```rust
let text = "a b a c b a";
let mut counts: HashMap<&str, i32> = HashMap::new();

for word in text.split_whitespace() {
    *counts.entry(word).or_insert(0) += 1;   // default 0, then increment
}
// counts == {"a": 3, "b": 2, "c": 1}
```

`entry(word).or_insert(0)` returns a mutable reference to the value — either the
existing one or a freshly inserted `0` — which you then `+= 1`.

## The same in Python

These two map almost directly onto Python's `list` and `dict` — probably the
Python types you use most:

```python
nums = [10, 20]
nums.append(30)               # ~ Vec::push
last = nums.pop()             # ~ Vec::pop

scores = {}
scores["Ana"] = 10            # ~ HashMap::insert
scores.get("Bo")             # None if absent — like HashMap::get → Option
```

The counting pattern `*counts.entry(word).or_insert(0) += 1` is Python's
`counts[word] = counts.get(word, 0) + 1`, or more idiomatically a
`collections.Counter`. The main difference: Rust makes you handle the
"key missing" case (via the `Option` from `.get`), whereas Python's `d[k]` just
raises a `KeyError` at runtime.

## Choosing a collection

| Type          | Use for                              | Ordered?          |
| ------------- | ------------------------------------ | ----------------- |
| `Vec<T>`      | an ordered list, a stack, a sequence | yes (insertion)   |
| `HashMap<K,V>`| fast key → value lookup              | no                |
| `BTreeMap<K,V>`| key → value, kept **sorted** by key | yes (by key)      |
| `HashSet<T>`  | membership / uniqueness              | no                |

Reach for `Vec` by default; use `HashMap` when you look things up by a key; pick
the `BTree*` variants when you need sorted order.

## Why collection access often borrows

`v.get(0)` returns `Option<&T>`, not `Option<T>`. Returning an owned `T` would
have to move the element out or clone it; a reference lets you inspect it cheaply
while the vector remains its owner. Similarly, `map.get(key)` returns a borrowed
value. The reference cannot outlive or conflict with mutation of its collection.

```rust
let mut names = vec![String::from("Ana")];
let first = &names[0];
println!("{first}"); // last use of the borrow
names.push(String::from("Bo")); // now mutation is allowed
```

This matters because `push` may reallocate the vector and move its buffer. Rust
will not allow a reference into the old buffer to remain usable afterward.

## Try it

1. Create a `Vec`, push three numbers, then pop one off.
2. Use `.get(10)` on a short vector and handle the `None` case.
3. Count words in a short string using `HashMap` and the `entry` API.

> **Takeaway:** `Vec<T>` is your growable list (index with `[i]` for bugs, `.get`
> for maybes), and `HashMap<K,V>` is your key-value store (get returns an
> `Option`; the `entry` API makes update-or-insert a one-liner).
