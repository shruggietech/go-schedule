package notification

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type dispatcherStore struct {
	mu         sync.Mutex
	delivery   domain.NotificationDelivery
	completed  bool
	retried    bool
	claimLimit int
}

func (s *dispatcherStore) RecoverNotificationDeliveries(time.Time) (int64, error) { return 0, nil }
func (s *dispatcherStore) ClaimNotificationDeliveries(limit int, _ time.Time) ([]domain.NotificationDelivery, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.claimLimit = limit
	if s.completed || s.retried {
		return nil, nil
	}
	s.delivery.Attempts++
	s.delivery.State = domain.NotificationDeliveryClaimed
	return []domain.NotificationDelivery{s.delivery}, nil
}
func (s *dispatcherStore) RetryNotificationDelivery(_ string, _ time.Time, _ int, _ string) error {
	s.mu.Lock()
	s.retried = true
	s.mu.Unlock()
	return nil
}
func (s *dispatcherStore) CompleteNotificationDelivery(_ string, succeeded bool, _ int, _ string, _ time.Time) error {
	s.mu.Lock()
	s.completed = succeeded
	s.mu.Unlock()
	return nil
}

type dispatcherSender struct{ err error }

func (s dispatcherSender) Send(context.Context, domain.NotificationDelivery) (int, string, error) {
	if s.err != nil {
		return 503, "receiver returned HTTP 503", s.err
	}
	return 204, "", nil
}

func TestDispatcherCompletesOutsideCallerAndSchedulesBoundedRetry(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	success := &dispatcherStore{delivery: domain.NotificationDelivery{ID: "ok", Attempts: 0}}
	dispatcher := NewDispatcher(success, dispatcherSender{}, log, 1)
	if err := dispatcher.process(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !success.completed {
		t.Fatal("delivery was not completed")
	}
	if success.claimLimit != 1 {
		t.Fatalf("claim limit=%d, want one available worker", success.claimLimit)
	}
	failure := &dispatcherStore{delivery: domain.NotificationDelivery{ID: "retry", Attempts: 0}}
	dispatcher = NewDispatcher(failure, dispatcherSender{err: context.DeadlineExceeded}, log, 1)
	if err := dispatcher.process(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !failure.retried || failure.completed {
		t.Fatalf("retry=%t completed=%t", failure.retried, failure.completed)
	}
}
