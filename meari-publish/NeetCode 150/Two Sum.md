---
created: "2026-06-07"
id: nc-two-sum
source: imported:neetcode-150
study:
  answer: |-
    One pass with a hash map of value → index: for each x at i, if target − x is already in the map return its index and i; otherwise store x. The single pass guarantees you never match an element with itself.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Two Sum" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Two Sum
---

**Pattern:** Arrays & Hashing · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/two-sum)

Given an array of integers `nums` and an integer `target`, return the **indices** of the two numbers that add up to `target`. Exactly one solution exists, and you may not use the same element twice.

**Example 1:**

    Input: nums = [2,7,11,15], target = 9
    Output: [0,1]
    Explanation: nums[0] + nums[1] == 9.

**Example 2:**

    Input: nums = [3,2,4], target = 6
    Output: [1,2]

**Constraints:**

- `2 <= nums.length <= 10^4`
- `-10^9 <= nums[i], target <= 10^9`
- Exactly one valid answer exists

---

**Hints — try each one before reading on:**
1. For each number, what exactly are you looking for in the rest of the array?
2. Store value → index as you scan; look up target − current before inserting.

**Target:** O(n) time, O(n) space
