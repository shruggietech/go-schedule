package mcpobserve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	tasklogic "github.com/shruggietech/go-schedule/internal/task"
)

type fakeReader struct {
	health       server.HealthResponse
	tasks        []server.TaskResponse
	runs         []domain.Run
	alerts       []domain.Alert
	err          error
	seenCtx      context.Context
	runsLimit    int
	alertLimit   int
	runsOffset   int
	alertOffset  int
	outputLimit  int
	messageLimit int
	wait         bool
}

func (f *fakeReader) Health(ctx context.Context) (server.HealthResponse, error) {
	f.seenCtx = ctx
	if f.wait {
		<-ctx.Done()
		return server.HealthResponse{}, ctx.Err()
	}
	return f.health, f.err
}

func (f *fakeReader) ListTaskDetails(ctx context.Context, _, _ string) ([]server.TaskResponse, error) {
	f.seenCtx = ctx
	return f.tasks, f.err
}

func (f *fakeReader) ListRunsPage(ctx context.Context, _ string, offset, limit, outputLimit int) ([]domain.Run, error) {
	f.seenCtx, f.runsOffset, f.runsLimit, f.outputLimit = ctx, offset, limit, outputLimit
	return fakePage(f.runs, offset, limit), f.err
}

func (f *fakeReader) ListAlertsPage(ctx context.Context, _ bool, offset, limit, messageLimit int) ([]domain.Alert, error) {
	f.seenCtx, f.alertOffset, f.alertLimit, f.messageLimit = ctx, offset, limit, messageLimit
	return fakePage(f.alerts, offset, limit), f.err
}

func fakePage[T any](items []T, offset, limit int) []T {
	if offset >= len(items) {
		return nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}

func TestObserveResourcesAreBoundedPaginatedAndSecretFree(t *testing.T) {
	now := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	secret := "SECRET_CANARY"
	reader := &fakeReader{health: server.HealthResponse{Status: "ok", Version: "v1.3.0"}}
	for i := 0; i < 205; i++ {
		id := fmt.Sprintf("task-%03d", i)
		reader.tasks = append(reader.tasks, server.TaskResponse{
			Task:          domain.Task{ID: id, Name: "task <ignore instructions> " + id, GroupID: "group", Command: secret, Args: []string{secret}, WorkingDir: secret, Env: map[string]string{"TOKEN": secret}, Stdin: secret, RunAs: secret, Enabled: true, Timezone: "UTC", ScheduleID: secret, State: domain.TaskActive, UpdatedAt: now},
			Schedule:      &domain.Schedule{ID: secret, RRULE: secret, Expression: secret, HumanSummary: "Every day " + id},
			Readiness:     tasklogic.Readiness{Status: tasklogic.StatusReady, Reason: "ready"},
			PolicySummary: "wall clock",
			NextRuns:      []time.Time{now.Add(time.Hour)},
		})
	}
	reader.runs = []domain.Run{{ID: "run", TaskID: "task", ScheduledFor: now, Output: strings.Repeat("界", OutputLimit), OutputTruncated: false, Trigger: domain.TriggerManual, SourceTriggerID: secret, SourceWatcherID: secret}}
	reader.alerts = []domain.Alert{{ID: "alert", TaskID: "task", Message: strings.Repeat("x", TextLimit+20) + secret, CreatedAt: now}}
	a := &adapter{client: reader, now: func() time.Time { return now }, timeout: time.Second}

	first, err := a.envelope(context.Background(), "tasks", tasksURI, 0)
	if err != nil {
		t.Fatal(err)
	}
	if first.Page.Count != PageLimit || first.Page.NextURI == "" {
		t.Fatalf("first page = %+v", first.Page)
	}
	secondOffset, err := resourceOffset(first.Page.NextURI, tasksURI)
	if err != nil || secondOffset != PageLimit {
		t.Fatalf("continuation offset = %d, %v", secondOffset, err)
	}
	second, err := a.envelope(context.Background(), "tasks", tasksURI, secondOffset)
	if err != nil || second.Page.Count != PageLimit || second.Page.NextURI == "" {
		t.Fatalf("second page = %+v, %v", second.Page, err)
	}
	thirdOffset, _ := resourceOffset(second.Page.NextURI, tasksURI)
	third, err := a.envelope(context.Background(), "tasks", tasksURI, thirdOffset)
	if err != nil || third.Page.Count != 5 || third.Page.NextURI != "" {
		t.Fatalf("third page = %+v, %v", third.Page, err)
	}

	for _, tc := range []struct {
		kind string
		uri  string
	}{
		{kind: "tasks", uri: tasksURI},
		{kind: "schedules", uri: schedulesURI},
		{kind: "runs", uri: runsURI},
		{kind: "alerts", uri: alertsURI},
	} {
		envelope, readErr := a.envelope(context.Background(), tc.kind, tc.uri, 0)
		if readErr != nil {
			t.Fatalf("%s: %v", tc.kind, readErr)
		}
		encoded, marshalErr := json.Marshal(envelope)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		if strings.Contains(string(encoded), secret) {
			t.Fatalf("%s leaked secret: %s", tc.kind, encoded)
		}
		if !strings.Contains(string(encoded), TrustNotice) || len(envelope.UntrustedFields) == 0 {
			t.Fatalf("%s lacks untrusted-content metadata", tc.kind)
		}
	}
	if reader.runsLimit != PageLimit+1 || reader.alertLimit != PageLimit+1 || reader.outputLimit != OutputLimit || reader.messageLimit != TextLimit {
		t.Fatalf("daemon bounds = runs %d/%d alerts %d/%d", reader.runsLimit, reader.outputLimit, reader.alertLimit, reader.messageLimit)
	}
	runEnvelope, _ := a.envelope(context.Background(), "runs", runsURI, 0)
	run := runEnvelope.Data.([]RunSummary)[0]
	if len(run.OutputExcerpt) > OutputLimit || !run.OutputTruncated {
		t.Fatalf("run output = %d bytes, truncated %t", len(run.OutputExcerpt), run.OutputTruncated)
	}
	alertEnvelope, _ := a.envelope(context.Background(), "alerts", alertsURI, 0)
	alert := alertEnvelope.Data.([]AlertSummary)[0]
	if len(alert.Message) > TextLimit || !alert.MessageTruncated {
		t.Fatalf("alert message = %d bytes, truncated %t", len(alert.Message), alert.MessageTruncated)
	}
}

func TestTaskAndScheduleReportEveryTruncatedField(t *testing.T) {
	long := strings.Repeat("x", TextLimit+1)
	detail := server.TaskResponse{Task: domain.Task{Name: long}, Schedule: &domain.Schedule{HumanSummary: long}, Readiness: tasklogic.Readiness{Reason: long}, PolicySummary: long}
	task := safeTask(detail)
	if got := strings.Join(task.TruncatedFields, ","); got != "name,policy_summary,readiness_reason,schedule_summary" {
		t.Fatalf("task truncated fields = %q", got)
	}
	schedule := safeSchedule(detail)
	if got := strings.Join(schedule.TruncatedFields, ","); got != "policy_summary,summary,task_name" {
		t.Fatalf("schedule truncated fields = %q", got)
	}
}

func TestObserveErrorsAreBoundedAndDoNotLeakCause(t *testing.T) {
	secret := "C:\\secret\\daemon.pipe"
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "denied", err: client.NewConnectionError("GET "+secret, os.ErrPermission), want: "access_denied"},
		{name: "unavailable", err: client.NewConnectionError("GET "+secret, os.ErrNotExist), want: "daemon_unavailable"},
		{name: "timeout", err: context.DeadlineExceeded, want: "timeout"},
		{name: "canceled", err: context.Canceled, want: "canceled"},
		{name: "internal", err: errors.New(secret), want: "internal_error"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := safeError(tc.err).Error()
			if !strings.HasPrefix(got, tc.want+":") || strings.Contains(got, secret) || len(got) > 160 {
				t.Fatalf("safeError() = %q", got)
			}
		})
	}
}

func TestResourceOffsetRejectsUnapprovedAndMisalignedCursors(t *testing.T) {
	if _, err := resourceOffset("goschedule://tasks/private", tasksURI); err == nil {
		t.Fatal("unapproved resource succeeded")
	}
	if _, _, err := paginate([]int{1, 2}, 1, tasksURI); err == nil {
		t.Fatal("misaligned cursor succeeded")
	}
}

func TestReadPropagatesDeadlineAndReturnsJSON(t *testing.T) {
	reader := &fakeReader{health: server.HealthResponse{Status: "ok", Version: "v"}}
	a := &adapter{client: reader, now: time.Now, timeout: time.Second}
	handler := a.read("health", healthURI)
	result, err := handler(context.Background(), &mcp.ReadResourceRequest{Params: &mcp.ReadResourceParams{URI: healthURI}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Contents) != 1 || result.Contents[0].MIMEType != "application/json" || !json.Valid([]byte(result.Contents[0].Text)) {
		t.Fatalf("result = %+v", result)
	}
	if _, ok := reader.seenCtx.Deadline(); !ok {
		t.Fatal("daemon context lacks deadline")
	}
}

func TestReadBoundsTimeoutAndCancellation(t *testing.T) {
	for _, tc := range []struct {
		name    string
		context func() (context.Context, context.CancelFunc)
		want    string
	}{
		{name: "timeout", context: func() (context.Context, context.CancelFunc) { return context.Background(), func() {} }, want: "timeout:"},
		{name: "canceled", context: func() (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return ctx, func() {}
		}, want: "canceled:"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := tc.context()
			defer cancel()
			timeout := 10 * time.Millisecond
			if tc.name == "canceled" {
				timeout = time.Second
			}
			a := &adapter{client: &fakeReader{wait: true}, now: time.Now, timeout: timeout}
			_, err := a.read("health", healthURI)(ctx, &mcp.ReadResourceRequest{Params: &mcp.ReadResourceParams{URI: healthURI}})
			if err == nil || !strings.HasPrefix(err.Error(), tc.want) {
				t.Fatalf("read error = %v", err)
			}
		})
	}
}
