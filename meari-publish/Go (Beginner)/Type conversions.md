---
created: "2026-06-07"
id: go-b-conv
source: imported:builtin-go
study:
  answer: |
    import "strconv"

    func FizzBuzzAt(n int) string {
    	switch {
    	case n%15 == 0:
    		return "FizzBuzz"
    	case n%3 == 0:
    		return "Fizz"
    	case n%5 == 0:
    		return "Buzz"
    	default:
    		return strconv.Itoa(n)
    	}
    }
  kind: code
  lang: go
  prompt: 'Write FizzBuzzAt(n int) string: "Fizz" if n is divisible by 3, "Buzz" by 5, "FizzBuzz" by both, otherwise n itself as a string (strconv.Itoa).'
  starter: |
    import "strconv"

    func FizzBuzzAt(n int) string {
    	return ""
    }
  tests:
    - if FizzBuzzAt(3) != "Fizz" { t.Fatal("3") }
    - if FizzBuzzAt(5) != "Buzz" { t.Fatal("5") }
    - if FizzBuzzAt(15) != "FizzBuzz" { t.Fatal("15") }
    - if FizzBuzzAt(7) != "7" { t.Fatalf("7 -> %q", FizzBuzzAt(7)) }
title: Type conversions
---

Go never converts between types implicitly. Mixing types is a compile error;
you convert explicitly with T(value):
    var age int = 41
    var marsAge float64 = float64(age)

Numbers and strings need the strconv package, not a plain conversion
(string(65) gives the character "A", not "65"):
    s := strconv.Itoa(10)       // "10"

To use a package, import it above your functions:
    import "strconv"

(strconv.Atoi parses the other way, string to int — it returns a value AND an
error, a two-result pattern you'll learn with functions later.)

Worked example — labeling a count:
    msg := "items: " + strconv.Itoa(count)
