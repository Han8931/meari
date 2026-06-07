---
created: "2026-06-07"
id: nc-minimum-window-substring
source: imported:neetcode-150
study:
  answer: |-
    Classic expand-then-contract window: need = Counter(t); extend right updating have (number of characters whose count requirement is met); while have == need, record the window and shrink from the left. The best recorded window is the answer.

    Complexity: O(|s| + |t|) time, O(alphabet) space
  kind: essay
  prompt: 'Solve "Minimum Window Substring" (Sliding Window): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Minimum Window Substring
---

**Pattern:** Sliding Window · **Difficulty:** Hard · [LeetCode ↗](https://leetcode.com/problems/minimum-window-substring)

Given strings `s` and `t`, return the **smallest window** of `s` that contains every character of `t` (counting multiplicity), or `""` if none exists.

**Example 1:**

    Input: s = "ADOBECODEBANC", t = "ABC"
    Output: "BANC"

**Example 2:**

    Input: s = "a", t = "aa"
    Output: ""

**Constraints:**

- `1 <= s.length, t.length <= 10^5`
- `s` and `t` consist of English letters
- The answer is guaranteed unique

---

**Hints — try each one before reading on:**
1. Expand right until the window covers t, then shrink left while it still does.
2. Track have/need counts so coverage checks are O(1).

**Target:** O(|s| + |t|) time, O(alphabet) space
