---
created: "2026-06-07"
id: nc-number-of-connected-components-in-an-undirected-graph
source: imported:neetcode-150
study:
  answer: |-
    Union-Find starting at count = n: each union of two different roots decrements the count; the remainder is the answer. (Equivalent: DFS/BFS counting how many starts are needed.)

    Complexity: O(E α(V)) time, O(V) space
  kind: essay
  prompt: 'Solve "Number of Connected Components in an Undirected Graph" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Number of Connected Components in an Undirected Graph
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/number-of-connected-components-in-an-undirected-graph)

Given `n` nodes (0 to n−1) and an undirected edge list, return the number of **connected components**.

**Example 1:**

    Input: n = 5, edges = [[0,1],[1,2],[3,4]]
    Output: 2

**Example 2:**

    Input: n = 5, edges = [[0,1],[1,2],[2,3],[3,4]]
    Output: 1

**Constraints:**

- `1 <= n <= 2000`
- No repeated edges

---

**Hints — try each one before reading on:**
1. Start with n components; each successful union merges two.
2. Union-Find, or DFS from every unvisited node.

**Target:** O(E α(V)) time, O(V) space
