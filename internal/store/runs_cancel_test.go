package store

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestRunResultMetadataRetainsCancellationAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "current.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	task := chainTask(t, st, "owned")
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	runs := []domain.Run{{ID: "canceled", TaskID: task.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: domain.OutcomeCanceled, Output: strings.Repeat("private-output", 100000), Trigger: domain.TriggerManual}, {ID: "cleanup-failed", TaskID: task.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: domain.OutcomeCancellationFailed, Output: "retained", Trigger: domain.TriggerManual}}
	for i := range runs {
		if err := st.CreateRun(&runs[i]); err != nil {
			t.Fatal(err)
		}
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	for _, want := range runs {
		got, err := st.GetRunResultContext(context.Background(), want.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != want.ID || got.TaskID != task.ID || got.Outcome != want.Outcome || got.StartedAt == nil || !got.StartedAt.Equal(now) || got.EndedAt == nil || !got.EndedAt.Equal(now) || got.Output != "" {
			t.Fatalf("metadata=%+v", got)
		}
	}
	if _, err := st.GetRunResultContext(context.Background(), "absent"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing=%v", err)
	}
}
func TestRunResultMetadataCancelsWhileSQLPoolIsOccupied(t *testing.T) {
	st := openMem(t)
	conn, err := st.db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	before := st.db.Stats().WaitCount
	done := make(chan error, 1)
	go func() { _, err := st.GetRunResultContext(ctx, "not-loaded"); done <- err }()
	watch := time.NewTimer(2 * time.Second)
	defer watch.Stop()
	for st.db.Stats().WaitCount == before {
		select {
		case err := <-done:
			t.Fatalf("returned before waiting: %v", err)
		case <-watch.C:
			t.Fatal("did not reach SQL pool wait")
		default:
			runtime.Gosched()
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error=%v", err)
		}
	case <-watch.C:
		t.Fatal("request cancellation did not release SQL pool waiter")
	}
}
