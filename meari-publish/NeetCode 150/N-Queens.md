---
created: "2026-06-07"
id: nc-n-queens
source: imported:neetcode-150
study:
  answer: |-
    Backtrack over rows with three occupancy sets (cols, r+c, r−c): try each free column, place, recurse, remove. Materialize the board from the chosen column list at each complete placement.

    Complexity: O(n!) time, O(n) space
  kind: essay
  prompt: 'Solve "N-Queens" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: N-Queens
---

**Pattern:** Backtracking · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/n-queens)

Place `n` queens on an `n x n` chessboard so that no two attack each other. Return **all distinct** board configurations, using `'Q'` and `'.'`.

**Example 1:**

    Input: n = 4
    Output: [[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]

**Example 2:**

    Input: n = 1
    Output: [["Q"]]

**Constraints:**

- `1 <= n <= 9`

---

**Hints — try each one before reading on:**
1. One queen per row — recurse row by row choosing a column.
2. Constant-time safety: sets for used columns, r+c diagonals, r−c anti-diagonals.

**Target:** O(n!) time, O(n) space
