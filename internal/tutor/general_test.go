package tutor

import (
	"encoding/json"
	"testing"
)

func TestParseNoteContent(t *testing.T) {
	raw := "```json\n" +
		`{"title":"Derivatives","subject":"math","tags":["calculus"],"body":"A derivative measures change. See [[Limits]]."}` +
		"\n```"
	nc, err := parseNoteContent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if nc.Title != "Derivatives" || nc.Subject != "math" {
		t.Fatalf("metadata: %+v", nc)
	}
	if len(nc.Tags) != 1 || nc.Tags[0] != "calculus" {
		t.Fatalf("tags: %v", nc.Tags)
	}
	if nc.Body == "" {
		t.Fatalf("body lost")
	}
}

func TestParseEssayGrade(t *testing.T) {
	g, err := parseEssayGrade(`{"score":0.8,"feedback":"Good, but mention causes."}`)
	if err != nil {
		t.Fatal(err)
	}
	if g.Score != 0.8 || g.Feedback == "" {
		t.Fatalf("grade: %+v", g)
	}
}

func TestExtractJSONObject(t *testing.T) {
	if s, ok := extractJSONObject("noise {\"a\":1} trailing"); !ok || s != `{"a":1}` {
		t.Fatalf("plain: %q ok=%v", s, ok)
	}
	if _, ok := extractJSONObject("no object here"); ok {
		t.Fatalf("should not find an object")
	}

	// Quirks local models routinely emit.
	cases := []struct{ name, raw, want string }{
		{"fenced", "```json\n{\"a\":1}\n```", `{"a":1}`},
		{"chatty suffix", "{\"a\":1}\n\nHope that helps!", `{"a":1}`},
		{"brace in trailing prose", "{\"a\":1}\nuse {braces} freely", `{"a":1}`},
		{"trailing comma object", "{\"a\":1,}", `{"a":1}`},
		{"trailing comma array", "{\"a\":[1,2,],}", `{"a":[1,2]}`},
		{"trailing comma with newline", "{\n\"a\":1,\n}", "{\n\"a\":1\n}"},
		{"comma inside string kept", `{"a":"x,}"}`, `{"a":"x,}"}`},
		{"brace inside string", `{"a":"}{","b":2}`, `{"a":"}{","b":2}`},
		{"nested object", `{"a":{"b":2}} tail`, `{"a":{"b":2}}`},
	}
	for _, c := range cases {
		got, ok := extractJSONObject(c.raw)
		if !ok || got != c.want {
			t.Fatalf("%s: got %q ok=%v, want %q", c.name, got, ok, c.want)
		}
		var v any
		if err := json.Unmarshal([]byte(got), &v); err != nil {
			t.Fatalf("%s: extracted JSON does not parse: %v", c.name, err)
		}
	}

	// Truncated output (a cut-off local model) has no matching close brace.
	if _, ok := extractJSONObject(`{"a":1, "b":`); ok {
		t.Fatalf("truncated object should not be reported as found")
	}
}

func TestOfflineNoteAndGrade(t *testing.T) {
	nc := offlineNote("the cold war")
	if nc.Title != "the cold war" || nc.Body == "" {
		t.Fatalf("offline note: %+v", nc)
	}
	if g := offlineEssayGrade(""); g.Score != 0 {
		t.Fatalf("empty answer should score 0: %+v", g)
	}
	if g := offlineEssayGrade("some thoughtful answer"); g.Score != 1 {
		t.Fatalf("non-empty answer should pass offline: %+v", g)
	}
}
