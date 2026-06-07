---
created: "2026-06-07"
id: nc-merge-triplets-to-form-target
source: imported:neetcode-150
study:
  answer: |-
    Filter out triplets exceeding the target in any coordinate; the rest can be merged freely, so the answer is whether every coordinate of the target appears at full value among the survivors.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Merge Triplets to Form Target" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Merge Triplets to Form Target
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/merge-triplets-to-form-target)

You may repeatedly pick two triplets and replace one with their element-wise **maximum**. Return `true` if some sequence of merges can produce `target` among the triplets.

**Example 1:**

    Input: triplets = [[2,5,3],[1,8,4],[1,7,5]], target = [2,7,5]
    Output: true
    Explanation: merge [2,5,3] with [1,7,5] → [2,7,5].

**Example 2:**

    Input: triplets = [[3,4,5],[4,5,6]], target = [3,2,5]
    Output: false

**Constraints:**

- `1 <= triplets.length <= 10^5`
- All values in `[1, 1000]`

---

**Hints — try each one before reading on:**
1. A triplet exceeding the target ANYWHERE is unusable.
2. Among usable triplets, just check each target coordinate is achieved.

**Target:** O(n) time, O(1) space
