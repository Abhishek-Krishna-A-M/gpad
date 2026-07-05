package core

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

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

// AutoSave is called by 'open' after editing — errors are silently
// discarded so the TUI is never corrupted by git output.
func AutoSave(msg string) {
	_ = autoCommit(msg)
}

// AutoSaveDetached spawns a truly detached child process that stages,
// commits, and pushes changes. The process lives on even after gpad
// exits — no blocking, no goroutines, no risk of TUI corruption.
func AutoSaveDetached(msg string) {
	cfg, _ := config.Load()
	if !cfg.GitEnabled || !cfg.AutoPush {
		return
	}

	notesDir := storage.NotesDir()
	escapedDir := strings.ReplaceAll(notesDir, "'", "'\\''")
	escapedMsg := strings.ReplaceAll(msg, "'", "'\\''")

	shellCmd := fmt.Sprintf(
		"cd '%s' && git add . && git commit -m 'gpad: %s' && git push",
		escapedDir, escapedMsg,
	)

	cmd := exec.Command("sh", "-c", shellCmd)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Detach so the child lives in its own process group and survives
	// when the parent (gpad) exits.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	_ = cmd.Start() // fire-and-forget — never Wait()
}
