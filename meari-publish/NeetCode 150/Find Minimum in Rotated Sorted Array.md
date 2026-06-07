---
created: "2026-06-07"
id: nc-find-minimum-in-rotated-sorted-array
source: imported:neetcode-150
study:
  answer: |-
    Binary search: if nums[mid] > nums[r] the minimum lies right of mid (l = mid+1); otherwise it's mid or left (r = mid). Loop until l == r — that's the minimum.

    Complexity: O(log n) time, O(1) space
  kind: essay
  prompt: 'Solve "Find Minimum in Rotated Sorted Array" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Find Minimum in Rotated Sorted Array
---

**Pattern:** Binary Search · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array)

A sorted array of unique elements was rotated between 1 and n times. Return the **minimum** element in O(log n).

**Example 1:**

    Input: nums = [3,4,5,1,2]
    Output: 1

**Example 2:**

    Input: nums = [4,5,6,7,0,1,2]
    Output: 0

**Constraints:**

- `1 <= nums.length <= 5000`
- All values are unique

---

**Hints — try each one before reading on:**
1. The minimum is the rotation point.
2. Compare mid against the RIGHT end: greater means the pivot is right of mid.

**Target:** O(log n) time, O(1) space
