---
created: "2026-06-07"
id: nc-valid-parentheses
source: imported:neetcode-150
study:
  answer: |-
    Push openers; on a closer, the stack top must be the matching opener (pop it), else invalid. Valid iff the stack ends empty.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Valid Parentheses" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Valid Parentheses
---

**Pattern:** Stack · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/valid-parentheses)

Given a string `s` containing just `()[]{}`, decide whether it is valid: every opener is closed by the same bracket type, in the correct order.

**Example 1:**

    Input: s = "()[]{}"
    Output: true

**Example 2:**

    Input: s = "(]"
    Output: false

**Constraints:**

- `1 <= s.length <= 10^4`
- `s` consists only of `()[]{}`

---

**Hints — try each one before reading on:**
1. The most recently opened bracket must close first — that's a stack.
2. Map each closer to its expected opener.

**Target:** O(n) time, O(n) space
