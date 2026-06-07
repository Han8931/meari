---
created: "2026-06-07"
id: nc-missing-number
source: imported:neetcode-150
study:
  answer: |-
    Either arithmetic (expected sum − actual sum) or XOR-fold of every index 0..n with every element; both isolate the missing value without extra space.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Missing Number" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Missing Number
---

**Pattern:** Bit Manipulation · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/missing-number)

`nums` contains `n` distinct numbers from `[0, n]` — one value of the range is missing. Return it in O(n) time and O(1) space.

**Example 1:**

    Input: nums = [3,0,1]
    Output: 2

**Example 2:**

    Input: nums = [9,6,4,2,3,5,7,0,1]
    Output: 8

**Constraints:**

- `1 <= n <= 10^4`
- All numbers are unique

---

**Hints — try each one before reading on:**
1. Sum formula n(n+1)/2 minus the actual sum.
2. Or XOR all indices and values — the missing one survives.

**Target:** O(n) time, O(1) space
