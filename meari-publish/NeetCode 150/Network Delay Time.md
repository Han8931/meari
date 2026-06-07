---
created: "2026-06-07"
id: nc-network-delay-time
source: imported:neetcode-150
study:
  answer: |-
    Dijkstra from k with a min-heap of (distance, node): pop the closest unsettled node, relax its outgoing edges. Return the maximum settled distance, or −1 if some node was never reached.

    Complexity: O(E log V) time, O(V + E) space
  kind: essay
  prompt: 'Solve "Network Delay Time" (Advanced Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Network Delay Time
---

**Pattern:** Advanced Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/network-delay-time)

A signal leaves node `k` in a directed weighted graph of `n` nodes (`times[i] = (u, v, w)`). Return the time for **all** nodes to receive it, or `-1`.

**Example 1:**

    Input: times = [[2,1,1],[2,3,1],[3,4,1]], n = 4, k = 2
    Output: 2

**Example 2:**

    Input: times = [[1,2,1]], n = 2, k = 2
    Output: -1

**Constraints:**

- `1 <= k <= n <= 100`
- `1 <= times.length <= 6000`
- `1 <= w <= 100`

---

**Hints — try each one before reading on:**
1. Single-source shortest paths with non-negative weights — Dijkstra.
2. The answer is the LARGEST shortest distance.

**Target:** O(E log V) time, O(V + E) space
