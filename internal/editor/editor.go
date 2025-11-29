package editor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
)

func Open(path string) error {
	cfg, _ := config.Load()

	// 1. Preferred editor from config
	if cfg.Editor != "" {
		return runCommand(cfg.Editor, path)
	}

	// 2. System $EDITOR
	if ed := os.Getenv("EDITOR"); ed != "" {
		return runCommand(ed, path)
	}

	// 3. Fallback editors
	choices := []string{"nvim", "vim", "micro", "nano"}

	for _, e := range choices {
		if exists(e) {
			return exec.Command(e, path).Run()
		}
	}

	// 4. Last fallback: raw text open (not ideal)
	return exec.Command("nano", path).Run()
}

func runCommand(cmdStr, file string) error {
	parts := strings.Split(cmdStr, " ")
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
