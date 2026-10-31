// Package store provides durable persistence for scheduler entities using an
// embedded SQLite database (pure-Go modernc.org/sqlite driver). All timestamps
// are stored in UTC (RFC 3339). The store is the only writer; clients reach it
// through the daemon's API.
package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite" // registers the "sqlite" database/sql driver
)

// Store wraps a SQLite database connection.
type Store struct {
	db  *sql.DB
	now func() time.Time
}

// Open opens (creating if needed) the current-format database at path and
// enables foreign keys. Use ":memory:" for an in-memory database in tests.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", path, err)
	}
	// modernc.org/sqlite is safe for concurrent use, but a single writer avoids
	// SQLITE_BUSY under our daemon's serialized write path.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("store: pragmas: %w", err)
	}
	s := &Store{db: db, now: time.Now}
	if err := s.initializeCurrentSchema(currentSchema); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("store: journal mode: %w", err)
	}
	if path != ":memory:" {
		if err := os.Chmod(path, 0o600); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("store: protect database permissions: %w", err)
		}
	}
	if err := s.initializeDaemonIdentity(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := s.initializeLocalActor(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the underlying database.
func (s *Store) Close() error { return s.db.Close() }

// PoolStats returns connection-pool metadata without acquiring a SQL connection.
func (s *Store) PoolStats() sql.DBStats { return s.db.Stats() }

const currentSchemaVersion = 24

//go:embed current_schema.sql
var currentSchema string

var schemaObjectName = regexp.MustCompile(`(?i)^CREATE (?:UNIQUE )?(?:TABLE|INDEX)\s+("[^"]+"|[^\s(]+)`)

func schemaDefinitions(schema string) map[string]string {
	definitions := make(map[string]string)
	for _, statement := range strings.Split(schema, ";") {
		normalized := strings.Join(strings.Fields(statement), " ")
		if normalized == "" {
			continue
		}
		match := schemaObjectName.FindStringSubmatch(normalized)
		if len(match) == 2 {
			definitions[strings.Trim(match[1], "\"")] = normalized
		}
	}
	return definitions
}

// initializeCurrentSchema creates only an empty database. Existing state must
// already have the complete current shape; the marker alone is insufficient.
func (s *Store) initializeCurrentSchema(schema string) (result error) {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin initialization: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			result = errors.Join(result, fmt.Errorf("store: rollback initialization: %w", err))
		}
	}()
	var objects int
	if err := tx.QueryRow("SELECT COUNT(*) FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%'").Scan(&objects); err != nil {
		return unsupportedStore(err)
	}
	if objects == 0 {
		if _, err := tx.Exec(schema); err != nil {
			return fmt.Errorf("store: initialize current schema: %w", err)
		}
		if _, err := tx.Exec("INSERT INTO schema_version(version) VALUES(?)", currentSchemaVersion); err != nil {
			return fmt.Errorf("store: record current schema: %w", err)
		}
	} else {
		var version int
		if err := tx.QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version); err != nil {
			return unsupportedStore(err)
		}
		if version != currentSchemaVersion {
			return unsupportedStore(fmt.Errorf("schema marker %d, expected %d", version, currentSchemaVersion))
		}
		expected := schemaDefinitions(currentSchema)
		rows, err := tx.Query("SELECT name,sql FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%'")
		if err != nil {
			return unsupportedStore(err)
		}
		for rows.Next() {
			var name string
			var definition sql.NullString
			if err := rows.Scan(&name, &definition); err != nil {
				rows.Close()
				return unsupportedStore(err)
			}
			want, ok := expected[name]
			if !ok || !definition.Valid || strings.Join(strings.Fields(definition.String), " ") != want {
				rows.Close()
				return unsupportedStore(fmt.Errorf("unexpected schema object %s", name))
			}
			delete(expected, name)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return unsupportedStore(err)
		}
		rows.Close()
		if len(expected) != 0 {
			return unsupportedStore(fmt.Errorf("incomplete current schema"))
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit initialization: %w", err)
	}
	return nil
}

func unsupportedStore(cause error) error {
	return fmt.Errorf("store: unsupported database (%v); use an empty database path or explicitly reset the disposable development database; no conversion was attempted", cause)
}
