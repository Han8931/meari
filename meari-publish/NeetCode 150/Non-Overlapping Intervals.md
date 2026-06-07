---
created: "2026-06-07"
id: nc-non-overlapping-intervals
source: imported:neetcode-150
study:
  answer: |-
    Sort by end time; sweep keeping intervals whose start ≥ the last kept end, counting the rest as removed. Earliest-deadline-first maximizes the kept set, so removals are minimal.

    Complexity: O(n log n) time, O(1) extra space
  kind: essay
  prompt: 'Solve "Non-Overlapping Intervals" (Intervals): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Non-Overlapping Intervals
---

**Pattern:** Intervals · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/non-overlapping-intervals)

Return the **minimum number of intervals to remove** so the rest don't overlap (touching endpoints don't overlap).

**Example 1:**

    Input: intervals = [[1,2],[2,3],[3,4],[1,3]]
    Output: 1
    Explanation: remove [1,3].

**Example 2:**

    Input: intervals = [[1,2],[1,2],[1,2]]
    Output: 2

**Constraints:**

- `1 <= intervals.length <= 10^5`
- `-5 * 10^4 <= start < end <= 5 * 10^4`

---

**Hints — try each one before reading on:**
1. Equivalent to KEEPING the most non-overlapping intervals.
2. Sort by END; greedily keep each interval that starts at or after the last kept end.

**Target:** O(n log n) time, O(1) extra space
