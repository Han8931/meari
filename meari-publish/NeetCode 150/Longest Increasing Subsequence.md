---
created: "2026-06-07"
id: nc-longest-increasing-subsequence
source: imported:neetcode-150
study:
  answer: |-
    Patience method: maintain tails[] where tails[k] is the smallest tail of any increasing subsequence of length k+1; for each x, bisect_left its position and replace (or append at the end). The final length of tails is the answer.

    Complexity: O(n log n) time, O(n) space
  kind: essay
  prompt: 'Solve "Longest Increasing Subsequence" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Increasing Subsequence
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/longest-increasing-subsequence)

Return the length of the longest **strictly increasing subsequence** of `nums` (elements keep their order but need not be contiguous).

**Example 1:**

    Input: nums = [10,9,2,5,3,7,101,18]
    Output: 4
    Explanation: [2,3,7,101].

**Example 2:**

    Input: nums = [0,1,0,3,2,3]
    Output: 4

**Constraints:**

- `1 <= nums.length <= 2500`
- `-10^4 <= nums[i] <= 10^4`

---

**Hints — try each one before reading on:**
1. O(n²) DP: best ending at i extends the best j < i with smaller value.
2. O(n log n): keep "tails" — the smallest possible tail of an LIS of each length — and binary-search where each number lands.

**Target:** O(n log n) time, O(n) space
