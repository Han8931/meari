---
created: "2026-06-07"
id: nc-house-robber-ii
source: imported:neetcode-150
study:
  answer: |-
    Answer = max(rob(nums[1:]), rob(nums[:-1])) using the linear House Robber helper, with the single-house edge case handled directly.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "House Robber II" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: House Robber II
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/house-robber-ii)

Same as House Robber, but the houses form a **circle** — the first and last are adjacent. Return the maximum loot.

**Example 1:**

    Input: nums = [2,3,2]
    Output: 3
    Explanation: houses 0 and 2 are adjacent on the circle, so take 3.

**Example 2:**

    Input: nums = [1,2,3,1]
    Output: 4

**Constraints:**

- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 1000`

---

**Hints — try each one before reading on:**
1. The circle only forbids taking BOTH ends.
2. Run the linear robber twice: without the first house, and without the last.

**Target:** O(n) time, O(1) space
