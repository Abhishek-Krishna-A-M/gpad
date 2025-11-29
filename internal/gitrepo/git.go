package gitrepo

import (
	"os"
	"os/exec"
)

func Exists() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func Clone(url string, dest string) error {
	cmd := exec.Command("git", "clone", url, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitRepo(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AddCommitPush(path, msg string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", msg)
	cmd.Dir = path
	_ = cmd.Run() // ignore "nothing to commit"

	cmd = exec.Command("git", "push")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SetRemote(path, url string) error {
	cmd := exec.Command("git", "remote", "add", "origin", url)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

