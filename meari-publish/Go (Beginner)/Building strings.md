---
created: "2026-06-07"
id: go-b-build
source: imported:builtin-go
study:
  answer: |
    func Repeat(word string, n int) string {
    	out := ""
    	for i := 0; i < n; i++ {
    		out += word
    	}
    	return out
    }
  kind: code
  lang: go
  prompt: Write Repeat(word string, n int) string returning word repeated n times using a loop ("" when n < 1).
  starter: |
    func Repeat(word string, n int) string {
    	return ""
    }
  tests:
    - if Repeat("ha", 3) != "hahaha" { t.Fatalf("got %q", Repeat("ha", 3)) }
    - if Repeat("x", 0) != "" { t.Fatal("n < 1") }
    - if Repeat("ab", 1) != "ab" { t.Fatal("once") }
title: Building strings
---

Join strings with +, and grow one in a loop with +=. Strings compare with ==
and order lexically with < and >.

    greeting := "Hello, " + name + "!"

    line := ""
    for i := 0; i < 3; i++ {
        line += "ab"
    }
    // line == "ababab"

An empty string "" is the zero value — the natural starting point for an
accumulator, just like 0 for a sum.

Worked example — a separated list without a trailing separator:
    out := ""
    for i := 1; i <= 3; i++ {
        if out != "" {
            out += "-"
        }
        out += "x"
    }
    // out == "x-x-x"
