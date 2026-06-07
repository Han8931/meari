---
created: "2026-06-07"
id: nc-validate-binary-search-tree
source: imported:neetcode-150
study:
  answer: |-
    DFS with bounds: valid(node, lo, hi) requires lo < node.val < hi, recursing left with (lo, val) and right with (val, hi). The in-order "strictly increasing" check is an equivalent formulation.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Validate Binary Search Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Validate Binary Search Tree
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/validate-binary-search-tree)

Return `true` if a binary tree is a **valid BST**: every node's left subtree holds only smaller values, its right subtree only larger ones, recursively.

**Example 1:**

    Input: root = [2,1,3]
    Output: true

**Example 2:**

    Input: root = [5,1,4,null,null,3,6]
    Output: false
    Explanation: 3 sits in 5's right subtree but 3 < 5.

**Constraints:**

- Node count is in `[1, 10^4]`
- `-2^31 <= Node.val <= 2^31 - 1`

---

**Hints — try each one before reading on:**
1. Checking only parent vs child misses violations across generations.
2. Pass down an open interval (lo, hi) each node must fall in — or check that in-order traversal is strictly increasing.

**Target:** O(n) time, O(h) space
