// Package logbus turns the daemon's structured logs into a queryable, streamable
// feed for the GUI Logs view. A single slog.Handler (Handler) tees every record
// to three sinks: a bounded in-memory ring (served by GET /v1/logs), a rotating
// on-disk JSONL file (the durable troubleshooting record), and the live event
// broker (so open clients update in real time).
package logbus

import (
	"sync"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Ring is a bounded, concurrency-safe buffer of the most recent log records,
// newest-first on read. When full, the oldest record is evicted.
type Ring struct {
	mu   sync.RWMutex
	buf  []domain.LogRecord
	next int // index of the next write (oldest when full)
	size int
	full bool
}

// NewRing creates a ring holding up to size records. size <= 0 defaults to 1000.
func NewRing(size int) *Ring {
	if size <= 0 {
		size = 1000
	}
	return &Ring{buf: make([]domain.LogRecord, size), size: size}
}

// Add appends a record, evicting the oldest if the ring is full.
func (r *Ring) Add(rec domain.LogRecord) {
	r.mu.Lock()
	r.buf[r.next] = rec
	r.next = (r.next + 1) % r.size
	if r.next == 0 {
		r.full = true
	}
	r.mu.Unlock()
}

// Snapshot returns recent records newest-first. If severity is non-empty, only
// records of that severity are returned. limit <= 0 means "no limit".
func (r *Ring) Snapshot(severity domain.AlertSeverity, limit int) []domain.LogRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.LogRecord, 0, r.count())
	// Walk from newest to oldest.
	for i := 0; i < r.count(); i++ {
		idx := (r.next - 1 - i + r.size) % r.size
		rec := r.buf[idx]
		if severity != "" && rec.Severity != severity {
			continue
		}
		out = append(out, rec)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

// count returns the number of stored records (caller holds the lock).
func (r *Ring) count() int {
	if r.full {
		return r.size
	}
	return r.next
}
