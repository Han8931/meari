---
created: "2026-06-07"
id: nc-combination-sum
source: imported:neetcode-150
study:
  answer: |-
    Backtrack(start, remaining, path): for each i ≥ start with candidate ≤ remaining, take it and recurse with the same i (reuse allowed); emit when remaining hits 0. The start index prevents permuted duplicates.

    Complexity: Exponential time, O(target/min) recursion depth
  kind: essay
  prompt: 'Solve "Combination Sum" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Combination Sum
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/combination-sum)

Given **distinct** candidates and a `target`, return all unique combinations summing to `target`. A candidate may be used **unlimited times**.

**Example 1:**

    Input: candidates = [2,3,6,7], target = 7
    Output: [[2,2,3],[7]]

**Example 2:**

    Input: candidates = [2,3,5], target = 8
    Output: [[2,2,2,2],[2,3,3],[3,5]]

**Constraints:**

- `1 <= candidates.length <= 30`
- `2 <= candidates[i] <= 40`, all distinct
- `1 <= target <= 40`

---

**Hints — try each one before reading on:**
1. Reuse is allowed — recurse on the SAME index after taking a candidate.
2. Pass a start index so combinations stay ordered and unique.

**Target:** Exponential time, O(target/min) recursion depth
