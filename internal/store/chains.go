package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// ErrChainCycle reports a completion relationship that would make the graph cyclic.
var ErrChainCycle = errors.New("store: completion chain would create a cycle")

// ErrDuplicateChain reports an existing source/target/outcome relationship.
var ErrDuplicateChain = errors.New("store: completion chain already exists")

// ErrInvalidOutcome reports an unsupported completion outcome selector.
var ErrInvalidOutcome = errors.New("store: invalid completion outcome")

// CreateCompletionChain validates and inserts a completion relationship.
func (s *Store) CreateCompletionChain(c *domain.CompletionChain) error {
	if c.ID == "" {
		c.ID = newID()
	}
	now := time.Now().UTC()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin create completion chain: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := validateCompletionChain(tx, c, ""); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO completion_chains(id,source_task_id,target_task_id,on_outcome,created_at,updated_at) VALUES(?,?,?,?,?,?)`,
		c.ID, c.SourceTaskID, c.TargetTaskID, string(c.OnOutcome), fmtTime(c.CreatedAt), fmtTime(c.UpdatedAt))
	if err != nil {
		if isUniqueConstraint(err) {
			return ErrDuplicateChain
		}
		return fmt.Errorf("store: create completion chain: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit completion chain: %w", err)
	}
	return nil
}

// UpdateCompletionChain validates and replaces a relationship's mutable fields.
func (s *Store) UpdateCompletionChain(c *domain.CompletionChain) error {
	existing, err := s.GetCompletionChain(c.ID)
	if err != nil {
		return err
	}
	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin update completion chain: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := validateCompletionChain(tx, c, c.ID); err != nil {
		return err
	}
	res, err := tx.Exec(`UPDATE completion_chains SET source_task_id=?,target_task_id=?,on_outcome=?,updated_at=? WHERE id=?`,
		c.SourceTaskID, c.TargetTaskID, string(c.OnOutcome), fmtTime(c.UpdatedAt), c.ID)
	if err != nil {
		if isUniqueConstraint(err) {
			return ErrDuplicateChain
		}
		return fmt.Errorf("store: update completion chain: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	if existing.TargetTaskID != c.TargetTaskID {
		if err := disableIfNotActivationReady(tx, existing.TargetTaskID); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit completion chain update: %w", err)
	}
	return nil
}

// GetCompletionChain returns one relationship enriched with current task names.
func (s *Store) GetCompletionChain(id string) (domain.CompletionChain, error) {
	return scanCompletionChain(s.db.QueryRow(completionChainSelect+` WHERE c.id=?`, id))
}

// ListCompletionChains returns all relationships enriched with task names.
func (s *Store) ListCompletionChains() ([]domain.CompletionChain, error) {
	rows, err := s.db.Query(completionChainSelect + ` ORDER BY s.name,t.name,c.on_outcome`)
	if err != nil {
		return nil, fmt.Errorf("store: list completion chains: %w", err)
	}
	defer rows.Close()
	var out []domain.CompletionChain
	for rows.Next() {
		c, err := scanCompletionChain(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// DeleteCompletionChain removes a relationship. Already-created deliveries
// retain the snapshotted identity so the engine can terminally resolve them
// instead of silently losing durable work.
func (s *Store) DeleteCompletionChain(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin delete completion chain: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID string
	if err := tx.QueryRow(`SELECT target_task_id FROM completion_chains WHERE id=?`, id).Scan(&targetID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("store: find completion chain target: %w", err)
	}
	res, err := tx.Exec(`DELETE FROM completion_chains WHERE id=?`, id)
	if err := affected(res, err, "delete completion chain"); err != nil {
		return err
	}
	if err := disableIfNotActivationReady(tx, targetID); err != nil {
		return err
	}
	return tx.Commit()
}

const completionChainSelect = `SELECT c.id,c.source_task_id,s.name,c.target_task_id,t.name,c.on_outcome,c.created_at,c.updated_at
 FROM completion_chains c JOIN tasks s ON s.id=c.source_task_id JOIN tasks t ON t.id=c.target_task_id`

func scanCompletionChain(sc scanner) (domain.CompletionChain, error) {
	var c domain.CompletionChain
	var outcome, created, updated string
	if err := sc.Scan(&c.ID, &c.SourceTaskID, &c.SourceTaskName, &c.TargetTaskID, &c.TargetTaskName, &outcome, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CompletionChain{}, ErrNotFound
		}
		return domain.CompletionChain{}, fmt.Errorf("store: scan completion chain: %w", err)
	}
	c.OnOutcome = domain.CompletionOutcome(outcome)
	c.CreatedAt, _ = parseTime(created)
	c.UpdatedAt, _ = parseTime(updated)
	return c, nil
}

func validateCompletionChain(tx *sql.Tx, c *domain.CompletionChain, excludeID string) error {
	if c.OnOutcome != domain.CompletionOnSuccess && c.OnOutcome != domain.CompletionOnFailure && c.OnOutcome != domain.CompletionOnAny {
		return ErrInvalidOutcome
	}
	if c.SourceTaskID == "" || c.TargetTaskID == "" {
		return ErrNotFound
	}
	if c.SourceTaskID == c.TargetTaskID {
		return ErrChainCycle
	}
	for _, id := range []string{c.SourceTaskID, c.TargetTaskID} {
		var one int
		if err := tx.QueryRow(`SELECT 1 FROM tasks WHERE id=?`, id).Scan(&one); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("store: validate completion task: %w", err)
		}
	}
	rows, err := tx.Query(`SELECT source_task_id,target_task_id FROM completion_chains WHERE id<>?`, excludeID)
	if err != nil {
		return fmt.Errorf("store: list completion graph: %w", err)
	}
	defer rows.Close()
	graph := map[string][]string{}
	for rows.Next() {
		var source, target string
		if err := rows.Scan(&source, &target); err != nil {
			return fmt.Errorf("store: scan completion graph: %w", err)
		}
		graph[source] = append(graph[source], target)
	}
	graph[c.SourceTaskID] = append(graph[c.SourceTaskID], c.TargetTaskID)
	if reaches(graph, c.TargetTaskID, c.SourceTaskID, map[string]bool{}) {
		return ErrChainCycle
	}
	return rows.Err()
}

func reaches(graph map[string][]string, from, target string, seen map[string]bool) bool {
	if from == target {
		return true
	}
	if seen[from] {
		return false
	}
	seen[from] = true
	for _, next := range graph[from] {
		if reaches(graph, next, target, seen) {
			return true
		}
	}
	return false
}

func isUniqueConstraint(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "constraint failed"))
}

// RecordRunAndCreateDeliveries atomically records a run, completes its incoming
// delivery, and creates matching outgoing deliveries. Only success and failure
// create downstream work.
func (s *Store) RecordRunAndCreateDeliveries(r *domain.Run, incomingDeliveryID string) error {
	if r.ID == "" {
		r.ID = newID()
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin record run: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.Exec(`INSERT INTO runs(id,task_id,scheduled_for,started_at,ended_at,outcome,exit_code,output,output_truncated,trigger,source_task_id,source_run_id,source_trigger_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		r.ID, r.TaskID, fmtTime(r.ScheduledFor), fmtTimePtr(r.StartedAt), fmtTimePtr(r.EndedAt), string(r.Outcome), nullInt(r.ExitCode), r.Output, boolToInt(r.OutputTruncated), string(r.Trigger), nullStr(r.SourceTaskID), nullStr(r.SourceRunID), nullStr(r.SourceTriggerID))
	if err != nil {
		return fmt.Errorf("store: record run: %w", err)
	}
	now := time.Now().UTC()
	if incomingDeliveryID != "" {
		res, err := tx.Exec(`UPDATE completion_deliveries SET state=?,completed_at=?,target_run_id=? WHERE id=? AND state=?`,
			string(domain.DeliveryCompleted), fmtTime(now), r.ID, incomingDeliveryID, string(domain.DeliveryClaimed))
		if err != nil {
			return fmt.Errorf("store: complete delivery: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("store: complete delivery %s: %w", incomingDeliveryID, ErrNotFound)
		}
	}
	if r.Outcome == domain.OutcomeSuccess || r.Outcome == domain.OutcomeFailure {
		rows, err := tx.Query(`SELECT id,target_task_id FROM completion_chains WHERE source_task_id=? AND (on_outcome=? OR on_outcome=?)`,
			r.TaskID, string(r.Outcome), string(domain.CompletionOnAny))
		if err != nil {
			return fmt.Errorf("store: match completion chains: %w", err)
		}
		type match struct{ chainID, targetID string }
		var matches []match
		for rows.Next() {
			var m match
			if err := rows.Scan(&m.chainID, &m.targetID); err != nil {
				_ = rows.Close()
				return fmt.Errorf("store: scan completion match: %w", err)
			}
			matches = append(matches, m)
		}
		if err := rows.Close(); err != nil {
			return fmt.Errorf("store: close completion matches: %w", err)
		}
		for _, m := range matches {
			_, err := tx.Exec(`INSERT OR IGNORE INTO completion_deliveries(id,chain_id,source_task_id,target_task_id,source_run_id,state,created_at) VALUES(?,?,?,?,?,?,?)`,
				newID(), m.chainID, r.TaskID, m.targetID, r.ID, string(domain.DeliveryPending), fmtTime(now))
			if err != nil {
				return fmt.Errorf("store: create completion delivery: %w", err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit run and deliveries: %w", err)
	}
	return nil
}

// RecoverCompletionDeliveries makes interrupted claims replayable.
func (s *Store) RecoverCompletionDeliveries() (int64, error) {
	res, err := s.db.Exec(`UPDATE completion_deliveries SET state=?,claimed_at=NULL,resolution='recovered after unclean shutdown' WHERE state=?`,
		string(domain.DeliveryPending), string(domain.DeliveryClaimed))
	if err != nil {
		return 0, fmt.Errorf("store: recover completion deliveries: %w", err)
	}
	n, err := res.RowsAffected()
	return n, err
}

// ClaimCompletionDeliveries claims up to limit pending deliveries in creation order.
func (s *Store) ClaimCompletionDeliveries(limit int) ([]domain.CompletionDelivery, error) {
	if limit <= 0 {
		limit = 100
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("store: begin claim deliveries: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.Query(`SELECT id,chain_id,source_task_id,target_task_id,source_run_id,state,attempts,created_at,claimed_at,completed_at,target_run_id,resolution FROM completion_deliveries WHERE state=? ORDER BY created_at,id LIMIT ?`, string(domain.DeliveryPending), limit)
	if err != nil {
		return nil, fmt.Errorf("store: list pending deliveries: %w", err)
	}
	var out []domain.CompletionDelivery
	for rows.Next() {
		d, err := scanDelivery(rows)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		out = append(out, d)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("store: close pending deliveries: %w", err)
	}
	now := time.Now().UTC()
	for i := range out {
		res, err := tx.Exec(`UPDATE completion_deliveries SET state=?,attempts=attempts+1,claimed_at=? WHERE id=? AND state=?`,
			string(domain.DeliveryClaimed), fmtTime(now), out[i].ID, string(domain.DeliveryPending))
		if err != nil {
			return nil, fmt.Errorf("store: claim delivery: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return nil, fmt.Errorf("store: claim delivery %s: %w", out[i].ID, ErrNotFound)
		}
		out[i].State = domain.DeliveryClaimed
		out[i].Attempts++
		out[i].ClaimedAt = &now
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("store: commit delivery claims: %w", err)
	}
	return out, nil
}

// ResolveCompletionDelivery terminally records a delivery that did not execute.
func (s *Store) ResolveCompletionDelivery(id, reason string) error {
	res, err := s.db.Exec(`UPDATE completion_deliveries SET state=?,completed_at=?,resolution=? WHERE id=? AND state=?`,
		string(domain.DeliveryResolved), fmtTime(time.Now().UTC()), reason, id, string(domain.DeliveryClaimed))
	return affected(res, err, "resolve completion delivery")
}

// ListCompletionDeliveries returns delivery state for tests and diagnostics.
func (s *Store) ListCompletionDeliveries() ([]domain.CompletionDelivery, error) {
	rows, err := s.db.Query(`SELECT id,chain_id,source_task_id,target_task_id,source_run_id,state,attempts,created_at,claimed_at,completed_at,target_run_id,resolution FROM completion_deliveries ORDER BY created_at,id`)
	if err != nil {
		return nil, fmt.Errorf("store: list completion deliveries: %w", err)
	}
	defer rows.Close()
	var out []domain.CompletionDelivery
	for rows.Next() {
		d, err := scanDelivery(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func scanDelivery(sc scanner) (domain.CompletionDelivery, error) {
	var d domain.CompletionDelivery
	var state, created string
	var claimed, completed, targetRun sql.NullString
	if err := sc.Scan(&d.ID, &d.ChainID, &d.SourceTaskID, &d.TargetTaskID, &d.SourceRunID, &state, &d.Attempts, &created, &claimed, &completed, &targetRun, &d.Resolution); err != nil {
		return d, fmt.Errorf("store: scan completion delivery: %w", err)
	}
	d.State = domain.DeliveryState(state)
	d.CreatedAt, _ = parseTime(created)
	d.ClaimedAt, _ = parseTimePtr(claimed)
	d.CompletedAt, _ = parseTimePtr(completed)
	d.TargetRunID = targetRun.String
	return d, nil
}
