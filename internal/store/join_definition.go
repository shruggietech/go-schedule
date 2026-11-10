package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// GetParallelJoinDefinitionPageContext reads lossless editable intent for Manage authoring.
// It pages branch slots without truncating accepted names or loading readonly task display labels.
func (s *Store) GetParallelJoinDefinitionPageContext(ctx context.Context, id string, branchOffset, branchLimit int) (domain.ParallelJoin, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return domain.ParallelJoin{}, fmt.Errorf("store: begin parallel join definition read: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var j domain.ParallelJoin
	var created, updated string
	err = tx.QueryRowContext(ctx, "SELECT j.id,j.name,j.entry_task_id,j.target_task_id,j.created_at,j.updated_at,(SELECT COUNT(*) FROM parallel_join_branches b WHERE b.join_id=j.id) FROM parallel_joins j WHERE j.id=?", id).Scan(&j.ID, &j.Name, &j.EntryTaskID, &j.TargetTaskID, &created, &updated, &j.BranchCount)
	if errors.Is(err, sql.ErrNoRows) {
		return j, ErrNotFound
	}
	if err != nil {
		return j, fmt.Errorf("store: read parallel join definition: %w", err)
	}
	j.CreatedAt, err = parseTime(created)
	if err != nil {
		return j, fmt.Errorf("store: parse parallel join definition creation time: %w", err)
	}
	j.UpdatedAt, err = parseTime(updated)
	if err != nil {
		return j, fmt.Errorf("store: parse parallel join definition update time: %w", err)
	}
	query, args := pageSQL("SELECT id,name,task_id FROM parallel_join_branches WHERE join_id=? ORDER BY position", branchOffset, branchLimit)
	args = append([]any{id}, args...)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return j, fmt.Errorf("store: read parallel join definition branches: %w", err)
	}
	defer rows.Close()
	j.Branches = []domain.ParallelJoinBranch{}
	for rows.Next() {
		var b domain.ParallelJoinBranch
		if err := rows.Scan(&b.ID, &b.Name, &b.TaskID); err != nil {
			return j, fmt.Errorf("store: scan parallel join definition branch: %w", err)
		}
		j.Branches = append(j.Branches, b)
	}
	if err := rows.Err(); err != nil {
		return j, fmt.Errorf("store: read parallel join definition branches: %w", err)
	}
	if err := rows.Close(); err != nil {
		return j, fmt.Errorf("store: close parallel join definition branches: %w", err)
	}
	j.BranchOffset = max(branchOffset, 0)
	j.HasMoreBranches = j.BranchOffset+len(j.Branches) < j.BranchCount
	if err := tx.Commit(); err != nil {
		return j, fmt.Errorf("store: commit parallel join definition read: %w", err)
	}
	return j, nil
}
