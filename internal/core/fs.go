package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Move handles moving one or more files/folders to a destination.
// It also updates Markdown headers if a note is moved/renamed.
func Move(targets []string, destRel string) error {
	notesRoot := storage.NotesDir()
	destAbs := filepath.Clean(filepath.Join(notesRoot, destRel))

	for _, srcRel := range targets {
		srcAbs := filepath.Clean(filepath.Join(notesRoot, srcRel))

		// Security check: stay inside notes dir
		if !strings.HasPrefix(srcAbs, notesRoot) || !strings.HasPrefix(destAbs, notesRoot) {
			return fmt.Errorf("operation restricted to notes directory")
		}

		// If moving to a directory, update dest path
		actualDest := destAbs
		fi, err := os.Stat(destAbs)
		if err == nil && fi.IsDir() {
			actualDest = filepath.Join(destAbs, filepath.Base(srcAbs))
		}

		if err := os.MkdirAll(filepath.Dir(actualDest), 0755); err != nil {
			return err
		}

		if err := os.Rename(srcAbs, actualDest); err != nil {
			return fmt.Errorf("failed to move %s: %w", srcRel, err)
		}

		// Smart header update for .md files
		if strings.HasSuffix(actualDest, ".md") {
			updateNoteHeader(actualDest)
		}
	}

	return autoCommit("move " + strings.Join(targets, ", "))
}

// Delete removes multiple targets
func Delete(targets []string, recursive bool) error {
	notesRoot := storage.NotesDir()
	for _, t := range targets {
		abs := filepath.Join(notesRoot, t)
		if !strings.HasPrefix(filepath.Clean(abs), notesRoot) {
			continue
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
	return autoCommit("delete targets")
}

// Internal helper to keep things DRY (Don't Repeat Yourself)
func updateNoteHeader(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	newName := strings.TrimSuffix(filepath.Base(path), ".md")
	
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			lines[i] = "# " + newName
			os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
			break
		}
	}
}

