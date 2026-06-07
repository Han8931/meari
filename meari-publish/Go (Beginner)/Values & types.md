---
created: "2026-06-07"
id: go-b-values
source: imported:builtin-go
study:
  answer: |
    func Greet(name string) string {
    	return "Hello, " + name + "!"
    }
  kind: code
  lang: go
  prompt: Write Greet(name string) string returning "Hello, <name>!" — for example Greet("Ada") is "Hello, Ada!". Join the pieces with the + operator.
  starter: |
    func Greet(name string) string {
    	return ""
    }
  tests:
    - if Greet("Ada") != "Hello, Ada!" { t.Fatalf("got %q", Greet("Ada")) }
    - if Greet("Go") != "Hello, Go!" { t.Fatal("Greet(\"Go\") should be \"Hello, Go!\"") }
title: Values & types
---

Before writing much code it helps to know Go's basic value types. Because Go is
statically typed, every value is one specific type, and the compiler won't
silently mix them.

The everyday built-in types are:
    int      whole numbers like 42 or -7    (also sized: int8, int32, int64…)
    float64  numbers with a fraction, 3.14
    string   text in double quotes, "hello"  (immutable)
    bool     true or false

You can compute directly with literal values:
    1 + 1          // 2          integer arithmetic
    7.0 / 2.0      // 3.5        floating-point division
    "go" + "pher"  // "gopher"   + joins two strings
    3 > 2          // true       a comparison yields a bool

Go won't add an int to a float64 for you — you convert one first (a later topic).
That strictness is on purpose: it keeps surprising bugs out of your programs.
