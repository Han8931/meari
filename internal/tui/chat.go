package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// chatRole tags a transcript block so it can be styled and (for some kinds)
// fed back to the AI as conversation history.
type chatRole int

const (
	roleSystem chatRole = iota // app notices: "— now on: … —", errors
	roleLesson                 // the lesson text
	roleUser                   // the learner's typed question
	roleTutor                  // the AI's reply / feedback
	roleOK                     // test pass line
	roleFail                   // test fail line
)

type chatBlock struct {
	role chatRole
	text string
}

// chatModel is the right pane: a scrollable transcript above a single input
// line. The transcript is kept as structured blocks so it can be re-wrapped
// whenever the pane is resized.
type chatModel struct {
	vp      viewport.Model
	input   textinput.Model
	blocks  []chatBlock
	w, h    int
	focused bool
}

func newChat() chatModel {
	in := textinput.New()
	in.Prompt = "› "
	in.Placeholder = "ask the tutor…"

	return chatModel{
		vp:    viewport.New(0, 0),
		input: in,
	}
}

// setSize lays the pane out within w×h: the last row is the input line, the
// rest is the scrollable transcript. The transcript is re-wrapped to the new
// width and kept pinned to the bottom.
func (c *chatModel) setSize(w, h int) {
	c.w, c.h = w, h
	c.input.Width = w - 2 // room for the "› " prompt
	c.vp.Width = w
	if h > 1 {
		c.vp.Height = h - 1
	} else {
		c.vp.Height = 1
	}
	c.reflow()
}

// append adds a block to the transcript, re-wraps, and scrolls to the bottom.
func (c *chatModel) append(role chatRole, text string) {
	c.blocks = append(c.blocks, chatBlock{role: role, text: strings.TrimRight(text, "\n")})
	c.reflow()
	c.vp.GotoBottom()
}

// reflow renders all blocks wrapped to the current width and loads them into
// the viewport. Called on resize and on every append.
func (c *chatModel) reflow() {
	if c.w <= 0 {
		return
	}
	var b strings.Builder
	for i, blk := range c.blocks {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(renderBlock(blk, c.w))
	}
	c.vp.SetContent(b.String())
}

// renderBlock styles one transcript block and wraps it to width w.
func renderBlock(blk chatBlock, w int) string {
	var style lipgloss.Style
	var prefix string
	switch blk.role {
	case roleUser:
		style, prefix = chatUserStyle, "you  "
	case roleTutor:
		style, prefix = chatTutorStyle, "tutor  "
	case roleLesson:
		style, prefix = chatTutorStyle, "lesson  "
	case roleOK:
		style, prefix = chatOkStyle, ""
	case roleFail:
		style, prefix = chatFailStyle, ""
	default:
		style, prefix = chatSystemStyle, ""
	}
	body := prefix + blk.text
	return style.Width(w).Render(body)
}

func (c *chatModel) focus() tea.Cmd {
	c.focused = true
	return c.input.Focus()
}

func (c *chatModel) blur() {
	c.focused = false
	c.input.Blur()
}

// submit returns the trimmed input value and clears the field. ok is false when
// the input is empty/whitespace.
func (c *chatModel) submit() (text string, ok bool) {
	v := strings.TrimSpace(c.input.Value())
	c.input.SetValue("")
	if v == "" {
		return "", false
	}
	return v, true
}

// Update routes keys to the input (when focused) and lets the viewport handle
// scrolling. The parent only forwards messages here when this pane is focused.
func (c chatModel) Update(msg tea.Msg) (chatModel, tea.Cmd) {
	var cmds []tea.Cmd

	var icmd tea.Cmd
	c.input, icmd = c.input.Update(msg)
	cmds = append(cmds, icmd)

	var vcmd tea.Cmd
	c.vp, vcmd = c.vp.Update(msg)
	cmds = append(cmds, vcmd)

	return c, tea.Batch(cmds...)
}

func (c chatModel) view() string {
	return c.vp.View() + "\n" + c.input.View()
}
