---
created: "2026-06-07"
id: nc-subsets
source: imported:neetcode-150
study:
  answer: |-
    DFS(i, path): record path at every call, then for each j ≥ i append nums[j], recurse with j+1, and pop. (Iteratively: start with [[]] and for each number append it to a copy of every subset so far.)

    Complexity: O(n · 2ⁿ) time, O(n) recursion space
  kind: essay
  prompt: 'Solve "Subsets" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Subsets
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/subsets)

Given an array `nums` of **unique** integers, return all possible subsets (the power set), without duplicates, in any order.

**Example 1:**

    Input: nums = [1,2,3]
    Output: [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]

**Example 2:**

    Input: nums = [0]
    Output: [[],[0]]

**Constraints:**

- `1 <= nums.length <= 10`
- All numbers are unique

---

**Hints — try each one before reading on:**
1. Each element is either in or out — recurse on both choices.
2. Or iterate: every new element doubles the existing subsets.

**Target:** O(n · 2ⁿ) time, O(n) recursion space
