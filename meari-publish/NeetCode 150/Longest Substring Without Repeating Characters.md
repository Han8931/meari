---
created: "2026-06-07"
id: nc-longest-substring-without-repeating-characters
source: imported:neetcode-150
study:
  answer: |-
    Sliding window with a set of characters in the window: extend right; while s[right] is already present, evict s[left] and advance left. Track the max window length. A last-index map lets left jump instead of crawling.

    Complexity: O(n) time, O(min(n, alphabet)) space
  kind: essay
  prompt: 'Solve "Longest Substring Without Repeating Characters" (Sliding Window): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Substring Without Repeating Characters
---

**Pattern:** Sliding Window · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/longest-substring-without-repeating-characters)

Given a string `s`, return the length of the longest **substring** with no repeated characters.

**Example 1:**

    Input: s = "abcabcbb"
    Output: 3
    Explanation: "abc".

**Example 2:**

    Input: s = "bbbbb"
    Output: 1

**Constraints:**

- `0 <= s.length <= 5 * 10^4`
- `s` consists of English letters, digits, symbols and spaces

---

**Hints — try each one before reading on:**
1. Grow a window right; on a repeat, shrink from the left.
2. A set (or last-seen index map) tells you when and how far to shrink.

**Target:** O(n) time, O(min(n, alphabet)) space
