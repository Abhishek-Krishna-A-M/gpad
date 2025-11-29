package notes

import (
	"os"
	"path/filepath"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

func Open(relPath string) error {
	base := storage.NotesDir()
	full := filepath.Join(base, relPath)

	if _, err := os.Stat(full); err == nil {
		return editor.Open(full)
	}

	return Create(relPath)
}

