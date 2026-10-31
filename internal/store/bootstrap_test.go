package store

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurrentBootstrapRecordsOneMarker(t *testing.T) {
	st := openMem(t)
	var count, version int
	if err := st.db.QueryRow("SELECT COUNT(*), MAX(version) FROM schema_version").Scan(&count, &version); err != nil {
		t.Fatal(err)
	}
	if count != 1 || version != 24 {
		t.Fatalf("fresh marker rows=%d version=%d; want one marker24", count, version)
	}
	var enabled int
	if err := st.db.QueryRow("PRAGMA foreign_keys").Scan(&enabled); err != nil || enabled != 1 {
		t.Fatalf("foreign keys=%d: %v", enabled, err)
	}
}

func TestUnsupportedStoreRefusesWithoutChangingFile(t *testing.T) {
	for _, seed := range []string{
		"CREATE TABLE unrelated(value TEXT); INSERT INTO unrelated VALUES('keep');",
		"CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version VALUES(25);",
		"CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version VALUES(23);",
		"CREATE TABLE schema_version(version INTEGER NOT NULL); INSERT INTO schema_version VALUES(24);",
	} {
		t.Run(seed, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "unsupported.db")
			db, err := sql.Open("sqlite", path)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(seed); err != nil {
				t.Fatal(err)
			}
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			st, err := Open(path)
			if st != nil {
				st.Close()
			}
			if err == nil || !strings.Contains(err.Error(), "reset") {
				t.Fatalf("want explicit reset instruction, got %v", err)
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatal("refused database was modified")
			}
		})
	}
}

func BenchmarkFreshOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		st, err := Open(":memory:")
		if err != nil {
			b.Fatal(err)
		}
		if err := st.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func createCurrentStore(t *testing.T, path string) {
	t.Helper()
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCurrentBootstrapIsAtomic(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	st := &Store{db: db}
	if err := st.initializeCurrentSchema(currentSchema + "invalid SQL;"); err == nil {
		t.Fatal("expected injected initialization failure")
	}
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%'").Scan(&n); err != nil || n != 0 {
		t.Fatalf("partial bootstrap objects=%d: %v", n, err)
	}
	if err := st.initializeCurrentSchema(currentSchema); err != nil {
		t.Fatal(err)
	}
}

func TestCurrentStorePreservesAcceptedMarkerHistory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "current.db")
	createCurrentStore(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	for version := 1; version < 24; version++ {
		if _, err := db.Exec("INSERT INTO schema_version VALUES(?)", version); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec("INSERT INTO groups(id,name,enabled,created_at,updated_at) VALUES('persisted','Persisted',1,'now','now')"); err != nil {
		t.Fatal(err)
	}
	db.Close()
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	var markers, groups int
	if err := st.db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&markers); err != nil || markers != 24 {
		t.Fatalf("marker history changed: %d %v", markers, err)
	}
	if err := st.db.QueryRow("SELECT COUNT(*) FROM groups WHERE id='persisted'").Scan(&groups); err != nil || groups != 1 {
		t.Fatalf("persisted row changed: %d %v", groups, err)
	}
}
