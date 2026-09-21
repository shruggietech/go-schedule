package store

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSystemSummaryFactsBoundsWindowsAndRedactsRepresentatives(t *testing.T) {
	st := openMem(t)
	now := time.Date(2026, 9, 21, 15, 0, 0, 0, time.UTC)
	schedule := &domain.Schedule{Kind: domain.ScheduleEvent, TriggerID: domain.StartupEventID, HumanSummary: "startup"}
	if err := st.CreateSchedule(schedule); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{ID: "task-1", Name: "Archive", Command: "secret-command", Args: []string{"secret-arg"}, Env: map[string]string{"TOKEN": "secret"}, Enabled: true, Timezone: "UTC", ScheduleID: schedule.ID, State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	recentEnd := now.Add(-time.Hour)
	recent := &domain.Run{ID: "run-recent", TaskID: task.ID, ScheduledFor: recentEnd, EndedAt: &recentEnd, Outcome: domain.OutcomeFailure, Output: "secret-output", Trigger: domain.TriggerSchedule}
	oldEnd := now.Add(-25 * time.Hour)
	old := &domain.Run{ID: "run-old", TaskID: task.ID, ScheduledFor: oldEnd, EndedAt: &oldEnd, Outcome: domain.OutcomeFailure, Trigger: domain.TriggerSchedule}
	if err := st.CreateRun(recent); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateRun(old); err != nil {
		t.Fatal(err)
	}
	alert := &domain.Alert{ID: "alert-1", TaskID: task.ID, RunID: recent.ID, Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "secret-alert", CreatedAt: recentEnd}
	if err := st.CreateAlert(alert); err != nil {
		t.Fatal(err)
	}
	_, err := st.db.Exec(`INSERT INTO notification_deliveries(id,channel_name,destination_summary,endpoint,authorization,event_kind,task_id,run_id,task_name,group_id,group_name,payload,state,attempts,next_attempt_at,created_at,last_status,last_error,condition_kind,condition_summary) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, "delivery-1", "Operations", "secret-destination", "https://secret.invalid", "secret-auth", "run.completed", task.ID, recent.ID, task.Name, "", "", []byte(`{"secret":true}`), string(domain.NotificationDeliveryFailed), 3, fmtTime(now), fmtTime(recentEnd), 500, "secret-error", "failure", "secret-condition")
	if err != nil {
		t.Fatal(err)
	}

	summary, err := st.SystemSummaryFacts(now)
	if err != nil {
		t.Fatal(err)
	}
	if summary.ActiveTaskCount != 1 || summary.RecentFailureCount != 1 || summary.RecentFailure == nil || summary.RecentFailure.RunID != recent.ID {
		t.Fatalf("failure facts = %+v", summary)
	}
	if summary.UnacknowledgedAlertCount != 1 || summary.UnacknowledgedAlert == nil || summary.UnacknowledgedAlert.AlertID != alert.ID {
		t.Fatalf("alert facts = %+v", summary)
	}
	if summary.NotificationProblemCount != 1 || summary.NotificationProblem == nil || summary.NotificationProblem.DeliveryID != "delivery-1" {
		t.Fatalf("notification facts = %+v", summary)
	}
}
