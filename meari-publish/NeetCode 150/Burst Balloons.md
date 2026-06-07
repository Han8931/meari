---
created: "2026-06-07"
id: nc-burst-balloons
source: imported:neetcode-150
study:
  answer: |-
    Pad with 1s; dp(l, r) = best for the open interval (l, r): choose k as the final burst, earning nums[l]·nums[k]·nums[r] + dp(l, k) + dp(k, r); evaluate intervals by increasing length (or memoized recursion). dp(0, n+1) is the answer.

    Complexity: O(n³) time, O(n²) space
  kind: essay
  prompt: 'Solve "Burst Balloons" (2-D Dynamic Programming): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Burst Balloons
---

**Pattern:** 2-D Dynamic Programming · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/burst-balloons)

Burst all `n` balloons; bursting balloon i earns `nums[left] * nums[i] * nums[right]` where left/right are the **still-unburst** neighbors (out-of-bounds counts as 1). Return the maximum total coins.

**Example 1:**

    Input: nums = [3,1,5,8]
    Output: 167
    Explanation: [3,1,5,8] → [3,5,8] → [3,8] → [8] → []; 3·1·5 + 3·5·8 + 1·3·8 + 1·8·1.

**Example 2:**

    Input: nums = [1,5]
    Output: 10

**Constraints:**

- `1 <= n <= 300`
- `0 <= nums[i] <= 100`

---

**Hints — try each one before reading on:**
1. Think about the LAST balloon burst in a range, not the first — its neighbors are then the range's fixed borders.
2. Interval DP on the array padded with 1s.

**Target:** O(n³) time, O(n²) space
