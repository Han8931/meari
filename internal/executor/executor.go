// Package executor runs a learner's code against a challenge's test cases.
//
// The tests are the ground truth for correctness. For Python, the solution and
// the test assertions are concatenated into one script and run as a subprocess.
// For Go, the solution becomes a package and the tests become a generated
// _test.go file run with "go test". Both are guarded by a timeout so infinite
// loops can't hang the program.
//
// Safety note: this executes code on the host with only a timeout guard. That is
// acceptable for a single trusted local learner. Do NOT run untrusted code this
// way; a real sandbox (container / seccomp) would be needed first.
package executor

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// builder constructs the command to run, given the timeout context and temp dir.
type builder func(ctx context.Context, dir string) *exec.Cmd

// Result captures the outcome of running a solution against its tests.
type Result struct {
	Passed   bool
	TimedOut bool
	Output   string // combined stdout+stderr, useful for hints on failure
}

// Run executes code+tests for the given language ("python" or "go"). Unknown or
// empty languages default to Python (the historical behavior). Non-programming
// curricula use "physics" or "essay" for prose exercises: a non-empty response
// counts as submitted so the tutor can discuss/grade it.
func Run(lang, code string, tests []string) (Result, error) {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return runGo(code, tests)
	case "physics", "essay":
		return runReflection(code)
	default:
		return runPython(code, tests)
	}
}

func runReflection(code string) (Result, error) {
	if strings.TrimSpace(code) == "" {
		return Result{Output: "write a short response before submitting"}, nil
	}
	return Result{Passed: true, Output: "response submitted"}, nil
}

// runPython writes the solution + assert tests to a temp .py file and runs it.
// Exit code 0 means every assertion passed.
func runPython(code string, tests []string) (Result, error) {
	var b strings.Builder
	b.WriteString(code)
	b.WriteString("\n\n# --- tutor test harness ---\n")
	for _, t := range tests {
		b.WriteString(t)
		b.WriteString("\n")
	}
	return runFile("python", map[string]string{"solution.py": b.String()},
		5*time.Second, func(ctx context.Context, dir string) *exec.Cmd {
			return exec.CommandContext(ctx, "python3", filepath.Join(dir, "solution.py"))
		})
}

// runGo writes the solution as a package and the tests as a generated test file,
// then runs "go test". reflect and fmt are pre-imported (and blank-referenced so
// Go doesn't complain) so tests can compare slices/maps and format values.
func runGo(code string, tests []string) (Result, error) {
	files := map[string]string{
		"go.mod":      "module sol\n\ngo 1.21\n",
		"sol.go":      "package sol\n\n" + code + "\n",
		"sol_test.go": goTestHarness(tests),
	}
	// Go's first compile can be slow; give it more headroom than Python.
	return runFile("go", files, 30*time.Second, func(ctx context.Context, dir string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, "go", "test", "./...")
		cmd.Dir = dir
		// Keep the build self-contained and offline.
		cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GO111MODULE=on")
		return cmd
	})
}

func goTestHarness(tests []string) string {
	var b strings.Builder
	b.WriteString("package sol\n\n")
	b.WriteString("import (\n\t\"testing\"\n\t\"reflect\"\n\t\"fmt\"\n\t\"bytes\"\n\t\"context\"\n\t\"encoding/json\"\n\t\"errors\"\n\t\"io\"\n\t\"net/http\"\n\t\"net/http/httptest\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\n")
	// Blank references so the imports are always \"used\" even if a given test set
	// doesn't touch them.
	b.WriteString("var (\n\t_ = reflect.DeepEqual\n\t_ = fmt.Sprint\n\t_ = bytes.Buffer{}\n\t_ = context.Background\n\t_ = json.NewDecoder\n\t_ = errors.Is\n\t_ io.Writer\n\t_ = http.MethodGet\n\t_ = httptest.NewRecorder\n\t_ = os.Open\n\t_ = strings.Contains\n\t_ = time.Millisecond\n)\n\n")
	b.WriteString("func TestAll(t *testing.T) {\n")
	// Each test gets its own block scope so local variable names (e.g. a ":=")
	// declared in different tests don't collide.
	for _, t := range tests {
		b.WriteString("\t{\n\t\t")
		b.WriteString(t)
		b.WriteString("\n\t}\n")
	}
	b.WriteString("}\n")
	return b.String()
}

// runFile writes files into a temp dir and runs the built command with a
// timeout. Exit code 0 means the tests passed.
func runFile(prefix string, files map[string]string, timeout time.Duration, build builder) (Result, error) {
	tmp, err := os.MkdirTemp("", "meari-"+prefix+"-")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(tmp)

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(content), 0o644); err != nil {
			return Result{}, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := build(ctx, tmp)
	out, runErr := cmd.CombinedOutput()

	res := Result{Output: strings.TrimSpace(string(out))}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		res.TimedOut = true
		if res.Output == "" {
			res.Output = "execution timed out (possible infinite loop)"
		}
		return res, nil
	}
	res.Passed = runErr == nil
	return res, nil
}
