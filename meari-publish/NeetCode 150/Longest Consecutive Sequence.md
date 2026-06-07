---
created: "2026-06-07"
id: nc-longest-consecutive-sequence
source: imported:neetcode-150
study:
  answer: |-
    Put everything in a set. For each value x with x−1 not in the set (a run's start), walk x+1, x+2… counting the streak. Each element is visited at most twice, so it's linear despite the nested loop.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Longest Consecutive Sequence" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Longest Consecutive Sequence
---

**Pattern:** Arrays & Hashing · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/longest-consecutive-sequence)

Given an unsorted array `nums`, return the length of the longest run of **consecutive integers** (the elements' positions don't matter). Your algorithm must run in O(n).

**Example 1:**

    Input: nums = [100,4,200,1,3,2]
    Output: 4
    Explanation: the run 1,2,3,4.

**Example 2:**

    Input: nums = [0,3,7,2,5,8,4,6,0,1]
    Output: 9

**Constraints:**

- `0 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

**Hints — try each one before reading on:**
1. Sorting gives O(n log n) — a set can answer "does x−1 exist?" instantly.
2. Only start counting from numbers that begin a run (x−1 absent).

**Target:** O(n) time, O(n) space
