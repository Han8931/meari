---
created: "2026-06-07"
id: nc-maximum-subarray
source: imported:neetcode-150
study:
  answer: |-
    Kadane's algorithm: keep cur = max(x, cur + x) (equivalently drop the prefix when it goes negative) and best = max(best, cur) in one pass.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Maximum Subarray" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Maximum Subarray
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/maximum-subarray)

Return the largest **sum** over all contiguous subarrays of `nums`.

**Example 1:**

    Input: nums = [-2,1,-3,4,-1,2,1,-5,4]
    Output: 6
    Explanation: [4,-1,2,1].

**Example 2:**

    Input: nums = [5,4,-1,7,8]
    Output: 23

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`

---

**Hints — try each one before reading on:**
1. A negative running prefix can never help what follows.
2. Kadane: reset the running sum when it drops below zero.

**Target:** O(n) time, O(1) space
