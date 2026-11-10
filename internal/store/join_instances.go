package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// GetParallelJoinInstance returns exact immutable bindings and final native evidence.
func (s *Store) GetParallelJoinInstance(id string) (domain.ParallelJoinInstance, error) {
	return s.GetParallelJoinInstancePage(id, 0, 0, 0)
}

// GetParallelJoinInstancePage returns a bounded exact-instance projection.
func (s *Store) GetParallelJoinInstancePage(id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoinInstance, error) {
	return s.GetParallelJoinInstancePageContext(context.Background(), id, branchOffset, branchLimit, textLimit)
}

// GetParallelJoinInstancePageContext cancels both pool waits and nested observation reads.
func (s *Store) GetParallelJoinInstancePageContext(ctx context.Context, id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoinInstance, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return domain.ParallelJoinInstance{}, err
	}
	defer func() { _ = tx.Rollback() }()
	j, err := readJoinInstance(ctx, tx, id, branchOffset, branchLimit, textLimit)
	if err != nil {
		return j, err
	}
	return j, tx.Commit()
}

// ListParallelJoinInstances returns the requested stable instance and nested branch page.
func (s *Store) ListParallelJoinInstances(filter domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, error) {
	return s.ListParallelJoinInstancesContext(context.Background(), filter)
}

// ListParallelJoinInstancesContext bounds observations and cancels storage waits.
func (s *Store) ListParallelJoinInstancesContext(ctx context.Context, f domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q := "SELECT id FROM parallel_join_instances WHERE 1=1"
	var args []any
	if f.DefinitionID != "" {
		q += " AND definition_id=?"
		args = append(args, f.DefinitionID)
	}
	if f.EntryRunID != "" {
		q += " AND entry_run_id=?"
		args = append(args, f.EntryRunID)
	}
	if f.State != "" {
		q += " AND state=?"
		args = append(args, string(f.State))
	}
	q += " ORDER BY created_at DESC,id DESC"
	q, pageArgs := pageSQL(q, f.Offset, f.Limit)
	args = append(args, pageArgs...)
	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	out := make([]domain.ParallelJoinInstance, 0, len(ids))
	for _, id := range ids {
		j, err := readJoinInstance(ctx, tx, id, f.BranchOffset, f.BranchLimit, f.TextLimit)
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return out, nil
}

type joinStageProjection struct {
	prefix                        string
	deliveryID                    *string
	deliveryState                 *domain.DeliveryState
	runID                         *string
	outcome                       *domain.RunOutcome
	settlement                    *domain.JoinSettlement
	reason                        *string
	started, ended, settled       sql.NullString
	startedAt, endedAt, settledAt **time.Time
}

func readJoinInstance(ctx context.Context, tx *sql.Tx, id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoinInstance, error) {
	var j domain.ParallelJoinInstance
	var created, updated, recoveryJSON string
	var completed sql.NullString
	stages := []joinStageProjection{
		{prefix: "target", deliveryID: &j.TargetDeliveryID, deliveryState: &j.TargetDeliveryState, runID: &j.TargetRunID, outcome: &j.TargetOutcome, settlement: &j.TargetSettlement, reason: &j.TargetReason, startedAt: &j.TargetStartedAt, endedAt: &j.TargetEndedAt, settledAt: &j.TargetSettledAt},
		{prefix: "recovery", deliveryID: &j.RecoveryDeliveryID, deliveryState: &j.RecoveryDeliveryState, runID: &j.RecoveryRunID, outcome: &j.RecoveryOutcome, settlement: &j.RecoverySettlement, reason: &j.RecoveryReason, startedAt: &j.RecoveryStartedAt, endedAt: &j.RecoveryEndedAt, settledAt: &j.RecoverySettledAt},
		{prefix: "cleanup", deliveryID: &j.CleanupDeliveryID, deliveryState: &j.CleanupDeliveryState, runID: &j.CleanupRunID, outcome: &j.CleanupOutcome, settlement: &j.CleanupSettlement, reason: &j.CleanupReason, startedAt: &j.CleanupStartedAt, endedAt: &j.CleanupEndedAt, settledAt: &j.CleanupSettledAt},
	}
	columns := []string{"i.id", "i.definition_id", joinTextSQL("i.definition_name", textLimit), "i.entry_task_id", joinTextSQL("i.entry_task_name", textLimit), "i.entry_run_id", "i.target_task_id", joinTextSQL("i.target_task_name", textLimit), "i.recovery_task_id", joinTextSQL("i.recovery_task_name", textLimit), "i.recovery_on_json", "i.cleanup_task_id", joinTextSQL("i.cleanup_task_name", textLimit), "i.selected_route", "i.state", joinTextSQL("i.reason", textLimit), "i.created_at", "i.updated_at", "i.completed_at", "(SELECT COUNT(*) FROM parallel_join_invocations WHERE instance_id=i.id)"}
	dest := []any{&j.ID, &j.DefinitionID, &j.DefinitionName, &j.EntryTaskID, &j.EntryTaskName, &j.EntryRunID, &j.TargetTaskID, &j.TargetTaskName, &j.RecoveryTaskID, &j.RecoveryTaskName, &recoveryJSON, &j.CleanupTaskID, &j.CleanupTaskName, &j.SelectedRoute, &j.State, &j.Reason, &created, &updated, &completed, &j.BranchCount}
	joins := ""
	for n := range stages {
		v := &stages[n]
		alias := fmt.Sprintf("d%d", n)
		prefix := "i." + v.prefix
		columns = append(columns, prefix+"_delivery_id", "COALESCE("+alias+".state,'')", prefix+"_run_id", prefix+"_outcome", prefix+"_started_at", prefix+"_ended_at", prefix+"_settlement", prefix+"_settled_at", joinTextSQL("COALESCE("+alias+".resolution,'')", textLimit))
		dest = append(dest, v.deliveryID, v.deliveryState, v.runID, v.outcome, &v.started, &v.ended, v.settlement, &v.settled, v.reason)
		joins += " LEFT JOIN completion_deliveries " + alias + " ON " + alias + ".id=" + prefix + "_delivery_id"
	}
	if err := tx.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM parallel_join_instances i"+joins+" WHERE i.id=?", id).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return j, ErrNotFound
		}
		return j, err
	}
	if err := json.Unmarshal([]byte(recoveryJSON), &j.RecoveryOn); err != nil {
		return j, fmt.Errorf("store: decode instance recovery conditions: %w", err)
	}
	var err error
	if j.CreatedAt, err = parseTime(created); err != nil {
		return j, err
	}
	if j.UpdatedAt, err = parseTime(updated); err != nil {
		return j, err
	}
	if j.CompletedAt, err = parseTimePtr(completed); err != nil {
		return j, err
	}
	for n := range stages {
		v := &stages[n]
		if *v.startedAt, err = parseTimePtr(v.started); err != nil {
			return j, err
		}
		if *v.endedAt, err = parseTimePtr(v.ended); err != nil {
			return j, err
		}
		if *v.settledAt, err = parseTimePtr(v.settled); err != nil {
			return j, err
		}
		truncateJoinText(v.reason, v.prefix+"_reason", textLimit, &j.TruncatedFields)
	}
	for _, v := range []struct {
		value *string
		field string
	}{{&j.DefinitionName, "definition_name"}, {&j.EntryTaskName, "entry_task_name"}, {&j.TargetTaskName, "target_task_name"}, {&j.RecoveryTaskName, "recovery_task_name"}, {&j.CleanupTaskName, "cleanup_task_name"}, {&j.Reason, "reason"}} {
		truncateJoinText(v.value, v.field, textLimit, &j.TruncatedFields)
	}
	q, args := pageSQL("SELECT v.id,v.branch_id,"+joinTextSQL("v.name", textLimit)+",v.task_id,"+joinTextSQL("v.task_name", textLimit)+",v.delivery_id,d.state,v.run_id,v.outcome,v.settlement,v.settled_at,"+joinTextSQL("v.reason", textLimit)+",v.started_at,v.ended_at FROM parallel_join_invocations v JOIN completion_deliveries d ON d.id=v.delivery_id WHERE v.instance_id=? ORDER BY v.position", branchOffset, branchLimit)
	args = append([]any{id}, args...)
	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return j, err
	}
	defer rows.Close()
	j.Branches = []domain.ParallelJoinInvocation{}
	for rows.Next() {
		var v domain.ParallelJoinInvocation
		var started, ended, settled sql.NullString
		if err := rows.Scan(&v.ID, &v.BranchID, &v.Name, &v.TaskID, &v.TaskName, &v.DeliveryID, &v.DeliveryState, &v.RunID, &v.Outcome, &v.Settlement, &settled, &v.Reason, &started, &ended); err != nil {
			return j, err
		}
		if v.StartedAt, err = parseTimePtr(started); err != nil {
			return j, err
		}
		if v.EndedAt, err = parseTimePtr(ended); err != nil {
			return j, err
		}
		if v.SettledAt, err = parseTimePtr(settled); err != nil {
			return j, err
		}
		truncateJoinText(&v.Name, "name", textLimit, &v.TruncatedFields)
		truncateJoinText(&v.TaskName, "task_name", textLimit, &v.TruncatedFields)
		truncateJoinText(&v.Reason, "reason", textLimit, &v.TruncatedFields)
		j.Branches = append(j.Branches, v)
	}
	j.BranchOffset = max(branchOffset, 0)
	j.HasMoreBranches = j.BranchOffset+len(j.Branches) < j.BranchCount
	return j, rows.Err()
}

func (s *Store) createJoinInstances(tx *sql.Tx, r domain.Run, now time.Time, ids []string) error {
	for _, id := range ids {
		j, err := readParallelJoin(context.Background(), tx, id, 0, 0, 0)
		if err != nil {
			return err
		}
		conditions, err := json.Marshal(j.RecoveryOn)
		if err != nil {
			return fmt.Errorf("store: encode snapshot recovery conditions: %w", err)
		}
		instanceID := newID()
		res, err := tx.Exec("INSERT OR IGNORE INTO parallel_join_instances(id,definition_id,definition_name,entry_task_id,entry_task_name,entry_run_id,target_task_id,target_task_name,recovery_task_id,recovery_task_name,recovery_on_json,cleanup_task_id,cleanup_task_name,state,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", instanceID, j.ID, j.Name, j.EntryTaskID, j.EntryTaskName, r.ID, j.TargetTaskID, j.TargetTaskName, j.RecoveryTaskID, j.RecoveryTaskName, string(conditions), j.CleanupTaskID, j.CleanupTaskName, string(domain.JoinInstanceWaiting), fmtTime(now), fmtTime(now))
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		for pos, b := range j.Branches {
			invocationID, deliveryID := newID(), newID()
			if _, err := tx.Exec("INSERT INTO parallel_join_invocations(id,instance_id,branch_id,position,name,task_id,task_name,delivery_id) VALUES(?,?,?,?,?,?,?,?)", invocationID, instanceID, b.ID, pos, b.Name, b.TaskID, b.TaskName, deliveryID); err != nil {
				return err
			}
			if _, err := tx.Exec("INSERT INTO completion_deliveries(id,source_task_id,target_task_id,source_run_id,role,join_instance_id,join_branch_id,state,created_at) VALUES(?,?,?,?,?,?,?,?,?)", deliveryID, j.EntryTaskID, b.TaskID, r.ID, string(domain.DeliveryRoleJoinBranch), instanceID, invocationID, string(domain.DeliveryPending), fmtTime(now)); err != nil {
				return err
			}
		}
	}
	return nil
}

func recordJoinResult(tx *sql.Tx, d domain.CompletionDelivery, r domain.Run, now time.Time) error {
	if d.Role == domain.DeliveryRoleChain {
		return nil
	}
	settlement, err := nativeJoinSettlement(r)
	if err != nil {
		return err
	}
	decision, err := readJoinDecision(tx, d.JoinInstanceID)
	if err != nil {
		return err
	}
	if d.Role == domain.DeliveryRoleJoinBranch {
		var name string
		if err := tx.QueryRow("SELECT name FROM parallel_join_invocations WHERE instance_id=? AND id=? AND delivery_id=? AND task_id=?", d.JoinInstanceID, d.JoinBranchID, d.ID, r.TaskID).Scan(&name); err != nil {
			return fmt.Errorf("store: exact join branch binding: %w", err)
		}
		reason := ""
		if settlement != domain.JoinSettlementSuccess {
			reason = "branch " + name + " (" + d.JoinBranchID + ") finished " + string(r.Outcome)
			if r.StartFailed {
				reason += " before native command start"
			}
		}
		if _, err := tx.Exec("UPDATE parallel_join_invocations SET run_id=?,outcome=?,settlement=?,settled_at=?,reason=CASE WHEN ?='' THEN reason ELSE ? END,started_at=?,ended_at=? WHERE id=?", r.ID, string(r.Outcome), string(settlement), fmtTime(now), reason, reason, fmtTimePtr(r.StartedAt), fmtTimePtr(r.EndedAt), d.JoinBranchID); err != nil {
			return err
		}
		if decision.state == domain.JoinInstanceResolved {
			return nil
		}
		if settlement != domain.JoinSettlementSuccess {
			if err := blockLegacyJoin(tx, &decision, reason, now); err != nil {
				return err
			}
		}
	} else {
		prefix, _, err := joinStage(d.Role)
		if err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE parallel_join_instances SET "+prefix+"_run_id=?,"+prefix+"_outcome=?,"+prefix+"_started_at=?,"+prefix+"_ended_at=?,"+prefix+"_settlement=?,"+prefix+"_settled_at=?,updated_at=? WHERE id=?", r.ID, string(r.Outcome), fmtTimePtr(r.StartedAt), fmtTimePtr(r.EndedAt), string(settlement), fmtTime(now), fmtTime(now), d.JoinInstanceID); err != nil {
			return err
		}
		if err := decision.setSettlement(d.Role, settlement); err != nil {
			return err
		}
	}
	return advanceJoinInstance(tx, &decision, now)
}

func invalidateJoinInstances(tx *sql.Tx, definitionID, taskID, reason string, now time.Time) error {
	q := "SELECT id,target_task_id FROM parallel_join_instances WHERE state<>? AND state<>?"
	args := []any{string(domain.JoinInstanceCompleted), string(domain.JoinInstanceResolved)}
	if definitionID != "" {
		q += " AND definition_id=?"
		args = append(args, definitionID)
	}
	if taskID != "" {
		q += " AND (entry_task_id=? OR target_task_id=? OR recovery_task_id=? OR cleanup_task_id=? OR EXISTS(SELECT 1 FROM parallel_join_invocations v WHERE v.instance_id=parallel_join_instances.id AND v.task_id=?))"
		args = append(args, taskID, taskID, taskID, taskID, taskID)
	}
	rows, err := tx.Query(q, args...)
	if err != nil {
		return err
	}
	type affectedInstance struct{ id, target string }
	var instances []affectedInstance
	for rows.Next() {
		var i affectedInstance
		if err := rows.Scan(&i.id, &i.target); err != nil {
			_ = rows.Close()
			return err
		}
		instances = append(instances, i)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, i := range instances {
		if _, err := tx.Exec("UPDATE parallel_join_instances SET state=?,reason=?,updated_at=?,completed_at=? WHERE id=?", string(domain.JoinInstanceResolved), reason, fmtTime(now), fmtTime(now), i.id); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE completion_deliveries SET state=?,resolution=?,completed_at=? WHERE join_instance_id=? AND state IN (?,?)", string(domain.DeliveryResolved), reason, fmtTime(now), i.id, string(domain.DeliveryPending), string(domain.DeliveryClaimed)); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE parallel_join_invocations SET reason=? WHERE instance_id=? AND outcome=''", reason, i.id); err != nil {
			return err
		}
		if err := disableInstanceDestinations(tx, i.id); err != nil {
			return err
		}
	}
	return nil
}

func disableInstanceDestinations(tx *sql.Tx, instanceID string) error {
	rows, err := tx.Query("SELECT task_id FROM parallel_join_invocations WHERE instance_id=? UNION SELECT target_task_id FROM parallel_join_instances WHERE id=? UNION SELECT recovery_task_id FROM parallel_join_instances WHERE id=? AND recovery_task_id<>'' UNION SELECT cleanup_task_id FROM parallel_join_instances WHERE id=? AND cleanup_task_id<>''", instanceID, instanceID, instanceID, instanceID)
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, id := range ids {
		if err := disableIfNotActivationReady(tx, id); err != nil {
			return err
		}
	}
	return nil
}

// GetCompletionDelivery reads the exact durable work identity.
func (s *Store) GetCompletionDelivery(id string) (domain.CompletionDelivery, error) {
	return scanDelivery(s.db.QueryRow(deliverySelect+" WHERE id=?", id))
}

// CompletionDeliveryReady validates ownership and immutable bindings before native admission.
func (s *Store) CompletionDeliveryReady(id string) (bool, string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, "", err
	}
	defer func() { _ = tx.Rollback() }()
	d, err := scanDelivery(tx.QueryRow(deliverySelect+" WHERE id=?", id))
	if errors.Is(err, ErrNotFound) {
		return false, "completion delivery no longer exists", nil
	}
	if err != nil {
		return false, "", err
	}
	ready, reason, err := completionDeliveryReady(tx, d)
	if err != nil {
		return false, "", err
	}
	return ready, reason, tx.Commit()
}
func completionDeliveryReady(tx *sql.Tx, d domain.CompletionDelivery) (bool, string, error) {
	if d.State != domain.DeliveryClaimed {
		return false, "completion delivery is no longer claimed", nil
	}
	if d.Role == domain.DeliveryRoleChain {
		var exists int
		err := tx.QueryRow("SELECT 1 FROM completion_chains WHERE id=?", d.ChainID).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return false, "completion chain no longer exists", nil
		}
		return err == nil, "", err
	}
	var state domain.JoinInstanceState
	var definitionID, entryID, sourceRun string
	err := tx.QueryRow("SELECT state,definition_id,entry_task_id,entry_run_id FROM parallel_join_instances WHERE id=?", d.JoinInstanceID).Scan(&state, &definitionID, &entryID, &sourceRun)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "parallel join instance no longer exists", nil
	}
	if err != nil {
		return false, "", err
	}
	if state == domain.JoinInstanceResolved || state == domain.JoinInstanceCompleted {
		return false, "parallel join instance is " + string(state), nil
	}
	var exists int
	err = tx.QueryRow("SELECT 1 FROM parallel_joins WHERE id=?", definitionID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "parallel join definition no longer exists", nil
	}
	if err != nil {
		return false, "", err
	}
	if d.SourceTaskID != entryID || d.SourceRunID != sourceRun {
		return false, "parallel join source binding does not match", nil
	}
	if d.Role != domain.DeliveryRoleJoinBranch {
		prefix, expected, err := joinStage(d.Role)
		if err != nil {
			return false, "unsupported completion delivery role", err
		}
		var targetID, deliveryID string
		var settlement domain.JoinSettlement
		if err := tx.QueryRow("SELECT "+prefix+"_task_id,"+prefix+"_delivery_id,"+prefix+"_settlement FROM parallel_join_instances WHERE id=?", d.JoinInstanceID).Scan(&targetID, &deliveryID, &settlement); err != nil {
			return false, "", err
		}
		return state == expected && d.TargetTaskID == targetID && d.ID == deliveryID && settlement == domain.JoinSettlementPending, "parallel join " + prefix + " is not eligible for this instance", nil
	}
	if d.Role != domain.DeliveryRoleJoinBranch {
		return false, "unsupported completion delivery role", nil
	}
	var taskID, deliveryID, outcome string
	err = tx.QueryRow("SELECT task_id,delivery_id,outcome FROM parallel_join_invocations WHERE instance_id=? AND id=?", d.JoinInstanceID, d.JoinBranchID).Scan(&taskID, &deliveryID, &outcome)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "parallel join invocation no longer exists", nil
	}
	if err != nil {
		return false, "", err
	}
	return (state == domain.JoinInstanceWaiting || state == domain.JoinInstanceBlocked) && taskID == d.TargetTaskID && deliveryID == d.ID && outcome == "", "parallel join branch is no longer eligible", nil
}

func resolveJoinDelivery(tx *sql.Tx, d domain.CompletionDelivery, reason string, now time.Time) error {
	if d.Role == domain.DeliveryRoleChain {
		return nil
	}
	decision, err := readJoinDecision(tx, d.JoinInstanceID)
	if err != nil {
		return err
	}
	if decision.state == domain.JoinInstanceResolved {
		return nil
	}
	if d.Role == domain.DeliveryRoleJoinBranch {
		var name string
		if err := tx.QueryRow("SELECT name FROM parallel_join_invocations WHERE id=? AND instance_id=?", d.JoinBranchID, d.JoinInstanceID).Scan(&name); err != nil {
			return err
		}
		reason = "branch " + name + " (" + d.JoinBranchID + "): " + reason
		if _, err := tx.Exec("UPDATE parallel_join_invocations SET reason=?,settlement=?,settled_at=? WHERE id=? AND outcome='' AND settlement=''", reason, string(domain.JoinSettlementNotStarted), fmtTime(now), d.JoinBranchID); err != nil {
			return err
		}
		if err := blockLegacyJoin(tx, &decision, reason, now); err != nil {
			return err
		}
	} else {
		prefix, _, err := joinStage(d.Role)
		if err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE parallel_join_instances SET "+prefix+"_settlement=?,"+prefix+"_settled_at=?,updated_at=? WHERE id=?", string(domain.JoinSettlementNotStarted), fmtTime(now), fmtTime(now), d.JoinInstanceID); err != nil {
			return err
		}
		if err := decision.setSettlement(d.Role, domain.JoinSettlementNotStarted); err != nil {
			return err
		}
		if decision.recovery == "" && decision.cleanup == "" {
			if _, err := tx.Exec("UPDATE parallel_join_instances SET state=?,reason=?,updated_at=?,completed_at=? WHERE id=?", string(domain.JoinInstanceResolved), reason, fmtTime(now), fmtTime(now), d.JoinInstanceID); err != nil {
				return err
			}
			decision.state = domain.JoinInstanceResolved
		}
	}
	if err := advanceJoinInstance(tx, &decision, now); err != nil {
		return err
	}
	return disableInstanceDestinations(tx, d.JoinInstanceID)
}

func correlateDeliveryRun(r *domain.Run, d domain.CompletionDelivery) {
	r.DeliveryID = d.ID
	r.SourceTaskID = d.SourceTaskID
	r.SourceRunID = d.SourceRunID
	r.JoinInstanceID = d.JoinInstanceID
	r.JoinBranchID = d.JoinBranchID
	switch d.Role {
	case domain.DeliveryRoleJoinBranch:
		r.JoinRole = "branch"
	case domain.DeliveryRoleJoinTarget:
		r.JoinRole = "target"
	case domain.DeliveryRoleJoinRecovery:
		r.JoinRole = "recovery"
	case domain.DeliveryRoleJoinCleanup:
		r.JoinRole = "cleanup"
	default:
		r.JoinRole = ""
	}
}

func sameCommittedRun(a, b domain.Run) bool {
	a.CancellationStatus = ""
	a.CancellationError = ""
	b.CancellationStatus = ""
	b.CancellationError = ""
	// Stored timestamps are UTC-normalized and do not retain monotonic clock components.
	a.ScheduledFor = a.ScheduledFor.UTC()
	b.ScheduledFor = b.ScheduledFor.UTC()
	if a.StartedAt != nil {
		x := a.StartedAt.UTC()
		a.StartedAt = &x
	}
	if b.StartedAt != nil {
		x := b.StartedAt.UTC()
		b.StartedAt = &x
	}
	if a.EndedAt != nil {
		x := a.EndedAt.UTC()
		a.EndedAt = &x
	}
	if b.EndedAt != nil {
		x := b.EndedAt.UTC()
		b.EndedAt = &x
	}
	return reflect.DeepEqual(a, b)
}
