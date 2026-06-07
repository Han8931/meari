---
created: "2026-06-07"
id: nc-palindrome-partitioning
source: imported:neetcode-150
study:
  answer: |-
    Backtrack(start, path): for each end > start where s[start:end] is a palindrome, push it and recurse from end; emit the path when start reaches the string's end.

    Complexity: O(n · 2ⁿ) time, O(n) space
  kind: essay
  prompt: 'Solve "Palindrome Partitioning" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Palindrome Partitioning
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/palindrome-partitioning)

Partition a string `s` so that **every substring** of the partition is a palindrome. Return all such partitions.

**Example 1:**

    Input: s = "aab"
    Output: [["a","a","b"],["aa","b"]]

**Example 2:**

    Input: s = "a"
    Output: [["a"]]

**Constraints:**

- `1 <= s.length <= 16`
- `s` contains only lowercase English letters

---

**Hints — try each one before reading on:**
1. Choose a palindromic prefix, recurse on the rest.
2. Check palindromes with two pointers (or precompute a DP table).

**Target:** O(n · 2ⁿ) time, O(n) space
