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

// safePrefix returns the notes root with a trailing separator,
// so /notes/daily passes but /notes-other does not.
func safePrefix() string {
	return filepath.Clean(storage.NotesDir()) + string(filepath.Separator)
}

// inVault reports whether abs is inside (or equal to) the notes root.
func inVault(abs string) bool {
	clean := filepath.Clean(abs)
	root := filepath.Clean(storage.NotesDir())
	return clean == root || strings.HasPrefix(clean, root+string(filepath.Separator))
}

// Move moves one or more notes/folders to destination.
func Move(targets []string, destRel string) error {
	destAbs := filepath.Clean(filepath.Join(storage.NotesDir(), destRel))
	if !inVault(destAbs) {
		return fmt.Errorf("destination is outside the notes directory")
	}

	for _, srcRel := range targets {
		srcAbs := filepath.Clean(filepath.Join(storage.NotesDir(), srcRel))
		if !inVault(srcAbs) {
			return fmt.Errorf("source %q is outside the notes directory", srcRel)
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
	for _, t := range targets {
		abs := filepath.Clean(filepath.Join(storage.NotesDir(), t))
		if !inVault(abs) {
			return fmt.Errorf("%q is outside the notes directory", t)
		}

		info, err := os.Stat(abs)
		if err != nil {
			return fmt.Errorf("%q not found", t)
		}

		if info.IsDir() && !recursive {
			return fmt.Errorf("%q is a directory — use -r to delete directories", t)
		}

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

	if !inVault(src) || !inVault(dest) {
		return fmt.Errorf("operation restricted to notes directory")
	}

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
