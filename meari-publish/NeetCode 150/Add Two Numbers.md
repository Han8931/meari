---
created: "2026-06-07"
id: nc-add-two-numbers
source: imported:neetcode-150
study:
  answer: |-
    Walk both lists with a carry, emitting (a + b + carry) % 10 nodes onto a dummy-headed result and carrying the quotient; continue while any list or a carry survives.

    Complexity: O(max(m, n)) time, O(1) extra space
  kind: essay
  prompt: 'Solve "Add Two Numbers" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Add Two Numbers
---

**Pattern:** Linked List · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/add-two-numbers)

Two non-negative integers are stored as linked lists in **reverse digit order**. Add them and return the sum as a list in the same format.

**Example 1:**

    Input: l1 = [2,4,3], l2 = [5,6,4]
    Output: [7,0,8]
    Explanation: 342 + 465 = 807.

**Example 2:**

    Input: l1 = [9,9,9], l2 = [1]
    Output: [0,0,0,1]

**Constraints:**

- Each list has `[1, 100]` nodes
- No leading zeros except the number 0 itself

---

**Hints — try each one before reading on:**
1. Grade-school addition: digit + digit + carry.
2. Loop while either list OR the carry remains.

**Target:** O(max(m, n)) time, O(1) extra space
