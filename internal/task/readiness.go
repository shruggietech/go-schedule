package task

import (
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type AutomaticSource string

const (
	SourceSchedule   AutomaticSource = "schedule"
	SourceCompletion AutomaticSource = "completion"
)

type ReadinessStatus string

const (
	StatusNotRunnable ReadinessStatus = "not_runnable"
	StatusManualOnly  ReadinessStatus = "manual_only"
	StatusReady       ReadinessStatus = "ready"
	StatusDisabled    ReadinessStatus = "disabled"
	StatusTerminal    ReadinessStatus = "terminal"
)

type Readiness struct {
	CommandReady     bool              `json:"command_ready"`
	ActivationReady  bool              `json:"activation_ready"`
	AutomaticSources []AutomaticSource `json:"automatic_sources"`
	Status           ReadinessStatus   `json:"status"`
	Reason           string            `json:"reason"`
}

func DisplayName(t domain.Task) string {
	if strings.TrimSpace(t.Name) == "" {
		return "unnamed"
	}
	return t.Name
}

func EvaluateReadiness(t domain.Task, hasCompletionSource bool) Readiness {
	r := Readiness{CommandReady: strings.TrimSpace(t.Command) != ""}
	if t.ScheduleID != "" {
		r.AutomaticSources = append(r.AutomaticSources, SourceSchedule)
	}
	if hasCompletionSource {
		r.AutomaticSources = append(r.AutomaticSources, SourceCompletion)
	}
	r.ActivationReady = r.CommandReady && t.State == domain.TaskActive && len(r.AutomaticSources) > 0
	switch {
	case t.State == domain.TaskCompleted || t.State == domain.TaskDisabled:
		r.Status = StatusTerminal
		r.Reason = "Task lifecycle is " + string(t.State) + "."
	case !r.CommandReady:
		r.Status = StatusNotRunnable
		r.Reason = "No command is configured."
	case !r.ActivationReady:
		r.Status = StatusManualOnly
		r.Reason = "No automatic activation source is configured."
	case !t.Enabled:
		r.Status = StatusDisabled
		r.Reason = "Automatic activation is disabled."
	default:
		r.Status = StatusReady
	}
	return r
}
