---
created: "2026-06-07"
id: nc-daily-temperatures
source: imported:neetcode-150
study:
  answer: |-
    Monotonic decreasing stack of indices: when today's temperature beats the stack top, pop it and record today − popped as its wait; push today. Unresolved indices stay 0.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Daily Temperatures" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Daily Temperatures
---

**Pattern:** Stack · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/daily-temperatures)

Given `temperatures`, return an array `answer` where `answer[i]` is the number of days you must wait after day i for a warmer temperature, or 0 if it never comes.

**Example 1:**

    Input: temperatures = [73,74,75,71,69,72,76,73]
    Output: [1,1,4,2,1,1,0,0]

**Constraints:**

- `1 <= temperatures.length <= 10^5`
- `30 <= temperatures[i] <= 100`

---

**Hints — try each one before reading on:**
1. "Next greater element" — a monotonic stack problem.
2. Keep indices of a decreasing stack; a warmer day resolves everything cooler on top.

**Target:** O(n) time, O(n) space
