---
created: "2026-06-07"
id: nc-multiply-strings
source: imported:neetcode-150
study:
  answer: |-
    Grade-school multiplication into a m+n result array: add digit products at index i+j+1 (reversed indexing), propagate carries, strip leading zeros, and join. Handle a zero operand early.

    Complexity: O(m·n) time, O(m+n) space
  kind: essay
  prompt: 'Solve "Multiply Strings" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Multiply Strings
---

**Pattern:** Math & Geometry · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/multiply-strings)

Multiply two non-negative integers given as strings, **without** big-integer libraries or direct conversion.

**Example 1:**

    Input: num1 = "2", num2 = "3"
    Output: "6"

**Example 2:**

    Input: num1 = "123", num2 = "456"
    Output: "56088"

**Constraints:**

- `1 <= num1.length, num2.length <= 200`
- No leading zeros except the number 0

---

**Hints — try each one before reading on:**
1. Digit i × digit j lands at result positions i+j and i+j+1 (from the right).
2. Accumulate into an int array of size m+n, then carry once.

**Target:** O(m·n) time, O(m+n) space
