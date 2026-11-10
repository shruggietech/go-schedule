package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestCancelRunRetainsRequestedAndFinishedAcknowledgement(t *testing.T) {
	for _, status := range []string{"requested", "already_finished"} {
		calls := 0
		c := NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if r.Method != "POST" || r.URL.Path != "/v1/runs/owned-id/cancel" || r.Header.Get(server.ExpectedDaemonHeader) != "owned-daemon" {
				t.Errorf("request=%s %s identity=%s", r.Method, r.URL.Path, r.Header.Get(server.ExpectedDaemonHeader))
			}
			if status == "requested" {
				w.WriteHeader(http.StatusAccepted)
			}
			outcome := domain.RunOutcome("")
			if status == "already_finished" {
				outcome = domain.OutcomeSuccess
			}
			json.NewEncoder(w).Encode(domain.RunCancellation{RunID: "owned-id", TaskID: "task", Status: status, Outcome: outcome})
		}))
		ack, err := c.CancelRun(ExpectDaemon(context.Background(), "owned-daemon"), "owned-id")
		if err != nil || ack.RunID != "owned-id" || ack.Status != status || calls != 1 {
			t.Fatalf("ack=%+v err=%v calls=%d", ack, err, calls)
		}
	}
}
func TestCancelRunPreservesAPIFailureWithoutReplay(t *testing.T) {
	calls := 0
	c := NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":{"code":"conflict","message":"run is not active"}}`))
	}))
	_, err := c.CancelRun(context.Background(), "queued")
	var status *StatusError
	if !errors.As(err, &status) || status.Code != server.CodeConflict || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}

func TestCancelRunRejectsAnUnattributedOrIncompleteAcknowledgement(t *testing.T) {
	for _, ack := range []domain.RunCancellation{{RunID: "different", Status: "requested"}, {RunID: "owned", Status: "unknown"}, {RunID: "owned", Status: "already_finished"}} {
		c := NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(ack) }))
		if _, err := c.CancelRun(context.Background(), "owned"); err == nil {
			t.Fatalf("accepted invalid ack %+v", ack)
		}
	}
}
