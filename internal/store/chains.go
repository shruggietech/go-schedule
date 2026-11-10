package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
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
	now := s.now().UTC()
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

// ReplaceCompletionChains atomically replaces existing chain definitions. It
// is used for reviewed portable imports so valid final graphs that swap two
// relationship identities do not fail against a transient unique constraint.
func (s *Store) ReplaceCompletionChains(chains []domain.CompletionChain) error {
	if len(chains) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin replace completion chains: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	oldTargets := make([]string, 0, len(chains))
	for i := range chains {
		var created string
		var targetID string
		if err := tx.QueryRow(`SELECT created_at,target_task_id FROM completion_chains WHERE id=?`, chains[i].ID).Scan(&created, &targetID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("store: read completion chain for replacement: %w", err)
		}
		oldTargets = append(oldTargets, targetID)
		chains[i].CreatedAt, _ = parseTime(created)
		chains[i].UpdatedAt = time.Now().UTC()
	}
	for _, chain := range chains {
		if _, err := tx.Exec(`DELETE FROM completion_chains WHERE id=?`, chain.ID); err != nil {
			return fmt.Errorf("store: stage completion chain replacement: %w", err)
		}
	}
	for _, chain := range chains {
		if err := validateCompletionChain(tx, &chain, ""); err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO completion_chains(id,source_task_id,target_task_id,on_outcome,created_at,updated_at) VALUES(?,?,?,?,?,?)`, chain.ID, chain.SourceTaskID, chain.TargetTaskID, string(chain.OnOutcome), fmtTime(chain.CreatedAt), fmtTime(chain.UpdatedAt)); err != nil {
			if isUniqueConstraint(err) {
				return ErrDuplicateChain
			}
			return fmt.Errorf("store: replace completion chain: %w", err)
		}
	}
	for _, targetID := range oldTargets {
		if err := disableIfNotActivationReady(tx, targetID); err != nil {
			return err
		}
	}
	return tx.Commit()
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
	graph, err := automationGraph(tx, excludeID)
	if err != nil {
		return err
	}
	graph[c.SourceTaskID] = append(graph[c.SourceTaskID], c.TargetTaskID)
	if reaches(graph, c.TargetTaskID, c.SourceTaskID, map[string]bool{}) {
		return ErrChainCycle
	}
	return nil
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
	return s.recordRunAndDeliveries(r, incomingDeliveryID)
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
	rows, err := tx.Query(deliverySelect+` WHERE state=? ORDER BY created_at,id LIMIT ?`, string(domain.DeliveryPending), limit)
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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	d, err := scanDelivery(tx.QueryRow(deliverySelect+" WHERE id=?", id))
	if err != nil {
		return err
	}
	if d.State == domain.DeliveryCompleted || d.State == domain.DeliveryResolved {
		return tx.Commit()
	}
	now := s.now().UTC()
	if _, err := tx.Exec("UPDATE completion_deliveries SET state=?,completed_at=?,resolution=? WHERE id=? AND state IN (?,?)", string(domain.DeliveryResolved), fmtTime(now), reason, id, string(domain.DeliveryPending), string(domain.DeliveryClaimed)); err != nil {
		return err
	}
	if err := resolveJoinDelivery(tx, d, reason, now); err != nil {
		return err
	}
	return tx.Commit()
}

// ListCompletionDeliveries returns delivery state for tests and diagnostics.
func (s *Store) ListCompletionDeliveries() ([]domain.CompletionDelivery, error) {
	rows, err := s.db.Query(deliverySelect + ` ORDER BY created_at,id`)
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

const deliverySelect = `SELECT id,chain_id,source_task_id,target_task_id,source_run_id,state,attempts,created_at,claimed_at,completed_at,target_run_id,resolution,role,join_instance_id,join_branch_id FROM completion_deliveries`

func scanDelivery(sc scanner) (domain.CompletionDelivery, error) {
	var d domain.CompletionDelivery
	var state, created string
	var claimed, completed, targetRun sql.NullString
	if err := sc.Scan(&d.ID, &d.ChainID, &d.SourceTaskID, &d.TargetTaskID, &d.SourceRunID, &state, &d.Attempts, &created, &claimed, &completed, &targetRun, &d.Resolution, &d.Role, &d.JoinInstanceID, &d.JoinBranchID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return d, ErrNotFound
		}
		return d, fmt.Errorf("store: scan completion delivery: %w", err)
	}
	d.State = domain.DeliveryState(state)
	d.CreatedAt, _ = parseTime(created)
	d.ClaimedAt, _ = parseTimePtr(claimed)
	d.CompletedAt, _ = parseTimePtr(completed)
	d.TargetRunID = targetRun.String
	return d, nil
}
