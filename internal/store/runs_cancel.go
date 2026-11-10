package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// GetRunResultContext reads only the identity and observed result needed for run control.
// Command output is deliberately excluded, and cancellation bounds SQL-pool acquisition.
func (s *Store) GetRunResultContext(ctx context.Context, id string) (domain.Run, error) {
	var run domain.Run
	var started, ended sql.NullString
	err := s.db.QueryRowContext(ctx, "SELECT id,task_id,outcome,started_at,ended_at FROM runs WHERE id=?", id).Scan(&run.ID, &run.TaskID, &run.Outcome, &started, &ended)
	if errors.Is(err, sql.ErrNoRows) {
		return run, ErrNotFound
	}
	if err != nil {
		return run, fmt.Errorf("store: read run result: %w", err)
	}
	run.StartedAt, err = parseTimePtr(started)
	if err != nil {
		return domain.Run{}, fmt.Errorf("store: read run start: %w", err)
	}
	run.EndedAt, err = parseTimePtr(ended)
	if err != nil {
		return domain.Run{}, fmt.Errorf("store: read run end: %w", err)
	}
	return run, nil
}
