package mcpobserve

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	healthURI    = "goschedule://daemon/health"
	tasksURI     = "goschedule://tasks/active"
	schedulesURI = "goschedule://schedules/upcoming"
	alertsURI    = "goschedule://alerts/recent"
	runsURI      = "goschedule://runs/recent"
)

type readClient interface {
	Health(context.Context) (server.HealthResponse, error)
	ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error)
	ListRunsPage(context.Context, string, int, int, int) ([]domain.Run, error)
	ListAlertsPage(context.Context, bool, int, int, int) ([]domain.Alert, error)
}

type adapter struct {
	client  readClient
	now     func() time.Time
	timeout time.Duration
}

func (a *adapter) read(kind, baseURI string) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		offset, err := resourceOffset(req.Params.URI, baseURI)
		if err != nil {
			return nil, err
		}
		readCtx, cancel := context.WithTimeout(ctx, a.timeout)
		defer cancel()
		envelope, err := a.envelope(readCtx, kind, baseURI, offset)
		if err != nil {
			return nil, safeError(err)
		}
		body, err := json.Marshal(envelope)
		if err != nil {
			return nil, errors.New("internal_error: resource could not be encoded")
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: req.Params.URI, MIMEType: "application/json", Text: string(body)}}}, nil
	}
}

func resourceOffset(uri, base string) (int, error) {
	if uri == base {
		return 0, nil
	}
	prefix := base + "/page/"
	if !strings.HasPrefix(uri, prefix) || strings.Contains(strings.TrimPrefix(uri, prefix), "/") {
		return 0, errors.New("invalid_resource: resource URI is not approved")
	}
	offset, err := decodeCursor(strings.TrimPrefix(uri, prefix))
	if err != nil || offset == 0 || offset%PageLimit != 0 {
		return 0, errInvalidCursor
	}
	return offset, nil
}

func (a *adapter) envelope(ctx context.Context, kind, baseURI string, offset int) (Envelope, error) {
	result := Envelope{SchemaVersion: SchemaVersion, Permission: Permission, GeneratedAt: timestamp(a.now()), Trust: TrustNotice}
	switch kind {
	case "health":
		if offset != 0 {
			return Envelope{}, errInvalidCursor
		}
		health, err := a.client.Health(ctx)
		if err != nil {
			return Envelope{}, err
		}
		status, _ := boundedText(health.Status, TextLimit)
		version, _ := boundedText(health.Version, TextLimit)
		result.Data = Health{Status: status, Version: version}
		return result, nil
	case "tasks", "schedules":
		details, err := a.client.ListTaskDetails(ctx, "", string(domain.TaskActive))
		if err != nil {
			return Envelope{}, err
		}
		sort.Slice(details, func(i, j int) bool { return details[i].Task.ID < details[j].Task.ID })
		if kind == "tasks" {
			items := make([]TaskSummary, 0, len(details))
			for _, detail := range details {
				items = append(items, safeTask(detail))
			}
			page, metadata, err := paginate(items, offset, baseURI)
			if err != nil {
				return Envelope{}, err
			}
			result.UntrustedFields = []string{"data[].name", "data[].readiness_reason", "data[].schedule_summary", "data[].policy_summary"}
			result.Page, result.Data = metadata, page
			return result, nil
		}
		items := make([]ScheduleSummary, 0, len(details))
		for _, detail := range details {
			if detail.Schedule != nil {
				items = append(items, safeSchedule(detail))
			}
		}
		page, metadata, err := paginate(items, offset, baseURI)
		if err != nil {
			return Envelope{}, err
		}
		result.UntrustedFields = []string{"data[].task_name", "data[].summary", "data[].policy_summary"}
		result.Page, result.Data = metadata, page
		return result, nil
	case "alerts":
		alerts, err := a.client.ListAlertsPage(ctx, false, offset, PageLimit+1, TextLimit)
		if err != nil {
			return Envelope{}, err
		}
		sort.Slice(alerts, func(i, j int) bool {
			if alerts[i].CreatedAt.Equal(alerts[j].CreatedAt) {
				return alerts[i].ID > alerts[j].ID
			}
			return alerts[i].CreatedAt.After(alerts[j].CreatedAt)
		})
		items := make([]AlertSummary, 0, len(alerts))
		for _, alert := range alerts {
			items = append(items, safeAlert(alert))
		}
		page, metadata := boundedPage(items, offset, baseURI)
		result.UntrustedFields = []string{"data[].message"}
		result.Page, result.Data = metadata, page
		return result, nil
	case "runs":
		runs, err := a.client.ListRunsPage(ctx, "", offset, PageLimit+1, OutputLimit)
		if err != nil {
			return Envelope{}, err
		}
		sort.Slice(runs, func(i, j int) bool {
			if runs[i].ScheduledFor.Equal(runs[j].ScheduledFor) {
				return runs[i].ID > runs[j].ID
			}
			return runs[i].ScheduledFor.After(runs[j].ScheduledFor)
		})
		items := make([]RunSummary, 0, len(runs))
		for _, run := range runs {
			items = append(items, safeRun(run))
		}
		page, metadata := boundedPage(items, offset, baseURI)
		result.UntrustedFields = []string{"data[].output_excerpt"}
		result.Page, result.Data = metadata, page
		return result, nil
	default:
		return Envelope{}, errors.New("invalid_resource: resource URI is not approved")
	}
}

func boundedPage[T any](items []T, offset int, baseURI string) ([]T, *Page) {
	hasMore := len(items) > PageLimit
	if hasMore {
		items = items[:PageLimit]
	}
	if items == nil {
		items = []T{}
	}
	metadata := &Page{Count: len(items), Limit: PageLimit}
	if hasMore {
		metadata.NextURI = baseURI + "/page/" + encodeCursor(offset+PageLimit)
	}
	return items, metadata
}

func paginate[T any](items []T, offset int, baseURI string) ([]T, *Page, error) {
	if offset < 0 || offset > len(items) || (offset > 0 && offset%PageLimit != 0) {
		return nil, nil, errInvalidCursor
	}
	end := offset + PageLimit
	if end > len(items) {
		end = len(items)
	}
	page := items[offset:end]
	if page == nil {
		page = []T{}
	}
	metadata := &Page{Count: len(page), Limit: PageLimit}
	if end < len(items) {
		metadata.NextURI = baseURI + "/page/" + encodeCursor(end)
	}
	return page, metadata, nil
}

func safeTask(detail server.TaskResponse) TaskSummary {
	name, nameTruncated := boundedText(detail.Task.Name, TextLimit)
	reason, reasonTruncated := boundedText(detail.Readiness.Reason, TextLimit)
	policy, policyTruncated := boundedText(detail.PolicySummary, TextLimit)
	summary := ""
	summaryTruncated := false
	if detail.Schedule != nil {
		summary, summaryTruncated = boundedText(detail.Schedule.HumanSummary, TextLimit)
	}
	return TaskSummary{ID: detail.Task.ID, Name: name, GroupID: detail.Task.GroupID, Enabled: detail.Task.Enabled, State: string(detail.Task.State), Timezone: detail.Task.Timezone, Readiness: string(detail.Readiness.Status), ReadinessReason: reason, ScheduleSummary: summary, PolicySummary: policy, NextRuns: safeTimes(detail.NextRuns), UpdatedAt: timestamp(detail.Task.UpdatedAt), TruncatedFields: truncatedFields(map[string]bool{"name": nameTruncated, "readiness_reason": reasonTruncated, "schedule_summary": summaryTruncated, "policy_summary": policyTruncated})}
}

func safeSchedule(detail server.TaskResponse) ScheduleSummary {
	name, nameTruncated := boundedText(detail.Task.Name, TextLimit)
	summary, summaryTruncated := boundedText(detail.Schedule.HumanSummary, TextLimit)
	policy, policyTruncated := boundedText(detail.PolicySummary, TextLimit)
	return ScheduleSummary{TaskID: detail.Task.ID, TaskName: name, GroupID: detail.Task.GroupID, Enabled: detail.Task.Enabled, Readiness: string(detail.Readiness.Status), Timezone: detail.Task.Timezone, Summary: summary, PolicySummary: policy, NextRuns: safeTimes(detail.NextRuns), TruncatedFields: truncatedFields(map[string]bool{"task_name": nameTruncated, "summary": summaryTruncated, "policy_summary": policyTruncated})}
}

func truncatedFields(fields map[string]bool) []string {
	result := make([]string, 0, len(fields))
	for name, truncated := range fields {
		if truncated {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func safeRun(run domain.Run) RunSummary {
	output, truncated := boundedText(run.Output, OutputLimit)
	return RunSummary{ID: run.ID, TaskID: run.TaskID, ScheduledFor: timestamp(run.ScheduledFor), StartedAt: optionalTimestamp(run.StartedAt), EndedAt: optionalTimestamp(run.EndedAt), Outcome: string(run.Outcome), ExitCode: run.ExitCode, TriggerKind: string(run.Trigger), OutputExcerpt: output, OutputTruncated: run.OutputTruncated || truncated}
}

func safeAlert(alert domain.Alert) AlertSummary {
	message, truncated := boundedText(alert.Message, TextLimit)
	return AlertSummary{ID: alert.ID, TaskID: alert.TaskID, RunID: alert.RunID, Severity: string(alert.Severity), Kind: string(alert.Kind), Message: message, MessageTruncated: alert.MessageTruncated || truncated, CreatedAt: timestamp(alert.CreatedAt), Acknowledged: alert.Acknowledged}
}

func safeTimes(values []time.Time) []string {
	result := make([]string, 0, len(values))
	for i, value := range values {
		if i == 5 {
			break
		}
		result = append(result, timestamp(value))
	}
	return result
}

func safeError(err error) error {
	if errors.Is(err, context.Canceled) {
		return errors.New("canceled: request was canceled")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("timeout: daemon read exceeded its deadline")
	}
	var connectionError *client.ConnectionError
	if errors.As(err, &connectionError) {
		switch connectionError.Kind {
		case client.ConnectionAccessDenied:
			return errors.New("access_denied: daemon IPC denied the current identity")
		case client.ConnectionUnavailable:
			return errors.New("daemon_unavailable: start the local goschedd service and verify access")
		case client.ConnectionTimeout:
			return errors.New("timeout: daemon read exceeded its deadline")
		default:
			return errors.New("daemon_unavailable: the local daemon could not be reached")
		}
	}
	var statusError *client.StatusError
	if errors.As(err, &statusError) {
		return errors.New("incompatible_response: the daemon could not provide this observation")
	}
	return errors.New("internal_error: observation failed")
}
