package store

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestBoundedHistoryPreservesEmptyText(t *testing.T) {
	st := openMem(t)
	task := domain.Task{Name: "silent task", Command: "true", Enabled: true, Timezone: "UTC", State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	for _, tc := range []struct {
		text      string
		truncated bool
	}{{"", false}, {"", true}, {"世界long", false}} {
		run := domain.Run{TaskID: task.ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual, Output: tc.text, OutputTruncated: tc.truncated}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
		values, err := st.ListRunsPage(task.ID, 0, 1, 6)
		if err != nil || len(values) != 1 {
			t.Fatalf("bounded run %v %+v", err, values)
		}
		want := tc.text
		if len(want) > 6 {
			want = want[:6]
		}
		if values[0].Output != want || values[0].OutputTruncated != (tc.truncated || len(tc.text) > 6) {
			t.Fatalf("run projection %+v", values[0])
		}
		alert := domain.Alert{TaskID: task.ID, Severity: domain.SeverityWarning, Kind: domain.AlertRunFailed, Message: tc.text, CreatedAt: now}
		if err := st.CreateAlert(&alert); err != nil {
			t.Fatal(err)
		}
		alerts, err := st.ListAlertsPage(false, 0, 1, 6)
		if err != nil || len(alerts) != 1 || alerts[0].Message != want || alerts[0].MessageTruncated != (len(tc.text) > 6) {
			t.Fatalf("alert projection %v %+v", err, alerts)
		}
		now = now.Add(time.Minute)
	}
	values, err := st.ListRunsPage(task.ID, 2, 1, 6)
	if err != nil || len(values) != 1 || values[0].Output != "" || values[0].OutputTruncated {
		t.Fatalf("empty-page offset %v %+v", err, values)
	}
}
