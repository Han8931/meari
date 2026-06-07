---
created: "2026-06-07"
id: nc-combination-sum-ii
source: imported:neetcode-150
study:
  answer: |-
    Sort, backtrack with start index moving to i+1 on take; inside the loop skip nums[i] == nums[i−1] when i > start; prune when the candidate exceeds the remaining target (sorted input makes the break safe).

    Complexity: Exponential time, O(n) space
  kind: essay
  prompt: 'Solve "Combination Sum II" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Combination Sum II
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/combination-sum-ii)

Given candidates (with duplicates) and a `target`, return all unique combinations summing to `target` where each candidate is used **at most once**.

**Example 1:**

    Input: candidates = [10,1,2,7,6,1,5], target = 8
    Output: [[1,1,6],[1,2,5],[1,7],[2,6]]

**Example 2:**

    Input: candidates = [2,5,2,1,2], target = 5
    Output: [[1,2,2],[5]]

**Constraints:**

- `1 <= candidates.length <= 100`
- `1 <= candidates[i] <= 50`
- `1 <= target <= 30`

---

**Hints — try each one before reading on:**
1. Sort; advance to i+1 when recursing (single use).
2. Skip equal neighbors at the same level, like Subsets II.

**Target:** Exponential time, O(n) space
