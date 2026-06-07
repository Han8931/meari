---
created: "2026-06-07"
id: nc-3sum
source: imported:neetcode-150
study:
  answer: |-
    Sort the array. For each index i (skipping repeated values), run the two-pointer scan on the suffix for pairs summing to −nums[i], skipping duplicates after each hit. Early-exit when nums[i] > 0.

    Complexity: O(n²) time, O(1) extra space
  kind: essay
  prompt: 'Solve "3Sum" (Two Pointers): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: 3Sum
---

**Pattern:** Two Pointers · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/3sum)

Given an integer array `nums`, return **all unique triplets** `[nums[i], nums[j], nums[k]]` (i, j, k distinct) whose sum is zero.

**Example 1:**

    Input: nums = [-1,0,1,2,-1,-4]
    Output: [[-1,-1,2],[-1,0,1]]

**Example 2:**

    Input: nums = [0,1,1]
    Output: []

**Constraints:**

- `3 <= nums.length <= 3000`
- `-10^5 <= nums[i] <= 10^5`

---

**Hints — try each one before reading on:**
1. Sort, then fix one element and the problem becomes Two Sum II.
2. Skip duplicate values at every level to keep triplets unique.

**Target:** O(n²) time, O(1) extra space
