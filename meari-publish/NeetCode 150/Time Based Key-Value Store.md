---
created: "2026-06-07"
id: nc-time-based-key-value-store
source: imported:neetcode-150
study:
  answer: |-
    Store key → list of (timestamp, value), appended in order. get binary searches that list for the last entry with ts ≤ query (bisect_right minus one) and returns its value, or "" if none.

    Complexity: O(1) set, O(log n) get; O(n) space
  kind: essay
  prompt: 'Solve "Time Based Key-Value Store" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Time Based Key-Value Store
---

**Pattern:** Binary Search · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/time-based-key-value-store)

Design a store with `set(key, value, timestamp)` and `get(key, timestamp)` returning the value set at the **largest timestamp ≤** the query (or `""`). Timestamps in `set` calls are strictly increasing per key.

**Example 1:**

    Input: set("foo","bar",1), get("foo",1), get("foo",3)
    Output: "bar", "bar"
    Explanation: at t=3 the latest value at or before 3 is still "bar".

**Constraints:**

- `1 <= key.length, value.length <= 100`
- `1 <= timestamp <= 10^7`
- Up to `2 * 10^5` calls

---

**Hints — try each one before reading on:**
1. Timestamps arrive strictly increasing — each key's list is already sorted.
2. get is a binary search for the rightmost timestamp ≤ query.

**Target:** O(1) set, O(log n) get; O(n) space
