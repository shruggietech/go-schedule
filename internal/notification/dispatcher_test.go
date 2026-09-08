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

type refillStore struct {
	mu    sync.Mutex
	queue []domain.NotificationDelivery
}

func (s *refillStore) RecoverNotificationDeliveries(time.Time) (int64, error) { return 0, nil }
func (s *refillStore) ClaimNotificationDeliveries(limit int, _ time.Time) ([]domain.NotificationDelivery, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit != 1 || len(s.queue) == 0 {
		return nil, nil
	}
	delivery := s.queue[0]
	s.queue = s.queue[1:]
	delivery.Attempts++
	return []domain.NotificationDelivery{delivery}, nil
}
func (*refillStore) RetryNotificationDelivery(string, time.Time, int, string) error { return nil }
func (*refillStore) CompleteNotificationDelivery(string, bool, int, string, time.Time) error {
	return nil
}

type refillSender struct {
	slowRelease chan struct{}
	started     chan string
}

func (s refillSender) Send(ctx context.Context, delivery domain.NotificationDelivery) (int, string, error) {
	s.started <- delivery.ID
	if delivery.ID == "slow" {
		select {
		case <-s.slowRelease:
		case <-ctx.Done():
			return 0, "cancelled", ctx.Err()
		}
	}
	return 204, "", nil
}

func TestDispatcherRefillsAvailableWorkerBeforeSlowAttemptFinishes(t *testing.T) {
	st := &refillStore{queue: []domain.NotificationDelivery{{ID: "slow"}, {ID: "fast-1"}, {ID: "fast-2"}}}
	sender := refillSender{slowRelease: make(chan struct{}), started: make(chan string, 3)}
	dispatcher := NewDispatcher(st, sender, slog.New(slog.NewTextHandler(io.Discard, nil)), 2)
	done := make(chan error, 1)
	go func() { done <- dispatcher.process(context.Background()) }()
	started := map[string]bool{}
	for len(started) < 3 {
		select {
		case id := <-sender.started:
			started[id] = true
		case <-time.After(time.Second):
			t.Fatalf("workers stalled before refill: started=%v", started)
		}
	}
	close(sender.slowRelease)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
