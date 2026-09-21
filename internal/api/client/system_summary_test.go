package client

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

type summaryDiscard struct{}

func (summaryDiscard) Write(value []byte) (int, error) { return len(value), nil }

func TestSystemSummaryDecodesAdditiveContract(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	handler := server.New(st, nil, nil, nil, "", config.NewLogger(config.Default(), summaryDiscard{})).Handler()
	summary, err := NewInProcess(handler).SystemSummary(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if summary.Schema != domain.SystemSummarySchema || summary.ActiveTaskCount != 0 || summary.ObservedAt.IsZero() {
		t.Fatalf("summary=%+v", summary)
	}
}
