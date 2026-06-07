---
created: "2026-06-07"
id: nc-word-search
source: imported:neetcode-150
study:
  answer: |-
    For each starting cell, DFS(r, c, i): fail on bounds/mismatch, succeed when i reaches the word's end; temporarily overwrite board[r][c] (e.g. '#') around the four recursive calls to prevent reuse.

    Complexity: O(cells · 3^L) time, O(L) recursion space
  kind: essay
  prompt: 'Solve "Word Search" (Backtracking): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Word Search
---

**Pattern:** Backtracking · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/word-search)

Given an `m x n` board of characters and a `word`, return `true` if the word exists as a path of **sequentially adjacent** cells (horizontal/vertical), each cell used at most once.

**Example 1:**

    Input: board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCCED"
    Output: true

**Example 2:**

    Input: same board, word = "ABCB"
    Output: false

**Constraints:**

- `1 <= m, n <= 6`
- `1 <= word.length <= 15`

---

**Hints — try each one before reading on:**
1. DFS from every cell matching the word character by character.
2. Mark a cell visited by mutating it during the recursion; restore on backtrack.

**Target:** O(cells · 3^L) time, O(L) recursion space
