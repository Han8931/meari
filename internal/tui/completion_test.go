package tui

import (
	"testing"

	"meari/internal/curriculum"
	"meari/internal/executor"
	"meari/internal/tutor"
)

// twoTopicCourse is a minimal in-memory curriculum for driving completion.
func twoTopicCourse() curriculum.Curriculum {
	return curriculum.Curriculum{
		Lang:  "demo",
		Level: "beginner",
		Modules: []curriculum.Module{{
			Name:   "Module",
			Topics: []curriculum.Topic{{ID: "t1", Title: "One"}, {ID: "t2", Title: "Two"}},
		}},
	}
}

// Finishing the LAST topic raises the completion card with the right stats and
// doesn't re-fire on a later pass of an already-done topic.
func TestCourseCompletionRaisesCardOnce(t *testing.T) {
	m := newModel(testDeps(t)) // no vault service → overlay only, no certificate
	m.curriculum = true
	m.curr = twoTopicCourse()
	m.deps.Progress.MarkTopicDone("t1") // first topic already done
	m.currentTopicID = "t2"

	_ = m.handleRunResult(runResultMsg{
		res: executor.Result{Passed: true}, ch: tutor.Challenge{ID: "t2"}, code: "x",
	})
	if m.overlay != overlayComplete {
		t.Fatalf("finishing the last topic should raise the completion overlay, got %d", m.overlay)
	}
	if m.completion.topics != 2 || m.completion.firstTry != 1 {
		t.Fatalf("stats = %+v, want topics 2 / firstTry 1", m.completion)
	}

	// Dismiss, then pass the (already done) topic again — no re-celebration.
	m.overlay = overlayNone
	_ = m.handleRunResult(runResultMsg{
		res: executor.Result{Passed: true}, ch: tutor.Challenge{ID: "t2"}, code: "x",
	})
	if m.overlay == overlayComplete {
		t.Fatal("re-finishing a completed course must not re-celebrate")
	}
}

// Finishing a NON-last topic does not celebrate.
func TestCourseCompletionMidCourseNoCard(t *testing.T) {
	m := newModel(testDeps(t))
	m.curriculum = true
	m.curr = twoTopicCourse()
	m.currentTopicID = "t1"

	_ = m.handleRunResult(runResultMsg{
		res: executor.Result{Passed: true}, ch: tutor.Challenge{ID: "t1"}, code: "x",
	})
	if m.overlay == overlayComplete {
		t.Fatal("finishing a mid-course topic should not celebrate")
	}
}

// With a vault service, completing a real course writes a certificate note.
func TestCourseCompletionWritesCertificate(t *testing.T) {
	d := testDepsSeeded(t)
	m := newModel(d)
	m.loadCurriculum("go-beginner", "beginner", "")
	m.curriculum = true

	topics := m.curr.Topics()
	for _, tp := range topics[:len(topics)-1] {
		m.deps.Progress.MarkTopicDone(tp.ID)
	}
	last := topics[len(topics)-1]
	m.currentTopicID = last.ID

	_ = m.handleRunResult(runResultMsg{
		res: executor.Result{Passed: true}, ch: tutor.Challenge{ID: last.ID}, code: "x",
	})
	if m.overlay != overlayComplete {
		t.Fatal("completing the course should raise the card")
	}
	if m.completion.certPath == "" {
		t.Fatal("a certificate should be written when a vault service is present")
	}
	if _, err := d.Svc.OpenNote(m.completion.certPath); err != nil {
		t.Fatalf("certificate note should exist at %q: %v", m.completion.certPath, err)
	}
	if m.completion.title != "Go (Beginner)" {
		t.Fatalf("completion title = %q, want the course title", m.completion.title)
	}
}
