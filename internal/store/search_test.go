package store

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSearchFactsBoundsKindsAndExcludesSecretFields(t *testing.T) {
	st := openMem(t)
	now := time.Date(2026, 9, 21, 18, 0, 0, 0, time.UTC)
	group := &domain.Group{ID: "group-archive", Name: "Archive group", Enabled: true}
	if err := st.CreateGroup(group); err != nil {
		t.Fatal(err)
	}
	schedule := &domain.Schedule{ID: "schedule-archive", Kind: domain.ScheduleRecurring, RRULE: "FREQ=DAILY", HumanSummary: "Daily"}
	if err := st.CreateSchedule(schedule); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{ID: "task-archive", Name: "Archive", GroupID: group.ID, Command: "secret-command", Args: []string{"secret-argument"}, Env: map[string]string{"TOKEN": "secret-environment"}, Enabled: true, Timezone: "UTC", ScheduleID: schedule.ID, State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	ended := now.Add(-time.Hour)
	if err := st.CreateRun(&domain.Run{ID: "run-archive", TaskID: task.ID, ScheduledFor: ended, EndedAt: &ended, Outcome: domain.OutcomeFailure, Output: "secret-output", Trigger: domain.TriggerSchedule}); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateAlert(&domain.Alert{ID: "alert-archive", TaskID: task.ID, RunID: "run-archive", Severity: domain.SeverityError, Kind: domain.AlertRunFailed, Message: "secret-alert-message", CreatedAt: ended}); err != nil {
		t.Fatal(err)
	}

	results, err := st.SearchFacts("archive", nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 4 {
		t.Fatalf("results = %+v, want four non-temporal kinds", results)
	}
	wantKinds := []domain.SearchKind{domain.SearchKindTask, domain.SearchKindGroup, domain.SearchKindFailure, domain.SearchKindAlert}
	for index, kind := range wantKinds {
		if results[index].Kind != kind {
			t.Fatalf("result %d kind = %q, want %q", index, results[index].Kind, kind)
		}
	}
	for _, secret := range []string{"secret-command", "secret-argument", "secret-environment", "secret-output", "secret-alert-message"} {
		matches, searchErr := st.SearchFacts(secret, nil, 50)
		if searchErr != nil {
			t.Fatal(searchErr)
		}
		if len(matches) != 0 {
			t.Fatalf("secret query %q returned %+v", secret, matches)
		}
	}
}

func TestSearchFactsTreatsWildcardsLiterallyAndAppliesGlobalBound(t *testing.T) {
	st := openMem(t)
	for _, task := range []*domain.Task{
		{ID: "literal-percent", Name: "100% ready", Timezone: "UTC", State: domain.TaskActive},
		{ID: "ordinary", Name: "1000 ready", Timezone: "UTC", State: domain.TaskActive},
		{ID: "second", Name: "100% second", Timezone: "UTC", State: domain.TaskActive},
	} {
		if err := st.CreateTask(task); err != nil {
			t.Fatal(err)
		}
	}

	results, err := st.SearchFacts("100%", map[domain.SearchKind]bool{domain.SearchKindTask: true}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("limit-plus-one results = %+v, want 2", results)
	}
	for _, result := range results {
		if result.ObjectID == "ordinary" {
			t.Fatalf("LIKE wildcard was not escaped: %+v", results)
		}
	}
}

func TestSearchFactsHandlesEmptyInputDisabledObjectsAndKindIsolation(t *testing.T) {
	st := openMem(t)
	if values, err := st.SearchFacts("  ", nil, 50); err != nil || len(values) != 0 {
		t.Fatalf("empty search = %+v, %v", values, err)
	}
	if values, err := st.SearchFacts("anything", nil, 0); err != nil || len(values) != 0 {
		t.Fatalf("zero-limit search = %+v, %v", values, err)
	}
	group := &domain.Group{ID: "disabled-group", Name: "Disabled group", Enabled: false}
	if err := st.CreateGroup(group); err != nil {
		t.Fatal(err)
	}
	task := &domain.Task{ID: "disabled-task", Name: "Disabled task", GroupID: group.ID, Enabled: false, Timezone: "UTC", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	alert := &domain.Alert{ID: "daemon-warning", Severity: domain.SeverityWarning, Kind: domain.AlertService, Message: "not searchable", CreatedAt: time.Now().UTC()}
	if err := st.CreateAlert(alert); err != nil {
		t.Fatal(err)
	}

	tasks, err := st.SearchFacts("disabled", map[domain.SearchKind]bool{domain.SearchKindTask: true}, 50)
	if err != nil || len(tasks) != 1 || tasks[0].Kind != domain.SearchKindTask || tasks[0].Enabled == nil || *tasks[0].Enabled || len(tasks[0].ActionHints) != 1 || tasks[0].ActionHints[0] != domain.SearchActionOpen {
		t.Fatalf("disabled task search = %+v, %v", tasks, err)
	}
	groups, err := st.SearchFacts("disabled", map[domain.SearchKind]bool{domain.SearchKindGroup: true}, 50)
	if err != nil || len(groups) != 1 || groups[0].Context != "disabled" {
		t.Fatalf("disabled group search = %+v, %v", groups, err)
	}
	alerts, err := st.SearchFacts("warning", map[domain.SearchKind]bool{domain.SearchKindAlert: true}, 50)
	if err != nil || len(alerts) != 1 || alerts[0].Name != "Daemon alert" || alerts[0].Kind != domain.SearchKindAlert {
		t.Fatalf("daemon alert search = %+v, %v", alerts, err)
	}
	if failures, searchErr := st.SearchFacts("disabled", map[domain.SearchKind]bool{domain.SearchKindFailure: true}, 50); searchErr != nil || len(failures) != 0 {
		t.Fatalf("failure-only search = %+v, %v", failures, searchErr)
	}
}
