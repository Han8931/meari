---
created: "2026-06-07"
id: nc-number-of-islands
source: imported:neetcode-150
study:
  answer: |-
    Scan the grid; for every '1' not yet visited, increment the count and flood-fill its whole component (DFS or BFS over 4-neighbors), marking cells as seen.

    Complexity: O(rows · cols) time, O(rows · cols) space worst case
  kind: essay
  prompt: 'Solve "Number of Islands" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Number of Islands
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/number-of-islands)

Given an `m x n` grid of `'1'` (land) and `'0'` (water), return the number of **islands** — groups of land connected horizontally or vertically.

**Example 1:**

    Input: grid = [["1","1","0","0"],["1","1","0","0"],["0","0","1","0"],["0","0","0","1"]]
    Output: 3

**Constraints:**

- `1 <= m, n <= 300`

---

**Hints — try each one before reading on:**
1. Each unvisited land cell starts one island — flood it entirely.
2. DFS/BFS marking visited (or sinking '1' to '0') prevents recounting.

**Target:** O(rows · cols) time, O(rows · cols) space worst case
