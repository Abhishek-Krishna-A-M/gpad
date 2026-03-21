package gitrepo

import (
	"os"
	"os/exec"
)

// Pull fetches and merges from remote.
func Pull(path string) error {
	cmd := exec.Command("git", "pull", "--no-rebase")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
