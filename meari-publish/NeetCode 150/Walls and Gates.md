---
created: "2026-06-07"
id: nc-walls-and-gates
source: imported:neetcode-150
study:
  answer: |-
    Multi-source BFS seeded with every gate at distance 0; expand outward writing distance+1 into untouched empty rooms (their INF value doubles as the visited check).

    Complexity: O(rows · cols) time, O(rows · cols) space
  kind: essay
  prompt: 'Solve "Walls and Gates" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Walls and Gates
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/walls-and-gates)

You are given a grid of rooms: `-1` is a wall, `0` a gate, and `INF` an empty room. Fill each empty room with the distance to its **nearest gate** (leave INF if unreachable). Modify in place.

**Example 1:**

    Input: rooms = [[INF,-1,0,INF],[INF,INF,INF,-1],[INF,-1,INF,-1],[0,-1,INF,INF]]
    Output: [[3,-1,0,1],[2,2,1,-1],[1,-1,2,-1],[0,-1,3,4]]

**Constraints:**

- `1 <= m, n <= 250`
- `INF = 2^31 - 1`

---

**Hints — try each one before reading on:**
1. Nearest-gate distance for ALL rooms = BFS from all gates at once.
2. First visit to a room is necessarily the shortest.

**Target:** O(rows · cols) time, O(rows · cols) space
