package logbus

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func rec(id string, sev domain.AlertSeverity) domain.LogRecord {
	return domain.LogRecord{ID: id, Severity: sev, Time: time.Unix(0, 0), Message: id}
}

func TestRing_BoundedNewestFirst(t *testing.T) {
	r := NewRing(3)
	for _, id := range []string{"1", "2", "3", "4"} {
		r.Add(rec(id, domain.SeverityInfo))
	}
	got := r.Snapshot("", 0)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3 (bounded)", len(got))
	}
	// Newest first; "1" evicted.
	wantOrder := []string{"4", "3", "2"}
	for i, w := range wantOrder {
		if got[i].ID != w {
			t.Errorf("pos %d = %q, want %q", i, got[i].ID, w)
		}
	}
}

func TestRing_SeverityFilterAndLimit(t *testing.T) {
	r := NewRing(10)
	r.Add(rec("a", domain.SeverityInfo))
	r.Add(rec("b", domain.SeverityError))
	r.Add(rec("c", domain.SeverityWarning))
	r.Add(rec("d", domain.SeverityError))

	errs := r.Snapshot(domain.SeverityError, 0)
	if len(errs) != 2 {
		t.Fatalf("error count = %d, want 2", len(errs))
	}
	for _, e := range errs {
		if e.Severity != domain.SeverityError {
			t.Errorf("got severity %q in error filter", e.Severity)
		}
	}

	limited := r.Snapshot("", 2)
	if len(limited) != 2 {
		t.Fatalf("limit=2 returned %d", len(limited))
	}
	if limited[0].ID != "d" {
		t.Errorf("newest = %q, want d", limited[0].ID)
	}
}

func TestRing_EmptyIsSafe(t *testing.T) {
	r := NewRing(5)
	if got := r.Snapshot("", 0); len(got) != 0 {
		t.Fatalf("empty ring returned %d", len(got))
	}
}
