---
created: "2026-06-07"
id: nc-min-stack
source: imported:neetcode-150
study:
  answer: |-
    Keep a parallel min-stack: on push, also push min(value, current min); on pop, pop both. getMin reads the min-stack's top. Every operation is O(1).

    Complexity: O(1) per operation, O(n) space
  kind: essay
  prompt: 'Solve "Min Stack" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Min Stack
---

**Pattern:** Stack · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/min-stack)

Design a stack supporting `push`, `pop`, `top`, and `getMin` — retrieving the minimum element — each in O(1) time.

**Example 1:**

    Input: push(-2), push(0), push(-3), getMin(), pop(), top(), getMin()
    Output: -3, 0, -2

**Constraints:**

- `-2^31 <= val <= 2^31 - 1`
- pop/top/getMin are always called on a non-empty stack
- Up to `3 * 10^4` calls

---

**Hints — try each one before reading on:**
1. Pops destroy information — store the minimum alongside each element.
2. A second stack holding "min so far at this depth" mirrors push/pop.

**Target:** O(1) per operation, O(n) space
