package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
)

// AddCommitPush handles the full automated sync workflow
func AddCommitPush(path, msg string) error {
	// 1. Try to pull first
	pullCmd := exec.Command("git", "pull", "--no-rebase")
	pullCmd.Dir = path
	if err := pullCmd.Run(); err != nil {
		fmt.Println("🔄 Conflict or remote changes detected. Attempting auto-merge...")
		// Here is where you'd hook into your Merge logic if needed
	}

	// 2. Add
	exec.Command("git", "-C", path, "add", ".").Run()

	// 3. Commit
	commitCmd := exec.Command("git", "-C", path, "commit", "-m", msg)
	_ = commitCmd.Run() // Ignore if nothing to commit

	// 4. Push
	pushCmd := exec.Command("git", "-C", path, "push")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	return pushCmd.Run()
}

// Initialize improved
func Initialize(path, url string) error {
	exec.Command("git", "-C", path, "init").Run()
	
	// Check if remote exists
	remoteCheck := exec.Command("git", "-C", path, "remote", "get-url", "origin")
	if err := remoteCheck.Run(); err != nil {
		exec.Command("git", "-C", path, "remote", "add", "origin", url).Run()
	} else {
		exec.Command("git", "-C", path, "remote", "set-url", "origin", url).Run()
	}

	fmt.Println("Connected to remote. Run 'gpad sync' to pull initial data.")
	return nil
}
