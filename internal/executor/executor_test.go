package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"reflect"
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

func TestExecutor_ProcessStartFailureNamesBoundaryWithoutSecrets(t *testing.T) {
	task := domain.Task{
		Command: "definitely-not-a-real-go-schedule-executable-xyz",
		Args:    []string{"--secret-argument"},
		Env:     map[string]string{"GO_SCHEDULE_SECRET": "secret-environment-value"},
		Stdin:   "secret-stdin-value",
	}
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeFailure || run.ExitCode != nil {
		t.Fatalf("run = %#v, want process-start failure without exit code", run)
	}
	if !strings.Contains(run.Output, `process start failed for "definitely-not-a-real-go-schedule-executable-xyz"`) {
		t.Fatalf("output = %q, want named process-start boundary", run.Output)
	}
	for _, secret := range []string{"--secret-argument", "secret-environment-value", "secret-stdin-value"} {
		if strings.Contains(run.Output, secret) {
			t.Fatalf("output disclosed %q: %q", secret, run.Output)
		}
	}
}

func TestExecutor_NonzeroExitPreservesChildOutput(t *testing.T) {
	task := failTask()
	if runtime.GOOS == "windows" {
		task.Args = []string{"/d", "/q", "/c", "echo controlled failure & exit /b 7"}
	} else {
		task.Args = []string{"-c", "printf 'controlled failure'; exit 7"}
	}
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeFailure || run.ExitCode == nil || *run.ExitCode != 7 {
		t.Fatalf("run = %#v, want child exit 7", run)
	}
	if !strings.Contains(run.Output, "controlled failure") || strings.Contains(run.Output, "process start failed") {
		t.Fatalf("output = %q, want unchanged child output", run.Output)
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
	if !run.OutputTruncated {
		t.Fatal("discarded output was not reported as truncated")
	}
}

func TestCapBufferReportsOnlyDiscardedBytesAsTruncated(t *testing.T) {
	exact := &capBuffer{cap: 4}
	if n, err := exact.Write([]byte("four")); err != nil || n != 4 {
		t.Fatalf("exact write n=%d err=%v", n, err)
	}
	if exact.Truncated() {
		t.Fatal("exactly capped output reported truncation")
	}
	if n, err := exact.Write([]byte("more")); err != nil || n != 4 {
		t.Fatalf("discarded write n=%d err=%v", n, err)
	}
	if !exact.Truncated() || exact.String() != "four" {
		t.Fatalf("buffer=%q truncated=%v", exact.String(), exact.Truncated())
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

func TestExecutor_DirectArgumentsRemainLiteral(t *testing.T) {
	if os.Getenv("GO_SCHEDULE_EXECUTOR_ARG_HELPER") == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				_ = json.NewEncoder(os.Stdout).Encode(os.Args[i+1:])
				return
			}
		}
		t.Fatal("helper separator missing")
	}

	want := []string{"", "hello world", "héllo 世界", "tabs\there", "lines\r\nhere", `C:\trail\`, "$HOME", "%PATH%", "|", ">", "*.txt", ";", "&"}
	task := domain.Task{
		Command: os.Args[0],
		Args:    append([]string{"-test.run=^TestExecutor_DirectArgumentsRemainLiteral$", "--"}, want...),
		Env:     map[string]string{"GO_SCHEDULE_EXECUTOR_ARG_HELPER": "1"},
	}
	run := New(0).Run(context.Background(), task, time.Now().UTC(), domain.TriggerManual)
	if run.Outcome != domain.OutcomeSuccess {
		t.Fatalf("helper run = %#v", run)
	}
	var got []string
	if err := json.NewDecoder(bytes.NewBufferString(run.Output)).Decode(&got); err != nil {
		t.Fatalf("decode helper output %q: %v", run.Output, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("received args = %#v, want %#v", got, want)
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
