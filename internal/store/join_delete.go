package store

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *Store) deleteTaskJoinBindings(tx *sql.Tx, id string) error {
	rows, err := tx.Query("SELECT j.id FROM parallel_joins j WHERE j.entry_task_id=? OR j.target_task_id=? OR EXISTS(SELECT 1 FROM parallel_join_branches b WHERE b.join_id=j.id AND b.task_id=?)", id, id, id)
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var joinID string
		if err := rows.Scan(&joinID); err != nil {
			_ = rows.Close()
			return err
		}
		ids = append(ids, joinID)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	var definitions []string
	for _, joinID := range ids {
		j, err := readParallelJoin(context.Background(), tx, joinID, 0, 0, 0)
		if err != nil {
			return err
		}
		for _, b := range j.Branches {
			definitions = append(definitions, b.TaskID)
		}
		definitions = append(definitions, j.TargetTaskID)
	}
	reason := fmt.Sprintf("required task %s was deleted", id)
	if err := invalidateJoinInstances(tx, "", id, reason, s.now().UTC()); err != nil {
		return err
	}
	for _, joinID := range ids {
		if err := invalidateJoinInstances(tx, joinID, "", reason, s.now().UTC()); err != nil {
			return err
		}
		if _, err := tx.Exec("DELETE FROM parallel_joins WHERE id=?", joinID); err != nil {
			return err
		}
	}
	for _, destination := range definitions {
		if destination != id {
			if err := disableIfNotActivationReady(tx, destination); err != nil {
				return err
			}
		}
	}
	return nil
}
