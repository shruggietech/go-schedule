package store

import "testing"

func TestMigrationV20AddsNotificationConditionState(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	checks := []struct {
		table  string
		column string
	}{
		{"runs", "start_failed"},
		{"notification_channels", "health_interval_seconds"},
		{"notification_assignments", "failure_threshold"},
		{"notification_assignments", "on_failure_to_start"},
		{"notification_assignments", "duration_threshold_seconds"},
		{"notification_assignments", "on_recovery"},
		{"notification_assignments", "reminder_interval_seconds"},
		{"notification_assignments", "quiet_period_seconds"},
		{"notification_deliveries", "condition_kind"},
		{"notification_deliveries", "condition_summary"},
	}
	for _, check := range checks {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?`, check.table, check.column).Scan(&count); err != nil || count != 1 {
			t.Fatalf("%s.%s count=%d err=%v", check.table, check.column, count, err)
		}
	}
	for _, table := range []string{"notification_condition_states", "notification_daemon_health_states"} {
		var count int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count); err != nil || count != 1 {
			t.Fatalf("table %s count=%d err=%v", table, count, err)
		}
	}
	var assignmentSchema string
	if err := st.db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='notification_assignments'`).Scan(&assignmentSchema); err != nil {
		t.Fatal(err)
	}
	if _, err := st.db.Exec(`INSERT INTO notification_channels(id,name,kind,endpoint,endpoint_summary,enabled,created_at,updated_at) VALUES('channel-advanced','Advanced','webhook','https://example.test','https://example.test',1,'2026-09-21T00:00:00Z','2026-09-21T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := st.db.Exec(`INSERT INTO groups(id,name,enabled,created_at,updated_at) VALUES('group-advanced','Advanced',1,'2026-09-21T00:00:00Z','2026-09-21T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := st.db.Exec(`INSERT INTO notification_assignments(id,channel_id,group_id,on_failure_to_start,created_at,updated_at) VALUES('assignment-advanced','channel-advanced','group-advanced',1,'2026-09-21T00:00:00Z','2026-09-21T00:00:00Z')`); err != nil {
		t.Fatalf("advanced-only assignment rejected by schema %q: %v", assignmentSchema, err)
	}
}
