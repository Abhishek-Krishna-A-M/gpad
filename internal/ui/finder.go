package ui

import (
	"os"
	"path/filepath"
	"strings"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// GetAllNotes returns a slice of relative paths to all .md files
func GetAllNotes() ([]string, error) {
	var notes []string
	root := storage.NotesDir()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() && info.Name() == ".git" { return filepath.SkipDir }
		
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			rel, _ := filepath.Rel(root, path)
			notes = append(notes, rel)
		}
		return nil
	})
	
	return notes, err
}
