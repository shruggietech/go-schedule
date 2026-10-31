package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/clievaluate"
)

func TestPublicCommandsRejectUnexpectedOperandsBeforeOpeningState(t *testing.T) {
	for _, command := range []string{"inventory", "build", "capture", "resume", "finalize", "validate", "render"} {
		t.Run(command, func(t *testing.T) {
			args := []string{command}
			if command != "inventory" {
				if command == "build" {
					args = append(args, "--bin-dir", "unused")
				} else {
					args = append(args, "--out", "unused")
				}
			}
			if command == "capture" || command == "resume" {
				args = append(args, "--bin-dir", "unused", "--wsl-distro", "unused", "--wsl-user", "unused")
			}
			args = append(args, "unexpected")
			var out, diagnostic bytes.Buffer
			if err := run(context.Background(), args, &out, &diagnostic); err == nil {
				t.Fatal("unexpected operand accepted")
			}
			if out.Len() != 0 {
				t.Fatal("rejected command emitted a success result")
			}
		})
	}
}

func TestEvaluationAndProcessingExitCodesAreDistinct(t *testing.T) {
	for _, test := range []struct {
		err  error
		want int
	}{{nil, 0}, {flag.ErrHelp, 0}, {errors.New("invalid evidence"), 1}, {outcomeError{outcome: "fail"}, 2}, {outcomeError{outcome: "blocked"}, 3}, {fmt.Errorf("wrapped: %w", outcomeError{outcome: "fail"}), 2}} {
		if got := exitStatus(test.err); got != test.want {
			t.Errorf("%v: exit=%d want=%d", test.err, got, test.want)
		}
	}
}

func TestPublicHelpAndInventoryNeedNoExternalState(t *testing.T) {
	for _, command := range []string{"inventory", "build", "capture", "resume", "finalize", "validate", "render"} {
		var out, diagnostic bytes.Buffer
		if err := run(context.Background(), []string{command, "--help"}, &out, &diagnostic); !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("%s help: %v", command, err)
		}
		if diagnostic.Len() == 0 {
			t.Fatalf("%s missing help", command)
		}
	}
	var out, diagnostic bytes.Buffer
	if err := run(context.Background(), []string{"inventory"}, &out, &diagnostic); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out.Bytes(), []byte(`"go-schedule.cli-qa-manifest"`)) {
		t.Fatal("missing initialized manifest")
	}
}

func TestTimingPolicyUsesTheSameValidatedSelectedBudget(t *testing.T) {
	budget := clievaluate.DefaultBudget()
	budget.TotalTimeoutMS = 7200000
	file := filepath.Join(t.TempDir(), "budget.json")
	data, err := json.Marshal(budget)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	if err := run(context.Background(), []string{"timing-policy", "--budget", file}, &out, &diagnostic); err != nil {
		t.Fatal(err)
	}
	var policy clievaluate.TimingPolicy
	if err := json.Unmarshal(out.Bytes(), &policy); err != nil {
		t.Fatal(err)
	}
	if policy.TotalTimeoutMS != 7200000 || policy.OuterTimeoutMS <= 7200000 {
		t.Fatalf("wrong selected policy: %+v", policy)
	}
	out.Reset()
	if err := run(context.Background(), []string{"timing-policy", "unexpected"}, &out, &diagnostic); err == nil || out.Len() != 0 {
		t.Fatal("unexpected policy operand accepted")
	}
	budget.TotalTimeoutMS = 7200001
	data, err = json.Marshal(budget)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := run(context.Background(), []string{"timing-policy", "--budget", file}, &out, &diagnostic); err == nil || out.Len() != 0 {
		t.Fatal("invalid policy accepted")
	}
}
