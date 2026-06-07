---
created: "2026-06-07"
id: nc-plus-one
source: imported:neetcode-150
study:
  answer: |-
    From the last digit: a 9 becomes 0 and the carry continues; anything else increments and returns immediately. If every digit was 9, prepend a 1.

    Complexity: O(n) time, O(1) extra space
  kind: essay
  prompt: 'Solve "Plus One" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Plus One
---

**Pattern:** Math & Geometry · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/plus-one)

A large integer is given as a digit array (most significant first, no leading zeros). Add one and return the digits.

**Example 1:**

    Input: digits = [1,2,3]
    Output: [1,2,4]

**Example 2:**

    Input: digits = [9,9,9]
    Output: [1,0,0,0]

**Constraints:**

- `1 <= digits.length <= 100`
- `0 <= digits[i] <= 9`

---

**Hints — try each one before reading on:**
1. Only 9s propagate a carry.
2. Walk from the right turning 9→0; a non-9 absorbs the carry and you're done.

**Target:** O(n) time, O(1) extra space
