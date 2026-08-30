package catchup

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/cron"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func hourly(anchor time.Time) domain.Schedule {
	return domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=HOURLY;INTERVAL=1", Anchor: &anchor}
}

func TestEvaluate_MissedRunTriggersCatchup(t *testing.T) {
	anchor := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)
	sch := hourly(anchor)
	last := time.Date(2026, 6, 19, 9, 0, 0, 0, time.UTC)  // last run at 09:00
	now := time.Date(2026, 6, 19, 12, 30, 0, 0, time.UTC) // 3+ hours later (downtime)

	dec, err := Evaluate(sch, "UTC", last, true, domain.CatchupOne, domain.MissingDateSkip, now)
	if err != nil {
		t.Fatal(err)
	}
	if !dec.ShouldCatchUp {
		t.Fatal("expected catch-up after missing runs during downtime")
	}
	if want := time.Date(2026, 6, 19, 10, 0, 0, 0, time.UTC); !dec.FirstMissed.Equal(want) {
		t.Fatalf("first missed = %v, want %v", dec.FirstMissed, want)
	}
}

func TestEvaluate_NoMissWhenNextIsFuture(t *testing.T) {
	anchor := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)
	sch := hourly(anchor)
	last := time.Date(2026, 6, 19, 9, 0, 0, 0, time.UTC)
	now := time.Date(2026, 6, 19, 9, 30, 0, 0, time.UTC) // before the next (10:00)

	dec, _ := Evaluate(sch, "UTC", last, true, domain.CatchupOne, domain.MissingDateSkip, now)
	if dec.ShouldCatchUp {
		t.Fatal("no catch-up expected when the next run is still in the future")
	}
}

func TestEvaluate_PolicyNone(t *testing.T) {
	anchor := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)
	sch := hourly(anchor)
	last := time.Date(2026, 6, 19, 9, 0, 0, 0, time.UTC)
	now := time.Date(2026, 6, 19, 15, 0, 0, 0, time.UTC) // many missed

	dec, _ := Evaluate(sch, "UTC", last, true, domain.CatchupNone, domain.MissingDateSkip, now)
	if dec.ShouldCatchUp {
		t.Fatal("policy 'none' must never catch up")
	}
}

func TestEvaluate_NoPriorRun(t *testing.T) {
	anchor := time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)
	sch := hourly(anchor)
	now := time.Date(2026, 6, 19, 15, 0, 0, 0, time.UTC)

	dec, _ := Evaluate(sch, "UTC", time.Time{}, false, domain.CatchupOne, domain.MissingDateSkip, now)
	if dec.ShouldCatchUp {
		t.Fatal("a never-run task has nothing to catch up")
	}
}

func TestEvaluate_NearestWeekdayUsesAdjustedOccurrence(t *testing.T) {
	anchor := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	sch := domain.Schedule{
		Kind:   domain.ScheduleRecurring,
		RRULE:  "FREQ=MONTHLY;BYMONTHDAY=15;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		Anchor: &anchor, CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday,
	}
	last := time.Date(2026, 7, 15, 9, 0, 0, 0, time.UTC)
	now := time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
	dec, err := Evaluate(sch, "UTC", last, true, domain.CatchupOne, domain.MissingDateSkip, now)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC) // August 15 is Saturday.
	if !dec.ShouldCatchUp || !dec.FirstMissed.Equal(want) {
		t.Fatalf("decision=%+v, want first missed %v", dec, want)
	}
}

func TestEvaluate_ReplacementScheduleDoesNotCatchUpBeforeAnchor(t *testing.T) {
	anchor := time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
	sch := domain.Schedule{
		Kind:               domain.ScheduleRecurring,
		RRULE:              "FREQ=MONTHLY;BYMONTHDAY=15;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		Anchor:             &anchor,
		CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday,
	}
	last := time.Date(2026, 7, 15, 9, 0, 0, 0, time.UTC)
	now := time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC)

	dec, err := Evaluate(sch, "UTC", last, true, domain.CatchupOne, domain.MissingDateSkip, now)
	if err != nil {
		t.Fatal(err)
	}
	if dec.ShouldCatchUp {
		t.Fatalf("decision=%+v, want no catch-up before replacement schedule anchor %v", dec, anchor)
	}
}

func TestEvaluate_CompositeCronFindsFirstMissedRun(t *testing.T) {
	anchor := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	sch, bad, err := cron.Compile("30 8-17 * * MON,WED,FRI", "UTC", anchor)
	if err != nil || bad.Reason != "" {
		t.Fatalf("compile: refusal=%q err=%v", bad.Reason, err)
	}
	last := time.Date(2026, 8, 28, 16, 30, 0, 0, time.UTC)
	now := time.Date(2026, 8, 31, 10, 0, 0, 0, time.UTC)

	dec, err := Evaluate(sch, "UTC", last, true, domain.CatchupOne, domain.MissingDateSkip, now)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 8, 28, 17, 30, 0, 0, time.UTC)
	if !dec.ShouldCatchUp || !dec.FirstMissed.Equal(want) {
		t.Fatalf("decision=%+v, want first missed %v", dec, want)
	}
}

func TestEvaluate_SecondsCronFindsFirstMissedRun(t *testing.T) {
	anchor := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	sch, bad, err := cron.Compile("15,45 * * * * *", "UTC", anchor)
	if err != nil || bad.Reason != "" {
		t.Fatalf("compile: refusal=%q err=%v", bad.Reason, err)
	}
	last := time.Date(2026, 8, 28, 12, 0, 15, 0, time.UTC)
	now := time.Date(2026, 8, 28, 12, 1, 10, 0, time.UTC)

	dec, err := Evaluate(sch, "UTC", last, true, domain.CatchupOne, domain.MissingDateSkip, now)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 8, 28, 12, 0, 45, 0, time.UTC)
	if !dec.ShouldCatchUp || !dec.FirstMissed.Equal(want) {
		t.Fatalf("decision=%+v, want first missed %v", dec, want)
	}
}

func TestEvaluateWithPolicyOverlapBothFindsSecondFold(t *testing.T) {
	anchor := time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC)
	sch, err := schedule.Parse("every day at 01:30", "America/New_York", anchor)
	if err != nil {
		t.Fatal(err)
	}
	last := time.Date(2026, time.November, 1, 5, 30, 0, 0, time.UTC)
	now := time.Date(2026, time.November, 1, 7, 0, 0, 0, time.UTC)
	dec, err := EvaluateWithPolicy(sch, "America/New_York", last, true, domain.CatchupOne,
		domain.SchedulePolicy{DSTOverlap: domain.DSTOverlapBoth}, now)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, time.November, 1, 6, 30, 0, 0, time.UTC)
	if !dec.ShouldCatchUp || !dec.FirstMissed.Equal(want) {
		t.Fatalf("decision=%+v, want one catch-up for %s", dec, want)
	}
}
