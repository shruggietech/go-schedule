package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// ErrInvalidJoin reports an invalid reusable parallel block or branch binding.
var ErrInvalidJoin = errors.New("store: invalid parallel join")

// ErrStaleJoin reports that another mutation changed the reviewed definition.
var ErrStaleJoin = errors.New("store: parallel join changed; refresh before editing or explicitly overwrite")

// CreateParallelJoin validates and atomically persists a reusable parallel block.
func (s *Store) CreateParallelJoin(j *domain.ParallelJoin) error {
	if j.ID == "" {
		j.ID = newID()
	}
	now := s.now().UTC()
	j.CreatedAt = now
	j.UpdatedAt = now
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin create parallel join: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := normalizeParallelJoin(j); err != nil {
		return err
	}
	if err := validateParallelJoinReferences(tx, *j); err != nil {
		return err
	}
	if err := insertParallelJoin(tx, *j); err != nil {
		return err
	}
	if err := validateAutomationGraph(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateParallelJoin explicitly overwrites a definition without a revision precondition.
func (s *Store) UpdateParallelJoin(j *domain.ParallelJoin) error {
	return s.UpdateParallelJoinIfUnchanged(j, nil)
}

// UpdateParallelJoinIfUnchanged applies an optional atomic reviewed-revision precondition.
func (s *Store) UpdateParallelJoinIfUnchanged(j *domain.ParallelJoin, expected *time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin update parallel join: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	existing, err := readParallelJoin(context.Background(), tx, j.ID, 0, 0, 0)
	if err != nil {
		return err
	}
	if expected != nil && !existing.UpdatedAt.Equal(*expected) {
		return ErrStaleJoin
	}
	if err := normalizeParallelJoin(j); err != nil {
		return err
	}
	if err := validateParallelJoinReferences(tx, *j); err != nil {
		return err
	}
	j.CreatedAt = existing.CreatedAt
	j.UpdatedAt = s.now().UTC()
	if !j.UpdatedAt.After(existing.UpdatedAt) {
		j.UpdatedAt = existing.UpdatedAt.Add(time.Nanosecond)
	}
	if _, err := tx.Exec("DELETE FROM parallel_joins WHERE id=?", j.ID); err != nil {
		return fmt.Errorf("store: stage parallel join update: %w", err)
	}
	if err := insertParallelJoin(tx, *j); err != nil {
		return err
	}
	if err := validateAutomationGraph(tx); err != nil {
		return err
	}
	if err := disableJoinDestinations(tx, existing); err != nil {
		return err
	}
	return tx.Commit()
}

func normalizeParallelJoin(j *domain.ParallelJoin) error {
	j.Name = strings.TrimSpace(j.Name)
	if j.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidJoin)
	}
	if len(j.Branches) == 0 {
		return fmt.Errorf("%w: at least one named branch is required", ErrInvalidJoin)
	}
	if j.BranchOffset != 0 || j.HasMoreBranches || j.BranchCount > len(j.Branches) {
		return fmt.Errorf("%w: load the complete branch definition before editing", ErrInvalidJoin)
	}
	names := map[string]bool{}
	ids := map[string]bool{}
	for i := range j.Branches {
		b := &j.Branches[i]
		b.Name = strings.TrimSpace(b.Name)
		if b.ID == "" {
			b.ID = newID()
		}
		if b.Name == "" || names[b.Name] || ids[b.ID] {
			return fmt.Errorf("%w: branch %d requires a unique ID and nonempty unique name", ErrInvalidJoin, i+1)
		}
		names[b.Name] = true
		ids[b.ID] = true
	}
	j.BranchCount = len(j.Branches)
	j.BranchOffset = 0
	j.HasMoreBranches = false
	j.TruncatedFields = nil
	return nil
}

func validateParallelJoinReferences(tx *sql.Tx, j domain.ParallelJoin) error {
	refs := []struct{ name, id string }{{"entry", j.EntryTaskID}, {"target", j.TargetTaskID}}
	for _, b := range j.Branches {
		refs = append(refs, struct{ name, id string }{"branch " + b.Name, b.TaskID})
	}
	for _, ref := range refs {
		var exists int
		err := tx.QueryRow("SELECT 1 FROM tasks WHERE id=?", ref.id).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %s task %q does not exist", ErrInvalidJoin, ref.name, ref.id)
		}
		if err != nil {
			return fmt.Errorf("store: validate parallel join reference: %w", err)
		}
	}
	return nil
}

func insertParallelJoin(tx *sql.Tx, j domain.ParallelJoin) error {
	if _, err := tx.Exec("INSERT INTO parallel_joins(id,name,entry_task_id,target_task_id,created_at,updated_at) VALUES(?,?,?,?,?,?)", j.ID, j.Name, j.EntryTaskID, j.TargetTaskID, fmtTime(j.CreatedAt), fmtTime(j.UpdatedAt)); err != nil {
		return fmt.Errorf("store: save parallel join: %w", err)
	}
	for i, b := range j.Branches {
		if _, err := tx.Exec("INSERT INTO parallel_join_branches(join_id,id,position,name,task_id) VALUES(?,?,?,?,?)", j.ID, b.ID, i, b.Name, b.TaskID); err != nil {
			return fmt.Errorf("store: save parallel join branch: %w", err)
		}
	}
	return nil
}

// ReplaceParallelJoins atomically replaces selected definitions, preserving snapshots.
func (s *Store) ReplaceParallelJoins(joins []domain.ParallelJoin) error {
	return s.ReplaceAutomationRelationships(nil, joins)
}

// ReplaceAutomationRelationships validates the final mixed graph and replaces only submitted identities.
func (s *Store) ReplaceAutomationRelationships(chains []domain.CompletionChain, joins []domain.ParallelJoin) error {
	if len(chains) == 0 && len(joins) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin automation replacement: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var oldJoins []domain.ParallelJoin
	var oldTargets []string
	seen := map[string]bool{}
	for i := range chains {
		c := &chains[i]
		if seen["chain:"+c.ID] {
			return fmt.Errorf("%w: duplicate replacement identity", ErrInvalidJoin)
		}
		seen["chain:"+c.ID] = true
		var created, target string
		if err := tx.QueryRow("SELECT created_at,target_task_id FROM completion_chains WHERE id=?", c.ID).Scan(&created, &target); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				if c.ID == "" {
					c.ID = newID()
				}
				c.CreatedAt = s.now().UTC()
				created = fmtTime(c.CreatedAt)
			} else {
				return err
			}
		}
		c.CreatedAt, err = parseTime(created)
		if err != nil {
			return err
		}
		c.UpdatedAt = s.now().UTC()
		oldTargets = append(oldTargets, target)
	}
	for i := range joins {
		j := &joins[i]
		if seen["join:"+j.ID] {
			return fmt.Errorf("%w: duplicate replacement identity", ErrInvalidJoin)
		}
		seen["join:"+j.ID] = true
		old, err := readParallelJoin(context.Background(), tx, j.ID, 0, 0, 0)
		if err != nil && !errors.Is(err, ErrNotFound) {
			return err
		}
		if errors.Is(err, ErrNotFound) {
			if j.ID == "" {
				j.ID = newID()
			}
			old.CreatedAt = s.now().UTC()
		} else {
			oldJoins = append(oldJoins, old)
		}
		if err := normalizeParallelJoin(j); err != nil {
			return err
		}
		if err := validateParallelJoinReferences(tx, *j); err != nil {
			return err
		}
		j.CreatedAt = old.CreatedAt
		j.UpdatedAt = s.now().UTC()
		if !j.UpdatedAt.After(old.UpdatedAt) {
			j.UpdatedAt = old.UpdatedAt.Add(time.Nanosecond)
		}
	}
	for _, c := range chains {
		if _, err := tx.Exec("DELETE FROM completion_chains WHERE id=?", c.ID); err != nil {
			return err
		}
	}
	for _, j := range joins {
		if _, err := tx.Exec("DELETE FROM parallel_joins WHERE id=?", j.ID); err != nil {
			return err
		}
	}
	for _, c := range chains {
		if err := validateCompletionChain(tx, &c, ""); err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT INTO completion_chains(id,source_task_id,target_task_id,on_outcome,created_at,updated_at) VALUES(?,?,?,?,?,?)", c.ID, c.SourceTaskID, c.TargetTaskID, string(c.OnOutcome), fmtTime(c.CreatedAt), fmtTime(c.UpdatedAt)); err != nil {
			if isUniqueConstraint(err) {
				return ErrDuplicateChain
			}
			return err
		}
	}
	for _, j := range joins {
		if err := insertParallelJoin(tx, j); err != nil {
			return err
		}
	}
	if err := validateAutomationGraph(tx); err != nil {
		return err
	}
	for _, id := range oldTargets {
		if err := disableIfNotActivationReady(tx, id); err != nil {
			return err
		}
	}
	for _, j := range oldJoins {
		if err := disableJoinDestinations(tx, j); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// DeleteParallelJoin resolves its outstanding decisions and retains immutable history.
func (s *Store) DeleteParallelJoin(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	j, err := readParallelJoin(context.Background(), tx, id, 0, 0, 0)
	if err != nil {
		return err
	}
	if err := invalidateJoinInstances(tx, id, "", "parallel join definition was deleted", s.now().UTC()); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM parallel_joins WHERE id=?", id); err != nil {
		return err
	}
	if err := disableJoinDestinations(tx, j); err != nil {
		return err
	}
	return tx.Commit()
}

func disableJoinDestinations(tx *sql.Tx, j domain.ParallelJoin) error {
	for _, b := range j.Branches {
		if err := disableIfNotActivationReady(tx, b.TaskID); err != nil {
			return err
		}
	}
	return disableIfNotActivationReady(tx, j.TargetTaskID)
}

// GetParallelJoin returns a complete reusable definition.
func (s *Store) GetParallelJoin(id string) (domain.ParallelJoin, error) {
	return s.GetParallelJoinPage(id, 0, 0, 0)
}

// GetParallelJoinPage returns a bounded branch projection with explicit completeness metadata.
func (s *Store) GetParallelJoinPage(id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoin, error) {
	return s.GetParallelJoinPageContext(context.Background(), id, branchOffset, branchLimit, textLimit)
}

// GetParallelJoinPageContext makes SQL pool acquisition and nested reads cancelable.
func (s *Store) GetParallelJoinPageContext(ctx context.Context, id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoin, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return domain.ParallelJoin{}, err
	}
	defer func() { _ = tx.Rollback() }()
	j, err := readParallelJoin(ctx, tx, id, branchOffset, branchLimit, textLimit)
	if err != nil {
		return j, err
	}
	return j, tx.Commit()
}

// ListParallelJoins returns complete definitions for internal execution and portable export.
func (s *Store) ListParallelJoins() ([]domain.ParallelJoin, error) {
	return s.ListParallelJoinsPage(0, 0, 0, 0, 0)
}

// ListParallelJoinsPage returns a stable bounded definition page.
func (s *Store) ListParallelJoinsPage(offset, limit, textLimit, branchOffset, branchLimit int) ([]domain.ParallelJoin, error) {
	return s.ListParallelJoinsPageContext(context.Background(), offset, limit, textLimit, branchOffset, branchLimit)
}

// ListParallelJoinsPageContext bounds reads and respects cancellation while waiting for storage.
func (s *Store) ListParallelJoinsPageContext(ctx context.Context, offset, limit, textLimit, branchOffset, branchLimit int) ([]domain.ParallelJoin, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q, args := pageSQL("SELECT id FROM parallel_joins ORDER BY created_at,id", offset, limit)
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
	out := make([]domain.ParallelJoin, 0, len(ids))
	for _, id := range ids {
		j, err := readParallelJoin(ctx, tx, id, branchOffset, branchLimit, textLimit)
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

func readParallelJoin(ctx context.Context, tx *sql.Tx, id string, branchOffset, branchLimit, textLimit int) (domain.ParallelJoin, error) {
	var j domain.ParallelJoin
	var created, updated string
	q := "SELECT j.id," + joinTextSQL("j.name", textLimit) + ",j.entry_task_id," + joinTextSQL("e.name", textLimit) + ",j.target_task_id," + joinTextSQL("t.name", textLimit) + ",j.created_at,j.updated_at,(SELECT COUNT(*) FROM parallel_join_branches WHERE join_id=j.id) FROM parallel_joins j JOIN tasks e ON e.id=j.entry_task_id JOIN tasks t ON t.id=j.target_task_id WHERE j.id=?"
	if err := tx.QueryRowContext(ctx, q, id).Scan(&j.ID, &j.Name, &j.EntryTaskID, &j.EntryTaskName, &j.TargetTaskID, &j.TargetTaskName, &created, &updated, &j.BranchCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return j, ErrNotFound
		}
		return j, fmt.Errorf("store: read parallel join: %w", err)
	}
	var err error
	j.CreatedAt, err = parseTime(created)
	if err != nil {
		return j, err
	}
	j.UpdatedAt, err = parseTime(updated)
	if err != nil {
		return j, err
	}
	truncateJoinText(&j.Name, "name", textLimit, &j.TruncatedFields)
	truncateJoinText(&j.EntryTaskName, "entry_task_name", textLimit, &j.TruncatedFields)
	truncateJoinText(&j.TargetTaskName, "target_task_name", textLimit, &j.TruncatedFields)
	q, args := pageSQL("SELECT b.id,"+joinTextSQL("b.name", textLimit)+",b.task_id,"+joinTextSQL("t.name", textLimit)+" FROM parallel_join_branches b JOIN tasks t ON t.id=b.task_id WHERE b.join_id=? ORDER BY b.position", branchOffset, branchLimit)
	args = append([]any{id}, args...)
	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return j, err
	}
	defer rows.Close()
	j.Branches = []domain.ParallelJoinBranch{}
	for rows.Next() {
		var b domain.ParallelJoinBranch
		if err := rows.Scan(&b.ID, &b.Name, &b.TaskID, &b.TaskName); err != nil {
			return j, err
		}
		truncateJoinText(&b.Name, "name", textLimit, &b.TruncatedFields)
		truncateJoinText(&b.TaskName, "task_name", textLimit, &b.TruncatedFields)
		j.Branches = append(j.Branches, b)
	}
	j.BranchOffset = max(branchOffset, 0)
	j.HasMoreBranches = j.BranchOffset+len(j.Branches) < j.BranchCount
	return j, rows.Err()
}

func pageSQL(q string, offset, limit int) (string, []any) {
	var args []any
	if limit > 0 {
		q += " LIMIT ?"
		args = append(args, limit)
	} else if offset > 0 {
		q += " LIMIT -1"
	}
	if offset > 0 {
		q += " OFFSET ?"
		args = append(args, offset)
	}
	return q, args
}
func joinTextSQL(column string, limit int) string {
	if limit <= 0 {
		return column
	}
	return fmt.Sprintf("substr(%s,1,%d)", column, limit+1)
}
func truncateJoinText(value *string, field string, limit int, fields *[]string) {
	if limit <= 0 || len(*value) <= limit {
		return
	}
	s := (*value)[:limit]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	*value = s
	*fields = append(*fields, field)
}

func validateAutomationGraph(tx *sql.Tx) error {
	graph, err := automationGraph(tx, "")
	if err != nil {
		return err
	}
	colors := map[string]int{}
	var visit func(string) bool
	visit = func(id string) bool {
		if colors[id] == 1 {
			return true
		}
		if colors[id] == 2 {
			return false
		}
		colors[id] = 1
		for _, next := range graph[id] {
			if visit(next) {
				return true
			}
		}
		colors[id] = 2
		return false
	}
	for id := range graph {
		if visit(id) {
			return ErrChainCycle
		}
	}
	return nil
}
func automationGraph(tx *sql.Tx, excludeChain string) (map[string][]string, error) {
	rows, err := tx.Query("SELECT source_task_id,target_task_id FROM completion_chains WHERE id<>? UNION ALL SELECT j.entry_task_id,b.task_id FROM parallel_joins j JOIN parallel_join_branches b ON b.join_id=j.id UNION ALL SELECT b.task_id,j.target_task_id FROM parallel_joins j JOIN parallel_join_branches b ON b.join_id=j.id", excludeChain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	graph := map[string][]string{}
	for rows.Next() {
		var from, to string
		if err := rows.Scan(&from, &to); err != nil {
			return nil, err
		}
		graph[from] = append(graph[from], to)
	}
	return graph, rows.Err()
}
