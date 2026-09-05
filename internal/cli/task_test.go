package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestGroupIntent_ThreeWay pins the CLI half of FR-014/FR-015. Before this
// change `--group ""` was indistinguishable from omitting the flag, so no
// client could take a task back out of a group.
func TestGroupIntent_ThreeWay(t *testing.T) {
	tests := []struct {
		name string
		argv []string
		want *string // nil = leave membership unchanged
	}{
		{"omitted leaves membership unchanged", []string{"some-id"}, nil},
		{"explicit empty removes from group", []string{"some-id", "--group", ""}, strp("")},
		{"named group assigns", []string{"some-id", "--group", "g-123"}, strp("g-123")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := taskEdit()
			// Parse flags only; RunE would need a live daemon.
			if err := cmd.ParseFlags(tt.argv); err != nil {
				t.Fatalf("parse %v: %v", tt.argv, err)
			}
			group, err := cmd.Flags().GetString("group")
			if err != nil {
				t.Fatal(err)
			}

			got := groupIntent(cmd, group)
			switch {
			case tt.want == nil && got != nil:
				t.Errorf("got %q, want nil (membership must be left unchanged)", *got)
			case tt.want != nil && got == nil:
				t.Errorf("got nil, want %q", *tt.want)
			case tt.want != nil && got != nil && *got != *tt.want:
				t.Errorf("got %q, want %q", *got, *tt.want)
			}
		})
	}
}

func strp(s string) *string { return &s }

func TestTaskScheduleHelpDescribesBothAcceptedSyntaxes(t *testing.T) {
	for _, cmd := range []*cobra.Command{taskAdd(), taskEdit()} {
		usage := cmd.Flags().Lookup("schedule").Usage
		lower := strings.ToLower(usage)
		if !strings.Contains(lower, "human") || !strings.Contains(lower, "cron") {
			t.Errorf("%s --schedule usage = %q, want human and cron", cmd.Name(), usage)
		}
	}
}

func TestTaskAddAcceptsUnnamedDraftAndEditExposesClearIntent(t *testing.T) {
	add := taskAdd()
	if err := add.Args(add, nil); err != nil {
		t.Fatalf("task add rejected an unnamed draft: %v", err)
	}
	edit := taskEdit()
	for _, name := range []string{"name", "command", "clear-schedule"} {
		if edit.Flags().Lookup(name) == nil {
			t.Errorf("task edit missing --%s", name)
		}
	}
}

func TestRootHelpUsesHumanFirstDualSyntaxPositioning(t *testing.T) {
	help := strings.ToLower(newRoot().Short)
	if !strings.Contains(help, "readable") || !strings.Contains(help, "cron") {
		t.Fatalf("root help = %q, want readable schedules and cron", newRoot().Short)
	}
}

func TestTaskDSTPolicyFlagsAreConsistent(t *testing.T) {
	for _, cmd := range []*cobra.Command{taskAdd(), taskEdit()} {
		for _, name := range []string{"time-basis", "dst-gap", "dst-overlap"} {
			flag := cmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("%s missing --%s", cmd.Name(), name)
				continue
			}
			if !strings.Contains(flag.Usage, "|") {
				t.Errorf("%s --%s usage does not name choices: %q", cmd.Name(), name, flag.Usage)
			}
		}
	}
}
