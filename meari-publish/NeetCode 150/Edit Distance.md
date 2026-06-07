---
created: "2026-06-07"
id: nc-edit-distance
source: imported:neetcode-150
study:
  answer: |-
    Levenshtein DP over suffixes (or prefixes): base rows are the remaining lengths; equal characters copy the diagonal, otherwise 1 + min of the three neighboring states. dp[0][0] is the distance.

    Complexity: O(m·n) time, O(min(m,n)) space
  kind: essay
  prompt: 'Solve "Edit Distance" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Edit Distance
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/edit-distance)

Return the minimum number of operations — **insert, delete, or replace** a character — to convert `word1` into `word2`.

**Example 1:**

    Input: word1 = "horse", word2 = "ros"
    Output: 3
    Explanation: horse → rorse → rose → ros.

**Example 2:**

    Input: word1 = "intention", word2 = "execution"
    Output: 5

**Constraints:**

- `0 <= word1.length, word2.length <= 500`

---

**Hints — try each one before reading on:**
1. Equal chars cost nothing — move diagonally.
2. Otherwise 1 + min(insert, delete, replace) = 1 + min(dp[i][j+1], dp[i+1][j], dp[i+1][j+1]).

**Target:** O(m·n) time, O(min(m,n)) space
