---
created: "2026-06-07"
id: nc-search-a-2d-matrix
source: imported:neetcode-150
study:
  answer: |-
    Treat the matrix as a flat sorted array of length n·m and binary search once, translating mid to (mid // cols, mid % cols). (Two stacked searches — row first, then within it — also works.)

    Complexity: O(log (n·m)) time, O(1) space
  kind: essay
  prompt: 'Solve "Search a 2D Matrix" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Search a 2D Matrix
---

**Pattern:** Binary Search · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/search-a-2d-matrix)

Search for `target` in an `m x n` matrix where each row is sorted and each row's first element is greater than the previous row's last — in O(log (m·n)).

**Example 1:**

    Input: matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
    Output: true

**Example 2:**

    Input: same matrix, target = 13
    Output: false

**Constraints:**

- `1 <= m, n <= 100`
- `-10^4 <= matrix[i][j], target <= 10^4`

---

**Hints — try each one before reading on:**
1. That ordering makes the whole matrix one sorted list.
2. Binary search over n·m virtual indices: row = mid // m, col = mid % m.

**Target:** O(log (n·m)) time, O(1) space
