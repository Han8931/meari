---
created: "2026-06-07"
id: nc-maximum-depth-of-binary-tree
source: imported:neetcode-150
study:
  answer: |-
    Recursion: depth(node) = 1 + max(depth(left), depth(right)), 0 for null. BFS level counting is the iterative alternative.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Maximum Depth of Binary Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Maximum Depth of Binary Tree
---

**Pattern:** Trees · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/maximum-depth-of-binary-tree)

Return a binary tree's **maximum depth**: the number of nodes on the longest path from the root down to a leaf.

**Example 1:**

    Input: root = [3,9,20,null,null,15,7]
    Output: 3

**Constraints:**

- Node count is in `[0, 10^4]`

---

**Hints — try each one before reading on:**
1. A node's depth is 1 + the deeper child.
2. Null contributes 0.

**Target:** O(n) time, O(h) space
