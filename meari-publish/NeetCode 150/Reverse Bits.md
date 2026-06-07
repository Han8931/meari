---
created: "2026-06-07"
id: nc-reverse-bits
source: imported:neetcode-150
study:
  answer: |-
    Build the reversal bit by bit: shift the result left and OR in n's lowest bit, shifting n right, 32 times.

    Complexity: O(32) ≈ O(1) time, O(1) space
  kind: essay
  prompt: 'Solve "Reverse Bits" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Reverse Bits
---

**Pattern:** Bit Manipulation · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/reverse-bits)

Reverse the bits of a 32-bit unsigned integer.

**Example 1:**

    Input: n = 0b00000010100101000001111010011100
    Output: 964176192 (0b00111001011110000010100101000000)

**Constraints:**

- The input is a 32-bit binary string's value

---

**Hints — try each one before reading on:**
1. Peel bits from one end, push onto the other.
2. 32 iterations: result = (result << 1) | (n & 1); n >>= 1.

**Target:** O(32) ≈ O(1) time, O(1) space
