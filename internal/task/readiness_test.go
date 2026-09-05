package task

import (
	"reflect"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestDisplayNameFallsBackWithoutChangingIdentity(t *testing.T) {
	if got := DisplayName(domain.Task{ID: "one"}); got != "unnamed" {
		t.Fatalf("DisplayName = %q, want unnamed", got)
	}
	if got := DisplayName(domain.Task{ID: "two", Name: "Named"}); got != "Named" {
		t.Fatalf("DisplayName = %q, want Named", got)
	}
}

func TestReadinessClassifiesCommandAndAutomaticSources(t *testing.T) {
	tests := []struct {
		name       string
		task       domain.Task
		completion bool
		want       Readiness
	}{
		{name: "missing command", task: domain.Task{ScheduleID: "s"}, want: Readiness{AutomaticSources: []AutomaticSource{SourceSchedule}, Status: StatusNotRunnable, Reason: "No command is configured."}},
		{name: "manual only", task: domain.Task{Command: "echo"}, want: Readiness{CommandReady: true, Status: StatusManualOnly, Reason: "No automatic activation source is configured."}},
		{name: "scheduled", task: domain.Task{Command: "echo", ScheduleID: "s", Enabled: true, State: domain.TaskActive}, want: Readiness{CommandReady: true, ActivationReady: true, AutomaticSources: []AutomaticSource{SourceSchedule}, Status: StatusReady}},
		{name: "completion", task: domain.Task{Command: "echo", Enabled: true, State: domain.TaskActive}, completion: true, want: Readiness{CommandReady: true, ActivationReady: true, AutomaticSources: []AutomaticSource{SourceCompletion}, Status: StatusReady}},
		{name: "terminal precedence", task: domain.Task{Command: "echo", ScheduleID: "s", State: domain.TaskCompleted}, want: Readiness{CommandReady: true, AutomaticSources: []AutomaticSource{SourceSchedule}, Status: StatusTerminal, Reason: "Task lifecycle is completed."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateReadiness(tt.task, tt.completion); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("EvaluateReadiness = %+v, want %+v", got, tt.want)
			}
		})
	}
}
