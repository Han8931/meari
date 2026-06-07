---
created: "2026-06-07"
id: nc-set-matrix-zeroes
source: imported:neetcode-150
study:
  answer: |-
    First pass stores zero-markers in the first row/column (with a separate flag for column 0); second pass zeroes inner cells whose row or column is marked; finally handle the first row and column from the markers and flag.

    Complexity: O(m·n) time, O(1) space
  kind: essay
  prompt: 'Solve "Set Matrix Zeroes" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Set Matrix Zeroes
---

**Pattern:** Math & Geometry · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/set-matrix-zeroes)

If a matrix element is 0, set its **entire row and column** to 0 — in place, ideally with O(1) extra space.

**Example 1:**

    Input: matrix = [[1,1,1],[1,0,1],[1,1,1]]
    Output: [[1,0,1],[0,0,0],[1,0,1]]

**Example 2:**

    Input: matrix = [[0,1,2,0],[3,4,5,2],[1,3,1,5]]
    Output: [[0,0,0,0],[0,4,5,0],[0,3,1,0]]

**Constraints:**

- `1 <= m, n <= 200`

---

**Hints — try each one before reading on:**
1. Marking immediately cascades — record which rows/columns die first.
2. Use row 0 and column 0 themselves as the marker arrays (one extra flag for column 0).

**Target:** O(m·n) time, O(1) space
