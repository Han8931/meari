---
created: "2026-06-07"
id: nc-binary-tree-maximum-path-sum
source: imported:neetcode-150
study:
  answer: |-
    Post-order DFS returning the best downward gain (node.val + max(0, leftGain, rightGain)) while updating a global best with node.val + max(0,leftGain) + max(0,rightGain) — the through-path that can't extend upward.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Binary Tree Maximum Path Sum" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Binary Tree Maximum Path Sum
---

**Pattern:** Trees · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/binary-tree-maximum-path-sum)

A **path** is any node sequence connected by edges, appearing at most once each; it need not pass the root. Return the maximum possible path sum.

**Example 1:**

    Input: root = [1,2,3]
    Output: 6
    Explanation: 2 + 1 + 3.

**Example 2:**

    Input: root = [-10,9,20,null,null,15,7]
    Output: 42
    Explanation: 15 + 20 + 7.

**Constraints:**

- Node count is in `[1, 3 * 10^4]`
- `-1000 <= Node.val <= 1000`

---

**Hints — try each one before reading on:**
1. At each node a path either passes THROUGH it (left gain + val + right gain) or extends upward (val + one side).
2. Clamp negative child gains to zero.

**Target:** O(n) time, O(h) space
