package core

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Sync handles the "Smart Pull" logic
func Sync() error {
	cfg, _ := config.Load()
	if !cfg.GitEnabled {
		return nil // Silently skip if git isn't set up
	}

	notesDir := storage.NotesDir()
	
	// Standard pull
	err := gitrepo.Pull(notesDir)
	if err == nil {
		return nil
	}

	// If there's a conflict, you could trigger your merge logic here
	fmt.Println("Sync conflict detected. Please resolve manually in the git repo.")
	return err
}

// autoCommit is the internal helper used by Move, Delete, and Copy
func autoCommit(msg string) error {
	cfg, _ := config.Load()
	if cfg.GitEnabled && cfg.AutoPush {
		return gitrepo.AddCommitPush(storage.NotesDir(), msg)
	}
	return nil
}

// AutoSave is used by the 'open' command for background pushes
func AutoSave(msg string) {
	_ = autoCommit(msg)
}
