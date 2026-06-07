---
created: "2026-06-07"
id: nc-meeting-rooms
source: imported:neetcode-150
study:
  answer: |-
    Sort by start time and verify each meeting starts at or after the previous one ends; any violation means overlap.

    Complexity: O(n log n) time, O(1) space
  kind: essay
  prompt: 'Solve "Meeting Rooms" (Intervals): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Meeting Rooms
---

**Pattern:** Intervals · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/meeting-rooms)

Given meeting time `intervals`, return `true` if one person could attend **all** of them (no overlaps).

**Example 1:**

    Input: intervals = [[0,30],[5,10],[15,20]]
    Output: false

**Example 2:**

    Input: intervals = [[7,10],[2,4]]
    Output: true

**Constraints:**

- `0 <= intervals.length <= 10^4`

---

**Hints — try each one before reading on:**
1. Conflicts are overlaps — sort and check neighbors.
2. After sorting by start, only adjacent pairs can reveal an overlap.

**Target:** O(n log n) time, O(1) space
