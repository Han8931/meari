---
created: "2026-06-07"
id: nc-course-schedule-ii
source: imported:neetcode-150
study:
  answer: |-
    Kahn's algorithm: compute in-degrees, queue all zero-in-degree courses, repeatedly pop one into the order and decrement its dependents (enqueueing new zeros). An order shorter than n means a cycle → [].

    Complexity: O(V + E) time, O(V + E) space
  kind: essay
  prompt: 'Solve "Course Schedule II" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Course Schedule II
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/course-schedule-ii)

Same setup as Course Schedule, but return a **valid order** in which to take all courses, or `[]` if impossible.

**Example 1:**

    Input: numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]
    Output: [0,1,2,3] (or [0,2,1,3])

**Example 2:**

    Input: numCourses = 1, prerequisites = []
    Output: [0]

**Constraints:**

- `1 <= numCourses <= 2000`
- All prerequisite pairs are distinct

---

**Hints — try each one before reading on:**
1. That's a topological sort.
2. Kahn's BFS emits the order directly; DFS post-order reversed also works.

**Target:** O(V + E) time, O(V + E) space
