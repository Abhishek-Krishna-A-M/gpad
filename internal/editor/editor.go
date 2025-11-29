package editor

import (
	"os"
	"os/exec"
)

func Open(path string) error {
	ed := os.Getenv("EDITOR")
	if ed == "" {
		ed = "vim" // fallback editor
	}

	cmd := exec.Command(ed, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

