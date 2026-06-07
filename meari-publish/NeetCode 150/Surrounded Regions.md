---
created: "2026-06-07"
id: nc-surrounded-regions
source: imported:neetcode-150
study:
  answer: |-
    DFS/BFS from every border 'O' marking its region as safe (e.g. 'T'); then sweep the grid flipping unmarked 'O'→'X' and restoring 'T'→'O'. Solving the complement avoids tracking enclosure directly.

    Complexity: O(rows · cols) time, O(rows · cols) space
  kind: essay
  prompt: 'Solve "Surrounded Regions" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Surrounded Regions
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/surrounded-regions)

Given an `m x n` board of `'X'` and `'O'`, **capture** every region of O's that is fully surrounded by X — flip those O's to X. Regions touching the border survive.

**Example 1:**

    Input: board = [["X","X","X","X"],["X","O","O","X"],["X","X","O","X"],["X","O","X","X"]]
    Output: [["X","X","X","X"],["X","X","X","X"],["X","X","X","X"],["X","O","X","X"]]
    Explanation: the bottom-left O touches the border and survives.

**Constraints:**

- `1 <= m, n <= 200`

---

**Hints — try each one before reading on:**
1. A region survives iff it touches the border.
2. Mark border-connected 'O's first; everything else flips.

**Target:** O(rows · cols) time, O(rows · cols) space
