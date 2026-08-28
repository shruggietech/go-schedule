package schedule

import (
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestParseMonthlyCalendarSelectors(t *testing.T) {
	tests := []struct {
		phrase, token, summary string
		adjustment             domain.CalendarAdjustment
	}{
		{"last day of every month at 09:00", "BYMONTHDAY=-1", "Last day of every month at 09:00", ""},
		{"nearest weekday to the 15th of every month at 09:00", "BYMONTHDAY=15", "Nearest weekday to the 15th of every month at 09:00", domain.CalendarAdjustmentNearestWeekday},
		{"last weekday of every month at 09:00", "BYSETPOS=-1", "Last weekday of every month at 09:00", ""},
	}
	for _, tt := range tests {
		t.Run(tt.phrase, func(t *testing.T) {
			sch, err := Parse(tt.phrase, "UTC", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(sch.RRULE, tt.token) {
				t.Fatalf("RRULE %q missing %q", sch.RRULE, tt.token)
			}
			if sch.HumanSummary != tt.summary {
				t.Fatalf("summary = %q, want %q", sch.HumanSummary, tt.summary)
			}
			if sch.CalendarAdjustment != tt.adjustment {
				t.Fatalf("adjustment = %q, want %q", sch.CalendarAdjustment, tt.adjustment)
			}
		})
	}
}

func TestNearestWeekdayCalendarSemantics(t *testing.T) {
	tests := []struct {
		phrase string
		after  time.Time
		want   time.Time
	}{
		{"nearest weekday to the 1st of every month at 09:00", time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC), time.Date(2026, 8, 3, 9, 0, 0, 0, time.UTC)},    // Saturday first
		{"nearest weekday to the 15th of every month at 09:00", time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 7, 15, 9, 0, 0, 0, time.UTC)},   // Wednesday
		{"nearest weekday to the 15th of every month at 09:00", time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC)},   // Saturday
		{"nearest weekday to the 15th of every month at 09:00", time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC)}, // Sunday
		{"nearest weekday to the 31st of every month at 09:00", time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 5, 29, 9, 0, 0, 0, time.UTC)},   // skip April; Sunday month-end
	}
	for _, tt := range tests {
		t.Run(tt.phrase+tt.after.Format("2006-01"), func(t *testing.T) {
			sch, err := Parse(tt.phrase, "UTC", tt.after)
			if err != nil {
				t.Fatal(err)
			}
			got, ok, err := NextRun(sch, "UTC", domain.MissingDateSkip, tt.after)
			if err != nil || !ok {
				t.Fatalf("NextRun: ok=%v err=%v", ok, err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNearestWeekdayDoesNotPrecedeScheduleAnchor(t *testing.T) {
	anchor := time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
	sch := domain.Schedule{
		Kind:               domain.ScheduleRecurring,
		RRULE:              "FREQ=MONTHLY;BYMONTHDAY=15;BYHOUR=9;BYMINUTE=0;BYSECOND=0",
		Anchor:             &anchor,
		CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday,
	}
	after := time.Date(2026, 7, 15, 9, 0, 0, 0, time.UTC)
	want := time.Date(2026, 9, 15, 9, 0, 0, 0, time.UTC)

	got, ok, err := NextRun(sch, "UTC", domain.MissingDateSkip, after)
	if err != nil || !ok || !got.Equal(want) {
		t.Fatalf("got=%v ok=%v err=%v, want %v", got, ok, err, want)
	}
}

func TestNearestWeekdayPreservesWallTimeAcrossDSTGap(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	anchor := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)
	sch := domain.Schedule{
		Kind:               domain.ScheduleRecurring,
		RRULE:              "FREQ=MONTHLY;BYMONTHDAY=8;BYHOUR=2;BYMINUTE=30;BYSECOND=0",
		Anchor:             &anchor,
		CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday,
	}
	after := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)
	want := time.Date(2026, 3, 9, 2, 30, 0, 0, loc).UTC()

	got, ok, err := NextRun(sch, "America/New_York", domain.MissingDateSkip, after)
	if err != nil || !ok || !got.Equal(want) {
		t.Fatalf("got=%v ok=%v err=%v, want %v", got, ok, err, want)
	}
}

func TestNearestWeekdayAllWeekdayPositions(t *testing.T) {
	anchor := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	sch, err := Parse("nearest weekday to the 15th of every month at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	wants := []time.Time{
		time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), // Thursday
		time.Date(2026, 2, 16, 9, 0, 0, 0, time.UTC), // Sunday -> Monday
		time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC), // Wednesday
		time.Date(2026, 5, 15, 9, 0, 0, 0, time.UTC), // Friday
		time.Date(2026, 6, 15, 9, 0, 0, 0, time.UTC), // Monday
		time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC), // Saturday -> Friday
		time.Date(2026, 9, 15, 9, 0, 0, 0, time.UTC), // Tuesday
	}
	for _, want := range wants {
		after := time.Date(want.Year(), want.Month(), 1, 0, 0, 0, 0, time.UTC)
		got, ok, err := NextRun(sch, "UTC", domain.MissingDateSkip, after)
		if err != nil || !ok || !got.Equal(want) {
			t.Fatalf("month %s: got=%v ok=%v err=%v want=%v", want.Month(), got, ok, err, want)
		}
	}
}

func TestLastDayAndLastWeekdayCalendarSemantics(t *testing.T) {
	anchor := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	lastDay, err := Parse("last day of every month at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	got, ok, err := NextRun(lastDay, "UTC", domain.MissingDateSkip, time.Date(2028, 2, 1, 0, 0, 0, 0, time.UTC))
	wantLeapDay := time.Date(2028, 2, 29, 9, 0, 0, 0, time.UTC)
	if err != nil || !ok || !got.Equal(wantLeapDay) {
		t.Fatalf("last day: got=%v ok=%v err=%v want=%v", got, ok, err, wantLeapDay)
	}

	lastWeekday, err := Parse("last weekday of every month at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	wants := []time.Time{
		time.Date(2026, 8, 31, 9, 0, 0, 0, time.UTC), // Monday
		time.Date(2026, 3, 31, 9, 0, 0, 0, time.UTC), // Tuesday
		time.Date(2026, 9, 30, 9, 0, 0, 0, time.UTC), // Wednesday
		time.Date(2026, 4, 30, 9, 0, 0, 0, time.UTC), // Thursday
		time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC), // Friday
		time.Date(2026, 1, 30, 9, 0, 0, 0, time.UTC), // Saturday -> Friday
		time.Date(2026, 5, 29, 9, 0, 0, 0, time.UTC), // Sunday -> Friday
	}
	for _, want := range wants {
		after := time.Date(want.Year(), want.Month(), 1, 0, 0, 0, 0, time.UTC)
		got, ok, err := NextRun(lastWeekday, "UTC", domain.MissingDateSkip, after)
		if err != nil || !ok || !got.Equal(want) {
			t.Fatalf("month %s: got=%v ok=%v err=%v want=%v", want.Month(), got, ok, err, want)
		}
	}
}

func TestNearestWeekdayMissingDatePolicyOrder(t *testing.T) {
	sch, err := Parse("nearest weekday to the 31st of every month at 09:00", "UTC", time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	after := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		policy domain.MissingDatePolicy
		want   time.Time
	}{
		{domain.MissingDateSkip, time.Date(2026, 5, 29, 9, 0, 0, 0, time.UTC)},
		{domain.MissingDateLastValid, time.Date(2026, 4, 30, 9, 0, 0, 0, time.UTC)},
		{domain.MissingDateNextValid, time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		got, ok, err := NextRun(sch, "UTC", tt.policy, after)
		if err != nil || !ok || !got.Equal(tt.want) {
			t.Fatalf("policy %s: got=%v ok=%v err=%v want=%v", tt.policy, got, ok, err, tt.want)
		}
	}
	if got := Describe(sch, domain.MissingDateLastValid); !strings.Contains(got, "nearest weekday to the last day") {
		t.Fatalf("last-valid description = %q", got)
	}
	if got := Describe(sch, domain.MissingDateNextValid); !strings.Contains(got, "first day of the following month") {
		t.Fatalf("next-valid description = %q", got)
	}
	runs, err := UpcomingRuns(sch, "UTC", domain.MissingDateNextValid, after, 4)
	if err != nil {
		t.Fatal(err)
	}
	wants := []time.Time{
		time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 29, 9, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC),
	}
	for i := range wants {
		if i >= len(runs) || !runs[i].Equal(wants[i]) {
			t.Fatalf("next_valid runs = %v, want %v", runs, wants)
		}
	}
}

func TestNearestWeekdayRejectsInvalidStoredShape(t *testing.T) {
	sch := domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=DAILY;BYHOUR=9;BYMINUTE=0;BYSECOND=0", CalendarAdjustment: domain.CalendarAdjustmentNearestWeekday}
	_, _, err := NextRun(sch, "UTC", domain.MissingDateSkip, time.Now())
	if err == nil || !strings.Contains(err.Error(), "nearest_weekday") {
		t.Fatalf("error = %v, want nearest_weekday shape error", err)
	}
}
