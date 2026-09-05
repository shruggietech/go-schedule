package gui

import (
	"math"
	"reflect"
	"testing"
)

func TestColumnProfilePersistsNormalizedPerViewLayouts(t *testing.T) {
	prefs := testApp.Preferences()
	const scheduleKey = "test.columns.schedule"
	const activityKey = "test.columns.activity"
	t.Cleanup(func() {
		prefs.RemoveValue(scheduleKey)
		prefs.RemoveValue(activityKey)
	})
	columns := []structuredColumn{
		{Header: "When", Minimum: 100, Preferred: 2},
		{Header: "Task", Minimum: 100, Preferred: 3},
		{Header: "State", Minimum: 80, Preferred: 1},
	}
	schedule := newColumnProfile(columns, prefs, scheduleKey)
	activity := newColumnProfile(columns, prefs, activityKey)
	schedule.setProportions([]float32{0.5, 0.3, 0.2}, true)
	activity.setProportions([]float32{0.2, 0.3, 0.5}, true)

	if got := newColumnProfile(columns, prefs, scheduleKey).proportions(); !closeFloat32s(got, []float32{0.5, 0.3, 0.2}) {
		t.Fatalf("restored Schedule proportions = %v", got)
	}
	if got := newColumnProfile(columns, prefs, activityKey).proportions(); !closeFloat32s(got, []float32{0.2, 0.3, 0.5}) {
		t.Fatalf("restored Activity proportions = %v", got)
	}
}

func TestColumnProfileRejectsInvalidAndObsoletePreferencesAtomically(t *testing.T) {
	prefs := testApp.Preferences()
	const key = "test.columns.invalid"
	t.Cleanup(func() { prefs.RemoveValue(key) })
	columns := []structuredColumn{
		{Header: "When", Preferred: 2},
		{Header: "Task", Preferred: 3},
		{Header: "State", Preferred: 1},
	}
	want := []float32{1.0 / 3, 0.5, 1.0 / 6}
	for name, raw := range map[string]string{
		"malformed":               "{",
		"wrong version":           `{"version":99,"columns":["When","Task","State"],"proportions":[0.5,0.3,0.2]}`,
		"wrong columns":           `{"version":1,"columns":["When","Other","State"],"proportions":[0.5,0.3,0.2]}`,
		"wrong length":            `{"version":1,"columns":["When","Task","State"],"proportions":[0.5,0.5]}`,
		"non positive":            `{"version":1,"columns":["When","Task","State"],"proportions":[0.5,0.5,0]}`,
		"non finite":              `{"version":1,"columns":["When","Task","State"],"proportions":[1e999,0.5,0.5]}`,
		"sum overflow":            `{"version":1,"columns":["When","Task","State"],"proportions":[3e38,3e38,1]}`,
		"normalization underflow": `{"version":1,"columns":["When","Task","State"],"proportions":[3e38,1e-45,1e-45]}`,
	} {
		t.Run(name, func(t *testing.T) {
			prefs.SetString(key, raw)
			if got := newColumnProfile(columns, prefs, key).proportions(); !closeFloat32s(got, want) {
				t.Fatalf("fallback proportions = %v, want %v", got, want)
			}
		})
	}
}

func TestColumnProfileAdjustsAdjacentColumnsAndResets(t *testing.T) {
	prefs := testApp.Preferences()
	const key = "test.columns.adjust"
	t.Cleanup(func() { prefs.RemoveValue(key) })
	columns := []structuredColumn{
		{Header: "A", Minimum: 100, Preferred: 1},
		{Header: "B", Minimum: 80, Preferred: 1},
		{Header: "C", Minimum: 60, Preferred: 1},
	}
	profile := newColumnProfile(columns, prefs, key)
	before := profile.widths(600, structuredColumnGap)
	if !profile.adjust(0, 40, 600, structuredColumnGap) {
		t.Fatal("expected adjustment")
	}
	after := profile.widths(600, structuredColumnGap)
	if math.Abs(float64(after[0]-before[0]-40)) > 0.01 || math.Abs(float64(after[1]-before[1]+40)) > 0.01 || after[2] != before[2] {
		t.Fatalf("before=%v after=%v, want adjacent 40-unit transfer", before, after)
	}
	profile.persist()
	if got := newColumnProfile(columns, prefs, key).widths(600, structuredColumnGap); !closeFloat32s(got, after) {
		t.Fatalf("restored widths = %v, want %v", got, after)
	}
	profile.reset()
	if got := profile.proportions(); !reflect.DeepEqual(got, []float32{1.0 / 3, 1.0 / 3, 1.0 / 3}) {
		t.Fatalf("reset proportions = %v", got)
	}
	if profile.adjust(-1, 10, 600, structuredColumnGap) || profile.adjust(2, 10, 600, structuredColumnGap) || profile.adjust(0, float32(math.NaN()), 600, structuredColumnGap) {
		t.Fatal("invalid adjustment changed the profile")
	}
	for profile.adjust(0, -100, 600, structuredColumnGap) {
	}
	clamped := profile.widths(600, structuredColumnGap)
	if clamped[0] < columns[0].Minimum || clamped[1] < columns[1].Minimum {
		t.Fatalf("adjustment crossed minimums: %v", clamped)
	}
	if profile.adjust(0, 10, 100, structuredColumnGap) {
		t.Fatal("below-minimum compressed layout should not accept adjustment")
	}
}

func closeFloat32s(got, want []float32) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if math.Abs(float64(got[i]-want[i])) > 0.001 {
			return false
		}
	}
	return true
}
