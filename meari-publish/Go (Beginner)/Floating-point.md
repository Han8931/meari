---
created: "2026-06-07"
id: go-b-floats
source: imported:builtin-go
study:
  answer: |
    import "math"

    func NearlyEqual(a, b, tol float64) bool {
    	return math.Abs(a-b) <= tol
    }
  kind: code
  lang: go
  prompt: Write NearlyEqual(a, b, tol float64) bool that reports whether a and b differ by at most tol (use math.Abs).
  starter: |
    import "math"

    func NearlyEqual(a, b, tol float64) bool {
    	return false
    }
  tests:
    - if !NearlyEqual(0.1+0.2, 0.3, 1e-9) { t.Fatal("0.1+0.2 should be ~0.3") }
    - if NearlyEqual(1, 2, 0.1) { t.Fatal("1 and 2 are not close") }
    - if !NearlyEqual(5, 5, 0) { t.Fatal("equal values") }
title: Floating-point
---

Real numbers use floating-point types. The default is float64 (8 bytes);
float32 (4 bytes) trades precision for memory. A literal with a decimal point
is inferred as float64:
    pi := 3.14159        // float64

Floating-point can't represent every decimal exactly, so arithmetic accrues
tiny errors:
    fmt.Println(0.1 + 0.2)          // 0.30000000000000004
    fmt.Println(0.1+0.2 == 0.3)     // false

Never compare floats with ==. Instead check that the absolute difference is
within a small tolerance, using math.Abs (add import "math" above your
function to use it):
    math.Abs(a-b) < 1e-9

Printf's %f verb controls formatting: %8.3f means width 8, 3 digits after the
point.
