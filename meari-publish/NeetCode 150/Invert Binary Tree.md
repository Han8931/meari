---
created: "2026-06-07"
id: nc-invert-binary-tree
source: imported:neetcode-150
study:
  answer: |-
    Recursively swap node.left and node.right and invert both subtrees; the base case is a null node. Iterative BFS/DFS with an explicit swap works identically.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Invert Binary Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Invert Binary Tree
---

**Pattern:** Trees · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/invert-binary-tree)

Given the root of a binary tree, **invert** it (mirror every left/right pair) and return the root.

**Example 1:**

    Input: root = [4,2,7,1,3,6,9]
    Output: [4,7,2,9,6,3,1]

**Example 2:**

    Input: root = []
    Output: []

**Constraints:**

- Node count is in `[0, 100]`

---

**Hints — try each one before reading on:**
1. Swap the children, then recurse into them.
2. Or BFS, swapping at each visited node.

**Target:** O(n) time, O(h) space
