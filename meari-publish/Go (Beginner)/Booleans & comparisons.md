---
created: "2026-06-07"
id: go-b-bools
source: imported:builtin-go
study:
  answer: |
    func InRange(n, lo, hi int) bool {
    	return lo <= n && n <= hi
    }
  kind: code
  lang: go
  prompt: Write InRange(n, lo, hi int) bool reporting whether n is between lo and hi inclusive (use && with two comparisons).
  starter: |
    func InRange(n, lo, hi int) bool {
    	return false
    }
  tests:
    - if !InRange(5, 1, 10) { t.Fatal("5 is in 1..10") }
    - if !InRange(1, 1, 10) || !InRange(10, 1, 10) { t.Fatal("bounds are inclusive") }
    - if InRange(0, 1, 10) || InRange(11, 1, 10) { t.Fatal("outside") }
title: Booleans & comparisons
---

A bool is true or false. Comparisons produce bools:
    ==  !=  <  <=  >  >=
    age >= 18          // true or false
    name == "Ada"      // strings compare with == too

Combine bools with && (and), || (or), and ! (not):
    age >= 13 && age <= 19      // a teenager
    day == "Sat" || day == "Sun"
    !done

&& and || short-circuit: the right side is only evaluated when it can still
change the answer. That makes guards safe and idiomatic:
    n != 0 && total/n > 10      // never divides by zero

Worked example — is a year a leap year?
    leap := year%4 == 0 && (year%100 != 0 || year%400 == 0)
