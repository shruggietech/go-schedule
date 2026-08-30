package schedule

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestTimeBasisAcrossSpringForward(t *testing.T) {
	created := time.Date(2026, time.March, 7, 14, 0, 0, 0, time.UTC)
	sch, err := Parse("every 6 hours starting at 09:00", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	after := time.Date(2026, time.March, 8, 2, 0, 0, 0, time.UTC) // Mar 7 21:00 EST

	wall, ok, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{}, after)
	if err != nil || !ok {
		t.Fatalf("wall next: ok=%v err=%v", ok, err)
	}
	if want := time.Date(2026, time.March, 8, 7, 0, 0, 0, time.UTC); !wall.Equal(want) {
		t.Fatalf("wall = %s, want %s", wall, want)
	}

	elapsedPolicy := domain.SchedulePolicy{TimeBasis: domain.TimeBasisElapsed}
	elapsed, ok, err := NextRunWithPolicy(sch, "America/New_York", elapsedPolicy, after)
	if err != nil || !ok {
		t.Fatalf("elapsed next: ok=%v err=%v", ok, err)
	}
	if want := time.Date(2026, time.March, 8, 8, 0, 0, 0, time.UTC); !elapsed.Equal(want) {
		t.Fatalf("elapsed = %s, want %s", elapsed, want)
	}
}

func TestTimeBasisAcrossFallBack(t *testing.T) {
	created := time.Date(2026, time.October, 31, 13, 0, 0, 0, time.UTC)
	sch, err := Parse("every 6 hours starting at 09:00", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	after := time.Date(2026, time.November, 1, 1, 0, 0, 0, time.UTC) // Oct 31 21:00 EDT
	wall, _, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{}, after)
	if err != nil || !wall.Equal(time.Date(2026, time.November, 1, 8, 0, 0, 0, time.UTC)) {
		t.Fatalf("wall = %s, err=%v", wall, err)
	}
	elapsed, _, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{TimeBasis: domain.TimeBasisElapsed}, after)
	if err != nil || !elapsed.Equal(time.Date(2026, time.November, 1, 7, 0, 0, 0, time.UTC)) {
		t.Fatalf("elapsed = %s, err=%v", elapsed, err)
	}
}

func TestElapsedEpochSurvivesPresentationTimezoneChange(t *testing.T) {
	created := time.Date(2026, time.March, 7, 12, 0, 0, 0, time.UTC)
	sch, err := Parse("every day at 09:00", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	policy := domain.SchedulePolicy{TimeBasis: domain.TimeBasisElapsed}
	if err := PrepareForPolicy(&sch, "America/New_York", policy); err != nil {
		t.Fatal(err)
	}
	if sch.ElapsedEpoch == nil || !sch.ElapsedEpoch.Equal(time.Date(2026, time.March, 7, 14, 0, 0, 0, time.UTC)) {
		t.Fatalf("elapsed epoch = %v", sch.ElapsedEpoch)
	}
	after := time.Date(2026, time.March, 8, 13, 30, 0, 0, time.UTC)
	ny, _, err := NextRunWithPolicy(sch, "America/New_York", policy, after)
	if err != nil {
		t.Fatal(err)
	}
	la, _, err := NextRunWithPolicy(sch, "America/Los_Angeles", policy, after)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, time.March, 8, 14, 0, 0, 0, time.UTC)
	if !ny.Equal(want) || !la.Equal(want) {
		t.Fatalf("next after timezone change: New York=%s Los Angeles=%s want=%s", ny, la, want)
	}
}

func TestUTCBasisKeepsUTCReading(t *testing.T) {
	created := time.Date(2026, time.March, 7, 12, 0, 0, 0, time.UTC)
	sch, err := Parse("every day at 09:00", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	got, ok, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{TimeBasis: domain.TimeBasisUTC}, time.Date(2026, time.March, 8, 8, 0, 0, 0, time.UTC))
	if err != nil || !ok {
		t.Fatalf("next: ok=%v err=%v", ok, err)
	}
	if want := time.Date(2026, time.March, 8, 9, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestWallClockGapPolicies(t *testing.T) {
	created := time.Date(2026, time.March, 7, 12, 0, 0, 0, time.UTC)
	sch, err := Parse("every day at 02:30", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	after := time.Date(2026, time.March, 7, 8, 0, 0, 0, time.UTC)

	next, _, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{DSTGap: domain.DSTGapNextValid}, after)
	if err != nil || !next.Equal(time.Date(2026, time.March, 8, 7, 0, 0, 0, time.UTC)) {
		t.Fatalf("next-valid = %s, err=%v", next, err)
	}
	skipped, _, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{DSTGap: domain.DSTGapSkip}, after)
	if err != nil || !skipped.Equal(time.Date(2026, time.March, 9, 6, 30, 0, 0, time.UTC)) {
		t.Fatalf("skip = %s, err=%v", skipped, err)
	}
}

func TestWallClockOverlapPoliciesAndBetweenFoldCursor(t *testing.T) {
	created := time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC)
	sch, err := Parse("every day at 01:30", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	after := time.Date(2026, time.November, 1, 5, 0, 0, 0, time.UTC)
	cases := []struct {
		policy domain.DSTOverlapPolicy
		want   []string
	}{
		{domain.DSTOverlapFirst, []string{"2026-11-01T05:30:00Z", "2026-11-02T06:30:00Z"}},
		{domain.DSTOverlapBoth, []string{"2026-11-01T05:30:00Z", "2026-11-01T06:30:00Z"}},
		{domain.DSTOverlapLast, []string{"2026-11-01T06:30:00Z", "2026-11-02T06:30:00Z"}},
	}
	for _, tc := range cases {
		t.Run(string(tc.policy), func(t *testing.T) {
			got, err := UpcomingRunsWithPolicy(sch, "America/New_York", domain.SchedulePolicy{DSTOverlap: tc.policy}, after, 2)
			if err != nil {
				t.Fatal(err)
			}
			for i, want := range tc.want {
				if got[i].Format(time.RFC3339) != want {
					t.Errorf("[%d] = %s, want %s", i, got[i].Format(time.RFC3339), want)
				}
			}
		})
	}

	between := time.Date(2026, time.November, 1, 5, 45, 0, 0, time.UTC)
	got, ok, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{DSTOverlap: domain.DSTOverlapBoth}, between)
	if err != nil || !ok || !got.Equal(time.Date(2026, time.November, 1, 6, 30, 0, 0, time.UTC)) {
		t.Fatalf("between folds = %s, ok=%v err=%v", got, ok, err)
	}
}

func TestOverlapBothDenseRuleHasBoundedWork(t *testing.T) {
	created := time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC)
	sch, err := Parse("every second", "America/New_York", created)
	if err != nil {
		t.Fatal(err)
	}
	after := time.Date(2026, time.November, 1, 5, 30, 0, 0, time.UTC)
	allocs := testing.AllocsPerRun(1, func() {
		got, ok, err := NextRunWithPolicy(sch, "America/New_York", domain.SchedulePolicy{DSTOverlap: domain.DSTOverlapBoth}, after)
		if err != nil || !ok || !got.Equal(time.Date(2026, time.November, 1, 6, 0, 0, 0, time.UTC)) {
			t.Fatalf("next = %s ok=%v err=%v", got, ok, err)
		}
	})
	if allocs > 1000 {
		t.Fatalf("dense overlap evaluation allocated %.0f objects, want bounded work", allocs)
	}
}

func TestElapsedRejectsCalendarSelectedRecurrence(t *testing.T) {
	sch, err := Parse("3rd wednesday monthly at 14:00", "UTC", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	err = ValidatePolicy(sch, domain.SchedulePolicy{TimeBasis: domain.TimeBasisElapsed})
	if err == nil {
		t.Fatal("elapsed monthly recurrence should be rejected")
	}
}

func TestMissingDatePathComposesWithTransitionPolicies(t *testing.T) {
	overlap, err := Parse("on the 1st of every month at 01:30", "America/New_York", time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	policy := domain.SchedulePolicy{MissingDate: domain.MissingDateLastValid, DSTOverlap: domain.DSTOverlapBoth}
	runs, err := UpcomingRunsWithPolicy(overlap, "America/New_York", policy, time.Date(2026, time.November, 1, 5, 0, 0, 0, time.UTC), 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 2 || runs[0].Format(time.RFC3339) != "2026-11-01T05:30:00Z" || runs[1].Format(time.RFC3339) != "2026-11-01T06:30:00Z" {
		t.Fatalf("overlap runs = %v", runs)
	}

	gap, err := Parse("on the 8th of every month at 02:30", "America/New_York", time.Date(2026, time.February, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	policy = domain.SchedulePolicy{MissingDate: domain.MissingDateLastValid, DSTGap: domain.DSTGapSkip}
	next, ok, err := NextRunWithPolicy(gap, "America/New_York", policy, time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC))
	if err != nil || !ok || next.Format(time.RFC3339) != "2026-04-08T06:30:00Z" {
		t.Fatalf("gap skip = %s ok=%v err=%v", next, ok, err)
	}
}

func TestNextValidCollisionIsEmittedOnce(t *testing.T) {
	anchor := time.Date(2026, time.March, 7, 0, 0, 0, 0, time.UTC)
	sch := domain.Schedule{
		Kind:  domain.ScheduleRecurring,
		RRULE: "FREQ=DAILY;BYHOUR=2,3;BYMINUTE=0;BYSECOND=0", Anchor: &anchor,
	}
	runs, err := UpcomingRunsWithPolicy(sch, "America/New_York", domain.SchedulePolicy{}, time.Date(2026, time.March, 8, 0, 0, 0, 0, time.UTC), 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"2026-03-08T07:00:00Z", "2026-03-09T06:00:00Z", "2026-03-09T07:00:00Z"}
	for i := range want {
		if runs[i].Format(time.RFC3339) != want[i] {
			t.Fatalf("runs=%v, want=%v", runs, want)
		}
	}
}
