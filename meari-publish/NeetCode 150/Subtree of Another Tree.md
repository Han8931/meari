---
created: "2026-06-07"
id: nc-subtree-of-another-tree
source: imported:neetcode-150
study:
  answer: |-
    For each node of root (DFS), run the Same Tree comparison against subRoot; true if any position matches. (Serialization with null markers + substring matching is the asymptotically faster alternative.)

    Complexity: O(m·n) time, O(h) space
  kind: essay
  prompt: 'Solve "Subtree of Another Tree" (Trees): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Subtree of Another Tree
---

**Pattern:** Trees · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/subtree-of-another-tree)

Given trees `root` and `subRoot`, return `true` if `root` has a subtree **identical** to `subRoot`.

**Example 1:**

    Input: root = [3,4,5,1,2], subRoot = [4,1,2]
    Output: true

**Example 2:**

    Input: root = [3,4,5,1,2,null,null,null,null,0], subRoot = [4,1,2]
    Output: false

**Constraints:**

- `root` has `[1, 2000]` nodes, `subRoot` has `[1, 1000]`

---

**Hints — try each one before reading on:**
1. Reuse Same Tree as the inner check.
2. Try it at every node of root.

**Target:** O(m·n) time, O(h) space
