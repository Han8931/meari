---
created: "2026-06-07"
id: nc-word-ladder
source: imported:neetcode-150
study:
  answer: |-
    BFS from beginWord over an adjacency built from wildcard patterns: map each pattern (word with one position starred) to its words, then neighbors of w are all words sharing a pattern. The BFS level at endWord (+1 for the start) is the answer; 0 if unreachable.

    Complexity: O(N · L²) time, O(N · L²) space
  kind: essay
  prompt: 'Solve "Word Ladder" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Word Ladder
---

**Pattern:** Graphs · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/word-ladder)

Given `beginWord`, `endWord`, and a `wordList`, return the length of the **shortest transformation sequence** changing one letter at a time, every intermediate word in the list — or 0.

**Example 1:**

    Input: beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
    Output: 5
    Explanation: hit → hot → dot → dog → cog.

**Example 2:**

    Input: same but wordList lacks "cog"
    Output: 0

**Constraints:**

- `1 <= wordList.length <= 5000`
- All words have the same length (≤ 10) and are lowercase

---

**Hints — try each one before reading on:**
1. Words are nodes; one-letter differences are edges — shortest path = BFS.
2. Wildcard buckets (h*t, ho*, *ot) find neighbors without 26·L probing per word.

**Target:** O(N · L²) time, O(N · L²) space
