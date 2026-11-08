//go:build windows

package integration

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/incident"
)

func TestPowerShellIncidentBeforeIdentityPrelude(t *testing.T) {
	var tw twin
	for _, candidate := range allTwins() {
		if candidate.name == "powershell" {
			tw = candidate
		}
	}
	requireTwin(t, tw)
	parent := os.Getenv("GOSCHEDULE_INCIDENT_DIR")
	if parent == "" {
		parent = canonicalScriptIncidentParent(t)
	}
	t.Setenv("GOSCHEDULE_INCIDENT_DIR", parent)
	// Process-local error mode is inherited before any PowerShell code executes.
	mode, _, _ := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetErrorMode").Call()
	old := windows.SetErrorMode(uint32(mode) | 3)
	defer windows.SetErrorMode(old)
	harness := filepath.Join(canonicalScriptIncidentParent(t), "early-crash.ps1")
	if err := os.WriteFile(harness, []byte("[Environment]::FailFast('owned S138 early-crash fixture')\n"), 0600); err != nil {
		t.Fatal(err)
	}
	got := observeScript(tw.command(harness), harness, os.Environ())
	if got.code == 0 || got.code == -1 {
		t.Fatalf("expected native crash, got %d: %s", got.code, got.output)
	}
	var record incident.Record
	if err := json.Unmarshal([]byte(got.incident), &record); err != nil {
		t.Fatal(err)
	}
	if record.Original.NativeExit != got.code || record.Target.PID <= 0 {
		t.Fatalf("original native result lost: %+v", record)
	}
	directory := filepath.Join(parent, "incident-"+record.ID)
	if _, err := os.Stat(filepath.Join(directory, "runtime.json")); !os.IsNotExist(err) {
		t.Fatalf("pre-entry fixture unexpectedly has prelude identity: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(directory, "runtime-events.json"))
	if err != nil {
		t.Fatalf("missing pre-entry runtime-event receipt: %v", err)
	}
	var events incident.RuntimeEventReceipt
	if err := json.Unmarshal(data, &events); err != nil {
		t.Fatal(err)
	}
	if events.PID != record.Target.PID || events.Application != "pwsh.exe" {
		t.Fatalf("event target mismatch: %s", data)
	}
	switch events.Status {
	case "collected":
		if len(events.Events) == 0 {
			t.Fatal("claimed events without evidence")
		}
	case "absent", "unavailable", "failed", "truncated":
		if events.Detail == "" {
			t.Fatal("unexplained collection outcome")
		}
	default:
		t.Fatalf("unknown event status: %s", data)
	}
	for _, event := range events.Events {
		if event.PID != record.Target.PID || !strings.EqualFold(event.Application, "pwsh.exe") {
			t.Fatalf("unrelated event retained: %+v", event)
		}
	}
	if record.Collection.Status != "collected" && record.Collection.Status != "absent" && record.Collection.Status != "oversized" {
		t.Fatalf("dump status changed: %+v", record)
	}
	t.Logf("controlled pre-identity native=0x%08x events=%s", uint32(got.code), data)
}

func TestPowerShellIncidentSuccessDoesNotQueryEvents(t *testing.T) {
	var tw twin
	for _, candidate := range allTwins() {
		if candidate.name == "powershell" {
			tw = candidate
		}
	}
	requireTwin(t, tw)
	parent := canonicalScriptIncidentParent(t)
	t.Setenv("GOSCHEDULE_INCIDENT_DIR", parent)
	harness := filepath.Join(canonicalScriptIncidentParent(t), "success.ps1")
	if err := os.WriteFile(harness, []byte("exit 0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	got := observeScript(tw.command(harness), harness, os.Environ())
	if got.code != 0 || got.err != nil {
		t.Fatalf("ordinary success changed: %+v", got)
	}
	var record incident.Record
	if err := json.Unmarshal([]byte(got.incident), &record); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range record.Artifacts {
		if artifact.Name == "runtime-events.json" {
			t.Fatal("successful interpreter queried runtime events")
		}
	}
}

func TestRuntimeEventQueryBoundsNative(t *testing.T) {
	directory := canonicalScriptIncidentParent(t)
	source := filepath.Join(directory, "query-helper.go")
	code := `package main
import("fmt";"os";"strings";"time")
func main(){switch os.Args[1]{case "overflow":fmt.Print("overflow-ready:"+strings.Repeat("x",327680));case "wait":fmt.Print("wait-ready");<-time.After(time.Hour);case "failure":os.Exit(37)}}
`
	if err := os.WriteFile(source, []byte(code), 0600); err != nil {
		t.Fatal(err)
	}
	helper := filepath.Join(directory, "query-helper.exe")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	goPath, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	goPath, err = filepath.Abs(goPath)
	if err != nil {
		t.Fatal(err)
	}
	workdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	stdout, stderr := filepath.Join(directory, "build.stdout"), filepath.Join(directory, "build.stderr")
	build, err := clifixture.StartProcess(clifixture.Command{Executable: goPath, Arguments: []string{"build", "-race=false", "-o", helper, source}, Directory: workdir, StdoutPath: stdout, StderrPath: stderr, Environment: append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GCM_INTERACTIVE=Never")})
	if err != nil {
		t.Fatal(err)
	}
	waitErr := build.Wait(ctx)
	stopErr := build.Stop()
	if waitErr != nil || stopErr != nil {
		read := func(path string) string {
			f, err := os.Open(path)
			if err != nil {
				return err.Error()
			}
			defer f.Close()
			b, _ := io.ReadAll(io.LimitReader(f, 64<<10))
			return string(b)
		}
		t.Fatalf("prepare bounded native helper: %v: %s %s", errors.Join(waitErr, stopErr), read(stdout), read(stderr))
	}

	for _, role := range []string{"overflow", "wait", "failure"} {
		t.Run(role, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			result := runRuntimeEventQuery(ctx, helper, []string{role}, directory)
			switch role {
			case "overflow":
				if !result.Truncated || len(result.Data) != incident.RuntimeEventByteLimit+1 || !strings.HasPrefix(string(result.Data), "overflow-ready:") {
					t.Fatalf("output limit did not reach its intended child: %+v", result)
				}
			case "wait":
				if result.Error != context.DeadlineExceeded || !strings.Contains(string(result.Data), "wait-ready") {
					t.Fatalf("cancellation did not reach the waiting child: %+v", result)
				}
			case "failure":
				if !result.Unavailable || result.Error == nil || !strings.Contains(result.Error.Error(), "37") {
					t.Fatalf("query native failure lost: %+v", result)
				}
			}
		})
	}
}
