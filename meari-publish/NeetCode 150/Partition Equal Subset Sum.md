---
created: "2026-06-07"
id: nc-partition-equal-subset-sum
source: imported:neetcode-150
study:
  answer: |-
    If the total is odd return False; otherwise track all reachable subset sums in a set, folding each number in (new = {s, s+num}); True iff total//2 appears. A bitmask (bits = bits | bits << num) is the compact form.

    Complexity: O(n · target) time, O(target) space
  kind: essay
  prompt: 'Solve "Partition Equal Subset Sum" (1-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Partition Equal Subset Sum
---

**Pattern:** 1-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/partition-equal-subset-sum)

Return `true` if `nums` can be split into two subsets with **equal sums**.

**Example 1:**

    Input: nums = [1,5,11,5]
    Output: true
    Explanation: [1,5,5] and [11].

**Example 2:**

    Input: nums = [1,2,3,5]
    Output: false

**Constraints:**

- `1 <= nums.length <= 200`
- `1 <= nums[i] <= 100`

---

**Hints — try each one before reading on:**
1. Equivalent to reaching sum/2 with a subset; odd total → no.
2. 0/1 subset-sum DP over a set of reachable sums (or a bitset).

**Target:** O(n · target) time, O(target) space
