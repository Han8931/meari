---
created: "2026-06-07"
id: nc-meeting-rooms-ii
source: imported:neetcode-150
study:
  answer: |-
    Sort starts and ends separately; two-pointer sweep: a start before the next end takes a room (count++), otherwise an end frees one. The peak count is the answer. (Min-heap of end times: push each meeting, popping ends ≤ its start; peak heap size.)

    Complexity: O(n log n) time, O(n) space
  kind: essay
  prompt: 'Solve "Meeting Rooms II" (Intervals): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Meeting Rooms II
---

**Pattern:** Intervals · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/meeting-rooms-ii)

Given meeting time `intervals`, return the **minimum number of conference rooms** required.

**Example 1:**

    Input: intervals = [[0,30],[5,10],[15,20]]
    Output: 2

**Example 2:**

    Input: intervals = [[7,10],[2,4]]
    Output: 1

**Constraints:**

- `1 <= intervals.length <= 10^4`

---

**Hints — try each one before reading on:**
1. The answer is the maximum number of meetings alive at one instant.
2. Sweep sorted start times against sorted end times (or a min-heap of end times).

**Target:** O(n log n) time, O(n) space
