---
created: "2026-06-07"
id: nc-best-time-to-buy-and-sell-stock-with-cooldown
source: imported:neetcode-150
study:
  answer: |-
    State-machine DP with three rolling values — hold = max(hold, rest − price), sold = hold + price, rest = max(rest, prevSold) — iterated over prices; answer is max(sold, rest).

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Best Time to Buy and Sell Stock with Cooldown" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Best Time to Buy and Sell Stock with Cooldown
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-cooldown)

Trade a stock as many times as you like, but after **selling** you must cool down one day before buying again. Return the maximum profit.

**Example 1:**

    Input: prices = [1,2,3,0,2]
    Output: 3
    Explanation: buy, sell, cooldown, buy, sell.

**Example 2:**

    Input: prices = [1]
    Output: 0

**Constraints:**

- `1 <= prices.length <= 5000`
- `0 <= prices[i] <= 1000`

---

**Hints — try each one before reading on:**
1. Model states: holding, just sold (cooldown), free to buy.
2. Each day's states derive from yesterday's three values.

**Target:** O(n) time, O(1) space
