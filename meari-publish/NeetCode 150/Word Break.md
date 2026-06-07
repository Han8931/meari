---
created: "2026-06-07"
id: nc-word-break
source: imported:neetcode-150
study:
  answer: |-
    DP over prefixes with the dictionary as a set: dp[0] = True; dp[i] = any(dp[i−len(w)] and s[i−len(w):i] == w for w in dict). Answer dp[n].

    Complexity: O(n · dict · L) time, O(n) space
  kind: essay
  prompt: 'Solve "Word Break" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Word Break
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/word-break)

Given a string `s` and a dictionary `wordDict`, return `true` if `s` can be segmented into a sequence of dictionary words (reuse allowed).

**Example 1:**

    Input: s = "leetcode", wordDict = ["leet","code"]
    Output: true

**Example 2:**

    Input: s = "catsandog", wordDict = ["cats","dog","sand","and","cat"]
    Output: false

**Constraints:**

- `1 <= s.length <= 300`
- `1 <= wordDict.length <= 1000`

---

**Hints — try each one before reading on:**
1. dp[i] = "s[:i] is breakable".
2. dp[i] is true if some word w ends at i with dp[i−len(w)] true.

**Target:** O(n · dict · L) time, O(n) space
