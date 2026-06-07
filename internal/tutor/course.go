package tutor

// Course-building capabilities: the multi-step pipeline that turns a vault
// note into a runnable course. Three calls, orchestrated by core's
// GenerateCourse: CourseOutline plans modules and topics over the source
// notes, CourseLesson writes a missing topic's lesson, and StudyItem authors
// one exercise per topic.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// NoteRef is a source note handed to the course planner.
type NoteRef struct {
	Title string
	Body  string
}

// OutlineTopic is one planned unit of study.
type OutlineTopic struct {
	Title   string `json:"title"`
	UseNote string `json:"use_note"` // an existing source note's title; "" = write a new lesson
	Kind    string `json:"kind"`     // "code" | "essay"
	Lang    string `json:"lang"`     // code topics: "python" | "go"
	Summary string `json:"summary"`  // what the lesson must cover
}

// OutlineModule groups planned topics under a heading.
type OutlineModule struct {
	Name   string         `json:"name"`
	Topics []OutlineTopic `json:"topics"`
}

// CourseOutline is the planner's result: the course shape before any content
// is written.
type CourseOutline struct {
	Title   string          `json:"title"`
	Level   string          `json:"level"`
	Modules []OutlineModule `json:"modules"`
}

// CourseOutline plans a course over the learner's notes. goal carries the
// learner's requirements (title/level/extras from the intake conversation, or
// empty for the defaults); seed is the note the course is built from and
// neighbors its linked notes. The defaults are baked into the prompt: order
// topics incrementally (fundamentals first) and cover the seed comprehensively.
func (t *Tutor) CourseOutline(ctx context.Context, goal, level string, seed NoteRef, neighbors []NoteRef) (CourseOutline, error) {
	if t.offline {
		return offlineOutline(seed, neighbors), nil
	}

	system := `You are designing a COURSE for a self-directed learner, built from their own
study notes. Respond with ONLY a JSON object, no prose, no fences, matching:
{
  "title": "a short, specific course title",
  "level": "beginner" | "intermediate" | "advanced",
  "modules": [
    { "name": "module heading",
      "topics": [
        { "title": "topic title",
          "use_note": "EXACT title of the source note that already covers this topic, or \"\" when a new lesson must be written",
          "kind": "code" | "essay",
          "lang": "python" | "go" | "",
          "summary": "1-2 sentences: what the lesson must cover" } ] } ]
}
Rules:
- DERIVE the modules and topics from the source note's ACTUAL content — its
  sections, claims, examples, and code. This is a course about what the
  learner wrote down, NOT a generic textbook syllabus for the subject. Every
  topic summary must name specific material from the note it builds on.
- 2-4 modules, 2-5 topics each, ordered INCREMENTALLY: fundamentals first, each
  topic building on the previous ones.
- Be COMPREHENSIVE about the primary source note: every significant idea in it
  must be covered by some topic.
- Prefer use_note when a source note already covers a topic well; plan a new
  lesson (use_note "") only for real gaps.
- kind "code" (with lang "python" or "go") ONLY for programming topics where
  writing code is the right exercise; everything else is "essay".` + t.levelClause()

	var b strings.Builder
	if goal != "" {
		b.WriteString("Learner requirements: " + goal + "\n")
	}
	if level != "" {
		b.WriteString("Required level: " + level + "\n")
	}
	b.WriteString("\nPrimary source note: " + seed.Title + "\n---\n" + clampRunes(seed.Body, 6000) + "\n---\n")
	for _, n := range neighbors {
		b.WriteString("\nLinked source note: " + n.Title + "\n---\n" + clampRunes(n.Body, 2500) + "\n---\n")
	}

	raw, err := t.chat(ctx, system, b.String())
	if err != nil {
		return CourseOutline{}, err
	}
	s, ok := extractJSONObject(raw)
	if !ok {
		return CourseOutline{}, fmt.Errorf("course outline: no JSON object in reply")
	}
	var out CourseOutline
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return CourseOutline{}, fmt.Errorf("course outline: %w", err)
	}
	if len(out.Modules) == 0 {
		return CourseOutline{}, fmt.Errorf("course outline: no modules planned")
	}
	return out, nil
}

// CourseLesson writes the lesson body (markdown) for a planned topic that no
// existing note covers. linkable names the only valid [[wikilink]] targets
// (this course's topics and existing notes) — anything else must stay plain
// text, and the caller verifies that with stripDeadLinks regardless.
func (t *Tutor) CourseLesson(ctx context.Context, courseTitle string, topic OutlineTopic, sourceContext string, linkable []string) (string, error) {
	if t.offline {
		return offlineNote(topic.Title).Body, nil
	}
	system := `You are writing one LESSON of the course "` + courseTitle + `" as a markdown
note for a self-directed learner. The course is built FROM THE LEARNER'S OWN NOTE:
ground the lesson in that source material — keep its terminology, build on its
examples and code, and deepen what it says rather than writing generic prose about
the subject. Extend beyond the note only where the topic requires it.
A few short sections of clear prose, optionally one worked example. Do not restate
the title as an H1 — start with the explanation.
You may wrap a concept in [[wikilinks]] ONLY when it is one of these existing notes:
` + clampRunes(strings.Join(linkable, ", "), 1500) + `
Everything else stays plain text — never invent a link. Respond with ONLY the markdown body.` + t.levelClause()
	user := "Topic: " + topic.Title + "\nIt must cover: " + topic.Summary
	if sourceContext != "" {
		user += "\n\nTHE SOURCE NOTE (the learner's own material — your starting point):\n---\n" +
			clampRunes(sourceContext, 4000) + "\n---"
	}
	return t.chat(ctx, system, user)
}

// StudyItem is one authored exercise for a course topic.
type StudyItem struct {
	Kind    string   `json:"kind"`
	Lang    string   `json:"lang"`
	Prompt  string   `json:"prompt"`
	Starter string   `json:"starter"`
	Tests   []string `json:"tests"`
	Answer  string   `json:"answer"`
}

// CourseStudyItem authors the exercise for one topic: a code challenge with
// hidden tests for kind "code", an essay prompt with a model answer otherwise.
func (t *Tutor) CourseStudyItem(ctx context.Context, topicTitle, kind, lang, lesson string) (StudyItem, error) {
	if t.offline {
		return StudyItem{Kind: "essay"}, nil // core fills the essay defaults
	}

	var system string
	if kind == "code" {
		if lang == "" {
			lang = "python"
		}
		system = fmt.Sprintf(`You are authoring ONE %s coding exercise for the lesson below.
Respond with ONLY a JSON object, no prose, no fences, matching:
{
  "kind": "code",
  "lang": %q,
  "prompt": "what to implement, 2-4 sentences",
  "starter": "a code stub the learner completes",
  "tests": ["each entry one %s assertion against the learner's code"],
  "answer": "a correct reference solution that passes every test"
}
The tests must be self-contained and deterministic.`, lang, lang, lang)
	} else {
		system = `You are authoring ONE essay exercise for the lesson below.
Respond with ONLY a JSON object, no prose, no fences, matching:
{
  "kind": "essay",
  "prompt": "a question that makes the learner actively reconstruct the lesson's key ideas (not yes/no)",
  "answer": "a complete model answer, 1-2 short paragraphs"
}` + t.levelClause()
	}

	raw, err := t.chat(ctx, system, "Topic: "+topicTitle+"\n\nLesson:\n"+clampRunes(lesson, 5000))
	if err != nil {
		return StudyItem{}, err
	}
	s, ok := extractJSONObject(raw)
	if !ok {
		return StudyItem{}, fmt.Errorf("study item: no JSON object in reply")
	}
	var item StudyItem
	if err := json.Unmarshal([]byte(s), &item); err != nil {
		return StudyItem{}, fmt.Errorf("study item: %w", err)
	}
	if item.Kind == "" {
		item.Kind = kind
	}
	return item, nil
}

// ReviseOutline is the outline critic: it checks a planned outline against
// the seed note for coverage (every significant idea addressed), grounding,
// ordering (incremental, no forward dependencies), and granularity (one idea
// per topic), returning a revised outline. feedback (optional) carries the
// learner's own revision wishes (:revise) and licenses restructuring. Any
// failure returns the original — critics improve a course, they never block.
func (t *Tutor) ReviseOutline(ctx context.Context, outline CourseOutline, seed NoteRef, feedback string) CourseOutline {
	if t.offline {
		return outline
	}
	planned, err := json.Marshal(outline)
	if err != nil {
		return outline
	}
	system := `You are reviewing a planned course outline against its source note. Check:
1. COVERAGE — every significant idea in the source note is addressed by a topic.
2. GROUNDING — every topic traces to specific content in the source note; replace
   generic textbook-syllabus filler with topics about what the note actually says.
3. ORDER — topics are incremental; nothing depends on a later topic.
4. GRANULARITY — one idea per topic; split topics that span several.
Respond with ONLY the corrected outline as a JSON object in exactly the same schema
you received (title/level/modules/topics with use_note, kind, lang, summary). If the
outline is already good, return it unchanged. Keep use_note values exactly as given
for topics you keep; new topics use use_note "".`
	if feedback != "" {
		system += "\n\nThe learner asked for this revision — follow it, restructuring (adding," +
			" removing, reordering, re-leveling topics) as it requires:\n" + feedback
	}
	user := "Source note: " + seed.Title + "\n---\n" + clampRunes(seed.Body, 6000) +
		"\n---\n\nPlanned outline JSON:\n" + string(planned)
	raw, err := t.chat(ctx, system, user)
	if err != nil {
		return outline
	}
	s, ok := extractJSONObject(raw)
	if !ok {
		return outline
	}
	var revised CourseOutline
	if err := json.Unmarshal([]byte(s), &revised); err != nil || len(revised.Modules) == 0 {
		return outline
	}
	return revised
}

// RepairStudyItem feeds an executor failure back to the model so it can fix
// the reference solution or the tests (whichever is wrong).
func (t *Tutor) RepairStudyItem(ctx context.Context, topicTitle string, item StudyItem, failure string) (StudyItem, error) {
	if t.offline {
		return item, fmt.Errorf("offline: cannot repair")
	}
	current, err := json.Marshal(item)
	if err != nil {
		return item, err
	}
	system := `Your reference solution for a coding exercise FAILED its own tests. Fix the
solution or the tests — whichever is wrong — so the solution passes every test.
Respond with ONLY the corrected exercise as a JSON object in the same schema
(kind/lang/prompt/starter/tests/answer). Keep the prompt's intent.`
	user := "Topic: " + topicTitle + "\n\nCurrent exercise JSON:\n" + string(current) +
		"\n\nTest run output:\n" + clampRunes(failure, 2000)
	raw, err := t.chat(ctx, system, user)
	if err != nil {
		return item, err
	}
	s, ok := extractJSONObject(raw)
	if !ok {
		return item, fmt.Errorf("repair: no JSON object in reply")
	}
	var fixed StudyItem
	if err := json.Unmarshal([]byte(s), &fixed); err != nil {
		return item, err
	}
	if fixed.Kind == "" {
		fixed.Kind = item.Kind
	}
	if fixed.Lang == "" {
		fixed.Lang = item.Lang
	}
	return fixed, nil
}

// JudgeEssayItem checks an essay exercise: answerable from the lesson alone,
// and actually answered by its model answer. Uncertainty counts as ok — an
// imperfect essay prompt is annoying, not broken.
func (t *Tutor) JudgeEssayItem(ctx context.Context, topicTitle, lesson string, item StudyItem) bool {
	if t.offline {
		return true
	}
	system := `You are checking ONE essay exercise. Respond with ONLY a JSON object:
{ "ok": true|false, "reason": "one sentence" }
ok is false ONLY when the prompt cannot be answered from the lesson alone, or the
model answer does not actually answer the prompt. When unsure, ok is true.`
	user := "Topic: " + topicTitle + "\n\nLesson:\n" + clampRunes(lesson, 4000) +
		"\n\nPrompt: " + item.Prompt + "\n\nModel answer:\n" + clampRunes(item.Answer, 2000)
	raw, err := t.chat(ctx, system, user)
	if err != nil {
		return true
	}
	s, ok := extractJSONObject(raw)
	if !ok {
		return true
	}
	var verdict struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal([]byte(s), &verdict); err != nil {
		return true
	}
	return verdict.OK
}

// UncoveredIdeas is the completeness critic: after the course is filled, it
// names the seed note's significant ideas no topic covers, as ready-to-fill
// topics. An empty result (or any failure) means "complete".
func (t *Tutor) UncoveredIdeas(ctx context.Context, seed NoteRef, coveredTopics []string) []OutlineTopic {
	if t.offline {
		return nil
	}
	system := `You are checking a finished course for completeness against its source note.
Respond with ONLY a JSON object:
{ "missing": [ { "title": "...", "use_note": "", "kind": "code"|"essay", "lang": "python"|"go"|"", "summary": "what the lesson must cover" } ] }
List ONLY significant ideas from the source note that none of the covered topics
address. Minor details and asides do not count. An empty list is the expected
answer for a well-built course.`
	user := "Source note: " + seed.Title + "\n---\n" + clampRunes(seed.Body, 6000) +
		"\n---\n\nCovered topics:\n- " + strings.Join(coveredTopics, "\n- ")
	raw, err := t.chat(ctx, system, user)
	if err != nil {
		return nil
	}
	s, ok := extractJSONObject(raw)
	if !ok {
		return nil
	}
	var out struct {
		Missing []OutlineTopic `json:"missing"`
	}
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	const maxAddenda = 4 // a long tail means the outline was bad, not the course incomplete
	if len(out.Missing) > maxAddenda {
		out.Missing = out.Missing[:maxAddenda]
	}
	return out.Missing
}

// offlineOutline degrades gracefully: one module of essay topics, one per
// source note, so :course still yields a runnable course with no provider.
func offlineOutline(seed NoteRef, neighbors []NoteRef) CourseOutline {
	mod := OutlineModule{Name: "Topics"}
	for _, n := range append([]NoteRef{seed}, neighbors...) {
		mod.Topics = append(mod.Topics, OutlineTopic{
			Title: n.Title, UseNote: n.Title, Kind: "essay",
		})
	}
	return CourseOutline{
		Title:   seed.Title + " Course",
		Level:   "beginner",
		Modules: []OutlineModule{mod},
	}
}

// clampRunes bounds s to n runes for model-friendly context sizes.
func clampRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "\n…(truncated)"
}
