---
created: "2026-06-07"
id: nc-single-number
source: imported:neetcode-150
study:
  answer: |-
    XOR-fold the whole array: the pairs annihilate, leaving the unique number.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Single Number" (Bit Manipulation): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Single Number
---

**Pattern:** Bit Manipulation · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/single-number)

Every element of `nums` appears **twice** except one. Find that one, in linear time and O(1) space.

**Example 1:**

    Input: nums = [2,2,1]
    Output: 1

**Example 2:**

    Input: nums = [4,1,2,1,2]
    Output: 4

**Constraints:**

- `1 <= nums.length <= 3 * 10^4`
- Exactly one element appears once

---

**Hints — try each one before reading on:**
1. x ^ x = 0 and XOR is commutative.
2. XOR everything together.

**Target:** O(n) time, O(1) space
