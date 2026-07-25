// Package quotes serves a short, reflective line for the launch screen — a
// learning-themed epigraph that rotates once a day. The choice is deterministic
// for a given date (a hash of the day), so it's stable within a session and
// changes at midnight, the same trick page-sage uses for its daily fortune.
package quotes

import (
	"hash/fnv"
	"time"
)

// Quote is one epigraph: the line, who said it, and (optionally) the work it's
// from.
type Quote struct {
	Text   string
	Author string
	Work   string
}

// Daily returns the quote for the day of t. The same date always yields the
// same quote.
func Daily(t time.Time) Quote {
	if len(builtin) == 0 {
		return Quote{}
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(t.Format("2006-01-02")))
	return builtin[int(h.Sum32())%len(builtin)]
}

// builtin is the curated set — aphorisms on learning, curiosity, and knowing,
// chosen to suit a tutor that meets you at the start of a session.
var builtin = []Quote{
	{Text: "The only true wisdom is in knowing you know nothing.", Author: "Socrates"},
	{Text: "Live as if you were to die tomorrow. Learn as if you were to live forever.", Author: "Mahatma Gandhi"},
	{Text: "I have no special talent. I am only passionately curious.", Author: "Albert Einstein"},
	{Text: "The more that you read, the more things you will know. The more that you learn, the more places you'll go.", Author: "Dr. Seuss", Work: "I Can Read With My Eyes Shut!"},
	{Text: "Develop a passion for learning. If you do, you will never cease to grow.", Author: "Anthony J. D'Angelo"},
	{Text: "Tell me and I forget. Teach me and I remember. Involve me and I learn.", Author: "Benjamin Franklin"},
	{Text: "Learning never exhausts the mind.", Author: "Leonardo da Vinci"},
	{Text: "The beautiful thing about learning is that no one can take it away from you.", Author: "B.B. King"},
	{Text: "Education is not the filling of a pail, but the lighting of a fire.", Author: "W. B. Yeats"},
	{Text: "The capacity to learn is a gift; the ability to learn is a skill; the willingness to learn is a choice.", Author: "Brian Herbert"},
	{Text: "Anyone who stops learning is old, whether at twenty or eighty.", Author: "Henry Ford"},
	{Text: "Change is the end result of all true learning.", Author: "Leo Buscaglia"},
	{Text: "I am still learning.", Author: "Michelangelo"},
	{Text: "The mind is not a vessel to be filled, but a fire to be kindled.", Author: "Plutarch"},
	{Text: "It is the mark of an educated mind to be able to entertain a thought without accepting it.", Author: "Aristotle"},
	{Text: "What we learn with pleasure we never forget.", Author: "Alfred Mercier"},
	{Text: "Real knowledge is to know the extent of one's ignorance.", Author: "Confucius"},
	{Text: "He who learns but does not think, is lost. He who thinks but does not learn is in great danger.", Author: "Confucius"},
	{Text: "The expert in anything was once a beginner.", Author: "Helen Hayes"},
	{Text: "Study without desire spoils the memory, and it retains nothing that it takes in.", Author: "Leonardo da Vinci"},
	{Text: "A little learning is a dangerous thing.", Author: "Alexander Pope", Work: "An Essay on Criticism"},
	{Text: "The roots of education are bitter, but the fruit is sweet.", Author: "Aristotle"},
	{Text: "Curiosity is the wick in the candle of learning.", Author: "William Arthur Ward"},
	{Text: "Learn from yesterday, live for today, hope for tomorrow. The important thing is not to stop questioning.", Author: "Albert Einstein"},
	{Text: "An investment in knowledge pays the best interest.", Author: "Benjamin Franklin"},
	{Text: "Spoon feeding in the long run teaches us nothing but the shape of the spoon.", Author: "E. M. Forster"},
	{Text: "The whole art of teaching is only the art of awakening the natural curiosity of young minds.", Author: "Anatole France"},
	{Text: "Nothing in life is to be feared, it is only to be understood. Now is the time to understand more, so that we may fear less.", Author: "Marie Curie"},
	{Text: "Somewhere, something incredible is waiting to be known.", Author: "Carl Sagan"},
	{Text: "There is no end to education. It is not that you read a book, pass an examination, and finish with it.", Author: "Jiddu Krishnamurti"},
	{Text: "Wonder is the beginning of wisdom.", Author: "Socrates"},
	{Text: "By three methods we may learn wisdom: by reflection, the noblest; by imitation, the easiest; and by experience, the bitterest.", Author: "Confucius"},
	{Text: "The important thing is not to stop questioning. Curiosity has its own reason for existing.", Author: "Albert Einstein"},
	{Text: "Knowing yourself is the beginning of all wisdom.", Author: "Aristotle"},
}
