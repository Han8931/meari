---
created: "2026-06-07"
id: nc-median-of-two-sorted-arrays
source: imported:neetcode-150
study:
  answer: |-
    Binary search the partition of the shorter array: cut A at i and B at half − i; valid when Aleft ≤ Bright and Bleft ≤ Aright. Then the median is max(lefts) for odd totals, or the average of max(lefts) and min(rights) for even. Move the cut by comparing the offending border elements.

    Complexity: O(log min(m, n)) time, O(1) space
  kind: essay
  prompt: 'Solve "Median of Two Sorted Arrays" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Median of Two Sorted Arrays
---

**Pattern:** Binary Search · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/median-of-two-sorted-arrays)

Given two sorted arrays `nums1` and `nums2`, return the **median** of their merged order — in O(log (m+n)).

**Example 1:**

    Input: nums1 = [1,3], nums2 = [2]
    Output: 2.0

**Example 2:**

    Input: nums1 = [1,2], nums2 = [3,4]
    Output: 2.5

**Constraints:**

- `0 <= m, n <= 1000`
- `-10^6 <= values <= 10^6`

---

**Hints — try each one before reading on:**
1. Binary search the smaller array for a partition; the other array's cut follows from half = (m+n)//2.
2. A partition is correct when every left element ≤ every right element across both arrays.

**Target:** O(log min(m, n)) time, O(1) space
