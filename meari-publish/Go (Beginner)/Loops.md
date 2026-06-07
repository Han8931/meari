---
created: "2026-06-07"
id: go-b-loops
source: imported:builtin-go
study:
  answer: |
    func SumTo(n int) int {
    	total := 0
    	for i := 1; i <= n; i++ {
    		total += i
    	}
    	return total
    }
  kind: code
  lang: go
  prompt: Write SumTo(n int) int returning 1 + 2 + ... + n computed with a for loop (return 0 when n < 1).
  starter: |
    func SumTo(n int) int {
    	return 0
    }
  tests:
    - if SumTo(5) != 15 { t.Fatalf("SumTo(5) -> %d", SumTo(5)) }
    - if SumTo(1) != 1 { t.Fatal("SumTo(1)") }
    - if SumTo(0) != 0 { t.Fatal("n < 1 should give 0") }
    - if SumTo(100) != 5050 { t.Fatal("SumTo(100)") }
title: Loops
---

Go has one loop keyword — for — covering every case:
    for i := 0; i < n; i++ { ... }   // classic counter
    for cond { ... }                 // like while
    for { ... }                      // infinite (use break to leave)

A fourth form, "for ... range", walks a sequence and hands you each position and
value — useful for text now, and for slices and maps later:
    for i, c := range "hi" {         // i = 0,1   c = each character
        fmt.Println(i, c)
    }

break leaves the loop; continue skips to the next iteration:
    for i := 1; i <= 100; i++ {
        if i%2 == 1 {
            continue        // skip odd numbers
        }
        if i > 40 {
            break           // stop entirely
        }
    }

Worked example — multiplying the numbers 1..n:
    fact := 1
    for i := 2; i <= n; i++ {
        fact *= i
    }
