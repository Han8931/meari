package executor

import (
	"os/exec"
	"testing"
)

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

// rustReady reports whether rustc can actually compile — a bare rustup shim
// with no default toolchain is on PATH but fails, so LookPath alone isn't enough.
func rustReady(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("rustc"); err != nil {
		t.Skip("rustc not installed")
	}
	if err := exec.Command("rustc", "--version").Run(); err != nil {
		t.Skip("rustc has no usable toolchain (run: rustup default stable)")
	}
}

func TestRunRustPass(t *testing.T) {
	rustReady(t)
	code := "fn add(a: i32, b: i32) -> i32 {\n    a + b\n}\n"
	tests := []string{
		`assert_eq!(add(2, 3), 5, "add(2,3)");`,
		`assert_eq!(add(-1, 1), 0, "add(-1,1)");`,
	}
	res, err := Run("rust", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected pass, got output:\n%s", res.Output)
	}
}

func TestRunRustFail(t *testing.T) {
	rustReady(t)
	code := "fn add(a: i32, b: i32) -> i32 {\n    a - b\n}\n" // wrong
	tests := []string{`assert_eq!(add(2, 3), 5, "add(2,3)");`}
	res, err := Run("rust", code, tests)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Fatal("expected failure for wrong Rust solution")
	}
}

func TestRunPhysicsReflection(t *testing.T) {
	res, err := Run("physics", "Energy changes form but is conserved.", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected submitted reflection to pass, got %q", res.Output)
	}

	res, err = Run("physics", "   ", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed {
		t.Fatal("empty physics response should not pass")
	}

	res, err = Run("quiz", "1A", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected submitted quiz answer to pass, got %q", res.Output)
	}
}
