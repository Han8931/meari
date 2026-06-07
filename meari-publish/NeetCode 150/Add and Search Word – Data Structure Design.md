---
created: "2026-06-07"
id: nc-add-and-search-word-data-structure-design
source: imported:neetcode-150
study:
  answer: |-
    Standard trie for addWord. search does DFS(node, i): a literal character follows its single child; '.' branches into every child, succeeding if any branch completes with isWord set.

    Complexity: O(L) add; search worst-case O(26^d · L) with dots, O(L) without
  kind: essay
  prompt: 'Solve "Add and Search Word – Data Structure Design" (Tries): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Add and Search Word – Data Structure Design
---

**Pattern:** Tries · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/add-and-search-word-data-structure-design)

Design a structure with `addWord(word)` and `search(word)`, where search words may contain `'.'` matching **any single letter**.

**Example 1:**

    Input: addWord("bad"); addWord("dad"); addWord("mad"); search("pad"); search("bad"); search(".ad"); search("b..")
    Output: false, true, true, true

**Constraints:**

- `1 <= word.length <= 25`
- At most 2 dots per search word; up to `10^4` calls

---

**Hints — try each one before reading on:**
1. '.' means: try every child at this position.
2. DFS through the trie with the word index.

**Target:** O(L) add; search worst-case O(26^d · L) with dots, O(L) without
