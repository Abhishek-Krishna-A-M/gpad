//go:build !unix

package core

import (
	"os/exec"
)

// detachProcess is a no-op on non-Unix platforms that lack Setpgid.
func detachProcess(cmd *exec.Cmd) {}
