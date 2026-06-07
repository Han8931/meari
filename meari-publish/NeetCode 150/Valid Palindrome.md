---
created: "2026-06-07"
id: nc-valid-palindrome
source: imported:neetcode-150
study:
  answer: |-
    Left and right pointers move inward, each skipping non-alphanumerics; compare lowercased characters and fail on mismatch until the pointers cross.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Valid Palindrome" (Two Pointers): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Valid Palindrome
---

**Pattern:** Two Pointers · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/valid-palindrome)

A phrase is a palindrome if, after lowercasing and removing all non-alphanumeric characters, it reads the same forward and backward. Given a string `s`, return `true` if it is a palindrome.

**Example 1:**

    Input: s = "A man, a plan, a canal: Panama"
    Output: true
    Explanation: "amanaplanacanalpanama" is a palindrome.

**Example 2:**

    Input: s = "race a car"
    Output: false

**Constraints:**

- `1 <= s.length <= 2 * 10^5`
- `s` consists of printable ASCII characters

---

**Hints — try each one before reading on:**
1. Two pointers from both ends; skip the junk characters.
2. str.isalnum() and lower() do the filtering.

**Target:** O(n) time, O(1) space
