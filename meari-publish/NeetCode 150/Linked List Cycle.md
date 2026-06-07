---
created: "2026-06-07"
id: nc-linked-list-cycle
source: imported:neetcode-150
study:
  answer: |-
    Floyd's tortoise and hare: advance slow by one and fast by two; if they ever coincide there's a cycle, and reaching null proves there isn't.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Linked List Cycle" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Linked List Cycle
---

**Pattern:** Linked List · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/linked-list-cycle)

Given the head of a linked list, return `true` if following `next` pointers ever loops back to an earlier node.

**Example 1:**

    Input: head = [3,2,0,-4] with tail connecting to index 1
    Output: true

**Example 2:**

    Input: head = [1] with no cycle
    Output: false

**Constraints:**

- List length is in `[0, 10^4]`
- Follow-up: O(1) memory

---

**Hints — try each one before reading on:**
1. Two runners at different speeds must meet inside a loop.
2. Floyd: slow +1, fast +2; meeting ⇒ cycle, null ⇒ no cycle.

**Target:** O(n) time, O(1) space
