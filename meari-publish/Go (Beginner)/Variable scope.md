---
created: "2026-06-07"
id: go-b-scope
source: imported:builtin-go
study:
  answer: |
    func Classify(n int) string {
    	if r := n % 2; r == 0 {
    		return "even"
    	}
    	return "odd"
    }
  kind: code
  lang: go
  prompt: 'Write Classify(n int) string using an if with a short statement: return "even" when n is divisible by 2, otherwise "odd".'
  starter: |
    func Classify(n int) string {
    	return ""
    }
  tests:
    - if Classify(4) != "even" { t.Fatal("4") }
    - if Classify(7) != "odd" { t.Fatal("7") }
    - if Classify(0) != "even" { t.Fatal("0") }
title: Variable scope
---

A variable exists only within the block — the { } — where it's declared, and
in nested blocks. Leaving the block ends its life. This keeps names local and
state contained.

A short statement in if/for/switch creates variables scoped to that construct:
    if r := n % 10; r != 0 {
        fmt.Println(r)      // r lives only inside this if/else
    }

Declaring a name in an inner block can shadow an outer one — a common source of
bugs:
    x := 1
    {
        x := 2      // a different variable, shadows the outer x
        _ = x
    }
    // x is still 1 here

Prefer the smallest scope that works; it makes code easier to reason about.
