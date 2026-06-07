---
created: "2026-06-07"
id: nc-container-with-most-water
source: imported:neetcode-150
study:
  answer: |-
    Two pointers at the ends; record the area, then move the shorter line inward — moving the taller one can only shrink both width and the limiting height, so it's never useful.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Container With Most Water" (Two Pointers): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Container With Most Water
---

**Pattern:** Two Pointers · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/container-with-most-water)

Given `height`, where `height[i]` is the height of a vertical line at x = i, choose two lines that, with the x-axis, contain the most water. Return that maximum area.

**Example 1:**

    Input: height = [1,8,6,2,5,4,8,3,7]
    Output: 49
    Explanation: lines at i = 1 and i = 8: min(8,7) * (8-1) = 49.

**Example 2:**

    Input: height = [1,1]
    Output: 1

**Constraints:**

- `2 <= height.length <= 10^5`
- `0 <= height[i] <= 10^4`

---

**Hints — try each one before reading on:**
1. Area = min(height[l], height[r]) × (r − l).
2. Which pointer is ever worth moving — the taller or the shorter?

**Target:** O(n) time, O(1) space
