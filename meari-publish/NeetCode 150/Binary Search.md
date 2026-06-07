---
created: "2026-06-07"
id: nc-binary-search
source: imported:neetcode-150
study:
  answer: |-
    Standard binary search: l, r = 0, n−1; mid = (l+r)//2; go left or right by comparison until found or l > r.

    Complexity: O(log n) time, O(1) space
  kind: essay
  prompt: 'Solve "Binary Search" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Binary Search
---

**Pattern:** Binary Search · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/binary-search)

Given a **sorted** integer array `nums` and a `target`, return the index of `target`, or `-1` if absent — in O(log n).

**Example 1:**

    Input: nums = [-1,0,3,5,9,12], target = 9
    Output: 4

**Example 2:**

    Input: nums = [-1,0,3,5,9,12], target = 2
    Output: -1

**Constraints:**

- `1 <= nums.length <= 10^4`
- All elements are unique and sorted ascending

---

**Hints — try each one before reading on:**
1. Halve the search space by comparing the middle element.
2. Watch the loop condition (l ≤ r) and mid updates to avoid an infinite loop.

**Target:** O(log n) time, O(1) space
