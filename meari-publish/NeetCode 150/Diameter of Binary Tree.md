---
created: "2026-06-07"
id: nc-diameter-of-binary-tree
source: imported:neetcode-150
study:
  answer: |-
    Post-order DFS returning each subtree's height while updating best = max(best, hL + hR) at every node. The answer is the global best, not the root's value.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Diameter of Binary Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Diameter of Binary Tree
---

**Pattern:** Trees · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/diameter-of-binary-tree)

Return the **diameter** of a binary tree: the number of edges on the longest path between any two nodes (the path need not pass the root).

**Example 1:**

    Input: root = [1,2,3,4,5]
    Output: 3
    Explanation: the path 4 → 2 → 1 → 3 (or 5 → 2 → 1 → 3).

**Constraints:**

- Node count is in `[1, 10^4]`

---

**Hints — try each one before reading on:**
1. Through any node, the longest path is leftHeight + rightHeight.
2. Compute heights bottom-up; update a global best at every node.

**Target:** O(n) time, O(h) space
