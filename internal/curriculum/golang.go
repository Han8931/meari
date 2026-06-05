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
fmt.Sprintf returns the formatted string instead of printing it.

How challenges work: the editor gives you a function stub like
    func Describe(name string, age int) string {
        return ""
    }
Fill in the body and return the result — the tests call your function and
check its return value.`,
					Challenge: Challenge{
						Prompt:      "Write Describe(name string, age int) string returning e.g. \"Ada is 36 years old\". Use fmt.Sprintf.",
						StarterCode: "import \"fmt\"\n\nfunc Describe(name string, age int) string {\n\treturn \"\"\n}\n",
						Solution:    "import \"fmt\"\n\nfunc Describe(name string, age int) string {\n\treturn fmt.Sprintf(\"%s is %d years old\", name, age)\n}\n",
						Tests: []string{
							`if Describe("Ada", 36) != "Ada is 36 years old" { t.Fatalf("got %q", Describe("Ada", 36)) }`,
							`if Describe("Sam", 7) != "Sam is 7 years old" { t.Fatal("Sam") }`,
						},
					},
				},
				{
					ID:    "go-b-arith",
					Title: "Arithmetic",
					Lesson: `The numeric operators are + - * / and % (remainder). On integers, / TRUNCATES
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
    m := minutes % 60       // 10`,
					Challenge: Challenge{
						Prompt:      "Write Average3(a, b, c int) int returning the integer average of the three values (the / operator truncates, which is fine here).",
						StarterCode: "func Average3(a, b, c int) int {\n\treturn 0\n}\n",
						Solution:    "func Average3(a, b, c int) int {\n\treturn (a + b + c) / 3\n}\n",
						Tests: []string{
							`if Average3(1, 2, 3) != 2 { t.Fatalf("got %d", Average3(1, 2, 3)) }`,
							`if Average3(10, 10, 10) != 10 { t.Fatal("equal values") }`,
							`if Average3(1, 1, 2) != 1 { t.Fatal("truncation: 4/3 = 1") }`,
							`if Average3(-3, 0, 3) != 0 { t.Fatal("negatives") }`,
						},
					},
				},
				{
					ID:    "go-b-bools",
					Title: "Booleans & comparisons",
					Lesson: `A bool is true or false. Comparisons produce bools:
    ==  !=  <  <=  >  >=
    age >= 18          // true or false
    name == "Ada"      // strings compare with == too

Combine bools with && (and), || (or), and ! (not):
    age >= 13 && age <= 19      // a teenager
    day == "Sat" || day == "Sun"
    !done

&& and || short-circuit: the right side is only evaluated when it can still
change the answer. That makes guards safe and idiomatic:
    n != 0 && total/n > 10      // never divides by zero

Worked example — is a year a leap year?
    leap := year%4 == 0 && (year%100 != 0 || year%400 == 0)`,
					Challenge: Challenge{
						Prompt:      "Write InRange(n, lo, hi int) bool reporting whether n is between lo and hi inclusive (use && with two comparisons).",
						StarterCode: "func InRange(n, lo, hi int) bool {\n\treturn false\n}\n",
						Solution:    "func InRange(n, lo, hi int) bool {\n\treturn lo <= n && n <= hi\n}\n",
						Tests: []string{
							`if !InRange(5, 1, 10) { t.Fatal("5 is in 1..10") }`,
							`if !InRange(1, 1, 10) || !InRange(10, 1, 10) { t.Fatal("bounds are inclusive") }`,
							`if InRange(0, 1, 10) || InRange(11, 1, 10) { t.Fatal("outside") }`,
						},
					},
				},
				{
					ID:    "go-b-if",
					Title: "Branches: if & switch",
					Lesson: `Branching uses if/else — no parentheses around the condition, braces always
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

A switch can also match values directly: switch day { case "Sat", "Sun": ... }`,
					Challenge: Challenge{
						Prompt:      "Write Grade(score int) string returning \"A\" for 90+, \"B\" for 80+, \"C\" for 70+, and \"F\" below that (a condition switch works well).",
						StarterCode: "func Grade(score int) string {\n\treturn \"\"\n}\n",
						Solution:    "func Grade(score int) string {\n\tswitch {\n\tcase score >= 90:\n\t\treturn \"A\"\n\tcase score >= 80:\n\t\treturn \"B\"\n\tcase score >= 70:\n\t\treturn \"C\"\n\tdefault:\n\t\treturn \"F\"\n\t}\n}\n",
						Tests: []string{
							`if Grade(95) != "A" { t.Fatal("95") }`,
							`if Grade(90) != "A" { t.Fatal("90 is an A") }`,
							`if Grade(85) != "B" || Grade(71) != "C" { t.Fatal("middle grades") }`,
							`if Grade(69) != "F" { t.Fatal("69") }`,
						},
					},
				},
				{
					ID:    "go-b-loops",
					Title: "Loops",
					Lesson: `Go has one loop keyword — for — covering every case:
    for i := 0; i < n; i++ { ... }   // classic counter
    for cond { ... }                 // like while
    for { ... }                      // infinite (use break to leave)
(There is also "for ... range" for walking collections — you'll meet it when
collections are introduced.)

break leaves the loop; continue skips to the next iteration:
    for i := 1; i <= 100; i++ {
        if i%2 == 1 {
            continue        // skip odd numbers
        }
        if i > 40 {
            break           // stop entirely
        }
    }

Worked example — multiplying the numbers 1..n:
    fact := 1
    for i := 2; i <= n; i++ {
        fact *= i
    }`,
					Challenge: Challenge{
						Prompt:      "Write SumTo(n int) int returning 1 + 2 + ... + n computed with a for loop (return 0 when n < 1).",
						StarterCode: "func SumTo(n int) int {\n\treturn 0\n}\n",
						Solution:    "func SumTo(n int) int {\n\ttotal := 0\n\tfor i := 1; i <= n; i++ {\n\t\ttotal += i\n\t}\n\treturn total\n}\n",
						Tests: []string{
							`if SumTo(5) != 15 { t.Fatalf("SumTo(5) -> %d", SumTo(5)) }`,
							`if SumTo(1) != 1 { t.Fatal("SumTo(1)") }`,
							`if SumTo(0) != 0 { t.Fatal("n < 1 should give 0") }`,
							`if SumTo(100) != 5050 { t.Fatal("SumTo(100)") }`,
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
    if r := n % 10; r != 0 {
        fmt.Println(r)      // r lives only inside this if/else
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
				{
					ID:    "go-b-evens",
					Title: "Practice: loops + branches",
					Lesson: `Everything so far combines: a loop walks the numbers, a branch picks the ones
that matter, a variable accumulates the answer. This shape — loop, test,
accumulate — solves an enormous number of small problems.

Worked example — counting multiples of 3 from 1 to n:
    count := 0
    for i := 1; i <= n; i++ {
        if i%3 == 0 {
            count++
        }
    }

Variations on the same shape: sum instead of count, track the largest value
seen so far, or stop early with break when something is found. Try re-writing
the example as a sum in your head before doing the challenge.`,
					Challenge: Challenge{
						Prompt:      "Write CountEvens(from, to int) int counting the even numbers between from and to inclusive (return 0 when from > to).",
						StarterCode: "func CountEvens(from, to int) int {\n\treturn 0\n}\n",
						Solution:    "func CountEvens(from, to int) int {\n\tcount := 0\n\tfor i := from; i <= to; i++ {\n\t\tif i%2 == 0 {\n\t\t\tcount++\n\t\t}\n\t}\n\treturn count\n}\n",
						Tests: []string{
							`if CountEvens(1, 10) != 5 { t.Fatalf("1..10 -> %d", CountEvens(1, 10)) }`,
							`if CountEvens(2, 2) != 1 { t.Fatal("2..2") }`,
							`if CountEvens(3, 3) != 0 { t.Fatal("3..3") }`,
							`if CountEvens(5, 1) != 0 { t.Fatal("from > to") }`,
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
within a small tolerance, using math.Abs (add import "math" above your
function to use it):
    math.Abs(a-b) < 1e-9

Printf's %f verb controls formatting: %8.3f means width 8, 3 digits after the
point.`,
					Challenge: Challenge{
						Prompt:      "Write NearlyEqual(a, b, tol float64) bool that reports whether a and b differ by at most tol (use math.Abs).",
						StarterCode: "import \"math\"\n\nfunc NearlyEqual(a, b, tol float64) bool {\n\treturn false\n}\n",
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

Types don't mix implicitly — you must convert explicitly with T(value):
    var i int = 300
    var b uint8 = uint8(i)   // 300 wraps to 44
    var back int = int(b)    // 44

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
and big.Float (arbitrary precision).`,
					Challenge: Challenge{
						Prompt:      "Write Power(base, exp int64) string returning base raised to exp as a decimal string, using math/big so it works far beyond int64 (e.g. Power(2, 64)).",
						StarterCode: "import \"math/big\"\n\nfunc Power(base, exp int64) string {\n\treturn \"\"\n}\n",
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
					ID:    "go-b-build",
					Title: "Building strings",
					Lesson: `Join strings with +, and grow one in a loop with +=. Strings compare with ==
and order lexically with < and >.

    greeting := "Hello, " + name + "!"

    line := ""
    for i := 0; i < 3; i++ {
        line += "ab"
    }
    // line == "ababab"

An empty string "" is the zero value — the natural starting point for an
accumulator, just like 0 for a sum.

Worked example — a separated list without a trailing separator:
    out := ""
    for i := 1; i <= 3; i++ {
        if out != "" {
            out += "-"
        }
        out += "x"
    }
    // out == "x-x-x"`,
					Challenge: Challenge{
						Prompt:      "Write Repeat(word string, n int) string returning word repeated n times using a loop (\"\" when n < 1).",
						StarterCode: "func Repeat(word string, n int) string {\n\treturn \"\"\n}\n",
						Solution:    "func Repeat(word string, n int) string {\n\tout := \"\"\n\tfor i := 0; i < n; i++ {\n\t\tout += word\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`if Repeat("ha", 3) != "hahaha" { t.Fatalf("got %q", Repeat("ha", 3)) }`,
							`if Repeat("x", 0) != "" { t.Fatal("n < 1") }`,
							`if Repeat("ab", 1) != "ab" { t.Fatal("once") }`,
						},
					},
				},
				{
					ID:    "go-b-strings",
					Title: "Strings",
					Lesson: `A string is an immutable sequence of bytes, usually UTF-8 text. You can read a
byte by index but never assign to one:
    s := "shalom"
    fmt.Println(s[0])     // 115 (the byte 's')
    // s[0] = 'x'         // compile error: strings are immutable

Slicing takes a substring by byte positions — s[i:j] is bytes i up to (not
including) j; omit an end to mean "from the start" or "to the end":
    s[0:3]   // "sha"
    s[:1]    // "s"  (the first byte, as a string)
    s[3:]    // "lom"

len(s) returns the number of BYTES, not characters. Interpreted literals use
double quotes and honor escapes (\n, \t); raw literals use back-quotes and
take the bytes verbatim across multiple lines.

The strings package has the everyday helpers: Fields (split on whitespace),
Split, Join, ToUpper/ToLower, Contains, HasPrefix, ReplaceAll. Fields returns
a collection you can walk with range:
    for _, w := range strings.Fields("ada lovelace") {
        fmt.Println(w)        // "ada", then "lovelace"
    }`,
					Challenge: Challenge{
						Prompt:      "Write Initials(name string) string returning the uppercased first letter of each space-separated word (e.g. \"ada lovelace\" -> \"AL\"). Use strings.Fields, slicing, and strings.ToUpper.",
						StarterCode: "import \"strings\"\n\nfunc Initials(name string) string {\n\treturn \"\"\n}\n",
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
						StarterCode: "import \"unicode/utf8\"\n\nfunc RuneCount(s string) int {\n\treturn 0\n}\n",
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

Numbers and strings need the strconv package, not a plain conversion
(string(65) gives the character "A", not "65"):
    s := strconv.Itoa(10)       // "10"

To use a package, import it above your functions:
    import "strconv"

(strconv.Atoi parses the other way, string to int — it returns a value AND an
error, a two-result pattern you'll learn with functions later.)

Worked example — labeling a count:
    msg := "items: " + strconv.Itoa(count)`,
					Challenge: Challenge{
						Prompt:      "Write FizzBuzzAt(n int) string: \"Fizz\" if n is divisible by 3, \"Buzz\" by 5, \"FizzBuzz\" by both, otherwise n itself as a string (strconv.Itoa).",
						StarterCode: "import \"strconv\"\n\nfunc FizzBuzzAt(n int) string {\n\treturn \"\"\n}\n",
						Solution:    "import \"strconv\"\n\nfunc FizzBuzzAt(n int) string {\n\tswitch {\n\tcase n%15 == 0:\n\t\treturn \"FizzBuzz\"\n\tcase n%3 == 0:\n\t\treturn \"Fizz\"\n\tcase n%5 == 0:\n\t\treturn \"Buzz\"\n\tdefault:\n\t\treturn strconv.Itoa(n)\n\t}\n}\n",
						Tests: []string{
							`if FizzBuzzAt(3) != "Fizz" { t.Fatal("3") }`,
							`if FizzBuzzAt(5) != "Buzz" { t.Fatal("5") }`,
							`if FizzBuzzAt(15) != "FizzBuzz" { t.Fatal("15") }`,
							`if FizzBuzzAt(7) != "7" { t.Fatalf("7 -> %q", FizzBuzzAt(7)) }`,
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
		{
			Name: "App development",
			Topics: []Topic{
				{
					ID:    "go-a-errors",
					Title: "Error wrapping",
					Lesson: `Production Go code treats errors as values with context. Return early when an
operation fails, and wrap lower-level errors so callers can still inspect them:
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("parse port %q: %w", s, err)
    }

For domain errors, define a sentinel with errors.New and wrap it with %w:
    var ErrInvalidPort = errors.New("invalid port")
    return 0, fmt.Errorf("%w: %d", ErrInvalidPort, n)

Callers can then use errors.Is(err, ErrInvalidPort), even when the error has
extra detail. Good app errors say what failed and preserve the cause.`,
					Challenge: Challenge{
						Prompt:      "Write ParsePort(s string) (int, error). It should parse a TCP port, require 1..65535, wrap parse errors with context, and wrap ErrInvalidPort for out-of-range values.",
						StarterCode: "import \"errors\"\n\nvar ErrInvalidPort = errors.New(\"invalid port\")\n\nfunc ParsePort(s string) (int, error) {\n\treturn 0, nil\n}\n",
						Solution:    "import (\n\t\"errors\"\n\t\"fmt\"\n\t\"strconv\"\n)\n\nvar ErrInvalidPort = errors.New(\"invalid port\")\n\nfunc ParsePort(s string) (int, error) {\n\tn, err := strconv.Atoi(s)\n\tif err != nil {\n\t\treturn 0, fmt.Errorf(\"parse port %q: %w\", s, err)\n\t}\n\tif n < 1 || n > 65535 {\n\t\treturn 0, fmt.Errorf(\"%w: %d\", ErrInvalidPort, n)\n\t}\n\treturn n, nil\n}\n",
						Tests: []string{
							`got, err := ParsePort("8080"); if err != nil || got != 8080 { t.Fatalf("8080 -> %d %v", got, err) }`,
							`_, err := ParsePort("70000"); if !errors.Is(err, ErrInvalidPort) { t.Fatalf("want ErrInvalidPort, got %v", err) }`,
							`_, err := ParsePort("abc"); if err == nil || !strings.Contains(err.Error(), "parse port") { t.Fatalf("want contextual parse error, got %v", err) }`,
						},
					},
				},
				{
					ID:    "go-a-io-writer",
					Title: "io.Writer",
					Lesson: `The io package gives Go apps small interfaces that compose well. io.Writer is:
    type Writer interface {
        Write([]byte) (int, error)
    }

Files, network connections, buffers, HTTP responses, and compressors can all be
writers. Code that accepts io.Writer is easy to test because a bytes.Buffer can
stand in for a file or socket.

Always return write errors. App code often fails at the boundaries: disk full,
client disconnected, permission denied. Propagating those errors keeps failures
visible to the caller.`,
					Challenge: Challenge{
						Prompt:      "Write WriteLines(w io.Writer, lines []string) error. Write each line followed by \"\\n\" and stop immediately if a write fails.",
						StarterCode: "import \"io\"\n\nfunc WriteLines(w io.Writer, lines []string) error {\n\treturn nil\n}\n",
						Solution:    "import (\n\t\"fmt\"\n\t\"io\"\n)\n\nfunc WriteLines(w io.Writer, lines []string) error {\n\tfor _, line := range lines {\n\t\tif _, err := fmt.Fprintln(w, line); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\treturn nil\n}\n",
						Tests: []string{
							`var b bytes.Buffer; err := WriteLines(&b, []string{"alpha", "beta"}); if err != nil || b.String() != "alpha\nbeta\n" { t.Fatalf("%q %v", b.String(), err) }`,
							`var b bytes.Buffer; if err := WriteLines(&b, nil); err != nil || b.String() != "" { t.Fatalf("empty -> %q %v", b.String(), err) }`,
						},
					},
				},
				{
					ID:    "go-a-http-json",
					Title: "HTTP JSON handlers",
					Lesson: `A Go web API is often just functions with this shape:
    func handler(w http.ResponseWriter, r *http.Request)

Read from the request, set response headers, choose a status code, and write a
body. For JSON APIs, set Content-Type and use json.NewEncoder(w).Encode(value).

Handlers are easy to test with net/http/httptest:
    req := httptest.NewRequest(http.MethodGet, "/hello?name=Ada", nil)
    rec := httptest.NewRecorder()
    handler(rec, req)

That lets you verify status, headers, and body without opening a real port.`,
					Challenge: Challenge{
						Prompt:      "Write GreetHandler(w http.ResponseWriter, r *http.Request). It should respond with JSON {\"message\":\"Hello, <name>\"}, using query parameter name or \"friend\" by default.",
						StarterCode: "import \"net/http\"\n\nfunc GreetHandler(w http.ResponseWriter, r *http.Request) {\n}\n",
						Solution:    "import (\n\t\"encoding/json\"\n\t\"net/http\"\n)\n\nfunc GreetHandler(w http.ResponseWriter, r *http.Request) {\n\tname := r.URL.Query().Get(\"name\")\n\tif name == \"\" {\n\t\tname = \"friend\"\n\t}\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t_ = json.NewEncoder(w).Encode(map[string]string{\"message\": \"Hello, \" + name})\n}\n",
						Tests: []string{
							`req := httptest.NewRequest(http.MethodGet, "/hello?name=Ada", nil); rec := httptest.NewRecorder(); GreetHandler(rec, req); if rec.Code != http.StatusOK { t.Fatalf("status %d", rec.Code) }; if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") { t.Fatalf("content-type %q", rec.Header().Get("Content-Type")) }; if !strings.Contains(rec.Body.String(), "Hello, Ada") { t.Fatalf("body %q", rec.Body.String()) }`,
							`req := httptest.NewRequest(http.MethodGet, "/hello", nil); rec := httptest.NewRecorder(); GreetHandler(rec, req); if !strings.Contains(rec.Body.String(), "Hello, friend") { t.Fatalf("default body %q", rec.Body.String()) }`,
						},
					},
				},
				{
					ID:    "go-a-channels",
					Title: "Goroutines & channels",
					Lesson: `A goroutine runs a function concurrently:
    go doWork()

Channels let goroutines communicate. Send with ch <- v, receive with v := <-ch,
and close a channel when no more values will be sent. A range loop over a
channel receives until it is closed:
    for v := range ch {
        use(v)
    }

Prefer simple ownership: one goroutine sends and closes, another receives. This
keeps concurrent code readable and avoids races over shared memory.`,
					Challenge: Challenge{
						Prompt:      "Write Collect(ch <-chan int) []int that receives values until ch is closed and returns them in receive order.",
						StarterCode: "func Collect(ch <-chan int) []int {\n\treturn nil\n}\n",
						Solution:    "func Collect(ch <-chan int) []int {\n\tout := []int{}\n\tfor n := range ch {\n\t\tout = append(out, n)\n\t}\n\treturn out\n}\n",
						Tests: []string{
							`ch := make(chan int, 3); ch <- 1; ch <- 2; ch <- 3; close(ch); if !reflect.DeepEqual(Collect(ch), []int{1, 2, 3}) { t.Fatalf("got %v", Collect(ch)) }`,
							`ch := make(chan int); close(ch); if !reflect.DeepEqual(Collect(ch), []int{}) { t.Fatal("closed empty channel") }`,
						},
					},
				},
				{
					ID:    "go-a-context",
					Title: "Context cancellation",
					Lesson: `context.Context carries cancellation, deadlines, and request-scoped values
through app code. The most important rule: if your function may block, accept a
context and stop when ctx.Done() is closed.

The pattern is usually a select:
    select {
    case v := <-work:
        return v, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }

HTTP servers, databases, queues, and CLIs all use context so shutdowns and
timeouts can unwind work instead of leaking goroutines.`,
					Challenge: Challenge{
						Prompt:      "Write WaitForValue(ctx context.Context, ch <-chan string) (string, error). Return the value from ch, or return ctx.Err() if the context is canceled first.",
						StarterCode: "import \"context\"\n\nfunc WaitForValue(ctx context.Context, ch <-chan string) (string, error) {\n\treturn \"\", nil\n}\n",
						Solution:    "import \"context\"\n\nfunc WaitForValue(ctx context.Context, ch <-chan string) (string, error) {\n\tselect {\n\tcase v := <-ch:\n\t\treturn v, nil\n\tcase <-ctx.Done():\n\t\treturn \"\", ctx.Err()\n\t}\n}\n",
						Tests: []string{
							`ch := make(chan string, 1); ch <- "ready"; got, err := WaitForValue(context.Background(), ch); if err != nil || got != "ready" { t.Fatalf("%q %v", got, err) }`,
							`ctx, cancel := context.WithCancel(context.Background()); cancel(); got, err := WaitForValue(ctx, make(chan string)); if got != "" || !errors.Is(err, context.Canceled) { t.Fatalf("%q %v", got, err) }`,
						},
					},
				},
			},
		},
		{
			Name: "Bubble Tea TUI apps",
			Topics: []Topic{
				{
					ID:    "go-a-bubbletea-model",
					Title: "Model, Update, View",
					Lesson: `Bubble Tea uses The Elm Architecture. A TUI is a loop around three ideas:
    Model   // all state needed to draw the screen
    Update  // takes a message, returns the next model
    View    // renders the model as a string

Real Bubble Tea code uses tea.Model and tea.Msg, but the core idea is plain Go:
keep state in a struct, make state transitions explicit, and render from state
instead of printing from random places.

This makes TUIs testable. You can call Update with fake messages, then assert
on the new Model and View output.`,
					Challenge: Challenge{
						Prompt:      "Implement CounterModel.Update(msg Msg) CounterModel and CounterModel.View() string. \"up\" increments, \"down\" decrements, \"q\" sets Quit. View should render \"count: N\".",
						StarterCode: "type Msg string\n\ntype CounterModel struct {\n\tCount int\n\tQuit  bool\n}\n\nfunc (m CounterModel) Update(msg Msg) CounterModel {\n\treturn m\n}\n\nfunc (m CounterModel) View() string {\n\treturn \"\"\n}\n",
						Solution:    "import \"fmt\"\n\ntype Msg string\n\ntype CounterModel struct {\n\tCount int\n\tQuit  bool\n}\n\nfunc (m CounterModel) Update(msg Msg) CounterModel {\n\tswitch msg {\n\tcase \"up\":\n\t\tm.Count++\n\tcase \"down\":\n\t\tm.Count--\n\tcase \"q\":\n\t\tm.Quit = true\n\t}\n\treturn m\n}\n\nfunc (m CounterModel) View() string {\n\treturn fmt.Sprintf(\"count: %d\", m.Count)\n}\n",
						Tests: []string{
							`m := CounterModel{}.Update("up").Update("up").Update("down"); if m.Count != 1 || m.Quit { t.Fatalf("%+v", m) }`,
							`m := CounterModel{Count: 4}.Update("q"); if !m.Quit || m.View() != "count: 4" { t.Fatalf("%+v view=%q", m, m.View()) }`,
						},
					},
				},
				{
					ID:    "go-a-bubbletea-commands",
					Title: "Commands & messages",
					Lesson: `Bubble Tea keeps Update pure-ish by representing side effects as commands.
A command is a function that eventually returns a message:
    type Cmd func() Msg

Update can return a new model and a command. Bubble Tea runs the command and
feeds its message back into Update. That is how apps do timers, HTTP requests,
file reads, and subprocesses without freezing the UI.

The app stays understandable because state changes still happen in Update; the
command only reports what happened.`,
					Challenge: Challenge{
						Prompt:      "Implement Loader.Update(msg Msg) (Loader, Cmd). \"load\" should set Loading and return a command that emits \"loaded\". \"loaded\" should set Ready and clear Loading.",
						StarterCode: "type Msg string\ntype Cmd func() Msg\n\ntype Loader struct {\n\tLoading bool\n\tReady   bool\n}\n\nfunc (m Loader) Update(msg Msg) (Loader, Cmd) {\n\treturn m, nil\n}\n",
						Solution:    "type Msg string\ntype Cmd func() Msg\n\ntype Loader struct {\n\tLoading bool\n\tReady   bool\n}\n\nfunc (m Loader) Update(msg Msg) (Loader, Cmd) {\n\tswitch msg {\n\tcase \"load\":\n\t\tm.Loading = true\n\t\treturn m, func() Msg { return \"loaded\" }\n\tcase \"loaded\":\n\t\tm.Loading = false\n\t\tm.Ready = true\n\t}\n\treturn m, nil\n}\n",
						Tests: []string{
							`m, cmd := (Loader{}).Update("load"); if !m.Loading || m.Ready || cmd == nil { t.Fatalf("%+v cmd=%v", m, cmd) }`,
							`m, cmd := (Loader{}).Update("load"); m, cmd = m.Update(cmd()); if m.Loading || !m.Ready || cmd != nil { t.Fatalf("%+v cmd=%v", m, cmd) }`,
						},
					},
				},
				{
					ID:    "go-a-bubbletea-list",
					Title: "Lists & selection",
					Lesson: `Many Bubble Tea apps are lists: files, menu choices, search results, tasks.
The model usually stores the items and the selected index. Update handles keys
like up/down and clamps the cursor so it never points outside the slice.

View then renders every row, using a marker for the selected item:
    > current
      other

This pattern scales: later you can add scrolling, filtering, status badges, and
enter-to-open behavior without changing the basic state shape.`,
					Challenge: Challenge{
						Prompt:      "Implement Menu.Update(msg Msg) Menu and Menu.View() string. \"down\"/\"up\" move Selection within bounds. View should prefix the selected row with \"> \" and others with \"  \".",
						StarterCode: "type Msg string\n\ntype Menu struct {\n\tItems     []string\n\tSelection int\n}\n\nfunc (m Menu) Update(msg Msg) Menu {\n\treturn m\n}\n\nfunc (m Menu) View() string {\n\treturn \"\"\n}\n",
						Solution:    "import \"strings\"\n\ntype Msg string\n\ntype Menu struct {\n\tItems     []string\n\tSelection int\n}\n\nfunc (m Menu) Update(msg Msg) Menu {\n\tswitch msg {\n\tcase \"down\":\n\t\tif m.Selection < len(m.Items)-1 {\n\t\t\tm.Selection++\n\t\t}\n\tcase \"up\":\n\t\tif m.Selection > 0 {\n\t\t\tm.Selection--\n\t\t}\n\t}\n\treturn m\n}\n\nfunc (m Menu) View() string {\n\trows := make([]string, 0, len(m.Items))\n\tfor i, item := range m.Items {\n\t\tprefix := \"  \"\n\t\tif i == m.Selection {\n\t\t\tprefix = \"> \"\n\t\t}\n\t\trows = append(rows, prefix+item)\n\t}\n\treturn strings.Join(rows, \"\\n\")\n}\n",
						Tests: []string{
							`m := Menu{Items: []string{"one", "two", "three"}}.Update("down").Update("down").Update("down"); if m.Selection != 2 { t.Fatalf("selection %d", m.Selection) }`,
							`m := Menu{Items: []string{"one", "two"}, Selection: 1}.Update("up"); if m.Selection != 0 { t.Fatalf("selection %d", m.Selection) }`,
							`view := (Menu{Items: []string{"one", "two"}, Selection: 1}).View(); if view != "  one\n> two" { t.Fatalf("view %q", view) }`,
						},
					},
				},
			},
		},
		{
			Name: "Deferral and cleanup",
			Topics: []Topic{
				{
					ID:    "go-a-defer",
					Title: "Defer for cleanup",
					Lesson: `The defer statement schedules a function call to run just before the surrounding
function returns. This is the Go idiom for cleanup: closing files, unlocking
mutexes, releasing resources.

    file, err := os.Open("data.txt")
    if err != nil { return err }
    defer file.Close()   // Runs when enclosing function returns

Multiple defers run in LIFO order (last in, first out). You can also defer
recovery from panics:

    defer func() {
        if r := recover(); r != nil {
            log.Println("recovered:", r)
        }
    }()

Defer arguments are evaluated immediately (when the defer runs, not when called):
    defer fmt.Println("t", t)   // t captured at defer time

Use defer for cleanup at function entry - it's cleaner than multiple close()
calls at every return path.`,
					Challenge: Challenge{
						Prompt:      "Write ProcessFile(path string) error that opens a file, defers its Close(), reads all bytes, and returns an error if the file can't be opened or read. Use defer.",
						StarterCode: "import \"os\"\n\nfunc ProcessFile(path string) error {\n\treturn nil\n}\n",
						Solution:    "import (\n\t\"io\"\n\t\"os\"\n)\n\nfunc ProcessFile(path string) error {\n\tf, err := os.Open(path)\n\tif err != nil {\n\t\treturn err\n\t}\n\tdefer f.Close()\n\t_, err = io.ReadAll(f)\n\treturn err\n}\n",
						Tests: []string{
							`tmp, _ := os.CreateTemp("", "test-*"); tmp.WriteString("hello"); tmp.Close(); err := ProcessFile(tmp.Name()); os.Remove(tmp.Name()); if err != nil { t.Fatalf("%v", err) }`,
							`if err := ProcessFile("/nonexistent/path/file.txt"); err == nil { t.Fatal("expected error") }`,
						},
					},
				},
			},
		},
		{
			Name: "Synchronization primitives",
			Topics: []Topic{
				{
					ID:    "go-a-waitgroup",
					Title: "sync.WaitGroup",
					Lesson: `sync.WaitGroup waits for a collection of goroutines to finish. It has three
operations:

    var wg sync.WaitGroup
    wg.Add(1)           // Increment counter (before launching goroutine)
    go func() {
        defer wg.Done() // Decrement counter when done
        // work...
    }()
    wg.Wait()           // Block until counter is zero

Add() should be called before the goroutine starts; Done() is usually deferred
inside the goroutine. Wait() blocks until all expected goroutines have called Done().

This pattern lets you launch many concurrent tasks and wait for them all:

    for i := range jobs {
        wg.Add(1)
        go func(job int) {
            defer wg.Done()
            process(job)
        }(i)
    }
    wg.Wait()

When several goroutines touch the same variable, you also need a sync.Mutex so
only one enters the critical section at a time. Lock before the shared access
and Unlock after (defer is handy):
    var mu sync.Mutex
    mu.Lock()
    total += n       // safe: one goroutine here at a time
    mu.Unlock()

Without the lock, concurrent writes race and the total comes out wrong (run the
tests with -race to see it). WaitGroup answers "are they all done?"; Mutex
answers "who may touch the shared state right now?"`,
					Challenge: Challenge{
						Prompt:      "Write SumConcurrently(nums []int) int that spawns one goroutine per number, each adding its value to a shared total. Use sync.WaitGroup and sync.Mutex to coordinate.",
						StarterCode: "import \"sync\"\n\nfunc SumConcurrently(nums []int) int {\n\treturn 0\n}\n",
						Solution:    "import \"sync\"\n\nfunc SumConcurrently(nums []int) int {\n\tvar mu sync.Mutex\n\ttotal := 0\n\tvar wg sync.WaitGroup\n\tfor _, n := range nums {\n\t\twg.Add(1)\n\t\tgo func(val int) {\n\t\t\tdefer wg.Done()\n\t\t\tmu.Lock()\n\t\t\ttotal += val\n\t\t\tmu.Unlock()\n\t\t}(n)\n\t}\n\twg.Wait()\n\treturn total\n}\n",
						Tests: []string{
							`if got := SumConcurrently([]int{1, 2, 3, 4, 5}); got != 15 { t.Fatalf("got %d", got) }`,
							`if got := SumConcurrently([]int{}); got != 0 { t.Fatal("empty") }`,
							`if got := SumConcurrently([]int{100}); got != 100 { t.Fatal("single") }`,
						},
					},
				},
			},
		},
	}
}
