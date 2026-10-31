package integration

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/notification"
	"github.com/shruggietech/go-schedule/v2/internal/store"
	"github.com/shruggietech/go-schedule/v2/internal/testsmtp"
)

type smtpCompletionStore struct {
	*store.Store
	completed chan string
}

func (s *smtpCompletionStore) CompleteNotificationDelivery(id string, ok bool, status int, diagnostic string, now time.Time) error {
	err := s.Store.CompleteNotificationDelivery(id, ok, status, diagnostic, now)
	if err == nil {
		s.completed <- id
	}
	return err
}
func (s *smtpCompletionStore) FinishSMTPAttempt(d domain.NotificationDelivery, result domain.SMTPAttemptResult, next, now time.Time) error {
	err := s.Store.FinishSMTPAttempt(d, result, next, now)
	if err == nil {
		s.completed <- d.ID
	}
	return err
}

func TestSMTPDispatcherHealthTestAndFailureIsolation(t *testing.T) {
	receiver, err := testsmtp.Start(testsmtp.Options{Reject: map[string]int{"bad@example.test": 550}})
	if err != nil {
		t.Fatal(err)
	}
	defer receiver.Close()
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	defer webhook.Close()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	host, portText, err := net.SplitHostPort(receiver.Addr())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	cfg := domain.SMTPConfig{Host: host, Port: port, TLSMode: "plain", From: "sender@example.test", Recipients: []string{"bad@example.test"}, TimeoutSeconds: 2}
	failing := domain.NotificationChannel{Name: "mail failure", Kind: domain.NotificationChannelSMTP, SMTP: &cfg, Enabled: true}
	if err := st.CreateNotificationChannel(&failing); err != nil {
		t.Fatal(err)
	}
	cfg.Recipients = []string{"good@example.test"}
	health := domain.NotificationChannel{Name: "mail health", Kind: domain.NotificationChannelSMTP, SMTP: &cfg, Enabled: true, HealthIntervalSeconds: 60}
	if err := st.CreateNotificationChannel(&health); err != nil {
		t.Fatal(err)
	}
	hook := domain.NotificationChannel{Name: "webhook", Kind: domain.NotificationChannelWebhook, Endpoint: webhook.URL, EndpointSummary: webhook.URL, Enabled: true}
	if err := st.CreateNotificationChannel(&hook); err != nil {
		t.Fatal(err)
	}
	task := domain.Task{Name: "unaffected task", Command: "true", Enabled: true, Timezone: "UTC", State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: failing.ID, OnSuccess: true}, {ChannelID: hook.ID, OnSuccess: true}}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	exit := 0
	run := domain.Run{TaskID: task.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: domain.OutcomeSuccess, ExitCode: &exit, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTestNotificationDelivery(health.ID); err != nil {
		t.Fatal(err)
	}
	observed := &smtpCompletionStore{Store: st, completed: make(chan string, 8)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- notification.NewDispatcher(observed, notification.NewSender(nil), slog.New(slog.NewTextHandler(io.Discard, nil)), 2).Run(ctx)
	}()
	defer func() {
		cancel()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Error("dispatcher did not stop")
		}
	}()
	for range 4 {
		select {
		case <-observed.completed:
		case <-time.After(3 * time.Second):
			t.Fatal("notification did not reach a recorded outcome")
		}
	}
	values, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[domain.NotificationEventKind]bool{}
	for _, d := range values {
		if d.Kind == domain.NotificationChannelSMTP {
			seen[d.EventKind] = true
			expected := domain.NotificationDeliverySucceeded
			if d.ChannelID == failing.ID {
				expected = domain.NotificationDeliveryFailed
			}
			if d.Diagnostic == "" || (d.State == domain.NotificationDeliverySucceeded && d.LastError != "") || (d.State != domain.NotificationDeliverySucceeded && d.LastError != d.Diagnostic) {
				t.Fatalf("transport detail mislabeled: %+v", d)
			}
			if d.State != expected || len(d.Recipients) != 1 {
				t.Fatalf("SMTP result %+v", d)
			}
		} else if d.State != domain.NotificationDeliverySucceeded {
			t.Fatalf("webhook affected %+v", d)
		}
	}
	if len(values) != 4 || !seen[domain.NotificationEventTest] || !seen[domain.NotificationEventDaemonHealth] || !seen[domain.NotificationEventRunCompleted] || len(receiver.Messages()) != 2 {
		t.Fatalf("missing dispatch evidence %+v", values)
	}
	persisted, err := st.GetRun(run.ID)
	if err != nil || persisted.Outcome != domain.OutcomeSuccess || persisted.ExitCode == nil || *persisted.ExitCode != 0 {
		t.Fatalf("mail failure changed run %+v %v", persisted, err)
	}
}
