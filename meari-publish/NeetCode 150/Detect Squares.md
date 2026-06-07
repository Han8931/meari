---
created: "2026-06-07"
id: nc-detect-squares
source: imported:neetcode-150
study:
  answer: |-
    Keep a Counter of points. For query (x, y), scan distinct stored points (x2, y2) that are diagonal (|x−x2| == |y−y2| > 0) and add count(diagonal) · count((x, y2)) · count((x2, y)). Scanning distinct points keeps it fast.

    Complexity: O(distinct points) per count, O(points) space
  kind: essay
  prompt: 'Solve "Detect Squares" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Detect Squares
---

**Pattern:** Math & Geometry · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/detect-squares)

Design a structure with `add(point)` and `count(point)`: count returns how many **axis-aligned squares** with positive area can be formed using the query point and three previously added points (duplicates allowed and counted).

**Example 1:**

    Input: add([3,10]); add([11,2]); add([3,2]); count([11,10])
    Output: 1
    Explanation: the square (3,10),(11,2),(3,2),(11,10).

**Constraints:**

- `0 <= x, y <= 1000`
- Up to `3000` calls

---

**Hints — try each one before reading on:**
1. A square is fixed by the query and any DIAGONAL point.
2. For each stored point (x2, y2) with |dx| == |dy| ≠ 0, check the two completing corners; multiply counts.

**Target:** O(distinct points) per count, O(points) space
