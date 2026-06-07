---
created: "2026-06-07"
id: nc-construct-binary-tree-from-preorder-and-inorder-traversal
source: imported:neetcode-150
study:
  answer: |-
    Recursive split: take the next preorder element as root, look up its inorder position via a prebuilt index map, and recurse on the left and right inorder ranges with a shared preorder cursor.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Construct Binary Tree from Preorder and Inorder Traversal" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Construct Binary Tree from Preorder and Inorder Traversal
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal)

Given the `preorder` and `inorder` traversals of a binary tree with **unique values**, rebuild the tree and return its root.

**Example 1:**

    Input: preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
    Output: [3,9,20,null,null,15,7]

**Constraints:**

- `1 <= preorder.length <= 3000`
- Values are unique; both arrays describe the same tree

---

**Hints — try each one before reading on:**
1. Preorder's first element is the root; find it in inorder to split left/right.
2. A value → inorder-index map removes the linear search.

**Target:** O(n) time, O(n) space
