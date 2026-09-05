// Package executor runs a task's command. On Windows the child process is
// spawned with no console window (and the GUI binary is built windowless), so
// running tasks never flash or leave a visible command prompt. Output is
// captured up to a configured cap and a Run record is produced.
package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/platform"
)

// Executor runs commands and produces Run records.
type Executor struct {
	capBytes int
}

// New returns an Executor capping captured output at capBytes (per run).
func New(capBytes int) *Executor {
	if capBytes <= 0 {
		capBytes = 1 << 20
	}
	return &Executor{capBytes: capBytes}
}

// Run executes the task's command and returns a populated Run. The trigger is
// recorded so history distinguishes scheduled/manual/catch-up/event runs.
func (e *Executor) Run(ctx context.Context, task domain.Task, scheduledFor time.Time, trigger domain.RunTrigger) domain.Run {
	start := time.Now().UTC()
	run := domain.Run{TaskID: task.ID, ScheduledFor: scheduledFor, StartedAt: &start, Trigger: trigger}

	cmd := exec.CommandContext(ctx, task.Command, task.Args...)
	if task.WorkingDir != "" {
		cmd.Dir = task.WorkingDir
	}
	if len(task.Env) > 0 {
		cmd.Env = append(os.Environ(), envSlice(task.Env)...)
	}
	if task.Stdin != "" {
		cmd.Stdin = strings.NewReader(task.Stdin)
	}
	platform.HideConsole(cmd) // no console window for the child process
	buf := &capBuffer{cap: e.capBytes}

	_, explicitHome := task.Env["HOME"]
	if err := applyRunAs(cmd, task.RunAs, explicitHome); err != nil {
		end := time.Now().UTC()
		run.EndedAt = &end
		run.Outcome = domain.OutcomeFailure
		captureDiagnostic(buf, &run, "run_as: "+err.Error())
		return run
	}

	cmd.Stdout = buf
	cmd.Stderr = buf

	err := cmd.Run()
	end := time.Now().UTC()
	run.EndedAt = &end
	run.Output = buf.String()
	run.OutputTruncated = buf.Truncated()

	if err != nil {
		run.Outcome = domain.OutcomeFailure
		if ee, ok := err.(*exec.ExitError); ok {
			code := ee.ExitCode()
			run.ExitCode = &code
		} else if run.Output == "" {
			captureDiagnostic(buf, &run, fmt.Sprintf("process start failed for %q: %v", task.Command, err))
		}
		return run
	}
	zero := 0
	run.Outcome = domain.OutcomeSuccess
	run.ExitCode = &zero
	return run
}

func captureDiagnostic(buf *capBuffer, run *domain.Run, message string) {
	_, _ = buf.Write([]byte(message))
	run.Output = buf.String()
	run.OutputTruncated = buf.Truncated()
}

func envSlice(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}

// capBuffer is a thread-safe buffer that retains at most cap bytes. os/exec may
// write stdout and stderr from separate goroutines, so writes are guarded.
type capBuffer struct {
	cap       int
	mu        sync.Mutex
	buf       bytes.Buffer
	truncated bool
}

func (b *capBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	remaining := b.cap - b.buf.Len()
	if remaining <= 0 {
		b.truncated = true
		return len(p), nil // silently drop beyond the cap
	}
	if len(p) > remaining {
		b.buf.Write(p[:remaining])
		b.truncated = true
		return len(p), nil
	}
	return b.buf.Write(p)
}

// Truncated reports whether Write discarded any bytes after reaching the cap.
func (b *capBuffer) Truncated() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.truncated
}

func (b *capBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
