package tui

// chat_store.go persists the vault TUI's per-note chat transcripts (and the
// tutor conversation history behind them) to <DataDir>/chats.json, so a
// :discuss thread survives quitting the app. Saving is best-effort — chat
// history must never block or fail the UI — and skipped entirely when no
// DataDir is configured (tests, embedded use).

import (
	"encoding/json"
	"os"
	"path/filepath"

	"meari/internal/fsutil"
	"meari/internal/tutor"
)

// chatStoreMaxBlocks bounds how much of one note's transcript is persisted;
// the newest blocks win. Keeps chats.json from growing without limit.
const chatStoreMaxBlocks = 200

type persistedBlock struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type persistedChat struct {
	Blocks []persistedBlock `json:"blocks"`
	Turns  []tutor.ChatTurn `json:"turns,omitempty"`
}

// Role names are persisted as strings so the on-disk form survives enum
// reordering. Unknown names load as roleSystem.
var roleNames = map[chatRole]string{
	roleSystem: "system", roleLesson: "lesson", roleUser: "user",
	roleTutor: "tutor", roleQuiz: "quiz", roleOK: "ok", roleFail: "fail",
}

func roleFromName(s string) chatRole {
	for r, n := range roleNames {
		if n == s {
			return r
		}
	}
	return roleSystem
}

func chatStorePath(dataDir string) string { return filepath.Join(dataDir, "chats.json") }

// loadChats reads the persisted transcripts. A missing or unreadable file is
// an empty history, never an error.
func loadChats(dataDir string) (map[string][]chatBlock, map[string][]tutor.ChatTurn) {
	blocks := map[string][]chatBlock{}
	turns := map[string][]tutor.ChatTurn{}
	if dataDir == "" {
		return blocks, turns
	}
	raw, err := os.ReadFile(chatStorePath(dataDir))
	if err != nil {
		return blocks, turns
	}
	var byNote map[string]persistedChat
	if err := json.Unmarshal(raw, &byNote); err != nil {
		return blocks, turns
	}
	for path, pc := range byNote {
		bs := make([]chatBlock, 0, len(pc.Blocks))
		for _, b := range pc.Blocks {
			bs = append(bs, chatBlock{role: roleFromName(b.Role), text: b.Text})
		}
		blocks[path] = bs
		if len(pc.Turns) > 0 {
			turns[path] = pc.Turns
		}
	}
	return blocks, turns
}

// saveChats writes every note's transcript atomically.
func saveChats(dataDir string, blocks map[string][]chatBlock, turns map[string][]tutor.ChatTurn) error {
	if dataDir == "" {
		return nil
	}
	byNote := map[string]persistedChat{}
	for path, bs := range blocks {
		if len(bs) == 0 {
			continue
		}
		if len(bs) > chatStoreMaxBlocks {
			bs = bs[len(bs)-chatStoreMaxBlocks:]
		}
		pc := persistedChat{Blocks: make([]persistedBlock, 0, len(bs)), Turns: turns[path]}
		for _, b := range bs {
			pc.Blocks = append(pc.Blocks, persistedBlock{Role: roleNames[b.role], Text: b.text})
		}
		byNote[path] = pc
	}
	b, err := json.MarshalIndent(byNote, "", " ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(chatStorePath(dataDir), b, 0o644)
}

// persistChats stashes the live transcript and saves everything. Best-effort:
// called on note switches and on quit.
func (m *VaultModel) persistChats() {
	if m.cfg.DataDir == "" {
		return
	}
	if m.current != "" {
		m.chatByNote[m.current] = m.chat.snapshot()
		m.histByNote[m.current] = m.chatHist
	}
	_ = saveChats(m.cfg.DataDir, m.chatByNote, m.histByNote)
}
