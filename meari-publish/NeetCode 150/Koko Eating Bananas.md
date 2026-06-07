---
created: "2026-06-07"
id: nc-koko-eating-bananas
source: imported:neetcode-150
study:
  answer: |-
    Binary search on speed: for mid, compute total hours as the sum of ceil(pile/mid); if it fits in h search lower, else higher. The smallest feasible speed is the answer — the canonical "binary search on a monotonic predicate".

    Complexity: O(n log max(pile)) time, O(1) space
  kind: essay
  prompt: 'Solve "Koko Eating Bananas" (Binary Search): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Koko Eating Bananas
---

**Pattern:** Binary Search · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/koko-eating-bananas)

Koko has `piles` of bananas and `h` hours. Each hour she picks one pile and eats `k` bananas from it (finishing a smaller pile ends that hour). Return the **minimum** integer speed `k` that finishes everything within `h` hours.

**Example 1:**

    Input: piles = [3,6,7,11], h = 8
    Output: 4

**Example 2:**

    Input: piles = [30,11,23,4,20], h = 5
    Output: 30

**Constraints:**

- `1 <= piles.length <= 10^4`
- `piles.length <= h <= 10^9`
- `1 <= piles[i] <= 10^9`

---

**Hints — try each one before reading on:**
1. You can't search positions — search the ANSWER: speeds from 1 to max(pile).
2. hours(k) = Σ ceil(pile / k) is monotonic in k.

**Target:** O(n log max(pile)) time, O(1) space
