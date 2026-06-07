---
created: "2026-06-07"
id: nc-valid-anagram
source: imported:neetcode-150
study:
  answer: |-
    Count character frequencies of s (increment) and t (decrement) in one 26-slot array (or Counter); all zeros means anagram. Counter(s) == Counter(t) is the idiomatic Python.

    Complexity: O(n) time, O(1) space (fixed alphabet)
  kind: essay
  prompt: 'Solve "Valid Anagram" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Valid Anagram
---

**Pattern:** Arrays & Hashing · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/valid-anagram)

Given two strings `s` and `t`, return `true` if `t` is an **anagram** of `s` — the same letters with the same counts, rearranged.

**Example 1:**

    Input: s = "anagram", t = "nagaram"
    Output: true

**Example 2:**

    Input: s = "rat", t = "car"
    Output: false

**Constraints:**

- `1 <= s.length, t.length <= 5 * 10^4`
- `s` and `t` consist of lowercase English letters

---

**Hints — try each one before reading on:**
1. Sorting works but costs O(n log n) — can counts do better?
2. 26 lowercase letters fit in a fixed-size array.

**Target:** O(n) time, O(1) space (fixed alphabet)
