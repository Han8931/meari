---
created: "2026-06-07"
id: nc-reorder-list
source: imported:neetcode-150
study:
  answer: |-
    Find the middle (slow/fast), split and reverse the second half, then weave the two halves alternately. Each phase is a standard O(n) pass.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Reorder List" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Reorder List
---

**Pattern:** Linked List · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/reorder-list)

Given a list L0 → L1 → … → Ln, reorder it **in place** to L0 → Ln → L1 → Ln−1 → L2 → Ln−2 → … (only links may change, not values).

**Example 1:**

    Input: head = [1,2,3,4]
    Output: [1,4,2,3]

**Example 2:**

    Input: head = [1,2,3,4,5]
    Output: [1,5,2,4,3]

**Constraints:**

- List length is in `[1, 5 * 10^4]`

---

**Hints — try each one before reading on:**
1. Three known tricks compose: find the middle, reverse the second half, interleave.
2. Slow/fast pointers find the middle in one pass.

**Target:** O(n) time, O(1) space
