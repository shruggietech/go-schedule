package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

const runSelect = "SELECT id,task_id,scheduled_for,started_at,ended_at,outcome,exit_code,output,output_truncated,trigger,source_task_id,source_run_id,source_trigger_id,source_watcher_id,start_failed,delivery_id,join_instance_id,join_branch_id,join_role FROM runs"

type runExecer interface {
	Exec(string, ...any) (sql.Result, error)
}

const runInsertSQL = "INSERT INTO runs(id,task_id,scheduled_for,started_at,ended_at,outcome,exit_code,output,output_truncated,trigger,source_task_id,source_run_id,source_trigger_id,source_watcher_id,start_failed,delivery_id,join_instance_id,join_branch_id,join_role) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

func runValues(r domain.Run) []any {
	return []any{r.ID, r.TaskID, fmtTime(r.ScheduledFor), fmtTimePtr(r.StartedAt), fmtTimePtr(r.EndedAt), string(r.Outcome), nullInt(r.ExitCode), r.Output, boolToInt(r.OutputTruncated), string(r.Trigger), nullStr(r.SourceTaskID), nullStr(r.SourceRunID), nullStr(r.SourceTriggerID), nullStr(r.SourceWatcherID), boolToInt(r.StartFailed), r.DeliveryID, r.JoinInstanceID, r.JoinBranchID, r.JoinRole}
}

func insertRun(q runExecer, r domain.Run) (sql.Result, error) {
	return q.Exec(runInsertSQL, runValues(r)...)
}

func (s *Store) recordRunAndDeliveries(r *domain.Run, incomingID string) error {
	if r.ID == "" {
		r.ID = newID()
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin record run: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var incoming domain.CompletionDelivery
	if incomingID != "" {
		incoming, err = scanDelivery(tx.QueryRow(deliverySelect+" WHERE id=?", incomingID))
		if err != nil {
			return err
		}
		if incoming.TargetTaskID != r.TaskID {
			return fmt.Errorf("store: delivery target does not match run: %w", ErrInvalidJoin)
		}
		correlateDeliveryRun(r, incoming)
	}
	// Insert first: fresh executions are the common path. Read immutable history only after an insertion conflict.
	liveTask, duplicate, err := s.insertObservedRun(tx, *r, incoming)
	if err != nil {
		return err
	}
	if duplicate {
		return tx.Commit()
	}
	if incomingID != "" {
		if incoming.State == domain.DeliveryCompleted {
			if incoming.TargetRunID != r.ID {
				return fmt.Errorf("store: delivery already completed by a different run")
			}
			if incoming.Role != domain.DeliveryRoleChain {
				var outcome domain.RunOutcome
				var started, ended sql.NullString
				if incoming.Role == domain.DeliveryRoleJoinBranch {
					err = tx.QueryRow("SELECT outcome,started_at,ended_at FROM parallel_join_invocations WHERE id=? AND instance_id=?", incoming.JoinBranchID, incoming.JoinInstanceID).Scan(&outcome, &started, &ended)
				} else {
					err = tx.QueryRow("SELECT target_outcome,target_started_at,target_ended_at FROM parallel_join_instances WHERE id=?", incoming.JoinInstanceID).Scan(&outcome, &started, &ended)
				}
				if err != nil {
					return err
				}
				if outcome != r.Outcome || started != fmtTimePtr(r.StartedAt) || ended != fmtTimePtr(r.EndedAt) {
					return fmt.Errorf("store: conflicting late native result %s", r.ID)
				}
			}
			return tx.Commit()
		}
		if incoming.State != domain.DeliveryClaimed && (incoming.State != domain.DeliveryResolved || incoming.Role == domain.DeliveryRoleChain) {
			return fmt.Errorf("store: delivery is not claimed: %w", ErrNotFound)
		}
		if r.Outcome == domain.OutcomeQueued {
			return fmt.Errorf("store: queued marker is not a terminal delivery result")
		}
	}
	now := s.now().UTC()
	if incomingID != "" {
		if _, err := tx.Exec("UPDATE completion_deliveries SET state=?,completed_at=?,target_run_id=? WHERE id=?", string(domain.DeliveryCompleted), fmtTime(now), r.ID, incomingID); err != nil {
			return err
		}
		if err := recordJoinResult(tx, incoming, *r, now); err != nil {
			return err
		}
	}
	if liveTask && (r.Outcome == domain.OutcomeSuccess || r.Outcome == domain.OutcomeFailure) {
		rows, err := tx.Stmt(s.automationMatches).Query(r.TaskID, string(r.Outcome), string(domain.CompletionOnAny), r.TaskID, string(r.Outcome))
		if err != nil {
			return err
		}
		type match struct{ id, target string }
		var joinIDs []string
		var matches []match
		for rows.Next() {
			var m match
			var joinMatch bool
			if err := rows.Scan(&m.id, &m.target, &joinMatch); err != nil {
				_ = rows.Close()
				return err
			}
			if joinMatch {
				joinIDs = append(joinIDs, m.id)
			} else {
				matches = append(matches, m)
			}
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		if err := rows.Close(); err != nil {
			return err
		}
		if len(matches) > 0 {
			deliveryInsert := tx.Stmt(s.chainDeliveryInsert)
			for _, m := range matches {
				if _, err := deliveryInsert.Exec(newID(), m.id, r.TaskID, m.target, r.ID, string(domain.DeliveryPending), fmtTime(now)); err != nil {
					return err
				}
			}
		}
		if len(joinIDs) > 0 {
			if err := s.createJoinInstances(tx, *r, now, joinIDs); err != nil {
				return err
			}
		}
	}
	if liveTask {
		if err := s.createRunNotificationDeliveries(tx, *r); err != nil {
			return fmt.Errorf("store: create notification deliveries: %w", err)
		}
	}
	if incoming.Role != domain.DeliveryRoleChain && incoming.JoinInstanceID != "" {
		if err := disableInstanceDestinations(tx, incoming.JoinInstanceID); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit run and deliveries: %w", err)
	}
	return nil
}

// insertObservedRun avoids loading output or checking task existence for fresh successful inserts.
func (s *Store) insertObservedRun(tx *sql.Tx, r domain.Run, incoming domain.CompletionDelivery) (live, duplicate bool, result error) {
	if _, err := tx.Stmt(s.runInsert).Exec(runValues(r)...); err != nil {
		previous, readErr := scanRun(tx.QueryRow(runSelect+" WHERE id=?", r.ID))
		if readErr == nil {
			if !sameCommittedRun(previous, r) {
				return false, false, fmt.Errorf("store: conflicting run identity %s", r.ID)
			}
			return true, true, nil
		}
		if !errors.Is(readErr, ErrNotFound) {
			return false, false, readErr
		}
		// Task deletion keeps late owned evidence in immutable join snapshots instead of recreating deleted task history.
		var exists bool
		if queryErr := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id=?)", r.TaskID).Scan(&exists); queryErr != nil {
			return false, false, queryErr
		}
		if !exists {
			if incoming.ID != "" && incoming.Role != domain.DeliveryRoleChain {
				return false, false, nil
			}
			return false, false, ErrNotFound
		}
		return false, false, fmt.Errorf("store: record run: %w", err)
	}
	return true, false, nil
}
