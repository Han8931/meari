---
created: "2026-06-07"
id: nc-target-sum
source: imported:neetcode-150
study:
  answer: |-
    Reduce to counting subsets summing to (total + target)//2 (impossible if that's negative or odd) and run counting subset-sum DP: dp[s] += dp[s − num] per number, s descending. (Direct memoized DFS over (index, running sum) also passes.)

    Complexity: O(n · sum) time, O(sum) space
  kind: essay
  prompt: 'Solve "Target Sum" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Target Sum
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/target-sum)

Assign `+` or `-` to every element of `nums`; return how many assignments make the expression equal `target`.

**Example 1:**

    Input: nums = [1,1,1,1,1], target = 3
    Output: 5

**Example 2:**

    Input: nums = [1], target = 1
    Output: 1

**Constraints:**

- `1 <= nums.length <= 20`
- `0 <= nums[i] <= 1000`
- `-1000 <= target <= 1000`

---

**Hints — try each one before reading on:**
1. Choosing the + set P: sum(P) − (total − sum(P)) = target ⇒ sum(P) = (total + target) / 2.
2. It reduces to subset-sum COUNTING; parity/negative checks first.

**Target:** O(n · sum) time, O(sum) space
