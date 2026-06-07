---
created: "2026-06-07"
id: nc-coin-change
source: imported:neetcode-150
study:
  answer: |-
    Bottom-up: dp array of size amount+1 initialized to infinity, dp[0]=0; for each amount a, dp[a] = min over coins of dp[a−c] + 1. Return dp[amount] or −1 when it stayed infinite.

    Complexity: O(amount · coins) time, O(amount) space
  kind: essay
  prompt: 'Solve "Coin Change" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Coin Change
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/coin-change)

Given coin denominations and an `amount`, return the **fewest coins** needed to make the amount exactly (unlimited supply), or `-1`.

**Example 1:**

    Input: coins = [1,2,5], amount = 11
    Output: 3
    Explanation: 5 + 5 + 1.

**Example 2:**

    Input: coins = [2], amount = 3
    Output: -1

**Constraints:**

- `1 <= coins.length <= 12`
- `0 <= amount <= 10^4`

---

**Hints — try each one before reading on:**
1. Greedy fails (try 1, 3, 4 for 6) — it's unbounded-knapsack DP.
2. dp[a] = 1 + min(dp[a − coin]) over all coins.

**Target:** O(amount · coins) time, O(amount) space
