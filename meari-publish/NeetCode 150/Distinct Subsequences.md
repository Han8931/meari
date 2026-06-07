---
created: "2026-06-07"
id: nc-distinct-subsequences
source: imported:neetcode-150
study:
  answer: |-
    Suffix DP: dp[i][j] counts ways to build t[j:] from s[i:], with dp[i][len(t)] = 1; recurrence adds the skip-s branch and, when characters match, the take-both branch. One rolling row over j works.

    Complexity: O(m·n) time, O(n) space
  kind: essay
  prompt: 'Solve "Distinct Subsequences" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Distinct Subsequences
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/distinct-subsequences)

Return the number of **distinct subsequences** of `s` that equal `t`.

**Example 1:**

    Input: s = "rabbbit", t = "rabbit"
    Output: 3

**Example 2:**

    Input: s = "babgbag", t = "bag"
    Output: 5

**Constraints:**

- `1 <= s.length, t.length <= 1000`
- The answer fits in 32 bits

---

**Hints — try each one before reading on:**
1. At s[i] vs t[j]: always allowed to skip s[i]; on a match you may also consume both.
2. dp[i][j] = dp[i+1][j] + (match ? dp[i+1][j+1] : 0).

**Target:** O(m·n) time, O(n) space
