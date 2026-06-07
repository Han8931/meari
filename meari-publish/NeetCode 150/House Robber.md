---
created: "2026-06-07"
id: nc-house-robber
source: imported:neetcode-150
study:
  answer: |-
    Iterate with rob1, rob2: new = max(rob2, rob1 + nums[i]); shift. rob2 at the end is the maximum without ever taking adjacent houses.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "House Robber" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: House Robber
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/house-robber)

Houses on a street hold money `nums[i]`, but robbing two **adjacent** houses trips the alarm. Return the maximum you can rob.

**Example 1:**

    Input: nums = [1,2,3,1]
    Output: 4
    Explanation: rob houses 0 and 2.

**Example 2:**

    Input: nums = [2,7,9,3,1]
    Output: 12
    Explanation: rob 2 + 9 + 1.

**Constraints:**

- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 400`

---

**Hints — try each one before reading on:**
1. At each house: rob it (skip the previous) or don't.
2. dp[i] = max(dp[i−1], dp[i−2] + nums[i]) — two rolling values.

**Target:** O(n) time, O(1) space
