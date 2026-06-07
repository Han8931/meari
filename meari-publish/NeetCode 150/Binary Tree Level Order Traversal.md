---
created: "2026-06-07"
id: nc-binary-tree-level-order-traversal
source: imported:neetcode-150
study:
  answer: |-
    BFS: each round pops exactly len(queue) nodes (one level), collecting values and enqueueing children; append each level's list to the result.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Binary Tree Level Order Traversal" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Binary Tree Level Order Traversal
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/binary-tree-level-order-traversal)

Return the **level order traversal** of a binary tree's values — one list per level, left to right.

**Example 1:**

    Input: root = [3,9,20,null,null,15,7]
    Output: [[3],[9,20],[15,7]]

**Example 2:**

    Input: root = []
    Output: []

**Constraints:**

- Node count is in `[0, 2000]`

---

**Hints — try each one before reading on:**
1. BFS with a queue.
2. Freeze the queue's length per round to know where the level ends.

**Target:** O(n) time, O(n) space
