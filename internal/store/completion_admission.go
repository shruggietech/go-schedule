package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// GetCompletionAdmission reads exact delivery ownership and current task/group eligibility in one snapshot.
// Every admission boundary performs this query again, including queue release and worker acquisition.
func (s *Store) GetCompletionAdmission(id string) (domain.CompletionDelivery, domain.Task, bool, string, error) {
	var d domain.CompletionDelivery
	var created string
	var claimed, completed, targetRun sql.NullString
	var owned, groupsEnabled bool
	var instanceState, instanceReason string
	extra := []any{&d.ID, &d.ChainID, &d.SourceTaskID, &d.TargetTaskID, &d.SourceRunID, &d.State, &d.Attempts, &created, &claimed, &completed, &targetRun, &d.Resolution, &d.Role, &d.JoinInstanceID, &d.JoinBranchID, &owned, &groupsEnabled, &instanceState, &instanceReason}
	task, err := scanTask(extendedAdmissionScanner{row: s.completionAdmission.QueryRow(id), extra: extra})
	if errors.Is(err, ErrNotFound) {
		return d, domain.Task{}, false, "completion delivery no longer exists", nil
	}
	if err != nil {
		return d, domain.Task{}, false, "", fmt.Errorf("store: read completion admission: %w", err)
	}
	d.CreatedAt, err = parseTime(created)
	if err != nil {
		return d, task, false, "", err
	}
	d.ClaimedAt, err = parseTimePtr(claimed)
	if err != nil {
		return d, task, false, "", err
	}
	d.CompletedAt, err = parseTimePtr(completed)
	if err != nil {
		return d, task, false, "", err
	}
	d.TargetRunID = targetRun.String
	if d.State != domain.DeliveryClaimed {
		reason := d.Resolution
		if reason == "" {
			reason = "completion delivery is no longer claimed"
		}
		return d, task, false, reason, nil
	}
	if !owned {
		reason := "completion chain no longer exists"
		if d.Role != domain.DeliveryRoleChain {
			reason = "parallel join definition or exact invocation is no longer eligible"
			if instanceState == string(domain.JoinInstanceResolved) || instanceState == string(domain.JoinInstanceBlocked) {
				reason = "parallel join instance is " + instanceState
				if instanceReason != "" {
					reason += ": " + instanceReason
				}
			}
		}
		return d, task, false, reason, nil
	}
	switch {
	case task.ID == "":
		return d, task, false, "target task no longer exists", nil
	case !task.Enabled || task.State != domain.TaskActive:
		return d, task, false, "target task is disabled or inactive", nil
	case strings.TrimSpace(task.Command) == "":
		return d, task, false, "target task has no command", nil
	case !groupsEnabled:
		return d, task, false, "target group chain is disabled or unavailable", nil
	default:
		return d, task, true, "", nil
	}
}

type extendedAdmissionScanner struct {
	row   scanner
	extra []any
}

func (r extendedAdmissionScanner) Scan(dest ...any) error {
	return r.row.Scan(append(dest, r.extra...)...)
}

// The recursive projection terminates for cycles. An enabled ancestry must reach a root.
const completionAdmissionSQL = `SELECT
 COALESCE(t.id,''),COALESCE(t.name,''),t.group_id,COALESCE(t.command,''),COALESCE(t.args_json,'[]'),COALESCE(t.working_dir,''),COALESCE(t.env_json,'{}'),COALESCE(t.stdin,''),COALESCE(t.run_as,''),COALESCE(t.enabled,0),COALESCE(t.timezone,''),t.schedule_id,COALESCE(t.overlap_policy,''),COALESCE(t.catchup_policy,''),COALESCE(t.missing_date_policy,''),COALESCE(t.time_basis,''),COALESCE(t.dst_gap_policy,''),COALESCE(t.dst_overlap_policy,''),COALESCE(t.state,''),COALESCE(t.created_at,''),COALESCE(t.updated_at,''),
 d.id,d.chain_id,d.source_task_id,d.target_task_id,d.source_run_id,d.state,d.attempts,d.created_at,d.claimed_at,d.completed_at,d.target_run_id,d.resolution,d.role,d.join_instance_id,d.join_branch_id,
 CASE d.role
 WHEN 'chain' THEN EXISTS(SELECT 1 FROM completion_chains c WHERE c.id=d.chain_id)
 WHEN 'join_branch' THEN COALESCE(j.id IS NOT NULL AND i.state IN ('waiting','blocked') AND d.source_task_id=i.entry_task_id AND d.source_run_id=i.entry_run_id AND v.task_id=d.target_task_id AND v.delivery_id=d.id AND v.outcome='' AND v.settlement='',0)
 WHEN 'join_target' THEN COALESCE(j.id IS NOT NULL AND i.state='target_pending' AND d.source_task_id=i.entry_task_id AND d.source_run_id=i.entry_run_id AND d.target_task_id=i.target_task_id AND d.id=i.target_delivery_id AND i.target_settlement='',0)
 WHEN 'join_recovery' THEN COALESCE(j.id IS NOT NULL AND i.state='recovery_pending' AND d.source_task_id=i.entry_task_id AND d.source_run_id=i.entry_run_id AND d.target_task_id=i.recovery_task_id AND d.id=i.recovery_delivery_id AND i.recovery_settlement='',0)
 WHEN 'join_cleanup' THEN COALESCE(j.id IS NOT NULL AND i.state='cleanup_pending' AND d.source_task_id=i.entry_task_id AND d.source_run_id=i.entry_run_id AND d.target_task_id=i.cleanup_task_id AND d.id=i.cleanup_delivery_id AND i.cleanup_settlement='',0)
 ELSE 0 END,
 CASE WHEN COALESCE(t.group_id,'')='' THEN 1 ELSE
 (WITH RECURSIVE ancestors(id,parent_id,enabled) AS (
 SELECT id,parent_id,enabled FROM groups WHERE id=t.group_id
 UNION SELECT g.id,g.parent_id,g.enabled FROM groups g JOIN ancestors a ON g.id=a.parent_id)
 SELECT EXISTS(SELECT 1 FROM ancestors WHERE COALESCE(parent_id,'')='') AND NOT EXISTS(SELECT 1 FROM ancestors WHERE enabled=0)) END,
 COALESCE(i.state,''),COALESCE(i.reason,'')
 FROM completion_deliveries d
 LEFT JOIN tasks t ON t.id=d.target_task_id
 LEFT JOIN parallel_join_instances i ON i.id=d.join_instance_id
 LEFT JOIN parallel_joins j ON j.id=i.definition_id
 LEFT JOIN parallel_join_invocations v ON v.id=d.join_branch_id AND v.instance_id=d.join_instance_id
 WHERE d.id=?`
