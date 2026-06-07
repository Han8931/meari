---
created: "2026-06-07"
id: nc-merge-intervals
source: imported:neetcode-150
study:
  answer: |-
    Sort by start and sweep: if the next interval starts at or before the current merged end, extend the end; otherwise emit and start fresh.

    Complexity: O(n log n) time, O(n) space
  kind: essay
  prompt: 'Solve "Merge Intervals" (Intervals): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Merge Intervals
---

**Pattern:** Intervals · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/merge-intervals)

Given an array of `intervals`, merge all overlapping ones and return the non-overlapping result.

**Example 1:**

    Input: intervals = [[1,3],[2,6],[8,10],[15,18]]
    Output: [[1,6],[8,10],[15,18]]

**Example 2:**

    Input: intervals = [[1,4],[4,5]]
    Output: [[1,5]]

**Constraints:**

- `1 <= intervals.length <= 10^4`
- `0 <= start <= end <= 10^4`

---

**Hints — try each one before reading on:**
1. Sort by start; overlaps become adjacent.
2. Extend the last merged interval when the next start is within it.

**Target:** O(n log n) time, O(n) space
