//go:build windows

package commandexample

import (
	"os/exec"
	"syscall"
)

const createNoWindow = 0x08000000

func configureTestCommand(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
}
