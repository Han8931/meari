---
created: "2026-06-07"
id: go-b-bignum
source: imported:builtin-go
study:
  answer: |
    import "math/big"

    func Power(base, exp int64) string {
    	r := new(big.Int).Exp(big.NewInt(base), big.NewInt(exp), nil)
    	return r.String()
    }
  kind: code
  lang: go
  prompt: Write Power(base, exp int64) string returning base raised to exp as a decimal string, using math/big so it works far beyond int64 (e.g. Power(2, 64)).
  starter: |
    import "math/big"

    func Power(base, exp int64) string {
    	return ""
    }
  tests:
    - if Power(2, 10) != "1024" { t.Fatalf("2^10 -> %s", Power(2, 10)) }
    - if Power(2, 64) != "18446744073709551616" { t.Fatalf("2^64 -> %s", Power(2, 64)) }
    - if Power(10, 20) != "100000000000000000000" { t.Fatal("10^20") }
title: Big numbers
---

int64 maxes out near 9.2 x 10^18. For larger exact integers, use the
math/big package, which grows until you run out of memory. big.Int values are
built and mutated through methods (calls written value.Method(...) — you'll
study methods properly later; here just follow the pattern).

Create them with big.NewInt(x) from an int64, or from a string for values too
large to write as a literal:
    a := big.NewInt(2)
    n := new(big.Int)
    n.Exp(a, big.NewInt(100), nil)   // 2^100
    fmt.Println(n.String())

new(big.Int) allocates a zero value and returns a pointer; big.NewInt both
allocates and initializes. The package also offers big.Rat (exact fractions)
and big.Float (arbitrary precision).
