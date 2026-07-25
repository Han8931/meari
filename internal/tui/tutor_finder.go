package tui

// tutor_finder.go is the tutor's fuzzy lecture finder (",ff"): a full-screen
// overlay listing every lecture in the running course so the learner can jump
// straight to one. It reuses the vault finder's fuzzy scorer (fuzzyScore) and
// visual shape; selecting a lecture records a jump so ⌃o returns.

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// topicResult is one row of the lecture finder: the topic to open plus its
// display label and dim right-hand detail (module / progress).
type topicResult struct {
	id    string
	title string
	meta  string
}

// openLessonFinder opens the fuzzy lecture finder over the current course.
func (m Model) openLessonFinder() (tea.Model, tea.Cmd) {
	if !m.curriculum || len(m.topicByID) == 0 {
		m.flash("no lectures to jump to")
		return m, nil
	}
	m.finderMode = true
	m.finderCursor = 0
	m.finderInput.SetValue("")
	m.finderInput.CursorEnd()
	m.refreshLessonFinder()
	return m, m.finderInput.Focus()
}

// updateLessonFinder drives the finder overlay: type to filter, ↑/↓ to move,
// Enter to jump, Esc to close.
func (m Model) updateLessonFinder(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, m.quit()
	case tea.KeyEsc:
		m.finderMode = false
		m.finderInput.Blur()
		return m, nil
	case tea.KeyEnter:
		if len(m.topicResults) == 0 {
			return m, nil
		}
		if m.finderCursor < 0 || m.finderCursor >= len(m.topicResults) {
			m.finderCursor = 0
		}
		id := m.topicResults[m.finderCursor].id
		m.finderMode = false
		m.finderInput.Blur()
		t, ok := m.topicByID[id]
		if !ok {
			return m, nil
		}
		m.recordLessonJump()
		return m, m.startTopicView(t, "lesson")
	case tea.KeyUp:
		if m.finderCursor > 0 {
			m.finderCursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.finderCursor < len(m.topicResults)-1 {
			m.finderCursor++
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.finderInput, cmd = m.finderInput.Update(msg)
	m.refreshLessonFinder()
	return m, cmd
}

// refreshLessonFinder recomputes the ranked lecture list for the current query,
// preserving course order when the query is empty.
func (m *Model) refreshLessonFinder() {
	q := strings.TrimSpace(m.finderInput.Value())
	type scored struct {
		result topicResult
		score  int
		order  int
	}
	var rows []scored
	i := 0
	for _, mod := range m.curr.Modules {
		for _, t := range mod.Topics {
			meta := mod.Name
			if st := m.deps.Progress.TopicStatus(t.ID); st == "done" {
				meta = "✓ " + meta
			}
			score := 0
			if q != "" {
				var ok bool
				score, ok = fuzzyScore(q, t.Title+" "+mod.Name)
				if !ok {
					i++
					continue
				}
			}
			rows = append(rows, scored{
				result: topicResult{id: t.ID, title: t.Title, meta: meta},
				score:  score,
				order:  i,
			})
			i++
		}
	}
	sort.SliceStable(rows, func(a, b int) bool {
		if rows[a].score != rows[b].score {
			return rows[a].score > rows[b].score
		}
		return rows[a].order < rows[b].order
	})
	out := make([]topicResult, 0, min(len(rows), 40))
	for i := 0; i < len(rows) && i < 40; i++ {
		out = append(out, rows[i].result)
	}
	m.topicResults = out
	if m.finderCursor >= len(m.topicResults) {
		m.finderCursor = clampMin(len(m.topicResults)-1, 0)
	}
}

// lessonFinderView renders the centered finder overlay.
func (m Model) lessonFinderView() string {
	w := clampRange(m.width-10, 40, 92)
	if m.width < 50 {
		w = clampMin(m.width-4, 24)
	}
	inner := clampMin(w-6, 12)
	maxRows := min(len(m.topicResults), 12)

	var b strings.Builder
	b.WriteString(titleBar.Render(" Jump to lecture "))
	b.WriteString("\n\n")
	b.WriteString(m.finderInput.View())
	b.WriteString("\n\n")
	if len(m.topicResults) == 0 {
		b.WriteString(hintStyle.Render("no matching lectures"))
	} else {
		for i := 0; i < maxRows; i++ {
			r := m.topicResults[i]
			label := r.title
			if r.meta != "" {
				label += "  " + r.meta
			}
			label = truncate(label, inner-2)
			if i == m.finderCursor {
				b.WriteString(selectedRow.Width(inner).Render("▸ " + label))
			} else {
				b.WriteString("  " + label)
			}
			if i < maxRows-1 {
				b.WriteString("\n")
			}
		}
		if len(m.topicResults) > maxRows {
			b.WriteString("\n" + hintStyle.Render("  +"+itoa(len(m.topicResults)-maxRows)+" more"))
		}
	}
	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render(",ff · type to filter · enter jump · esc close"))
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2).
		Width(w).
		Render(b.String())
}
