package mcpobserve

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/schedule"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestPersistedScheduleDescriptionsReachAuthenticatedResources(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	api := server.New(st, nil, nil, nil, "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	base := client.NewInProcess(api.Handler())
	ctx := context.Background()
	credential, err := base.CreateMCPSession(ctx, "schedule-description-regression", domain.CapabilityObserve)
	if err != nil {
		t.Fatal(err)
	}
	observed := base.WithMCPSession(credential.Credential)
	now := time.Now().UTC()
	future := now.Add(24 * time.Hour)
	secret := "EXECUTION_INPUT_CANARY"
	want := map[string]string{}
	for _, tc := range []struct {
		name   string
		sch    *domain.Schedule
		policy domain.MissingDatePolicy
	}{
		{name: "recurring", sch: &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=DAILY", Anchor: &now, HumanSummary: "Every day"}},
		{name: "one-off", sch: &domain.Schedule{Kind: domain.ScheduleOneOff, RunAt: &future, HumanSummary: "Once tomorrow"}},
		{name: "startup", sch: &domain.Schedule{Kind: domain.ScheduleEvent, TriggerID: domain.StartupEventID, HumanSummary: "At scheduler startup"}},
		{name: "missing-date", sch: &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MONTHLY;BYMONTHDAY=31", Anchor: &now, HumanSummary: "The 31st of every month"}, policy: domain.MissingDateLastValid},
		{name: "long-unicode", sch: &domain.Schedule{Kind: domain.ScheduleEvent, TriggerID: domain.StartupEventID, HumanSummary: strings.Repeat("界", TextLimit)}},
		{name: "long-ascii", sch: &domain.Schedule{Kind: domain.ScheduleEvent, TriggerID: domain.StartupEventID, HumanSummary: strings.Repeat("a", TextLimit+1)}},
		{name: "manual"},
	} {
		task := domain.Task{ID: tc.name, Name: tc.name, Command: secret, Args: []string{secret}, Env: map[string]string{"TOKEN": secret}, Stdin: secret, Timezone: "UTC", State: domain.TaskActive, MissingDatePolicy: tc.policy}
		if tc.sch != nil {
			if err := st.CreateSchedule(tc.sch); err != nil {
				t.Fatal(err)
			}
			task.ScheduleID = tc.sch.ID
			want[task.ID] = schedule.Describe(*tc.sch, tc.policy)
		}
		if err := st.CreateTask(&task); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := observed.ListTaskObservations(ctx, "", string(domain.TaskActive), false, 0, 100, TextLimit)
	if err != nil || len(rows) != 7 {
		t.Fatalf("observations=%+v err=%v", rows, err)
	}
	for _, row := range rows {
		expected := want[row.ID]
		if row.ID == "long-unicode" {
			expected = strings.Repeat("界", TextLimit/3)
		}
		if row.ID == "long-ascii" {
			expected = strings.Repeat("a", TextLimit)
		}
		if row.ScheduleSummary != expected || !utf8.ValidString(row.ScheduleSummary) || len(row.ScheduleSummary) > TextLimit || row.ScheduleSummaryTruncated != strings.HasPrefix(row.ID, "long-") {
			t.Errorf("%s summary=%q truncated=%v want=%q", row.ID, row.ScheduleSummary, row.ScheduleSummaryTruncated, expected)
		}
	}
	a := &adapter{client: observed, now: time.Now, timeout: 30 * time.Second}
	for _, kind := range []string{"tasks", "schedules"} {
		uri := tasksURI
		if kind == "schedules" {
			uri = schedulesURI
		}
		result, err := a.read(kind, uri)(ctx, &mcp.ReadResourceRequest{Params: &mcp.ReadResourceParams{URI: uri}})
		if err != nil {
			t.Fatal(err)
		}
		body := result.Contents[0].Text
		if strings.Contains(body, secret) {
			t.Fatalf("%s exposed execution inputs", kind)
		}
		var envelope struct {
			Data []struct {
				ID              string
				TaskID          string `json:"task_id"`
				Summary         string
				ScheduleSummary string   `json:"schedule_summary"`
				TruncatedFields []string `json:"truncated_fields"`
			}
		}
		if err := json.Unmarshal([]byte(body), &envelope); err != nil {
			t.Fatal(err)
		}
		expectedCount := 7
		if kind == "schedules" {
			expectedCount = 6
		}
		if len(envelope.Data) != expectedCount {
			t.Fatalf("%s count=%d", kind, len(envelope.Data))
		}
		for _, row := range envelope.Data {
			id, summary, field := row.ID, row.ScheduleSummary, "schedule_summary"
			if kind == "schedules" {
				id, summary, field = row.TaskID, row.Summary, "summary"
			}
			expected := want[id]
			if id == "long-unicode" {
				expected = strings.Repeat("界", TextLimit/3)
			}
			if id == "long-ascii" {
				expected = strings.Repeat("a", TextLimit)
			}
			if summary != expected || slices.Contains(row.TruncatedFields, field) != strings.HasPrefix(id, "long-") {
				t.Errorf("%s/%s summary=%q truncated=%v want=%q", kind, id, summary, row.TruncatedFields, expected)
			}
		}
	}
}
