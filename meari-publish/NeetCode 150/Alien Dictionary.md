---
created: "2026-06-07"
id: nc-alien-dictionary
source: imported:neetcode-150
study:
  answer: |-
    Compare adjacent words: the first mismatched characters give edge a → b; reject if a longer word precedes its own prefix. Topologically sort the letter graph (Kahn's or DFS); any cycle → "", else the order is the answer.

    Complexity: O(total characters) time, O(1) graph space (≤26 nodes)
  kind: essay
  prompt: 'Solve "Alien Dictionary" (Advanced Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Alien Dictionary
---

**Pattern:** Advanced Graphs · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/alien-dictionary)

Given `words` sorted lexicographically in an **unknown alphabet**, return any string of that alphabet's letters in a consistent order, or `""` if the ordering is contradictory.

**Example 1:**

    Input: words = ["wrt","wrf","er","ett","rftt"]
    Output: "wertf"

**Example 2:**

    Input: words = ["z","x","z"]
    Output: ""
    Explanation: z before x and x before z is a contradiction.

**Constraints:**

- `1 <= words.length <= 100`
- `1 <= words[i].length <= 100`, lowercase letters

---

**Hints — try each one before reading on:**
1. Each adjacent word pair's first differing letter yields one ordering edge.
2. Topologically sort those edges; a prefix conflict ("abc" before "ab") or a cycle means "".

**Target:** O(total characters) time, O(1) graph space (≤26 nodes)
