package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shruggietech/go-scheduler/internal/domain"
)

// CreateTrigger inserts a trigger, assigning an ID when empty.
func (s *Store) CreateTrigger(t *domain.Trigger) error {
	if t.ID == "" {
		t.ID = newID()
	}
	if t.OnOutcome == "" {
		t.OnOutcome = domain.OnSuccess
	}
	_, err := s.db.Exec(
		`INSERT INTO triggers(id,source_task_id,target_task_id,on_outcome,dedup_key,dedup_window_ns)
		 VALUES(?,?,?,?,?,?)`,
		t.ID, t.SourceTaskID, t.TargetTaskID, string(t.OnOutcome), t.DedupKey, int64(t.DedupWindow),
	)
	if err != nil {
		return fmt.Errorf("store: create trigger: %w", err)
	}
	return nil
}

// ListTriggers returns all triggers.
func (s *Store) ListTriggers() ([]domain.Trigger, error) {
	return s.queryTriggers(`SELECT id,source_task_id,target_task_id,on_outcome,dedup_key,dedup_window_ns FROM triggers ORDER BY id`)
}

// ListTriggersBySource returns triggers fired by the given source task.
func (s *Store) ListTriggersBySource(sourceTaskID string) ([]domain.Trigger, error) {
	return s.queryTriggers(
		`SELECT id,source_task_id,target_task_id,on_outcome,dedup_key,dedup_window_ns FROM triggers WHERE source_task_id=?`,
		sourceTaskID)
}

func (s *Store) queryTriggers(q string, args ...any) ([]domain.Trigger, error) {
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list triggers: %w", err)
	}
	defer rows.Close()
	var out []domain.Trigger
	for rows.Next() {
		var t domain.Trigger
		var outcome string
		var windowNS int64
		if err := rows.Scan(&t.ID, &t.SourceTaskID, &t.TargetTaskID, &outcome, &t.DedupKey, &windowNS); err != nil {
			return nil, fmt.Errorf("store: scan trigger: %w", err)
		}
		t.OnOutcome = domain.TriggerOutcome(outcome)
		t.DedupWindow = time.Duration(windowNS)
		out = append(out, t)
	}
	return out, rows.Err()
}

// DeleteTrigger removes a trigger.
func (s *Store) DeleteTrigger(id string) error {
	res, err := s.db.Exec(`DELETE FROM triggers WHERE id=?`, id)
	return affected(res, err, "delete trigger")
}

// ClaimEvent attempts to claim a logical event for (triggerID, key). It returns
// claimed=true when this is a fresh event (no prior claim within window),
// inserting a ledger row with executed=false. A repeat within the window returns
// claimed=false (deduplicated). A window of 0 means "dedupe forever for this key".
func (s *Store) ClaimEvent(triggerID, key string, window time.Duration, now time.Time) (bool, error) {
	var firstSeen string
	row := s.db.QueryRow(`SELECT first_seen_at FROM dedup_ledger WHERE trigger_id=? AND dedup_key=?`, triggerID, key)
	switch err := row.Scan(&firstSeen); {
	case err == nil:
		seen, perr := parseTime(firstSeen)
		if perr != nil {
			return false, fmt.Errorf("store: dedup first_seen: %w", perr)
		}
		if window <= 0 || now.Sub(seen) < window {
			return false, nil // duplicate within window (or forever)
		}
		// Window expired: refresh the claim for a new logical event.
		if _, err := s.db.Exec(`UPDATE dedup_ledger SET first_seen_at=?, executed=0 WHERE trigger_id=? AND dedup_key=?`,
			fmtTime(now), triggerID, key); err != nil {
			return false, fmt.Errorf("store: refresh dedup: %w", err)
		}
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		if _, err := s.db.Exec(`INSERT INTO dedup_ledger(trigger_id,dedup_key,first_seen_at,executed) VALUES(?,?,?,0)`,
			triggerID, key, fmtTime(now)); err != nil {
			return false, fmt.Errorf("store: claim event: %w", err)
		}
		return true, nil
	default:
		return false, fmt.Errorf("store: claim lookup: %w", err)
	}
}

// MarkExecuted records that the target run for a claimed event was dispatched.
func (s *Store) MarkExecuted(triggerID, key string) error {
	_, err := s.db.Exec(`UPDATE dedup_ledger SET executed=1 WHERE trigger_id=? AND dedup_key=?`, triggerID, key)
	if err != nil {
		return fmt.Errorf("store: mark executed: %w", err)
	}
	return nil
}

// PendingClaims returns claimed-but-not-executed ledger entries (for
// at-least-once recovery after a restart).
func (s *Store) PendingClaims() ([]domain.DedupLedger, error) {
	rows, err := s.db.Query(`SELECT trigger_id,dedup_key,first_seen_at,executed FROM dedup_ledger WHERE executed=0`)
	if err != nil {
		return nil, fmt.Errorf("store: pending claims: %w", err)
	}
	defer rows.Close()
	var out []domain.DedupLedger
	for rows.Next() {
		var d domain.DedupLedger
		var seen string
		var exec int
		if err := rows.Scan(&d.TriggerID, &d.DedupKey, &seen, &exec); err != nil {
			return nil, fmt.Errorf("store: scan claim: %w", err)
		}
		d.FirstSeenAt, _ = parseTime(seen)
		d.Executed = exec != 0
		out = append(out, d)
	}
	return out, rows.Err()
}

// GetTrigger returns a trigger by ID.
func (s *Store) GetTrigger(id string) (domain.Trigger, error) {
	ts, err := s.queryTriggers(
		`SELECT id,source_task_id,target_task_id,on_outcome,dedup_key,dedup_window_ns FROM triggers WHERE id=?`, id)
	if err != nil {
		return domain.Trigger{}, err
	}
	if len(ts) == 0 {
		return domain.Trigger{}, ErrNotFound
	}
	return ts[0], nil
}
