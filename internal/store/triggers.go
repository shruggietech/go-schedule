package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// ErrInvalidTrigger reports a trigger with a missing name or target.
var ErrInvalidTrigger = errors.New("store: invalid external trigger")

// ErrTriggerSetMemberTarget reports an individual target change that would
// violate a Trigger Set's shared-target invariant.
var ErrTriggerSetMemberTarget = errors.New("store: set member target changes require set retarget")

// GenerateTriggerKey returns an opaque shell-safe 256-bit credential.
func GenerateTriggerKey() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("store: generate trigger key: %w", err)
	}
	return "gst_" + base64.RawURLEncoding.EncodeToString(raw), nil
}

func (s *Store) CreateExternalTrigger(t *domain.ExternalTrigger) error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" || t.TargetTaskID == "" {
		return ErrInvalidTrigger
	}
	if _, err := s.GetTask(t.TargetTaskID); err != nil {
		return err
	}
	if t.ID == "" {
		t.ID = newID()
	}
	if t.Key == "" {
		key, err := GenerateTriggerKey()
		if err != nil {
			return err
		}
		t.Key = key
	}
	now := time.Now().UTC()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
	_, err := s.db.Exec(`INSERT INTO external_triggers(id,name,key,target_task_id,enabled,created_at,updated_at) VALUES(?,?,?,?,?,?,?)`, t.ID, t.Name, t.Key, t.TargetTaskID, boolToInt(t.Enabled), fmtTime(t.CreatedAt), fmtTime(t.UpdatedAt))
	if err != nil {
		return fmt.Errorf("store: create external trigger: %w", err)
	}
	return nil
}

func (s *Store) GetExternalTrigger(id string) (domain.ExternalTrigger, error) {
	return scanExternalTrigger(s.db.QueryRow(externalTriggerSelect+` WHERE e.id=?`, id))
}

func (s *Store) GetExternalTriggerByKey(key string) (domain.ExternalTrigger, error) {
	return scanExternalTrigger(s.db.QueryRow(externalTriggerSelect+` WHERE e.key=?`, key))
}

func (s *Store) ListExternalTriggers() ([]domain.ExternalTrigger, error) {
	rows, err := s.db.Query(externalTriggerSelect + ` ORDER BY e.name,e.id`)
	if err != nil {
		return nil, fmt.Errorf("store: list external triggers: %w", err)
	}
	defer rows.Close()
	var out []domain.ExternalTrigger
	for rows.Next() {
		item, err := scanExternalTrigger(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) UpdateExternalTrigger(t *domain.ExternalTrigger) error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" || t.TargetTaskID == "" {
		return ErrInvalidTrigger
	}
	existing, err := s.GetExternalTrigger(t.ID)
	if err != nil {
		return err
	}
	if existing.SetID != "" && existing.TargetTaskID != t.TargetTaskID {
		return ErrTriggerSetMemberTarget
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin update external trigger: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var one int
	if err := tx.QueryRow(`SELECT 1 FROM tasks WHERE id=?`, t.TargetTaskID).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	t.UpdatedAt = time.Now().UTC()
	res, err := tx.Exec(`UPDATE external_triggers SET name=?,target_task_id=?,updated_at=? WHERE id=?`, t.Name, t.TargetTaskID, fmtTime(t.UpdatedAt), t.ID)
	if err := affected(res, err, "update external trigger"); err != nil {
		return err
	}
	if existing.TargetTaskID != t.TargetTaskID && existing.Enabled {
		if err := disableIfNotActivationReady(tx, existing.TargetTaskID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) SetExternalTriggerEnabled(id string, enabled bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin set external trigger enabled: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID string
	if err := tx.QueryRow(`SELECT target_task_id FROM external_triggers WHERE id=?`, id).Scan(&targetID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	res, err := tx.Exec(`UPDATE external_triggers SET enabled=?,updated_at=? WHERE id=?`, boolToInt(enabled), fmtTime(time.Now().UTC()), id)
	if err := affected(res, err, "set external trigger enabled"); err != nil {
		return err
	}
	if !enabled {
		if err := disableIfNotActivationReady(tx, targetID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) RotateExternalTrigger(id string) (string, error) {
	key, err := GenerateTriggerKey()
	if err != nil {
		return "", err
	}
	res, err := s.db.Exec(`UPDATE external_triggers SET key=?,updated_at=? WHERE id=?`, key, fmtTime(time.Now().UTC()), id)
	if err := affected(res, err, "rotate external trigger"); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) DeleteExternalTrigger(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin delete external trigger: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID, setID string
	var enabled int
	if err := tx.QueryRow(`SELECT target_task_id,enabled,COALESCE(set_id,'') FROM external_triggers WHERE id=?`, id).Scan(&targetID, &enabled, &setID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	res, err := tx.Exec(`DELETE FROM external_triggers WHERE id=?`, id)
	if err := affected(res, err, "delete external trigger"); err != nil {
		return err
	}
	if enabled != 0 {
		if err := disableIfNotActivationReady(tx, targetID); err != nil {
			return err
		}
	}
	if setID != "" {
		var remaining int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM external_triggers WHERE set_id=?`, setID).Scan(&remaining); err != nil {
			return fmt.Errorf("store: count remaining set members: %w", err)
		}
		if remaining == 0 {
			if _, err := tx.Exec(`DELETE FROM external_trigger_sets WHERE id=?`, setID); err != nil {
				return fmt.Errorf("store: delete empty trigger set: %w", err)
			}
		}
	}
	return tx.Commit()
}

const externalTriggerSelect = `SELECT e.id,e.name,e.key,e.target_task_id,t.name,COALESCE(e.set_id,''),COALESCE(s.name,''),COALESCE(e.set_position,0),e.enabled,e.created_at,e.updated_at FROM external_triggers e JOIN tasks t ON t.id=e.target_task_id LEFT JOIN external_trigger_sets s ON s.id=e.set_id`

func scanExternalTrigger(sc scanner) (domain.ExternalTrigger, error) {
	var t domain.ExternalTrigger
	var enabled int
	var created, updated string
	if err := sc.Scan(&t.ID, &t.Name, &t.Key, &t.TargetTaskID, &t.TargetTaskName, &t.SetID, &t.SetName, &t.SetPosition, &enabled, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNotFound
		}
		return t, fmt.Errorf("store: scan external trigger: %w", err)
	}
	t.Enabled = enabled != 0
	t.CreatedAt, _ = parseTime(created)
	t.UpdatedAt, _ = parseTime(updated)
	return t, nil
}
