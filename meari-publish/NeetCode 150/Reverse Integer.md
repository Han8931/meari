---
created: "2026-06-07"
id: nc-reverse-integer
source: imported:neetcode-150
study:
  answer: |-
    Pop-and-push digits while guarding the 32-bit range before each multiply-add (compare to (2³¹−1)//10 and the trailing digit); handle the sign separately since Python's % differs on negatives. Return 0 when a push would overflow.

    Complexity: O(log n) time, O(1) space
  kind: essay
  prompt: 'Solve "Reverse Integer" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Reverse Integer
---

**Pattern:** Bit Manipulation · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/reverse-integer)

Reverse the digits of a signed 32-bit integer `x`. If the result overflows the 32-bit range, return 0 (assume you cannot store 64-bit values).

**Example 1:**

    Input: x = 123
    Output: 321

**Example 2:**

    Input: x = -123
    Output: -321

**Example 3:**

    Input: x = 120
    Output: 21

**Constraints:**

- `-2^31 <= x <= 2^31 - 1`

---

**Hints — try each one before reading on:**
1. Pop digits with % 10 and / 10, push with result · 10 + digit.
2. Check against INT_MAX // 10 BEFORE each push — that's where overflow happens.

**Target:** O(log n) time, O(1) space
