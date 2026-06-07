---
created: "2026-06-07"
id: nc-binary-tree-right-side-view
source: imported:neetcode-150
study:
  answer: |-
    Level-order BFS recording the last value of every level. (Equivalently DFS visiting right children first and recording the first node seen at each new depth.)

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Binary Tree Right Side View" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Binary Tree Right Side View
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/binary-tree-right-side-view)

Standing to the **right** of a binary tree, return the values you can see, ordered top to bottom.

**Example 1:**

    Input: root = [1,2,3,null,5,null,4]
    Output: [1,3,4]

**Example 2:**

    Input: root = [1,null,3]
    Output: [1,3]

**Constraints:**

- Node count is in `[0, 100]`

---

**Hints — try each one before reading on:**
1. That's each level's last node.
2. BFS per level, keep the final value; or DFS right-first taking the first node per depth.

**Target:** O(n) time, O(n) space
