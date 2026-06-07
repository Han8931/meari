---
created: "2026-06-07"
id: nc-lru-cache
source: imported:neetcode-150
study:
  answer: |-
    Hash map from key to a node in a doubly linked list ordered by recency, with sentinel head/tail: get/put unlink the node and reinsert at the MRU end; on overflow remove the LRU end's node and its map entry. OrderedDict with move_to_end/popitem(last=False) is the concise Python form.

    Complexity: O(1) per operation, O(capacity) space
  kind: essay
  prompt: 'Solve "LRU Cache" (Linked List): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: LRU Cache
---

**Pattern:** Linked List · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/lru-cache)

Design a **Least Recently Used cache**: `LRUCache(capacity)`, `get(key)` returning the value or -1, and `put(key, value)` which evicts the least recently used key when full. Both operations must run in O(1) average time.

**Example 1:**

    Input: LRUCache(2); put(1,1); put(2,2); get(1); put(3,3); get(2)
    Output: get(1) = 1, get(2) = -1
    Explanation: put(3,3) evicted key 2 (key 1 was touched more recently).

**Constraints:**

- `1 <= capacity <= 3000`
- Up to `2 * 10^5` calls

---

**Hints — try each one before reading on:**
1. O(1) lookup is a hash map; O(1) reordering is a doubly linked list.
2. Map key → node; move a touched node to the front; evict from the tail. (Python's OrderedDict packages exactly this.)

**Target:** O(1) per operation, O(capacity) space
