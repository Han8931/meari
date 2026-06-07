---
created: "2026-06-07"
id: nc-search-in-rotated-sorted-array
source: imported:neetcode-150
study:
  answer: |-
    Modified binary search: compare nums[l] ≤ nums[mid] to learn which half is sorted; if the target falls within the sorted half's bounds, recurse there, else into the other half.

    Complexity: O(log n) time, O(1) space
  kind: essay
  prompt: 'Solve "Search in Rotated Sorted Array" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Search in Rotated Sorted Array
---

**Pattern:** Binary Search · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/search-in-rotated-sorted-array)

A sorted array of **distinct** values was rotated at an unknown pivot. Given the rotated `nums` and a `target`, return its index or `-1`, in O(log n).

**Example 1:**

    Input: nums = [4,5,6,7,0,1,2], target = 0
    Output: 4

**Example 2:**

    Input: nums = [4,5,6,7,0,1,2], target = 3
    Output: -1

**Constraints:**

- `1 <= nums.length <= 5000`
- All values are unique

---

**Hints — try each one before reading on:**
1. At every mid, one half is still properly sorted — figure out which.
2. If the target lies inside that sorted half's range, go there; otherwise the other half.

**Target:** O(log n) time, O(1) space
