//go:build integration && !unix && !windows

package integration

import (
	"os/exec"
	"time"
)

func configureTestCommand(cmd *exec.Cmd) {
	cmd.WaitDelay = 5 * time.Second
}
