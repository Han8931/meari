package curriculum

// Pre-authored Python curriculum. Every challenge's Solution is verified against
// its Tests in curriculum_test.go.

func pythonBeginner() []Module {
	return []Module{
		{
			Name: "Basics",
			Topics: []Topic{
				{
					ID:    "py-b-variables",
					Title: "Variables & types",
					Lesson: `Python stores values in variables with =. You don't declare a type; the
value's type is inferred. The core types you'll meet first are int (whole
numbers), float (decimals), str (text in quotes), and bool (True/False).

Example:
    name = "Ada"        # str
    age = 36            # int
    height = 1.7        # float
    is_admin = True     # bool
    print(name, age)    # Ada 36

You can combine values with operators: + adds numbers (or joins strings),
* multiplies, and so on.`,
					Challenge: Challenge{
						Prompt:      "Write rectangle_area(width, height) that returns the area of a rectangle.",
						StarterCode: "def rectangle_area(width, height):\n    pass\n",
						Solution:    "def rectangle_area(width, height):\n    return width * height\n",
						Tests: []string{
							"assert rectangle_area(3, 4) == 12",
							"assert rectangle_area(10, 2) == 20",
							"assert rectangle_area(0, 5) == 0",
						},
					},
				},
				{
					ID:    "py-b-conditionals",
					Title: "if / elif / else",
					Lesson: `Conditionals run different code depending on whether a condition is True.
Use if, optional elif (else-if) branches, and an optional else. The body of
each branch is indented (4 spaces by convention).

Example:
    def grade(score):
        if score >= 90:
            return "A"
        elif score >= 80:
            return "B"
        else:
            return "C"

Comparisons (==, !=, <, >, <=, >=) produce booleans you can branch on.`,
					Challenge: Challenge{
						Prompt:      `Write sign(n) that returns "positive", "negative", or "zero".`,
						StarterCode: "def sign(n):\n    pass\n",
						Solution:    "def sign(n):\n    if n > 0:\n        return \"positive\"\n    elif n < 0:\n        return \"negative\"\n    return \"zero\"\n",
						Tests: []string{
							`assert sign(5) == "positive"`,
							`assert sign(-3) == "negative"`,
							`assert sign(0) == "zero"`,
						},
					},
				},
				{
					ID:    "py-b-loops",
					Title: "Loops",
					Lesson: `A for loop repeats once per item in a sequence. range(n) gives 0..n-1, and
range(a, b) gives a..b-1. A while loop repeats as long as a condition holds.

Example:
    total = 0
    for i in range(1, 4):   # 1, 2, 3
        total += i          # total becomes 6
    print(total)            # 6

Accumulating into a variable (like total) inside a loop is a very common
pattern.`,
					Challenge: Challenge{
						Prompt:      "Write sum_to(n) that returns 1 + 2 + ... + n (0 if n < 1).",
						StarterCode: "def sum_to(n):\n    pass\n",
						Solution:    "def sum_to(n):\n    total = 0\n    for i in range(1, n + 1):\n        total += i\n    return total\n",
						Tests: []string{
							"assert sum_to(5) == 15",
							"assert sum_to(1) == 1",
							"assert sum_to(0) == 0",
						},
					},
				},
				{
					ID:    "py-b-functions",
					Title: "Functions",
					Lesson: `A function packages reusable logic. Define it with def, list its parameters
in parentheses, and hand back a result with return. A function without a
return statement returns None.

Example:
    def greet(name):
        return "Hello, " + name

    print(greet("Sam"))   # Hello, Sam

Functions can take several parameters and call other functions.`,
					Challenge: Challenge{
						Prompt:      "Write max_of_three(a, b, c) that returns the largest of three numbers (without using max()).",
						StarterCode: "def max_of_three(a, b, c):\n    pass\n",
						Solution:    "def max_of_three(a, b, c):\n    biggest = a\n    if b > biggest:\n        biggest = b\n    if c > biggest:\n        biggest = c\n    return biggest\n",
						Tests: []string{
							"assert max_of_three(1, 2, 3) == 3",
							"assert max_of_three(9, 4, 1) == 9",
							"assert max_of_three(2, 8, 5) == 8",
						},
					},
				},
				{
					ID:    "py-b-lists",
					Title: "Lists",
					Lesson: `A list is an ordered, changeable collection written with square brackets.
Index from 0 with lst[0]; append with lst.append(x); get the length with
len(lst); and loop over items with for x in lst.

Example:
    nums = [3, 1, 2]
    nums.append(4)        # [3, 1, 2, 4]
    print(nums[0])        # 3
    print(len(nums))      # 4

You often build a new list by appending inside a loop.`,
					Challenge: Challenge{
						Prompt:      "Write evens(nums) that returns a new list containing only the even numbers from nums, in order.",
						StarterCode: "def evens(nums):\n    pass\n",
						Solution:    "def evens(nums):\n    out = []\n    for n in nums:\n        if n % 2 == 0:\n            out.append(n)\n    return out\n",
						Tests: []string{
							"assert evens([1, 2, 3, 4]) == [2, 4]",
							"assert evens([1, 3, 5]) == []",
							"assert evens([2, 4, 6]) == [2, 4, 6]",
						},
					},
				},
			},
		},
	}
}

func pythonIntermediate() []Module {
	return []Module{
		{
			Name: "Working with data",
			Topics: []Topic{
				{
					ID:    "py-i-comprehensions",
					Title: "Comprehensions",
					Lesson: `A list comprehension builds a list in one expression:
[expr for item in iterable if condition]. It replaces the build-with-a-loop
pattern and reads top to bottom.

Example:
    squares = [n * n for n in range(5)]        # [0, 1, 4, 9, 16]
    odds    = [n for n in range(10) if n % 2]  # [1, 3, 5, 7, 9]

Dict comprehensions work similarly: {k: v for ...}.`,
					Challenge: Challenge{
						Prompt:      "Write square_odds(nums) that returns a list of the squares of only the odd numbers in nums, using a list comprehension.",
						StarterCode: "def square_odds(nums):\n    pass\n",
						Solution:    "def square_odds(nums):\n    return [n * n for n in nums if n % 2 == 1]\n",
						Tests: []string{
							"assert square_odds([1, 2, 3, 4]) == [1, 9]",
							"assert square_odds([2, 4]) == []",
							"assert square_odds([5]) == [25]",
						},
					},
				},
				{
					ID:    "py-i-dicts",
					Title: "Dictionaries",
					Lesson: `A dict maps keys to values: {"a": 1}. Look up with d[key], set with
d[key] = value, and check membership with key in d. d.get(key, default)
avoids a KeyError when a key may be missing.

Example:
    counts = {}
    for ch in "banana":
        counts[ch] = counts.get(ch, 0) + 1
    print(counts)   # {'b': 1, 'a': 3, 'n': 2}

This "count occurrences" pattern is everywhere.`,
					Challenge: Challenge{
						Prompt:      "Write word_count(words) that returns a dict mapping each word to how many times it appears in the list words.",
						StarterCode: "def word_count(words):\n    pass\n",
						Solution:    "def word_count(words):\n    counts = {}\n    for w in words:\n        counts[w] = counts.get(w, 0) + 1\n    return counts\n",
						Tests: []string{
							`assert word_count(["a", "b", "a"]) == {"a": 2, "b": 1}`,
							`assert word_count([]) == {}`,
							`assert word_count(["x"]) == {"x": 1}`,
						},
					},
				},
				{
					ID:    "py-i-errors",
					Title: "Error handling",
					Lesson: `Operations can raise exceptions. Wrap risky code in try/except to handle
them instead of crashing. You can catch a specific type like
ZeroDivisionError, and an optional else/finally block can follow.

Example:
    def parse(text):
        try:
            return int(text)
        except ValueError:
            return None

    print(parse("42"))   # 42
    print(parse("x"))    # None`,
					Challenge: Challenge{
						Prompt:      "Write safe_div(a, b) that returns a / b, or None if b is zero (catch the exception, don't check b directly).",
						StarterCode: "def safe_div(a, b):\n    pass\n",
						Solution:    "def safe_div(a, b):\n    try:\n        return a / b\n    except ZeroDivisionError:\n        return None\n",
						Tests: []string{
							"assert safe_div(10, 2) == 5",
							"assert safe_div(7, 0) is None",
							"assert safe_div(9, 3) == 3",
						},
					},
				},
				{
					ID:    "py-i-classes",
					Title: "Classes",
					Lesson: `A class bundles data and behavior. __init__ runs when you create an instance
and sets up attributes on self. Methods are functions defined inside the
class whose first parameter is self.

Example:
    class Dog:
        def __init__(self, name):
            self.name = name
        def speak(self):
            return self.name + " says woof"

    d = Dog("Rex")
    print(d.speak())   # Rex says woof`,
					Challenge: Challenge{
						Prompt:      "Implement a class Stack with methods push(x), pop() (return and remove the top), and is_empty(). Use a list internally.",
						StarterCode: "class Stack:\n    def __init__(self):\n        self.items = []\n\n    def push(self, x):\n        pass\n\n    def pop(self):\n        pass\n\n    def is_empty(self):\n        pass\n",
						Solution:    "class Stack:\n    def __init__(self):\n        self.items = []\n\n    def push(self, x):\n        self.items.append(x)\n\n    def pop(self):\n        return self.items.pop()\n\n    def is_empty(self):\n        return len(self.items) == 0\n",
						Tests: []string{
							"s = Stack()\nassert s.is_empty() == True",
							"s = Stack()\ns.push(1)\ns.push(2)\nassert s.pop() == 2\nassert s.pop() == 1\nassert s.is_empty() == True",
						},
					},
				},
			},
		},
	}
}

func pythonAdvanced() []Module {
	return []Module{
		{
			Name: "Powerful Python",
			Topics: []Topic{
				{
					ID:    "py-a-recursion",
					Title: "Recursion",
					Lesson: `A recursive function calls itself on a smaller input until it reaches a base
case that returns directly. Every recursive function needs a base case, or it
recurses forever.

Example:
    def countdown(n):
        if n <= 0:           # base case
            return "done"
        return countdown(n - 1)

Think: what is the simplest input I can answer immediately, and how do I
shrink the problem toward it?`,
					Challenge: Challenge{
						Prompt:      "Write factorial(n) recursively: factorial(0) is 1 and factorial(n) is n * factorial(n-1).",
						StarterCode: "def factorial(n):\n    pass\n",
						Solution:    "def factorial(n):\n    if n <= 1:\n        return 1\n    return n * factorial(n - 1)\n",
						Tests: []string{
							"assert factorial(0) == 1",
							"assert factorial(5) == 120",
							"assert factorial(3) == 6",
						},
					},
				},
				{
					ID:    "py-a-generators",
					Title: "Generators",
					Lesson: `A generator function uses yield to produce a sequence lazily, one value at a
time, without building the whole list in memory. Calling it returns a
generator you can iterate or pass to list().

Example:
    def first_n(n):
        i = 0
        while i < n:
            yield i
            i += 1

    print(list(first_n(3)))   # [0, 1, 2]

Each yield pauses the function until the next value is requested.`,
					Challenge: Challenge{
						Prompt:      "Write fib(n) as a generator that yields the first n Fibonacci numbers starting 0, 1, 1, 2, 3, ...",
						StarterCode: "def fib(n):\n    pass\n",
						Solution:    "def fib(n):\n    a, b = 0, 1\n    for _ in range(n):\n        yield a\n        a, b = b, a + b\n",
						Tests: []string{
							"assert list(fib(5)) == [0, 1, 1, 2, 3]",
							"assert list(fib(0)) == []",
							"assert list(fib(1)) == [0]",
						},
					},
				},
				{
					ID:    "py-a-decorators",
					Title: "Decorators",
					Lesson: `A decorator is a function that takes a function and returns a new function
wrapping it, adding behavior. The @name syntax applies one.

Example:
    def shout(func):
        def wrapper(*args, **kwargs):
            return func(*args, **kwargs).upper()
        return wrapper

    @shout
    def greet(name):
        return "hi " + name

    print(greet("sam"))   # HI SAM

The wrapper forwards *args/**kwargs so it works for any signature.`,
					Challenge: Challenge{
						Prompt:      "Write a decorator double(func) that returns a wrapper calling func and returning twice its result. e.g. double(lambda x: x + 1)(3) == 8.",
						StarterCode: "def double(func):\n    pass\n",
						Solution:    "def double(func):\n    def wrapper(*args, **kwargs):\n        return func(*args, **kwargs) * 2\n    return wrapper\n",
						Tests: []string{
							"assert double(lambda x: x + 1)(3) == 8",
							"assert double(lambda a, b: a + b)(2, 3) == 10",
							"assert double(lambda: 7)() == 14",
						},
					},
				},
			},
		},
	}
}
