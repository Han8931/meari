---
created: "2026-06-07"
id: nc-decode-ways
source: imported:neetcode-150
study:
  answer: |-
    DP from the end with two rolling values: ways(i) = ways(i+1) when s[i] != '0', plus ways(i+2) when s[i:i+2] is 10–26; a leading zero contributes 0. Answer is ways(0).

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Decode Ways" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Decode Ways
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/decode-ways)

Letters map to numbers `'A' → 1 … 'Z' → 26`. Given a digit string `s`, return how many ways it can be **decoded** (e.g. "12" is "AB" or "L").

**Example 1:**

    Input: s = "12"
    Output: 2

**Example 2:**

    Input: s = "226"
    Output: 3
    Explanation: 2|26, 22|6, 2|2|6.

**Example 3:**

    Input: s = "06"
    Output: 0

**Constraints:**

- `1 <= s.length <= 100`
- `s` contains only digits (possibly leading zeros)

---

**Hints — try each one before reading on:**
1. At each position: take one digit (if not '0'), or two (if 10–26).
2. dp[i] = (s[i]!='0') · dp[i+1] + (valid pair) · dp[i+2].

**Target:** O(n) time, O(1) space
