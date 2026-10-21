package clievaluate

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

type memoryCapture struct {
	mu       sync.Mutex
	data     bytes.Buffer
	limit    int64
	exceeded bool
	cancel   func()
	once     sync.Once
}

func newMemoryCapture(limit int64, cancel func()) *memoryCapture {
	return &memoryCapture{limit: limit, cancel: cancel}
}
func (c *memoryCapture) Write(data []byte) (int, error) {
	n := len(data)
	c.mu.Lock()
	left := c.limit - int64(c.data.Len())
	if int64(len(data)) > left {
		c.exceeded = true
		data = data[:max(0, int(left))]
	}
	_, err := c.data.Write(data)
	overflow := c.exceeded
	c.mu.Unlock()
	if overflow {
		c.once.Do(func() {
			if c.cancel != nil {
				c.cancel()
			}
		})
	}
	return n, err
}
func (c *memoryCapture) Snapshot() ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]byte{}, c.data.Bytes()...), c.exceeded
}

type commandCapture struct {
	Code     int
	Out, Err []byte
	Exceeded bool
	Cleanup  error
	Failure  error
}

func captureCommand(ctx context.Context, command clifixture.Command, input []byte, limit int64) commandCapture {
	return captureManaged(ctx, command, input, limit, false)
}

func captureManaged(ctx context.Context, command clifixture.Command, input []byte, limit int64, keepAlive bool) commandCapture {
	return captureWithStart(ctx, command, input, limit, keepAlive, nil)
}

func captureWithStart(ctx context.Context, command clifixture.Command, input []byte, limit int64, keepAlive bool, onStart func(int) error) commandCapture {
	if err := ctx.Err(); err != nil {
		return commandCapture{Code: -1, Failure: err}
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	out, stderr := newMemoryCapture(limit, cancel), newMemoryCapture(limit, cancel)
	p, pipes, err := clifixture.StartPiped(command)
	if err != nil {
		return commandCapture{Code: -1, Failure: err}
	}
	if onStart != nil {
		if err := onStart(p.PID); err != nil {
			return commandCapture{Code: -1, Failure: errors.Join(err, p.Stop(), pipes.Close())}
		}
	}
	finished := make(chan error, 3)
	go func() { _, err := io.Copy(out, pipes.Stdout); finished <- err }()
	go func() { _, err := io.Copy(stderr, pipes.Stderr); finished <- err }()
	inputDone := make(chan struct{})
	go func() {
		_, err := io.Copy(pipes.Stdin, bytes.NewReader(input))
		if keepAlive && err == nil {
			<-inputDone
		}
		finished <- errors.Join(err, pipes.Stdin.Close())
	}()
	var waitErr error
	if keepAlive {
		waited := make(chan error, 1)
		go func() { waited <- p.Wait(context.Background()) }()
		select {
		case waitErr = <-waited:
			close(inputDone)
		case <-ctx.Done():
			// EOF cancels the owned worker and gives it a bounded chance to persist
			// its ledger and remove runtime state before exact-tree termination.
			close(inputDone)
			grace := time.NewTimer(7 * time.Second)
			select {
			case err := <-waited:
				waitErr = errors.Join(ctx.Err(), err)
			case <-grace.C:
				waitErr = errors.Join(ctx.Err(), p.Stop())
			}
			grace.Stop()
		}
	} else {
		waitErr = p.Wait(ctx)
		close(inputDone)
	}
	cleanup := p.Stop()
	var closeErr error
	if cleanup != nil {
		closeErr = pipes.Close()
	}
	var streamErr error
	for i := 0; i < 3; i++ {
		streamErr = errors.Join(streamErr, <-finished)
	}
	closeErr = errors.Join(closeErr, pipes.Close())
	code, done := p.ExitCode()
	if !done {
		code = -1
	}
	stdoutBytes, outExceeded := out.Snapshot()
	stderrBytes, errExceeded := stderr.Snapshot()
	if waitErr != nil && ctx.Err() != nil {
		streamErr = errors.Join(streamErr, ctx.Err())
	}
	return commandCapture{Code: code, Out: stdoutBytes, Err: stderrBytes, Exceeded: outExceeded || errExceeded, Cleanup: cleanup, Failure: errors.Join(streamErr, closeErr)}
}
