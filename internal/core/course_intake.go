package core

// course_intake.go supports the :course requirements interview. The intake is
// its own conversation — CourseIntakeStream gives the interview instructions
// to the model as THE system prompt (riding the regular tutor chat demoted
// them to ignorable context, and the model would drift into tutoring instead
// of interviewing). The front-end watches each reply for the course_request
// JSON that ends the interview (ParseCourseRequest).

import (
	"context"
	"encoding/json"
	"strings"

	"meari/internal/tutor"
)

// CourseIntakeStream continues the requirements interview for building a
// course from the note at seedPath, streaming the reply.
func (s *Service) CourseIntakeStream(ctx context.Context, seedPath string, history []tutor.ChatTurn, onDelta func(string)) (string, error) {
	return s.tutor.StreamConversation(ctx, s.courseIntakeSystem(seedPath), TrimTurns(history), onDelta)
}

// courseIntakeSystem builds the interview's system prompt: grounded in the
// seed note's actual content, forbidden from tutoring, and forced to converge
// on the course_request JSON within a couple of exchanges.
func (s *Service) courseIntakeSystem(seedPath string) string {
	title, body, linked := seedPath, "", "none"
	if seed, err := s.readNote(seedPath); err == nil {
		title, body = seed.Title, ClampContext(seed.Body)
		if notes, err := s.allNotes(); err == nil {
			var titles []string
			for _, n := range linkedNotes(seed, notes) {
				titles = append(titles, n.Title)
			}
			if len(titles) > 0 {
				linked = strings.Join(titles, ", ")
			}
		}
	}
	return `You are running a SHORT requirements interview to configure an automated
course builder. The course will be built from the learner's note "` + title + `".
You are NOT teaching: never explain the material, never answer questions about it —
only settle how the course should be built.

The note's content, so your questions can reference what it actually covers:
---
` + body + `
---
Linked notes that could also be included: ` + linked + `.

Interview rules:
- FIRST message: at most three short numbered questions — (1) difficulty:
  beginner / intermediate / advanced, (2) scope: only this note, or also the
  linked notes, (3) a course title, or you pick one. Mention they can answer
  any subset or just say "defaults".
- Settle everything within TWO exchanges. Vague or partial answers: decide the
  rest yourself. Defaults: incremental ordering, comprehensive coverage of the
  note, you pick the title and difficulty.
- Once settled (or the learner says "defaults"), reply with ONE line summarizing
  what you will build, then ON ITS OWN LINE exactly this JSON (no code fences):
{"course_request": {"title": "<title or empty>", "level": "<beginner|intermediate|advanced or empty>", "includeLinked": <true|false>, "extra": "<other requirements, or empty>"}}
- Never emit the JSON in your first message. Keep every message short.`
}

// ParseCourseRequest detects the course_request JSON that ends the intake
// conversation, anywhere in the reply.
func ParseCourseRequest(reply string) (CourseRequest, bool) {
	i := strings.Index(reply, `"course_request"`)
	if i < 0 {
		return CourseRequest{}, false
	}
	start := strings.LastIndexByte(reply[:i], '{')
	if start < 0 {
		return CourseRequest{}, false
	}
	depth := 0
	for j := start; j < len(reply); j++ {
		switch reply[j] {
		case '{':
			depth++
		case '}':
			if depth--; depth == 0 {
				var w struct {
					Req CourseRequest `json:"course_request"`
				}
				if err := json.Unmarshal([]byte(reply[start:j+1]), &w); err != nil {
					return CourseRequest{}, false
				}
				return w.Req, true
			}
		}
	}
	return CourseRequest{}, false
}
