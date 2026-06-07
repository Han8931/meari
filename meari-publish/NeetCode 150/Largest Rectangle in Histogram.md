---
created: "2026-06-07"
id: nc-largest-rectangle-in-histogram
source: imported:neetcode-150
study:
  answer: |-
    Monotonic increasing stack of (start index, height). When a new bar is shorter than the top, pop and compute area = height × (i − start), remembering the popped start as the new bar's extended start. Flush the stack at the end against width n.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Largest Rectangle in Histogram" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Largest Rectangle in Histogram
---

**Pattern:** Stack · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/largest-rectangle-in-histogram)

Given `heights` describing a histogram with bars of width 1, return the area of the **largest rectangle** that fits inside it.

**Example 1:**

    Input: heights = [2,1,5,6,2,3]
    Output: 10
    Explanation: the 5 and 6 bars form a 5 x 2 rectangle.

**Example 2:**

    Input: heights = [2,4]
    Output: 4

**Constraints:**

- `1 <= heights.length <= 10^5`
- `0 <= heights[i] <= 10^4`

---

**Hints — try each one before reading on:**
1. A bar defines a rectangle extending while neighbors are at least as tall.
2. Monotonic increasing stack: a shorter bar finalizes the rectangles of taller bars on the stack.

**Target:** O(n) time, O(n) space
