---
created: "2026-06-07"
id: nc-cheapest-flights-within-k-stops
source: imported:neetcode-150
study:
  answer: |-
    Bellman-Ford variant: run k+1 relaxation rounds, each relaxing every edge against a COPY of last round's distances (so a round uses at most one extra hop). dist[dst] after the rounds is the answer, or −1.

    Complexity: O(k · E) time, O(V) space
  kind: essay
  prompt: 'Solve "Cheapest Flights Within K Stops" (Advanced Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Cheapest Flights Within K Stops
---

**Pattern:** Advanced Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/cheapest-flights-within-k-stops)

Given flights `(from, to, price)`, find the cheapest price from `src` to `dst` with **at most `k` stops** in between, or `-1`.

**Example 1:**

    Input: n = 4, flights = [[0,1,100],[1,2,100],[2,0,100],[1,3,600],[2,3,200]], src = 0, dst = 3, k = 1
    Output: 700
    Explanation: 0 → 1 → 3; the cheaper 0 → 1 → 2 → 3 needs 2 stops.

**Example 2:**

    Input: n = 3, flights = [[0,1,100],[1,2,100],[0,2,500]], src = 0, dst = 2, k = 0
    Output: 500

**Constraints:**

- `1 <= n <= 100`
- `0 <= k < n`
- No duplicate flights

---

**Hints — try each one before reading on:**
1. Plain Dijkstra breaks: a pricier path with fewer stops can still win.
2. Bellman-Ford limited to k+1 rounds — relax all edges per round on a snapshot.

**Target:** O(k · E) time, O(V) space
