---
created: "2026-06-07"
id: nc-max-area-of-island
source: imported:neetcode-150
study:
  answer: |-
    Flood-fill every unvisited land cell, where DFS returns the component's cell count (1 plus recursive neighbor sums, 0 for water/visited); track the maximum.

    Complexity: O(rows · cols) time, O(rows · cols) space
  kind: essay
  prompt: 'Solve "Max Area of Island" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Max Area of Island
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/max-area-of-island)

Given a binary grid, return the **area of the largest island** (4-directionally connected 1s), or 0 if there is no land.

**Example 1:**

    Input: grid = [[0,1,0],[1,1,0],[0,0,1]]
    Output: 3
    Explanation: the L-shaped island of three 1s.

**Constraints:**

- `1 <= m, n <= 50`

---

**Hints — try each one before reading on:**
1. Same flood fill as Number of Islands, but return the size.
2. DFS returns 1 + the four neighbor areas.

**Target:** O(rows · cols) time, O(rows · cols) space
