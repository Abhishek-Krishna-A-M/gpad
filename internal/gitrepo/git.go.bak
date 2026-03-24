package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Initialize sets up git in the notes directory with a remote.
// Works whether or not the directory already has a repo.
func Initialize(path, url string) error {
	// init (idempotent)
	run(path, "git", "init")

	// configure remote
	remoteCheck := exec.Command("git", "-C", path, "remote", "get-url", "origin")
	if err := remoteCheck.Run(); err != nil {
		run(path, "git", "remote", "add", "origin", url)
	} else {
		run(path, "git", "remote", "set-url", "origin", url)
	}

	// create .gitignore if absent
	ignore := path + "/.gitignore"
	if _, err := os.Stat(ignore); os.IsNotExist(err) {
		_ = os.WriteFile(ignore, []byte("*.tmp\n.DS_Store\n"), 0644)
	}

	// initial commit if repo has no commits yet
	logCmd := exec.Command("git", "-C", path, "log", "--oneline", "-1")
	if logCmd.Run() != nil {
		run(path, "git", "add", ".")
		run(path, "git", "commit", "--allow-empty", "-m", "gpad: init vault")
	}

	// attempt initial pull (non-fatal — remote may be empty)
	pullCmd := exec.Command("git", "-C", path, "pull", "--no-rebase", "origin")
	pullOut, _ := pullCmd.CombinedOutput()
	if strings.Contains(string(pullOut), "error") {
		fmt.Println("Remote is empty or unreachable — will push on first sync.")
	}

	return nil
}

// AddCommitPush stages everything, commits, and pushes.
// It pulls first to minimize conflicts.
func AddCommitPush(path, msg string) error {
	// pull latest
	pullCmd := exec.Command("git", "-C", path, "pull", "--no-rebase")
	pullCmd.Dir = path
	if out, err := pullCmd.CombinedOutput(); err != nil {
		outStr := string(out)
		if strings.Contains(outStr, "CONFLICT") {
			fmt.Println("⚠  Merge conflict detected — resolve manually in", path)
			return fmt.Errorf("merge conflict")
		}
		// non-fatal (offline, empty remote, etc.)
	}

	run(path, "git", "add", ".")

	commitCmd := exec.Command("git", "-C", path, "commit", "-m", "gpad: "+msg)
	commitOut, _ := commitCmd.CombinedOutput()
	if strings.Contains(string(commitOut), "nothing to commit") {
		return nil
	}

	pushCmd := exec.Command("git", "-C", path, "push", "--set-upstream", "origin", "HEAD")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	return pushCmd.Run()
}

// run executes a git command silently (errors ignored — callers decide).
func run(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	_ = cmd.Run()
}
