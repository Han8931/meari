---
created: "2026-06-07"
id: nc-sum-of-two-integers
source: imported:neetcode-150
study:
  answer: |-
    Loop: sum = a ^ b, carry = (a & b) << 1, repeat with (sum, carry) until carry is 0. Python's unbounded ints need a 0xFFFFFFFF mask each round and a final two's-complement fix for negative results.

    Complexity: O(32) ≈ O(1) time, O(1) space
  kind: essay
  prompt: 'Solve "Sum of Two Integers" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Sum of Two Integers
---

**Pattern:** Bit Manipulation · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/sum-of-two-integers)

Return the sum of two integers `a` and `b` **without** using `+` or `-`.

**Example 1:**

    Input: a = 1, b = 2
    Output: 3

**Example 2:**

    Input: a = 2, b = 3
    Output: 5

**Constraints:**

- `-1000 <= a, b <= 1000`

---

**Hints — try each one before reading on:**
1. XOR adds without carries; AND << 1 is exactly the carries.
2. Repeat until the carry dies; in Python, mask to 32 bits and fix negatives.

**Target:** O(32) ≈ O(1) time, O(1) space
