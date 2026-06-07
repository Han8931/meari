---
created: "2026-06-07"
id: nc-group-anagrams
source: imported:neetcode-150
study:
  answer: |-
    Bucket words in a dict keyed by their canonical form — tuple of 26 letter counts (O(k) per word) or the sorted word (O(k log k)). Return the dict's values.

    Complexity: O(n·k) time, O(n·k) space
  kind: essay
  prompt: 'Solve "Group Anagrams" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Group Anagrams
---

**Pattern:** Arrays & Hashing · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/group-anagrams)

Given an array of strings `strs`, group the anagrams together. Return the groups in any order.

**Example 1:**

    Input: strs = ["eat","tea","tan","ate","nat","bat"]
    Output: [["eat","tea","ate"],["tan","nat"],["bat"]]

**Constraints:**

- `1 <= strs.length <= 10^4`
- `0 <= strs[i].length <= 100`
- `strs[i]` consists of lowercase English letters

---

**Hints — try each one before reading on:**
1. Anagrams share something canonical — what key collapses them together?
2. Sorted word works; a 26-count tuple avoids the sort.

**Target:** O(n·k) time, O(n·k) space
