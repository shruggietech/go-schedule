package cron

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

// TestExplain_Supported is the conversion table. Each case pairs a cron
// expression with the phrase a user would have typed, and then asserts the
// phrase actually parses — so a phrase this package invents but the grammar
// cannot read fails here rather than at import time.
func TestExplain_Supported(t *testing.T) {
	cases := []struct {
		expr, phrase string
	}{
		{"*/15 * * * *", "every 15 minutes starting at 00:00"},
		{"*/5 * * * *", "every 5 minutes starting at 00:00"},
		{"* * * * *", "every minute starting at 00:00"},
		{"0 * * * *", "every hour starting at 00:00"},
		{"0 */6 * * *", "every 6 hours starting at 00:00"},
		{"0 9 * * *", "every day at 09:00"},
		{"30 2 * * *", "every day at 02:30"},
		{"0 9 * * 1-5", "weekdays at 09:00"},
		{"0 9 * * MON-FRI", "weekdays at 09:00"},
		{"0 10 * * 0,6", "weekends at 10:00"},
		{"0 14 * * 3", "every wednesday at 14:00"},
		{"0 14 * * WED", "every wednesday at 14:00"},
		{"0 9 * * 5#3", "3rd friday monthly at 09:00"},
		{"30 14 * * WED#2", "2nd wednesday monthly at 14:30"},
		{"0 8 * * 7#1", "1st sunday monthly at 08:00"},
		{"0 9 * * 5L", "last friday of the month at 09:00"},
		{"30 14 * * WEDL", "last wednesday of the month at 14:30"},
		{"0 8 * * 7L", "last sunday of the month at 08:00"},
		{"0 9 1 * *", "on the 1st of every month at 09:00"},
		{"0 9 31 * *", "on the 31st of every month at 09:00"},
		{"0 0 29 2 *", "every year on february 29 at 00:00"},
		{"0 0 4 7 *", "every year on july 4 at 00:00"},
		// Shorthands expand to their documented five-field equivalents.
		{"@daily", "every day at 00:00"},
		{"@midnight", "every day at 00:00"},
		{"@hourly", "every hour starting at 00:00"},
		{"@weekly", "every sunday at 00:00"},
		{"@monthly", "on the 1st of every month at 00:00"},
		{"@yearly", "every year on january 1 at 00:00"},
		{"@annually", "every year on january 1 at 00:00"},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			phrase, bad, err := Explain(c.expr)
			if err != nil {
				t.Fatalf("Explain(%q): %v", c.expr, err)
			}
			if bad.Reason != "" {
				t.Fatalf("Explain(%q) refused: %s", c.expr, bad.Reason)
			}
			if phrase != c.phrase {
				t.Fatalf("phrase = %q, want %q", phrase, c.phrase)
			}
			// The phrase must be readable by the grammar — this is the
			// single-route guarantee (FR-003a).
			if _, err := schedule.Parse(phrase, "UTC", time.Now().UTC()); err != nil {
				t.Fatalf("phrase %q does not parse: %v", phrase, err)
			}
		})
	}
}

// TestExplain_Declines covers the remaining semantic boundaries: everything
// this package still cannot represent is named, not approximated or dropped.
func TestExplain_Declines(t *testing.T) {
	cases := []struct {
		expr, contains string
	}{
		{"@reboot", "boot"},
		{"0 0 * * * *", "six-field"},
		{"0 0 * JAN 5#3", "month"},
		{"0 0 1 * 5#3", "either"},
		{"0 0 * * 5#3,2#4", "one weekday"},
		{"0 0 * * 1-5#3", "one weekday"},
		{"0 0 * * 5#3/2", "one weekday"},
		{"0#2 0 * * *", "#"},
		{"0 0 * JAN 5L", "month"},
		{"0 0 1 * 5L", "either"},
		{"0 0 * * 5L,2L", "one last weekday"},
		{"0 0 * * 1-5L", "one last weekday"},
		{"0 0 * * 5L/2", "one last weekday"},
		{"0 0 * * 5L#2", "one last weekday"},
		{"0 0 13 * 5", "either"},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			phrase, bad, err := Explain(c.expr)
			if err != nil {
				t.Fatalf("Explain(%q) returned an error, want a refusal: %v", c.expr, err)
			}
			if phrase != "" {
				t.Fatalf("Explain(%q) produced a phrase %q, want a refusal", c.expr, phrase)
			}
			if !strings.Contains(bad.Reason, c.contains) {
				t.Fatalf("reason = %q, want it to mention %q", bad.Reason, c.contains)
			}
		})
	}
}

func TestExplain_LastWeekdayMatrix(t *testing.T) {
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	for day, name := range days {
		expr := fmt.Sprintf("0 9 * * %dL", day)
		want := fmt.Sprintf("last %s of the month at 09:00", name)
		phrase, bad, err := Explain(expr)
		if err != nil || bad.Reason != "" || phrase != want {
			t.Errorf("Explain(%q) = %q, refusal=%q, err=%v; want %q", expr, phrase, bad.Reason, err, want)
		}
	}
	for _, expr := range []string{"0 9 * * 0L", "0 9 * * 7L", "0 9 * * SUNL", "0 9 * * sundayL"} {
		phrase, bad, err := Explain(expr)
		if err != nil || bad.Reason != "" || phrase != "last sunday of the month at 09:00" {
			t.Errorf("Explain(%q) = %q, refusal=%q, err=%v", expr, phrase, bad.Reason, err)
		}
	}
}

func TestExplain_OrdinalWeekdayMatrix(t *testing.T) {
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	for day, name := range days {
		for occurrence := 1; occurrence <= 5; occurrence++ {
			expr := fmt.Sprintf("0 9 * * %d#%d", day, occurrence)
			want := fmt.Sprintf("%s %s monthly at 09:00", ordinal(occurrence), name)
			phrase, bad, err := Explain(expr)
			if err != nil || bad.Reason != "" || phrase != want {
				t.Errorf("Explain(%q) = %q, refusal=%q, err=%v; want %q", expr, phrase, bad.Reason, err, want)
			}
		}
	}
	for _, expr := range []string{"0 9 * * 0#2", "0 9 * * 7#2", "0 9 * * SUN#2", "0 9 * * sunday#2"} {
		phrase, bad, err := Explain(expr)
		if err != nil || bad.Reason != "" || phrase != "2nd sunday monthly at 09:00" {
			t.Errorf("Explain(%q) = %q, refusal=%q, err=%v", expr, phrase, bad.Reason, err)
		}
	}
}

func TestParse_CalendarWildcardStepOneRemainsUnrestricted(t *testing.T) {
	for _, expr := range []string{"0 9 */1 * *", "0 9 * */1 *", "0 9 * * */1"} {
		t.Run(expr, func(t *testing.T) {
			got, err := Parse(expr)
			if err != nil || !got.OK {
				t.Fatalf("Parse(%q): result=%+v err=%v", expr, got, err)
			}
		})
	}
}

// TestParse_Malformed covers the other half of the distinction: a typo is the
// operator's mistake and must be an error naming the field, not a refusal.
func TestParse_Malformed(t *testing.T) {
	for _, expr := range []string{
		"", "0 9 * *", "0 9 * * * * *", "@nonsense",
		"99 * * * *", "0 99 * * *", "0 9 32 * *", "0 9 * 13 *",
		"0 9 * * smarch", "*/0 * * * *", "5-1 * * * *",
		"0 9 * * #3", "0 9 * * 5#", "0 9 * * 5#0", "0 9 * * 5#6",
		"0 9 * * 5#third", "0 9 * * 5#3#2",
		"0 9 * * 8L", "0 9 * * FUNDAYL", "0 9 * * 5LL",
	} {
		t.Run(expr, func(t *testing.T) {
			if _, err := Parse(expr); err == nil {
				t.Fatalf("Parse(%q) succeeded, want an error", expr)
			}
		})
	}
}

// TestExplain_RunTimesMatchCron is the substance behind the conversion table:
// the phrase must produce the run times the cron expression describes, not
// merely read plausibly. Times are checked against hand-computed instants.
func TestExplain_RunTimesMatchCron(t *testing.T) {
	// 2026-06-01 is a Monday.
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		expr  string
		after time.Time
		want  time.Time
	}{
		{"*/15 * * * *", time.Date(2026, 6, 1, 0, 7, 0, 0, time.UTC), time.Date(2026, 6, 1, 0, 15, 0, 0, time.UTC)},
		{"0 9 * * 1-5", base, time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC)},
		{"0 14 * * 3", base, time.Date(2026, 6, 3, 14, 0, 0, 0, time.UTC)},
		{"0 9 1 * *", time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC), time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)},
		{"0 0 4 7 *", base, time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC)},
		{"0 9 * * 5#3", base, time.Date(2026, 6, 19, 9, 0, 0, 0, time.UTC)},
		{"0 9 * * 1#5", base, time.Date(2026, 6, 29, 9, 0, 0, 0, time.UTC)},
		{"0 9 * * 5L", base, time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			phrase, bad, err := Explain(c.expr)
			if err != nil || bad.Reason != "" {
				t.Fatalf("Explain(%q): err=%v refusal=%q", c.expr, err, bad.Reason)
			}
			sch, err := schedule.Parse(phrase, "UTC", base)
			if err != nil {
				t.Fatalf("Parse(%q): %v", phrase, err)
			}
			got, ok, err := schedule.NextRun(sch, "UTC", domain.MissingDateSkip, c.after)
			if err != nil || !ok {
				t.Fatalf("NextRun: ok=%v err=%v", ok, err)
			}
			if !got.Equal(c.want) {
				t.Fatalf("cron %q via %q: got %v, want %v", c.expr, phrase, got, c.want)
			}
		})
	}
}
