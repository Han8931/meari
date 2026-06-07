---
created: "2026-06-07"
id: nc-minimum-interval-to-include-each-query
source: imported:neetcode-150
study:
  answer: |-
    Sort intervals by start and queries ascending (remembering original positions). For each query, push every interval starting ≤ q as (size, end), then pop while the top's end < q; the surviving top's size answers the query. Each interval enters and leaves the heap once.

    Complexity: O((n + q) log n + q log q) time, O(n + q) space
  kind: essay
  prompt: 'Solve "Minimum Interval to Include Each Query" (Intervals): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Minimum Interval to Include Each Query
---

**Pattern:** Intervals · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/minimum-interval-to-cover-each-query)

For each query `q`, find the **size** (right − left + 1) of the smallest interval containing `q`, or `-1`. Return the answers in query order.

**Example 1:**

    Input: intervals = [[1,4],[2,4],[3,6],[4,4]], queries = [2,3,4,5]
    Output: [3,3,1,4]

**Example 2:**

    Input: intervals = [[2,3],[2,5],[1,8],[20,25]], queries = [2,19,5,22]
    Output: [2,-1,4,6]

**Constraints:**

- `1 <= intervals.length, queries.length <= 10^5`

---

**Hints — try each one before reading on:**
1. Process queries in sorted order so intervals can be added once.
2. Min-heap by interval size: push intervals whose start ≤ query, lazily pop those that already ended.

**Target:** O((n + q) log n + q log q) time, O(n + q) space
