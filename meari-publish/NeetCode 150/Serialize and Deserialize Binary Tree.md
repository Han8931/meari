---
created: "2026-06-07"
id: nc-serialize-and-deserialize-binary-tree
source: imported:neetcode-150
study:
  answer: |-
    Serialize via preorder DFS emitting values and 'N' for nulls, comma-joined. Deserialize with an iterator over the tokens: read one token, return None for the marker, else build the node and recurse left then right — the same shape as the writer.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Serialize and Deserialize Binary Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Serialize and Deserialize Binary Tree
---

**Pattern:** Trees · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/serialize-and-deserialize-binary-tree)

Design `serialize(root)` encoding a binary tree to a string, and `deserialize(data)` reconstructing exactly the original tree. Any format works as long as the round trip is lossless.

**Example 1:**

    Input: root = [1,2,3,null,null,4,5]
    Output (after round trip): [1,2,3,null,null,4,5]

**Constraints:**

- Node count is in `[0, 10^4]`
- `-1000 <= Node.val <= 1000`

---

**Hints — try each one before reading on:**
1. Preorder with explicit null markers is unambiguous.
2. Deserialize by consuming tokens with an iterator in the same preorder.

**Target:** O(n) time, O(n) space
