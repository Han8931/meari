---
created: "2026-06-07"
id: nc-merge-two-sorted-lists
source: imported:neetcode-150
study:
  answer: |-
    Dummy node + tail pointer: repeatedly attach the smaller head of the two lists and advance; attach whatever remains at the end. Return dummy.next.

    Complexity: O(n + m) time, O(1) space
  kind: essay
  prompt: 'Solve "Merge Two Sorted Lists" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Merge Two Sorted Lists
---

**Pattern:** Linked List · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/merge-two-sorted-lists)

Merge two sorted linked lists `list1` and `list2` into one sorted list by splicing their nodes, and return its head.

**Example 1:**

    Input: list1 = [1,2,4], list2 = [1,3,4]
    Output: [1,1,2,3,4,4]

**Example 2:**

    Input: list1 = [], list2 = [0]
    Output: [0]

**Constraints:**

- Each list has `[0, 50]` nodes
- Both lists are sorted in non-decreasing order

---

**Hints — try each one before reading on:**
1. A dummy head removes every "is this the first node?" special case.
2. Append the smaller of the two fronts, advance that list.

**Target:** O(n + m) time, O(1) space
