---
created: "2026-06-07"
id: nc-climbing-stairs
source: imported:neetcode-150
study:
  answer: |-
    Bottom-up with two variables: iterate n−1 times folding (a, b) → (b, a+b) from base cases 1, 1. It's the Fibonacci recurrence; no memo table needed.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Climbing Stairs" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Climbing Stairs
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/climbing-stairs)

You climb a staircase of `n` steps, taking **1 or 2** steps at a time. In how many distinct ways can you reach the top?

**Example 1:**

    Input: n = 2
    Output: 2
    Explanation: 1+1 or 2.

**Example 2:**

    Input: n = 3
    Output: 3
    Explanation: 1+1+1, 1+2, 2+1.

**Constraints:**

- `1 <= n <= 45`

---

**Hints — try each one before reading on:**
1. ways(n) = ways(n−1) + ways(n−2) — Fibonacci in disguise.
2. Two rolling variables beat an array.

**Target:** O(n) time, O(1) space
