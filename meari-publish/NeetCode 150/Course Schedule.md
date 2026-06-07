---
created: "2026-06-07"
id: nc-course-schedule
source: imported:neetcode-150
study:
  answer: |-
    Build the adjacency list and detect a cycle: DFS marking nodes "visiting" on entry and "done" on exit — meeting a "visiting" node is a cycle. (Kahn's algorithm: repeatedly remove in-degree-0 nodes; leftovers mean a cycle.)

    Complexity: O(V + E) time, O(V + E) space
  kind: essay
  prompt: 'Solve "Course Schedule" (Graphs): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Course Schedule
---

**Pattern:** Graphs · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/course-schedule)

There are `numCourses` courses; `prerequisites[i] = [a, b]` means course `a` requires course `b` first. Return `true` if **all** courses can be finished.

**Example 1:**

    Input: numCourses = 2, prerequisites = [[1,0]]
    Output: true

**Example 2:**

    Input: numCourses = 2, prerequisites = [[1,0],[0,1]]
    Output: false
    Explanation: a circular dependency.

**Constraints:**

- `1 <= numCourses <= 2000`
- `0 <= prerequisites.length <= 5000`

---

**Hints — try each one before reading on:**
1. Possible iff the prerequisite digraph has no cycle.
2. DFS with three colors (unvisited / in-stack / done), or Kahn's in-degree BFS.

**Target:** O(V + E) time, O(V + E) space
