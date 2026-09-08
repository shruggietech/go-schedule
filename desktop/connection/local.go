package connection

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/events"
)

var localCapabilities = []string{"tasks", "groups", "chains", "triggers", "watchers", "schedule", "activity", "notifications"}
var localPermissions = []string{"read", "manage"}

type daemonClient interface {
	Health(context.Context) (server.HealthResponse, error)
	StreamEvents(context.Context, func(events.Event)) error
}

// LocalBackend adapts the existing protected local IPC client.
type LocalBackend struct{ daemon daemonClient }

// NewLocalBackend creates the This computer adapter without opening a network listener.
func NewLocalBackend(daemon daemonClient) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) Health(ctx context.Context) (Health, error) {
	value, err := b.daemon.Health(ctx)
	if err != nil {
		return Health{}, safeFailure(err)
	}
	if value.Status != "ok" || !compatibleVersion(value.Version) {
		return Health{}, &Failure{State: StateIncompatible, Message: "This scheduler version is not compatible with the desktop app.", Action: "Update the scheduler service."}
	}
	return Health{Version: value.Version, Capabilities: append([]string(nil), localCapabilities...), Permissions: append([]string(nil), localPermissions...)}, nil
}

func (b *LocalBackend) StreamEvents(ctx context.Context, publish func(DomainEvent)) error {
	return b.daemon.StreamEvents(ctx, func(event events.Event) {
		kind, entityID := safeDomainEvent(event)
		if kind != "" {
			publish(DomainEvent{Kind: kind, EntityID: entityID})
		}
	})
}

func compatibleVersion(version string) bool {
	majorText := strings.SplitN(strings.TrimPrefix(version, "v"), ".", 2)[0]
	major, err := strconv.Atoi(majorText)
	return err == nil && (major == 0 || major == 1)
}

func safeFailure(err error) error {
	var connectionErr *client.ConnectionError
	if errors.As(err, &connectionErr) {
		switch connectionErr.Kind {
		case client.ConnectionAccessDenied:
			return &Failure{State: StateAccessDenied, Message: "Access to the local scheduler service was denied.", Action: "Check the service permissions for this account.", Cause: err}
		case client.ConnectionTimeout:
			return &Failure{State: StateTimedOut, Message: "The local scheduler service did not respond in time.", Action: "Try again.", Cause: err}
		case client.ConnectionUnavailable:
			return &Failure{State: StateUnavailable, Message: "The local scheduler service is unavailable.", Action: "Start the service, then try again.", Cause: err}
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &Failure{State: StateTimedOut, Message: "The local scheduler service did not respond in time.", Action: "Try again.", Cause: err}
	}
	return &Failure{State: StateUnavailable, Message: "The local scheduler service is unavailable.", Action: "Check the service, then try again.", Cause: err}
}

func safeDomainEvent(event events.Event) (string, string) {
	switch event.Kind {
	case events.KindTask:
		if event.Task != nil {
			return "task." + string(event.Task.Verb), event.Task.ID
		}
	case events.KindGroup:
		if event.Group != nil {
			return "group." + string(event.Group.Verb), event.Group.ID
		}
	case events.KindChain:
		if event.Chain != nil {
			return "chain." + string(event.Chain.Verb), event.Chain.ID
		}
	case events.KindTrigger:
		if event.Trigger != nil {
			return "trigger." + string(event.Trigger.Verb), event.Trigger.ID
		}
	case events.KindTriggerSet:
		if event.TriggerSet != nil {
			return "trigger_set." + string(event.TriggerSet.Verb), event.TriggerSet.ID
		}
	case events.KindWatcher:
		if event.Watcher != nil {
			return "filesystem_watcher." + string(event.Watcher.Verb), event.Watcher.ID
		}
	case events.KindRun, events.KindAlert, events.KindLog:
		return string(event.Kind) + ".changed", ""
	}
	return "", ""
}

func localTarget() Target {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform = "macos"
	}
	return Target{ID: "local", DisplayName: "This computer", Platform: platform, Capabilities: []string{}, Permissions: []string{}}
}

func eventMessage(kind string) string {
	return fmt.Sprintf("Scheduler %s changed.", strings.ReplaceAll(kind, "_", " "))
}
