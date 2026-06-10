package tutor

// General (subject-agnostic) tutor capabilities for the learning-vault model:
// generating a lesson as a markdown NOTE, and grading a free-text essay answer
// for any subject. These sit alongside the original coding-specific helpers in
// tutor.go; the existing methods stay until the core engine can update every
// call site at once.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// NoteContent is an AI-generated lesson shaped for the vault: structured
// metadata plus a markdown body the learner can keep, edit, and link from.
type NoteContent struct {
	Title   string   `json:"title"`
	Subject string   `json:"subject"`
	Tags    []string `json:"tags"`
	Body    string   `json:"body"` // markdown
}

// GenerateNote turns "I want to learn X" into a self-contained lesson note on
// any subject (languages, math, history, …) — not just programming. The body is
// markdown and may wrap key prerequisite concepts in [[wikilinks]] so the
// learner can branch out into linked notes.
func (t *Tutor) GenerateNote(ctx context.Context, request string) (NoteContent, error) {
	if t.offline {
		return offlineNote(request), nil
	}

	system := `You are a knowledgeable tutor creating a focused, self-contained lesson NOTE
for a self-directed learner. The subject can be anything (a language, math, history,
science, music, …), not only programming.
Respond with ONLY a JSON object, no prose, no markdown fences, matching:
{
  "title": "a short, specific note title",
  "subject": "the broad subject area, lowercase (e.g. \"math\", \"spanish\", \"history\")",
  "tags": ["2-4", "short", "kebab-case", "tags"],
  "body": "the lesson as MARKDOWN: a few short sections of clear prose, optionally one worked example. Wrap key prerequisite or related concepts in [[wikilinks]] so the learner can branch into linked notes. Do not restate the title as an H1 — start with the explanation."
}` + t.levelClause()

	raw, err := t.chat(ctx, system, "Teach me about: "+request)
	if err != nil {
		return NoteContent{}, err
	}
	nc, err := parseNoteContent(raw)
	if err != nil {
		return NoteContent{}, fmt.Errorf("could not parse generated note: %w", err)
	}
	if strings.TrimSpace(nc.Title) == "" {
		nc.Title = strings.TrimSpace(request)
	}
	if strings.TrimSpace(nc.Body) == "" {
		nc.Body = strings.TrimSpace(raw)
	}
	return nc, nil
}

// EssayGrade is the result of grading a free-text answer: a coarse 0..1 score
// (essays are not pass/fail like code tests) plus tutor feedback.
type EssayGrade struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// GradeEssay evaluates a learner's free-text response to a study prompt on any
// subject. It generalizes the old "reflection" path (any non-empty answer
// counts as submitted) by adding a rubric score and specific feedback.
func (t *Tutor) GradeEssay(ctx context.Context, prompt, answer string) (EssayGrade, error) {
	if t.offline {
		return offlineEssayGrade(answer), nil
	}

	system := `You are grading a learner's free-text answer to a study prompt on any subject.
Respond with ONLY a JSON object, no prose, no fences, matching:
{ "score": <number 0..1>, "feedback": "2-4 sentences" }
Be encouraging and specific: name what was correct, what was missing or wrong, and one
concrete way to improve. NEVER simply hand over the full model answer.` + t.levelClause()

	user := fmt.Sprintf("Prompt:\n%s\n\nLearner's answer:\n%s", prompt, answer)
	raw, err := t.chat(ctx, system, user)
	if err != nil {
		return EssayGrade{}, err
	}
	g, err := parseEssayGrade(raw)
	if err != nil {
		// Fall back to treating the whole reply as feedback rather than failing.
		return EssayGrade{Score: 0.5, Feedback: strings.TrimSpace(raw)}, nil
	}
	if g.Score < 0 {
		g.Score = 0
	}
	if g.Score > 1 {
		g.Score = 1
	}
	return g, nil
}

// ModelAnswer returns a reference answer to a study prompt, for when the
// learner wants to compare their attempt (or is stuck). Unlike Feedback, this
// deliberately reveals a full answer — the learner asked to see it.
func (t *Tutor) ModelAnswer(ctx context.Context, prompt string) (string, error) {
	if t.offline {
		return "I'm offline (no AI provider configured), so I can't write a model answer. " +
			"Re-read the note and compare it against your own answer instead.", nil
	}
	system := "You are a tutor writing a model answer to a study prompt, for a learner " +
		"who has attempted it and asked to see a reference. Write a clear, complete but " +
		"concise answer (1-2 short paragraphs, or a short list if the prompt calls for one). " +
		"Plain text, no markdown headers." + t.levelClause()
	return t.chat(ctx, system, "Provide a model answer to:\n"+prompt)
}

// --- parsing helpers ---

func parseNoteContent(raw string) (NoteContent, error) {
	s, ok := extractJSONObject(raw)
	if !ok {
		return NoteContent{}, fmt.Errorf("no JSON object found")
	}
	var nc NoteContent
	if err := json.Unmarshal([]byte(s), &nc); err != nil {
		return NoteContent{}, err
	}
	return nc, nil
}

func parseEssayGrade(raw string) (EssayGrade, error) {
	s, ok := extractJSONObject(raw)
	if !ok {
		return EssayGrade{}, fmt.Errorf("no JSON object found")
	}
	var g EssayGrade
	if err := json.Unmarshal([]byte(s), &g); err != nil {
		return EssayGrade{}, err
	}
	return g, nil
}

// extractJSONObject pulls the first {...} JSON object out of a model reply,
// tolerating markdown fences around it. The single fence/brace extractor for
// parseChallenge and the course parsers.
func extractJSONObject(raw string) (string, bool) {
	s := strings.TrimSpace(raw)
	if i := strings.Index(s, "```"); i >= 0 {
		s = s[i+3:]
		if nl := strings.IndexByte(s, '\n'); nl >= 0 {
			s = s[nl+1:]
		}
		if j := strings.LastIndex(s, "```"); j >= 0 {
			s = s[:j]
		}
	}
	start := strings.IndexByte(s, '{')
	end := strings.LastIndexByte(s, '}')
	if start < 0 || end <= start {
		return "", false
	}
	return s[start : end+1], true
}
