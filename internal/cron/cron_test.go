package cron

import (
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
		{"*/15 * * * *", "every 15 minutes"},
		{"*/5 * * * *", "every 5 minutes"},
		{"* * * * *", "every minute"},
		{"0 * * * *", "every hour"},
		{"0 */6 * * *", "every 6 hours starting at 00:00"},
		{"0 9 * * *", "every day at 09:00"},
		{"30 2 * * *", "every day at 02:30"},
		{"0 9 * * 1-5", "weekdays at 09:00"},
		{"0 9 * * MON-FRI", "weekdays at 09:00"},
		{"0 10 * * 0,6", "weekends at 10:00"},
		{"0 14 * * 3", "every wednesday at 14:00"},
		{"0 14 * * WED", "every wednesday at 14:00"},
		{"0 9 1 * *", "on the 1st of every month at 09:00"},
		{"0 9 31 * *", "on the 31st of every month at 09:00"},
		{"0 0 29 2 *", "every year on february 29 at 00:00"},
		{"0 0 4 7 *", "every year on july 4 at 00:00"},
		// Shorthands expand to their documented five-field equivalents.
		{"@daily", "every day at 00:00"},
		{"@midnight", "every day at 00:00"},
		{"@hourly", "every hour"},
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

// TestExplain_Declines covers FR-002 and FR-003b: everything this package will
// not represent is named, not approximated and not silently dropped. A refusal
// is a value, not an error.
func TestExplain_Declines(t *testing.T) {
	cases := []struct {
		expr, contains string
	}{
		{"@reboot", "boot"},
		{"0 0 * * * *", "six-field"},
		{"0 0 L * *", "L"},
		{"0 0 15W * *", "W"},
		{"0 0 * * 5#3", "#"},
		{"*/7 * * * *", "does not divide the hour evenly"},
		{"0 */7 * * *", "does not divide the day evenly"},
		{"0 0 13 * 5", "either"},
		{"0 9,17 * * *", "hour list"},
		{"0,30 9 * * *", "minute list"},
		{"0 9 1,15 * *", "day-of-month list"},
		{"0 9 * * 1,3", "weekdays has no phrase equivalent"},
		{"15 * * * *", "minute other than :00"},
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

// TestParse_Malformed covers the other half of the distinction: a typo is the
// operator's mistake and must be an error naming the field, not a refusal.
func TestParse_Malformed(t *testing.T) {
	for _, expr := range []string{
		"", "0 9 * *", "0 9 * * * * *", "@nonsense",
		"99 * * * *", "0 99 * * *", "0 9 32 * *", "0 9 * 13 *",
		"0 9 * * smarch", "*/0 * * * *", "5-1 * * * *",
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
