---
created: "2026-06-07"
id: nc-min-cost-climbing-stairs
source: imported:neetcode-150
study:
  answer: |-
    Walk backward (or forward with two variables): dp[i] = cost[i] + min(dp[i+1], dp[i+2]) with dp = 0 beyond the array; answer is min(dp[0], dp[1]).

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Min Cost Climbing Stairs" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Min Cost Climbing Stairs
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/min-cost-climbing-stairs)

`cost[i]` is the price of stepping on stair i; after paying you may climb 1 or 2 stairs. Starting from index 0 **or** 1, return the minimum cost to reach the top (just past the last stair).

**Example 1:**

    Input: cost = [10,15,20]
    Output: 15
    Explanation: start at 15, jump two to the top.

**Example 2:**

    Input: cost = [1,100,1,1,1,100,1,1,100,1]
    Output: 6

**Constraints:**

- `2 <= cost.length <= 1000`
- `0 <= cost[i] <= 999`

---

**Hints — try each one before reading on:**
1. minCost(i) = cost[i] + min(minCost(i+1), minCost(i+2)).
2. You may start at step 0 or 1; the "top" is past the last step.

**Target:** O(n) time, O(1) space
