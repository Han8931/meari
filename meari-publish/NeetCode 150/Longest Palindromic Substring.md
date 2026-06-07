---
created: "2026-06-07"
id: nc-longest-palindromic-substring
source: imported:neetcode-150
study:
  answer: |-
    For each index expand outward twice (odd center i,i and even center i,i+1) while the ends match, tracking the longest span seen. Simpler than the DP table and the same worst case.

    Complexity: O(n²) time, O(1) space
  kind: essay
  prompt: 'Solve "Longest Palindromic Substring" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Palindromic Substring
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/longest-palindromic-substring)

Given a string `s`, return the **longest palindromic substring**.

**Example 1:**

    Input: s = "babad"
    Output: "bab"
    Explanation: "aba" is also accepted.

**Example 2:**

    Input: s = "cbbd"
    Output: "bb"

**Constraints:**

- `1 <= s.length <= 1000`

---

**Hints — try each one before reading on:**
1. Every palindrome has a center — expand around all 2n−1 of them.
2. Handle odd (single center) and even (pair center) expansions.

**Target:** O(n²) time, O(1) space
