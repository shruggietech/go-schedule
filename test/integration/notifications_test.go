package integration

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/notification"
	"github.com/shruggietech/go-schedule/internal/store"
)

type observingNotificationStore struct {
	*store.Store
	completed chan struct{}
}

func (s *observingNotificationStore) CompleteNotificationDelivery(id string, succeeded bool, status int, diagnostic string, now time.Time) error {
	err := s.Store.CompleteNotificationDelivery(id, succeeded, status, diagnostic, now)
	if err == nil {
		close(s.completed)
	}
	return err
}

func TestWebhookNotificationEndToEndPreservesRunOutcome(t *testing.T) {
	received := make(chan string, 1)
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- r.Header.Get("X-Go-Schedule-Delivery")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer receiver.Close()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := domain.Task{Name: "failing task", Command: "false", Enabled: true, Timezone: "UTC", State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	channel := domain.NotificationChannel{Name: "receiver", Kind: domain.NotificationChannelWebhook, Endpoint: receiver.URL, EndpointSummary: receiver.URL, Enabled: true}
	if err := st.CreateNotificationChannel(&channel); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: channel.ID, OnFailure: true}}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	exit := 1
	run := domain.Run{TaskID: task.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: domain.OutcomeFailure, ExitCode: &exit, Output: "sensitive output", Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	observed := &observingNotificationStore{Store: st, completed: make(chan struct{})}
	dispatcher := notification.NewDispatcher(observed, notification.NewSender(nil), slog.New(slog.NewTextHandler(io.Discard, nil)), 1)
	go func() { done <- dispatcher.Run(ctx) }()
	select {
	case deliveryID := <-received:
		if deliveryID == "" {
			t.Fatal("receiver did not get stable delivery ID")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("webhook was not received")
	}
	select {
	case <-observed.completed:
	case <-time.After(2 * time.Second):
		t.Fatal("delivery did not reach terminal state")
	}
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{RunID: run.ID})
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.NotificationDeliverySucceeded {
		t.Fatalf("deliveries=%+v err=%v", deliveries, err)
	}
	persisted, err := st.GetRun(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if persisted.Outcome != domain.OutcomeFailure || persisted.ExitCode == nil || *persisted.ExitCode != 1 || persisted.Output != "sensitive output" {
		t.Fatalf("source run changed: %+v", persisted)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("dispatcher did not stop")
	}
}
