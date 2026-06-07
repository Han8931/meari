---
created: "2026-06-07"
id: nc-swim-in-rising-water
source: imported:neetcode-150
study:
  answer: |-
    Modified Dijkstra: heap keyed by the bottleneck value max(path so far, neighbor height); the first time the goal pops, its key is the answer. (Alternative: binary search t, BFS over cells ≤ t.)

    Complexity: O(n² log n) time, O(n²) space
  kind: essay
  prompt: 'Solve "Swim in Rising Water" (Advanced Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Swim in Rising Water
---

**Pattern:** Advanced Graphs · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/swim-in-rising-water)

In an `n x n` grid, `grid[i][j]` is that cell's elevation. At time `t` you can move through cells with elevation `≤ t`. Return the **least time** to travel from `(0,0)` to `(n-1,n-1)`.

**Example 1:**

    Input: grid = [[0,2],[1,3]]
    Output: 3

**Example 2:**

    Input: grid = [[0,1,2,3,4],[24,23,22,21,5],[12,13,14,15,16],[11,17,18,19,20],[10,9,8,7,6]]
    Output: 16

**Constraints:**

- `1 <= n <= 50`
- Elevations are a permutation of `0 … n² - 1`

---

**Hints — try each one before reading on:**
1. Minimize the MAXIMUM cell along the path, not the sum.
2. Dijkstra where a path's cost is max(so far, next cell); or binary search the time + BFS.

**Target:** O(n² log n) time, O(n²) space
