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
	answerMsg    struct{ text string }
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

// answerCmd asks the tutor for a model solution to the current challenge — the
// learner explicitly asked to see it (":answer"), unlike feedback which never
// reveals one.
func answerCmd(t *tutor.Tutor, ch tutor.Challenge) tea.Cmd {
	return func() tea.Msg {
		lang := ch.Lang
		if lang == "" {
			lang = "python"
		}
		prompt := "Provide a model solution (code plus a 1-2 sentence explanation) for this " +
			lang + " exercise:\n" + ch.Prompt
		s, err := t.ModelAnswer(context.Background(), prompt)
		if err != nil {
			return errMsg{kind: "answer", err: err}
		}
		return answerMsg{text: s}
	}
}

// --- streaming chat ---

// streamChunkMsg carries one streamed chat fragment (or its terminal state)
// from the worker goroutine into the Update loop.
type streamChunkMsg struct {
	delta string // one text fragment ("" on the final message)
	full  string // the assembled reply, set when done
	done  bool
	err   error
}

// startChatStream runs fn (a streaming chat call) on a goroutine, forwarding
// each delta through the returned channel, then a final done/err message. The
// returned cmd delivers the first message; the Update loop re-arms with
// listenStream until done.
func startChatStream(ctx context.Context, fn func(context.Context, func(string)) (string, error)) (chan streamChunkMsg, tea.Cmd) {
	ch := make(chan streamChunkMsg, 64)
	go func() {
		// Every send races ctx.Done so the goroutine can never wedge on a full
		// channel after the reader (the Update loop) goes away — e.g. when the
		// learner switches modes (:vault/:tutor) or quits mid-stream. Without
		// this, a cancelled stream whose buffer filled would leak the goroutine
		// (and its open HTTP request) for the rest of the process.
		full, err := fn(ctx, func(d string) {
			select {
			case ch <- streamChunkMsg{delta: d}:
			case <-ctx.Done():
			}
		})
		select {
		case ch <- streamChunkMsg{done: true, full: full, err: err}:
		case <-ctx.Done():
		}
	}()
	return ch, listenStream(ch)
}

func listenStream(ch chan streamChunkMsg) tea.Cmd {
	return func() tea.Msg { return <-ch }
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
