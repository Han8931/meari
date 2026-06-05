package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"meari/internal/editor"
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

// chatInputRows is the height of the typing area. Multi-row so longer
// questions wrap and stay fully visible while being written.
const chatInputRows = 3

// chatModel is the right pane: a scrollable transcript, an optional animated
// "working…" line while an AI call is in flight, and a multi-row input. The
// transcript is kept as structured blocks so it can be re-wrapped on resize.
type chatModel struct {
	vp      viewport.Model
	input   textarea.Model
	blocks  []chatBlock
	w, h    int
	focused bool

	// busy is the label of the in-flight async op ("" = idle); busyTick drives
	// the spinner animation, advanced by the parent's spinner tick.
	busy     string
	busyTick int

	// codeLang is the language assumed for UNLABELED ``` fences in tutor and
	// lesson messages (labeled fences always win). Empty renders them plain.
	codeLang string
}

// setCodeLang sets the default language for unlabeled code fences and
// re-renders the transcript with the new highlighting.
func (c *chatModel) setCodeLang(lang string) {
	if c.codeLang == lang {
		return
	}
	c.codeLang = lang
	c.reflow()
}

func newChat() chatModel {
	in := textarea.New()
	in.Placeholder = "ask the tutor…"
	in.Prompt = "│ "
	in.ShowLineNumbers = false
	in.CharLimit = 0
	in.SetHeight(chatInputRows)

	return chatModel{
		vp:    viewport.New(0, 0),
		input: in,
	}
}

// setSize lays the pane out within w×h: the input block at the bottom, an
// optional busy line above it, and the transcript in the remaining rows.
func (c *chatModel) setSize(w, h int) {
	c.w, c.h = w, h
	c.relayout()
	c.reflow()
}

// relayout recomputes the vertical split (transcript / busy line / input) from
// the stored size and current busy state.
func (c *chatModel) relayout() {
	if c.w <= 0 || c.h <= 0 {
		return
	}
	inputH := chatInputRows
	if c.h < 7 {
		inputH = 1 // tiny panes: give the transcript what little there is
	}
	c.input.SetWidth(c.w - 2) // room for the "│ " prompt
	c.input.SetHeight(inputH)

	vpH := c.h - inputH
	if c.busy != "" {
		vpH--
	}
	if vpH < 1 {
		vpH = 1
	}
	c.vp.Width = c.w
	c.vp.Height = vpH
}

// setBusy shows (or, with "", hides) the animated progress line. The label
// names the operation, e.g. "tutor thinking".
func (c *chatModel) setBusy(label string) {
	if c.busy == label {
		return
	}
	follow := c.vp.AtBottom()
	c.busy = label
	c.relayout()
	if follow {
		c.vp.GotoBottom()
	}
}

// tickBusy advances the spinner animation one frame.
func (c *chatModel) tickBusy() { c.busyTick++ }

var busyFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (c chatModel) busyLine() string {
	frame := busyFrames[c.busyTick%len(busyFrames)]
	line := chatBusyStyle.Render(frame + " " + c.busy + "…")
	return lipgloss.NewStyle().MaxWidth(c.w).Render(line)
}

// snapshot returns the transcript blocks so a parent can stash them away when
// the learner switches topics.
func (c *chatModel) snapshot() []chatBlock { return c.blocks }

// restore replaces the transcript (nil = start fresh) and jumps to its tail —
// used when switching back to a previously visited topic.
func (c *chatModel) restore(blocks []chatBlock) {
	c.blocks = blocks
	c.reflow()
	c.vp.GotoBottom()
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
			b.WriteString("\n\n") // a blank line between blocks so turns don't run together
		}
		b.WriteString(c.renderBlock(blk))
	}
	c.vp.SetContent(b.String())
}

// renderBlock styles one transcript block and wraps it to the pane width.
// Speaker turns get a colored badge on its own line with the body below, so
// who is talking is obvious at a glance; status lines (pass/fail and app
// notices) are short and keep a single tint.
func (c chatModel) renderBlock(blk chatBlock) string {
	w := c.w
	switch blk.role {
	case roleUser:
		return chatUserBadge.Render(" you ") + "\n" + chatBodyStyle.Width(w).Render(blk.text)
	case roleTutor:
		return chatTutorBadge.Render(" tutor ") + "\n" + c.renderRichBody(blk.text)
	case roleLesson:
		return chatLessonBadge.Render(" lesson ") + "\n" + c.renderRichBody(blk.text)
	case roleOK:
		return chatOkStyle.Width(w).Render(blk.text)
	case roleFail:
		return chatFailStyle.Width(w).Render(blk.text)
	default:
		return chatSystemStyle.Width(w).Render(blk.text)
	}
}

// renderRichBody renders a tutor/lesson body: prose is word-wrapped neutrally,
// and fenced ``` code blocks are syntax-highlighted (via the editor's
// highlighter) behind a gutter bar instead of being word-wrapped.
func (c chatModel) renderRichBody(text string) string {
	lines := strings.Split(text, "\n")
	var out, prose, code []string
	lang, inCode := "", false

	flushProse := func() {
		if len(prose) > 0 {
			out = append(out, chatBodyStyle.Width(c.w).Render(strings.Join(prose, "\n")))
			prose = nil
		}
	}
	flushCode := func() {
		if len(code) == 0 {
			return
		}
		l := lang
		if l == "" {
			l = c.codeLang
		}
		if l == "" {
			l = "plain"
		}
		// Hard-wrap each highlighted line to the pane (ANSI-aware) so no code is
		// ever clipped out of view; every visual row keeps the gutter bar.
		width := c.w - 2 // room for the "│ " gutter
		if width < 4 {
			width = 4
		}
		hl := editor.Highlight(l, strings.Join(code, "\n"))
		var rows []string
		for _, row := range strings.Split(hl, "\n") {
			for _, wr := range strings.Split(ansi.Hardwrap(row, width, true), "\n") {
				rows = append(rows, chatCodeGutter.Render("│ ")+wr)
			}
		}
		out = append(out, strings.Join(rows, "\n"))
		code = nil
	}

	for _, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				flushCode()
				inCode = false
			} else {
				flushProse()
				lang = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(trimmed, "```")))
				inCode = true
			}
			continue // the fence markers themselves are not shown
		}
		if inCode {
			code = append(code, ln)
		} else {
			prose = append(prose, ln)
		}
	}
	flushProse()
	flushCode() // tolerate an unterminated fence
	return strings.Join(out, "\n")
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
	c.input.Reset()
	if v == "" {
		return "", false
	}
	return v, true
}

// Update routes input to the transcript or the input area. Scroll keys and mouse
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
	parts := make([]string, 0, 3)
	parts = append(parts, c.vp.View())
	if c.busy != "" {
		parts = append(parts, c.busyLine())
	}
	parts = append(parts, c.input.View())
	return strings.Join(parts, "\n")
}
