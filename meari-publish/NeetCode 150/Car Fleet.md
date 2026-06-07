---
created: "2026-06-07"
id: nc-car-fleet
source: imported:neetcode-150
study:
  answer: |-
    Sort cars by position descending and compute arrival times. Sweep keeping the current fleet's arrival time (a stack of times works too): a later car with time ≤ the fleet ahead merges into it; a slower one starts a new fleet. Count the fleets.

    Complexity: O(n log n) time, O(n) space
  kind: essay
  prompt: 'Solve "Car Fleet" (Stack): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Car Fleet
---

**Pattern:** Stack · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/car-fleet)

`n` cars at positions `position[i]` with speeds `speed[i]` drive toward `target` (no passing — a faster car bumps into the one ahead and forms a **fleet** at its speed). Return how many fleets arrive at the target.

**Example 1:**

    Input: target = 12, position = [10,8,0,5,3], speed = [2,4,1,1,3]
    Output: 3
    Explanation: {10,8} meet at 12; {0} alone; {5,3} meet at 6.

**Constraints:**

- `n == position.length == speed.length`
- `1 <= n <= 10^5`
- All positions are distinct and `< target`

---

**Hints — try each one before reading on:**
1. A car's arrival time is (target − position) / speed; sort by starting position.
2. Scan from the car closest to the target: a car catches the fleet ahead iff its time ≤ that fleet's time.

**Target:** O(n log n) time, O(n) space
