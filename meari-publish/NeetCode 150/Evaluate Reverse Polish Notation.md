---
created: "2026-06-07"
id: nc-evaluate-reverse-polish-notation
source: imported:neetcode-150
study:
  answer: |-
    Scan tokens: push numbers; for an operator pop b then a, push a op b. For division use int(a / b) so truncation goes toward zero. The final stack element is the result.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Evaluate Reverse Polish Notation" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Evaluate Reverse Polish Notation
---

**Pattern:** Stack · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/evaluate-reverse-polish-notation)

Evaluate an arithmetic expression in **Reverse Polish (postfix) Notation** given as a token array. Operators are `+ - * /`; division truncates toward zero.

**Example 1:**

    Input: tokens = ["2","1","+","3","*"]
    Output: 9
    Explanation: ((2 + 1) * 3).

**Example 2:**

    Input: tokens = ["4","13","5","/","+"]
    Output: 6
    Explanation: (4 + (13 / 5)).

**Constraints:**

- `1 <= tokens.length <= 10^4`
- The expression is always valid

---

**Hints — try each one before reading on:**
1. Operands wait on a stack; an operator consumes the top two.
2. Mind the operand order for − and /, and truncate division toward zero.

**Target:** O(n) time, O(n) space
