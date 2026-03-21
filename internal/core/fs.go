package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Move moves one or more notes/folders to destination.
func Move(targets []string, destRel string) error {
	notesRoot := storage.NotesDir()
	destAbs := filepath.Clean(filepath.Join(notesRoot, destRel))

	for _, srcRel := range targets {
		srcAbs := filepath.Clean(filepath.Join(notesRoot, srcRel))

		if !strings.HasPrefix(srcAbs, notesRoot) || !strings.HasPrefix(destAbs, notesRoot) {
			return fmt.Errorf("operation restricted to notes directory")
		}

		actualDest := destAbs
		if fi, err := os.Stat(destAbs); err == nil && fi.IsDir() {
			actualDest = filepath.Join(destAbs, filepath.Base(srcAbs))
		}

		if err := os.MkdirAll(filepath.Dir(actualDest), 0755); err != nil {
			return err
		}
		if err := os.Rename(srcAbs, actualDest); err != nil {
			return fmt.Errorf("move %s: %w", srcRel, err)
		}

		if strings.HasSuffix(actualDest, ".md") {
			updateHeader(actualDest)
		}
	}
	return autoCommit("move " + strings.Join(targets, ", "))
}

// Delete removes targets (files or dirs with -r).
func Delete(targets []string, recursive bool) error {
	notesRoot := storage.NotesDir()
	for _, t := range targets {
		abs := filepath.Clean(filepath.Join(notesRoot, t))
		if !strings.HasPrefix(abs, notesRoot) {
			return fmt.Errorf("operation restricted to notes directory")
		}
		var err error
		if recursive {
			err = os.RemoveAll(abs)
		} else {
			err = os.Remove(abs)
		}
		if err != nil {
			return err
		}
	}
	return autoCommit("delete " + strings.Join(targets, ", "))
}

// Copy duplicates a note.
func Copy(srcRel, destRel string) error {
	notesRoot := storage.NotesDir()
	src := filepath.Join(notesRoot, srcRel)
	dest := filepath.Join(notesRoot, destRel)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return autoCommit("copy " + srcRel + " → " + destRel)
}

func updateHeader(path string) {
	meta, body, err := frontmatter.Parse(path)
	if err != nil {
		return
	}
	newTitle := strings.TrimSuffix(filepath.Base(path), ".md")
	if meta.Title == newTitle {
		return
	}
	meta.Title = newTitle
	_ = frontmatter.Write(path, meta, body)
}
