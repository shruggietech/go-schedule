package systems

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type profileStoreFake struct{ collection clientprofile.Collection }

func (f profileStoreFake) Load() (clientprofile.Collection, error) { return f.collection, nil }

type profileStoreFunc func() (clientprofile.Collection, error)

func (f profileStoreFunc) Load() (clientprofile.Collection, error) { return f() }

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

type summaryFake struct {
	value domain.SystemSummary
	err   error
}

func (f summaryFake) SystemSummary(context.Context) (domain.SystemSummary, error) {
	return f.value, f.err
}

func TestRefreshBoundsConcurrencyAndPreservesPartialResults(t *testing.T) {
	now := time.Date(2026, 9, 21, 15, 0, 0, 0, time.UTC)
	var active, peak atomic.Int32
	release := make(chan struct{})
	targets := make([]target, 8)
	for i := range targets {
		key := string(rune('a' + i))
		targets[i] = target{registration: Registration{Key: key, Kind: "remote", Label: key}, backend: backendFake{health: connection.Health{ID: key, DisplayName: key, Platform: "linux", Version: "1.0.0"}, active: &active, peak: &peak, release: release}, client: summaryFake{value: domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: now}}}
	}
	profiles := make([]clientprofile.Profile, 8)
	for i := range profiles {
		profiles[i].ID = targets[i].registration.Key
	}
	service := &Service{profiles: profileStoreFake{clientprofile.Collection{Version: 1, Profiles: profiles}}, cache: map[string]cachedSummary{}, now: func() time.Time { return now }, timeout: time.Second, concurrency: 4, loadTargets: func() ([]target, []Observation) { return targets, nil }}
	done := make(chan Snapshot, 1)
	go func() { done <- service.Refresh(context.Background()) }()
	deadline := time.After(time.Second)
	for peak.Load() < 4 {
		select {
		case <-deadline:
			t.Fatal("workers did not start")
		default:
		}
	}
	if peak.Load() != 4 {
		t.Fatalf("peak concurrency = %d", peak.Load())
	}
	close(release)
	snapshot := <-done
	if len(snapshot.Observations) != 8 {
		t.Fatalf("observations = %d", len(snapshot.Observations))
	}
}

func TestRefreshRetainsSessionSummaryAsStale(t *testing.T) {
	now := time.Date(2026, 9, 21, 15, 0, 0, 0, time.UTC)
	profile := clientprofile.Profile{ID: "remote"}
	store := profileStoreFake{clientprofile.Collection{Version: 1, Profiles: []clientprofile.Profile{profile}}}
	service := &Service{profiles: store, cache: map[string]cachedSummary{"remote": {summary: domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: now.Add(-time.Minute), ActiveTaskCount: 2}, at: now.Add(-time.Minute)}}, now: func() time.Time { return now }, timeout: time.Second, concurrency: 1}
	service.loadTargets = func() ([]target, []Observation) {
		return []target{{registration: Registration{Key: "remote", Kind: "remote", Label: "Remote"}, backend: backendFake{err: &connection.Failure{State: connection.StateTimedOut, Message: "Timed out.", Action: "Try again."}}, client: summaryFake{}}}, nil
	}
	snapshot := service.Refresh(context.Background())
	if len(snapshot.Observations) != 1 || !snapshot.Observations[0].Stale || snapshot.Observations[0].Summary.ActiveTaskCount != 2 || snapshot.Observations[0].State != connection.StateTimedOut {
		t.Fatalf("observation = %+v", snapshot.Observations)
	}
}

func TestRefreshDropsRemovedRegistration(t *testing.T) {
	store := profileStoreFake{clientprofile.Collection{Version: 1, Profiles: []clientprofile.Profile{}}}
	service := &Service{profiles: store, cache: map[string]cachedSummary{}, now: time.Now, timeout: time.Second, concurrency: 1}
	service.loadTargets = func() ([]target, []Observation) {
		return []target{{registration: Registration{Key: "removed", Kind: "remote", Label: "Removed"}, backend: backendFake{err: errors.New("offline")}, client: summaryFake{}}}, nil
	}
	if got := service.Refresh(context.Background()); len(got.Observations) != 0 {
		t.Fatalf("observations = %+v", got.Observations)
	}
}

func TestRefreshRetainsCompletedObservationsWhenMembershipRecheckFails(t *testing.T) {
	now := time.Date(2026, 9, 21, 15, 0, 0, 0, time.UTC)
	service := &Service{profiles: profileStoreFunc(func() (clientprofile.Collection, error) {
		return clientprofile.Collection{}, errors.New("profile store unavailable")
	}), cache: map[string]cachedSummary{}, now: func() time.Time { return now }, timeout: time.Second, concurrency: 1}
	service.loadTargets = func() ([]target, []Observation) {
		return []target{{registration: Registration{Key: "remote", Kind: "remote", Label: "Remote"}, backend: backendFake{health: connection.Health{ID: "remote", DisplayName: "Remote", Platform: "linux", Version: "1.0.0"}}, client: summaryFake{value: domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: now}}}}, nil
	}
	snapshot := service.Refresh(context.Background())
	if len(snapshot.Observations) != 1 || snapshot.Observations[0].Registration.Key != "remote" {
		t.Fatalf("observations = %+v", snapshot.Observations)
	}
}

func TestRefreshPublishesCompletedTargetsBeforeSlowTargetsFinish(t *testing.T) {
	now := time.Date(2026, 9, 21, 15, 0, 0, 0, time.UTC)
	release := make(chan struct{})
	targets := []target{
		{registration: Registration{Key: "fast", Kind: "remote", Label: "Fast"}, backend: backendFake{health: connection.Health{ID: "fast", DisplayName: "Fast", Platform: "linux", Version: "1.0.0"}}, client: summaryFake{value: domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: now}}},
		{registration: Registration{Key: "slow", Kind: "remote", Label: "Slow"}, backend: backendFake{health: connection.Health{ID: "slow", DisplayName: "Slow", Platform: "linux", Version: "1.0.0"}, release: release}, client: summaryFake{value: domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: now}}},
	}
	profiles := []clientprofile.Profile{{ID: "fast"}, {ID: "slow"}}
	service := &Service{profiles: profileStoreFake{clientprofile.Collection{Version: 1, Profiles: profiles}}, cache: map[string]cachedSummary{}, now: func() time.Time { return now }, timeout: time.Second, concurrency: 2, loadTargets: func() ([]target, []Observation) { return targets, nil }}
	updates := make(chan Snapshot, 4)
	done := make(chan Snapshot, 1)
	go func() { done <- service.Refresh(context.Background(), func(snapshot Snapshot) { updates <- snapshot }) }()
	select {
	case update := <-updates:
		if update.Complete || len(update.Observations) != 1 || update.Observations[0].Registration.Key != "fast" {
			t.Fatalf("first update = %+v", update)
		}
	case <-time.After(time.Second):
		t.Fatal("fast target was not published while slow target remained blocked")
	}
	close(release)
	if final := <-done; !final.Complete || len(final.Observations) != 2 {
		t.Fatalf("final snapshot = %+v", final)
	}
}
