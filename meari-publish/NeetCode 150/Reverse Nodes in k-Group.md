---
created: "2026-06-07"
id: nc-reverse-nodes-in-k-group
source: imported:neetcode-150
study:
  answer: |-
    With a dummy head, for each group: probe k ahead (stop if the tail is short), reverse the k nodes, then reconnect — groupPrev.next becomes the old k-th node and the old group head links to the next group's first node; advance groupPrev to the old head.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Reverse Nodes in k-Group" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Reverse Nodes in k-Group
---

**Pattern:** Linked List · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/reverse-nodes-in-k-group)

Reverse the nodes of a linked list `k` at a time and return the modified list. Nodes in a final group shorter than `k` keep their order. You may not alter node values.

**Example 1:**

    Input: head = [1,2,3,4,5], k = 2
    Output: [2,1,4,3,5]

**Example 2:**

    Input: head = [1,2,3,4,5], k = 3
    Output: [3,2,1,4,5]

**Constraints:**

- List length is `n`, `1 <= k <= n <= 5000`

---

**Hints — try each one before reading on:**
1. Per group: check k nodes exist, reverse them, splice into the neighbors.
2. A dummy head plus a groupPrev pointer keeps the splicing manageable.

**Target:** O(n) time, O(1) space
