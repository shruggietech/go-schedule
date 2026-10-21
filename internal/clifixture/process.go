package clifixture

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Command describes one hidden, noninteractive owned child.
type Command struct {
	Executable                        string
	Arguments, Environment            []string
	Directory, StdoutPath, StderrPath string
	Limits                            *ProcessLimits
}

// Process retains the exact process-tree capability created by StartProcess.
type Process struct {
	PID     int
	done    chan struct{}
	result  error
	stop    func() error
	once    sync.Once
	stopErr error
}

// StartProcess establishes process-tree containment before child execution.
// Call Stop during cleanup even when the direct child has already exited.
func StartProcess(command Command) (*Process, error) {
	if err := validateProcessLimits(command.Limits); err != nil {
		return nil, err
	}
	if !filepath.IsAbs(command.Executable) || !filepath.IsAbs(command.Directory) || !filepath.IsAbs(command.StdoutPath) || !filepath.IsAbs(command.StderrPath) {
		return nil, errors.New("owned child requires absolute executable, directory and output paths")
	}
	stdout, err := os.OpenFile(command.StdoutPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	stderr, err := os.OpenFile(command.StderrPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		_ = stdout.Close()
		return nil, err
	}
	input, err := os.Open(os.DevNull)
	if err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		return nil, err
	}
	p, wait, err := startContained(command, input, stdout, stderr)
	if err != nil {
		_ = input.Close()
		_ = stdout.Close()
		_ = stderr.Close()
		return nil, err
	}
	p.done = make(chan struct{})
	go func() { p.result = wait(); _ = input.Close(); _ = stdout.Close(); _ = stderr.Close(); close(p.done) }()
	return p, nil
}

// Wait is bounded by the caller's context; cancellation terminates only this tree.
func (p *Process) Wait(ctx context.Context) error {
	select {
	case <-p.done:
		return p.result
	case <-ctx.Done():
		_ = p.Stop()
		return ctx.Err()
	}
}

// Stop terminates only the process tree owned by this handle, once.
func (p *Process) Stop() error {
	p.once.Do(func() {
		p.stopErr = p.stop()
		select {
		case <-p.done:
		case <-time.After(5 * time.Second):
			p.stopErr = errors.Join(p.stopErr, errors.New("owned child did not exit within cleanup deadline"))
		}
	})
	return p.stopErr
}
