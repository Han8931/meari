---
created: "2026-06-07"
id: nc-count-good-nodes-in-binary-tree
source: imported:neetcode-150
study:
  answer: |-
    DFS passing down maxSoFar: a node is good when node.val ≥ maxSoFar; recurse with max(maxSoFar, node.val) and sum the counts.

    Complexity: O(n) time, O(h) space
  kind: essay
  prompt: 'Solve "Count Good Nodes in Binary Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Count Good Nodes in Binary Tree
---

**Pattern:** Trees · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/count-good-nodes-in-binary-tree)

A node X is **good** if the path from the root to X contains no node with a value greater than X's. Return the number of good nodes.

**Example 1:**

    Input: root = [3,1,4,3,null,1,5]
    Output: 4
    Explanation: 3 (root), 4, 5, and the leaf 3.

**Example 2:**

    Input: root = [3,3,null,4,2]
    Output: 3

**Constraints:**

- Node count is in `[1, 10^5]`
- `-10^4 <= Node.val <= 10^4`

---

**Hints — try each one before reading on:**
1. Carry the max value seen along the path down the recursion.
2. A node counts iff node.val ≥ that running max.

**Target:** O(n) time, O(h) space
