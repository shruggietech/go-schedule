package schedule

import (
	"strings"
	"testing"
	"time"
)

var now = time.Date(2026, 6, 19, 8, 0, 0, 0, time.UTC)

func TestParse_Forms(t *testing.T) {
	tests := []struct {
		input       string
		wantTokens  []string // every token must appear in the RRULE
		wantSummary string
	}{
		{"every 15 minutes", []string{"FREQ=MINUTELY", "INTERVAL=15"}, "Every 15 minutes"},
		{"every 30s", []string{"FREQ=SECONDLY", "INTERVAL=30"}, "Every 30 seconds"},
		{"every 2 hours", []string{"FREQ=HOURLY", "INTERVAL=2"}, "Every 2 hours"},
		{"every day at 09:00", []string{"FREQ=DAILY", "INTERVAL=1", "BYHOUR=9", "BYMINUTE=0"}, "Every day at 09:00"},
		{"every 3 days", []string{"FREQ=DAILY", "INTERVAL=3"}, "Every 3 days"},
		{"weekdays at 9:00 AM", []string{"FREQ=WEEKLY", "BYDAY=MO,TU,WE,TH,FR", "BYHOUR=9"}, "Every weekday at 09:00"},
		{"every monday at 9am", []string{"FREQ=WEEKLY", "BYDAY=MO", "BYHOUR=9"}, "Every Monday at 09:00"},
		{"3rd wednesday monthly at 14:00", []string{"FREQ=MONTHLY", "BYDAY=+3WE", "BYHOUR=14"}, "The 3rd Wednesday of every month at 14:00"},
		{"last friday of the month", []string{"FREQ=MONTHLY", "BYDAY=-1FR"}, "The last Friday of every month"},
		// By-date monthly and yearly forms (008-cron-interop). Without these,
		// ordinary cron lines like "0 9 1 * *" have no target representation.
		{"on the 15th of every month", []string{"FREQ=MONTHLY", "BYMONTHDAY=15"}, "The 15th of every month"},
		{"the 31st monthly at 09:00", []string{"FREQ=MONTHLY", "BYMONTHDAY=31", "BYHOUR=9"}, "The 31st of every month at 09:00"},
		{"on the 1st of each month at 6am", []string{"FREQ=MONTHLY", "BYMONTHDAY=1", "BYHOUR=6"}, "The 1st of every month at 06:00"},
		{"on the 22nd of every month", []string{"FREQ=MONTHLY", "BYMONTHDAY=22"}, "The 22nd of every month"},
		{"every year on february 29", []string{"FREQ=YEARLY", "BYMONTH=2", "BYMONTHDAY=29"}, "Every year on February 29"},
		{"annually on 4 july at 12:00", []string{"FREQ=YEARLY", "BYMONTH=7", "BYMONTHDAY=4", "BYHOUR=12"}, "Every year on July 4 at 12:00"},
		{"every 12 months", []string{"FREQ=MONTHLY", "INTERVAL=12"}, "Every 12 months"},
		{"every year", []string{"FREQ=YEARLY", "INTERVAL=1"}, "Every year"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sch, err := Parse(tt.input, "UTC", now)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			for _, want := range tt.wantTokens {
				if !strings.Contains(sch.RRULE, want) {
					t.Fatalf("RRULE %q missing token %q", sch.RRULE, want)
				}
			}
			if sch.HumanSummary != tt.wantSummary {
				t.Fatalf("summary = %q, want %q", sch.HumanSummary, tt.wantSummary)
			}
		})
	}
}

func TestParse_Rejects(t *testing.T) {
	for _, bad := range []string{
		"", "soon", "every banana", "every 15 minutes at 09:00", "3rd wednesday monthly at 99:99",
		// By-date and yearly forms reject impossible dates at the grammar level.
		// A day that exists in *some* month (the 31st) is accepted here and left
		// to the missing-date policy; a day that exists in *no* month is not.
		"on the 32nd of every month", "on the 0th of every month",
		"every year on february 30", "every year on smarch 3", "every year on april 31",
	} {
		if _, err := Parse(bad, "UTC", now); err == nil {
			t.Fatalf("expected error for %q", bad)
		}
	}
}

func TestParse_IntervalAnchor(t *testing.T) {
	tests := []struct {
		input       string
		wantAnchor  time.Time
		wantSummary string
	}{
		{"every 15 minutes starting at 09:00", time.Date(2026, 6, 19, 9, 0, 0, 0, time.UTC), "Every 15 minutes starting at 09:00"},
		{"every 30 minutes from 9am", time.Date(2026, 6, 19, 9, 0, 0, 0, time.UTC), "Every 30 minutes starting at 09:00"},
		{"every 2 hours starting at 08:30", time.Date(2026, 6, 19, 8, 30, 0, 0, time.UTC), "Every 2 hours starting at 08:30"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sch, err := Parse(tt.input, "UTC", now)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.input, err)
			}
			if sch.Anchor == nil || !sch.Anchor.Equal(tt.wantAnchor) {
				t.Fatalf("anchor = %v, want %v", sch.Anchor, tt.wantAnchor)
			}
			if sch.HumanSummary != tt.wantSummary {
				t.Fatalf("summary = %q, want %q", sch.HumanSummary, tt.wantSummary)
			}
		})
	}
}

func TestParse_AnchorRejects(t *testing.T) {
	for _, bad := range []string{
		"every 15 minutes at 09:00",   // bare 'at' still invalid for sub-daily
		"every day starting at 09:00", // anchor not valid for daily
		"every week from 9am",         // anchor not valid for weekly
		"every 15 minutes starting at 99:99",
	} {
		if _, err := Parse(bad, "UTC", now); err == nil {
			t.Fatalf("expected error for %q", bad)
		}
	}
}

func TestParse_AnchorTimezone(t *testing.T) {
	// Anchor 09:00 in New York (EDT, -04:00) on now's date → 13:00 UTC.
	sch, err := Parse("every 15 minutes starting at 09:00", "America/New_York", now)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 6, 19, 13, 0, 0, 0, time.UTC)
	if sch.Anchor == nil || !sch.Anchor.Equal(want) {
		t.Fatalf("anchor = %v, want %v", sch.Anchor, want)
	}
}

func TestParse_TimeOfDayVariants(t *testing.T) {
	for _, in := range []string{"every day at 14:00", "every day at 2:00 PM", "every day at 2pm"} {
		sch, err := Parse(in, "UTC", now)
		if err != nil {
			t.Fatalf("Parse(%q): %v", in, err)
		}
		if !strings.Contains(sch.RRULE, "BYHOUR=14") {
			t.Fatalf("%q did not yield 14:00, got %q", in, sch.RRULE)
		}
	}
}

// TestParse_RetainsExpression pins the round-trip half of the parser's contract:
// the phrase the user typed is kept on the schedule so an editing client can put
// their own words back in front of them. It never feeds back into evaluation.
func TestParse_RetainsExpression(t *testing.T) {
	sch, err := Parse("  Weekdays at 09:00  ", "UTC", now)
	if err != nil {
		t.Fatal(err)
	}
	if sch.Expression != "Weekdays at 09:00" {
		t.Errorf("Expression = %q, want the trimmed original phrase", sch.Expression)
	}
	// One-offs carry no phrase: their time is recovered from RunAt.
	if got := NewOneOff(now).Expression; got != "" {
		t.Errorf("one-off Expression = %q, want empty", got)
	}
}

// TestParse_NewFormsRoundTrip is FR-018 for the forms this feature adds: the
// phrase the operator typed is retained verbatim, and re-parsing the retained
// phrase yields the same rule. Without this, a client that offers the stored
// phrase for editing would hand back something that no longer parses.
func TestParse_NewFormsRoundTrip(t *testing.T) {
	for _, phrase := range []string{
		"on the 15th of every month at 09:00",
		"the 31st monthly",
		"every year on february 29 at 09:00",
		"annually on 4 july",
		"every 12 months",
	} {
		t.Run(phrase, func(t *testing.T) {
			first, err := Parse(phrase, "UTC", now)
			if err != nil {
				t.Fatalf("Parse(%q): %v", phrase, err)
			}
			if first.Expression != phrase {
				t.Fatalf("Expression = %q, want %q", first.Expression, phrase)
			}
			second, err := Parse(first.Expression, "UTC", now)
			if err != nil {
				t.Fatalf("re-parsing the retained phrase failed: %v", err)
			}
			if second.RRULE != first.RRULE {
				t.Fatalf("round trip changed the rule: %q -> %q", first.RRULE, second.RRULE)
			}
			if second.HumanSummary != first.HumanSummary {
				t.Fatalf("round trip changed the summary: %q -> %q", first.HumanSummary, second.HumanSummary)
			}
		})
	}
}
