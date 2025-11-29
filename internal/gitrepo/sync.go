package gitrepo

import (
	"os"
	"os/exec"
)

func Pull(path string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Log(path string) error {
	cmd := exec.Command("git", "--no-pager", "log", "--oneline", "-n", "10")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

