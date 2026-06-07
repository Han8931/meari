---
created: "2026-06-07"
id: nc-longest-repeating-character-replacement
source: imported:neetcode-150
study:
  answer: |-
    Sliding window with letter counts and the max single-letter count seen. If window size − maxFreq > k, slide left by one (no need to recompute maxFreq — a stale value only keeps the window from growing, never produces a wrong answer). The largest valid window is the answer.

    Complexity: O(n) time, O(26) space
  kind: essay
  prompt: 'Solve "Longest Repeating Character Replacement" (Sliding Window): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Repeating Character Replacement
---

**Pattern:** Sliding Window · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/longest-repeating-character-replacement)

Given a string `s` of uppercase letters and an integer `k`, you may change at most `k` characters. Return the length of the longest substring containing a single repeated letter you can achieve.

**Example 1:**

    Input: s = "ABAB", k = 2
    Output: 4
    Explanation: replace both A's (or both B's).

**Example 2:**

    Input: s = "AABABBA", k = 1
    Output: 4

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of uppercase English letters
- `0 <= k <= s.length`

---

**Hints — try each one before reading on:**
1. A window is valid while window_len − count(most frequent letter) ≤ k.
2. maxFreq never needs to decrease — the answer only improves when it grows.

**Target:** O(n) time, O(26) space
