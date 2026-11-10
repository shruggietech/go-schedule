package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/platform"
	"github.com/shruggietech/go-schedule/v2/internal/runcontrol"
)

func cancellationTask(mode, directory string) domain.Task {
	return domain.Task{ID: "same-task", Command: os.Args[0], Args: []string{"-test.run=^TestExecutorCancellationHelper$"}, Env: map[string]string{"GOSCHED_EXECUTOR_CANCEL_MODE": mode, "GOSCHED_EXECUTOR_CANCEL_DIR": directory, "GORACE": "atexit_sleep_ms=0"}}
}

func waitCancellationFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		select {
		case <-deadline.C:
			t.Fatalf("owned helper did not become ready: %s", filepath.Base(path))
		case <-ticker.C:
		}
	}
}

func TestExecutorExplicitCancellationBeforeLaunch(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(runcontrol.ErrRequested)
	run := New(0).Run(ctx, echoTask("must-not-execute"), time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeCanceled || run.ExitCode != nil || run.StartFailed || strings.Contains(run.Output, "must-not-execute") {
		t.Fatalf("prelaunch cancellation=%+v", run)
	}
}

func TestExecutorCancellationStopsOwnedDescendant(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(runcontrol.ErrRequested)
	result := make(chan domain.Run, 1)
	go func() {
		result <- New(0).Run(ctx, cancellationTask("parent", dir), time.Now().UTC(), domain.TriggerManual)
	}()
	waitCancellationFile(t, filepath.Join(dir, "ready"))
	cancel(runcontrol.ErrRequested)
	select {
	case run := <-result:
		if run.Outcome != domain.OutcomeCanceled || run.ExitCode == nil || !strings.Contains(run.Output, "partial-output") {
			t.Fatalf("canceled run=%+v", run)
		}
		if _, err := os.Stat(filepath.Join(dir, "finished")); err == nil {
			t.Fatal("descendant completed naturally instead of being stopped")
		}
	case <-time.After(7 * time.Second):
		t.Fatal("owned cancellation did not finish within its cleanup allowance")
	}
}

func TestExecutorNaturalCompletionPreservesBackgroundChild(t *testing.T) {
	dir := t.TempDir()
	run := New(0).Run(context.Background(), cancellationTask("background", dir), time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess {
		t.Fatalf("natural parent=%+v", run)
	}
	waitCancellationFile(t, filepath.Join(dir, "ready"))
	if err := os.WriteFile(filepath.Join(dir, "release"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	waitCancellationFile(t, filepath.Join(dir, "finished"))
}

func TestExecutorCancellationHelper(t *testing.T) {
	mode := os.Getenv("GOSCHED_EXECUTOR_CANCEL_MODE")
	if mode == "" {
		return
	}
	dir := os.Getenv("GOSCHED_EXECUTOR_CANCEL_DIR")
	if mode == "parent" || mode == "background" {
		child := exec.Command(os.Args[0], "-test.run=^TestExecutorCancellationHelper$")
		child.Env = append(os.Environ(), "GOSCHED_EXECUTOR_CANCEL_MODE=child")
		if mode == "parent" {
			child.Stdout = os.Stdout
			child.Stderr = os.Stderr
		}
		platform.HideConsole(child)
		if err := child.Start(); err != nil {
			t.Fatal(err)
		}
		if mode == "background" {
			os.Exit(0)
		}
		if err := child.Wait(); err != nil {
			t.Fatal(err)
		}
		os.Exit(0)
	}
	if err := os.WriteFile(filepath.Join(dir, "ready"), []byte(fmt.Sprint(os.Getpid())), 0600); err != nil {
		t.Fatal(err)
	}
	fmt.Fprintln(os.Stdout, "partial-output")
	deadline := time.NewTimer(3 * time.Second)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(filepath.Join(dir, "release")); err == nil {
			break
		}
		select {
		case <-deadline.C:
			if err := os.WriteFile(filepath.Join(dir, "finished"), nil, 0600); err != nil {
				t.Fatal(err)
			}
			os.Exit(0)
		case <-ticker.C:
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "finished"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func TestExecutorCancellationIsolatesConcurrentSameTaskRuns(t *testing.T) {
	dirs := []string{t.TempDir(), t.TempDir()}
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(runcontrol.ErrRequested)
	first, second := make(chan domain.Run, 1), make(chan domain.Run, 1)
	go func() {
		first <- New(0).Run(ctx, cancellationTask("child", dirs[0]), time.Now().UTC(), domain.TriggerManual)
	}()
	go func() {
		second <- New(0).Run(context.Background(), cancellationTask("child", dirs[1]), time.Now().UTC(), domain.TriggerManual)
	}()
	for _, dir := range dirs {
		waitCancellationFile(t, filepath.Join(dir, "ready"))
	}
	cancel(runcontrol.ErrRequested)
	select {
	case run := <-first:
		if run.Outcome != domain.OutcomeCanceled {
			t.Fatalf("selected run=%+v", run)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("selected run did not stop")
	}
	select {
	case run := <-second:
		t.Fatalf("other same-task run stopped: %+v", run)
	default:
	}
	if err := os.WriteFile(filepath.Join(dirs[1], "release"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	select {
	case run := <-second:
		if run.Outcome != domain.OutcomeSuccess {
			t.Fatalf("other run=%+v", run)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("other run did not complete")
	}
}

type controlledOwned struct {
	terminateFn func(time.Time) error
	emptyFn     func() (bool, error)
}

func (p controlledOwned) wait() nativeResult          { return nativeResult{} }
func (p controlledOwned) terminate(d time.Time) error { return p.terminateFn(d) }
func (p controlledOwned) empty() (bool, error)        { return p.emptyFn() }
func (p controlledOwned) release() error              { return nil }

func TestCancellationDeadlineReportsWithoutReleasingLiveOwnership(t *testing.T) {
	failures := make(chan error, 4)
	ctx, cancel := context.WithCancelCause(runcontrol.WithFailureReporter(context.Background(), func(err error) { failures <- err }))
	cancel(runcontrol.ErrRequested)
	native := make(chan nativeResult, 1)
	copies := make(chan error, 1)
	closed := make(chan struct{})
	var empty atomic.Bool
	p := controlledOwned{terminateFn: func(deadline time.Time) error {
		timer := time.NewTimer(max(0, time.Until(deadline)))
		defer timer.Stop()
		<-timer.C
		return errors.New("controlled termination remains unconfirmed")
	}, emptyFn: func() (bool, error) { return empty.Load(), nil }}
	done := make(chan executionResult, 1)
	go func() {
		done <- awaitOwned(ctx, p, native, copies, 1, func() error { close(closed); copies <- nil; return nil }, 30*time.Millisecond)
	}()
	select {
	case err := <-failures:
		if err == nil {
			t.Fatal("missing failure")
		}
	case <-time.After(time.Second):
		t.Fatal("cleanup failure feedback not bounded")
	}
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("owned collection pipes were not closed")
	}
	select {
	case <-done:
		t.Fatal("live native work was released at deadline")
	default:
	}
	empty.Store(true)
	code := -1
	native <- nativeResult{code: &code, err: errors.New("controlled native exit")}
	select {
	case r := <-done:
		if !r.canceled || !r.cleanupFailed || r.code == nil || *r.code != -1 {
			t.Fatalf("result=%+v", r)
		}
	case <-time.After(time.Second):
		t.Fatal("confirmed ownership did not retire")
	}
}

func TestObservedNaturalCompletionWinsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(runcontrol.ErrRequested)
	native := make(chan nativeResult, 1)
	code := 7
	native <- nativeResult{code: &code, err: errors.New("exit seven")}
	copies := make(chan error, 1)
	copies <- nil
	var terminated atomic.Bool
	p := controlledOwned{terminateFn: func(time.Time) error { terminated.Store(true); return nil }, emptyFn: func() (bool, error) { return true, nil }}
	r := awaitOwned(ctx, p, native, copies, 1, func() error { return nil }, time.Second)
	if r.canceled || r.cleanupFailed || r.code == nil || *r.code != 7 || terminated.Load() {
		t.Fatalf("natural result rewritten: %+v terminated=%v", r, terminated.Load())
	}
}

func TestNativeCreationDeadlineReportsAndRetainsPendingCreation(t *testing.T) {
	failures := make(chan error, 1)
	ctx, cancel := context.WithCancelCause(runcontrol.WithFailureReporter(context.Background(), func(err error) { failures <- err }))
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan startupObservation, 1)
	go func() {
		done <- observeStart(ctx, func() (ownedProcess, error) { close(entered); <-release; return nil, context.Canceled }, 30*time.Millisecond)
	}()
	<-entered
	cancel(runcontrol.ErrRequested)
	select {
	case <-failures:
	case <-time.After(time.Second):
		t.Fatal("slow native creation did not report cancellation deadline")
	}
	select {
	case <-done:
		t.Fatal("pending native creation was forgotten at deadline")
	default:
	}
	close(release)
	select {
	case r := <-done:
		if r.errorFeedback == nil || r.deadline.IsZero() || !errors.Is(r.err, context.Canceled) {
			t.Fatalf("startup=%+v", r)
		}
	case <-time.After(time.Second):
		t.Fatal("completed native creation did not retire")
	}
}
