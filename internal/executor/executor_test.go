package executor

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// helper: a command that succeeds and prints, cross-platform.
func echoTask(msg string) domain.Task {
	if runtime.GOOS == "windows" {
		return domain.Task{ID: "t", Command: "cmd", Args: []string{"/c", "echo " + msg}}
	}
	return domain.Task{ID: "t", Command: "sh", Args: []string{"-c", "echo " + msg}}
}

func failTask() domain.Task {
	if runtime.GOOS == "windows" {
		return domain.Task{ID: "t", Command: "cmd", Args: []string{"/c", "exit 3"}}
	}
	return domain.Task{ID: "t", Command: "sh", Args: []string{"-c", "exit 3"}}
}

func TestExecutor_Success(t *testing.T) {
	e := New(0)
	run := e.Run(context.Background(), echoTask("hello"), time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess {
		t.Fatalf("outcome = %v, output=%q", run.Outcome, run.Output)
	}
	if run.ExitCode == nil || *run.ExitCode != 0 {
		t.Fatalf("exit code = %v", run.ExitCode)
	}
	if !strings.Contains(run.Output, "hello") {
		t.Fatalf("output = %q, want to contain hello", run.Output)
	}
	if run.StartedAt == nil || run.EndedAt == nil {
		t.Fatal("timestamps should be set")
	}
}

func TestExecutor_Failure(t *testing.T) {
	e := New(0)
	run := e.Run(context.Background(), failTask(), time.Now().UTC(), domain.TriggerSchedule)
	if run.Outcome != domain.OutcomeFailure {
		t.Fatalf("outcome = %v", run.Outcome)
	}
	if run.ExitCode == nil || *run.ExitCode != 3 {
		t.Fatalf("exit code = %v, want 3", run.ExitCode)
	}
}

func TestExecutor_OutputCap(t *testing.T) {
	e := New(10) // cap to 10 bytes
	var task domain.Task
	if runtime.GOOS == "windows" {
		task = domain.Task{Command: "cmd", Args: []string{"/c", "echo aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	} else {
		task = domain.Task{Command: "sh", Args: []string{"-c", "printf 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'"}}
	}
	run := e.Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if len(run.Output) > 10 {
		t.Fatalf("output not capped: %d bytes", len(run.Output))
	}
}

func TestExecutor_SuppliesExactStdin(t *testing.T) {
	task := domain.Task{Command: "sh", Args: []string{"-c", "cat"}}
	if runtime.GOOS == "windows" {
		task = domain.Task{Command: "powershell.exe", Args: []string{
			"-NoProfile", "-NonInteractive", "-Command",
			"[Console]::OpenStandardInput().CopyTo([Console]::OpenStandardOutput())",
		}}
	}
	task.Stdin = "first line\nsecond % line\n"
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess || run.Output != task.Stdin {
		t.Fatalf("run = %#v, want exact stdin %q", run, task.Stdin)
	}
}

func TestExecutor_TaskEnvironmentOverridesParent(t *testing.T) {
	t.Setenv("GO_SCHEDULE_ENV_VALUE", "parent")
	task := domain.Task{
		Command: "sh", Args: []string{"-c", `printf %s "$GO_SCHEDULE_ENV_VALUE"`},
		Env: map[string]string{"GO_SCHEDULE_ENV_VALUE": "task"},
	}
	if runtime.GOOS == "windows" {
		task.Command = "powershell.exe"
		task.Args = []string{
			"-NoProfile", "-NonInteractive", "-Command",
			"[Console]::Out.Write($env:GO_SCHEDULE_ENV_VALUE)",
		}
	}
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess || run.Output != "task" {
		t.Fatalf("run = %#v, want task environment override", run)
	}
}

func TestRunAs_EmptyIsNoOp(t *testing.T) {
	if err := ValidateRunAs(""); err != nil {
		t.Fatalf("empty run_as should be valid: %v", err)
	}
}

func TestRunAs_InvalidRejected(t *testing.T) {
	// On Windows any run_as is unsupported; on Unix an unknown user is rejected.
	err := ValidateRunAs("definitely-not-a-real-user-xyz")
	if err == nil {
		t.Fatal("expected run_as validation to reject an unsupported/unknown user")
	}
}
