package operations

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type cancellationBackend struct {
	stubBackend
	cancelCalls int
	cancelID    string
	ack         domain.RunCancellation
	cancelErr   error
}

func (b *cancellationBackend) CancelRun(_ context.Context, id string) (domain.RunCancellation, error) {
	b.cancelCalls++
	b.cancelID = id
	return b.ack, b.cancelErr
}
func (b *cancellationBackend) GetRun(ctx context.Context, id string) (domain.Run, error) {
	b.calls = append(b.calls, "exact")
	for _, r := range b.runs {
		if r.ID == id {
			return r, nil
		}
	}
	return domain.Run{}, &client.StatusError{Code: server.CodeNotFound, Message: "run missing"}
}
func TestCancelRunKeepsAcknowledgementSeparateFromObservedCompletion(t *testing.T) {
	b := &cancellationBackend{ack: domain.RunCancellation{RunID: "run-1", TaskID: "task-1", Status: "requested"}}
	r := NewService(b).CancelRun(context.Background(), "run-1", "daemon-1")
	if r.Outcome != "accepted" || r.Cancellation == nil || r.Cancellation.Status != "requested" || b.cancelCalls != 1 || b.cancelID != "run-1" || r.Activity != nil {
		t.Fatalf("result=%+v calls=%d", r, b.cancelCalls)
	}
}
func TestExactActiveRecordPrefersPersistedCompletion(t *testing.T) {
	now := time.Now()
	b := &cancellationBackend{stubBackend: stubBackend{active: []domain.Run{{ID: "run-1", TaskID: "task", StartedAt: &now}}, runs: []domain.Run{{ID: "run-1", TaskID: "task", EndedAt: &now, Outcome: domain.OutcomeCanceled}}}}
	s := NewService(b)
	r := s.ActivityRecord(context.Background(), "run", "run-1")
	if r.Outcome != "accepted" || r.Activity.Runs[0].State != "canceled" || !reflect.DeepEqual(b.calls, []string{"active", "exact"}) {
		t.Fatalf("r=%+v calls=%v", r, b.calls)
	}
	b.runs = nil
	b.calls = nil
	r = s.ActivityRecord(context.Background(), "run", "run-1")
	if r.Outcome != "accepted" || r.Activity.Runs[0].State != "running" {
		t.Fatalf("active result=%+v", r)
	}
}
