---
created: "2026-06-07"
id: nc-rotate-image
source: imported:neetcode-150
study:
  answer: |-
    Transpose the matrix (swap across the diagonal), then reverse every row — together they effect a 90° clockwise rotation with O(1) extra space.

    Complexity: O(n²) time, O(1) space
  kind: essay
  prompt: 'Solve "Rotate Image" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Rotate Image
---

**Pattern:** Math & Geometry · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/rotate-image)

Rotate an `n x n` matrix 90° clockwise, **in place**.

**Example 1:**

    Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]
    Output: [[7,4,1],[8,5,2],[9,6,3]]

**Constraints:**

- `1 <= n <= 20`

---

**Hints — try each one before reading on:**
1. Transpose, then reverse each row.
2. Or rotate four cells at a time around the rings.

**Target:** O(n²) time, O(1) space
