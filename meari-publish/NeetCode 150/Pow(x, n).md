---
created: "2026-06-07"
id: nc-pow-x-n
source: imported:neetcode-150
study:
  answer: |-
    Fast exponentiation: loop while n > 0 multiplying the result by x when the low bit is set, then squaring x and halving n; for negative exponents compute on |n| and take the reciprocal.

    Complexity: O(log n) time, O(1) space
  kind: essay
  prompt: 'Solve "Pow(x, n)" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Pow(x, n)
---

**Pattern:** Math & Geometry · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/powx-n)

Implement `pow(x, n)` computing `x` raised to the integer power `n` (which may be negative).

**Example 1:**

    Input: x = 2.0, n = 10
    Output: 1024.0

**Example 2:**

    Input: x = 2.0, n = -2
    Output: 0.25

**Constraints:**

- `-100.0 < x < 100.0`
- `-2^31 <= n <= 2^31 - 1`

---

**Hints — try each one before reading on:**
1. x²ᵏ = (x²)ᵏ — halve the exponent each step.
2. Negative n: invert x and use |n|.

**Target:** O(log n) time, O(1) space
