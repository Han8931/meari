---
created: "2026-06-07"
id: go-b-arith
source: imported:builtin-go
study:
  answer: |
    func Average3(a, b, c int) int {
    	return (a + b + c) / 3
    }
  kind: code
  lang: go
  prompt: Write Average3(a, b, c int) int returning the integer average of the three values (the / operator truncates, which is fine here).
  starter: |
    func Average3(a, b, c int) int {
    	return 0
    }
  tests:
    - if Average3(1, 2, 3) != 2 { t.Fatalf("got %d", Average3(1, 2, 3)) }
    - if Average3(10, 10, 10) != 10 { t.Fatal("equal values") }
    - 'if Average3(1, 1, 2) != 1 { t.Fatal("truncation: 4/3 = 1") }'
    - if Average3(-3, 0, 3) != 0 { t.Fatal("negatives") }
title: Arithmetic
---

The numeric operators are + - * / and % (remainder). On integers, / TRUNCATES
toward zero — the remainder is what % is for:
    7 / 2    // 3   (not 3.5)
    7 % 2    // 1
    17 / 5   // 3
    17 % 5   // 2

Precedence works as in math (* / % before + -); use parentheses to be explicit:
    (a + b) / 2

Compound assignment updates a variable in place, and ++/-- add or subtract one
(they are statements, not expressions):
    total := 0
    total += 5      // total = total + 5
    total++         // 6

Worked example — splitting minutes into hours and minutes:
    minutes := 130
    h := minutes / 60       // 2
    m := minutes % 60       // 10
