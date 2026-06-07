---
created: "2026-06-07"
id: nc-two-sum-ii-input-array-is-sorted
source: imported:neetcode-150
study:
  answer: |-
    Pointers at both ends: if the sum is less than target advance the left pointer, if greater retreat the right; sortedness guarantees this never skips the answer.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Two Sum II – Input Array Is Sorted" (Two Pointers): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Two Sum II – Input Array Is Sorted
---

**Pattern:** Two Pointers · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/two-sum-ii-input-array-is-sorted)

Given a **1-indexed**, **sorted** array `numbers` and a `target`, return the 1-based indices of the two numbers adding to `target`. Exactly one solution exists; you must use only O(1) extra space.

**Example 1:**

    Input: numbers = [2,7,11,15], target = 9
    Output: [1,2]

**Example 2:**

    Input: numbers = [2,3,4], target = 6
    Output: [1,3]

**Constraints:**

- `2 <= numbers.length <= 3 * 10^4`
- `numbers` is sorted in non-decreasing order
- Exactly one solution exists

---

**Hints — try each one before reading on:**
1. Sorted input begs for pointers at both ends.
2. Sum too small → move left up; too big → move right down.

**Target:** O(n) time, O(1) space
