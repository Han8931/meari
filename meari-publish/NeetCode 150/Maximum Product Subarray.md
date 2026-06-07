---
created: "2026-06-07"
id: nc-maximum-product-subarray
source: imported:neetcode-150
study:
  answer: |-
    One pass keeping curMax and curMin ending here: each new x updates them from (x, x·curMax, x·curMin) — swap roles when x < 0 — and the global best tracks curMax. A zero resets both to x.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Maximum Product Subarray" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Maximum Product Subarray
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/maximum-product-subarray)

Return the largest **product** over all contiguous subarrays of `nums`.

**Example 1:**

    Input: nums = [2,3,-2,4]
    Output: 6
    Explanation: [2,3].

**Example 2:**

    Input: nums = [-2,0,-1]
    Output: 0

**Constraints:**

- `1 <= nums.length <= 2 * 10^4`
- `-10 <= nums[i] <= 10`
- Products fit in 32 bits

---

**Hints — try each one before reading on:**
1. A big negative becomes the biggest positive after one more negative.
2. Track BOTH the max and min product ending at each index.

**Target:** O(n) time, O(1) space
