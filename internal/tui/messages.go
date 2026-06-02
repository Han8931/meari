package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"meari/internal/executor"
	"meari/internal/tutor"
)

// Async result messages. Each blocking tutor/executor call is wrapped in a
// tea.Cmd that runs on its own goroutine and delivers one of these back to the
// Update loop, so the UI never blocks on the network or a subprocess.
type (
	lessonMsg    struct{ text string }
	challengeMsg struct{ ch tutor.Challenge }
	feedbackMsg  struct{ text string }
	chatReplyMsg struct{ text string }
	runResultMsg struct {
		res  executor.Result
		ch   tutor.Challenge
		code string
	}
	errMsg struct {
		kind string // which operation failed, for the status line
		err  error
	}
)

// lessonCmd fetches a topic lesson.
func lessonCmd(t *tutor.Tutor, topic string) tea.Cmd {
	return func() tea.Msg {
		s, err := t.Lesson(context.Background(), topic)
		if err != nil {
			return errMsg{kind: "lesson", err: err}
		}
		return lessonMsg{text: s}
	}
}

// challengeCmd generates the next challenge for a topic.
func challengeCmd(t *tutor.Tutor, topic string) tea.Cmd {
	return func() tea.Msg {
		ch, err := t.Challenge(context.Background(), topic)
		if err != nil {
			return errMsg{kind: "challenge", err: err}
		}
		return challengeMsg{ch: ch}
	}
}

// feedbackCmd asks the tutor for feedback on a run.
func feedbackCmd(t *tutor.Tutor, ch tutor.Challenge, code, out string, passed bool) tea.Cmd {
	return func() tea.Msg {
		s, err := t.Feedback(context.Background(), ch, code, out, passed)
		if err != nil {
			return errMsg{kind: "feedback", err: err}
		}
		return feedbackMsg{text: s}
	}
}

// chatCmd continues the free-form conversation.
func chatCmd(t *tutor.Tutor, history []tutor.ChatTurn) tea.Cmd {
	return func() tea.Msg {
		s, err := t.Chat(context.Background(), history)
		if err != nil {
			return errMsg{kind: "chat", err: err}
		}
		return chatReplyMsg{text: s}
	}
}

// runCmd executes the learner's code against the challenge's tests in the given
// language. The timeout lives inside executor.Run; running it in a Cmd goroutine
// keeps the UI responsive while it works.
func runCmd(lang, code string, ch tutor.Challenge) tea.Cmd {
	return func() tea.Msg {
		res, err := executor.Run(lang, code, ch.Tests)
		if err != nil {
			return errMsg{kind: "run", err: err}
		}
		return runResultMsg{res: res, ch: ch, code: code}
	}
}
