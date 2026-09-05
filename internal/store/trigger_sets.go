package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

const maxTriggerSetMembers = 99

// ErrInvalidTriggerSet reports invalid set identity, name, target, or count.
var ErrInvalidTriggerSet = errors.New("store: invalid trigger set")

// CreateTriggerSet atomically creates a set and its ordered trigger members.
func (s *Store) CreateTriggerSet(set *domain.TriggerSet, count int, enabled bool) error {
	set.Name = strings.TrimSpace(set.Name)
	if set.Name == "" || set.TargetTaskID == "" || count < 1 || count > maxTriggerSetMembers {
		return ErrInvalidTriggerSet
	}
	if set.ID == "" {
		set.ID = newID()
	}
	now := time.Now().UTC()
	if set.CreatedAt.IsZero() {
		set.CreatedAt = now
	}
	set.UpdatedAt = now
	members := make([]domain.ExternalTrigger, count)
	for index := range members {
		key, err := GenerateTriggerKey()
		if err != nil {
			return err
		}
		position := index + 1
		members[index] = domain.ExternalTrigger{ID: newID(), Name: fmt.Sprintf("%s %02d", set.Name, position), Key: key, TargetTaskID: set.TargetTaskID, SetID: set.ID, SetName: set.Name, SetPosition: position, Enabled: enabled, CreatedAt: now, UpdatedAt: now}
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin create trigger set: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var one int
	if err := tx.QueryRow(`SELECT 1 FROM tasks WHERE id=?`, set.TargetTaskID).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("store: validate trigger set target: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO external_trigger_sets(id,name,target_task_id,created_at,updated_at) VALUES(?,?,?,?,?)`, set.ID, set.Name, set.TargetTaskID, fmtTime(set.CreatedAt), fmtTime(set.UpdatedAt)); err != nil {
		return fmt.Errorf("store: create trigger set: %w", err)
	}
	for index := range members {
		member := &members[index]
		if _, err := tx.Exec(`INSERT INTO external_triggers(id,name,key,target_task_id,set_id,set_position,enabled,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, member.ID, member.Name, member.Key, member.TargetTaskID, member.SetID, member.SetPosition, boolToInt(member.Enabled), fmtTime(member.CreatedAt), fmtTime(member.UpdatedAt)); err != nil {
			return fmt.Errorf("store: create trigger set member %d: %w", member.SetPosition, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit create trigger set: %w", err)
	}
	set.Members = members
	return nil
}

// GetTriggerSet returns one set with members ordered by permanent position.
func (s *Store) GetTriggerSet(id string) (domain.TriggerSet, error) {
	set, err := scanTriggerSet(s.db.QueryRow(triggerSetSelect+` WHERE s.id=?`, id))
	if err != nil {
		return set, err
	}
	set.Members, err = s.listTriggerSetMembers(id)
	return set, err
}

// ListTriggerSets returns all sets with ordered members.
func (s *Store) ListTriggerSets() ([]domain.TriggerSet, error) {
	rows, err := s.db.Query(triggerSetSelect + ` ORDER BY s.name,s.id`)
	if err != nil {
		return nil, fmt.Errorf("store: list trigger sets: %w", err)
	}
	defer rows.Close()
	var sets []domain.TriggerSet
	for rows.Next() {
		set, err := scanTriggerSet(rows)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("store: close trigger set rows: %w", err)
	}
	for index := range sets {
		sets[index].Members, err = s.listTriggerSetMembers(sets[index].ID)
		if err != nil {
			return nil, err
		}
	}
	return sets, nil
}

// RetargetTriggerSet atomically moves every member to one existing task.
func (s *Store) RetargetTriggerSet(id, targetTaskID string) error {
	if targetTaskID == "" {
		return ErrInvalidTriggerSet
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin retarget trigger set: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var oldTarget string
	if err := tx.QueryRow(`SELECT target_task_id FROM external_trigger_sets WHERE id=?`, id).Scan(&oldTarget); err != nil {
		return triggerSetLookupError(err)
	}
	var one int
	if err := tx.QueryRow(`SELECT 1 FROM tasks WHERE id=?`, targetTaskID).Scan(&one); err != nil {
		return triggerSetLookupError(err)
	}
	now := fmtTime(time.Now().UTC())
	if _, err := tx.Exec(`UPDATE external_trigger_sets SET target_task_id=?,updated_at=? WHERE id=?`, targetTaskID, now, id); err != nil {
		return fmt.Errorf("store: retarget trigger set: %w", err)
	}
	if _, err := tx.Exec(`UPDATE external_triggers SET target_task_id=?,updated_at=? WHERE set_id=?`, targetTaskID, now, id); err != nil {
		return fmt.Errorf("store: retarget trigger set members: %w", err)
	}
	if oldTarget != targetTaskID {
		if err := disableIfNotActivationReady(tx, oldTarget); err != nil {
			return err
		}
	}
	return commitTriggerSetTx(tx, "retarget")
}

// SetTriggerSetEnabled atomically changes every member's enabled state.
func (s *Store) SetTriggerSetEnabled(id string, enabled bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin set trigger set enabled: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID string
	if err := tx.QueryRow(`SELECT target_task_id FROM external_trigger_sets WHERE id=?`, id).Scan(&targetID); err != nil {
		return triggerSetLookupError(err)
	}
	now := fmtTime(time.Now().UTC())
	if _, err := tx.Exec(`UPDATE external_trigger_sets SET updated_at=? WHERE id=?`, now, id); err != nil {
		return fmt.Errorf("store: update trigger set timestamp: %w", err)
	}
	if _, err := tx.Exec(`UPDATE external_triggers SET enabled=?,updated_at=? WHERE set_id=?`, boolToInt(enabled), now, id); err != nil {
		return fmt.Errorf("store: set trigger set members enabled: %w", err)
	}
	if !enabled {
		if err := disableIfNotActivationReady(tx, targetID); err != nil {
			return err
		}
	}
	return commitTriggerSetTx(tx, "set enabled")
}

// RotateTriggerSet atomically replaces every member key and returns new secrets.
func (s *Store) RotateTriggerSet(id string) (domain.TriggerSet, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return domain.TriggerSet{}, fmt.Errorf("store: begin rotate trigger set: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	set, err := scanTriggerSet(tx.QueryRow(triggerSetSelect+` WHERE s.id=?`, id))
	if err != nil {
		return set, err
	}
	set.Members, err = listTriggerSetMembers(tx, id)
	if err != nil {
		return domain.TriggerSet{}, err
	}
	keys := make([]string, len(set.Members))
	for index := range keys {
		keys[index], err = GenerateTriggerKey()
		if err != nil {
			return domain.TriggerSet{}, err
		}
	}
	now := time.Now().UTC()
	for index := range set.Members {
		if _, err := tx.Exec(`UPDATE external_triggers SET key=?,updated_at=? WHERE id=? AND set_id=?`, keys[index], fmtTime(now), set.Members[index].ID, id); err != nil {
			return domain.TriggerSet{}, fmt.Errorf("store: rotate trigger set member %d: %w", set.Members[index].SetPosition, err)
		}
		set.Members[index].Key = keys[index]
		set.Members[index].UpdatedAt = now
	}
	if _, err := tx.Exec(`UPDATE external_trigger_sets SET updated_at=? WHERE id=?`, fmtTime(now), id); err != nil {
		return domain.TriggerSet{}, fmt.Errorf("store: rotate trigger set timestamp: %w", err)
	}
	if err := commitTriggerSetTx(tx, "rotate"); err != nil {
		return domain.TriggerSet{}, err
	}
	set.UpdatedAt = now
	return set, nil
}

// DeleteTriggerSet atomically removes a set and every member.
func (s *Store) DeleteTriggerSet(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin delete trigger set: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID string
	var enabledCount int
	if err := tx.QueryRow(`SELECT s.target_task_id,(SELECT COUNT(*) FROM external_triggers e WHERE e.set_id=s.id AND e.enabled<>0) FROM external_trigger_sets s WHERE s.id=?`, id).Scan(&targetID, &enabledCount); err != nil {
		return triggerSetLookupError(err)
	}
	if _, err := tx.Exec(`DELETE FROM external_trigger_sets WHERE id=?`, id); err != nil {
		return fmt.Errorf("store: delete trigger set: %w", err)
	}
	if enabledCount > 0 {
		if err := disableIfNotActivationReady(tx, targetID); err != nil {
			return err
		}
	}
	return commitTriggerSetTx(tx, "delete")
}

const triggerSetSelect = `SELECT s.id,s.name,s.target_task_id,t.name,s.created_at,s.updated_at FROM external_trigger_sets s JOIN tasks t ON t.id=s.target_task_id`

func scanTriggerSet(sc scanner) (domain.TriggerSet, error) {
	var set domain.TriggerSet
	var created, updated string
	if err := sc.Scan(&set.ID, &set.Name, &set.TargetTaskID, &set.TargetTaskName, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return set, ErrNotFound
		}
		return set, fmt.Errorf("store: scan trigger set: %w", err)
	}
	set.CreatedAt, _ = parseTime(created)
	set.UpdatedAt, _ = parseTime(updated)
	return set, nil
}

func (s *Store) listTriggerSetMembers(id string) ([]domain.ExternalTrigger, error) {
	return listTriggerSetMembers(s.db, id)
}

type rowsQueryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

func listTriggerSetMembers(query rowsQueryer, id string) ([]domain.ExternalTrigger, error) {
	rows, err := query.Query(externalTriggerSelect+` WHERE e.set_id=? ORDER BY e.set_position`, id)
	if err != nil {
		return nil, fmt.Errorf("store: list trigger set members: %w", err)
	}
	defer rows.Close()
	var members []domain.ExternalTrigger
	for rows.Next() {
		member, err := scanExternalTrigger(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func triggerSetLookupError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func commitTriggerSetTx(tx *sql.Tx, action string) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit trigger set %s: %w", action, err)
	}
	return nil
}
