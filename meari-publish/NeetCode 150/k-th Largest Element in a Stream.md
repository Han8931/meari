---
created: "2026-06-07"
id: nc-k-th-largest-element-in-a-stream
source: imported:neetcode-150
study:
  answer: |-
    Maintain a min-heap capped at k elements: push each new value and pop the smallest while the heap exceeds k; add() returns the heap root.

    Complexity: O(log k) per add, O(k) space
  kind: essay
  prompt: 'Solve "k-th Largest Element in a Stream" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: k-th Largest Element in a Stream
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/kth-largest-element-in-a-stream)

Design `KthLargest(k, nums)` with `add(val)` that appends a value and returns the `k`-th largest element seen so far (duplicates count).

**Example 1:**

    Input: KthLargest(3, [4,5,8,2]); add(3); add(5); add(10); add(9); add(4)
    Output: 4, 5, 5, 8, 8

**Constraints:**

- `1 <= k <= 10^4`
- `0 <= nums.length <= 10^4`
- At least k elements before queries

---

**Hints — try each one before reading on:**
1. Only the k largest matter — keep exactly them.
2. A min-heap of size k: its root IS the k-th largest.

**Target:** O(log k) per add, O(k) space
