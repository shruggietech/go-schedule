package schedule_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/cron"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func mustCompileCron(t *testing.T, expr, zone string, anchor time.Time) domain.Schedule {
	t.Helper()
	sch, bad, err := cron.Compile(expr, zone, anchor)
	if err != nil || bad.Reason != "" {
		t.Fatalf("Compile(%q): refusal=%q err=%v", expr, bad.Reason, err)
	}
	return sch
}

func TestCompositeCronNextRunMatrix(t *testing.T) {
	tests := []struct {
		name, expr string
		anchor     time.Time
		after      time.Time
		want       []time.Time
	}{
		{
			name: "two daily times", expr: "0 9,17 * * *",
			anchor: time.Date(2026, 8, 28, 8, 0, 0, 0, time.UTC),
			after:  time.Date(2026, 8, 28, 8, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2026, 8, 28, 9, 0, 0, 0, time.UTC),
				time.Date(2026, 8, 28, 17, 0, 0, 0, time.UTC),
				time.Date(2026, 8, 29, 9, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "uneven minute step resets", expr: "*/7 * * * *",
			anchor: time.Date(2026, 8, 28, 0, 55, 0, 0, time.UTC),
			after:  time.Date(2026, 8, 28, 0, 55, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2026, 8, 28, 0, 56, 0, 0, time.UTC),
				time.Date(2026, 8, 28, 1, 0, 0, 0, time.UTC),
				time.Date(2026, 8, 28, 1, 7, 0, 0, time.UTC),
			},
		},
		{
			name: "month and date lists", expr: "0 0 1,15 JAN,MAR *",
			anchor: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			after:  time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "weekday set", expr: "30 8-17 * * MON,WED,FRI",
			anchor: time.Date(2026, 8, 28, 16, 45, 0, 0, time.UTC), // Friday
			after:  time.Date(2026, 8, 28, 16, 45, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2026, 8, 28, 17, 30, 0, 0, time.UTC),
				time.Date(2026, 8, 31, 8, 30, 0, 0, time.UTC),
				time.Date(2026, 8, 31, 9, 30, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sch := mustCompileCron(t, tt.expr, "UTC", tt.anchor)
			got, err := schedule.UpcomingRuns(sch, "UTC", domain.MissingDateSkip, tt.after, len(tt.want))
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("runs=%v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if !got[i].Equal(tt.want[i]) {
					t.Errorf("run %d=%v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCompositeCronRepresentativeFirstRuns(t *testing.T) {
	day := func(year int, month time.Month, date, hour, minute int) time.Time {
		return time.Date(year, month, date, hour, minute, 0, 0, time.UTC)
	}
	tests := []struct {
		expr  string
		after time.Time
		want  time.Time
	}{
		{"0 9,17 * * *", day(2026, 1, 1, 8, 0), day(2026, 1, 1, 9, 0)},
		{"30 8-17 * * *", day(2026, 1, 1, 8, 0), day(2026, 1, 1, 8, 30)},
		{"*/10 9-17 * * *", day(2026, 1, 1, 8, 59), day(2026, 1, 1, 9, 0)},
		{"10-20/2 * * * *", day(2026, 1, 1, 0, 9), day(2026, 1, 1, 0, 10)},
		{"0 */5 * * *", day(2026, 1, 1, 0, 0), day(2026, 1, 1, 5, 0)},
		{"0 9 1,15 * *", day(2026, 1, 1, 9, 0), day(2026, 1, 15, 9, 0)},
		{"0 0 1 JAN,MAR *", day(2026, 1, 2, 0, 0), day(2026, 3, 1, 0, 0)},
		{"0 9 * * MON,WED,FRI", day(2026, 1, 2, 9, 0), day(2026, 1, 5, 9, 0)},
		{"0 9 * * MON-FRI", day(2026, 1, 2, 9, 0), day(2026, 1, 5, 9, 0)},
		{"0 9 * * */2", day(2026, 1, 4, 9, 0), day(2026, 1, 6, 9, 0)},
		{"0 9 */2 * *", day(2026, 1, 1, 9, 0), day(2026, 1, 3, 9, 0)},
		{"0 9 * */2 *", day(2026, 1, 31, 10, 0), day(2026, 3, 1, 9, 0)},
		{"0,1,1-2 * * * *", day(2026, 1, 1, 0, 0), day(2026, 1, 1, 0, 1)},
		{"0 9-11,10-12 * * *", day(2026, 1, 1, 9, 0), day(2026, 1, 1, 10, 0)},
		{"0 9 * * mon,wed", day(2026, 1, 5, 9, 0), day(2026, 1, 7, 9, 0)},
		{"0 9 * * 0,7", day(2026, 1, 2, 0, 0), day(2026, 1, 4, 9, 0)},
		{"0 8-16/4 * * *", day(2026, 1, 1, 8, 0), day(2026, 1, 1, 12, 0)},
		{"0,15 9,17 * * *", day(2026, 1, 1, 9, 0), day(2026, 1, 1, 9, 15)},
		{"0 9 10-12 * *", day(2026, 1, 1, 0, 0), day(2026, 1, 10, 9, 0)},
		{"0 0 29 FEB *", day(2026, 1, 1, 0, 0), day(2028, 2, 29, 0, 0)},
	}
	anchor := day(2025, 1, 1, 0, 0)
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			sch := mustCompileCron(t, tt.expr, "UTC", anchor)
			got, ok, err := schedule.NextRun(sch, "UTC", domain.MissingDateSkip, tt.after)
			if err != nil || !ok || !got.Equal(tt.want) {
				t.Fatalf("next=%v ok=%v err=%v, want %v", got, ok, err, tt.want)
			}
		})
	}
}

func TestCompositeCronStrictlyHonorsAnchor(t *testing.T) {
	anchor := time.Date(2026, 8, 28, 9, 30, 0, 0, time.UTC)
	sch := mustCompileCron(t, "0 9,17 * * *", "UTC", anchor)
	got, ok, err := schedule.NextRun(sch, "UTC", domain.MissingDateSkip, time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))
	want := time.Date(2026, 8, 28, 17, 0, 0, 0, time.UTC)
	if err != nil || !ok || !got.Equal(want) {
		t.Fatalf("got=%v ok=%v err=%v, want %v", got, ok, err, want)
	}
}

func TestCompositeCronDSTWallTime(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, expr string
		anchor     time.Time
		want       time.Time
	}{
		{
			name: "spring gap advances to first valid instant", expr: "30 2 * * SUN",
			anchor: time.Date(2026, 3, 2, 0, 0, 0, 0, loc),
			want:   time.Date(2026, 3, 8, 7, 0, 0, 0, time.UTC),
		},
		{
			name: "fall overlap uses first occurrence", expr: "30 1 * * SUN",
			anchor: time.Date(2026, 10, 26, 0, 0, 0, 0, loc),
			want:   time.Date(2026, 11, 1, 5, 30, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sch := mustCompileCron(t, tt.expr, "America/New_York", tt.anchor)
			got, ok, err := schedule.NextRun(sch, "America/New_York", domain.MissingDateSkip, tt.anchor)
			if err != nil || !ok || !got.Equal(tt.want) {
				t.Fatalf("got=%v ok=%v err=%v, want %v", got, ok, err, tt.want)
			}
		})
	}
}

func TestCompositeCronMissingDatePoliciesAndCollisionSuppression(t *testing.T) {
	anchor := time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)
	sch := mustCompileCron(t, "0 9,17 30,31 * *", "UTC", anchor)
	tests := []struct {
		policy domain.MissingDatePolicy
		want   []time.Time
	}{
		{domain.MissingDateSkip, []time.Time{
			time.Date(2027, 3, 30, 9, 0, 0, 0, time.UTC),
			time.Date(2027, 3, 30, 17, 0, 0, 0, time.UTC),
			time.Date(2027, 3, 31, 9, 0, 0, 0, time.UTC),
		}},
		{domain.MissingDateLastValid, []time.Time{
			time.Date(2027, 2, 28, 9, 0, 0, 0, time.UTC),
			time.Date(2027, 2, 28, 17, 0, 0, 0, time.UTC),
			time.Date(2027, 3, 30, 9, 0, 0, 0, time.UTC),
		}},
		{domain.MissingDateNextValid, []time.Time{
			time.Date(2027, 3, 1, 9, 0, 0, 0, time.UTC),
			time.Date(2027, 3, 1, 17, 0, 0, 0, time.UTC),
			time.Date(2027, 3, 30, 9, 0, 0, 0, time.UTC),
		}},
	}
	for _, tt := range tests {
		t.Run(string(tt.policy), func(t *testing.T) {
			got, err := schedule.UpcomingRuns(sch, "UTC", tt.policy, anchor, len(tt.want))
			if err != nil {
				t.Fatal(err)
			}
			for i := range tt.want {
				if i >= len(got) || !got[i].Equal(tt.want[i]) {
					t.Fatalf("runs=%v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestCompositeCronDescriptionNamesEffectiveMissingDatePolicy(t *testing.T) {
	anchor := time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)
	sch := mustCompileCron(t, "0 9 1,31 * *", "UTC", anchor)
	for policy, want := range map[domain.MissingDatePolicy]string{
		domain.MissingDateSkip:      "skipped in months that have no such date",
		domain.MissingDateLastValid: "last day of the month when there is no such date",
		domain.MissingDateNextValid: "rolling into the next period in months that have no such date",
	} {
		if got := schedule.Describe(sch, policy); !strings.Contains(got, want) {
			t.Errorf("Describe(%s)=%q, want policy clause %q", policy, got, want)
		}
	}

	inert := mustCompileCron(t, "0 9 1,15 * *", "UTC", anchor)
	if got := schedule.Describe(inert, domain.MissingDateLastValid); got != inert.HumanSummary {
		t.Fatalf("policy-inert description=%q, want %q", got, inert.HumanSummary)
	}
}
