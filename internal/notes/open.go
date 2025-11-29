package notes

import (
	"os"
	"path/filepath"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
)

func Open(relPath string) error {
	base := storage.NotesDir()
	full := filepath.Join(base, relPath)

	if _, err := os.Stat(full); err == nil {
		if err := editor.Open(full); err != nil {
			return err
		}
		return maybeSync(full)
	}

	// Else → create
	if err := Create(relPath); err != nil {
		return err
	}
	return maybeSync(full)
}

func maybeSync(path string) error {
	cfg, _ := config.Load()
	if !cfg.GitEnabled || !cfg.AutoPush {
		return nil
	}

	notesDir := storage.NotesDir()
	return gitrepo.AddCommitPush(notesDir, "Update: "+filepath.Base(path))
}

