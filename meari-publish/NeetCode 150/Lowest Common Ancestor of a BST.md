---
created: "2026-06-07"
id: nc-lowest-common-ancestor-of-a-bst
source: imported:neetcode-150
study:
  answer: |-
    Walk from the root: if both values are smaller go left, both larger go right; otherwise the current node is where they diverge — the LCA. No parent pointers or extra storage needed.

    Complexity: O(h) time, O(1) space
  kind: essay
  prompt: 'Solve "Lowest Common Ancestor of a BST" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Lowest Common Ancestor of a BST
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree)

Given a **binary search tree** and two of its nodes `p` and `q`, return their **lowest common ancestor** — the deepest node having both as descendants (a node counts as its own descendant).

**Example 1:**

    Input: root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 8
    Output: 6

**Example 2:**

    Input: same tree, p = 2, q = 4
    Output: 2

**Constraints:**

- Node count is in `[2, 10^5]`
- All values unique; `p != q`, both exist in the tree

---

**Hints — try each one before reading on:**
1. BST ordering tells you which side both nodes live on.
2. The first node where p and q split (or equals one of them) is the LCA.

**Target:** O(h) time, O(1) space
