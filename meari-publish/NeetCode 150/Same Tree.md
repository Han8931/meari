---
created: "2026-06-07"
id: nc-same-tree
source: imported:neetcode-150
study:
  answer: |-
    Simultaneous recursion: same(p, q) = both null, or values equal and same(left, left) and same(right, right).

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Same Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Same Tree
---

**Pattern:** Trees · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/same-tree)

Given the roots of two binary trees `p` and `q`, return `true` if the trees are structurally identical with equal node values.

**Example 1:**

    Input: p = [1,2,3], q = [1,2,3]
    Output: true

**Example 2:**

    Input: p = [1,2], q = [1,null,2]
    Output: false

**Constraints:**

- Node count is in `[0, 100]`

---

**Hints — try each one before reading on:**
1. Compare roots, then recurse pairwise.
2. Both null → true; one null or values differ → false.

**Target:** O(n) time, O(h) space
