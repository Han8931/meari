---
created: "2026-06-07"
id: nc-unique-paths
source: imported:neetcode-150
study:
  answer: |-
    DP over the grid where each cell is the sum of the cell above and to the left, edges seeded with 1 — kept as a single row updated in place. (Or compute the binomial coefficient directly.)

    Complexity: O(m·n) time, O(n) space
  kind: essay
  prompt: 'Solve "Unique Paths" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Unique Paths
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/unique-paths)

A robot starts at the top-left of an `m x n` grid and moves only **right or down**. How many distinct paths reach the bottom-right?

**Example 1:**

    Input: m = 3, n = 7
    Output: 28

**Example 2:**

    Input: m = 3, n = 2
    Output: 3

**Constraints:**

- `1 <= m, n <= 100`

---

**Hints — try each one before reading on:**
1. paths(cell) = paths(above) + paths(left).
2. One rolling row suffices; the closed form is C(m+n−2, m−1).

**Target:** O(m·n) time, O(n) space
