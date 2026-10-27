package store

import (
	"context"
	"database/sql"
)

// HoldConnectionForTest owns the single pool slot until its caller releases it.
func (s *Store) HoldConnectionForTest(ctx context.Context) (*sql.Conn, error) { return s.db.Conn(ctx) }

// PoolStatsForTest exposes counters without adding a production SQL capability.
func (s *Store) PoolStatsForTest() sql.DBStats { return s.db.Stats() }
