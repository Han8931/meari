---
created: "2026-06-07"
id: nc-longest-common-subsequence
source: imported:neetcode-150
study:
  answer: |-
    dp[i][j] = LCS of suffixes: equal characters give 1 + dp[i+1][j+1], otherwise max(dp[i+1][j], dp[i][j+1]); fill the table bottom-up (two rolling rows suffice). dp[0][0] is the answer.

    Complexity: O(m·n) time, O(min(m,n)) space
  kind: essay
  prompt: 'Solve "Longest Common Subsequence" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Common Subsequence
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/longest-common-subsequence)

Given `text1` and `text2`, return the length of their longest **common subsequence** (characters in order, not necessarily contiguous), or 0.

**Example 1:**

    Input: text1 = "abcde", text2 = "ace"
    Output: 3
    Explanation: "ace".

**Example 2:**

    Input: text1 = "abc", text2 = "def"
    Output: 0

**Constraints:**

- `1 <= text1.length, text2.length <= 1000`

---

**Hints — try each one before reading on:**
1. Match → 1 + diagonal; mismatch → max(skip either character).
2. Classic 2-D table over the two string prefixes.

**Target:** O(m·n) time, O(min(m,n)) space
