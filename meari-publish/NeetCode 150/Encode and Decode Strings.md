---
created: "2026-06-07"
id: nc-encode-and-decode-strings
source: imported:neetcode-150
study:
  answer: |-
    Encode each string as len + '#' + content. Decode by reading digits up to the '#', then consuming exactly that many characters; repeat. Length prefixes make the format self-describing.

    Complexity: O(n) time, O(1) extra space
  kind: essay
  prompt: 'Solve "Encode and Decode Strings" (Arrays & Hashing): describe the optimal approach — the key data structure or pattern and why it works — state time and space complexity, then write the Python solution.'
title: Encode and Decode Strings
---

**Pattern:** Arrays & Hashing · **Difficulty:** Medium · [LeetCode ↗](https://leetcode.com/problems/encode-and-decode-strings)

Design `encode(strs)` turning a list of strings into one string, and `decode(s)` turning it back into the original list. Any character — including your delimiter — may appear in the input strings.

**Example 1:**

    Input: ["neet","code","love","you"]
    Output (after encode + decode): ["neet","code","love","you"]

**Example 2:**

    Input: ["we","say",":","yes"]
    Output: ["we","say",":","yes"]

**Constraints:**

- `0 <= strs.length < 100`
- Strings may contain any of the 256 ASCII characters

---

**Hints — try each one before reading on:**
1. Delimiters fail because any character can appear in the data.
2. Length-prefix each string: "4#abcd" is unambiguous no matter the content.

**Target:** O(n) time, O(1) extra space
