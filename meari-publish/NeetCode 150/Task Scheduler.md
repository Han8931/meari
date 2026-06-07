---
created: "2026-06-07"
id: nc-task-scheduler
source: imported:neetcode-150
study:
  answer: |-
    Counting argument: with maxFreq the highest count and c tasks tied at it, answer = max(len(tasks), (maxFreq − 1) · (n + 1) + c). The greedy max-heap + cooldown-queue simulation derives the same number.

    Complexity: O(n) time (counting), O(1) space beyond counts
  kind: essay
  prompt: 'Solve "Task Scheduler" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Task Scheduler
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/task-scheduler)

Given CPU `tasks` (letters) and a cooldown `n` — identical tasks must be at least `n` ticks apart — return the **minimum** number of ticks to run everything (idle ticks count).

**Example 1:**

    Input: tasks = ["A","A","A","B","B","B"], n = 2
    Output: 8
    Explanation: A → B → idle → A → B → idle → A → B.

**Example 2:**

    Input: tasks = ["A","A","A","B","B","B"], n = 0
    Output: 6

**Constraints:**

- `1 <= tasks.length <= 10^4`
- `0 <= n <= 100`

---

**Hints — try each one before reading on:**
1. The most frequent task dictates the frame: (maxFreq − 1) blocks of size n+1.
2. Idle slots fill with other tasks; the answer can't be less than len(tasks).

**Target:** O(n) time (counting), O(1) space beyond counts
