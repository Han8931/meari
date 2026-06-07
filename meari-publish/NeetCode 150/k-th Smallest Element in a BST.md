---
created: "2026-06-07"
id: nc-k-th-smallest-element-in-a-bst
source: imported:neetcode-150
study:
  answer: |-
    Iterative in-order traversal with an explicit stack: push left spine, pop, decrement k; when k hits zero the popped node is the answer. No need to traverse the rest.

    Complexity: O(h + k) time, O(h) space
  kind: essay
  prompt: 'Solve "k-th Smallest Element in a BST" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: k-th Smallest Element in a BST
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/kth-smallest-element-in-a-bst)

Given the root of a BST and an integer `k`, return the `k`-th smallest value (1-indexed).

**Example 1:**

    Input: root = [3,1,4,null,2], k = 1
    Output: 1

**Example 2:**

    Input: root = [5,3,6,2,4,null,null,1], k = 3
    Output: 3

**Constraints:**

- Node count is `n`, `1 <= k <= n <= 10^4`

---

**Hints — try each one before reading on:**
1. In-order traversal of a BST yields sorted order.
2. Stop the traversal at the k-th visit — iterative with a stack makes early exit easy.

**Target:** O(h + k) time, O(h) space
