---
created: "2026-06-07"
id: go-b-ints
source: imported:builtin-go
study:
  answer: |
    func WrapUint8(n int) int {
    	return int(uint8(n))
    }
  kind: code
  lang: go
  prompt: Write WrapUint8(n int) int that converts n to a uint8 and back to int, demonstrating 8-bit wraparound (e.g. 256 -> 0, -1 -> 255).
  starter: |
    func WrapUint8(n int) int {
    	return 0
    }
  tests:
    - if WrapUint8(200) != 200 { t.Fatal("200 fits") }
    - if WrapUint8(256) != 0 { t.Fatalf("256 -> %d", WrapUint8(256)) }
    - if WrapUint8(257) != 1 { t.Fatal("257") }
    - if WrapUint8(-1) != 255 { t.Fatalf("-1 -> %d", WrapUint8(-1)) }
title: Integers & wraparound
---

Go has many integer types. Signed: int8, int16, int32, int64; unsigned:
uint8..uint64. int and uint are the machine word size (usually 64-bit). A whole
-number literal is inferred as int.

Sized integers have a fixed range, and exceeding it wraps around rather than
erroring:
    var b uint8 = 255
    b++                 // 0  (uint8 holds 0..255)

Types don't mix implicitly — you must convert explicitly with T(value):
    var i int = 300
    var b uint8 = uint8(i)   // 300 wraps to 44
    var back int = int(b)    // 44

Converting to a smaller type can silently lose information, so convert with
care (or range-check first).
