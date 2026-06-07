---
created: "2026-06-07"
id: nc-trapping-rain-water
source: imported:neetcode-150
study:
  answer: |-
    Two pointers with running left-max and right-max. Always step the side whose max is smaller: water there equals that max minus the bar height (the other side's max can't matter since it's larger). Accumulate as you go.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Trapping Rain Water" (Two Pointers): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Trapping Rain Water
---

**Pattern:** Two Pointers · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/trapping-rain-water)

Given `height` describing an elevation map with bars of width 1, compute how much water it traps after raining.

**Example 1:**

    Input: height = [0,1,0,2,1,0,1,3,2,1,2,1]
    Output: 6

**Example 2:**

    Input: height = [4,2,0,3,2,5]
    Output: 9

**Constraints:**

- `1 <= height.length <= 2 * 10^4`
- `0 <= height[i] <= 10^5`

---

**Hints — try each one before reading on:**
1. Water above a bar = min(max-to-left, max-to-right) − height.
2. Two pointers: the side with the smaller running max is fully determined — process it.

**Target:** O(n) time, O(1) space
