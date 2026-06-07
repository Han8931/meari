---
created: "2026-06-07"
id: go-b-if
source: imported:builtin-go
study:
  answer: |
    func Grade(score int) string {
    	switch {
    	case score >= 90:
    		return "A"
    	case score >= 80:
    		return "B"
    	case score >= 70:
    		return "C"
    	default:
    		return "F"
    	}
    }
  kind: code
  lang: go
  prompt: Write Grade(score int) string returning "A" for 90+, "B" for 80+, "C" for 70+, and "F" below that (a condition switch works well).
  starter: |
    func Grade(score int) string {
    	return ""
    }
  tests:
    - if Grade(95) != "A" { t.Fatal("95") }
    - if Grade(90) != "A" { t.Fatal("90 is an A") }
    - if Grade(85) != "B" || Grade(71) != "C" { t.Fatal("middle grades") }
    - if Grade(69) != "F" { t.Fatal("69") }
title: 'Branches: if & switch'
---

Branching uses if/else — no parentheses around the condition, braces always
required:
    if score >= 50 {
        result = "pass"
    } else if score >= 40 {
        result = "retry"
    } else {
        result = "fail"
    }

if can begin with a short statement whose variables exist only in the branch:
    if r := n % 2; r == 0 { ... }

switch is a cleaner if/else-if chain. Cases don't fall through (no break
needed), and a bare "switch {" with condition cases reads top to bottom:
    switch {
    case score >= 90:
        grade = "A"
    case score >= 80:
        grade = "B"
    default:
        grade = "C"
    }

A switch can also match values directly: switch day { case "Sat", "Sun": ... }
