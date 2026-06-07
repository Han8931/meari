---
created: "2026-06-07"
id: nc-hand-of-straights
source: imported:neetcode-150
study:
  answer: |-
    Counter plus sorted unique values (or a min-heap): take the smallest value with remaining count and decrement counts for the k consecutive values, failing if any is missing. Succeed when all counts drain.

    Complexity: O(n log n) time, O(n) space
  kind: essay
  prompt: 'Solve "Hand of Straights" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Hand of Straights
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/hand-of-straights)

Given `hand` and `groupSize`, return `true` if the cards can be rearranged into groups of exactly `groupSize` **consecutive** cards.

**Example 1:**

    Input: hand = [1,2,3,6,2,3,4,7,8], groupSize = 3
    Output: true
    Explanation: [1,2,3], [2,3,4], [6,7,8].

**Example 2:**

    Input: hand = [1,2,3,4,5], groupSize = 4
    Output: false

**Constraints:**

- `1 <= hand.length <= 10^4`
- `1 <= groupSize <= hand.length`

---

**Hints — try each one before reading on:**
1. The smallest remaining card MUST start a run.
2. Count cards; repeatedly consume v, v+1, …, v+k−1 from the minimum.

**Target:** O(n log n) time, O(n) space
