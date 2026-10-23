package store

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/secretstore"
)

func TestMigrationV23PreservesWebhookState(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "prior.db")
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec("CREATE TABLE schema_version(version INTEGER NOT NULL)"); err != nil {
		t.Fatal(err)
	}
	for _, m := range migrations {
		if m.version > 22 {
			break
		}
		if _, err = db.Exec(m.stmts); err != nil {
			t.Fatalf("seed migration %d: %v", m.version, err)
		}
		if _, err = db.Exec("INSERT INTO schema_version(version) VALUES(?)", m.version); err != nil {
			t.Fatal(err)
		}
	}
	endpoint, err := secretstore.Protect("https://example.test/private")
	if err != nil {
		t.Fatal(err)
	}
	auth, err := secretstore.Protect("Bearer private")
	if err != nil {
		t.Fatal(err)
	}
	now := "2026-10-21T00:00:00Z"
	if _, err = db.Exec("INSERT INTO notification_channels(id,name,kind,endpoint,endpoint_summary,authorization,enabled,created_at,updated_at) VALUES('old','old','webhook',?,'https://example.test',?,1,?,?)", endpoint, auth, now, now); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO groups(id,name,enabled,created_at,updated_at) VALUES('group','group',1,?,?)", now, now); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO notification_assignments(id,channel_id,group_id,on_failure,created_at,updated_at) VALUES('assignment','old','group',1,?,?)", now, now); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO notification_daemon_health_states(channel_id,next_due_at,updated_at) VALUES('old',?,?)", now, now); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO notification_deliveries(id,channel_id,channel_name,destination_summary,endpoint,authorization,event_kind,payload,state,next_attempt_at,created_at) VALUES('queued','old','old','https://example.test',?,?,'test','{}','pending',?,?)", endpoint, auth, now, now); err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	st, err := Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	c, err := st.GetNotificationChannel("old")
	if err != nil || c.Endpoint != "https://example.test/private" || c.Authorization != "Bearer private" || c.Kind != domain.NotificationChannelWebhook {
		t.Fatalf("preserved channel: %+v %v", c, err)
	}
	values, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
	if err != nil || len(values) != 1 || values[0].Kind != domain.NotificationChannelWebhook || values[0].Authorization != "Bearer private" {
		t.Fatalf("preserved delivery: %v", err)
	}
	for _, table := range []string{"notification_assignments", "notification_daemon_health_states"} {
		var n int
		if err := st.db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n); err != nil || n != 1 {
			t.Fatalf("preserved %s: %d %v", table, n, err)
		}
	}
	var v int
	if err := st.db.QueryRow("SELECT MAX(version) FROM schema_version").Scan(&v); err != nil || v != 24 {
		t.Fatalf("version %d %v", v, err)
	}
	var table string
	var row int
	var parent string
	var fk int
	if err := st.db.QueryRow("PRAGMA foreign_key_check").Scan(&table, &row, &parent, &fk); err != sql.ErrNoRows {
		t.Fatalf("foreign keys: %v", err)
	}
}

func TestSMTPStoredSnapshotsAndRecipientRecovery(t *testing.T) {
	st := openMem(t)
	config := domain.SMTPConfig{Host: "mail.example.test", Port: 25, TLSMode: "plain", From: "sender@example.test", Recipients: []string{"one@example.test", "two@example.test"}, TimeoutSeconds: 10}
	c := domain.NotificationChannel{Name: "mail", Kind: domain.NotificationChannelSMTP, SMTP: &config, Authorization: "credential-snapshot", Enabled: true}
	if err := st.CreateNotificationChannel(&c); err != nil {
		t.Fatal(err)
	}
	queued, err := st.CreateTestNotificationDelivery(c.ID)
	if err != nil {
		t.Fatal(err)
	}
	config.Host = "changed.example.test"
	c.SMTP = &config
	if err := st.UpdateNotificationChannel(c); err != nil {
		t.Fatal(err)
	}
	claimed, err := st.ClaimNotificationDeliveries(1, time.Now().Add(time.Minute))
	if err != nil || len(claimed) != 1 {
		t.Fatalf("claim: %v", err)
	}
	d := claimed[0]
	if d.SMTP.Host != "mail.example.test" || d.Authorization != "credential-snapshot" || len(d.Recipients) != 2 {
		t.Fatal("snapshot was not frozen")
	}
	d.Recipients[0].State = "accepted"
	d.Recipients[1].State = "submitting"
	if err := st.CheckpointSMTPRecipients(d.ID, d.Recipients); err != nil {
		t.Fatal(err)
	}
	if _, err := st.RecoverNotificationDeliveries(time.Now()); err != nil {
		t.Fatal(err)
	}
	claimed, err = st.ClaimNotificationDeliveries(1, time.Now().Add(time.Minute))
	if err != nil || len(claimed) != 1 {
		t.Fatalf("reclaim: %v", err)
	}
	d = claimed[0]
	if d.Recipients[0].State != "accepted" || d.Recipients[1].State != "uncertain" {
		t.Fatalf("recovery: %+v", d.Recipients)
	}
	d.Recipients[1].State = "permanent_failure"
	d.Recipients[1].Code = 550
	if err := st.FinishSMTPAttempt(d, domain.SMTPAttemptResult{Recipients: d.Recipients, Status: 550, Diagnostic: "One recipient rejected the message."}, time.Now(), time.Now()); err != nil {
		t.Fatal(err)
	}
	records, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
	if err != nil || len(records) != 1 || records[0].ID != queued.ID || records[0].State != domain.NotificationDeliveryFailed || records[0].Authorization != "" || records[0].Recipients[0].State != "accepted" {
		t.Fatalf("terminal: %v %+v", err, records)
	}
}

func TestSMTPAttemptBudgetAndCredentialLifecycle(t *testing.T) {
	for _, eventualSuccess := range []bool{false, true} {
		t.Run(fmt.Sprint(eventualSuccess), func(t *testing.T) {
			st := openMem(t)
			cfg := domain.SMTPConfig{Host: "mail.example.test", Port: 25, TLSMode: "plain", From: "sender@example.test", Recipients: []string{"ops@example.test"}, TimeoutSeconds: 10}
			channel := domain.NotificationChannel{Name: "mail", Kind: domain.NotificationChannelSMTP, SMTP: &cfg, Authorization: "owned snapshot", Enabled: true}
			if err := st.CreateNotificationChannel(&channel); err != nil {
				t.Fatal(err)
			}
			if _, err := st.CreateTestNotificationDelivery(channel.ID); err != nil {
				t.Fatal(err)
			}
			now := time.Now().UTC()
			for attempt := 1; attempt <= 3; attempt++ {
				claimed, err := st.ClaimNotificationDeliveries(1, now.Add(time.Duration(attempt)*time.Minute))
				if err != nil || len(claimed) != 1 {
					t.Fatalf("claim %d %v %+v", attempt, err, claimed)
				}
				d := claimed[0]
				if d.Attempts != attempt || d.Authorization != "owned snapshot" {
					t.Fatalf("attempt snapshot %+v", d)
				}
				d.Recipients[0].State = "temporary_failure"
				status := 450
				if eventualSuccess && attempt == 3 {
					d.Recipients[0].State = "accepted"
					status = 250
				}
				result := domain.SMTPAttemptResult{Recipients: d.Recipients, Status: status, Diagnostic: "Authored result."}
				if err := st.FinishSMTPAttempt(d, result, now, now); err != nil {
					t.Fatal(err)
				}
				if err := st.FinishSMTPAttempt(d, result, now, now); err == nil {
					t.Fatal("stale claimed attempt overwrote completion")
				}
				records, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
				if err != nil || len(records) != 1 {
					t.Fatalf("records %v", err)
				}
				got := records[0]
				wantError := "Authored result."
				if eventualSuccess && attempt == 3 {
					wantError = ""
				}
				if got.Diagnostic != result.Diagnostic {
					t.Fatalf("diagnostic lost: %+v", got)
				}
				if got.LastError != wantError {
					t.Fatalf("error field confused transport success: state=%s error=%q", got.State, got.LastError)
				}
				if attempt < 3 {
					if got.State != domain.NotificationDeliveryPending || got.Authorization != "owned snapshot" || got.CompletedAt != nil {
						t.Fatalf("retry lifecycle %+v", got)
					}
				} else {
					want := domain.NotificationDeliveryFailed
					if eventualSuccess {
						want = domain.NotificationDeliverySucceeded
					}
					if got.State != want || got.Authorization != "" || got.CompletedAt == nil {
						t.Fatalf("terminal lifecycle %+v", got)
					}
				}
			}
			claimed, err := st.ClaimNotificationDeliveries(1, now.Add(time.Hour))
			if err != nil || len(claimed) != 0 {
				t.Fatalf("exhausted work reclaimed %v %+v", err, claimed)
			}
			if err := st.CheckpointSMTPRecipients("missing", nil); err == nil {
				t.Fatal("empty checkpoint accepted")
			}
			if err := st.FinishSMTPAttempt(domain.NotificationDelivery{}, domain.SMTPAttemptResult{}, now, now); err == nil {
				t.Fatal("empty results accepted")
			}
			if _, err := decodeSMTPConfig("{"); err == nil {
				t.Fatal("corrupt settings accepted")
			}
		})
	}
}
