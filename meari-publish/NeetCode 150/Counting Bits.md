---
created: "2026-06-07"
id: nc-counting-bits
source: imported:neetcode-150
study:
  answer: |-
    DP over the answers already computed: dp[i] = dp[i >> 1] + (i & 1) (or dp[i & (i−1)] + 1), filled 0..n.

    Complexity: O(n) time, O(n) output space
  kind: essay
  prompt: 'Solve "Counting Bits" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Counting Bits
---

**Pattern:** Bit Manipulation · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/counting-bits)

Given `n`, return an array `ans` of length `n + 1` where `ans[i]` is the number of 1s in the binary representation of i — ideally in one O(n) pass.

**Example 1:**

    Input: n = 2
    Output: [0,1,1]

**Example 2:**

    Input: n = 5
    Output: [0,1,1,2,1,2]

**Constraints:**

- `0 <= n <= 10^5`

---

**Hints — try each one before reading on:**
1. i and i >> 1 differ by just the lowest bit.
2. dp[i] = dp[i >> 1] + (i & 1).

**Target:** O(n) time, O(n) output space
