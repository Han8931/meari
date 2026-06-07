---
created: "2026-06-07"
id: nc-word-search-ii
source: imported:neetcode-150
study:
  answer: |-
    Insert all words into a trie, then DFS from every cell carrying the current trie node: step only into children that exist in the trie, collect words at isWord nodes (then unset to avoid duplicates), mark cells visited during the path, and prune leaf trie nodes after use for speed.

    Complexity: O(cells · 4 · 3^(L−1)) time worst case, O(dictionary) space
  kind: essay
  prompt: 'Solve "Word Search II" (Tries): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Word Search II
---

**Pattern:** Tries · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/word-search-ii)

Given an `m x n` board of letters and a list `words`, return every word constructible from sequentially adjacent cells (horizontally or vertically), using each cell at most once per word.

**Example 1:**

    Input: board = [["o","a","a","n"],["e","t","a","e"],["i","h","k","r"],["i","f","l","v"]], words = ["oath","pea","eat","rain"]
    Output: ["eat","oath"]

**Constraints:**

- `1 <= m, n <= 12`
- `1 <= words.length <= 3 * 10^4`, `1 <= words[i].length <= 10`

---

**Hints — try each one before reading on:**
1. Searching each word separately repeats grid work — search ALL words at once.
2. Build a trie of the dictionary and DFS the grid while walking the trie; prune branches with no trie child.

**Target:** O(cells · 4 · 3^(L−1)) time worst case, O(dictionary) space
