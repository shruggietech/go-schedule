package remote

import "testing"

func TestLimiterBoundsBucketsAndBurst(t *testing.T) {
	set := newLimiterSet(0.0001, 1, 1)
	if !set.allow("first") {
		t.Fatal("first request rejected")
	}
	if set.allow("first") {
		t.Fatal("burst was not enforced")
	}
	if set.allow("second") {
		t.Fatal("bounded map admitted another source")
	}
}
