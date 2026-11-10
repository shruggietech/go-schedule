package store

import "fmt"

// TaskHasJoinSource includes current definitions and immutable outstanding owned work.
func (s *Store) TaskHasJoinSource(id string) (bool, error) { return taskHasJoinSource(s.db, id) }
func taskHasJoinSource(q queryer, id string) (bool, error) {
	var exists bool
	if err := q.QueryRow("SELECT "+joinSourceExistsSQL("source.task_id")+" FROM (SELECT ? AS task_id) source", id).Scan(&exists); err != nil {
		return false, fmt.Errorf("store: read parallel join activation sources: %w", err)
	}
	return exists, nil
}
func joinSourceExistsSQL(taskID string) string {
	return "(EXISTS(SELECT 1 FROM parallel_joins j WHERE j.target_task_id=" + taskID + " OR j.recovery_task_id=" + taskID + " OR j.cleanup_task_id=" + taskID + ") OR EXISTS(SELECT 1 FROM parallel_join_branches b WHERE b.task_id=" + taskID + ") OR EXISTS(SELECT 1 FROM parallel_join_invocations v JOIN parallel_join_instances i ON i.id=v.instance_id JOIN completion_deliveries d ON d.id=v.delivery_id WHERE v.task_id=" + taskID + " AND i.state<>'resolved' AND d.state IN ('pending','claimed')) OR EXISTS(SELECT 1 FROM parallel_join_instances i WHERE (i.target_task_id=" + taskID + " AND i.state IN ('waiting','target_pending')) OR (i.recovery_task_id=" + taskID + " AND i.state IN ('waiting','recovery_pending')) OR (i.cleanup_task_id=" + taskID + " AND i.state IN ('waiting','target_pending','recovery_pending','cleanup_pending'))))"
}
