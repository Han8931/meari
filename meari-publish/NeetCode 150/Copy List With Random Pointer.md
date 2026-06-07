---
created: "2026-06-07"
id: nc-copy-list-with-random-pointer
source: imported:neetcode-150
study:
  answer: |-
    Two passes with a hash map old → clone: first create all clones, then set clone.next = map[old.next] and clone.random = map[old.random]. (The O(1)-space variant interleaves clones into the original list, wires randoms, then unweaves.)

    Complexity: O(n) time, O(n) space (O(1) with interleaving)
  kind: essay
  prompt: 'Solve "Copy List With Random Pointer" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Copy List With Random Pointer
---

**Pattern:** Linked List · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/copy-list-with-random-pointer)

Each node of a linked list has an extra `random` pointer to any node (or null). Return a **deep copy** of the list.

**Example 1:**

    Input: head = [[7,null],[13,0],[11,4],[10,2],[1,0]]  (val, random index)
    Output: a structurally identical list of brand-new nodes

**Constraints:**

- List length is in `[0, 1000]`
- `random` is null or points into the list

---

**Hints — try each one before reading on:**
1. You need old node → new node to translate random pointers.
2. Pass 1 builds the map (or interleaves clones); pass 2 wires next and random through it.

**Target:** O(n) time, O(n) space (O(1) with interleaving)
