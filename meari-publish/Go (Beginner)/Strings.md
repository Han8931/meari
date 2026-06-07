---
created: "2026-06-07"
id: go-b-strings
source: imported:builtin-go
study:
  answer: |
    import "strings"

    func Initials(name string) string {
    	out := ""
    	for _, w := range strings.Fields(name) {
    		out += strings.ToUpper(w[:1])
    	}
    	return out
    }
  kind: code
  lang: go
  prompt: Write Initials(name string) string returning the uppercased first letter of each space-separated word (e.g. "ada lovelace" -> "AL"). Use strings.Fields, slicing, and strings.ToUpper.
  starter: |
    import "strings"

    func Initials(name string) string {
    	return ""
    }
  tests:
    - if Initials("ada lovelace") != "AL" { t.Fatalf("got %q", Initials("ada lovelace")) }
    - if Initials("grace brewster hopper") != "GBH" { t.Fatal("GBH") }
    - if Initials("") != "" { t.Fatal("empty") }
title: Strings
---

A string is an immutable sequence of bytes, usually UTF-8 text. You can read a
byte by index but never assign to one:
    s := "shalom"
    fmt.Println(s[0])     // 115 (the byte 's')
    // s[0] = 'x'         // compile error: strings are immutable

Slicing takes a substring by byte positions — s[i:j] is bytes i up to (not
including) j; omit an end to mean "from the start" or "to the end":
    s[0:3]   // "sha"
    s[:1]    // "s"  (the first byte, as a string)
    s[3:]    // "lom"

len(s) returns the number of BYTES, not characters. Interpreted literals use
double quotes and honor escapes (\n, \t); raw literals use back-quotes and
take the bytes verbatim across multiple lines.

The strings package has the everyday helpers: Fields (split on whitespace),
Split, Join, ToUpper/ToLower, Contains, HasPrefix, ReplaceAll. Fields returns
a collection you can walk with range:
    for _, w := range strings.Fields("ada lovelace") {
        fmt.Println(w)        // "ada", then "lovelace"
    }
