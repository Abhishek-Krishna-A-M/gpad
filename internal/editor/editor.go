package editor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
)

// Open opens path in the best available editor.
func Open(path string) error {
	cfg, _ := config.Load()

	if cfg.Editor != "" {
		return run(cfg.Editor, path)
	}
	if ed := os.Getenv("EDITOR"); ed != "" {
		return run(ed, path)
	}
	if ed := os.Getenv("VISUAL"); ed != "" {
		return run(ed, path)
	}

	for _, e := range []string{"nvim", "vim", "micro", "nano"} {
		if exists(e) {
			return exec.Command(e, path).Run()
		}
	}
	return exec.Command("nano", path).Run()
}

func run(cmdStr, file string) error {
	parts := strings.Fields(cmdStr)
	parts = append(parts, file)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func exists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
