//go:build !windows

package executor

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/processgroup"
)

type unixOwned struct {
	cmd  *exec.Cmd
	pgid int
}

func startOwnedProcess(ctx context.Context, cmd *exec.Cmd) (ownedProcess, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &unixOwned{cmd: cmd, pgid: cmd.Process.Pid}, nil
}
func (p *unixOwned) wait() nativeResult {
	err := p.cmd.Wait()
	r := nativeResult{err: err}
	if p.cmd.ProcessState != nil {
		code := p.cmd.ProcessState.ExitCode()
		r.code = &code
	}
	return r
}
func (p *unixOwned) terminate(deadline time.Time) error {
	return processgroup.Terminate(p.pgid, max(0, time.Until(deadline)))
}
func (p *unixOwned) empty() (bool, error) {
	err := syscall.Kill(-p.pgid, 0)
	if errors.Is(err, syscall.ESRCH) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect owned process group: %w", err)
	}
	return false, nil
}
func (p *unixOwned) release() error { return nil }
func isBrokenInput(err error) bool  { return errors.Is(err, syscall.EPIPE) }
