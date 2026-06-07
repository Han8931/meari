---
created: "2026-06-07"
id: nc-product-of-array-except-self
source: imported:neetcode-150
study:
  answer: |-
    First pass left→right writes prefix products into the answer array; second pass right→left multiplies each slot by a running suffix product. No division, and the output array doesn't count as extra space.

    Complexity: O(n) time, O(1) extra space
  kind: essay
  prompt: 'Solve "Product of Array Except Self" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Product of Array Except Self
---

**Pattern:** Arrays & Hashing · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/product-of-array-except-self)

Given an integer array `nums`, return an array `answer` where `answer[i]` is the product of every element **except** `nums[i]`. You must run in O(n) and may **not** use division.

**Example 1:**

    Input: nums = [1,2,3,4]
    Output: [24,12,8,6]

**Example 2:**

    Input: nums = [-1,1,0,-3,3]
    Output: [0,0,9,0,0]

**Constraints:**

- `2 <= nums.length <= 10^5`
- Products fit in a 32-bit integer

---

**Hints — try each one before reading on:**
1. answer[i] = (product of everything left of i) × (product of everything right of i).
2. Two sweeps: build prefix products into the output, then multiply a running suffix product in from the right.

**Target:** O(n) time, O(1) extra space
