---
created: "2026-06-07"
id: go-b-hello
source: imported:builtin-go
study:
  answer: |
    func Greeting() string {
    	return "Hello, Go!"
    }
  kind: code
  lang: go
  prompt: Write Greeting() string that returns the text "Hello, Go!".
  starter: |
    func Greeting() string {
    	return ""
    }
  tests:
    - if Greeting() != "Hello, Go!" { t.Fatalf("got %q, want \"Hello, Go!\"", Greeting()) }
title: Hello, Go
---

Welcome to Go! Go is a small, compiled, statically typed language made at Google
for building fast, reliable tools and servers. "Compiled" means your code is
turned into a standalone executable before it runs; "statically typed" means
every value has a fixed type the compiler checks for you, catching many mistakes
before the program even starts.

Every Go file belongs to a package. A runnable program lives in package main and
begins at the function main:
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, Go!")
    }

Normally you'd run that with  go run hello.go  (compile and run) or build a
binary with  go build. The fmt package handles formatted I/O — fmt.Println
prints a line of text.

In these lessons you don't write the package or main — that scaffolding is taken
care of. Instead the editor hands you one small function to complete:
    func Greeting() string {
        return ""
    }
Fill in the body so it returns the right value; hidden tests then call your
function and check what comes back. Click "Check answer" (or press Ctrl-S) to
run them.
