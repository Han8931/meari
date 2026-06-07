---
created: "2026-06-07"
id: nc-implement-trie-prefix-tree
source: imported:neetcode-150
study:
  answer: |-
    Nodes hold a dict of child characters and an isWord flag. insert walks/creates child nodes and flags the last; search walks and requires the flag; startsWith walks and ignores it.

    Complexity: O(L) per operation, O(total characters) space
  kind: essay
  prompt: 'Solve "Implement Trie (Prefix Tree)" (Tries): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Implement Trie (Prefix Tree)
---

**Pattern:** Tries · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/implement-trie-prefix-tree)

Implement a **trie**: `insert(word)`, `search(word)` (was this exact word inserted?), and `startsWith(prefix)` (was any word with this prefix inserted?).

**Example 1:**

    Input: insert("apple"); search("apple"); search("app"); startsWith("app"); insert("app"); search("app")
    Output: true, false, true, true

**Constraints:**

- `1 <= word.length, prefix.length <= 2000`
- Lowercase English letters; up to `3 * 10^4` calls

---

**Hints — try each one before reading on:**
1. Each node: children map + end-of-word flag.
2. search and startsWith differ only in whether the final node must be flagged.

**Target:** O(L) per operation, O(total characters) space
