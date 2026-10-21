package clifixture

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// ProcessLimits applies optional hard limits to the entire owned environment tree.
// Windows CPUPercent uses aggregate system CPU, not a single-core denominator.
type ProcessLimits struct {
	CPUPercent   int
	MemoryBytes  uint64
	MaxProcesses uint32
}

// Pipes gives the caller bounded-capture and interactive protocol capabilities.
// Stop the owned tree before closing output readers, and always close this owner.
type Pipes struct {
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
	once   sync.Once
	err    error
}

// Close releases all parent endpoints; repeated closes are harmless.
func (p *Pipes) Close() error {
	p.once.Do(func() { p.err = errors.Join(closePipe(p.Stdin), closePipe(p.Stdout), closePipe(p.Stderr)) })
	return p.err
}
func closePipe(file *os.File) error {
	if file == nil {
		return nil
	}
	err := file.Close()
	if errors.Is(err, os.ErrClosed) {
		return nil
	}
	return err
}

// StartPiped launches an exact contained child with no raw output spool files.
func StartPiped(command Command) (*Process, *Pipes, error) {
	if !filepath.IsAbs(command.Executable) || !filepath.IsAbs(command.Directory) {
		return nil, nil, errors.New("owned piped child requires absolute executable and directory")
	}
	if err := validateProcessLimits(command.Limits); err != nil {
		return nil, nil, err
	}
	input, writer, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	output, childOutput, err := os.Pipe()
	if err != nil {
		return nil, nil, errors.Join(err, input.Close(), writer.Close())
	}
	diagnostic, childDiagnostic, err := os.Pipe()
	if err != nil {
		return nil, nil, errors.Join(err, input.Close(), writer.Close(), output.Close(), childOutput.Close())
	}
	pipes := &Pipes{Stdin: writer, Stdout: output, Stderr: diagnostic}
	p, wait, startErr := startContained(command, input, childOutput, childDiagnostic)
	closeErr := errors.Join(input.Close(), childOutput.Close(), childDiagnostic.Close())
	if startErr != nil {
		return nil, nil, errors.Join(startErr, closeErr, pipes.Close())
	}
	p.done = make(chan struct{})
	go func() { p.result = wait(); close(p.done) }()
	if closeErr != nil {
		return nil, nil, errors.Join(closeErr, p.Stop(), pipes.Close())
	}
	return p, pipes, nil
}

// ExitCode reports an actual completed child exit, independently of Wait cancellation.
func (p *Process) ExitCode() (int, bool) {
	select {
	case <-p.done:
		if p.result == nil {
			return 0, true
		}
		var exited interface{ ExitCode() int }
		if errors.As(p.result, &exited) {
			return exited.ExitCode(), true
		}
		return -1, true
	default:
		return 0, false
	}
}

func validateProcessLimits(limits *ProcessLimits) error {
	if limits != nil && (limits.CPUPercent < 1 || limits.CPUPercent > 100 || limits.MemoryBytes < 32<<20 || limits.MaxProcesses < 1) {
		return errors.New("invalid owned process hard limits")
	}
	return nil
}
