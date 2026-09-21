package store

import (
	"database/sql"
	"fmt"
)

// PortableID returns the stable bundle identity for one transferable record.
// It is intentionally separate from the record's daemon-local primary key.
func (s *Store) PortableID(kind, objectID string) (string, error) {
	if kind == "" || objectID == "" {
		return "", fmt.Errorf("store: portable identity requires object kind and id")
	}
	var id string
	err := s.db.QueryRow(`SELECT portable_id FROM portable_identities WHERE object_kind=? AND object_id=?`, kind, objectID).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("store: read portable identity: %w", err)
	}
	id = newID()
	if _, err := s.db.Exec(`INSERT OR IGNORE INTO portable_identities(object_kind,object_id,portable_id,created_at) VALUES(?,?,?,?)`, kind, objectID, id, fmtTime(s.now())); err != nil {
		return "", fmt.Errorf("store: create portable identity: %w", err)
	}
	if err := s.db.QueryRow(`SELECT portable_id FROM portable_identities WHERE object_kind=? AND object_id=?`, kind, objectID).Scan(&id); err != nil {
		return "", fmt.Errorf("store: confirm portable identity: %w", err)
	}
	return id, nil
}

// ObjectIDForPortableID resolves a portable identity only within one object
// kind. It never falls back to a display name.
func (s *Store) ObjectIDForPortableID(kind, portableID string) (string, error) {
	var objectID string
	if err := s.db.QueryRow(`SELECT object_id FROM portable_identities WHERE object_kind=? AND portable_id=?`, kind, portableID).Scan(&objectID); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("store: resolve portable identity: %w", err)
	}
	return objectID, nil
}

// BindPortableID records a reviewed portable identity for a newly created
// target object. Conflicting rebinding is rejected by the database constraint.
func (s *Store) BindPortableID(kind, objectID, portableID string) error {
	if kind == "" || objectID == "" || portableID == "" {
		return fmt.Errorf("store: bind portable identity requires values")
	}
	if _, err := s.db.Exec(`INSERT INTO portable_identities(object_kind,object_id,portable_id,created_at) VALUES(?,?,?,?)`, kind, objectID, portableID, fmtTime(s.now())); err != nil {
		return fmt.Errorf("store: bind portable identity: %w", err)
	}
	return nil
}
