---
created: "2026-06-07"
id: nc-min-cost-to-connect-all-points
source: imported:neetcode-150
study:
  answer: |-
    Prim's algorithm: keep a min-heap (or O(n²) array scan) of cheapest known distances from the growing tree to each unvisited point; repeatedly absorb the closest point, relaxing distances through it, summing edge costs.

    Complexity: O(n² log n) heap / O(n²) array time, O(n) space
  kind: essay
  prompt: 'Solve "Min Cost to Connect All Points" (Advanced Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Min Cost to Connect All Points
---

**Pattern:** Advanced Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/min-cost-to-connect-all-points)

Given 2-D `points`, the cost to connect two points is their **Manhattan distance**. Return the minimum total cost to make all points connected.

**Example 1:**

    Input: points = [[0,0],[2,2],[3,10],[5,2],[7,0]]
    Output: 20

**Example 2:**

    Input: points = [[3,12],[-2,5],[-4,1]]
    Output: 18

**Constraints:**

- `1 <= points.length <= 1000`
- `-10^6 <= x, y <= 10^6`
- All points are distinct

---

**Hints — try each one before reading on:**
1. That's a minimum spanning tree on the complete graph.
2. Prim's: grow from one point, always adding the cheapest edge to an outside point.

**Target:** O(n² log n) heap / O(n²) array time, O(n) space
