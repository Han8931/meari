---
created: "2026-06-07"
id: nc-permutations
source: imported:neetcode-150
study:
  answer: |-
    Backtracking with a used-marker: at each depth try every unused number, recurse, unmark. Emit when the path holds all n. (In-place index swapping avoids the marker array.)

    Complexity: O(n · n!) time, O(n) space
  kind: essay
  prompt: 'Solve "Permutations" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Permutations
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/permutations)

Given an array of **distinct** integers, return all possible permutations, in any order.

**Example 1:**

    Input: nums = [1,2,3]
    Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]

**Example 2:**

    Input: nums = [0,1]
    Output: [[0,1],[1,0]]

**Constraints:**

- `1 <= nums.length <= 6`
- All integers are unique

---

**Hints — try each one before reading on:**
1. At each level, pick any unused element.
2. Track usage with a set/boolean array, or swap in place.

**Target:** O(n · n!) time, O(n) space
