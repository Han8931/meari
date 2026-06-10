package core

// certificate.go turns finishing a course into an ARTIFACT the learner owns: a
// markdown "certificate of completion" written into the course's own folder.
// Because it's an ordinary vault note, it shows in the tree, survives in git,
// and can be :publish-ed — the reward is a file, not a transient popup.

import (
	"fmt"
	"strings"

	"meari/internal/vault"
)

// Certificate is the data a completion certificate records. The TUI fills it in
// from saved progress (Date is stamped by the caller for testability).
type Certificate struct {
	CourseTitle string
	Level       string
	Date        string // ISO date, e.g. "2026-06-10"
	Topics      int    // topics finished
	FirstTry    int    // topics solved on the first attempt
	Attempts    int    // total attempts across the course
	Flawless    bool   // every topic solved first try
}

// IssueCertificate writes (or refreshes) certificate.md in the finished
// course's folder and returns its metadata. courseKey is an id, title, or
// manifest path — the same keys LoadCourse accepts.
func (s *Service) IssueCertificate(courseKey string, c Certificate) (NoteMeta, error) {
	course, err := s.LoadCourse(courseKey)
	if err != nil {
		return NoteMeta{}, err
	}
	parts := strings.Split(course.Path, "/")
	if len(parts) < 2 {
		return NoteMeta{}, fmt.Errorf("unexpected course path %q", course.Path)
	}
	folder := strings.Join(parts[:len(parts)-1], "/") // drop the trailing "/course.md"
	note := vault.Note{
		RelPath: folder + "/certificate.md",
		Title:   c.CourseTitle + " — Certificate",
		Source:  "imported:certificate",
		Body:    renderCertificate(c),
	}
	saved, err := s.writeNote(note)
	if err != nil {
		return NoteMeta{}, err
	}
	return metaOf(saved), nil
}

// renderCertificate builds the certificate's markdown body.
func renderCertificate(c Certificate) string {
	var b strings.Builder
	b.WriteString("# 🎓 Certificate of Completion\n\n")
	b.WriteString("*Meari · 메아리*\n\n---\n\n")

	title := c.CourseTitle
	if c.Level != "" {
		title += " (" + c.Level + ")"
	}
	fmt.Fprintf(&b, "## %s\n\n", title)
	if c.Date != "" {
		fmt.Fprintf(&b, "Completed on **%s**.\n\n", c.Date)
	}
	fmt.Fprintf(&b, "- Topics finished: **%d**\n", c.Topics)
	if c.Topics > 0 {
		fmt.Fprintf(&b, "- Solved on the first try: **%d / %d**\n", c.FirstTry, c.Topics)
	}
	if c.Attempts > 0 {
		fmt.Fprintf(&b, "- Total attempts: **%d**\n", c.Attempts)
	}
	if c.Flawless {
		b.WriteString("- ⭐ **Flawless run** — every topic on the first try!\n")
	}
	b.WriteString("\n> 메아리 — the echo returns. What you learned has come back to you.\n")
	return b.String()
}
