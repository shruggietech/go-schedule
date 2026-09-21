package search

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/systems"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type backendFake struct {
	health  connection.Health
	err     error
	active  *atomic.Int32
	peak    *atomic.Int32
	release <-chan struct{}
}

func (f backendFake) Health(ctx context.Context) (connection.Health, error) {
	if f.active != nil {
		value := f.active.Add(1)
		defer f.active.Add(-1)
		for {
			old := f.peak.Load()
			if value <= old || f.peak.CompareAndSwap(old, value) {
				break
			}
		}
	}
	if f.release != nil {
		select {
		case <-f.release:
		case <-ctx.Done():
			return connection.Health{}, ctx.Err()
		}
	}
	return f.health, f.err
}
func (backendFake) StreamEvents(context.Context, func(connection.DomainEvent)) error { return nil }

type daemonFake struct {
	search domain.DaemonSearch
	task   server.TaskResponse
	alerts []domain.Alert
	err    error
	action func(domain.SearchAction, string) error
}

func (f *daemonFake) Search(context.Context, string, []domain.SearchKind, int) (domain.DaemonSearch, error) {
	return f.search, f.err
}
func (f *daemonFake) GetTask(context.Context, string) (server.TaskResponse, error) {
	return f.task, f.err
}
func (f *daemonFake) SetTaskEnabled(_ context.Context, id string, enabled bool) error {
	if enabled {
		return f.call(domain.SearchActionEnable, id)
	}
	return f.call(domain.SearchActionDisable, id)
}
func (f *daemonFake) RunNow(_ context.Context, id string) error {
	return f.call(domain.SearchActionRunNow, id)
}
func (f *daemonFake) ListAlertsLimited(context.Context, bool, int) ([]domain.Alert, error) {
	return f.alerts, f.err
}
func (f *daemonFake) AckAlert(_ context.Context, id string) error {
	return f.call(domain.SearchActionAcknowledge, id)
}
func (f *daemonFake) call(action domain.SearchAction, id string) error {
	if f.action != nil {
		return f.action(action, id)
	}
	return nil
}

func TestSearchBoundsConcurrencyAndPublishesSourceIdentity(t *testing.T) {
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	var active, peak atomic.Int32
	release := make(chan struct{})
	targets := make([]target, 12)
	for index := range targets {
		key := string(rune('a' + index))
		daemon := &daemonFake{search: domain.DaemonSearch{Schema: domain.DaemonSearchSchema, ObservedAt: now, Results: []domain.SearchMatch{{Kind: domain.SearchKindTask, ObjectID: "task-" + key, Name: "Duplicate", ActionHints: []domain.SearchAction{domain.SearchActionRunNow}}}}}
		targets[index] = target{registration: systems.Registration{Key: key, Kind: "remote", Label: "Same label"}, backend: backendFake{health: connection.Health{ID: "daemon-" + key, Permissions: []string{"read", "operate"}}, active: &active, peak: &peak, release: release}, client: daemon}
	}
	service := &Service{now: func() time.Time { return now }, timeout: time.Second, concurrency: 8, loadTargets: func() ([]target, []Observation) { return targets, nil }}
	done := make(chan Snapshot, 1)
	go func() { done <- service.Search(context.Background(), Request{Query: "duplicate"}) }()
	deadline := time.After(time.Second)
	for peak.Load() < 8 {
		select {
		case <-deadline:
			t.Fatal("workers did not reach the bound")
		default:
		}
	}
	if peak.Load() != 8 {
		t.Fatalf("peak=%d", peak.Load())
	}
	close(release)
	final := <-done
	if !final.Complete || len(final.Observations) != 12 {
		t.Fatalf("snapshot=%+v", final)
	}
	for _, observation := range final.Observations {
		if len(observation.Matches) != 1 || observation.Matches[0].ExpectedDaemonID == "" || observation.Matches[0].SourceShortID == "" {
			t.Fatalf("observation=%+v", observation)
		}
	}
}

func TestSearchCancelsPriorGeneration(t *testing.T) {
	release := make(chan struct{})
	service := &Service{now: time.Now, timeout: time.Second, concurrency: 1}
	service.loadTargets = func() ([]target, []Observation) {
		return []target{{registration: systems.Registration{Key: "remote"}, backend: backendFake{release: release}, client: &daemonFake{}}}, nil
	}
	first := make(chan Snapshot, 1)
	go func() { first <- service.Search(context.Background(), Request{Query: "one"}) }()
	time.Sleep(10 * time.Millisecond)
	second := service.Search(context.Background(), Request{Query: "two"})
	close(release)
	if second.Generation != 2 || (<-first).Generation != 1 {
		t.Fatalf("generations were not isolated")
	}
}

func TestObserveOnlyResultsExposeOpenButNotMutation(t *testing.T) {
	now := time.Now().UTC()
	daemon := &daemonFake{search: domain.DaemonSearch{Schema: domain.DaemonSearchSchema, ObservedAt: now, Results: []domain.SearchMatch{{Kind: domain.SearchKindTask, ObjectID: "task", Name: "Task", ActionHints: []domain.SearchAction{domain.SearchActionRunNow}}}}}
	service := &Service{now: time.Now, timeout: time.Second, concurrency: 1, loadTargets: func() ([]target, []Observation) {
		return []target{{registration: systems.Registration{Key: "remote"}, backend: backendFake{health: connection.Health{ID: "daemon", Permissions: []string{"read"}}}, client: daemon}}, nil
	}}
	match := service.Search(context.Background(), Request{Query: "task"}).Observations[0].Matches[0]
	if len(match.AvailableActions) != 1 || match.AvailableActions[0] != domain.SearchActionOpen || match.DisabledReason == "" {
		t.Fatalf("match=%+v", match)
	}
}
