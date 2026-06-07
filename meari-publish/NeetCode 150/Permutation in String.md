---
created: "2026-06-07"
id: nc-permutation-in-string
source: imported:neetcode-150
study:
  answer: |-
    Fixed window of len(s1) over s2 with letter counts: add the entering char, drop the leaving one, compare to s1's counts (track a "matched letters" counter to make the compare O(1)).

    Complexity: O(n) time, O(26) space
  kind: essay
  prompt: 'Solve "Permutation in String" (Sliding Window): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Permutation in String
---

**Pattern:** Sliding Window · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/permutation-in-string)

Given strings `s1` and `s2`, return `true` if `s2` contains a **permutation** of `s1` as a contiguous substring.

**Example 1:**

    Input: s1 = "ab", s2 = "eidbaooo"
    Output: true
    Explanation: s2 contains "ba".

**Example 2:**

    Input: s1 = "ab", s2 = "eidboaoo"
    Output: false

**Constraints:**

- `1 <= s1.length, s2.length <= 10^4`
- Both consist of lowercase English letters

---

**Hints — try each one before reading on:**
1. A permutation of s1 is any window of len(s1) with identical letter counts.
2. Slide a fixed-size window updating two counts — or a matches counter over 26 letters.

**Target:** O(n) time, O(26) space
