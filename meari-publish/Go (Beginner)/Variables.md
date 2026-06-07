---
created: "2026-06-07"
id: go-b-vars
source: imported:builtin-go
study:
  answer: |
    import "fmt"

    func Describe(name string, age int) string {
    	return fmt.Sprintf("%s is %d years old", name, age)
    }
  kind: code
  lang: go
  prompt: Write Describe(name string, age int) string returning e.g. "Ada is 36 years old". Use fmt.Sprintf.
  starter: |
    import "fmt"

    func Describe(name string, age int) string {
    	return ""
    }
  tests:
    - if Describe("Ada", 36) != "Ada is 36 years old" { t.Fatalf("got %q", Describe("Ada", 36)) }
    - if Describe("Sam", 7) != "Sam is 7 years old" { t.Fatal("Sam") }
title: Variables
---

A variable is named storage for a value. Go gives you three ways to declare one:
    var name string = "Ada"   // explicit type
    var age = 36              // type inferred from the value
    height := 1.7             // short form (only inside a function)

The := short form is the one you'll use most inside functions. Every type also
has a zero value, used when you declare without assigning: 0 for numbers, "" for
strings, false for bool. An unused local variable is a compile error — Go keeps
code tidy.

To build text from values, fmt.Printf prints using format verbs, and fmt.Sprintf
returns the result as a string instead of printing it:
    %v the value, %T its type, %d an integer, %s a string, %q a quoted string
    fmt.Sprintf("%s is %d", name, age)   // "Ada is 36"
