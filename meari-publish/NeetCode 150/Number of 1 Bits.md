---
created: "2026-06-07"
id: nc-number-of-1-bits
source: imported:neetcode-150
study:
  answer: |-
    Loop n &= n − 1, counting iterations — each removes exactly one set bit, so it runs once per 1-bit (Brian Kernighan's trick). bin(n).count('1') is the Python shortcut.

    Complexity: O(set bits) time, O(1) space
  kind: essay
  prompt: 'Solve "Number of 1 Bits" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Number of 1 Bits
---

**Pattern:** Bit Manipulation · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/number-of-1-bits)

Return the number of **set bits** (1s) in the binary representation of an unsigned 32-bit integer.

**Example 1:**

    Input: n = 11 (binary 1011)
    Output: 3

**Example 2:**

    Input: n = 128 (binary 10000000)
    Output: 1

**Constraints:**

- `0 <= n < 2^32`

---

**Hints — try each one before reading on:**
1. n & (n−1) clears the lowest set bit.
2. Count how many clears reach zero.

**Target:** O(set bits) time, O(1) space
