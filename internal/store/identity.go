package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
)

var ErrIdentityMismatch = errors.New("daemon installation identity does not match")

var installationIDGenerator = generateInstallationID

func generateInstallationID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("store: generate daemon installation identity: %w", err)
	}
	return id.String(), nil
}

func (s *Store) initializeDaemonIdentity() error {
	if _, err := s.DaemonIdentity(); err == nil {
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("store: verify daemon identity: %w", err)
	}
	id, err := installationIDGenerator()
	if err != nil {
		return err
	}
	now := fmtTime(time.Now())
	if _, err := s.db.Exec(`INSERT OR IGNORE INTO daemon_identity(singleton,installation_id,display_name,created_at,updated_at) VALUES(1,?,?,?,?)`, id, domain.DefaultDaemonDisplayName, now, now); err != nil {
		return fmt.Errorf("store: initialize daemon identity: %w", err)
	}
	if _, err := s.DaemonIdentity(); err != nil {
		return fmt.Errorf("store: verify daemon identity: %w", err)
	}
	return nil
}

// DaemonIdentity returns the singleton identity for this logical daemon.
func (s *Store) DaemonIdentity() (domain.DaemonIdentity, error) {
	var identity domain.DaemonIdentity
	var created, updated string
	if err := s.db.QueryRow(`SELECT installation_id,display_name,created_at,updated_at FROM daemon_identity WHERE singleton=1`).Scan(&identity.InstallationID, &identity.DisplayName, &created, &updated); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: read daemon identity: %w", err)
	}
	var err error
	if identity.CreatedAt, err = parseTime(created); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: parse daemon identity creation time: %w", err)
	}
	if identity.UpdatedAt, err = parseTime(updated); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: parse daemon identity update time: %w", err)
	}
	if _, err := uuid.Parse(identity.InstallationID); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: validate daemon installation identity: %w", err)
	}
	normalized, err := domain.NormalizeDaemonDisplayName(identity.DisplayName)
	if err != nil || normalized != identity.DisplayName {
		return domain.DaemonIdentity{}, fmt.Errorf("store: validate daemon display name: %w", domain.ErrInvalidDaemonDisplayName)
	}
	return identity, nil
}

// RenameDaemon changes only the operator-facing display name.
func (s *Store) RenameDaemon(name string) (domain.DaemonIdentity, error) {
	name, err := domain.NormalizeDaemonDisplayName(name)
	if err != nil {
		return domain.DaemonIdentity{}, err
	}
	if _, err := s.db.Exec(`UPDATE daemon_identity SET display_name=?,updated_at=? WHERE singleton=1`, name, fmtTime(time.Now())); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: rename daemon: %w", err)
	}
	return s.DaemonIdentity()
}

// ResetDaemonIdentity atomically replaces the installation ID when confirmation is current.
func (s *Store) ResetDaemonIdentity(confirm string) (domain.DaemonIdentity, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: begin daemon identity reset: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var current string
	if err := tx.QueryRow(`SELECT installation_id FROM daemon_identity WHERE singleton=1`).Scan(&current); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: read daemon identity for reset: %w", err)
	}
	if confirm != current {
		return domain.DaemonIdentity{}, ErrIdentityMismatch
	}
	replacement, err := installationIDGenerator()
	if err != nil {
		return domain.DaemonIdentity{}, err
	}
	for replacement == current {
		replacement, err = installationIDGenerator()
		if err != nil {
			return domain.DaemonIdentity{}, err
		}
	}
	result, err := tx.Exec(`UPDATE daemon_identity SET installation_id=?,updated_at=? WHERE singleton=1 AND installation_id=?`, replacement, fmtTime(time.Now()), current)
	if err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: replace daemon identity: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: confirm daemon identity reset: %w", err)
	}
	if changed != 1 {
		return domain.DaemonIdentity{}, ErrIdentityMismatch
	}
	if err := tx.Commit(); err != nil {
		return domain.DaemonIdentity{}, fmt.Errorf("store: commit daemon identity reset: %w", err)
	}
	return s.DaemonIdentity()
}
