package core

// mount.go mounts a SECOND markdown root — the courses directory — at the
// virtual "meari-course/" path prefix. Generated study material thereby lives
// in the app's project directory (config: [vault] course_dir, default
// <baseDir>/meari-course), not inside the learner's notes vault, while every
// front-end keeps addressing courses as ordinary vault paths: the tree shows
// meari-course/, OpenNote/SaveNote/Delete/Rename route by prefix, wikilinks
// resolve across both roots.
//
// The default mount (when SetCourseDir is never called: tests, embedded use)
// is <vault>/meari-course — the historical in-vault layout — and the mount is
// lazy: nothing creates the directory until a course is actually written.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"meari/internal/vault"
)

// SetCourseDir points the courses mount at dir (created on first write).
func (s *Service) SetCourseDir(dir string) {
	s.coursesDir = dir
}

// coursesVault opens the mount. With create false and no directory on disk,
// ok is false — listing and reading skip the mount instead of creating it.
func (s *Service) coursesVault(create bool) (*vault.Vault, bool) {
	if !create {
		if fi, err := os.Stat(s.coursesDir); err != nil || !fi.IsDir() {
			return nil, false
		}
	}
	v, err := vault.Open(s.coursesDir)
	if err != nil {
		return nil, false
	}
	return v, true
}

// inMount reports whether a virtual path addresses the courses mount, and
// strips the prefix.
func inMount(path string) (string, bool) {
	return strings.CutPrefix(path, CourseDir+"/")
}

// mountNote re-prefixes a courses-vault note into the virtual path space.
func mountNote(n vault.Note) vault.Note {
	n.RelPath = CourseDir + "/" + n.RelPath
	return n
}

// allNotes lists the notes vault and the courses mount together. Main-vault
// files under meari-course/ are served by the mount (identical virtual paths
// in the default layout), so they are skipped to avoid double listing.
func (s *Service) allNotes() ([]vault.Note, error) {
	notes, err := s.vault.List()
	if err != nil {
		return nil, err
	}
	out := notes[:0]
	for _, n := range notes {
		if _, ok := inMount(n.RelPath); !ok {
			out = append(out, n)
		}
	}
	if cv, ok := s.coursesVault(false); ok {
		cnotes, err := cv.List()
		if err != nil {
			return nil, err
		}
		for _, n := range cnotes {
			out = append(out, mountNote(n))
		}
	}
	return out, nil
}

// readNote loads a note by virtual path, from whichever root owns it.
func (s *Service) readNote(path string) (vault.Note, error) {
	if rel, ok := inMount(path); ok {
		cv, mounted := s.coursesVault(false)
		if !mounted {
			return vault.Note{}, fmt.Errorf("no course material yet (%s)", path)
		}
		n, err := cv.Read(rel)
		if err != nil {
			return vault.Note{}, err
		}
		return mountNote(n), nil
	}
	return s.vault.Read(path)
}

// writeNote saves a note by virtual path, creating the mount when needed.
func (s *Service) writeNote(n vault.Note) (vault.Note, error) {
	if rel, ok := inMount(n.RelPath); ok {
		cv, mounted := s.coursesVault(true)
		if !mounted {
			return vault.Note{}, fmt.Errorf("cannot open the course directory %s", s.coursesDir)
		}
		in := n
		in.RelPath = rel
		saved, err := cv.Write(in)
		if err != nil {
			return vault.Note{}, err
		}
		return mountNote(saved), nil
	}
	return s.vault.Write(n)
}

// deletePath removes a note or directory by virtual path.
func (s *Service) deletePath(path string) error {
	if rel, ok := inMount(path); ok {
		cv, mounted := s.coursesVault(false)
		if !mounted {
			return fmt.Errorf("nothing at %s", path)
		}
		return cv.Delete(rel)
	}
	return s.vault.Delete(path)
}

// renamePath moves a note or directory by virtual path; moves across the
// mount boundary are supported for single notes (read+write+delete).
func (s *Service) renamePath(oldPath, newPath string) error {
	oldRel, oldMounted := inMount(oldPath)
	newRel, newMounted := inMount(newPath)
	if oldMounted == newMounted {
		if !oldMounted {
			return s.vault.Rename(oldPath, newPath)
		}
		cv, ok := s.coursesVault(false)
		if !ok {
			return fmt.Errorf("nothing at %s", oldPath)
		}
		return cv.Rename(oldRel, newRel)
	}
	// Across the boundary: a markdown note round-trips; folders don't.
	n, err := s.readNote(oldPath)
	if err != nil {
		return fmt.Errorf("only single notes can move across %s/ (%w)", CourseDir, err)
	}
	n.RelPath = newPath
	if _, err := s.writeNote(n); err != nil {
		return err
	}
	return s.deletePath(oldPath)
}

// makeDirPath creates a directory by virtual path.
func (s *Service) makeDirPath(path string) error {
	if rel, ok := inMount(path); ok {
		cv, mounted := s.coursesVault(true)
		if !mounted {
			return fmt.Errorf("cannot open the course directory %s", s.coursesDir)
		}
		return cv.MakeDir(rel)
	}
	return s.vault.MakeDir(path)
}

// mountTree returns the courses mount's tree entries in virtual path space,
// including the meari-course/ root itself when the mount has content.
func (s *Service) mountTree() []TreeEntry {
	cv, ok := s.coursesVault(false)
	if !ok {
		return nil
	}
	dirs, err := cv.Dirs()
	if err != nil {
		return nil
	}
	notes, err := cv.List()
	if err != nil {
		return nil
	}
	if len(dirs)+len(notes) == 0 {
		return nil
	}
	out := []TreeEntry{{Path: CourseDir, Name: CourseDir, Dir: true}}
	for _, d := range dirs {
		out = append(out, TreeEntry{Path: CourseDir + "/" + d, Name: filepath.Base(d), Dir: true})
	}
	for _, n := range notes {
		name := strings.TrimSuffix(filepath.Base(n.RelPath), filepath.Ext(n.RelPath))
		out = append(out, TreeEntry{Path: CourseDir + "/" + n.RelPath, Name: name})
	}
	return out
}
