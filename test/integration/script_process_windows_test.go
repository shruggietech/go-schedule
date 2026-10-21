//go:build windows

package integration

import (
	"os/exec"
	"syscall"
)

func configureScriptProcess(cmd *exec.Cmd) {
	const createNoWindow = 0x08000000
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
}
