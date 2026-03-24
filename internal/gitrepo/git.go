package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Initialize sets up git in the notes directory with a remote.
func Initialize(path, url string) error {
	run(path, "git", "init", "-b", "main")
	run(path, "git", "branch", "-M", "main")

	remoteCheck := exec.Command("git", "-C", path, "remote", "get-url", "origin")
	if err := remoteCheck.Run(); err != nil {
		run(path, "git", "remote", "add", "origin", url)
	} else {
		run(path, "git", "remote", "set-url", "origin", url)
	}

	ignore := path + "/.gitignore"
	if _, err := os.Stat(ignore); os.IsNotExist(err) {
		_ = os.WriteFile(ignore, []byte("*.tmp\n.DS_Store\n"), 0644)
	}

	logCmd := exec.Command("git", "-C", path, "log", "--oneline", "-1")
	if logCmd.Run() != nil {
		run(path, "git", "add", ".")
		run(path, "git", "commit", "--allow-empty", "-m", "gpad: init vault")
	}

	// non-fatal pull — remote may be empty on first init
	pullCmd := exec.Command("git", "-C", path, "pull", "origin", "main",
		"--no-rebase", "--allow-unrelated-histories")
	pullOut, _ := pullCmd.CombinedOutput()
	if strings.Contains(string(pullOut), "error") {
		fmt.Println("Remote is empty or unreachable — will push on first sync.")
	}

	return nil
}

// AddCommitPush stages, commits, and pushes to main.
// All git output is suppressed — this is safe to call from a goroutine
// while the TUI is running. Errors are returned, never printed.
func AddCommitPush(path, msg string) error {
	// fix stale tracking branch
	_ = exec.Command("git", "-C", path, "branch",
		"--set-upstream-to=origin/main", "main").Run()

	// pull latest — capture output, never print it
	pullCmd := exec.Command("git", "-C", path, "pull", "origin", "main",
		"--no-rebase", "--allow-unrelated-histories")
	pullCmd.Dir = path
	if out, err := pullCmd.CombinedOutput(); err != nil {
		if strings.Contains(string(out), "CONFLICT") {
			return fmt.Errorf("merge conflict — resolve manually in %s", path)
		}
		// offline / empty remote / no-op — non-fatal, continue to push
	}

	run(path, "git", "add", ".")

	commitCmd := exec.Command("git", "-C", path, "commit", "-m", "gpad: "+msg)
	commitOut, _ := commitCmd.CombinedOutput()
	if strings.Contains(string(commitOut), "nothing to commit") {
		return nil
	}

	// push — stdout/stderr suppressed so TUI screen is never corrupted
	pushCmd := exec.Command("git", "-C", path, "push",
		"--set-upstream", "origin", "main")
	pushCmd.Stdout = nil
	pushCmd.Stderr = nil
	return pushCmd.Run()
}

// run executes a git command silently, discarding all output.
func run(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run()
}
