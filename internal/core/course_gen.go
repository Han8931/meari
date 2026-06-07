package core

// course_gen.go is the agentic course-building pipeline: it turns one vault
// note (and optionally its linked neighbors) into a runnable course under
// meari-course/. The flow is plan → critique → fill → verify/repair → assemble
// → completeness check, every step reported through the progress callback:
//
//	outline  (tutor.CourseOutline)    plan modules/topics over the source notes
//	critique (tutor.ReviseOutline)    coverage / ordering / granularity, 1 round
//	fill     (tutor.CourseLesson,     write lessons for gaps, author one study
//	          tutor.CourseStudyItem)  item per topic
//	verify   code: run the reference answer against the generated tests in the
//	          REAL executor; failures go back for repair (≤2), then the topic
//	          demotes to an essay — a broken challenge never ships.
//	          essay: one judge round, regenerate once, then accept.
//	assemble write the manifest
//	complete (tutor.UncoveredIdeas)   uncovered seed ideas become an Addenda
//	                                  module, 1 round
//
// Critics and verifiers degrade, never block: any failure keeps the previous
// result. Offline skips them entirely.

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"meari/internal/executor"
	"meari/internal/tutor"
	"meari/internal/vault"
)

// CourseRequest is what the learner asked for — gathered by the intake
// conversation, or left zero for the defaults (incremental & comprehensive
// over the seed note, title picked by the model).
type CourseRequest struct {
	NotePath      string `json:"notePath"`      // the seed note (required)
	Title         string `json:"title"`         // optional: fixed course title
	Level         string `json:"level"`         // optional: beginner/intermediate/advanced
	IncludeLinked bool   `json:"includeLinked"` // also feed the seed's linked notes to the planner
	Extra         string `json:"extra"`         // free-form requirements from the intake chat
}

// GenerateCourse builds a course from req's seed note and writes it into the
// vault. progress (optional) receives one human-readable line per step.
// Existing notes referenced by the outline are reused — they only gain a
// study: block when they don't have one; their content is never rewritten.
func (s *Service) GenerateCourse(ctx context.Context, req CourseRequest, progress func(string)) (CourseMeta, error) {
	report := func(format string, args ...any) {
		if progress != nil {
			progress(fmt.Sprintf(format, args...))
		}
	}

	seed, err := s.readNote(req.NotePath)
	if err != nil {
		return CourseMeta{}, fmt.Errorf("seed note: %w", err)
	}
	notes, err := s.allNotes()
	if err != nil {
		return CourseMeta{}, err
	}
	var neighbors []vault.Note
	if req.IncludeLinked {
		neighbors = linkedNotes(seed, notes)
	}

	// 1. Plan, then critique the plan.
	report("planning the course outline from %q…", seed.Title)
	outline, err := s.tutor.CourseOutline(ctx, req.Extra, req.Level,
		noteRef(seed), noteRefs(neighbors))
	if err != nil {
		return CourseMeta{}, err
	}
	if !s.tutor.Offline() {
		report("reviewing the outline (coverage · order · granularity)…")
	}
	outline = s.tutor.ReviseOutline(ctx, outline, noteRef(seed), "")
	// The learner's explicit choices always win over the model's.
	if req.Title != "" {
		outline.Title = req.Title
	}
	if strings.TrimSpace(outline.Title) == "" {
		outline.Title = seed.Title + " Course"
	}
	if req.Level != "" {
		outline.Level = req.Level
	}
	if outline.Level == "" {
		outline.Level = "beginner"
	}
	folder := CourseDir + "/" + vault.CleanFilename(outline.Title)
	courseID := vault.Slug(outline.Title)
	return s.buildCourse(ctx, outline, courseID, folder, notes, seed, false, report)
}

// buildCourse fills, verifies, and assembles a course from its outline — the
// shared engine behind GenerateCourse and ReviseCourse. revise additionally
// maintains existing material: dead links are stripped from previously
// generated lessons and existing code exercises are re-run and repaired.
func (s *Service) buildCourse(ctx context.Context, outline tutor.CourseOutline,
	courseID, folder string, notes []vault.Note, seed vault.Note, revise bool,
	report func(string, ...any)) (CourseMeta, error) {

	total := 0
	for _, m := range outline.Modules {
		total += len(m.Topics)
	}
	report("outline ready: %q — %d modules, %d topics", outline.Title, len(outline.Modules), total)

	// Valid wikilink targets for generated lessons: existing notes plus this
	// course's own (planned) topics. Anything else gets unlinked — models
	// otherwise spray links at notes that don't exist.
	planned := map[string]bool{}
	var linkable []string
	for _, mod := range outline.Modules {
		for _, t := range mod.Topics {
			planned[strings.ToLower(strings.TrimSpace(t.Title))] = true
			linkable = append(linkable, t.Title)
		}
	}
	linkOK := func(target string) bool {
		if planned[strings.ToLower(strings.TrimSpace(target))] {
			return true
		}
		_, ok := resolveLink(target, notes)
		return ok
	}

	// Fill + verify each topic.
	var manifest strings.Builder
	if seed.Title != "" {
		manifest.WriteString("Generated from [[" + seed.Title + "]] on " +
			time.Now().Format("2006-01-02") + ".\n")
	} else {
		manifest.WriteString("Revised on " + time.Now().Format("2006-01-02") + ".\n")
	}
	var covered []string
	done := 0
	for _, mod := range outline.Modules {
		wroteHeading := false
		for _, t := range mod.Topics {
			done++
			title, err := s.buildTopic(ctx, outline.Title, folder, t, notes, seed.Body,
				fmt.Sprintf("%d/%d", done, total), linkable, linkOK, revise, report)
			if err != nil {
				return CourseMeta{}, err
			}
			if !wroteHeading {
				manifest.WriteString("\n## " + mod.Name + "\n")
				wroteHeading = true
			}
			manifest.WriteString("- [[" + title + "]]\n")
			covered = append(covered, title)
		}
	}

	// Completeness critic: uncovered seed ideas become an Addenda module.
	if seed.Body != "" {
		if missing := s.tutor.UncoveredIdeas(ctx, noteRef(seed), covered); len(missing) > 0 {
			report("completeness check: %d uncovered idea(s) — adding addenda…", len(missing))
			manifest.WriteString("\n## Addenda\n")
			for _, t := range missing { // addenda titles become linkable too
				planned[strings.ToLower(strings.TrimSpace(t.Title))] = true
				linkable = append(linkable, t.Title)
			}
			for i, t := range missing {
				title, err := s.buildTopic(ctx, outline.Title, folder, t, notes, seed.Body,
					fmt.Sprintf("+%d/%d", i+1, len(missing)), linkable, linkOK, revise, report)
				if err != nil {
					return CourseMeta{}, err
				}
				manifest.WriteString("- [[" + title + "]]\n")
			}
		}
	}

	// Assemble the manifest.
	report("assembling %s/course.md…", folder)
	_, err := s.writeNote(vault.Note{
		RelPath: folder + "/course.md",
		ID:      courseID,
		Title:   outline.Title,
		Source:  "meari-course",
		Extra:   map[string]any{"level": outline.Level},
		Body:    manifest.String(),
	})
	if err != nil {
		return CourseMeta{}, err
	}
	report("course %q is ready — switch to the tutor and :topic %s", outline.Title, courseID)
	return CourseMeta{
		ID:    courseID,
		Title: outline.Title,
		Level: outline.Level,
		Path:  folder + "/course.md",
	}, nil
}

// buildTopic produces one course topic: reuse or write its lesson note,
// author its study item, verify it, and persist. Returns the note title the
// manifest links.
func (s *Service) buildTopic(ctx context.Context, courseTitle, folder string,
	t tutor.OutlineTopic, notes []vault.Note, seedBody, label string,
	linkable []string, linkOK func(string) bool, revise bool,
	report func(string, ...any)) (string, error) {

	if err := ctx.Err(); err != nil {
		return "", err // canceled mid-build: keep what's already written
	}

	n, reused := vault.Note{}, false
	if t.UseNote != "" {
		n, reused = resolveLink(t.UseNote, notes)
	}
	changed := !reused
	if !reused {
		report("✍ %s writing lesson: %s", label, t.Title)
		body, err := s.tutor.CourseLesson(ctx, courseTitle, t, seedBody, linkable)
		if err != nil {
			return "", err
		}
		// Link verification: a wikilink either resolves (existing note or a
		// topic of this course) or it becomes plain text. Generated lessons
		// only — the learner's own notes are never rewritten.
		body, dead := stripDeadLinks(body, linkOK)
		if dead > 0 {
			report("✂ %s: unlinked %d reference(s) to notes that don't exist", label, dead)
		}
		n = vault.Note{
			RelPath: folder + "/" + vault.CleanFilename(t.Title) + ".md",
			Title:   t.Title,
			Subject: courseTitle,
			Source:  "meari-course",
			Body:    body,
		}
	} else if revise && n.Source == "meari-course" {
		// Maintenance: previously generated lessons get their links
		// re-verified too. The learner's own notes are still never rewritten.
		if body, dead := stripDeadLinks(n.Body, linkOK); dead > 0 {
			report("✂ %s: unlinked %d dead reference(s) in %s", label, dead, n.Title)
			n.Body = body
			changed = true
		}
	}

	if _, has := n.Extra["study"]; !has {
		report("✎ %s authoring exercise: %s", label, t.Title)
		item, err := s.tutor.CourseStudyItem(ctx, t.Title, t.Kind, t.Lang, n.Body)
		if err != nil {
			return "", err
		}
		item = s.verifyStudyItem(ctx, t.Title, n.Body, item, report)
		if block := studyBlock(item); block != nil {
			if n.Extra == nil {
				n.Extra = map[string]any{}
			}
			n.Extra["study"] = block
			changed = true
		}
	} else if revise {
		// Maintenance: re-run existing code exercises against the executor;
		// broken ones are repaired or demoted, like at generation time.
		if item, ok := studyItemFromBlock(n.Extra["study"]); ok && item.Kind == "code" {
			report("⚙ %s re-checking exercise: %s", label, t.Title)
			verified := s.verifyCodeItem(ctx, t.Title, item, report)
			if !reflect.DeepEqual(verified, item) {
				n.Extra["study"] = studyBlock(verified)
				changed = true
			}
		}
	}
	if changed { // reused notes without new/changed material stay untouched
		if _, err := s.writeNote(n); err != nil {
			return "", err
		}
	}
	return n.Title, nil
}

// studyItemFromBlock parses a study: frontmatter block back into a StudyItem
// (the inverse of studyBlock), for revision-time re-verification.
func studyItemFromBlock(v any) (tutor.StudyItem, bool) {
	block, ok := v.(map[string]any)
	if !ok {
		return tutor.StudyItem{}, false
	}
	get := func(k string) string {
		if val, ok := block[k]; ok {
			return fmt.Sprint(val)
		}
		return ""
	}
	item := tutor.StudyItem{
		Kind:    get("kind"),
		Lang:    get("lang"),
		Prompt:  get("prompt"),
		Starter: get("starter"),
		Answer:  get("answer"),
	}
	if tests, ok := block["tests"].([]any); ok {
		for _, t := range tests {
			item.Tests = append(item.Tests, fmt.Sprint(t))
		}
	}
	return item, true
}

// verifyStudyItem is the verification agent. Code items run against the real
// executor with a repair loop; essay items get one judge round.
func (s *Service) verifyStudyItem(ctx context.Context, title, lesson string,
	item tutor.StudyItem, report func(string, ...any)) tutor.StudyItem {

	if item.Kind == "code" {
		return s.verifyCodeItem(ctx, title, item, report)
	}
	if s.tutor.Offline() || strings.TrimSpace(item.Prompt) == "" {
		return item
	}
	if !s.tutor.JudgeEssayItem(ctx, title, lesson, item) {
		report("✗ %s: essay prompt rejected by the judge — regenerating once…", title)
		if again, err := s.tutor.CourseStudyItem(ctx, title, "essay", "", lesson); err == nil {
			return again // accepted as-is: an imperfect essay is annoying, not broken
		}
	}
	return item
}

// verifyCodeItem runs the generated reference solution against the generated
// tests — ground truth, not an LLM opinion. Failures go back to the model for
// repair (≤2 rounds); an exercise that still fails demotes to an essay, so a
// challenge that would mark correct code wrong can never reach the learner.
func (s *Service) verifyCodeItem(ctx context.Context, title string,
	item tutor.StudyItem, report func(string, ...any)) tutor.StudyItem {

	if len(item.Tests) == 0 || strings.TrimSpace(item.Answer) == "" {
		report("✗ %s: code exercise has no tests/solution — demoting to an essay", title)
		return demoteToEssay(item)
	}
	const maxRepairs = 2
	for try := 0; ; try++ {
		res, err := executor.Run(item.Lang, item.Answer, item.Tests)
		if err == nil && res.Passed {
			if try > 0 {
				report("✓ %s: repaired exercise passes its tests", title)
			}
			return item
		}
		failure := res.Output
		if err != nil {
			failure = err.Error()
		}
		if try >= maxRepairs {
			report("✗ %s: tests still failing after %d repairs — demoting to an essay", title, maxRepairs)
			return demoteToEssay(item)
		}
		report("✗ %s: reference solution fails its tests — repairing (%d/%d)…", title, try+1, maxRepairs)
		fixed, rerr := s.tutor.RepairStudyItem(ctx, title, item, failure)
		if rerr != nil {
			report("✗ %s: repair unavailable — demoting to an essay", title)
			return demoteToEssay(item)
		}
		item = fixed
	}
}

var wikilinkRE = regexp.MustCompile(`\[\[([^\]\[]+?)\]\]`)

// stripDeadLinks unlinks every [[wikilink]] whose target ok rejects, keeping
// the display text ([[X|alias]] → alias, [[X]] → X), and reports how many it
// removed.
func stripDeadLinks(body string, ok func(target string) bool) (string, int) {
	dead := 0
	out := wikilinkRE.ReplaceAllStringFunc(body, func(match string) string {
		inner := strings.TrimSpace(match[2 : len(match)-2])
		target, text := inner, inner
		if i := strings.IndexByte(inner, '|'); i >= 0 {
			target = strings.TrimSpace(inner[:i])
			text = strings.TrimSpace(inner[i+1:])
		}
		if target != "" && ok(target) {
			return match
		}
		dead++
		return text
	})
	return out, dead
}

// demoteToEssay turns a failed code exercise into an essay on the same
// prompt: the learner still writes the code, the AI grades it as prose.
func demoteToEssay(item tutor.StudyItem) tutor.StudyItem {
	return tutor.StudyItem{Kind: "essay", Prompt: item.Prompt, Answer: item.Answer}
}

// linkedNotes returns the seed's 1-hop neighborhood: the notes it links to
// plus the notes linking to it, deduped, capped to keep the planner's context
// bounded.
func linkedNotes(seed vault.Note, notes []vault.Note) []vault.Note {
	const maxNeighbors = 8
	seen := map[string]bool{seed.RelPath: true}
	var out []vault.Note
	add := func(n vault.Note) {
		if !seen[n.RelPath] && len(out) < maxNeighbors {
			seen[n.RelPath] = true
			out = append(out, n)
		}
	}
	for _, l := range vault.ParseLinks(seed.Body) {
		if n, ok := resolveLink(l.Target, notes); ok {
			add(n)
		}
	}
	for _, n := range notes { // backlinks
		if seen[n.RelPath] {
			continue
		}
		for _, l := range vault.ParseLinks(n.Body) {
			if linkMatches(l.Target, n2target(seed)) {
				add(n)
				break
			}
		}
	}
	return out
}

// studyBlock converts an authored study item into the frontmatter form the
// course loader reads. An empty essay item returns nil: the loader's defaults
// (an essay on the note) already cover it.
func studyBlock(item tutor.StudyItem) map[string]any {
	if item.Kind != "code" && strings.TrimSpace(item.Prompt) == "" {
		return nil
	}
	block := map[string]any{"kind": item.Kind}
	put := func(k, v string) {
		if strings.TrimSpace(v) != "" {
			block[k] = v
		}
	}
	put("lang", item.Lang)
	put("prompt", item.Prompt)
	put("starter", item.Starter)
	put("answer", item.Answer)
	if len(item.Tests) > 0 {
		tests := make([]any, len(item.Tests))
		for i, t := range item.Tests {
			tests[i] = t
		}
		block["tests"] = tests
	}
	return block
}

func noteRef(n vault.Note) tutor.NoteRef {
	return tutor.NoteRef{Title: n.Title, Body: n.Body}
}

func noteRefs(notes []vault.Note) []tutor.NoteRef {
	out := make([]tutor.NoteRef, len(notes))
	for i, n := range notes {
		out[i] = noteRef(n)
	}
	return out
}
