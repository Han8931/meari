---
created: "2026-06-07"
id: nc-regular-expression-matching
source: imported:neetcode-150
study:
  answer: |-
    Memoized DP over (text index, pattern index): firstMatch = chars equal or pattern '.'; if p[j+1] == '*', result = dp(i, j+2) or (firstMatch and dp(i+1, j)); else firstMatch and dp(i+1, j+1). Base: both exhausted.

    Complexity: O(m·n) time, O(m·n) space
  kind: essay
  prompt: 'Solve "Regular Expression Matching" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Regular Expression Matching
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/regular-expression-matching)

Implement **full-string** regex matching where `'.'` matches any single character and `'*'` matches zero or more of the preceding element.

**Example 1:**

    Input: s = "aa", p = "a"
    Output: false

**Example 2:**

    Input: s = "aa", p = "a*"
    Output: true

**Example 3:**

    Input: s = "ab", p = ".*"
    Output: true

**Constraints:**

- `1 <= s.length <= 20`, `1 <= p.length <= 20`
- Every `'*'` follows a valid element

---

**Hints — try each one before reading on:**
1. Decide at (i, j): does the pattern's NEXT char form a star pair?
2. Star: skip the pair (zero uses) OR consume one matching char and stay on the pair.

**Target:** O(m·n) time, O(m·n) space
