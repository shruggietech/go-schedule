package watcher

import (
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clock"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func observationFields(t *testing.T, health domain.WatcherHealth) (string, time.Time) {
	t.Helper()
	data, err := json.Marshal(health)
	if err != nil {
		t.Fatal(err)
	}
	var fields struct {
		Revision string    `json:"observed_definition_revision"`
		At       time.Time `json:"observed_at"`
	}
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if fields.Revision == "" || fields.At.IsZero() {
		t.Fatalf("missing actual definition observation: %s", data)
	}
	return fields.Revision, fields.At
}

func TestObservationAdvancesForSameHealthAndCoalescedDefinitions(t *testing.T) {
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	clk := clock.NewFake(now)
	definition := domain.FilesystemWatcher{ID: "watcher-1", Name: "Inbox", Kind: domain.WatcherDirectory, Path: canonicalTempDir(t), Pattern: "*.csv", Enabled: true, UpdatedAt: now}
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{definition}}
	manager := newManager(store, &dispatcherStub{ch: make(chan string, 1)}, clk, slog.Default(), func() (observer, error) { return newFakeObserver(), nil })
	var reports []domain.WatcherHealth
	manager.SetHealthReporter(func(_ domain.FilesystemWatcher, health domain.WatcherHealth) { reports = append(reports, health) })
	state := runtimeState{manager: manager}
	t.Cleanup(state.closeObserver)
	state.rebuild()
	first := manager.Health(definition.ID)
	firstRevision, firstAt := observationFields(t, first)
	// Two saved edits share a timestamp and are coalesced into one runtime reload.
	store.items[0].Pattern = "*.txt"
	manager.Reload()
	store.items[0].Pattern = "*.json"
	manager.Reload()
	clk.Advance(time.Second)
	state.rebuild()
	second := manager.Health(definition.ID)
	secondRevision, secondAt := observationFields(t, second)
	if second.State != first.State || second.Reason != first.Reason || !second.ChangedAt.Equal(first.ChangedAt) {
		t.Fatalf("unchanged health became a transition: first=%+v second=%+v", first, second)
	}
	if secondRevision == firstRevision || !secondAt.After(firstAt) {
		t.Fatalf("new definition was not observed: first=%+v second=%+v", first, second)
	}
	if len(reports) != 2 || reports[1].ObservedDefinitionRevision != secondRevision {
		t.Fatalf("same-health reconciliation was not published to consumers: %+v", reports)
	}
	// A later event from the older generation must remain associated with its old definition.
	manager.setHealth(definition, domain.WatcherDegraded, "old generation event")
	staleRevision, _ := observationFields(t, manager.Health(definition.ID))
	if staleRevision != firstRevision {
		t.Fatal("old observation was incorrectly stamped with the latest saved definition")
	}
}
