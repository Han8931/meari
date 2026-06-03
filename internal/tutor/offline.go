package tutor

// Built-in content used when no AI provider is configured. It keeps the whole
// app runnable offline so the editor loop can be exercised without a key.

func offlineLesson(topic string) string {
	return "[offline tutor — no AI provider configured]\n\n" +
		"Topic: " + topic + "\n\n" +
		"A function groups reusable code under a name. In Python you define one " +
		"with `def`, give it parameters, and return a value with `return`.\n\n" +
		"Example:\n" +
		"    def square(n):\n" +
		"        return n * n\n\n" +
		"Calling square(5) gives 25. Try the challenge below.\n" +
		"(Set OPENAI_API_KEY, or point the config at Ollama, for AI-generated lessons.)"
}

func offlineChallenge(topic string) Challenge {
	return Challenge{
		ID:     "offline-is-even",
		Prompt: "Write a function `is_even(n)` that returns True when n is even and False otherwise.",
		StarterCode: "def is_even(n):\n" +
			"    pass\n",
		Tests: []string{
			"assert is_even(4) == True",
			"assert is_even(7) == False",
			"assert is_even(0) == True",
			"assert is_even(-3) == False",
		},
	}
}

func offlineFeedback(passed bool) string {
	if passed {
		return "Nice work — all tests passed! You used a clear return value. " +
			"(Configure an AI provider for personalized feedback.)"
	}
	return "Not quite. Think about how the modulo operator `%` tells you the " +
		"remainder when dividing by 2. (Configure an AI provider for tailored hints.)"
}

// offlineNote returns a placeholder lesson note when no AI provider is set, so
// the vault + study loop is still exercisable offline.
func offlineNote(request string) NoteContent {
	topic := request
	if topic == "" {
		topic = "Untitled"
	}
	body := "_[offline tutor — no AI provider configured]_\n\n" +
		"This is a placeholder note for **" + topic + "**. " +
		"Write what you already know here, link related ideas with [[wikilinks]], " +
		"and use the study modes to test yourself.\n\n" +
		"Set `OPENAI_API_KEY`, or point the config at Ollama, to generate real lessons."
	return NoteContent{
		Title:   topic,
		Subject: "general",
		Tags:    []string{"offline"},
		Body:    body,
	}
}

// offlineEssayGrade gives generic encouragement when no provider is configured.
// A non-empty answer "passes" (score 1) so the reflection loop still advances.
func offlineEssayGrade(answer string) EssayGrade {
	if len(answer) == 0 {
		return EssayGrade{Score: 0, Feedback: "Write a short response before submitting."}
	}
	return EssayGrade{
		Score: 1,
		Feedback: "Response submitted. (Configure an AI provider for a graded, " +
			"personalized critique of your answer.)",
	}
}
