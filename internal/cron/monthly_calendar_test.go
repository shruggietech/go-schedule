package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func TestExplainMonthlyCalendarSelectors(t *testing.T) {
	tests := []struct{ expr, want string }{
		{"0 9 L * *", "last day of every month at 09:00"},
		{"0 9 15W * *", "nearest weekday to the 15th of every month at 09:00"},
		{"0 9 LW * *", "last weekday of every month at 09:00"},
		{"0 9 15w * *", "nearest weekday to the 15th of every month at 09:00"},
		{"0 9 lw * *", "last weekday of every month at 09:00"},
	}
	for _, tt := range tests {
		got, bad, err := Explain(tt.expr)
		if err != nil || bad.Reason != "" || got != tt.want {
			t.Fatalf("Explain(%q) = %q, bad=%q err=%v; want %q", tt.expr, got, bad.Reason, err, tt.want)
		}
	}
}

func TestMonthlyCalendarSelectorRefusals(t *testing.T) {
	for _, expr := range []string{
		"0 9 L JAN *", "0 9 15W * 1", "0 9 L-1 * *", "0 9 15W,16W * *",
		"0 9 1-5W * *", "0 9 15W/2 * *",
	} {
		res, err := Parse(expr)
		if err == nil && res.OK {
			t.Fatalf("Parse(%q) unexpectedly accepted", expr)
		}
	}
	for _, expr := range []string{"0 9 0W * *", "0 9 32W * *", "0 9 W * *", "0 9 LWW * *"} {
		if res, err := Parse(expr); err == nil && res.OK {
			t.Fatalf("Parse(%q) unexpectedly accepted", expr)
		}
	}
}

func TestExportMonthlyCalendarSelectors(t *testing.T) {
	tests := []struct {
		phrase string
		policy domain.MissingDatePolicy
		want   string
		ok     bool
	}{
		{"last day of every month at 09:00", domain.MissingDateNextValid, "0 9 L * *", true},
		{"nearest weekday to the 15th of every month at 09:00", domain.MissingDateLastValid, "0 9 15W * *", true},
		{"nearest weekday to the 31st of every month at 09:00", domain.MissingDateSkip, "0 9 31W * *", true},
		{"nearest weekday to the 31st of every month at 09:00", domain.MissingDateLastValid, "", false},
		{"last weekday of every month at 09:00", domain.MissingDateNextValid, "0 9 LW * *", true},
	}
	for _, tt := range tests {
		sch, err := schedule.Parse(tt.phrase, "UTC", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		got, _, ok := ExportSchedule(sch, tt.policy)
		if ok != tt.ok || got != tt.want {
			t.Fatalf("ExportSchedule(%q,%s) = %q,%v; want %q,%v", tt.phrase, tt.policy, got, ok, tt.want, tt.ok)
		}
	}
	policies := []domain.MissingDatePolicy{domain.MissingDateSkip, domain.MissingDateLastValid, domain.MissingDateNextValid}
	for _, day := range []int{1, 15, 28, 29, 30, 31} {
		phrase := "nearest weekday to the " + ordinal(day) + " of every month at 09:00"
		sch, err := schedule.Parse(phrase, "UTC", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		for _, policy := range policies {
			got, _, ok := ExportSchedule(sch, policy)
			wantOK := day <= 28 || policy == domain.MissingDateSkip
			if ok != wantOK {
				t.Fatalf("day=%d policy=%s: ok=%v want=%v output=%q", day, policy, ok, wantOK, got)
			}
			if wantOK && got != fmt.Sprintf("0 9 %dW * *", day) {
				t.Fatalf("day=%d policy=%s output=%q", day, policy, got)
			}
		}
	}
}

func TestMonthlyCalendarTwelveMonthRoundTrip(t *testing.T) {
	anchor := time.Date(2027, 9, 1, 0, 0, 0, 0, time.UTC)
	for _, phrase := range []string{
		"last day of every month at 09:00",
		"nearest weekday to the 15th of every month at 09:00",
		"last weekday of every month at 09:00",
	} {
		native, err := schedule.Parse(phrase, "America/New_York", anchor)
		if err != nil {
			t.Fatal(err)
		}
		expr, bad, ok := ExportSchedule(native, domain.MissingDateSkip)
		if !ok {
			t.Fatalf("export %q: %s", phrase, bad.Reason)
		}
		explained, bad, err := Explain(expr)
		if err != nil || bad.Reason != "" {
			t.Fatalf("explain %q: bad=%s err=%v", expr, bad.Reason, err)
		}
		roundTrip, err := schedule.Parse(explained, "America/New_York", anchor)
		if err != nil {
			t.Fatal(err)
		}
		want, err := schedule.UpcomingRuns(native, "America/New_York", domain.MissingDateSkip, anchor, 12)
		if err != nil {
			t.Fatal(err)
		}
		got, err := schedule.UpcomingRuns(roundTrip, "America/New_York", domain.MissingDateSkip, anchor, 12)
		if err != nil || len(got) != len(want) {
			t.Fatalf("runs %q: got=%d want=%d err=%v", phrase, len(got), len(want), err)
		}
		for i := range want {
			if !got[i].Equal(want[i]) {
				t.Fatalf("%q run %d: got=%v want=%v", phrase, i, got[i], want[i])
			}
			local := got[i].In(mustLocation(t, "America/New_York"))
			if local.Hour() != 9 {
				t.Fatalf("%q run %d lost wall time: %v", phrase, i, local)
			}
		}
	}
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return loc
}
