---
created: "2026-06-07"
id: go-b-runes
source: imported:builtin-go
study:
  answer: |
    import "unicode/utf8"

    func RuneCount(s string) int {
    	return utf8.RuneCountInString(s)
    }
  kind: code
  lang: go
  prompt: Write RuneCount(s string) int returning the number of Unicode characters (runes) in s — not the number of bytes.
  starter: |
    import "unicode/utf8"

    func RuneCount(s string) int {
    	return 0
    }
  tests:
    - if RuneCount("Hello") != 5 { t.Fatal("ascii") }
    - if RuneCount("Hello世界") != 7 { t.Fatalf("cjk -> %d", RuneCount("Hello世界")) }
    - if RuneCount("¿Cómo?") != 6 { t.Fatalf("accents -> %d", RuneCount("¿Cómo?")) }
title: Runes & UTF-8
---

Go source is UTF-8, and so are Go strings. A byte (alias for uint8) is one
8-bit unit; a rune (alias for int32) is one Unicode code point — a "character".
Non-ASCII characters take more than one byte, so byte length and character
count differ:
    s := "Héllo"
    len(s)                          // 6 bytes (é is 2 bytes)
    utf8.RuneCountInString(s)       // 5 runes

Indexing a string gives bytes. To work with characters, range over the string
(it decodes runes for you) or use the unicode/utf8 package:
    for i, r := range s {           // i is the byte offset, r is a rune
        fmt.Printf("%d %c\n", i, r)
    }
