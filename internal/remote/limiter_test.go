package remote

import (
	"testing"
	"time"
)

func TestLimiterBoundsBucketsAndBurst(t *testing.T) {
	set := newLimiterSet(0.0001, 1, 1)
	now := time.Date(2026, time.September, 9, 12, 0, 0, 0, time.UTC)
	set.now = func() time.Time { return now }
	if !set.allow("first") {
		t.Fatal("first request rejected")
	}
	if set.allow("first") {
		t.Fatal("burst was not enforced")
	}
	if set.allow("second") {
		t.Fatal("bounded map admitted another source")
	}
	now = now.Add(set.maxIdle)
	if !set.allow("second") {
		t.Fatal("stale source prevented a new source from being admitted")
	}
	if _, exists := set.items["first"]; exists {
		t.Fatal("stale source was not evicted")
	}
}
