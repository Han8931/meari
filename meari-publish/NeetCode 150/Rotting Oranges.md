---
created: "2026-06-07"
id: nc-rotting-oranges
source: imported:neetcode-150
study:
  answer: |-
    Multi-source BFS: enqueue all rotten cells at minute 0, expand level by level rotting fresh neighbors, counting minutes per level. If fresh oranges remain at the end, return −1.

    Complexity: O(rows · cols) time, O(rows · cols) space
  kind: essay
  prompt: 'Solve "Rotting Oranges" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Rotting Oranges
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/rotting-oranges)

In a grid, `0` is empty, `1` a fresh orange, `2` a rotten one. Every minute, fresh oranges adjacent to rotten ones rot. Return the minutes until **no fresh orange remains**, or `-1` if impossible.

**Example 1:**

    Input: grid = [[2,1,1],[1,1,0],[0,1,1]]
    Output: 4

**Example 2:**

    Input: grid = [[2,1,1],[0,1,1],[1,0,1]]
    Output: -1
    Explanation: the bottom-left orange can never rot.

**Constraints:**

- `1 <= m, n <= 10`

---

**Hints — try each one before reading on:**
1. All rotten oranges spread simultaneously — multi-source BFS.
2. Seed the queue with every initially rotten orange; count levels.

**Target:** O(rows · cols) time, O(rows · cols) space
