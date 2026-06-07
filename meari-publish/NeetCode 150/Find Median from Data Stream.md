---
created: "2026-06-07"
id: nc-find-median-from-data-stream
source: imported:neetcode-150
study:
  answer: |-
    Two heaps: push into the max-heap (small half), move its top to the min-heap, rebalance so sizes differ by at most one. Median is the bigger heap's root, or the average of both roots when equal.

    Complexity: O(log n) add, O(1) median; O(n) space
  kind: essay
  prompt: 'Solve "Find Median from Data Stream" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Find Median from Data Stream
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/find-median-from-data-stream)

Design `MedianFinder`: `addNum(num)` adds a value from a stream, `findMedian()` returns the median of everything added (the mean of the two middle values for even counts).

**Example 1:**

    Input: addNum(1); addNum(2); findMedian(); addNum(3); findMedian()
    Output: 1.5, 2.0

**Constraints:**

- `-10^5 <= num <= 10^5`
- Up to `5 * 10^4` calls

---

**Hints — try each one before reading on:**
1. Split the numbers into a small half and a large half.
2. Max-heap for the small half, min-heap for the large; keep sizes within one.

**Target:** O(log n) add, O(1) median; O(n) space
