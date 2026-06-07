---
created: "2026-06-07"
id: nc-valid-parenthesis-string
source: imported:neetcode-150
study:
  answer: |-
    Greedy interval: maintain [lo, hi] of feasible open counts — '(' adds to both, ')' subtracts from both, '*' does lo−1/hi+1; fail when hi < 0, clamp lo to 0, and accept iff lo == 0 after the scan.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Valid Parenthesis String" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Valid Parenthesis String
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/valid-parenthesis-string)

Given a string of `'('`, `')'`, and `'*'` — where `'*'` may stand for `'('`, `')'`, or nothing — return `true` if it can be a valid parenthesis string.

**Example 1:**

    Input: s = "()"
    Output: true

**Example 2:**

    Input: s = "(*))"
    Output: true

**Constraints:**

- `1 <= s.length <= 100`

---

**Hints — try each one before reading on:**
1. Track a RANGE of possible open-bracket counts.
2. lo and hi: '*' moves them apart; hi < 0 fails; clamp lo at 0; valid iff lo == 0 at the end.

**Target:** O(n) time, O(1) space
