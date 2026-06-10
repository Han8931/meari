// Package progress persists what the learner has done: which challenges are
// solved and which are in progress, plus simple attempt counts for adaptivity.
package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// State is the on-disk record. Challenges is keyed by challenge ID; Topics is
// keyed by curriculum topic ID and tracks progress through the curriculum.
type State struct {
	Challenges map[string]*Entry `json:"challenges"`
	Topics     map[string]string `json:"topics"` // topicID -> "in_progress" | "done"
	Last       *Session          `json:"last,omitempty"`
	// Completions is the course-completion ledger, keyed by course id — the
	// durable record behind the celebration card and the :achievements screen.
	Completions map[string]Completion `json:"completions,omitempty"`

	path string
}

// Completion records one finished course (every topic done).
type Completion struct {
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
	Level    string `json:"level"`
	Date     string `json:"date"` // ISO date of the FIRST completion
	Topics   int    `json:"topics"`
	FirstTry int    `json:"first_try"` // topics solved on the first attempt
	Attempts int    `json:"attempts"`  // total attempts across the course
	Flawless bool   `json:"flawless"`  // every topic solved first try
}

// Session records where the learner last was, so they can resume on relaunch.
type Session struct {
	Lang    string `json:"lang"`
	Level   string `json:"level"`
	TopicID string `json:"topic_id"`
	Title   string `json:"title"` // human-readable topic title for the resume prompt
}

// Entry tracks one challenge's status.
type Entry struct {
	Status   string `json:"status"` // "in_progress" | "done"
	Attempts int    `json:"attempts"`
	Passes   int    `json:"passes"`
}

// Load reads progress.json from dataDir, returning empty state if absent.
func Load(dataDir string) (*State, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}
	p := filepath.Join(dataDir, "progress.json")
	s := &State{Challenges: map[string]*Entry{}, Topics: map[string]string{}, Completions: map[string]Completion{}, path: p}

	b, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	if s.Challenges == nil {
		s.Challenges = map[string]*Entry{}
	}
	if s.Topics == nil {
		s.Topics = map[string]string{}
	}
	if s.Completions == nil {
		s.Completions = map[string]Completion{}
	}
	s.path = p
	return s, nil
}

// RecordCompletion logs (or refreshes) a finished course. The FIRST completion
// date is preserved across re-completions, but the stats are updated.
func (s *State) RecordCompletion(c Completion) {
	if s.Completions == nil {
		s.Completions = map[string]Completion{}
	}
	if prev, ok := s.Completions[c.CourseID]; ok && prev.Date != "" {
		c.Date = prev.Date
	}
	s.Completions[c.CourseID] = c
}

// CompletionOf returns the recorded completion for a course, if any.
func (s *State) CompletionOf(courseID string) (Completion, bool) {
	c, ok := s.Completions[courseID]
	return c, ok
}

// CompletedCourses returns the completion ledger, newest first (then by title).
func (s *State) CompletedCourses() []Completion {
	out := make([]Completion, 0, len(s.Completions))
	for _, c := range s.Completions {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Date != out[j].Date {
			return out[i].Date > out[j].Date
		}
		return out[i].Title < out[j].Title
	})
	return out
}

// TopicStatus returns "done", "in_progress", or "" for a curriculum topic.
func (s *State) TopicStatus(id string) string { return s.Topics[id] }

// MarkTopicInProgress records that the learner has started a topic (unless it is
// already done).
func (s *State) MarkTopicInProgress(id string) {
	if s.Topics[id] != "done" {
		s.Topics[id] = "in_progress"
	}
}

// MarkTopicDone records that the learner has completed a topic.
func (s *State) MarkTopicDone(id string) { s.Topics[id] = "done" }

// SetLast records the learner's current spot for resuming next session.
func (s *State) SetLast(lang, level, topicID, title string) {
	s.Last = &Session{Lang: lang, Level: level, TopicID: topicID, Title: title}
}

// Reset wipes all recorded learning history — challenge attempts, topic
// completion, and the resume point — and persists the cleared state. It backs
// the ":clear progress" command, so the change is durable, not just in-memory.
func (s *State) Reset() error {
	s.Challenges = map[string]*Entry{}
	s.Topics = map[string]string{}
	s.Last = nil
	return s.Save()
}

func (s *State) entry(id string) *Entry {
	e, ok := s.Challenges[id]
	if !ok {
		e = &Entry{}
		s.Challenges[id] = e
	}
	return e
}

// RecordAttempt increments the attempt count and, on success, the pass count
// and "done" status.
func (s *State) RecordAttempt(id string, passed bool) {
	e := s.entry(id)
	e.Attempts++
	if passed {
		e.Passes++
		e.Status = "done"
	} else if e.Status == "" {
		e.Status = "in_progress"
	}
}

// MarkInProgress flags a challenge the learner saved but hasn't solved.
func (s *State) MarkInProgress(id string) {
	e := s.entry(id)
	if e.Status != "done" {
		e.Status = "in_progress"
	}
}

// Done reports whether a challenge is solved.
func (s *State) Done(id string) bool {
	e, ok := s.Challenges[id]
	return ok && e.Status == "done"
}

// Save writes the state back to disk.
func (s *State) Save() error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
