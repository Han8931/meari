---
created: "2026-06-07"
id: nc-generate-parentheses
source: imported:neetcode-150
study:
  answer: |-
    DFS/backtracking carrying (current string, open, close): append '(' if open < n, append ')' if close < open, emit when the string reaches 2n. The two rules make every produced string valid by construction.

    Complexity: O(4ⁿ/√n) outputs (Catalan), O(n) recursion depth
  kind: essay
  prompt: 'Solve "Generate Parentheses" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Generate Parentheses
---

**Pattern:** Stack · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/generate-parentheses)

Given `n` pairs of parentheses, generate **all combinations** of well-formed parentheses.

**Example 1:**

    Input: n = 3
    Output: ["((()))","(()())","(())()","()(())","()()()"]

**Example 2:**

    Input: n = 1
    Output: ["()"]

**Constraints:**

- `1 <= n <= 8`

---

**Hints — try each one before reading on:**
1. Build left to right: you may add '(' while open < n, and ')' while closed < open.
2. Backtrack with the running string and both counts.

**Target:** O(4ⁿ/√n) outputs (Catalan), O(n) recursion depth
