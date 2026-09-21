package store

import (
	"database/sql"
	"strings"
	"testing"
	"time"
)

func TestMigrationV21IndexesBoundedFailureSummaryQueries(t *testing.T) {
	st := openMem(t)
	var version int
	if err := st.db.QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version != 22 {
		t.Fatalf("schema version=%d, want 22", version)
	}
	rows, err := st.db.Query(`EXPLAIN QUERY PLAN SELECT COUNT(*) FROM runs WHERE outcome=? AND COALESCE(ended_at,scheduled_for)>=? AND COALESCE(ended_at,scheduled_for)<=?`, "failure", fmtTime(time.Now().Add(-24*time.Hour)), fmtTime(time.Now()))
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var plan strings.Builder
	for rows.Next() {
		var id, parent, unused int
		var detail string
		if err := rows.Scan(&id, &parent, &unused, &detail); err != nil {
			t.Fatal(err)
		}
		plan.WriteString(detail)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(plan.String(), "idx_runs_outcome_effective_time") {
		t.Fatalf("query plan does not use summary index: %s", plan.String())
	}
	var indexName string
	if err := st.db.QueryRow(`SELECT name FROM sqlite_master WHERE type='index' AND name='idx_runs_outcome_effective_time'`).Scan(&indexName); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if indexName == "" {
		t.Fatal("summary index was not created")
	}
}
