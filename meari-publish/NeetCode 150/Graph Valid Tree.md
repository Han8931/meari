---
created: "2026-06-07"
id: nc-graph-valid-tree
source: imported:neetcode-150
study:
  answer: |-
    Reject immediately unless len(edges) == n−1; then verify connectivity (DFS/BFS from node 0 reaches all n, or Union-Find ends at one component). Both conditions together imply acyclicity.

    Complexity: O(V + E) time, O(V) space
  kind: essay
  prompt: 'Solve "Graph Valid Tree" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Graph Valid Tree
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/graph-valid-tree)

Given `n` nodes and an undirected edge list, return `true` if the edges form a **valid tree**: connected and acyclic.

**Example 1:**

    Input: n = 5, edges = [[0,1],[0,2],[0,3],[1,4]]
    Output: true

**Example 2:**

    Input: n = 5, edges = [[0,1],[1,2],[2,3],[1,3],[1,4]]
    Output: false

**Constraints:**

- `1 <= n <= 2000`
- No self-loops or repeated edges

---

**Hints — try each one before reading on:**
1. A tree has exactly n−1 edges AND is fully connected.
2. Either check beats checking for cycles directly when combined with the other.

**Target:** O(V + E) time, O(V) space
