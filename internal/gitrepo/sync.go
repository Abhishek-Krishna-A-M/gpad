package gitrepo

import (
	"os"
	"os/exec"
)

// Pull fetches and merges from origin/main.
// Also fixes any stale tracking branch pointing to master.
func Pull(path string) error {
	// fix stale tracking config that still points to master
	fixCmd := exec.Command("git", "-C", path, "branch", "--set-upstream-to=origin/main", "main")
	_ = fixCmd.Run()

	cmd := exec.Command("git", "-C", path, "pull", "origin", "main", "--no-rebase")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
