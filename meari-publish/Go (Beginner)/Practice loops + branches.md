---
created: "2026-06-07"
id: go-b-evens
source: imported:builtin-go
study:
  answer: |
    func CountEvens(from, to int) int {
    	count := 0
    	for i := from; i <= to; i++ {
    		if i%2 == 0 {
    			count++
    		}
    	}
    	return count
    }
  kind: code
  lang: go
  prompt: Write CountEvens(from, to int) int counting the even numbers between from and to inclusive (return 0 when from > to).
  starter: |
    func CountEvens(from, to int) int {
    	return 0
    }
  tests:
    - if CountEvens(1, 10) != 5 { t.Fatalf("1..10 -> %d", CountEvens(1, 10)) }
    - if CountEvens(2, 2) != 1 { t.Fatal("2..2") }
    - if CountEvens(3, 3) != 0 { t.Fatal("3..3") }
    - if CountEvens(5, 1) != 0 { t.Fatal("from > to") }
title: 'Practice: loops + branches'
---

Everything so far combines: a loop walks the numbers, a branch picks the ones
that matter, a variable accumulates the answer. This shape — loop, test,
accumulate — solves an enormous number of small problems.

Worked example — counting multiples of 3 from 1 to n:
    count := 0
    for i := 1; i <= n; i++ {
        if i%3 == 0 {
            count++
        }
    }

Variations on the same shape: sum instead of count, track the largest value
seen so far, or stop early with break when something is found. Try re-writing
the example as a sum in your head before doing the challenge.
