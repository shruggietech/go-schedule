package operations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const callTimeout = 3 * time.Second
const activityLimit = 200

type Service struct {
	backend Backend
	now     func() time.Time
}

func NewService(backend Backend) *Service { return &Service{backend: backend, now: time.Now} }

func (s *Service) ScheduleWindow(ctx context.Context, days int) OperationResult {
	if days != 1 && days != 7 && days != 30 {
		return OperationResult{Action: "load_schedule", Outcome: "rejected", Field: "days", Message: "Choose a 1-day, 7-day, or 30-day schedule window."}
	}
	now := s.now()
	from, to := now.Add(-24*time.Hour), now.Add(time.Duration(days)*24*time.Hour)
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	response, err := s.backend.GetCalendar(c, from, to)
	if err != nil {
		return failure("load_schedule", err)
	}
	snapshot := buildSchedule(response, s.now())
	return OperationResult{Action: "load_schedule", Outcome: "accepted", Message: "Schedule is up to date.", Schedule: &snapshot}
}

func buildSchedule(response server.CalendarResponse, loadedAt time.Time) ScheduleSnapshot {
	items := make([]ScheduleOccurrence, 0, len(response.Occurrences))
	identities := map[string]int{}
	for _, value := range response.Occurrences {
		kind, state := "prediction", "upcoming"
		base := "prediction:" + value.TaskID + ":" + value.Time.UTC().Format(time.RFC3339Nano)
		if value.Kind == "past" || value.RunID != "" {
			kind, state = "recorded", runState(value.Outcome, false, false)
			base = "run:" + value.RunID
		}
		ordinal := identities[base]
		identities[base] = ordinal + 1
		items = append(items, ScheduleOccurrence{ID: fmt.Sprintf("%s:%d", base, ordinal), TaskID: value.TaskID, TaskName: fallback(value.TaskName, "Unnamed task"), RunID: value.RunID, Time: value.Time.Format(time.RFC3339Nano), Kind: kind, State: state, Outcome: string(value.Outcome)})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Time < items[j].Time })
	return ScheduleSnapshot{From: response.From.Format(time.RFC3339Nano), To: response.To.Format(time.RFC3339Nano), Occurrences: items, LoadedAt: loadedAt.UTC().Format(time.RFC3339Nano)}
}

func (s *Service) ActivityWorkspace(ctx context.Context) OperationResult {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	activeRuns, err := s.backend.ListActiveRuns(c)
	if err != nil {
		return failure("load_activity", err)
	}
	runs, err := s.backend.ListRuns(c, "", activityLimit)
	if err != nil {
		return failure("load_activity", err)
	}
	logs, err := s.backend.ListLogs(c, "", activityLimit)
	if err != nil {
		return failure("load_activity", err)
	}
	alerts, err := s.backend.ListAlertsLimited(c, false, activityLimit)
	if err != nil {
		return failure("load_activity", err)
	}
	workspace := buildActivity(mergeRuns(runs, activeRuns), logs, alerts, s.now())
	return OperationResult{Action: "load_activity", Outcome: "accepted", Message: "Activity is up to date.", Activity: &workspace}
}

func mergeRuns(persisted, active []domain.Run) []domain.Run {
	merged := append([]domain.Run(nil), persisted...)
	seen := make(map[string]bool, len(persisted))
	for _, run := range persisted {
		seen[run.ID] = true
	}
	for _, run := range active {
		if !seen[run.ID] {
			merged = append(merged, run)
		}
	}
	return merged
}

func buildActivity(runs []domain.Run, logs server.LogsResponse, alerts []domain.Alert, loadedAt time.Time) ActivityWorkspace {
	workspace := ActivityWorkspace{Runs: make([]RunRecord, 0, len(runs)), Logs: make([]LogRecord, 0, len(logs.Logs)), Alerts: make([]AlertRecord, 0, len(alerts)), LogPath: logs.LogPath, LoadedAt: loadedAt.UTC().Format(time.RFC3339Nano)}
	for _, value := range runs {
		workspace.Runs = append(workspace.Runs, RunRecord{ID: value.ID, TaskID: value.TaskID, ScheduledFor: timestamp(value.ScheduledFor), StartedAt: optionalTimestamp(value.StartedAt), EndedAt: optionalTimestamp(value.EndedAt), State: runState(value.Outcome, value.StartedAt != nil, value.EndedAt != nil), Outcome: string(value.Outcome), ExitCode: value.ExitCode, Output: value.Output, OutputTruncated: value.OutputTruncated, Trigger: string(value.Trigger), SourceTaskID: value.SourceTaskID, SourceRunID: value.SourceRunID, SourceTriggerID: value.SourceTriggerID, SourceWatcherID: value.SourceWatcherID})
	}
	for _, value := range logs.Logs {
		workspace.Logs = append(workspace.Logs, LogRecord{ID: value.ID, Time: timestamp(value.Time), Severity: severity(value.Severity), Source: fallback(value.Source, "daemon"), Message: fallback(value.Message, "No message"), TaskID: value.TaskID, RunID: value.RunID, Detail: attributeDetail(value.Attrs)})
	}
	for _, value := range alerts {
		workspace.Alerts = append(workspace.Alerts, AlertRecord{ID: value.ID, TaskID: value.TaskID, RunID: value.RunID, Time: timestamp(value.CreatedAt), Severity: severity(value.Severity), Kind: string(value.Kind), Message: fallback(value.Message, "No message"), Acknowledged: value.Acknowledged})
	}
	return workspace
}

func (s *Service) AcknowledgeAlert(ctx context.Context, id string) OperationResult {
	return s.AcknowledgeAlerts(ctx, []string{id})
}

func (s *Service) AcknowledgeAlerts(ctx context.Context, ids []string) OperationResult {
	seen := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		c, cancel := context.WithTimeout(ctx, callTimeout)
		err := s.backend.AckAlert(c, id)
		cancel()
		if err != nil {
			return failure("acknowledge_alerts", err)
		}
	}
	loaded := s.ActivityWorkspace(ctx)
	if loaded.Outcome != "accepted" {
		return OperationResult{Action: "acknowledge_alerts", Outcome: "accepted", Message: "Visible alerts were acknowledged. Refresh Activity to see the latest state."}
	}
	return OperationResult{Action: "acknowledge_alerts", Outcome: "accepted", Message: "Visible alerts were acknowledged.", Activity: loaded.Activity}
}

func runState(outcome domain.RunOutcome, started, ended bool) string {
	switch outcome {
	case domain.OutcomeSuccess, domain.OutcomeFailure, domain.OutcomeSkipped, domain.OutcomeCaughtUp, domain.OutcomeQueued:
		return string(outcome)
	}
	if started && !ended {
		return "running"
	}
	return "unavailable"
}
func severity(value domain.AlertSeverity) string {
	if value == "" {
		return "info"
	}
	return string(value)
}
func timestamp(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }
func optionalTimestamp(value *time.Time) string {
	if value == nil {
		return ""
	}
	return timestamp(*value)
}
func fallback(value, fallbackValue string) string {
	if strings.TrimSpace(value) == "" {
		return fallbackValue
	}
	return value
}
func attributeDetail(attrs map[string]any) string {
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		value, err := json.Marshal(attrs[key])
		if err != nil {
			value = []byte(fmt.Sprint(attrs[key]))
		}
		lines = append(lines, fmt.Sprintf("%s: %s", key, value))
	}
	return strings.Join(lines, "\n")
}
func failure(action string, err error) OperationResult {
	var uncertain *client.MutationUncertainError
	if errors.As(err, &uncertain) {
		return OperationResult{Action: action, Outcome: "uncertain", Message: "The remote request may have completed. Refresh the selected scheduler before deciding whether to try again."}
	}
	var status *client.StatusError
	if errors.As(err, &status) && status.Code == server.CodeValidation {
		return OperationResult{Action: action, Outcome: "rejected", Message: "The request was rejected. Review the selected filters and try again."}
	}
	return OperationResult{Action: action, Outcome: "unavailable", Message: "The local scheduler service is unavailable. Try again."}
}
