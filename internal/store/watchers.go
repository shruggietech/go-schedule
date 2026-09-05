package store

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	DefaultWatcherDebounce  = 250 * time.Millisecond
	DefaultWatcherStability = 500 * time.Millisecond
	MinWatcherDuration      = 25 * time.Millisecond
	MaxWatcherDuration      = time.Hour
)

// ErrInvalidWatcher reports invalid watcher identity, selection, timing, or target data.
var ErrInvalidWatcher = errors.New("store: invalid filesystem watcher")

// CreateFilesystemWatcher validates and inserts one watcher.
func (s *Store) CreateFilesystemWatcher(w *domain.FilesystemWatcher) error {
	if err := normalizeFilesystemWatcher(w); err != nil {
		return err
	}
	if w.ID == "" {
		w.ID = newID()
	}
	now := time.Now().UTC()
	if w.CreatedAt.IsZero() {
		w.CreatedAt = now
	}
	w.UpdatedAt = now
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin create filesystem watcher: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := requireTask(tx, w.TargetTaskID); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO filesystem_watchers(id,name,kind,path,pattern,recursive,debounce_ns,stability_ns,target_task_id,enabled,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, w.ID, w.Name, string(w.Kind), w.Path, w.Pattern, boolToInt(w.Recursive), int64(w.Debounce), int64(w.Stability), w.TargetTaskID, boolToInt(w.Enabled), fmtTime(w.CreatedAt), fmtTime(w.UpdatedAt))
	if err != nil {
		return fmt.Errorf("store: create filesystem watcher: %w", err)
	}
	return commitWatcherTx(tx, "create")
}

// GetFilesystemWatcher returns one watcher with its current target name.
func (s *Store) GetFilesystemWatcher(id string) (domain.FilesystemWatcher, error) {
	return scanFilesystemWatcher(s.db.QueryRow(filesystemWatcherSelect+` WHERE w.id=?`, id))
}

// ListFilesystemWatchers returns all watchers ordered by name and ID.
func (s *Store) ListFilesystemWatchers() ([]domain.FilesystemWatcher, error) {
	rows, err := s.db.Query(filesystemWatcherSelect + ` ORDER BY w.name,w.id`)
	if err != nil {
		return nil, fmt.Errorf("store: list filesystem watchers: %w", err)
	}
	defer rows.Close()
	var out []domain.FilesystemWatcher
	for rows.Next() {
		watcher, err := scanFilesystemWatcher(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, watcher)
	}
	return out, rows.Err()
}

// UpdateFilesystemWatcher atomically replaces mutable watcher fields.
func (s *Store) UpdateFilesystemWatcher(w *domain.FilesystemWatcher) error {
	if w.ID == "" {
		return ErrInvalidWatcher
	}
	if err := normalizeFilesystemWatcher(w); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin update filesystem watcher: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	existing, err := scanFilesystemWatcher(tx.QueryRow(filesystemWatcherSelect+` WHERE w.id=?`, w.ID))
	if err != nil {
		return err
	}
	if err := requireTask(tx, w.TargetTaskID); err != nil {
		return err
	}
	w.CreatedAt = existing.CreatedAt
	w.UpdatedAt = time.Now().UTC()
	res, err := tx.Exec(`UPDATE filesystem_watchers SET name=?,kind=?,path=?,pattern=?,recursive=?,debounce_ns=?,stability_ns=?,target_task_id=?,enabled=?,updated_at=? WHERE id=?`, w.Name, string(w.Kind), w.Path, w.Pattern, boolToInt(w.Recursive), int64(w.Debounce), int64(w.Stability), w.TargetTaskID, boolToInt(w.Enabled), fmtTime(w.UpdatedAt), w.ID)
	if err := affected(res, err, "update filesystem watcher"); err != nil {
		return err
	}
	if existing.TargetTaskID != w.TargetTaskID || (existing.Enabled && !w.Enabled) {
		if err := disableIfNotActivationReady(tx, existing.TargetTaskID); err != nil {
			return err
		}
	}
	return commitWatcherTx(tx, "update")
}

// SetFilesystemWatcherEnabled changes one watcher's enabled state atomically.
func (s *Store) SetFilesystemWatcherEnabled(id string, enabled bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin set filesystem watcher enabled: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID string
	if err := tx.QueryRow(`SELECT target_task_id FROM filesystem_watchers WHERE id=?`, id).Scan(&targetID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("store: get filesystem watcher target: %w", err)
	}
	res, err := tx.Exec(`UPDATE filesystem_watchers SET enabled=?,updated_at=? WHERE id=?`, boolToInt(enabled), fmtTime(time.Now().UTC()), id)
	if err := affected(res, err, "set filesystem watcher enabled"); err != nil {
		return err
	}
	if !enabled {
		if err := disableIfNotActivationReady(tx, targetID); err != nil {
			return err
		}
	}
	return commitWatcherTx(tx, "set enabled")
}

// DeleteFilesystemWatcher removes one watcher and updates its former target readiness.
func (s *Store) DeleteFilesystemWatcher(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin delete filesystem watcher: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var targetID string
	if err := tx.QueryRow(`SELECT target_task_id FROM filesystem_watchers WHERE id=?`, id).Scan(&targetID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("store: get filesystem watcher target: %w", err)
	}
	res, err := tx.Exec(`DELETE FROM filesystem_watchers WHERE id=?`, id)
	if err := affected(res, err, "delete filesystem watcher"); err != nil {
		return err
	}
	if err := disableIfNotActivationReady(tx, targetID); err != nil {
		return err
	}
	return commitWatcherTx(tx, "delete")
}

// TaskHasEnabledWatcher reports whether a task has an enabled watcher source.
func (s *Store) TaskHasEnabledWatcher(id string) (bool, error) {
	return taskHasEnabledWatcher(s.db, id)
}

func taskHasEnabledWatcher(q queryer, id string) (bool, error) {
	var count int
	if err := q.QueryRow(`SELECT COUNT(*) FROM filesystem_watchers WHERE target_task_id=? AND enabled=1`, id).Scan(&count); err != nil {
		return false, fmt.Errorf("store: count task watcher sources: %w", err)
	}
	return count > 0, nil
}

func normalizeFilesystemWatcher(w *domain.FilesystemWatcher) error {
	w.Name = strings.TrimSpace(w.Name)
	w.Path = strings.TrimSpace(w.Path)
	w.Pattern = strings.TrimSpace(w.Pattern)
	if w.Name == "" || w.Path == "" || w.TargetTaskID == "" {
		return ErrInvalidWatcher
	}
	abs, err := filepath.Abs(w.Path)
	if err != nil {
		return fmt.Errorf("%w: path: %v", ErrInvalidWatcher, err)
	}
	w.Path = filepath.Clean(abs)
	if w.Debounce == 0 {
		w.Debounce = DefaultWatcherDebounce
	}
	if w.Stability == 0 {
		w.Stability = DefaultWatcherStability
	}
	if w.Debounce < MinWatcherDuration || w.Debounce > MaxWatcherDuration || w.Stability < MinWatcherDuration || w.Stability > MaxWatcherDuration {
		return ErrInvalidWatcher
	}
	switch w.Kind {
	case domain.WatcherFile:
		w.Pattern = ""
		w.Recursive = false
	case domain.WatcherDirectory:
		if w.Pattern == "" {
			w.Pattern = "*"
		}
		if strings.ContainsAny(w.Pattern, `/\\`) {
			return ErrInvalidWatcher
		}
		if _, err := filepath.Match(w.Pattern, "candidate"); err != nil {
			return ErrInvalidWatcher
		}
	default:
		return ErrInvalidWatcher
	}
	return nil
}

func requireTask(q queryer, id string) error {
	var one int
	if err := q.QueryRow(`SELECT 1 FROM tasks WHERE id=?`, id).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("store: validate filesystem watcher target: %w", err)
	}
	return nil
}

const filesystemWatcherSelect = `SELECT w.id,w.name,w.kind,w.path,w.pattern,w.recursive,w.debounce_ns,w.stability_ns,w.target_task_id,t.name,w.enabled,w.created_at,w.updated_at FROM filesystem_watchers w JOIN tasks t ON t.id=w.target_task_id`

func scanFilesystemWatcher(sc scanner) (domain.FilesystemWatcher, error) {
	var watcher domain.FilesystemWatcher
	var kind, created, updated string
	var recursive, enabled int
	var debounce, stability int64
	if err := sc.Scan(&watcher.ID, &watcher.Name, &kind, &watcher.Path, &watcher.Pattern, &recursive, &debounce, &stability, &watcher.TargetTaskID, &watcher.TargetTaskName, &enabled, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return watcher, ErrNotFound
		}
		return watcher, fmt.Errorf("store: scan filesystem watcher: %w", err)
	}
	watcher.Kind = domain.WatcherKind(kind)
	watcher.Recursive = recursive != 0
	watcher.Debounce = time.Duration(debounce)
	watcher.Stability = time.Duration(stability)
	watcher.Enabled = enabled != 0
	watcher.CreatedAt, _ = parseTime(created)
	watcher.UpdatedAt, _ = parseTime(updated)
	return watcher, nil
}

func commitWatcherTx(tx *sql.Tx, action string) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit filesystem watcher %s: %w", action, err)
	}
	return nil
}
