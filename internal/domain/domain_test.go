package domain

import "testing"

func TestSchedulePolicyEffectiveDefaults(t *testing.T) {
	got := (SchedulePolicy{}).Effective()
	if got.TimeBasis != TimeBasisWallClock || got.DSTGap != DSTGapNextValid ||
		got.DSTOverlap != DSTOverlapFirst || got.MissingDate != MissingDateSkip {
		t.Fatalf("effective defaults = %#v", got)
	}

	explicit := SchedulePolicy{
		TimeBasis: TimeBasisElapsed, DSTGap: DSTGapSkip,
		DSTOverlap: DSTOverlapBoth, MissingDate: MissingDateLastValid,
	}
	if got := explicit.Effective(); got != explicit {
		t.Fatalf("explicit policy changed: got %#v want %#v", got, explicit)
	}
}

func TestTaskSchedulePolicy(t *testing.T) {
	task := Task{
		TimeBasis: TimeBasisUTC, DSTGapPolicy: DSTGapSkip,
		DSTOverlapPolicy: DSTOverlapLast, MissingDatePolicy: MissingDateNextValid,
	}
	want := SchedulePolicy{
		TimeBasis: TimeBasisUTC, DSTGap: DSTGapSkip,
		DSTOverlap: DSTOverlapLast, MissingDate: MissingDateNextValid,
	}
	if got := task.SchedulePolicy(); got != want {
		t.Fatalf("SchedulePolicy = %#v, want %#v", got, want)
	}
}
