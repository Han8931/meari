---
created: "2026-06-07"
id: nc-valid-sudoku
source: imported:neetcode-150
study:
  answer: |-
    Single scan keeping seen-digit sets (or 9-bit bitmasks) per row, per column, and per box keyed by (r//3, c//3); a digit already present in any of its three units makes the board invalid. The bitmask version tests membership with (1 << d) & mask and inserts with |=.

    Complexity: O(81) ≈ O(1) time, O(1) space
  kind: essay
  prompt: 'Solve "Valid Sudoku" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Valid Sudoku
---

**Pattern:** Arrays & Hashing · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/valid-sudoku)

Determine whether a `9 x 9` Sudoku board is **valid**: each row, each column, and each of the nine `3 x 3` boxes must contain the digits 1-9 without repetition. Empty cells are `'.'` — only filled cells are checked, and a valid board need not be solvable. (Your `Array_Hashing.md` note compares the brute-force, hash-set, and bitmask versions of this one.)

**Example 1:**

    Input: board[0] = ["5","3",".",".","7",".",".",".","."], ... (9 rows)
    Output: true

**Example 2:**

    Input: same board but board[0][0] = "8" (clashes with the 8 below in its column/box)
    Output: false

**Constraints:**

- `board.length == board[i].length == 9`
- `board[i][j]` is a digit `1-9` or `'.'`

---

**Hints — try each one before reading on:**
1. One pass: every cell belongs to exactly one row set, one column set, one box set.
2. Box index is (r // 3, c // 3); a 9-bit int per unit replaces each set.

**Target:** O(81) ≈ O(1) time, O(1) space
