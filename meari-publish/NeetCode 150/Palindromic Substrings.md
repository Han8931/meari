---
created: "2026-06-07"
id: nc-palindromic-substrings
source: imported:neetcode-150
study:
  answer: |-
    Expand around all 2n−1 centers, incrementing a counter for every matching expansion (each is a distinct palindrome).

    Complexity: O(n²) time, O(1) space
  kind: essay
  prompt: 'Solve "Palindromic Substrings" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Palindromic Substrings
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/palindromic-substrings)

Return the number of palindromic **substrings** in `s` (substrings with different positions count separately).

**Example 1:**

    Input: s = "abc"
    Output: 3
    Explanation: "a", "b", "c".

**Example 2:**

    Input: s = "aaa"
    Output: 6
    Explanation: "a"x3, "aa"x2, "aaa".

**Constraints:**

- `1 <= s.length <= 1000`

---

**Hints — try each one before reading on:**
1. Same center-expansion as Longest Palindromic Substring.
2. Each successful expansion step counts one palindrome.

**Target:** O(n²) time, O(1) space
