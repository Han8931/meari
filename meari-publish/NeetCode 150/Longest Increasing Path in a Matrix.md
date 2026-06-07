---
created: "2026-06-07"
id: nc-longest-increasing-path-in-a-matrix
source: imported:neetcode-150
study:
  answer: |-
    DFS from every cell with a memo table: lip(r, c) = 1 + max(lip of neighbors with larger value, default 0); the memo makes each cell computed once. The answer is the max over all cells.

    Complexity: O(m·n) time, O(m·n) space
  kind: essay
  prompt: 'Solve "Longest Increasing Path in a Matrix" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Increasing Path in a Matrix
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/longest-increasing-path-in-a-matrix)

Return the length of the longest **strictly increasing path** in a matrix, moving up/down/left/right.

**Example 1:**

    Input: matrix = [[9,9,4],[6,6,8],[2,1,1]]
    Output: 4
    Explanation: 1 → 2 → 6 → 9.

**Example 2:**

    Input: matrix = [[3,4,5],[3,2,6],[2,2,1]]
    Output: 4

**Constraints:**

- `1 <= m, n <= 200`
- `0 <= matrix[i][j] <= 2^31 - 1`

---

**Hints — try each one before reading on:**
1. Increasing paths can't revisit — the implicit graph is a DAG.
2. Memoized DFS: each cell's answer is 1 + best strictly-larger neighbor.

**Target:** O(m·n) time, O(m·n) space
