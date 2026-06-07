---
created: "2026-06-07"
id: nc-gas-station
source: imported:neetcode-150
study:
  answer: |-
    One pass: if the running tank dips below zero at station i, no start ≤ i works — set the candidate to i+1 and reset the tank. With total gas ≥ total cost, the surviving candidate is the answer.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Gas Station" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Gas Station
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/gas-station)

Around a circular route, station i provides `gas[i]` and driving to the next costs `cost[i]`. Return the starting station allowing a full clockwise circuit, or `-1` (the answer is unique).

**Example 1:**

    Input: gas = [1,2,3,4,5], cost = [3,4,5,1,2]
    Output: 3

**Example 2:**

    Input: gas = [2,3,4], cost = [3,4,3]
    Output: -1

**Constraints:**

- `1 <= n <= 10^5`
- If an answer exists it is unique

---

**Hints — try each one before reading on:**
1. If total gas < total cost, no answer; otherwise exactly one exists.
2. A negative running tank invalidates every start in that stretch — restart from the next station.

**Target:** O(n) time, O(1) space
