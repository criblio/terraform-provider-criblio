//go:build integration && windows

package integration

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func configureTestCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	cmd.Cancel = func() error {
		if cmd.Process == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			return os.ErrProcessDone
		}
		return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
	}
	cmd.WaitDelay = 5 * time.Second
}
