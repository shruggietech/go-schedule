package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/store"
)

func newBrokerServer(t *testing.T) (*Server, *events.Broker) {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	broker := events.NewBroker()
	s := New(st, nil, broker, nil, config.NewLogger(config.Default(), discard{}))
	return s, broker
}

// nextEvent waits briefly for an event of the given kind, draining others.
func nextEvent(t *testing.T, ch <-chan events.Event, kind events.Kind) events.Event {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		select {
		case e := <-ch:
			if e.Kind == kind {
				return e
			}
		case <-deadline:
			t.Fatalf("timed out waiting for %s event", kind)
		}
	}
}

func TestPublish_TaskLifecycle(t *testing.T) {
	s, broker := newBrokerServer(t)
	ch, cancel := broker.Subscribe()
	defer cancel()

	// Create
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "n", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC",
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", rec.Code, rec.Body.String())
	}
	var resp TaskResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	id := resp.Task.ID

	ev := nextEvent(t, ch, events.KindTask)
	if ev.Task == nil || ev.Task.Verb != events.VerbCreated || ev.Task.ID != id {
		t.Fatalf("create event = %+v", ev.Task)
	}

	// Disable (update)
	if rec := doJSON(t, s, http.MethodPost, "/v1/tasks/"+id+"/disable", nil); rec.Code != http.StatusNoContent {
		t.Fatalf("disable status %d", rec.Code)
	}
	if ev := nextEvent(t, ch, events.KindTask); ev.Task.Verb != events.VerbUpdated {
		t.Fatalf("expected updated, got %+v", ev.Task)
	}

	// Delete
	if rec := doJSON(t, s, http.MethodDelete, "/v1/tasks/"+id, nil); rec.Code != http.StatusNoContent {
		t.Fatalf("delete status %d", rec.Code)
	}
	if ev := nextEvent(t, ch, events.KindTask); ev.Task.Verb != events.VerbDeleted || ev.Task.ID != id {
		t.Fatalf("expected deleted, got %+v", ev.Task)
	}
}

func TestPublish_GroupLifecycle(t *testing.T) {
	s, broker := newBrokerServer(t)
	ch, cancel := broker.Subscribe()
	defer cancel()

	rec := doJSON(t, s, http.MethodPost, "/v1/groups", GroupCreateRequest{Name: "G"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create group status %d: %s", rec.Code, rec.Body.String())
	}
	ev := nextEvent(t, ch, events.KindGroup)
	if ev.Group == nil || ev.Group.Verb != events.VerbCreated {
		t.Fatalf("group create event = %+v", ev.Group)
	}
	id := ev.Group.ID

	if rec := doJSON(t, s, http.MethodDelete, "/v1/groups/"+id, nil); rec.Code != http.StatusNoContent {
		t.Fatalf("delete group status %d", rec.Code)
	}
	if ev := nextEvent(t, ch, events.KindGroup); ev.Group.Verb != events.VerbDeleted {
		t.Fatalf("expected group deleted, got %+v", ev.Group)
	}
}
