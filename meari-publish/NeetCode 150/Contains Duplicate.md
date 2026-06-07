---
created: "2026-06-07"
id: nc-contains-duplicate
source: imported:neetcode-150
study:
  answer: |-
    Insert every element into a hash set; if an element is already present, return true. Equivalently, return len(set(nums)) < len(nums).

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Contains Duplicate" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Contains Duplicate
---

**Pattern:** Arrays & Hashing · **Difficulty:** Easy · [LeetCode ↗](https://leetcode.com/problems/contains-duplicate)

Given an integer array `nums`, return `true` if any value appears **at least twice**, and `false` if every element is distinct.

**Example 1:**

    Input: nums = [1,2,3,1]
    Output: true
    Explanation: 1 appears at indices 0 and 3.

**Example 2:**

    Input: nums = [1,2,3,4]
    Output: false

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

**Hints — try each one before reading on:**
1. What lets you ask "have I seen this before?" in O(1)?
2. Compare the array's length with the size of its set.

**Target:** O(n) time, O(n) space
