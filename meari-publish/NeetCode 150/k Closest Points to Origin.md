---
created: "2026-06-07"
id: nc-k-closest-points-to-origin
source: imported:neetcode-150
study:
  answer: |-
    Keep a size-k max-heap keyed on negated squared distance: push each point, pop when over k; the heap holds the answer. (heapq.nsmallest(k, points, key=dist2) is the one-liner.)

    Complexity: O(n log k) time, O(k) space
  kind: essay
  prompt: 'Solve "k Closest Points to Origin" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: k Closest Points to Origin
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/k-closest-points-to-origin)

Given `points` on a plane and an integer `k`, return the `k` points **closest to the origin** (Euclidean distance), in any order.

**Example 1:**

    Input: points = [[1,3],[-2,2]], k = 1
    Output: [[-2,2]]
    Explanation: sqrt(8) < sqrt(10).

**Example 2:**

    Input: points = [[3,3],[5,-1],[-2,4]], k = 2
    Output: [[3,3],[-2,4]]

**Constraints:**

- `1 <= k <= points.length <= 10^4`
- `-10^4 <= x, y <= 10^4`

---

**Hints — try each one before reading on:**
1. Compare squared distances — no sqrt needed.
2. A max-heap of size k keeps the closest seen; or heapify all and pop k.

**Target:** O(n log k) time, O(k) space
