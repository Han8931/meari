---
created: "2026-06-07"
id: nc-design-twitter
source: imported:neetcode-150
study:
  answer: |-
    Keep userId → list of (timestamp, tweetId) and userId → set of followees (including self). getNewsFeed seeds a heap with each followee's most recent tweet and pops 10, pushing each list's predecessor as it goes — the merge-k pattern.

    Complexity: O(f log f) per feed, O(users + tweets) space
  kind: essay
  prompt: 'Solve "Design Twitter" (Heap / Priority Queue): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Design Twitter
---

**Pattern:** Heap / Priority Queue · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/design-twitter)

Design a simplified Twitter: `postTweet(userId, tweetId)`, `getNewsFeed(userId)` returning the **10 most recent** tweet ids from the user and everyone they follow, `follow(a, b)`, and `unfollow(a, b)`.

**Example 1:**

    Input: postTweet(1,5); getNewsFeed(1); follow(1,2); postTweet(2,6); getNewsFeed(1); unfollow(1,2); getNewsFeed(1)
    Output: [5], [6,5], [5]

**Constraints:**

- Up to `3 * 10^4` calls
- A user cannot follow themselves (self-feed is implicit)

---

**Hints — try each one before reading on:**
1. Per-user tweet lists + follow sets cover everything but the feed.
2. The feed is "merge k sorted lists, take 10" — a heap over each followee's latest tweet.

**Target:** O(f log f) per feed, O(users + tweets) space
