package tui

// capture.go implements ":capture" — saving the tutor Q&A about a lecture into
// the learner's OWN companion note, so questions asked while studying become
// durable personal notes instead of scrolling away.
//
// Captures never touch the lecture itself. They land in a companion note under
// "My Notes/<Lecture Title>.md" in the learner's vault, which links back to the
// lecture with a [[wikilink]] — so the lecture's backlinks panel surfaces it
// without the lecture being modified. That keeps generated course material
// pristine and safe from :revise regenerating over a learner's work.

import (
	"strings"

	"meari/internal/tutor"
	"meari/internal/vault"
)

// companionDir is the vault folder holding capture notes.
const companionDir = "My Notes"

// notesSection is the heading captures are appended under.
const notesSection = "## Notes & Questions"

// companionPath returns the vault path of the companion note for a lecture.
func companionPath(lectureTitle string) string {
	return companionDir + "/" + vault.CleanFilename(lectureTitle) + ".md"
}

// qaExchange is one question and the reply it drew.
type qaExchange struct {
	Q, A string
}

// lastExchange returns the most recent completed question/answer pair: the last
// assistant turn plus the user turn that prompted it. ok is false when the
// conversation has no answered question yet.
func lastExchange(hist []tutor.ChatTurn) (qaExchange, bool) {
	for i := len(hist) - 1; i >= 0; i-- {
		if hist[i].Role != "assistant" {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if hist[j].Role == "user" {
				return qaExchange{Q: hist[j].Content, A: hist[i].Content}, true
			}
		}
		return qaExchange{}, false // an answer with no question before it
	}
	return qaExchange{}, false
}

// allExchanges returns every completed question/answer pair, oldest first.
// Consecutive user turns keep only the last one before an answer, mirroring how
// the tutor actually replied.
func allExchanges(hist []tutor.ChatTurn) []qaExchange {
	var out []qaExchange
	question := ""
	for _, t := range hist {
		switch t.Role {
		case "user":
			question = t.Content
		case "assistant":
			if strings.TrimSpace(question) != "" {
				out = append(out, qaExchange{Q: question, A: t.Content})
				question = ""
			}
		}
	}
	return out
}

// formatExchanges renders exchanges as the markdown appended to a companion
// note: a bolded question followed by the reply.
func formatExchanges(exs []qaExchange) string {
	var b strings.Builder
	for _, ex := range exs {
		q := strings.TrimSpace(ex.Q)
		a := strings.TrimSpace(ex.A)
		if q == "" && a == "" {
			continue
		}
		b.WriteString("**Q:** " + q + "\n\n")
		b.WriteString(a + "\n\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// newCompanionBody is the starting body for a companion note: a line naming the
// lecture (the [[wikilink]] that makes this note a backlink of it) and the
// section captures append to.
func newCompanionBody(lectureTitle string) string {
	return "Notes captured while studying [[" + lectureTitle + "]].\n\n" + notesSection + "\n"
}

// appendCapture adds block under the notes section of body, creating the
// section when it is missing. Appending twice must not duplicate the heading,
// so an existing section is reused and the block goes at the end of the note.
func appendCapture(body, block string) string {
	block = strings.TrimSpace(block)
	if block == "" {
		return body
	}
	trimmed := strings.TrimRight(body, "\n")
	if !hasNotesSection(trimmed) {
		if trimmed != "" {
			trimmed += "\n\n"
		}
		trimmed += notesSection
	}
	return trimmed + "\n\n" + block + "\n"
}

// hasNotesSection reports whether body already carries the captures heading.
func hasNotesSection(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) == notesSection {
			return true
		}
	}
	return false
}
