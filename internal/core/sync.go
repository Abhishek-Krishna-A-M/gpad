package core

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Sync pulls from remote. Silent no-op when git is not configured.
func Sync() error {
	cfg, _ := config.Load()
	if !cfg.GitEnabled {
		return nil
	}
	return gitrepo.Pull(storage.NotesDir())
}

// autoCommit stages, commits, and pushes if autopush is on.
func autoCommit(msg string) error {
	cfg, _ := config.Load()
	if !cfg.GitEnabled || !cfg.AutoPush {
		return nil
	}
	return gitrepo.AddCommitPush(storage.NotesDir(), msg)
}

// AutoSave is called by 'open' after editing — errors are swallowed so
// the user never sees a push failure block their workflow.
func AutoSave(msg string) {
	if err := autoCommit(msg); err != nil {
		fmt.Println("sync:", err)
	}
}
