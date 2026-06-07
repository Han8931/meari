---
created: "2026-06-07"
id: nc-k-th-largest-element-in-an-array
source: imported:neetcode-150
study:
  answer: |-
    Either a size-k min-heap (root is the answer after one pass), or quickselect: partition around a pivot and recurse only into the side containing index n−k, giving average linear time.

    Complexity: O(n log k) heap / O(n) average quickselect
  kind: essay
  prompt: 'Solve "k-th Largest Element in an Array" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: k-th Largest Element in an Array
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/kth-largest-element-in-an-array)

Return the `k`-th largest element of `nums` — the k-th in **sorted order**, not the k-th distinct — ideally without fully sorting.

**Example 1:**

    Input: nums = [3,2,1,5,6,4], k = 2
    Output: 5

**Example 2:**

    Input: nums = [3,2,3,1,2,4,5,5,6], k = 4
    Output: 4

**Constraints:**

- `1 <= k <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`

---

**Hints — try each one before reading on:**
1. A size-k min-heap over the array does it in n log k.
2. Quickselect partitions toward the (n−k)-th index for average O(n).

**Target:** O(n log k) heap / O(n) average quickselect
