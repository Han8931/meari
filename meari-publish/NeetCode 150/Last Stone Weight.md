---
created: "2026-06-07"
id: nc-last-stone-weight
source: imported:neetcode-150
study:
  answer: |-
    Max-heap (negated values in heapq): pop two, push back their difference if nonzero, repeat; the survivor (or 0) is the answer.

    Complexity: O(n log n) time, O(n) space
  kind: essay
  prompt: 'Solve "Last Stone Weight" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Last Stone Weight
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/last-stone-weight)

Repeatedly smash the two heaviest stones: equal weights destroy both; otherwise the heavier survives with the difference. Return the last stone's weight, or 0.

**Example 1:**

    Input: stones = [2,7,4,1,8,1]
    Output: 1
    Explanation: 8&7→1, 4&2→2, 2&1→1, 1&1→0; remaining [1].

**Constraints:**

- `1 <= stones.length <= 30`
- `1 <= stones[i] <= 1000`

---

**Hints — try each one before reading on:**
1. Always need the two largest — a max-heap.
2. Python's heapq is a min-heap: store negatives.

**Target:** O(n log n) time, O(n) space
