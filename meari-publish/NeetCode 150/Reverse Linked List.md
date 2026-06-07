---
created: "2026-06-07"
id: nc-reverse-linked-list
source: imported:neetcode-150
study:
  answer: |-
    Iterate with prev = None: save cur.next, point cur.next at prev, advance both. When cur runs out, prev is the new head. (Recursion mirrors the same idea with O(n) stack.)

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Reverse Linked List" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Reverse Linked List
---

**Pattern:** Linked List · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/reverse-linked-list)

Given the head of a singly linked list, reverse it and return the new head.

**Example 1:**

    Input: head = [1,2,3,4,5]
    Output: [5,4,3,2,1]

**Example 2:**

    Input: head = []
    Output: []

**Constraints:**

- List length is in `[0, 5000]`
- `-5000 <= Node.val <= 5000`

---

**Hints — try each one before reading on:**
1. Walk the list re-pointing each node's next at the previous node.
2. Three names: prev, cur, and a saved next.

**Target:** O(n) time, O(1) space
