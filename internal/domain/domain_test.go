package domain

import (
	"encoding/json"
	"testing"
)

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

func TestExecutionDiagnosticJSONCompatibility(t *testing.T) {
	legacyRun, err := json.Marshal(Run{ID: "run-legacy"})
	if err != nil {
		t.Fatal(err)
	}
	if jsonContainsField(t, legacyRun, "output_truncated") {
		t.Fatalf("zero-value run unexpectedly changed legacy JSON: %s", legacyRun)
	}
	legacyAlert, err := json.Marshal(Alert{ID: "alert-legacy"})
	if err != nil {
		t.Fatal(err)
	}
	if jsonContainsField(t, legacyAlert, "run_id") {
		t.Fatalf("zero-value alert unexpectedly changed legacy JSON: %s", legacyAlert)
	}

	runJSON, err := json.Marshal(Run{ID: "run-truncated", OutputTruncated: true})
	if err != nil {
		t.Fatal(err)
	}
	if !jsonContainsField(t, runJSON, "output_truncated") {
		t.Fatalf("truncated run omitted diagnostic metadata: %s", runJSON)
	}
	alertJSON, err := json.Marshal(Alert{ID: "alert-correlated", RunID: "run-truncated"})
	if err != nil {
		t.Fatal(err)
	}
	if !jsonContainsField(t, alertJSON, "run_id") {
		t.Fatalf("correlated alert omitted run identity: %s", alertJSON)
	}
}

func jsonContainsField(t *testing.T, data []byte, field string) bool {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	_, ok := value[field]
	return ok
}
