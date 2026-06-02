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

// append adds a block to the transcript and re-wraps. It follows the tail only
// when the view was already pinned to the bottom, so a new message can't yank
// the reader away while they're scrolled up in the history.
func (c *chatModel) append(role chatRole, text string) {
	follow := c.vp.AtBottom()
	c.blocks = append(c.blocks, chatBlock{role: role, text: strings.TrimRight(text, "\n")})
	c.reflow()
	if follow {
		c.vp.GotoBottom()
	}
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

// Update routes input to the transcript or the text field. Scroll keys and mouse
// wheel events drive the viewport; everything else is typing. We deliberately do
// NOT forward keystrokes to the viewport's own keymap — its defaults bind j/k/f/b
// etc., which would scroll the transcript while the learner types those letters.
func (c chatModel) Update(msg tea.Msg) (chatModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		var cmd tea.Cmd
		c.vp, cmd = c.vp.Update(msg)
		return c, cmd
	case tea.KeyMsg:
		if c.scrollKey(msg) {
			return c, nil
		}
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return c, cmd
}

// scrollKey handles transcript scrolling and reports whether it consumed the key.
// Bindings are vim-flavored — Ctrl+D/U for a half page, Ctrl+F/B for a full page —
// because the focused input box owns j/k, g/G, and the plain arrows. PgUp/PgDn and
// Shift+↑/↓ (single line) are accepted too for non-vim muscle memory.
func (c *chatModel) scrollKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "ctrl+d":
		c.vp.HalfViewDown()
	case "ctrl+u":
		c.vp.HalfViewUp()
	case "ctrl+f", "pgdown":
		c.vp.ViewDown()
	case "ctrl+b", "pgup":
		c.vp.ViewUp()
	case "shift+down":
		c.vp.ScrollDown(1)
	case "shift+up":
		c.vp.ScrollUp(1)
	default:
		return false
	}
	return true
}

func (c chatModel) view() string {
	return c.vp.View() + "\n" + c.input.View()
}
