package timezone

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestResolveWallTimeSpringPolicies(t *testing.T) {
	ny, err := Resolve("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	next := ResolveWallTime(ny, 2026, time.March, 8, 2, 30, 0, domain.DSTGapNextValid, domain.DSTOverlapFirst)
	assertInstants(t, next, "2026-03-08T07:00:00Z")
	skipped := ResolveWallTime(ny, 2026, time.March, 8, 2, 30, 0, domain.DSTGapSkip, domain.DSTOverlapFirst)
	assertInstants(t, skipped)
}

func TestResolveWallTimeFallPolicies(t *testing.T) {
	ny, err := Resolve("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		policy domain.DSTOverlapPolicy
		want   []string
	}{
		{domain.DSTOverlapFirst, []string{"2026-11-01T05:30:00Z"}},
		{domain.DSTOverlapBoth, []string{"2026-11-01T05:30:00Z", "2026-11-01T06:30:00Z"}},
		{domain.DSTOverlapLast, []string{"2026-11-01T06:30:00Z"}},
	}
	for _, tc := range cases {
		t.Run(string(tc.policy), func(t *testing.T) {
			got := ResolveWallTime(ny, 2026, time.November, 1, 1, 30, 0, domain.DSTGapNextValid, tc.policy)
			assertInstants(t, got, tc.want...)
		})
	}
}

func TestResolveWallTimeLordHoweHalfHourTransition(t *testing.T) {
	loc, err := Resolve("Australia/Lord_Howe")
	if err != nil {
		t.Fatal(err)
	}
	got := ResolveWallTime(loc, 2026, time.October, 4, 2, 15, 0, domain.DSTGapNextValid, domain.DSTOverlapFirst)
	assertInstants(t, got, "2026-10-03T15:30:00Z")
}

func assertInstants(t *testing.T, got []time.Time, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i].UTC().Format(time.RFC3339) != want[i] {
			t.Errorf("[%d] = %s, want %s", i, got[i].UTC().Format(time.RFC3339), want[i])
		}
	}
}
