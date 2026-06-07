---
created: "2026-06-07"
id: nc-letter-combinations-of-a-phone-number
source: imported:neetcode-150
study:
  answer: |-
    Map each digit to its letters and DFS by position: for digits[i] try every letter, recurse to i+1, emit at the end. (itertools.product is the standard-library equivalent.)

    Complexity: O(n · 4ⁿ) time, O(n) space
  kind: essay
  prompt: 'Solve "Letter Combinations of a Phone Number" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Letter Combinations of a Phone Number
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/letter-combinations-of-a-phone-number)

Given a string of digits `2-9`, return all letter combinations the number could represent on a phone keypad, in any order.

**Example 1:**

    Input: digits = "23"
    Output: ["ad","ae","af","bd","be","bf","cd","ce","cf"]

**Example 2:**

    Input: digits = ""
    Output: []

**Constraints:**

- `0 <= digits.length <= 4`
- `digits[i]` is in `['2', '9']`

---

**Hints — try each one before reading on:**
1. It's a cartesian product over each digit's 3-4 letters.
2. DFS over digit positions appending each mapped letter.

**Target:** O(n · 4ⁿ) time, O(n) space
