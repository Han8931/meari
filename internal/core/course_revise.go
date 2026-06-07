package core

// course_revise.go is :revise — improve and polish an existing course. The
// outline critic re-reviews the course (steered by the learner's feedback,
// which may add, drop, reorder, or re-level topics), then the shared
// buildCourse engine rebuilds it in maintenance mode: new gaps are filled,
// dead links are stripped from previously generated lessons, and every code
// exercise is re-run against the executor and repaired or demoted. The
// course's id and folder never change, so study progress survives revision.

import (
	"context"
	"fmt"
	"strings"

	"meari/internal/tutor"
	"meari/internal/vault"
)

// ReviseCourse improves the course matching key (id, title, or path),
// following the learner's free-form feedback when given.
func (s *Service) ReviseCourse(ctx context.Context, key, feedback string, progress func(string)) (CourseMeta, error) {
	report := func(format string, args ...any) {
		if progress != nil {
			progress(fmt.Sprintf(format, args...))
		}
	}

	c, err := s.LoadCourse(key)
	if err != nil {
		return CourseMeta{}, err
	}
	notes, err := s.allNotes()
	if err != nil {
		return CourseMeta{}, err
	}
	seed := s.courseSeed(c, notes)

	outline := outlineOf(c)
	if feedback != "" {
		report("revising %q (%s)…", c.Title, feedback)
	} else {
		report("revising %q…", c.Title)
	}
	outline = s.tutor.ReviseOutline(ctx, outline, noteRef(seed), feedback)
	// Revision never renames or relevels a course away from its identity:
	// the folder, id, and title stay, so :topic and progress keep working.
	outline.Title = c.Title
	if outline.Level == "" {
		outline.Level = c.Level
	}
	folder := strings.TrimSuffix(c.Path, "/course.md")
	return s.buildCourse(ctx, outline, c.ID, folder, notes, seed, true, report)
}

// courseSeed recovers the note the course was generated from: the first
// resolvable wikilink in the manifest's preamble ("Generated from [[X]]…").
// A hand-written manifest without one yields a zero note — revision then
// skips the seed-grounded steps (outline grounding context, completeness).
func (s *Service) courseSeed(c Course, notes []vault.Note) vault.Note {
	manifest, err := s.readNote(c.Path)
	if err != nil {
		return vault.Note{}
	}
	preamble := manifest.Body
	if i := strings.Index(preamble, "\n##"); i >= 0 {
		preamble = preamble[:i]
	}
	for _, l := range vault.ParseLinks(preamble) {
		if n, ok := resolveLink(l.Target, notes); ok {
			return n
		}
	}
	return vault.Note{}
}

// outlineOf reconstructs the planner-form outline from a loaded course, so
// the critic can rework it. Every topic references its existing note; the
// study prompt stands in as the summary.
func outlineOf(c Course) tutor.CourseOutline {
	out := tutor.CourseOutline{Title: c.Title, Level: c.Level}
	for _, m := range c.Modules {
		mod := tutor.OutlineModule{Name: m.Name}
		for _, t := range m.Topics {
			mod.Topics = append(mod.Topics, tutor.OutlineTopic{
				Title:   t.Title,
				UseNote: t.Title,
				Kind:    t.Kind,
				Lang:    t.Lang,
				Summary: t.Prompt,
			})
		}
		out.Modules = append(out.Modules, mod)
	}
	return out
}
