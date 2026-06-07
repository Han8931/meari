---
created: "2026-06-07"
id: nc-remove-nth-node-from-end-of-list
source: imported:neetcode-150
study:
  answer: |-
    Dummy head; advance a fast pointer n+1 steps, then move fast and slow together until fast is null — slow.next is the node to remove; splice it out and return dummy.next.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Remove Nth Node From End of List" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Remove Nth Node From End of List
---

**Pattern:** Linked List · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/remove-nth-node-from-end-of-list)

Remove the `n`-th node **from the end** of a linked list and return the head — ideally in one pass.

**Example 1:**

    Input: head = [1,2,3,4,5], n = 2
    Output: [1,2,3,5]

**Example 2:**

    Input: head = [1], n = 1
    Output: []

**Constraints:**

- List length `sz` is in `[1, 30]`
- `1 <= n <= sz`

---

**Hints — try each one before reading on:**
1. Two pointers n apart hit the end together.
2. Start the leading pointer n steps ahead from a dummy; when it hits the end, the trailing one sits just before the victim.

**Target:** O(n) time, O(1) space
