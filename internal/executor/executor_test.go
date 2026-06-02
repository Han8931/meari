package executor

import "testing"

func TestRunPythonPass(t *testing.T) {
	code := "def is_even(n):\n    return n % 2 == 0\n"
	tests := []string{"assert is_even(4) == True", "assert is_even(7) == False"}
	res, err := Run("python", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected pass, got output:\n%s", res.Output)
	}
}

func TestRunPythonFail(t *testing.T) {
	code := "def is_even(n):\n    return False\n" // wrong
	tests := []string{"assert is_even(4) == True"}
	res, err := Run("python", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Fatal("expected failure for wrong solution")
	}
}

func TestRunPythonTimeout(t *testing.T) {
	code := "def is_even(n):\n    while True:\n        pass\n"
	tests := []string{"assert is_even(4) == True"}
	res, err := Run("python", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if !res.TimedOut {
		t.Fatal("expected timeout for infinite loop")
	}
}

func TestRunGoPass(t *testing.T) {
	code := "func Add(a, b int) int {\n\treturn a + b\n}\n"
	tests := []string{
		`if Add(2, 3) != 5 { t.Fatalf("Add(2,3)=%d", Add(2,3)) }`,
		`if Add(-1, 1) != 0 { t.Fatal("Add(-1,1)") }`,
	}
	res, err := Run("go", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected pass, got output:\n%s", res.Output)
	}
}

func TestRunGoFail(t *testing.T) {
	code := "func Add(a, b int) int {\n\treturn a - b\n}\n" // wrong
	tests := []string{`if Add(2, 3) != 5 { t.Fatal("Add(2,3)") }`}
	res, err := Run("go", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Fatal("expected failure for wrong Go solution")
	}
}

func TestRunGoUsesReflect(t *testing.T) {
	code := "func Evens(n int) []int {\n\tout := []int{}\n\tfor i := 0; i < n; i++ {\n\t\tif i%2 == 0 {\n\t\t\tout = append(out, i)\n\t\t}\n\t}\n\treturn out\n}\n"
	tests := []string{`if !reflect.DeepEqual(Evens(6), []int{0, 2, 4}) { t.Fatalf("got %v", Evens(6)) }`}
	res, err := Run("go", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected pass, got output:\n%s", res.Output)
	}
}
