package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/runcontrol"
)

type nativeResult struct {
	code *int
	err  error
}
type ownedProcess interface {
	wait() nativeResult
	terminate(time.Time) error
	empty() (bool, error)
	release() error
}
type executionResult struct {
	nativeResult
	started, canceled, cleanupFailed bool
	diagnostic                       error
}

// processIO owns the parent's pipe endpoints. Only an explicit cancellation
// deadline closes them early; normal completion retains ordinary pipe draining.
type processIO struct {
	input, output           *os.File
	childInput, childOutput *os.File
	source                  io.Reader
	copies                  int
}

func openProcessIO(cmd *exec.Cmd) (*processIO, error) {
	s := &processIO{source: cmd.Stdin, copies: 1}
	var err error
	if s.source == nil {
		s.childInput, err = os.Open(os.DevNull)
	} else {
		s.childInput, s.input, err = os.Pipe()
		s.copies++
	}
	if err != nil {
		return nil, err
	}
	s.output, s.childOutput, err = os.Pipe()
	if err != nil {
		return nil, errors.Join(err, s.closeAll())
	}
	cmd.Stdin = s.childInput
	cmd.Stdout = s.childOutput
	cmd.Stderr = s.childOutput
	return s, nil
}
func closeOwnedFile(f *os.File) error {
	if f == nil {
		return nil
	}
	err := f.Close()
	if errors.Is(err, os.ErrClosed) {
		return nil
	}
	return err
}
func (s *processIO) closeChild() error {
	return errors.Join(closeOwnedFile(s.childInput), closeOwnedFile(s.childOutput))
}
func (s *processIO) closeParent() error {
	return errors.Join(closeOwnedFile(s.input), closeOwnedFile(s.output))
}
func (s *processIO) closeAll() error { return errors.Join(s.closeParent(), s.closeChild()) }
func (s *processIO) collect(buf io.Writer) <-chan error {
	done := make(chan error, s.copies)
	go func() { _, err := io.Copy(buf, s.output); done <- err }()
	if s.source != nil {
		go func() {
			_, err := io.Copy(s.input, s.source)
			if isBrokenInput(err) {
				err = nil
			}
			done <- errors.Join(err, closeOwnedFile(s.input))
		}()
	}
	return done
}

type startupObservation struct {
	process            ownedProcess
	err, errorFeedback error
	deadline           time.Time
}

// observeStart keeps ownership of a slow native creation call while reporting
// cancellation failure at its original deadline. It cannot fabricate an exit
// merely because the OS has not returned a process handle yet.
func observeStart(ctx context.Context, start func() (ownedProcess, error), allowance time.Duration) startupObservation {
	created := make(chan startupObservation, 1)
	go func() { p, err := start(); created <- startupObservation{process: p, err: err} }()
	select {
	case result := <-created:
		return result
	case <-ctx.Done():
	}
	end := time.Now().Add(allowance)
	timer := time.NewTimer(allowance)
	defer timer.Stop()
	select {
	case result := <-created:
		result.deadline = end
		return result
	case <-timer.C:
		err := errors.New("native process creation remains unconfirmed at cancellation deadline")
		reportCancellationFailure(ctx, err)
		result := <-created
		result.deadline = end
		result.errorFeedback = err
		return result
	}
}

func executeOwned(ctx context.Context, cmd *exec.Cmd, buf io.Writer) executionResult {
	streams, err := openProcessIO(cmd)
	if err != nil {
		return executionResult{nativeResult: nativeResult{err: err}}
	}
	startup := observeStart(ctx, func() (ownedProcess, error) { return startOwnedProcess(ctx, cmd) }, runcontrol.CleanupTimeout)
	if startup.err != nil && startup.process == nil {
		return executionResult{nativeResult: nativeResult{err: errors.Join(startup.err, streams.closeAll())}, canceled: explicitStartupCancellation(ctx, startup.err), cleanupFailed: startup.errorFeedback != nil, diagnostic: startup.errorFeedback}
	}
	childClose := streams.closeChild()
	copies := streams.collect(buf)
	native := make(chan nativeResult, 1)
	go func() { native <- startup.process.wait() }()
	waitCtx := ctx
	if startup.err != nil && ctx.Err() == nil {
		c, cancel := context.WithCancelCause(ctx)
		cancel(startup.err)
		waitCtx = c
	}
	allowance := runcontrol.CleanupTimeout
	if !startup.deadline.IsZero() {
		allowance = max(0, time.Until(startup.deadline))
	}
	result := awaitOwned(waitCtx, startup.process, native, copies, streams.copies, streams.closeParent, allowance)
	result.started = startup.err == nil
	if startup.err != nil {
		result.err = errors.Join(startup.err, result.err)
		result.canceled = explicitStartupCancellation(ctx, startup.err)
	}
	if startup.errorFeedback != nil {
		result.cleanupFailed = true
		result.canceled = errors.Is(context.Cause(ctx), runcontrol.ErrRequested)
	}
	result.diagnostic = errors.Join(result.diagnostic, startup.errorFeedback, childClose, streams.closeAll(), startup.process.release())
	return result
}
func explicitStartupCancellation(ctx context.Context, err error) bool {
	return ctx.Err() != nil && errors.Is(err, ctx.Err()) && errors.Is(context.Cause(ctx), runcontrol.ErrRequested)
}
func reportCancellationFailure(ctx context.Context, err error) {
	if errors.Is(context.Cause(ctx), runcontrol.ErrRequested) {
		runcontrol.ReportFailure(ctx, err)
	}
}

type ownershipWait struct {
	ctx                                     context.Context
	process                                 ownedProcess
	native                                  <-chan nativeResult
	copies                                  <-chan error
	closePipes                              func() error
	allowance                               time.Duration
	result                                  executionResult
	remaining                               int
	rootDone, stopping, stopDone, treeEmpty bool
	stop                                    <-chan error
	deadline                                <-chan time.Time
	timer                                   *time.Timer
	ticker                                  *time.Ticker
	poll                                    <-chan time.Time
	ctxDone                                 <-chan struct{}
}

// awaitOwned retains the native waiter and exact tree capability until observed
// completion. The deadline bounds feedback and collection, not live work's
// lifetime: an unreaped process continues owning its worker capacity.
func awaitOwned(ctx context.Context, process ownedProcess, native <-chan nativeResult, copies <-chan error, count int, closePipes func() error, allowance time.Duration) executionResult {
	w := ownershipWait{ctx: ctx, process: process, native: native, copies: copies, remaining: count, closePipes: closePipes, allowance: allowance, ctxDone: ctx.Done()}
	defer w.closeTimers()
	for !w.finished() {
		select {
		case result := <-w.native:
			w.observeNative(result)
		case err := <-w.copies:
			w.observeCopy(err)
		case <-w.ctxDone:
			w.beginCancellation()
		case err := <-w.stop:
			w.finishStop(err)
		case <-w.deadline:
			w.expire()
		case <-w.poll:
			w.observeTree()
		}
	}
	return w.result
}
func (w *ownershipWait) finished() bool {
	return w.rootDone && w.remaining == 0 && (!w.stopping || (w.stopDone && w.treeEmpty))
}
func (w *ownershipWait) observeNative(result nativeResult) {
	w.result.nativeResult = result
	w.rootDone = true
	w.native = nil
}
func (w *ownershipWait) observeCopy(err error) {
	w.remaining--
	if err != nil && !errors.Is(err, os.ErrClosed) {
		w.result.diagnostic = errors.Join(w.result.diagnostic, fmt.Errorf("process output/input collection: %w", err))
	}
	if w.remaining == 0 {
		w.copies = nil
	}
}
func (w *ownershipWait) observeReady() {
	if !w.rootDone {
		select {
		case result := <-w.native:
			w.observeNative(result)
		default:
		}
	}
	for w.remaining > 0 {
		select {
		case err := <-w.copies:
			w.observeCopy(err)
		default:
			return
		}
	}
}
func (w *ownershipWait) beginCancellation() {
	w.observeReady()
	if w.finished() {
		return
	}
	w.stopping = true
	w.ctxDone = nil
	w.result.canceled = errors.Is(context.Cause(w.ctx), runcontrol.ErrRequested)
	end := time.Now().Add(w.allowance)
	w.timer = time.NewTimer(w.allowance)
	w.deadline = w.timer.C
	stopped := make(chan error, 1)
	w.stop = stopped
	go func() { stopped <- w.process.terminate(end) }()
	w.ticker = time.NewTicker(10 * time.Millisecond)
	w.poll = w.ticker.C
}
func (w *ownershipWait) report(err error) {
	if err == nil {
		return
	}
	w.result.cleanupFailed = true
	w.result.diagnostic = errors.Join(w.result.diagnostic, err)
	reportCancellationFailure(w.ctx, err)
}
func (w *ownershipWait) finishStop(err error) {
	w.stop = nil
	w.stopDone = true
	w.report(err)
	var queryErr error
	w.treeEmpty, queryErr = w.process.empty()
	w.report(queryErr)
}
func (w *ownershipWait) expire() {
	w.deadline = nil
	w.report(errors.New("owned execution termination or collection remains unconfirmed at cancellation deadline"))
	w.report(w.closePipes())
}
func (w *ownershipWait) observeTree() {
	if w.stopDone && !w.treeEmpty {
		empty, err := w.process.empty()
		if err == nil {
			w.treeEmpty = empty
		}
	}
}
func (w *ownershipWait) closeTimers() {
	if w.timer != nil {
		w.timer.Stop()
	}
	if w.ticker != nil {
		w.ticker.Stop()
	}
}
