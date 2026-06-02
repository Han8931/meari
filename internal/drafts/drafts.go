// Package drafts persists in-progress solutions so a learner can stop and
// resume a challenge later. Drafts are keyed by challenge ID.
package drafts

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Store reads and writes drafts under a directory.
type Store struct {
	dir string
}

// New returns a Store rooted at <workspaceDir>/drafts, creating it if needed.
func New(workspaceDir string) (*Store, error) {
	dir := filepath.Join(workspaceDir, "drafts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

var unsafe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func (s *Store) path(challengeID string) string {
	safe := unsafe.ReplaceAllString(challengeID, "_")
	safe = strings.Trim(safe, "_")
	if safe == "" {
		safe = "untitled"
	}
	return filepath.Join(s.dir, safe+".py")
}

// Load returns the saved draft for a challenge, or ("", false) if none exists.
func (s *Store) Load(challengeID string) (string, bool) {
	b, err := os.ReadFile(s.path(challengeID))
	if err != nil {
		return "", false
	}
	return string(b), true
}

// Save writes the in-progress code for a challenge.
func (s *Store) Save(challengeID, code string) error {
	return os.WriteFile(s.path(challengeID), []byte(code), 0o644)
}

// Clear removes a draft (e.g. after the challenge is solved). Missing is OK.
func (s *Store) Clear(challengeID string) error {
	err := os.Remove(s.path(challengeID))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Has reports whether a draft exists for a challenge.
func (s *Store) Has(challengeID string) bool {
	_, ok := s.Load(challengeID)
	return ok
}

// IDs returns the challenge IDs that have a saved draft on disk, so the UI can
// repopulate the learner's work across restarts. The IDs are the draft file
// names with the ".py" suffix stripped.
func (s *Store) IDs() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".py") {
			ids = append(ids, strings.TrimSuffix(name, ".py"))
		}
	}
	return ids, nil
}
