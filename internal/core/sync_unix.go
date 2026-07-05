//go:build unix

package core

import (
	"os/exec"
	"syscall"
)

// detachProcess sets the child process to live in its own process group
// so it survives when gpad exits.
func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
