---
created: "2026-06-07"
id: nc-redundant-connection
source: imported:neetcode-150
study:
  answer: |-
    Union-Find: process edges in order, union the endpoints; the first edge whose endpoints already share a root is the answer (the problem guarantees returning the last such edge in input order works).

    Complexity: O(E α(V)) time, O(V) space
  kind: essay
  prompt: 'Solve "Redundant Connection" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Redundant Connection
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/redundant-connection)

A tree of `n` nodes had **one extra edge** added. Given the resulting edge list, return the edge that can be removed to restore a tree (the last such edge in input order).

**Example 1:**

    Input: edges = [[1,2],[1,3],[2,3]]
    Output: [2,3]

**Example 2:**

    Input: edges = [[1,2],[2,3],[3,4],[1,4],[1,5]]
    Output: [1,4]

**Constraints:**

- `3 <= n <= 1000`
- `edges.length == n`, no repeated edges

---

**Hints — try each one before reading on:**
1. Adding edges one by one, the first edge joining two already-connected nodes is the culprit.
2. Union-Find with path compression answers connectivity in near-O(1).

**Target:** O(E α(V)) time, O(V) space
