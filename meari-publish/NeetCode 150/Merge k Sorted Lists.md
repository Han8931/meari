---
created: "2026-06-07"
id: nc-merge-k-sorted-lists
source: imported:neetcode-150
study:
  answer: |-
    Min-heap seeded with each list's head; repeatedly pop the smallest, append it, and push its successor. (Divide-and-conquer pairwise merging achieves the same bound without a heap.)

    Complexity: O(N log k) time, O(k) space
  kind: essay
  prompt: 'Solve "Merge k Sorted Lists" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Merge k Sorted Lists
---

**Pattern:** Linked List · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/merge-k-sorted-lists)

Given an array of `k` sorted linked lists, merge them into one sorted list and return it.

**Example 1:**

    Input: lists = [[1,4,5],[1,3,4],[2,6]]
    Output: [1,1,2,3,4,4,5,6]

**Example 2:**

    Input: lists = []
    Output: []

**Constraints:**

- `0 <= k <= 10^4`
- Total nodes up to `10^4`

---

**Hints — try each one before reading on:**
1. Always take the smallest current head — a min-heap of k candidates.
2. Push (val, i, node) so ties don't compare nodes.

**Target:** O(N log k) time, O(k) space
