package viewmodel

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/events"
)

type fakeAPI struct {
	tasks  []domain.Task
	groups []domain.Group
	alerts []domain.Alert
	logs   []domain.LogRecord
}

func (f *fakeAPI) ListTasks(context.Context, string, string) ([]domain.Task, error) {
	return f.tasks, nil
}
func (f *fakeAPI) ListGroups(context.Context) ([]domain.Group, error)       { return f.groups, nil }
func (f *fakeAPI) ListAlerts(context.Context, bool) ([]domain.Alert, error) { return f.alerts, nil }
func (f *fakeAPI) ListLogs(context.Context, string, int) ([]domain.LogRecord, error) {
	return f.logs, nil
}

func TestRefresh_LoadsState(t *testing.T) {
	api := &fakeAPI{
		tasks:  []domain.Task{{ID: "t1", Name: "A"}},
		groups: []domain.Group{{ID: "g1", Name: "G"}},
		alerts: []domain.Alert{{ID: "a1"}},
	}
	m := New(api)
	changed := 0
	m.OnChange = func() { changed++ }

	if err := m.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	s := m.Snapshot()
	if len(s.Tasks) != 1 || len(s.Groups) != 1 || len(s.Alerts) != 1 {
		t.Fatalf("state not loaded: %+v", s)
	}
	if changed != 1 {
		t.Fatalf("OnChange should fire once on refresh, fired %d", changed)
	}
}

func TestApplyEvent_AlertPrependedAndDeduped(t *testing.T) {
	m := New(&fakeAPI{})
	m.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "a1", Message: "first"}})
	m.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "a2", Message: "second"}})
	// Duplicate ID should be ignored.
	m.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "a1", Message: "dup"}})

	s := m.Snapshot()
	if len(s.Alerts) != 2 {
		t.Fatalf("want 2 alerts after dedup, got %d", len(s.Alerts))
	}
	if s.Alerts[0].ID != "a2" {
		t.Fatalf("newest alert should be first, got %s", s.Alerts[0].ID)
	}
}

func TestApplyEvent_LogPrependedAndDeduped(t *testing.T) {
	m := New(&fakeAPI{})
	m.ApplyEvent(events.Event{Kind: events.KindLog, Log: &domain.LogRecord{ID: "l1", Message: "first"}})
	m.ApplyEvent(events.Event{Kind: events.KindLog, Log: &domain.LogRecord{ID: "l2", Message: "second"}})
	m.ApplyEvent(events.Event{Kind: events.KindLog, Log: &domain.LogRecord{ID: "l1", Message: "dup"}})

	s := m.Snapshot()
	if len(s.Logs) != 2 {
		t.Fatalf("want 2 logs after dedup, got %d", len(s.Logs))
	}
	if s.Logs[0].ID != "l2" {
		t.Fatalf("newest log should be first, got %s", s.Logs[0].ID)
	}
}

func TestApplyEvent_TaskUpsertAndRemove(t *testing.T) {
	m := New(&fakeAPI{})
	// created
	m.ApplyEvent(events.Event{Kind: events.KindTask, Task: &events.TaskEvent{
		Verb: events.VerbCreated, ID: "t1", Task: &domain.Task{ID: "t1", Name: "one"}}})
	if got := m.Snapshot().Tasks; len(got) != 1 || got[0].Name != "one" {
		t.Fatalf("after create: %+v", got)
	}
	// updated (same ID, no duplicate)
	m.ApplyEvent(events.Event{Kind: events.KindTask, Task: &events.TaskEvent{
		Verb: events.VerbUpdated, ID: "t1", Task: &domain.Task{ID: "t1", Name: "renamed"}}})
	if got := m.Snapshot().Tasks; len(got) != 1 || got[0].Name != "renamed" {
		t.Fatalf("after update: %+v", got)
	}
	// deleted
	m.ApplyEvent(events.Event{Kind: events.KindTask, Task: &events.TaskEvent{
		Verb: events.VerbDeleted, ID: "t1"}})
	if got := m.Snapshot().Tasks; len(got) != 0 {
		t.Fatalf("after delete: %+v", got)
	}
}

func TestApplyEvent_GroupUpsertAndRemove(t *testing.T) {
	m := New(&fakeAPI{})
	m.ApplyEvent(events.Event{Kind: events.KindGroup, Group: &events.GroupEvent{
		Verb: events.VerbCreated, ID: "g1", Group: &domain.Group{ID: "g1", Name: "G"}}})
	if got := m.Snapshot().Groups; len(got) != 1 {
		t.Fatalf("after create: %+v", got)
	}
	m.ApplyEvent(events.Event{Kind: events.KindGroup, Group: &events.GroupEvent{
		Verb: events.VerbDeleted, ID: "g1"}})
	if got := m.Snapshot().Groups; len(got) != 0 {
		t.Fatalf("after delete: %+v", got)
	}
}

func TestRefresh_LoadsLogs(t *testing.T) {
	m := New(&fakeAPI{logs: []domain.LogRecord{{ID: "x", Message: "hi"}}})
	if err := m.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := m.Snapshot().Logs; len(got) != 1 || got[0].ID != "x" {
		t.Fatalf("logs not loaded on refresh: %+v", got)
	}
}

func TestApplyEvent_RunsCappedAndOrdered(t *testing.T) {
	m := New(&fakeAPI{})
	for i := 0; i < maxRecentRuns+10; i++ {
		m.ApplyEvent(events.Event{Kind: events.KindRun, Run: &domain.Run{ID: string(rune('a' + (i % 26))), TaskID: "t"}})
	}
	s := m.Snapshot()
	if len(s.RecentRuns) != maxRecentRuns {
		t.Fatalf("recent runs should be capped at %d, got %d", maxRecentRuns, len(s.RecentRuns))
	}
}

func TestUnacknowledgedAlerts(t *testing.T) {
	m := New(&fakeAPI{})
	m.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "a1", Acknowledged: false}})
	m.ApplyEvent(events.Event{Kind: events.KindAlert, Alert: &domain.Alert{ID: "a2", Acknowledged: true}})
	if got := m.UnacknowledgedAlerts(); got != 1 {
		t.Fatalf("want 1 unacked, got %d", got)
	}
}
