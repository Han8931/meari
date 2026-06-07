---
created: "2026-06-07"
id: nc-best-time-to-buy-and-sell-stock
source: imported:neetcode-150
study:
  answer: |-
    Scan once keeping the lowest price so far; at each price, profit = price − minSoFar, keep the best. No window bookkeeping actually needed.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Best Time to Buy and Sell Stock" (Sliding Window): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Best Time to Buy and Sell Stock
---

**Pattern:** Sliding Window · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/best-time-to-buy-and-sell-stock)

Given `prices` where `prices[i]` is a stock's price on day i, choose one day to buy and a **later** day to sell. Return the maximum profit, or 0 if no profit is possible.

**Example 1:**

    Input: prices = [7,1,5,3,6,4]
    Output: 5
    Explanation: buy at 1 (day 2), sell at 6 (day 5).

**Example 2:**

    Input: prices = [7,6,4,3,1]
    Output: 0

**Constraints:**

- `1 <= prices.length <= 10^5`
- `0 <= prices[i] <= 10^4`

---

**Hints — try each one before reading on:**
1. For each day, the best buy is the minimum price seen so far.
2. One pass: track min price and best profit together.

**Target:** O(n) time, O(1) space
