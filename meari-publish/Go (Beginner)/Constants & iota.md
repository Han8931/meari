---
created: "2026-06-07"
id: go-b-consts
source: imported:builtin-go
study:
  answer: |
    const (
    	Sun = iota
    	Mon
    	Tue
    	Wed
    	Thu
    	Fri
    	Sat
    )

    func IsWeekend(d int) bool {
    	return d == Sun || d == Sat
    }
  kind: code
  lang: go
  prompt: The weekday constants Sun..Sat (0..6) are defined for you with iota. Write IsWeekend(d int) bool returning true only for Sun or Sat.
  starter: |
    const (
    	Sun = iota
    	Mon
    	Tue
    	Wed
    	Thu
    	Fri
    	Sat
    )

    func IsWeekend(d int) bool {
    	return false
    }
  tests:
    - if Sat != 6 { t.Fatalf("Sat = %d, want 6 (iota counts from 0)", Sat) }
    - if !IsWeekend(Sun) || !IsWeekend(Sat) { t.Fatal("Sun and Sat are weekend days") }
    - if IsWeekend(Wed) { t.Fatal("Wed is a weekday") }
title: Constants & iota
---

Go's const declares a value fixed at compile time. A constant can be typed or
"untyped" — an untyped constant adapts to whatever numeric type the expression
needs, so the same 3 works as an int or a float64:
    const greeting = "hi"
    const pi = 3.14159
    const shift = 1 << 20      // 1048576

iota counts up automatically inside a const block — this is how Go writes enums.
Each line increments iota, starting from 0:
    const (
        Sun = iota   // 0
        Mon          // 1
        Tue          // 2
    )

You can skip or scale values, since iota is just the line's index in the block:
    const (
        _  = iota
        KB = 1 << (10 * iota)   // 1024
        MB                      // 1048576
    )
