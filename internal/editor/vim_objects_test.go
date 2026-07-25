package editor

import "testing"

// --- text objects ---

func TestDeleteInnerWord(t *testing.T) {
	m := New("alpha beta gamma", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("w")) // cursor on "beta"
	m = apply(m, key("d"), key("i"), key("w"))
	if got := m.Value(); got != "alpha  gamma" {
		t.Fatalf("diw: %q", got)
	}
	if m.register != "beta" {
		t.Fatalf("diw register = %q", m.register)
	}
}

func TestChangeInnerWord(t *testing.T) {
	m := New("alpha beta", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("c"), key("i"), key("w"), key("h"), key("i"))
	if got := m.Value(); got != "hi beta" {
		t.Fatalf("ciw then type: %q", got)
	}
}

func TestDeleteAroundWord(t *testing.T) {
	m := New("alpha beta gamma", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("w")) // "beta"
	m = apply(m, key("d"), key("a"), key("w"))
	// aw takes the word and its trailing whitespace.
	if got := m.Value(); got != "alpha gamma" {
		t.Fatalf("daw: %q", got)
	}
}

func TestDeleteInnerQuote(t *testing.T) {
	m := New(`say "hello" now`, true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("d"), key("i"), key("\""))
	if got := m.Value(); got != `say "" now` {
		t.Fatalf(`di": %q`, got)
	}
	if m.register != "hello" {
		t.Fatalf(`di" register = %q`, m.register)
	}
}

func TestChangeInnerParen(t *testing.T) {
	m := New("foo(bar, baz)", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("f"), key("b")) // cursor inside the parens
	m = apply(m, key("c"), key("i"), key("("), key("x"))
	if got := m.Value(); got != "foo(x)" {
		t.Fatalf("ci( then type: %q", got)
	}
}

func TestDeleteAroundBrace(t *testing.T) {
	m := New("x{a}y", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("f"), key("a")) // inside the braces
	m = apply(m, key("d"), key("a"), key("{"))
	if got := m.Value(); got != "xy" {
		t.Fatalf("da{: %q", got)
	}
}

func TestYankInnerWord(t *testing.T) {
	m := New("alpha beta", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("y"), key("i"), key("w"))
	if m.register != "alpha" || m.regLinewise {
		t.Fatalf("yiw register = %q linewise=%v", m.register, m.regLinewise)
	}
}

// --- y with motions ---

func TestYankToWord(t *testing.T) {
	m := New("alpha beta", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("y"), key("w"))
	if m.register != "alpha " {
		t.Fatalf("yw register = %q", m.register)
	}
}

func TestYankToEnd(t *testing.T) {
	m := New("alpha beta", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("w"), key("y"), key("$"))
	if m.register != "beta" {
		t.Fatalf("y$ register = %q", m.register)
	}
}

func TestYankLineDown(t *testing.T) {
	m := New("one\ntwo\nthree", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("y"), key("j"))
	if m.register != "one\ntwo" || !m.regLinewise {
		t.Fatalf("yj register = %q linewise=%v", m.register, m.regLinewise)
	}
}

// --- d with the new motions ---

func TestDeleteBackWord(t *testing.T) {
	m := New("alpha beta", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("$"), key("d"), key("b"))
	if got := m.Value(); got != "alpha a" {
		t.Fatalf("db: %q", got)
	}
}

func TestDeleteLineDown(t *testing.T) {
	m := New("one\ntwo\nthree", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("d"), key("j"))
	if got := m.Value(); got != "three" {
		t.Fatalf("dj: %q", got)
	}
}

// --- X / s / S ---

func TestDeleteCharBefore(t *testing.T) {
	m := New("abc", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("l"), key("l"), key("X"))
	if got := m.Value(); got != "ac" {
		t.Fatalf("X: %q", got)
	}
}

func TestSubstituteChar(t *testing.T) {
	m := New("abc", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("s"), key("X"))
	if got := m.Value(); got != "Xbc" {
		t.Fatalf("s then type: %q", got)
	}
}

func TestSubstituteLine(t *testing.T) {
	m := New("one\ntwo", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("S"), key("X"))
	if got := m.Value(); got != "X\ntwo" {
		t.Fatalf("S then type: %q", got)
	}
}

// --- search: * and ? ---

func TestStarSearchesWord(t *testing.T) {
	m := New("foo bar\nbaz foo end", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("*"))
	if row, _ := m.cursorPos(); row != 1 {
		t.Fatalf("* should jump to the next 'foo' on line 1, row=%d", row)
	}
}

func TestBackwardSearch(t *testing.T) {
	m := New("foo\nbar\nfoo\nbar", true, nil)
	m.SetSize(60, 10)
	m = apply(m, key("G"))                                        // last line
	m = apply(m, key("?"), key("f"), key("o"), key("o"), enter()) // ?foo
	if row, _ := m.cursorPos(); row != 2 {
		t.Fatalf("?foo from the bottom should land on line 2, row=%d", row)
	}
}
