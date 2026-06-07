---
created: "2026-06-07"
id: nc-subsets-ii
source: imported:neetcode-150
study:
  answer: |-
    Sort, then subset-DFS with the rule: within one level's loop, skip nums[j] if j > start and nums[j] == nums[j−1] — equal values may only extend, not restart, so each multiset appears once.

    Complexity: O(n · 2ⁿ) time, O(n) space
  kind: essay
  prompt: 'Solve "Subsets II" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Subsets II
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/subsets-ii)

Given an array `nums` that **may contain duplicates**, return all possible subsets without duplicate subsets.

**Example 1:**

    Input: nums = [1,2,2]
    Output: [[],[1],[1,2],[1,2,2],[2],[2,2]]

**Constraints:**

- `1 <= nums.length <= 10`
- `-10 <= nums[i] <= 10`

---

**Hints — try each one before reading on:**
1. Sort so duplicates sit together.
2. At each level, skip a value equal to its left neighbor (when that neighbor wasn't taken at this level).

**Target:** O(n · 2ⁿ) time, O(n) space
