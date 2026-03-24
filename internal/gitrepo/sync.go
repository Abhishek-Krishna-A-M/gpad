package gitrepo

import (
	"os"
	"os/exec"
)

// Pull fetches and merges from origin/main.
// --allow-unrelated-histories handles first-time sync when the local vault
// and remote were initialised separately.
func Pull(path string) error {
	// fix any stale tracking branch still pointing to master
	_ = exec.Command("git", "-C", path, "branch", "--set-upstream-to=origin/main", "main").Run()

	cmd := exec.Command("git", "-C", path, "pull", "origin", "main",
		"--no-rebase",
		"--allow-unrelated-histories",
	)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
