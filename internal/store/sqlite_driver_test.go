package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// This qualifies the driver's supported prepared path, not a new application
// statement cache. Alternating 0/1 bindings exercises the upstream reprepare fix.
func preparedDriverWorkload(tb testing.TB, mode string) *sql.Stmt {
	tb.Helper()
	path := ":memory:"
	if mode == "wal" {
		path = filepath.Join(tb.TempDir(), "prepared.db")
	}
	st, err := Open(path)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := st.Close(); err != nil {
			tb.Error(err)
		}
	})
	for groupIndex := range 2 {
		group := domain.Group{Name: fmt.Sprintf("group-%d", groupIndex), Enabled: groupIndex == 1}
		if err := st.CreateGroup(&group); err != nil {
			tb.Fatal(err)
		}
		for i := range 32 {
			task := domain.Task{ID: fmt.Sprintf("task-%d-%02d", groupIndex, i), Name: "prepared", Command: "echo", GroupID: group.ID, Timezone: "UTC", State: domain.TaskActive}
			if err := st.CreateTask(&task); err != nil {
				tb.Fatal(err)
			}
		}
	}
	stmt, err := st.db.Prepare(`SELECT t.id FROM tasks t JOIN groups g ON g.id=t.group_id WHERE g.enabled=? ORDER BY t.id`)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			tb.Error(err)
		}
	})
	return stmt
}

func preparedDriverQuery(stmt *sql.Stmt, binding int) (count int, result error) {
	rows, err := stmt.Query(binding)
	if err != nil {
		return 0, err
	}
	defer func() { result = errors.Join(result, rows.Close()) }()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return count, err
		}
		if want := fmt.Sprintf("task-%d-%02d", binding, count); id != want {
			return count, fmt.Errorf("prepared row %q, want %q", id, want)
		}
		count++
	}
	return count, rows.Err()
}

func TestSQLitePreparedBindingsMemoryAndWAL(t *testing.T) {
	for _, mode := range []string{"memory", "wal"} {
		t.Run(mode, func(t *testing.T) {
			stmt := preparedDriverWorkload(t, mode)
			for i := range 64 {
				if count, err := preparedDriverQuery(stmt, i%2); err != nil || count != 32 {
					t.Fatalf("iteration %d: rows=%d error=%v", i, count, err)
				}
			}
		})
	}
}

func BenchmarkSQLitePreparedBindings(b *testing.B) {
	for _, mode := range []string{"memory", "wal"} {
		b.Run(mode, func(b *testing.B) {
			stmt := preparedDriverWorkload(b, mode)
			b.ReportAllocs()
			b.ResetTimer()
			for i := range b.N {
				if count, err := preparedDriverQuery(stmt, i%2); err != nil || count != 32 {
					b.Fatalf("rows=%d error=%v", count, err)
				}
			}
		})
	}
}

// Pool-acquisition cancellation is covered separately. This fixture holds a real
// SQLite read lock while a writer tries to commit in rollback-journal mode.
func TestSQLiteCommitLockCancellation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "commit-lock.db")
	writerDB := driverTestDB(t, path)
	readerDB := driverTestDB(t, path)
	if _, err := writerDB.Exec(`PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=5000; CREATE TABLE owned(value INTEGER); INSERT INTO owned VALUES(1)`); err != nil {
		t.Fatal(err)
	}
	writer, err := writerDB.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := writer.Close(); err != nil && !errors.Is(err, sql.ErrConnDone) {
			t.Error(err)
		}
	})
	reader, err := readerDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := reader.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Error(err)
		}
	})
	var value int
	if err := reader.QueryRow("SELECT value FROM owned").Scan(&value); err != nil || value != 1 {
		t.Fatalf("read lock: value=%d error=%v", value, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := writer.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO owned VALUES(2)"); err != nil {
		t.Fatal(err)
	}
	// All SQL setup and the blocking read lock exist before cancellation is
	// scheduled. This timer exercises real I/O waiting, not scheduler time.
	// Keep the reader locked until Commit returns, including on the old driver.
	started := time.Now()
	cancelFired := make(chan struct{})
	timer := time.AfterFunc(100*time.Millisecond, func() {
		cancel()
		close(cancelFired)
	})
	defer func() {
		if !timer.Stop() {
			<-cancelFired
		}
	}()
	if err := ctx.Err(); err != nil {
		t.Fatalf("fixture context ended before commit: %v", err)
	}
	commitErr := tx.Commit()
	if !errors.Is(commitErr, context.Canceled) {
		t.Errorf("commit lost context cancellation: %v", commitErr)
	}
	if elapsed := time.Since(started); elapsed >= 4*time.Second {
		t.Errorf("cancelled commit waited %v for the 5s SQLite busy timeout", elapsed)
	}
	if err := reader.Rollback(); err != nil {
		t.Fatal(err)
	}
	if err := writer.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM owned").Scan(&value); err != nil || value != 1 {
		t.Fatalf("cancelled commit left a partial transaction: count=%d error=%v", value, err)
	}
	if _, err := writer.ExecContext(context.Background(), "INSERT INTO owned VALUES(3)"); err != nil {
		t.Fatalf("connection unusable after cancelled commit: %v", err)
	}
	if err := readerDB.QueryRow("SELECT COUNT(*) FROM owned").Scan(&value); err != nil || value != 2 {
		t.Fatalf("next write was not independently committed: count=%d error=%v", value, err)
	}
}

func driverTestDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}
