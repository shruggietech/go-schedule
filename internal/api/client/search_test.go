package client

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestSearchDecodesAdditiveContract(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	task := &domain.Task{Name: "visible task", Command: "secret command", Timezone: "UTC", OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne, MissingDatePolicy: domain.MissingDateSkip, TimeBasis: domain.TimeBasisWallClock, DSTGapPolicy: domain.DSTGapNextValid, DSTOverlapPolicy: domain.DSTOverlapFirst, State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	handler := server.New(st, nil, nil, nil, "", config.NewLogger(config.Default(), summaryDiscard{})).Handler()
	result, err := NewInProcess(handler).Search(context.Background(), "visible", []domain.SearchKind{domain.SearchKindTask}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != domain.DaemonSearchSchema || result.Query != "visible" || result.ObservedAt.IsZero() || len(result.Results) != 1 {
		t.Fatalf("result=%+v", result)
	}
	if result.Results[0].ObjectID != task.ID || result.Results[0].Kind != domain.SearchKindTask {
		t.Fatalf("match=%+v", result.Results[0])
	}
}
