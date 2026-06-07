package core

// course_publish.go exports a course as a SELF-CONTAINED folder tree under a
// plain directory — meant to be (or live inside) a git repository — so
// learners can share courses. The published copy carries the manifest, every
// file already inside the course's folder, and a copy of every linked topic
// note that lives elsewhere in the vault; wikilinks resolve by title/id/stem,
// so the copies keep the outline working without rewriting the manifest.
//
// The original course is left untouched: publishing is a copy, and running
// :publish again refreshes the shared folder. A recipient drops the published
// folder into their own meari-course/ directory (or points [vault] course_dir
// at the clone) and the course loads as-is.

import (
	"fmt"
	"path/filepath"
	"strings"

	"meari/internal/vault"
)

// PublishResult reports where a course was published and how many notes the
// shared copy holds.
type PublishResult struct {
	Dir   string `json:"dir"`   // the published course folder
	Notes int    `json:"notes"` // markdown files written
}

// PublishCourse copies the course identified by key (id, title, or manifest
// path) into destDir/<course folder>. Notes are round-tripped through the
// vault format, so frontmatter — including study: blocks — survives intact.
func (s *Service) PublishCourse(key, destDir string) (PublishResult, error) {
	if strings.TrimSpace(destDir) == "" {
		return PublishResult{}, fmt.Errorf("no publish directory configured")
	}
	c, err := s.LoadCourse(key)
	if err != nil {
		return PublishResult{}, err
	}
	// meari-course/<folder>/course.md → <folder>
	parts := strings.Split(c.Path, "/")
	if len(parts) < 3 {
		return PublishResult{}, fmt.Errorf("unexpected course manifest path %q", c.Path)
	}
	folder := parts[len(parts)-2]

	dv, err := vault.Open(destDir)
	if err != nil {
		return PublishResult{}, fmt.Errorf("cannot open the publish directory %s: %w", destDir, err)
	}
	written := map[string]bool{}
	publish := func(n vault.Note, rel string) error {
		if written[rel] {
			return nil
		}
		written[rel] = true
		n.RelPath = rel
		_, err := dv.Write(n)
		return err
	}

	// 1) Everything already inside the course's folder, structure preserved —
	// the manifest, its generated lessons, and any hand-added material.
	notes, err := s.allNotes()
	if err != nil {
		return PublishResult{}, err
	}
	prefix := CourseDir + "/" + folder + "/"
	for _, n := range notes {
		if rel, ok := strings.CutPrefix(n.RelPath, prefix); ok {
			if err := publish(n, folder+"/"+rel); err != nil {
				return PublishResult{}, err
			}
		}
	}
	// 2) Topic notes linked from elsewhere in the vault, copied into the
	// folder root so the published course resolves without the home vault.
	for _, mod := range c.Modules {
		for _, t := range mod.Topics {
			if strings.HasPrefix(t.NotePath, prefix) {
				continue
			}
			n, err := s.readNote(t.NotePath)
			if err != nil {
				continue // the outline tolerates unresolvable links; so does publishing
			}
			name := vault.CleanFilename(n.Title)
			if name == "" {
				name = strings.TrimSuffix(filepath.Base(n.RelPath), filepath.Ext(n.RelPath))
			}
			if err := publish(n, folder+"/"+name+".md"); err != nil {
				return PublishResult{}, err
			}
		}
	}
	return PublishResult{Dir: filepath.Join(destDir, folder), Notes: len(written)}, nil
}
