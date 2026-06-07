---
created: "2026-06-07"
id: nc-reconstruct-itinerary
source: imported:neetcode-150
study:
  answer: |-
    Hierholzer's algorithm: adjacency lists sorted (use a heap or reversed sort for pop efficiency); DFS from JFK always taking the smallest unused edge, appending each airport to the route when its edges are exhausted; the reversed post-order is the itinerary.

    Complexity: O(E log E) time, O(E) space
  kind: essay
  prompt: 'Solve "Reconstruct Itinerary" (Advanced Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Reconstruct Itinerary
---

**Pattern:** Advanced Graphs · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/reconstruct-itinerary)

Given airline `tickets` as `[from, to]` pairs, reconstruct the full trip starting at `"JFK"`, using **every ticket exactly once**; among valid itineraries return the lexicographically smallest.

**Example 1:**

    Input: tickets = [["MUC","LHR"],["JFK","MUC"],["SFO","SJC"],["LHR","SFO"]]
    Output: ["JFK","MUC","LHR","SFO","SJC"]

**Example 2:**

    Input: tickets = [["JFK","SFO"],["JFK","ATL"],["SFO","ATL"],["ATL","JFK"],["ATL","SFO"]]
    Output: ["JFK","ATL","JFK","SFO","ATL","SFO"]

**Constraints:**

- `1 <= tickets.length <= 300`
- A valid itinerary always exists

---

**Hints — try each one before reading on:**
1. Using every edge exactly once is an Eulerian path.
2. Hierholzer: DFS consuming edges (smallest destination first), append airports on backtrack, reverse.

**Target:** O(E log E) time, O(E) space
