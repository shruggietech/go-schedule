// Package desktopcontrol projects the installed local service state for the
// Windows desktop surfaces. It never confuses a selected remote daemon with
// the service on this computer.
package desktopcontrol

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shruggietech/go-schedule/internal/service"
)

// Snapshot is one bounded observation of the installed local service.
type Snapshot struct {
	State      string `json:"state"`
	SCMState   string `json:"scmState"`
	Detail     string `json:"detail"`
	ObservedAt string `json:"observedAt"`
}

// Monitor owns read-only sources. It is safe to use concurrently.
type Monitor struct {
	Query   func() (service.State, error)
	Health  func(context.Context) error
	Execute func(string) error
}

// NewMonitor constructs a monitor with the installed-service query.
func NewMonitor(health func(context.Context) error) Monitor {
	return Monitor{Query: service.QueryState, Health: health, Execute: requestElevation}
}

// Observe reads the SCM first. A Running service requires a fresh health
// response before it can be represented as healthy.
func (m Monitor) Observe(ctx context.Context) Snapshot {
	at := time.Now().UTC().Format(time.RFC3339)
	if m.Query == nil {
		return Snapshot{State: "unknown", SCMState: "unknown", Detail: "Service status is unavailable.", ObservedAt: at}
	}
	st, err := m.Query()
	if err != nil {
		return Snapshot{State: "unknown", SCMState: "unknown", Detail: fmt.Sprintf("Cannot read local service status: %v", err), ObservedAt: at}
	}
	out := Snapshot{State: string(st), SCMState: string(st), ObservedAt: at}
	switch st {
	case service.StateRunning:
		if m.Health == nil {
			out.State = "unreachable"
			out.Detail = "The service is running, but daemon health is unavailable."
			return out
		}
		checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := m.Health(checkCtx); err != nil {
			out.State = "unreachable"
			out.Detail = fmt.Sprintf("The service is running, but the local daemon is unreachable: %v", err)
			return out
		}
		out.Detail = "Local daemon is running and healthy."
	case service.StateStopped:
		out.Detail = "Local daemon stopped. Scheduled tasks on this computer are not running."
	case service.StateNotInstalled:
		out.Detail = "The local Windows service is not installed."
	case service.StateStarting:
		out.Detail = "The local Windows service is starting."
	case service.StateStopping:
		out.Detail = "The local Windows service is stopping."
	default:
		out.State = "unknown"
		out.Detail = "The local Windows service state is unknown."
	}
	return out
}

// ActionResult includes only an observed outcome, never optimistic success.
type ActionResult struct {
	Action   string   `json:"action"`
	Outcome  string   `json:"outcome"`
	Message  string   `json:"message"`
	Snapshot Snapshot `json:"snapshot"`
}

// RequestAction starts a narrow privileged service operation and waits for
// the observed target state. The caller owns Stop/Restart confirmation.
func (m Monitor) RequestAction(ctx context.Context, action string) ActionResult {
	result := ActionResult{Action: action}
	if action != "start" && action != "stop" && action != "restart" {
		result.Outcome, result.Message = "rejected", "Unsupported service action."
		return result
	}
	before := m.Observe(ctx)
	result.Snapshot = before
	if before.SCMState == string(service.StateNotInstalled) {
		result.Outcome, result.Message = "unavailable", "The local service is not installed. Install the desktop product first."
		return result
	}
	if before.SCMState == string(service.StateUnknown) {
		result.Outcome, result.Message = "unavailable", "The local service state cannot be read. Check Windows service permissions."
		return result
	}
	if m.Execute == nil {
		result.Outcome, result.Message = "unavailable", "Windows service control is unavailable."
		return result
	}
	if err := m.Execute(action); err != nil {
		result.Outcome = "failed"
		if errors.Is(err, errElevationCancelled) {
			result.Outcome, result.Message = "cancelled", "Windows elevation was cancelled. The service was not changed."
		} else {
			result.Message = fmt.Sprintf("Could not request %s: %v", action, err)
		}
		result.Snapshot = m.Observe(ctx)
		return result
	}
	target := "running"
	if action == "stop" {
		target = "stopped"
	}
	waitCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()
	for {
		result.Snapshot = m.Observe(waitCtx)
		if result.Snapshot.State == target {
			result.Outcome = "accepted"
			result.Message = fmt.Sprintf("Local service %s confirmed.", target)
			return result
		}
		select {
		case <-waitCtx.Done():
			result.Outcome = "failed"
			result.Message = fmt.Sprintf("The %s request did not reach a healthy %s state. Check Windows Services and try again.", action, target)
			return result
		case <-ticker.C:
		}
	}
}
