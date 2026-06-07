---
created: "2026-06-07"
id: nc-interleaving-string
source: imported:neetcode-150
study:
  answer: |-
    Length check first; then dp[i][j] = (dp[i−1][j] and s1[i−1] == s3[i+j−1]) or (dp[i][j−1] and s2[j−1] == s3[i+j−1]), filled with one rolling row. dp[m][n] answers it.

    Complexity: O(m·n) time, O(n) space
  kind: essay
  prompt: 'Solve "Interleaving String" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Interleaving String
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/interleaving-string)

Return `true` if `s3` is an **interleaving** of `s1` and `s2`: it merges both strings completely while preserving each one's character order.

**Example 1:**

    Input: s1 = "aabcc", s2 = "dbbca", s3 = "aadbbcbcac"
    Output: true

**Example 2:**

    Input: s1 = "aabcc", s2 = "dbbca", s3 = "aadbbbaccc"
    Output: false

**Constraints:**

- `0 <= s1.length, s2.length <= 100`
- `s3.length == s1.length + s2.length` (else false)

---

**Hints — try each one before reading on:**
1. State = (i chars of s1 used, j chars of s2 used); s3's next char is forced to index i+j.
2. dp[i][j] true if either source can supply s3[i+j].

**Target:** O(m·n) time, O(n) space
