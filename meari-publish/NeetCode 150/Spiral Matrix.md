---
created: "2026-06-07"
id: nc-spiral-matrix
source: imported:neetcode-150
study:
  answer: |-
    Peel layers with four boundary pointers: left→right along top (top++), top→bottom along right (right--), then — if rows/cols remain — right→left along bottom (bottom--) and bottom→top along left (left++). Stop when boundaries cross.

    Complexity: O(m·n) time, O(1) extra space
  kind: essay
  prompt: 'Solve "Spiral Matrix" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Spiral Matrix
---

**Pattern:** Math & Geometry · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/spiral-matrix)

Return all elements of an `m x n` matrix in **spiral order**.

**Example 1:**

    Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]
    Output: [1,2,3,6,9,8,7,4,5]

**Example 2:**

    Input: matrix = [[1,2,3,4],[5,6,7,8],[9,10,11,12]]
    Output: [1,2,3,4,8,12,11,10,9,5,6,7]

**Constraints:**

- `1 <= m, n <= 10`

---

**Hints — try each one before reading on:**
1. Maintain four boundaries: top, bottom, left, right.
2. Walk a side, shrink its boundary, and re-check boundaries between sides.

**Target:** O(m·n) time, O(1) extra space
