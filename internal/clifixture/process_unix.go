//go:build !windows

package clifixture

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

func startContained(command Command, input, stdout, stderr *os.File) (*Process, func() error, error) {
	cmd := exec.Command(command.Executable, command.Arguments...)
	cmd.Dir = command.Directory
	cmd.Stdin = input
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if command.Environment != nil {
		cmd.Env = command.Environment
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	p := &Process{PID: cmd.Process.Pid}
	p.stop = func() error {
		err := syscall.Kill(-p.PID, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	return p, cmd.Wait, nil
}
