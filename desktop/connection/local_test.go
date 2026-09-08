package connection

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/events"
)

type localDaemonFake struct {
	health server.HealthResponse
	err    error
	event  events.Event
}

func (f localDaemonFake) Health(context.Context) (server.HealthResponse, error) {
	return f.health, f.err
}
func (f localDaemonFake) StreamEvents(ctx context.Context, publish func(events.Event)) error {
	if f.event.Kind != "" {
		publish(f.event)
	}
	<-ctx.Done()
	return ctx.Err()
}

func TestLocalBackendNegotiatesSafeThisComputerContract(t *testing.T) {
	backend := NewLocalBackend(localDaemonFake{health: server.HealthResponse{Status: "ok", Version: "1.2.0"}})
	health, err := backend.Health(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if health.Version != "1.2.0" || len(health.Capabilities) == 0 || len(health.Permissions) == 0 {
		t.Fatalf("health=%+v", health)
	}
	if !contains(health.Capabilities, "notifications") {
		t.Fatalf("notifications capability missing: %+v", health.Capabilities)
	}
	target := localTarget()
	if target.ID != "local" || target.DisplayName != "This computer" {
		t.Fatalf("target=%+v", target)
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func TestLocalBackendRejectsIncompatibleVersion(t *testing.T) {
	backend := NewLocalBackend(localDaemonFake{health: server.HealthResponse{Status: "ok", Version: "2.0.0"}})
	_, err := backend.Health(context.Background())
	var failure *Failure
	if !errors.As(err, &failure) || failure.State != StateIncompatible {
		t.Fatalf("error=%v", err)
	}
}

func TestLocalBackendMapsErrorsWithoutSensitiveDetails(t *testing.T) {
	sensitive := `C:\private\daemon.pipe`
	err := client.NewConnectionError("health "+sensitive, os.ErrPermission)
	_, got := NewLocalBackend(localDaemonFake{err: err}).Health(context.Background())
	var failure *Failure
	if !errors.As(got, &failure) || failure.State != StateAccessDenied || strings.Contains(failure.Message+failure.Action, sensitive) {
		t.Fatalf("failure=%+v", failure)
	}
}

func TestLocalBackendSanitizesEvents(t *testing.T) {
	raw := events.Event{Kind: events.KindTrigger, Trigger: &events.TriggerEvent{Verb: events.VerbUpdated, ID: "trigger-1"}}
	backend := NewLocalBackend(localDaemonFake{event: raw})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	seen := make(chan DomainEvent, 1)
	go func() { _ = backend.StreamEvents(ctx, func(event DomainEvent) { seen <- event; cancel() }) }()
	event := <-seen
	if event.Kind != "trigger.updated" || event.EntityID != "trigger-1" {
		t.Fatalf("event=%+v", event)
	}
}
