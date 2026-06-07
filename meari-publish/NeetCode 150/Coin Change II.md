---
created: "2026-06-07"
id: nc-coin-change-ii
source: imported:neetcode-150
study:
  answer: |-
    Unbounded-knapsack counting: dp[0] = 1; for each coin, for a from coin to amount, dp[a] += dp[a − coin]. The coin-major loop order counts combinations, not permutations.

    Complexity: O(amount · coins) time, O(amount) space
  kind: essay
  prompt: 'Solve "Coin Change II" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Coin Change II
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/coin-change-ii)

Given coin denominations and an `amount`, return the number of **combinations** that make the amount (unlimited supply; order doesn't matter).

**Example 1:**

    Input: amount = 5, coins = [1,2,5]
    Output: 4
    Explanation: 5, 2+2+1, 2+1+1+1, 1+1+1+1+1.

**Example 2:**

    Input: amount = 3, coins = [2]
    Output: 0

**Constraints:**

- `1 <= coins.length <= 300`
- `0 <= amount <= 5000`
- Coins are distinct

---

**Hints — try each one before reading on:**
1. Loop coins on the OUTSIDE so orderings aren't double-counted.
2. dp[a] += dp[a − coin] for each coin in turn.

**Target:** O(amount · coins) time, O(amount) space
