---
created: "2026-06-07"
id: nc-sliding-window-maximum
source: imported:neetcode-150
study:
  answer: |-
    Monotonic deque of indices whose values decrease front→back. For each new element: pop smaller values off the back (they can never be a future max), push it, drop the front if it left the window, and read the front as the window's max once i ≥ k−1.

    Complexity: O(n) time, O(k) space
  kind: essay
  prompt: 'Solve "Sliding Window Maximum" (Sliding Window): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Sliding Window Maximum
---

**Pattern:** Sliding Window · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/sliding-window-maximum)

Given an array `nums` and a window size `k`, return the maximum of each of the `n - k + 1` windows as it slides left to right.

**Example 1:**

    Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
    Output: [3,3,5,5,6,7]

**Example 2:**

    Input: nums = [1], k = 1
    Output: [1]

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `1 <= k <= nums.length`

---

**Hints — try each one before reading on:**
1. A max-heap works but holds stale values; what structure keeps candidates sorted AND evictable from both ends?
2. Monotonically decreasing deque of indices: front is the max; pop the back while smaller than the new value.

**Target:** O(n) time, O(k) space
