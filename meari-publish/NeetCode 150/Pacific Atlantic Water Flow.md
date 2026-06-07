---
created: "2026-06-07"
id: nc-pacific-atlantic-water-flow
source: imported:neetcode-150
study:
  answer: |-
    Reverse the flow: DFS/BFS from all Pacific-edge cells moving only to higher-or-equal neighbors marks everything that drains to the Pacific; repeat for the Atlantic; answer = cells in both sets.

    Complexity: O(rows · cols) time, O(rows · cols) space
  kind: essay
  prompt: 'Solve "Pacific Atlantic Water Flow" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Pacific Atlantic Water Flow
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/pacific-atlantic-water-flow)

An `m x n` height grid touches the **Pacific** on its top/left edges and the **Atlantic** on its bottom/right. Rain flows to neighbors of **equal or lower** height. Return all cells from which water can reach **both** oceans.

**Example 1:**

    Input: heights = [[1,2,2,3,5],[3,2,3,4,4],[2,4,5,3,1],[6,7,1,4,5],[5,1,1,2,4]]
    Output: [[0,4],[1,3],[1,4],[2,2],[3,0],[3,1],[4,0]]

**Constraints:**

- `1 <= m, n <= 200`
- `0 <= heights[r][c] <= 10^5`

---

**Hints — try each one before reading on:**
1. Don't simulate from every cell — flow BACKWARD from each ocean.
2. BFS/DFS uphill from both coastlines; intersect the two reachable sets.

**Target:** O(rows · cols) time, O(rows · cols) space
