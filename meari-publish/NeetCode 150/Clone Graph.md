---
created: "2026-06-07"
id: nc-clone-graph
source: imported:neetcode-150
study:
  answer: |-
    DFS (or BFS) with a hash map old → clone: visiting a node creates its clone if absent, then appends clone(neighbor) for every neighbor, returning the memoized clone on revisits.

    Complexity: O(V + E) time, O(V) space
  kind: essay
  prompt: 'Solve "Clone Graph" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Clone Graph
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/clone-graph)

Given a reference node of a **connected undirected graph** (each node has a value and a neighbor list), return a **deep copy** of the entire graph.

**Example 1:**

    Input: adjList = [[2,4],[1,3],[2,4],[1,3]]
    Output: the same adjacency structure built from new nodes

**Constraints:**

- Node count is in `[0, 100]`
- No self-loops or repeated edges

---

**Hints — try each one before reading on:**
1. old → clone mapping prevents infinite loops on cycles.
2. DFS: clone the node, then clone/look up each neighbor.

**Target:** O(V + E) time, O(V) space
