---
created: "2026-06-07"
id: nc-happy-number
source: imported:neetcode-150
study:
  answer: |-
    Iterate n → sum of squared digits, recording values in a set; reaching 1 is happy, revisiting a value means a cycle (unhappy). Floyd's tortoise-hare replaces the set for O(1) space.

    Complexity: O(log n) per step, O(1) cycle-bounded space with Floyd
  kind: essay
  prompt: 'Solve "Happy Number" (Math & Geometry): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Happy Number
---

**Pattern:** Math & Geometry · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/happy-number)

Repeatedly replace a number with the **sum of the squares of its digits**. Return `true` if the process reaches 1 (a "happy" number), `false` if it loops forever.

**Example 1:**

    Input: n = 19
    Output: true
    Explanation: 1²+9²=82 → 68 → 100 → 1.

**Example 2:**

    Input: n = 2
    Output: false

**Constraints:**

- `1 <= n <= 2^31 - 1`

---

**Hints — try each one before reading on:**
1. The sequence either hits 1 or enters a cycle.
2. Detect the cycle with a seen-set — or Floyd's slow/fast on the digit function.

**Target:** O(log n) per step, O(1) cycle-bounded space with Floyd
