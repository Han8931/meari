package curriculum

// Pre-authored Go curriculum, deepened to mirror a full introductory Go course
// (imperative basics & types → functions & collections → state, behavior &
// pointers). Challenge StarterCode/Solution are spliced into "package sol", and
// Tests run inside a generated test function where reflect and fmt are
// available. Every Solution is verified against its Tests in curriculum_test.go.

func goBeginner() []Module {
	return []Module{
		{
			Name: "Imperative programming",
			Topics: []Topic{
				{
					ID:    "go-b-vars",
					Title: "Print & variables",
					Lesson: `Go is compiled and statically typed. A program is a package; execution starts
at func main in package main. You print with the fmt package.

Declare variables three ways:
    var name string = "Ada"   // explicit type
    var age = 36              // type inferred from the value
    height := 1.7            // short declaration (only inside functions)

Every type has a zero value used when you declare without assigning: 0 for
numbers, "" for strings, false for bool. An unused local variable is a compile
error — Go keeps code clean.

fmt.Println prints values separated by spaces; fmt.Printf uses format verbs:
    %v the value, %T its type, %d an integer, %s a string, %q a quoted string.
    fmt.Printf("%s is %d\n", name, age)   // Ada is 36
fmt.Sprintf returns the formatted string instead of printing it.`,
					Challenge: Challenge{
						Prompt:      "Write Describe(name string, age int) string returning e.g. \"Ada is 36 years old\". Use fmt.Sprintf.",
						StarterCode: "func Describe(name string, age int) string {\n\treturn \"\"\n}\n",
						Solution:    "import \"fmt\"\n\nfunc Describe(name string, age int) string {\n\treturn fmt.Sprintf(\"%s is %d years old\", name, age)\n}\n",
						Tests: []string{
							`if Describe("Ada", 36) != "Ada is 36 years old" { t.Fatalf("got %q", Describe("Ada", 36)) }`,
							`if Describe("Sam", 7) != "Sam is 7 years old" { t.Fatal("Sam") }`,
						},
					},
				},
				{
					ID:    "go-b-loops",
					Title: "Loops & branches",
					Lesson: `Go has one loop keyword — for — covering every case:
    for i := 0; i < n; i++ { ... }   // classic
    for cond { ... }                 // like while
    for { ... }                      // infinite (use break)
    for i, v := range xs { ... }     // over a collection

Branching uses if/else (no parentheses, braces required). if can begin with a
short statement scoped to the branch:
    if r := n % 2; r == 0 { ... }

switch needs no break (no fall-through by default); a "switch true" with bare
case conditions reads like an if/elif chain:
    switch {
    case score >= 90: grade = "A"
    case score >= 80: grade = "B"
    default:          grade = "C"
    }`,
					Challenge: Challenge{
						Prompt:      "Write FizzBuzz(n int) []string for 1..n: \"Fizz\" if divisible by 3, \"Buzz\" by 5, \"FizzBuzz\" by both, otherwise the number as a string. Return an empty slice for n < 1.",
						StarterCode: "func FizzBuzz(n int) []string {\n\treturn nil\n}\n",
						Solution:    "import \"strconv\"\n\nfunc FizzBuzz(n int) []string {\n\tout := []string{}\n\tfor i := 1; i <= n; i++ {\n\t\tswitch {\n\t\tcase i%15 == 0:\n\t\t\tout = append(out, \"FizzBuzz\")\n\t\tcase i%3 == 0:\n\t\t\tout = append(out, \"Fizz\")\n\t\tcase i%5 == 0:\n\t\t\tout = append(out, \"Buzz\")\n\t\tdefault:\n\t\t\tout = append(out, strconv.Itoa(i))\n\t\t}\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`got := FizzBuzz(5); if !reflect.DeepEqual(got, []string{"1", "2", "Fizz", "4", "Buzz"}) { t.Fatalf("got %v", got) }`,
							`got := FizzBuzz(15); if got[14] != "FizzBuzz" { t.Fatalf("15 -> %q", got[14]) }`,
							`if !reflect.DeepEqual(FizzBuzz(0), []string{}) { t.Fatal("n<1 should be empty") }`,
						},
					},
				},
				{
					ID:    "go-b-scope",
					Title: "Variable scope",
					Lesson: `A variable exists only within the block — the { } — where it's declared, and
in nested blocks. Leaving the block ends its life. This keeps names local and
state contained.

A short statement in if/for/switch creates variables scoped to that construct:
    if v, err := strconv.Atoi(s); err == nil {
        use(v)        // v and err live only here
    }

Declaring a name in an inner block can shadow an outer one — a common source of
bugs:
    x := 1
    {
        x := 2      // a different variable, shadows the outer x
        _ = x
    }
    // x is still 1 here

Prefer the smallest scope that works; it makes code easier to reason about.`,
					Challenge: Challenge{
						Prompt:      "Write Classify(n int) string using an if with a short statement: return \"even\" when n is divisible by 2, otherwise \"odd\".",
						StarterCode: "func Classify(n int) string {\n\treturn \"\"\n}\n",
						Solution:    "func Classify(n int) string {\n\tif r := n % 2; r == 0 {\n\t\treturn \"even\"\n\t}\n\treturn \"odd\"\n}\n",
						Tests: []string{
							`if Classify(4) != "even" { t.Fatal("4") }`,
							`if Classify(7) != "odd" { t.Fatal("7") }`,
							`if Classify(0) != "even" { t.Fatal("0") }`,
						},
					},
				},
			},
		},
		{
			Name: "Numbers",
			Topics: []Topic{
				{
					ID:    "go-b-floats",
					Title: "Floating-point",
					Lesson: `Real numbers use floating-point types. The default is float64 (8 bytes);
float32 (4 bytes) trades precision for memory. A literal with a decimal point
is inferred as float64:
    pi := 3.14159        // float64

Floating-point can't represent every decimal exactly, so arithmetic accrues
tiny errors:
    fmt.Println(0.1 + 0.2)          // 0.30000000000000004
    fmt.Println(0.1+0.2 == 0.3)     // false

Never compare floats with ==. Instead check that the absolute difference is
within a small tolerance, using math.Abs:
    math.Abs(a-b) < 1e-9

Printf's %f verb controls formatting: %8.3f means width 8, 3 digits after the
point.`,
					Challenge: Challenge{
						Prompt:      "Write NearlyEqual(a, b, tol float64) bool that reports whether a and b differ by at most tol (use math.Abs).",
						StarterCode: "func NearlyEqual(a, b, tol float64) bool {\n\treturn false\n}\n",
						Solution:    "import \"math\"\n\nfunc NearlyEqual(a, b, tol float64) bool {\n\treturn math.Abs(a-b) <= tol\n}\n",
						Tests: []string{
							`if !NearlyEqual(0.1+0.2, 0.3, 1e-9) { t.Fatal("0.1+0.2 should be ~0.3") }`,
							`if NearlyEqual(1, 2, 0.1) { t.Fatal("1 and 2 are not close") }`,
							`if !NearlyEqual(5, 5, 0) { t.Fatal("equal values") }`,
						},
					},
				},
				{
					ID:    "go-b-ints",
					Title: "Integers & wraparound",
					Lesson: `Go has many integer types. Signed: int8, int16, int32, int64; unsigned:
uint8..uint64. int and uint are the machine word size (usually 64-bit). A whole
-number literal is inferred as int.

Sized integers have a fixed range, and exceeding it wraps around rather than
erroring:
    var b uint8 = 255
    b++                 // 0  (uint8 holds 0..255)

Types don't mix implicitly — you must convert explicitly:
    var i int = 300
    var b uint8 = uint8(i)   // 300 wraps to 44

Converting to a smaller type can silently lose information, so convert with
care (or range-check first).`,
					Challenge: Challenge{
						Prompt:      "Write WrapUint8(n int) int that converts n to a uint8 and back to int, demonstrating 8-bit wraparound (e.g. 256 -> 0, -1 -> 255).",
						StarterCode: "func WrapUint8(n int) int {\n\treturn 0\n}\n",
						Solution:    "func WrapUint8(n int) int {\n\treturn int(uint8(n))\n}\n",
						Tests: []string{
							`if WrapUint8(200) != 200 { t.Fatal("200 fits") }`,
							`if WrapUint8(256) != 0 { t.Fatalf("256 -> %d", WrapUint8(256)) }`,
							`if WrapUint8(257) != 1 { t.Fatal("257") }`,
							`if WrapUint8(-1) != 255 { t.Fatalf("-1 -> %d", WrapUint8(-1)) }`,
						},
					},
				},
				{
					ID:    "go-b-bignum",
					Title: "Big numbers",
					Lesson: `int64 maxes out near 9.2 x 10^18. For larger exact integers, use the
math/big package, which grows until you run out of memory. big.Int values are
built and mutated through methods.

Create them with big.NewInt(x) from an int64, or from a string for values too
large to write as a literal:
    a := big.NewInt(2)
    n := new(big.Int)
    n.Exp(a, big.NewInt(100), nil)   // 2^100
    fmt.Println(n.String())

new(big.Int) allocates a zero value and returns a pointer; big.NewInt both
allocates and initializes. The package also offers big.Rat (exact fractions)
and big.Float (arbitrary precision).`,
					Challenge: Challenge{
						Prompt:      "Write Power(base, exp int64) string returning base raised to exp as a decimal string, using math/big so it works far beyond int64 (e.g. Power(2, 64)).",
						StarterCode: "func Power(base, exp int64) string {\n\treturn \"\"\n}\n",
						Solution:    "import \"math/big\"\n\nfunc Power(base, exp int64) string {\n\tr := new(big.Int).Exp(big.NewInt(base), big.NewInt(exp), nil)\n\treturn r.String()\n}\n",
						Tests: []string{
							`if Power(2, 10) != "1024" { t.Fatalf("2^10 -> %s", Power(2, 10)) }`,
							`if Power(2, 64) != "18446744073709551616" { t.Fatalf("2^64 -> %s", Power(2, 64)) }`,
							`if Power(10, 20) != "100000000000000000000" { t.Fatal("10^20") }`,
						},
					},
				},
			},
		},
		{
			Name: "Text",
			Topics: []Topic{
				{
					ID:    "go-b-strings",
					Title: "Strings",
					Lesson: `A string is an immutable sequence of bytes, usually UTF-8 text. You can read a
byte by index but never assign to one:
    s := "shalom"
    fmt.Println(s[0])     // 115 (the byte 's')
    // s[0] = 'x'         // compile error: strings are immutable

len(s) returns the number of BYTES, not characters. Join strings with +.
Interpreted literals use double quotes and honor escapes (\n, \t); raw literals
use back-quotes and take the bytes verbatim across multiple lines (handy for
paths and templates).

The strings package has the everyday helpers: Fields (split on whitespace),
Split, Join, ToUpper/ToLower, Contains, HasPrefix, ReplaceAll.`,
					Challenge: Challenge{
						Prompt:      "Write Initials(name string) string returning the uppercased first letter of each space-separated word (e.g. \"ada lovelace\" -> \"AL\"). Use the strings package.",
						StarterCode: "func Initials(name string) string {\n\treturn \"\"\n}\n",
						Solution:    "import \"strings\"\n\nfunc Initials(name string) string {\n\tout := \"\"\n\tfor _, w := range strings.Fields(name) {\n\t\tout += strings.ToUpper(w[:1])\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`if Initials("ada lovelace") != "AL" { t.Fatalf("got %q", Initials("ada lovelace")) }`,
							`if Initials("grace brewster hopper") != "GBH" { t.Fatal("GBH") }`,
							`if Initials("") != "" { t.Fatal("empty") }`,
						},
					},
				},
				{
					ID:    "go-b-runes",
					Title: "Runes & UTF-8",
					Lesson: `Go source is UTF-8, and so are Go strings. A byte (alias for uint8) is one
8-bit unit; a rune (alias for int32) is one Unicode code point — a "character".
Non-ASCII characters take more than one byte, so byte length and character
count differ:
    s := "Héllo"
    len(s)                          // 6 bytes (é is 2 bytes)
    utf8.RuneCountInString(s)       // 5 runes

Indexing a string gives bytes. To work with characters, range over the string
(it decodes runes for you) or use the unicode/utf8 package:
    for i, r := range s {           // i is the byte offset, r is a rune
        fmt.Printf("%d %c\n", i, r)
    }`,
					Challenge: Challenge{
						Prompt:      "Write RuneCount(s string) int returning the number of Unicode characters (runes) in s — not the number of bytes.",
						StarterCode: "func RuneCount(s string) int {\n\treturn 0\n}\n",
						Solution:    "import \"unicode/utf8\"\n\nfunc RuneCount(s string) int {\n\treturn utf8.RuneCountInString(s)\n}\n",
						Tests: []string{
							`if RuneCount("Hello") != 5 { t.Fatal("ascii") }`,
							`if RuneCount("Hello世界") != 7 { t.Fatalf("cjk -> %d", RuneCount("Hello世界")) }`,
							`if RuneCount("¿Cómo?") != 6 { t.Fatalf("accents -> %d", RuneCount("¿Cómo?")) }`,
						},
					},
				},
				{
					ID:    "go-b-conv",
					Title: "Type conversions",
					Lesson: `Go never converts between types implicitly. Mixing types is a compile error;
you convert explicitly with T(value):
    var age int = 41
    var marsAge float64 = float64(age)

Strings need the strconv package, not a plain conversion (string(65) gives the
character "A", not "65"):
    s := strconv.Itoa(10)            // "10"
    n, err := strconv.Atoi("10")     // 10, nil  (err != nil if not a number)

Because parsing can fail, Atoi returns a value AND an error. The idiom is to
check err immediately:
    n, err := strconv.Atoi(s)
    if err != nil { /* handle bad input */ }`,
					Challenge: Challenge{
						Prompt:      "Write SafeAtoiSum(a, b string) (int, error) that parses both strings as integers and returns their sum, or a non-nil error if either fails to parse.",
						StarterCode: "func SafeAtoiSum(a, b string) (int, error) {\n\treturn 0, nil\n}\n",
						Solution:    "import \"strconv\"\n\nfunc SafeAtoiSum(a, b string) (int, error) {\n\tx, err := strconv.Atoi(a)\n\tif err != nil {\n\t\treturn 0, err\n\t}\n\ty, err := strconv.Atoi(b)\n\tif err != nil {\n\t\treturn 0, err\n\t}\n\treturn x + y, nil\n}\n",
						Tests: []string{
							`got, err := SafeAtoiSum("2", "3"); if err != nil || got != 5 { t.Fatalf("2+3: %d %v", got, err) }`,
							`if _, err := SafeAtoiSum("x", "3"); err == nil { t.Fatal("expected a parse error") }`,
						},
					},
				},
			},
		},
	}
}

func goIntermediate() []Module {
	return []Module{
		{
			Name: "Functions",
			Topics: []Topic{
				{
					ID:    "go-i-functions",
					Title: "Functions & multiple returns",
					Lesson: `A function lists its parameters and return types. Parameters of the same type
can share it:
    func add(a, b int) int { return a + b }

Go functions can return more than one value — used constantly for (result,
error) and for returning related values together:
    func divmod(a, b int) (int, int) {
        return a / b, a % b
    }
    q, r := divmod(17, 5)   // 3, 2

You can name the return values; a bare "return" then returns their current
values. Use named returns sparingly, where they aid clarity.`,
					Challenge: Challenge{
						Prompt:      "Write DivMod(a, b int) (int, int) returning the quotient and remainder of a divided by b.",
						StarterCode: "func DivMod(a, b int) (int, int) {\n\treturn 0, 0\n}\n",
						Solution:    "func DivMod(a, b int) (int, int) {\n\treturn a / b, a % b\n}\n",
						Tests: []string{
							`q, r := DivMod(17, 5); if q != 3 || r != 2 { t.Fatalf("17/5 -> %d %d", q, r) }`,
							`q, r := DivMod(10, 2); if q != 5 || r != 0 { t.Fatal("10/2") }`,
						},
					},
				},
				{
					ID:    "go-i-methods",
					Title: "Methods & defined types",
					Lesson: `You can define your own named types on top of existing ones, then attach
methods to them. A method is a function with a receiver written before its name:
    type Celsius float64

    func (c Celsius) Fahrenheit() float64 {
        return float64(c)*9/5 + 32
    }

    t := Celsius(100)
    fmt.Println(t.Fahrenheit())   // 212

The receiver (c here) is the value the method is called on. A defined type is
distinct from its underlying type, so Celsius and float64 won't mix by
accident — the type system encodes meaning.`,
					Challenge: Challenge{
						Prompt:      "Given the defined type Celsius (int), implement the method Fahrenheit() int returning c*9/5 + 32.",
						StarterCode: "type Celsius int\n\nfunc (c Celsius) Fahrenheit() int {\n\treturn 0\n}\n",
						Solution:    "type Celsius int\n\nfunc (c Celsius) Fahrenheit() int {\n\treturn int(c)*9/5 + 32\n}\n",
						Tests: []string{
							`if Celsius(100).Fahrenheit() != 212 { t.Fatalf("100C -> %d", Celsius(100).Fahrenheit()) }`,
							`if Celsius(0).Fahrenheit() != 32 { t.Fatal("0C") }`,
							`if Celsius(-40).Fahrenheit() != -40 { t.Fatal("-40C") }`,
						},
					},
				},
				{
					ID:    "go-i-firstclass",
					Title: "First-class functions",
					Lesson: `Functions in Go are values: you can store them in variables, pass them as
arguments, and return them. This enables flexible, reusable code.
    op := func(a, b int) int { return a + b }
    fmt.Println(op(2, 3))   // 5

A function that takes another function can apply behavior supplied by the
caller:
    func apply(xs []int, f func(int) int) []int {
        out := make([]int, 0, len(xs))
        for _, x := range xs {
            out = append(out, f(x))
        }
        return out
    }

You can name a function type (type BinOp func(int, int) int) to make signatures
read clearly.`,
					Challenge: Challenge{
						Prompt:      "Write Apply(nums []int, f func(int) int) []int returning a new slice with f applied to each element, in order.",
						StarterCode: "func Apply(nums []int, f func(int) int) []int {\n\treturn nil\n}\n",
						Solution:    "func Apply(nums []int, f func(int) int) []int {\n\tout := make([]int, 0, len(nums))\n\tfor _, n := range nums {\n\t\tout = append(out, f(n))\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`got := Apply([]int{1, 2, 3}, func(x int) int { return x * x }); if !reflect.DeepEqual(got, []int{1, 4, 9}) { t.Fatalf("got %v", got) }`,
							`got := Apply([]int{}, func(x int) int { return x }); if !reflect.DeepEqual(got, []int{}) { t.Fatal("empty") }`,
						},
					},
				},
				{
					ID:    "go-i-closures",
					Title: "Closures",
					Lesson: `An anonymous function can capture variables from the scope where it's defined
— that combination is a closure. The captured variables persist between calls,
giving the function private, mutable state:
    func counter() func() int {
        n := 0
        return func() int {
            n++          // captures n
            return n
        }
    }
    c := counter()
    fmt.Println(c(), c())   // 1 2

Each call to the outer function creates fresh captured state, so two closures
are independent.`,
					Challenge: Challenge{
						Prompt:      "Write MakeAdder(start int) func(int) int returning a closure that keeps a running total beginning at start: each call adds its argument and returns the new total.",
						StarterCode: "func MakeAdder(start int) func(int) int {\n\treturn nil\n}\n",
						Solution:    "func MakeAdder(start int) func(int) int {\n\ttotal := start\n\treturn func(x int) int {\n\t\ttotal += x\n\t\treturn total\n\t}\n}\n",
						Tests: []string{
							`a := MakeAdder(10); if a(5) != 15 || a(3) != 18 { t.Fatal("running total") }`,
							`a := MakeAdder(0); b := MakeAdder(0); a(100); if b(1) != 1 { t.Fatal("closures must be independent") }`,
						},
					},
				},
			},
		},
		{
			Name: "Collections",
			Topics: []Topic{
				{
					ID:    "go-i-arrays",
					Title: "Arrays",
					Lesson: `An array has a fixed length that is part of its type: [3]int and [4]int are
different types. Index from 0; out-of-range access is a compile error (constant
index) or a runtime panic.
    var a [3]int
    a[0] = 10
    b := [3]int{1, 2, 3}     // composite literal

Arrays are values: assigning one or passing it to a function copies every
element. Because the length is fixed, arrays are less common in everyday Go
than slices (next topic), which are built on top of them.`,
					Challenge: Challenge{
						Prompt:      "Write Sum3(a [3]int) int returning the sum of the three elements.",
						StarterCode: "func Sum3(a [3]int) int {\n\treturn 0\n}\n",
						Solution:    "func Sum3(a [3]int) int {\n\treturn a[0] + a[1] + a[2]\n}\n",
						Tests: []string{
							`got := Sum3([3]int{1, 2, 3}); if got != 6 { t.Fatalf("got %d", got) }`,
							`got := Sum3([3]int{0, 0, 0}); if got != 0 { t.Fatal("zeros") }`,
						},
					},
				},
				{
					ID:    "go-i-slices",
					Title: "Slices",
					Lesson: `A slice is a flexible view into an underlying array: it holds a pointer, a
length, and a capacity. Slice an array or another slice with [low:high] (high
is exclusive):
    xs := []int{10, 20, 30, 40}
    mid := xs[1:3]            // [20 30]

A slice literal []int{...} (no length) creates the backing array for you.
Because a slice points into shared memory, two slices of the same array see
each other's changes. The empty slice []int{} has length 0.`,
					Challenge: Challenge{
						Prompt:      "Write Tail(s []int) []int returning all elements except the first. Return an empty slice when s has 0 or 1 elements.",
						StarterCode: "func Tail(s []int) []int {\n\treturn nil\n}\n",
						Solution:    "func Tail(s []int) []int {\n\tif len(s) <= 1 {\n\t\treturn []int{}\n\t}\n\treturn s[1:]\n}\n",
						Tests: []string{
							`if !reflect.DeepEqual(Tail([]int{1, 2, 3}), []int{2, 3}) { t.Fatalf("got %v", Tail([]int{1, 2, 3})) }`,
							`if !reflect.DeepEqual(Tail([]int{5}), []int{}) { t.Fatal("single") }`,
							`if !reflect.DeepEqual(Tail([]int{}), []int{}) { t.Fatal("empty") }`,
						},
					},
				},
				{
					ID:    "go-i-append",
					Title: "append, len & cap",
					Lesson: `Grow a slice with append, which returns a (possibly newly allocated) slice —
always assign the result back:
    s := []int{}
    s = append(s, 1, 2, 3)

len(s) is the number of elements; cap(s) is how many the backing array can hold
before append must allocate a bigger one. When you know the size, preallocate
to avoid repeated growth:
    s := make([]int, 0, 100)   // len 0, cap 100

Building a result slice by appending in a loop is the standard Go pattern.`,
					Challenge: Challenge{
						Prompt:      "Write Filter(s []int, keep func(int) bool) []int returning a new slice of the elements for which keep returns true (empty slice if none).",
						StarterCode: "func Filter(s []int, keep func(int) bool) []int {\n\treturn nil\n}\n",
						Solution:    "func Filter(s []int, keep func(int) bool) []int {\n\tout := []int{}\n\tfor _, n := range s {\n\t\tif keep(n) {\n\t\t\tout = append(out, n)\n\t\t}\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`got := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 }); if !reflect.DeepEqual(got, []int{2, 4}) { t.Fatalf("got %v", got) }`,
							`got := Filter([]int{1, 3}, func(n int) bool { return n%2 == 0 }); if !reflect.DeepEqual(got, []int{}) { t.Fatal("none") }`,
						},
					},
				},
				{
					ID:    "go-i-variadic",
					Title: "Variadic functions",
					Lesson: `A variadic function accepts any number of trailing arguments of a type,
written with ... before the type. Inside, the parameter is a slice:
    func sum(nums ...int) int {
        total := 0
        for _, n := range nums {
            total += n
        }
        return total
    }
    sum(1, 2, 3)   // 6
    sum()          // 0

You can pass an existing slice to a variadic function by "spreading" it with
...:
    xs := []int{1, 2, 3}
    sum(xs...)`,
					Challenge: Challenge{
						Prompt:      "Write Max(nums ...int) int returning the largest argument, or 0 when called with no arguments.",
						StarterCode: "func Max(nums ...int) int {\n\treturn 0\n}\n",
						Solution:    "func Max(nums ...int) int {\n\tif len(nums) == 0 {\n\t\treturn 0\n\t}\n\tm := nums[0]\n\tfor _, n := range nums[1:] {\n\t\tif n > m {\n\t\t\tm = n\n\t\t}\n\t}\n\treturn m\n}\n",
						Tests: []string{
							`if Max(3, 1, 2) != 3 { t.Fatal("3,1,2") }`,
							`if Max() != 0 { t.Fatal("none -> 0") }`,
							`if Max(-5, -2, -9) != -2 { t.Fatalf("negatives -> %d", Max(-5, -2, -9)) }`,
						},
					},
				},
				{
					ID:    "go-i-maps",
					Title: "Maps & sets",
					Lesson: `A map stores key/value pairs with fast lookup: map[K]V. Reading a missing key
returns the value type's zero value, which makes counting trivial:
    counts := map[string]int{}
    counts["a"]++            // missing key starts at 0 -> becomes 1

Use the comma-ok form to tell "absent" from "zero":
    v, ok := counts["b"]     // ok is false if "b" isn't present

A map[T]bool is the idiomatic set: store true for members and test with the
comma-ok form. Maps are reference types — assigning one doesn't copy it.`,
					Challenge: Challenge{
						Prompt:      "Write Unique(s []int) []int returning the values of s with duplicates removed, preserving first-seen order. Use a set (map[int]bool).",
						StarterCode: "func Unique(s []int) []int {\n\treturn nil\n}\n",
						Solution:    "func Unique(s []int) []int {\n\tseen := map[int]bool{}\n\tout := []int{}\n\tfor _, n := range s {\n\t\tif !seen[n] {\n\t\t\tseen[n] = true\n\t\t\tout = append(out, n)\n\t\t}\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`got := Unique([]int{1, 2, 2, 3, 1}); if !reflect.DeepEqual(got, []int{1, 2, 3}) { t.Fatalf("got %v", got) }`,
							`if !reflect.DeepEqual(Unique([]int{}), []int{}) { t.Fatal("empty") }`,
						},
					},
				},
			},
		},
	}
}

func goAdvanced() []Module {
	return []Module{
		{
			Name: "State & behavior",
			Topics: []Topic{
				{
					ID:    "go-a-structs",
					Title: "Structs",
					Lesson: `A struct groups named fields into a single value type:
    type Rect struct {
        W, H int
    }
    r := Rect{W: 2, H: 3}    // or Rect{2, 3} positionally
    fmt.Println(r.W)         // 2

Structs are values: assigning a struct or passing it to a function copies every
field. A slice of structs ([]Rect) is a common way to hold many records, and
you can build one with a composite literal:
    rs := []Rect{{1, 2}, {3, 4}}`,
					Challenge: Challenge{
						Prompt:      "Given the Rect struct, write TotalWidth(rs []Rect) int returning the sum of every rectangle's W.",
						StarterCode: "type Rect struct {\n\tW, H int\n}\n\nfunc TotalWidth(rs []Rect) int {\n\treturn 0\n}\n",
						Solution:    "type Rect struct {\n\tW, H int\n}\n\nfunc TotalWidth(rs []Rect) int {\n\ttotal := 0\n\tfor _, r := range rs {\n\t\ttotal += r.W\n\t}\n\treturn total\n}\n",
						Tests: []string{
							`rs := []Rect{{2, 3}, {4, 5}}; if TotalWidth(rs) != 6 { t.Fatalf("got %d", TotalWidth(rs)) }`,
							`if TotalWidth([]Rect{}) != 0 { t.Fatal("empty") }`,
						},
					},
				},
				{
					ID:    "go-a-json",
					Title: "JSON & struct tags",
					Lesson: `The encoding/json package converts Go values to and from JSON. json.Marshal
returns the JSON bytes:
    b, err := json.Marshal(value)

Only exported (capitalized) fields are encoded. Struct tags rename fields in
the output — written in back-quotes after the field type (shown here with
parentheses):
    type Point struct {
        X int  (json:"x")
        Y int  (json:"y")
    }
Marshalling Point{1, 2} then yields {"x":1,"y":2}. json.Unmarshal goes the
other way, filling a struct from JSON bytes.`,
					Challenge: Challenge{
						Prompt:      "Given Point (with json tags \"x\" and \"y\"), write ToJSON(p Point) (string, error) returning its JSON encoding as a string, e.g. {\"x\":1,\"y\":2}.",
						StarterCode: "type Point struct {\n\tX int `json:\"x\"`\n\tY int `json:\"y\"`\n}\n\nfunc ToJSON(p Point) (string, error) {\n\treturn \"\", nil\n}\n",
						Solution:    "import \"encoding/json\"\n\ntype Point struct {\n\tX int `json:\"x\"`\n\tY int `json:\"y\"`\n}\n\nfunc ToJSON(p Point) (string, error) {\n\tb, err := json.Marshal(p)\n\tif err != nil {\n\t\treturn \"\", err\n\t}\n\treturn string(b), nil\n}\n",
						Tests: []string{
							`got, err := ToJSON(Point{1, 2}); if err != nil || got != "{\"x\":1,\"y\":2}" { t.Fatalf("got %q err %v", got, err) }`,
							`got, _ := ToJSON(Point{0, 0}); if got != "{\"x\":0,\"y\":0}" { t.Fatalf("zero -> %q", got) }`,
						},
					},
				},
				{
					ID:    "go-a-constructors",
					Title: "Constructors & pointer receivers",
					Lesson: `Go has no classes, but structs + methods cover the same ground. A constructor
is just a function (by convention NewXxx) that builds and returns a value,
often a pointer:
    func NewAccount(initial int) *Account {
        return &Account{balance: initial}
    }

A method with a pointer receiver (*T) can modify the value; a value receiver
(T) gets a copy and can only read:
    func (a *Account) Deposit(n int) { a.balance += n }   // mutates
    func (a Account) Balance() int   { return a.balance } // reads

Call them the same way: a.Deposit(50); Go takes the address automatically.`,
					Challenge: Challenge{
						Prompt:      "Implement NewAccount(initial int) *Account, the pointer-receiver method Deposit(n int) (adds to the balance), and Balance() int.",
						StarterCode: "type Account struct {\n\tbalance int\n}\n\nfunc NewAccount(initial int) *Account {\n\treturn nil\n}\n\nfunc (a *Account) Deposit(n int) {\n}\n\nfunc (a Account) Balance() int {\n\treturn 0\n}\n",
						Solution:    "type Account struct {\n\tbalance int\n}\n\nfunc NewAccount(initial int) *Account {\n\treturn &Account{balance: initial}\n}\n\nfunc (a *Account) Deposit(n int) {\n\ta.balance += n\n}\n\nfunc (a Account) Balance() int {\n\treturn a.balance\n}\n",
						Tests: []string{
							`a := NewAccount(100); a.Deposit(50); if a.Balance() != 150 { t.Fatalf("got %d", a.Balance()) }`,
							`a := NewAccount(0); if a.Balance() != 0 { t.Fatal("fresh account") }`,
						},
					},
				},
				{
					ID:    "go-a-composition",
					Title: "Composition & embedding",
					Lesson: `Go favors composition over inheritance. Embedding a type in a struct (writing
it with no field name) promotes its fields and methods to the outer struct:
    type Animal struct{ Name string }
    func (a Animal) Describe() string { return a.Name }

    type Dog struct {
        Animal        // embedded
        Breed string
    }

    d := Dog{Animal{"Rex"}, "Lab"}
    d.Describe()      // "Rex" — promoted from Animal
    d.Name            // also promoted

If the outer type declares a method with the same name, it overrides (shadows)
the embedded one; you can still reach the inner via d.Animal.Describe().`,
					Challenge: Challenge{
						Prompt:      "Implement Animal's Describe() string to return its Name. Dog embeds Animal, so Dog values should answer Describe() via promotion.",
						StarterCode: "type Animal struct {\n\tName string\n}\n\nfunc (a Animal) Describe() string {\n\treturn \"\"\n}\n\ntype Dog struct {\n\tAnimal\n\tBreed string\n}\n",
						Solution:    "type Animal struct {\n\tName string\n}\n\nfunc (a Animal) Describe() string {\n\treturn a.Name\n}\n\ntype Dog struct {\n\tAnimal\n\tBreed string\n}\n",
						Tests: []string{
							`d := Dog{Animal{"Rex"}, "Lab"}; if d.Describe() != "Rex" { t.Fatalf("got %q", d.Describe()) }`,
							`a := Animal{"Cat"}; if a.Describe() != "Cat" { t.Fatal("animal") }`,
						},
					},
				},
				{
					ID:    "go-a-interfaces",
					Title: "Interfaces",
					Lesson: `An interface is a set of method signatures. A type satisfies an interface
simply by having those methods — there is no "implements" keyword, so types and
interfaces can be defined independently:
    type Shape interface {
        Area() int
    }

If Rect and Circle both have an Area() int method, both are Shapes, and one
function can handle either:
    func totalArea(shapes []Shape) int {
        sum := 0
        for _, s := range shapes {
            sum += s.Area()
        }
        return sum
    }

This implicit, structural satisfaction is the heart of Go's polymorphism.`,
					Challenge: Challenge{
						Prompt:      "Implement Rect.Area() (W*H) and Square.Area() (Side*Side) so both satisfy Shape, then write TotalArea(shapes []Shape) int summing every shape's area.",
						StarterCode: "type Shape interface {\n\tArea() int\n}\n\ntype Rect struct {\n\tW, H int\n}\n\nfunc (r Rect) Area() int {\n\treturn 0\n}\n\ntype Square struct {\n\tSide int\n}\n\nfunc (s Square) Area() int {\n\treturn 0\n}\n\nfunc TotalArea(shapes []Shape) int {\n\treturn 0\n}\n",
						Solution:    "type Shape interface {\n\tArea() int\n}\n\ntype Rect struct {\n\tW, H int\n}\n\nfunc (r Rect) Area() int {\n\treturn r.W * r.H\n}\n\ntype Square struct {\n\tSide int\n}\n\nfunc (s Square) Area() int {\n\treturn s.Side * s.Side\n}\n\nfunc TotalArea(shapes []Shape) int {\n\tsum := 0\n\tfor _, s := range shapes {\n\t\tsum += s.Area()\n\t}\n\treturn sum\n}\n",
						Tests: []string{
							`if (Rect{2, 3}).Area() != 6 { t.Fatal("Rect.Area") }`,
							`if (Square{4}).Area() != 16 { t.Fatal("Square.Area") }`,
							`shapes := []Shape{Rect{2, 3}, Square{4}}; if TotalArea(shapes) != 22 { t.Fatalf("got %d", TotalArea(shapes)) }`,
						},
					},
				},
			},
		},
		{
			Name: "Pointers",
			Topics: []Topic{
				{
					ID:    "go-a-pointers",
					Title: "Pointers",
					Lesson: `A pointer holds the address of a value. &x takes the address of x; *p
dereferences p to read or write the value it points at:
    x := 10
    p := &x          // p is a *int
    fmt.Println(*p)  // 10
    *p = 20          // writes through the pointer
    fmt.Println(x)   // 20

The zero value of a pointer is nil; dereferencing nil panics. Pointers let
functions refer to the same value instead of a copy — essential for sharing and
mutating state across calls.`,
					Challenge: Challenge{
						Prompt:      "Write Swap(a, b *int) that exchanges the two integers its pointers refer to.",
						StarterCode: "func Swap(a, b *int) {\n}\n",
						Solution:    "func Swap(a, b *int) {\n\t*a, *b = *b, *a\n}\n",
						Tests: []string{
							`x, y := 1, 2; Swap(&x, &y); if x != 2 || y != 1 { t.Fatalf("got %d %d", x, y) }`,
							`a, b := 5, 5; Swap(&a, &b); if a != 5 || b != 5 { t.Fatal("equal") }`,
						},
					},
				},
				{
					ID:    "go-a-mutation",
					Title: "Pointers for mutation",
					Lesson: `Because Go passes arguments by value, a function that takes a struct by value
gets a copy — changes to it don't affect the caller's value. To mutate the
original, pass a pointer:
    func grow(r *Rect, factor int) {
        r.W *= factor      // r.W is shorthand for (*r).W
        r.H *= factor
    }
    r := Rect{2, 3}
    grow(&r, 2)            // r is now {4, 6}

Go lets you write r.W on a pointer directly (auto-dereference). This is also why
mutating methods use pointer receivers. Slices and maps are already
reference-like, so functions can modify their contents without a pointer.`,
					Challenge: Challenge{
						Prompt:      "Given Rect, write Grow(r *Rect, factor int) that multiplies both W and H of the pointed-to rectangle by factor, in place.",
						StarterCode: "type Rect struct {\n\tW, H int\n}\n\nfunc Grow(r *Rect, factor int) {\n}\n",
						Solution:    "type Rect struct {\n\tW, H int\n}\n\nfunc Grow(r *Rect, factor int) {\n\tr.W *= factor\n\tr.H *= factor\n}\n",
						Tests: []string{
							`r := Rect{2, 3}; Grow(&r, 2); if r.W != 4 || r.H != 6 { t.Fatalf("got %+v", r) }`,
							`r := Rect{5, 5}; Grow(&r, 0); if r.W != 0 || r.H != 0 { t.Fatal("factor 0") }`,
						},
					},
				},
			},
		},
	}
}
