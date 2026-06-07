---
created: "2026-06-07"
id: nc-insert-interval
source: imported:neetcode-150
study:
  answer: |-
    Append intervals ending before the new one; absorb every overlapping interval into it via min(start)/max(end); append it, then the rest. The input's sortedness makes one pass enough.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Insert Interval" (Intervals): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Insert Interval
---

**Pattern:** Intervals · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/insert-interval)

Given non-overlapping `intervals` sorted by start and a `newInterval`, insert it and merge so the result stays sorted and non-overlapping.

**Example 1:**

    Input: intervals = [[1,3],[6,9]], newInterval = [2,5]
    Output: [[1,5],[6,9]]

**Example 2:**

    Input: intervals = [[1,2],[3,5],[6,7],[8,10],[12,16]], newInterval = [4,8]
    Output: [[1,2],[3,10],[12,16]]

**Constraints:**

- `0 <= intervals.length <= 10^4`
- Starts/ends up to `10^5`, sorted by start

---

**Hints — try each one before reading on:**
1. Three phases: intervals entirely before, the overlapping middle, entirely after.
2. Merge the middle into the new interval by min/max.

**Target:** O(n) time, O(n) space
