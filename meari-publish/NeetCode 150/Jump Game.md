---
created: "2026-06-07"
id: nc-jump-game
source: imported:neetcode-150
study:
  answer: |-
    Greedy reach: scan left to right maintaining reach = max(reach, i + nums[i]); fail if i ever exceeds reach, succeed when reach covers the last index.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Jump Game" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Jump Game
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/jump-game)

`nums[i]` is your **maximum** jump length from index i. Starting at index 0, return `true` if you can reach the last index.

**Example 1:**

    Input: nums = [2,3,1,1,4]
    Output: true

**Example 2:**

    Input: nums = [3,2,1,0,4]
    Output: false
    Explanation: you always land on index 3 (jump 0).

**Constraints:**

- `1 <= nums.length <= 10^4`
- `0 <= nums[i] <= 10^5`

---

**Hints — try each one before reading on:**
1. Track the farthest index reachable so far.
2. Or walk backward shrinking the "goal" to any index that reaches it.

**Target:** O(n) time, O(1) space
