---
created: "2026-06-07"
id: nc-balanced-binary-tree
source: imported:neetcode-150
study:
  answer: |-
    DFS returning height but propagating a sentinel (−1) the moment a subtree is unbalanced or children differ by more than one — one bottom-up pass instead of recomputing heights per node.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Balanced Binary Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Balanced Binary Tree
---

**Pattern:** Trees · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/balanced-binary-tree)

Return `true` if the tree is **height-balanced**: at every node, the two subtree heights differ by at most one.

**Example 1:**

    Input: root = [3,9,20,null,null,15,7]
    Output: true

**Example 2:**

    Input: root = [1,2,2,3,3,null,null,4,4]
    Output: false

**Constraints:**

- Node count is in `[0, 5000]`

---

**Hints — try each one before reading on:**
1. Height and balance can be computed in the same pass.
2. Return −1 (or a flag) upward as soon as any subtree is unbalanced.

**Target:** O(n) time, O(h) space
