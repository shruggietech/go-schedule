//go:build !windows

package clifixture

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/shruggietech/go-schedule/internal/processgroup"
)

func startContained(command Command, input, stdout, stderr *os.File) (*Process, func() error, error) {
	if command.Limits != nil {
		return nil, nil, errors.New("environment limits on Unix require a verified owned cgroup worker")
	}
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
	p.stop = func() error { return processgroup.Terminate(p.PID, 5*time.Second) }
	return p, cmd.Wait, nil
}
