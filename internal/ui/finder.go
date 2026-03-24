package ui

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/index"
)

// GetAllNotes returns relative paths to all .md files in the vault.
// Uses the persistent index cache — O(changed notes) not O(vault).
func GetAllNotes() ([]string, error) {
	notes := index.AllNotes()
	paths := make([]string, 0, len(notes))
	for _, n := range notes {
		paths = append(paths, n.RelPath)
	}
	return paths, nil
}
