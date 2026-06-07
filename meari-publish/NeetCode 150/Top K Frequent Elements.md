---
created: "2026-06-07"
id: nc-top-k-frequent-elements
source: imported:neetcode-150
study:
  answer: |-
    Count frequencies, then bucket sort: an array of buckets indexed by count (0..n), each holding the values with that frequency; walk buckets from high to low collecting k values. Avoids the heap's log factor.

    Complexity: O(n) time, O(n) space
  kind: essay
  prompt: 'Solve "Top K Frequent Elements" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Top K Frequent Elements
---

**Pattern:** Arrays & Hashing · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/top-k-frequent-elements)

Given an integer array `nums` and an integer `k`, return the `k` most frequent elements, in any order.

**Example 1:**

    Input: nums = [1,1,1,2,2,3], k = 2
    Output: [1,2]

**Example 2:**

    Input: nums = [1], k = 1
    Output: [1]

**Constraints:**

- `1 <= nums.length <= 10^5`
- `k` is between 1 and the number of distinct elements
- The answer is guaranteed unique

---

**Hints — try each one before reading on:**
1. A heap gives O(n log k) — but frequencies are bounded by n.
2. Bucket index = frequency: buckets[count] holds the values with that count.

**Target:** O(n) time, O(n) space
