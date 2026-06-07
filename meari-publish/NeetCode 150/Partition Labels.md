---
created: "2026-06-07"
id: nc-partition-labels
source: imported:neetcode-150
study:
  answer: |-
    Precompute each letter's last index; sweep extending end = max(end, last[s[i]]); when i == end the part is forced closed — record its length and start anew.

    Complexity: O(n) time, O(1) space
  kind: essay
  prompt: 'Solve "Partition Labels" (Greedy): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Partition Labels
---

**Pattern:** Greedy · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/partition-labels)

Split `s` into as many parts as possible so that **each letter appears in at most one part**; return the part sizes in order.

**Example 1:**

    Input: s = "ababcbacadefegdehijhklij"
    Output: [9,7,8]
    Explanation: "ababcbaca" | "defegde" | "hijhklij".

**Constraints:**

- `1 <= s.length <= 500`
- `s` consists of lowercase English letters

---

**Hints — try each one before reading on:**
1. A part can't end before the last occurrence of any letter inside it.
2. Track the running max of last-occurrence indices; close a part when i reaches it.

**Target:** O(n) time, O(1) space
