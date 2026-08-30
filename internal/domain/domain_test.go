package domain

import "testing"

func TestStartupEventIdentityAndRunOrigin(t *testing.T) {
	if StartupEventID != "scheduler_startup" {
		t.Fatalf("StartupEventID = %q, want scheduler_startup", StartupEventID)
	}
	if TriggerStartup != "startup" {
		t.Fatalf("TriggerStartup = %q, want startup", TriggerStartup)
	}
	if TriggerStartup == TriggerSchedule || TriggerStartup == TriggerEvent {
		t.Fatalf("startup origin must remain distinct: startup=%q schedule=%q event=%q", TriggerStartup, TriggerSchedule, TriggerEvent)
	}
}

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
