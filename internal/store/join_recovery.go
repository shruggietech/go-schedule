package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func normalizeJoinRecovery(j *domain.ParallelJoin) error {
	j.RecoveryTaskID = strings.TrimSpace(j.RecoveryTaskID)
	j.CleanupTaskID = strings.TrimSpace(j.CleanupTaskID)
	if j.RecoveryTaskID == "" && len(j.RecoveryOn) > 0 {
		return fmt.Errorf("%w: recovery_on requires a recovery task", ErrInvalidJoin)
	}
	if j.RecoveryOn == nil {
		j.RecoveryOn = []domain.JoinRecoveryCondition{}
		if j.RecoveryTaskID != "" {
			j.RecoveryOn = append(j.RecoveryOn, domain.JoinRecoveryOnFailure)
		}
	}
	selected := map[domain.JoinRecoveryCondition]bool{}
	for _, c := range j.RecoveryOn {
		switch c {
		case domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnCanceled, domain.JoinRecoveryOnSkipped, domain.JoinRecoveryOnNotStarted:
			selected[c] = true
		default:
			return fmt.Errorf("%w: unknown recovery condition %q; choose failure, canceled, skipped or not_started", ErrInvalidJoin, c)
		}
	}
	j.RecoveryOn = []domain.JoinRecoveryCondition{}
	for _, c := range []domain.JoinRecoveryCondition{domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnCanceled, domain.JoinRecoveryOnSkipped, domain.JoinRecoveryOnNotStarted} {
		if selected[c] {
			j.RecoveryOn = append(j.RecoveryOn, c)
		}
	}
	return nil
}

// nativeJoinSettlement classifies only an observed owned terminal result.
func nativeJoinSettlement(r domain.Run) (domain.JoinSettlement, error) {
	if r.StartFailed {
		return domain.JoinSettlementNotStarted, nil
	}
	switch r.Outcome {
	case domain.OutcomeSuccess:
		return domain.JoinSettlementSuccess, nil
	case domain.OutcomeFailure, domain.OutcomeCancellationFailed:
		return domain.JoinSettlementFailure, nil
	case domain.OutcomeCanceled:
		return domain.JoinSettlementCanceled, nil
	case domain.OutcomeSkipped:
		return domain.JoinSettlementSkipped, nil
	default:
		return domain.JoinSettlementPending, fmt.Errorf("store: outcome %q is not an owned terminal result", r.Outcome)
	}
}

// joinDecisionSnapshot belongs to one result transaction, never to another call or admission check.
type joinDecisionSnapshot struct {
	id                                                                              string
	state                                                                           domain.JoinInstanceState
	route                                                                           domain.JoinSelectedRoute
	entry, entryRun, target, recovery, cleanup, conditions, reason, cleanupDelivery string
	targetSettlement, recoverySettlement, cleanupSettlement                         domain.JoinSettlement
}

func readJoinDecision(tx *sql.Tx, id string) (joinDecisionSnapshot, error) {
	d := joinDecisionSnapshot{id: id}
	err := tx.QueryRow("SELECT state,selected_route,entry_task_id,entry_run_id,target_task_id,recovery_task_id,cleanup_task_id,recovery_on_json,target_settlement,recovery_settlement,cleanup_settlement,reason,cleanup_delivery_id FROM parallel_join_instances WHERE id=?", id).Scan(&d.state, &d.route, &d.entry, &d.entryRun, &d.target, &d.recovery, &d.cleanup, &d.conditions, &d.targetSettlement, &d.recoverySettlement, &d.cleanupSettlement, &d.reason, &d.cleanupDelivery)
	return d, err
}

// blockLegacyJoin preserves the first exact unsuccessful cause without changing optional-route flows.
func blockLegacyJoin(tx *sql.Tx, d *joinDecisionSnapshot, reason string, now time.Time) error {
	if d.recovery != "" || d.cleanup != "" || (d.state != domain.JoinInstanceWaiting && d.state != domain.JoinInstanceBlocked) {
		return nil
	}
	if _, err := tx.Exec("UPDATE parallel_join_instances SET state=?,reason=CASE WHEN state=? THEN reason ELSE ? END,updated_at=? WHERE id=? AND state IN (?,?)", string(domain.JoinInstanceBlocked), string(domain.JoinInstanceBlocked), reason, fmtTime(now), d.id, string(domain.JoinInstanceWaiting), string(domain.JoinInstanceBlocked)); err != nil {
		return err
	}
	if d.state != domain.JoinInstanceBlocked {
		d.reason = reason
	}
	d.state = domain.JoinInstanceBlocked
	return nil
}

func (d *joinDecisionSnapshot) setSettlement(role domain.CompletionDeliveryRole, settlement domain.JoinSettlement) error {
	switch role {
	case domain.DeliveryRoleJoinTarget:
		d.targetSettlement = settlement
	case domain.DeliveryRoleJoinRecovery:
		d.recoverySettlement = settlement
	case domain.DeliveryRoleJoinCleanup:
		d.cleanupSettlement = settlement
	default:
		return fmt.Errorf("store: unsupported parallel settlement role %q", role)
	}
	return nil
}

// advanceJoinInstance uses only the current transaction snapshot and known writes to it.
// Branch settlement is aggregated after the latest write; delivery state never implies termination.
func advanceJoinInstance(tx *sql.Tx, d *joinDecisionSnapshot, now time.Time) error {
	if d.state == domain.JoinInstanceResolved || d.state == domain.JoinInstanceCompleted {
		return nil
	}
	route := d.route
	if route == domain.JoinRoutePending {
		var pending, unsuccessful int
		var matching [4]bool
		query := "SELECT COALESCE(SUM(settlement=''),0),COALESCE(SUM(settlement<>'' AND settlement<>'success'),0)"
		dest := []any{&pending, &unsuccessful}
		if d.recovery != "" {
			query += ",COALESCE(MAX(settlement='failure'),0),COALESCE(MAX(settlement='canceled'),0),COALESCE(MAX(settlement='skipped'),0),COALESCE(MAX(settlement='not_started'),0)"
			dest = append(dest, &matching[0], &matching[1], &matching[2], &matching[3])
		}
		if err := tx.QueryRow(query+" FROM parallel_join_invocations WHERE instance_id=?", d.id).Scan(dest...); err != nil {
			return err
		}
		if pending != 0 {
			return nil
		}
		route = domain.JoinRouteNone
		reason := "no unsuccessful branch condition selected for recovery"
		var role domain.CompletionDeliveryRole
		taskID := ""
		if unsuccessful == 0 {
			route = domain.JoinRouteSuccess
			reason = "all owned branches succeeded"
			role = domain.DeliveryRoleJoinTarget
			taskID = d.target
		} else if d.recovery != "" {
			var causes []domain.JoinRecoveryCondition
			if err := json.Unmarshal([]byte(d.conditions), &causes); err != nil {
				return fmt.Errorf("store: decode instance recovery conditions: %w", err)
			}
			present := map[domain.JoinRecoveryCondition]bool{domain.JoinRecoveryOnFailure: matching[0], domain.JoinRecoveryOnCanceled: matching[1], domain.JoinRecoveryOnSkipped: matching[2], domain.JoinRecoveryOnNotStarted: matching[3]}
			for _, cause := range causes {
				if present[cause] {
					route = domain.JoinRouteRecovery
					reason = "configured unsuccessful branch condition selected recovery"
					role = domain.DeliveryRoleJoinRecovery
					taskID = d.recovery
					break
				}
			}
		}
		if route == domain.JoinRouteNone && d.recovery == "" && d.cleanup == "" && d.reason != "" {
			reason = d.reason
		}
		if _, err := tx.Exec("UPDATE parallel_join_instances SET selected_route=?,reason=?,updated_at=? WHERE id=? AND selected_route=''", string(route), reason, fmtTime(now), d.id); err != nil {
			return err
		}
		if taskID != "" {
			return createJoinStage(tx, d.id, d.entry, d.entryRun, taskID, role, now)
		}
	}
	if route == domain.JoinRouteSuccess && d.targetSettlement == domain.JoinSettlementPending {
		return nil
	}
	if route == domain.JoinRouteRecovery && d.recoverySettlement == domain.JoinSettlementPending {
		return nil
	}
	if d.cleanup != "" {
		if d.cleanupDelivery == "" {
			return createJoinStage(tx, d.id, d.entry, d.entryRun, d.cleanup, domain.DeliveryRoleJoinCleanup, now)
		}
		if d.cleanupSettlement == domain.JoinSettlementPending {
			return nil
		}
	}
	final := domain.JoinInstanceCompleted
	if route == domain.JoinRouteNone && d.cleanup == "" {
		final = domain.JoinInstanceBlocked
	}
	_, err := tx.Exec("UPDATE parallel_join_instances SET state=?,updated_at=?,completed_at=? WHERE id=?", string(final), fmtTime(now), fmtTime(now), d.id)
	return err
}

func createJoinStage(tx *sql.Tx, id, entry, entryRun, taskID string, role domain.CompletionDeliveryRole, now time.Time) error {
	prefix, state, err := joinStage(role)
	if err != nil {
		return err
	}
	delivery := newID()
	result, err := tx.Exec("INSERT OR IGNORE INTO completion_deliveries(id,source_task_id,target_task_id,source_run_id,role,join_instance_id,state,created_at) VALUES(?,?,?,?,?,?,?,?)", delivery, entry, taskID, entryRun, string(role), id, string(domain.DeliveryPending), fmtTime(now))
	if err != nil {
		return err
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if inserted == 0 {
		if err := tx.QueryRow("SELECT id FROM completion_deliveries WHERE join_instance_id=? AND role=?", id, string(role)).Scan(&delivery); err != nil {
			return err
		}
	}
	_, err = tx.Exec("UPDATE parallel_join_instances SET state=?,"+prefix+"_delivery_id=?,updated_at=? WHERE id=?", string(state), delivery, fmtTime(now), id)
	return err
}

// joinStage returns only fixed SQL identifiers, never user-authored text.
func joinStage(role domain.CompletionDeliveryRole) (string, domain.JoinInstanceState, error) {
	switch role {
	case domain.DeliveryRoleJoinTarget:
		return "target", domain.JoinInstanceTargetPending, nil
	case domain.DeliveryRoleJoinRecovery:
		return "recovery", domain.JoinInstanceRecoveryPending, nil
	case domain.DeliveryRoleJoinCleanup:
		return "cleanup", domain.JoinInstanceCleanupPending, nil
	default:
		return "", "", fmt.Errorf("store: unsupported parallel stage %q", role)
	}
}
