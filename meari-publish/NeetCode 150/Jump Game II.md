---
created: "2026-06-07"
id: nc-jump-game-ii
source: imported:neetcode-150
study:
  answer: |-
    Implicit BFS: track the current jump's window end and the farthest reachable inside it; on hitting the window end, increment jumps and set the new end to that farthest. Done when the window covers the last index.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Jump Game II" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Jump Game II
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/jump-game-ii)

Same jumping rule, but reaching the last index is guaranteed. Return the **minimum number of jumps**.

**Example 1:**

    Input: nums = [2,3,1,1,4]
    Output: 2
    Explanation: 0 → 1 → 4.

**Example 2:**

    Input: nums = [2,3,0,1,4]
    Output: 2

**Constraints:**

- `1 <= nums.length <= 10^4`
- Reaching the end is guaranteed

---

**Hints — try each one before reading on:**
1. Think in BFS levels over indices: one jump expands a window.
2. When you reach the current window's end, a jump is forced — extend to the farthest seen.

**Target:** O(n) time, O(1) space
