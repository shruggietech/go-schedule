package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/secretstore"
)

func TestCurrentSchemaPreservesSMTPHistoryAndWebhookErrors(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "v23.db")
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(currentSchema + "INSERT INTO schema_version(version) VALUES(24);"); err != nil {
		t.Fatal(err)
	}

	endpoint, err := secretstore.Protect("owned endpoint")
	if err != nil {
		t.Fatal(err)
	}
	auth, err := secretstore.Protect("owned credential")
	if err != nil {
		t.Fatal(err)
	}
	now := "2026-10-23T01:02:03Z"
	for _, tc := range []struct{ id, kind, state string }{{"accepted", "smtp", "succeeded"}, {"rejected", "smtp", "failed"}, {"retry", "smtp", "pending"}, {"webhook", "webhook", "succeeded"}} {
		_, err := db.Exec(`INSERT INTO notification_deliveries(id,channel_name,destination_summary,endpoint,authorization,event_kind,payload,state,attempts,next_attempt_at,created_at,completed_at,last_status,last_error,kind,recipients,diagnostic) VALUES(?,?,?,?,?,'test','{}',?,2,?,?,?,250,?,?,'[{"address":"ops@example.test","state":"accepted"}]',?)`, tc.id, "mail", "mail summary", endpoint, auth, tc.state, now, now, now, map[bool]string{true: "", false: "Original authored result."}[tc.id == "accepted"], tc.kind, map[bool]string{true: "", false: "Original authored result."}[tc.kind == "webhook"])
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	for reopen := 0; reopen < 2; reopen++ {
		st, err := Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		values, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
		if err != nil {
			t.Fatal(err)
		}
		if len(values) != 4 {
			t.Fatalf("history count %d", len(values))
		}
		for _, d := range values {
			wantDiagnostic, wantError := "Original authored result.", "Original authored result."
			if d.ID == "accepted" {
				wantError = ""
			}
			if d.Kind == domain.NotificationChannelWebhook {
				wantDiagnostic = ""
			}
			if d.Diagnostic != wantDiagnostic || d.LastError != wantError {
				t.Fatalf("history changed: %+v", d)
			}
			if d.Attempts != 2 || d.LastStatus != 250 || d.Authorization != "owned credential" || d.Endpoint != "owned endpoint" || len(d.Recipients) != 1 || d.CreatedAt.Format("2006-01-02T15:04:05Z") != now || d.CompletedAt == nil {
				t.Fatalf("historical fields lost: %+v", d)
			}
		}
		var n int
		if err := st.db.QueryRow("SELECT COUNT(*) FROM schema_version WHERE version=24").Scan(&n); err != nil || n != 1 {
			t.Fatalf("bootstrap repeated: %d %v", n, err)
		}
		rows, err := st.db.Query("PRAGMA foreign_key_check")
		if err != nil {
			t.Fatal(err)
		}
		if rows.Next() {
			t.Fatal("foreign key violation")
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
		if err := st.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
