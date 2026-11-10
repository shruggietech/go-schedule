package mcpobserve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestActiveRunsDiscoveryIsBoundedAndOmitsExecutionContent(t *testing.T) {
	reader := &fakeReader{}
	for i := 0; i < PageLimit+1; i++ {
		reader.active = append(reader.active, domain.Run{ID: fmt.Sprintf("run-%03d", i), TaskID: "task", ScheduledFor: time.Unix(int64(i), 0), Output: "SECRET_OUTPUT", SourceTriggerID: "SECRET_TRIGGER", CancellationStatus: "failed", CancellationError: "SECRET_PROCESS_ERROR"})
	}
	a := &adapter{client: reader, now: time.Now, timeout: time.Second}
	first, err := a.envelope(context.Background(), "active_runs", "goschedule://runs/active", 0)
	if err != nil {
		t.Fatal(err)
	}
	if first.Page.Count != PageLimit || first.Page.NextURI == "" {
		t.Fatalf("page=%+v", first.Page)
	}
	body, _ := json.Marshal(first)
	if strings.Contains(string(body), "SECRET_") || strings.Contains(string(body), "output_") {
		t.Fatalf("execution content leaked: %s", body)
	}
	if !strings.Contains(string(body), "cancellation_status") || !strings.Contains(string(body), "termination could not be confirmed") {
		t.Fatalf("missing cleanup guidance: %s", body)
	}
	offset, err := resourceOffset(first.Page.NextURI, "goschedule://runs/active")
	if err != nil {
		t.Fatal(err)
	}
	second, err := a.envelope(context.Background(), "active_runs", "goschedule://runs/active", offset)
	if err != nil || second.Page.Count != 1 || second.Page.NextURI != "" {
		t.Fatalf("continuation=%+v, %v", second, err)
	}
}

func TestActiveRunsUnavailableDoesNotSubstitutePersistedHistory(t *testing.T) {
	reader := &fakeReader{activeErr: errors.New("SECRET_TRANSPORT"), runs: []domain.Run{{ID: "old-finished"}}}
	a := &adapter{client: reader, now: time.Now, timeout: time.Second}
	_, err := a.read("active_runs", "goschedule://runs/active")(context.Background(), &mcp.ReadResourceRequest{Params: &mcp.ReadResourceParams{URI: "goschedule://runs/active"}})
	if err == nil || !strings.HasPrefix(err.Error(), "internal_error:") || strings.Contains(err.Error(), "SECRET_") || reader.runsLimit != 0 {
		t.Fatalf("error=%v history limit=%d", err, reader.runsLimit)
	}
}
