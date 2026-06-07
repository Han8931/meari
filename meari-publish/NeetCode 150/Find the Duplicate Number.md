---
created: "2026-06-07"
id: nc-find-the-duplicate-number
source: imported:neetcode-150
study:
  answer: |-
    Run Floyd's cycle detection on the function i → nums[i]: phase 1 finds a meeting point inside the cycle; phase 2 restarts one pointer at index 0 and advances both by one — they meet exactly at the duplicate (the cycle's entrance).

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Find the Duplicate Number" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Find the Duplicate Number
---

**Pattern:** Linked List · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/find-the-duplicate-number)

`nums` holds `n + 1` integers, each in `[1, n]`, so at least one value repeats; exactly one value is duplicated (possibly multiple times). Find it **without modifying the array** and using O(1) space.

**Example 1:**

    Input: nums = [1,3,4,2,2]
    Output: 2

**Example 2:**

    Input: nums = [3,1,3,4,2]
    Output: 3

**Constraints:**

- `1 <= n <= 10^5`
- `nums.length == n + 1`
- The array is read-only; O(1) extra space

---

**Hints — try each one before reading on:**
1. Treat i → nums[i] as a linked list; a duplicate value is two arrows into one node — a cycle.
2. Floyd again: find the meeting point, then reset one pointer to the start; they re-meet at the cycle entrance.

**Target:** O(n) time, O(1) space
