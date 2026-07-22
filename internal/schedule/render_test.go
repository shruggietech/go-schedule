package schedule

import (
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// roundTripPhrases is every phrase form the parser accepts, drawn from the
// tables in parse_test.go. Render must produce a phrase that re-parses to the
// same rule for each of them — that is the whole correctness contract (FR-004).
var roundTripPhrases = []string{
	"every 15 minutes",
	"every 30s",
	"every 2 hours",
	"every day at 09:00",
	"every 3 days",
	"every week",
	"weekdays at 9:00 AM",
	"weekends at 18:00",
	"every monday at 9am",
	"every sunday",
	"3rd wednesday monthly at 14:00",
	"last friday of the month",
	"1st monday monthly",
	// Anchored forms: the anchor is deliberately not rendered (FR-005), but the
	// recurrence itself must still round-trip.
	"every 15 minutes starting at 09:00",
	"every 2 hours from 8:30",
}

// TestRender_RoundTripsToSameRule is the FR-004 gate: a rendered phrase must
// describe the same recurrence as the schedule it came from.
func TestRender_RoundTripsToSameRule(t *testing.T) {
	for _, phrase := range roundTripPhrases {
		t.Run(phrase, func(t *testing.T) {
			original, err := Parse(phrase, "UTC", now)
			if err != nil {
				t.Fatalf("Parse(%q): %v", phrase, err)
			}

			rendered := Render(original, "UTC")
			if rendered == "" {
				t.Fatalf("Render produced nothing for a schedule this package created (RRULE %q)", original.RRULE)
			}

			reparsed, err := Parse(rendered, "UTC", now)
			if err != nil {
				t.Fatalf("rendered phrase %q does not re-parse: %v", rendered, err)
			}
			if reparsed.RRULE != original.RRULE {
				t.Errorf("round-trip changed the recurrence:\n  phrase:   %q\n  rendered: %q\n  before:   %s\n  after:    %s",
					phrase, rendered, original.RRULE, reparsed.RRULE)
			}
		})
	}
}

// TestRender_NeverSynthesizesAnAnchor pins FR-005. A stored Anchor is
// indistinguishable from the creation timestamp, so rendering one would put a
// time into the user's Start-at field that they never typed.
func TestRender_NeverSynthesizesAnAnchor(t *testing.T) {
	for _, phrase := range roundTripPhrases {
		sch, err := Parse(phrase, "UTC", now)
		if err != nil {
			t.Fatalf("Parse(%q): %v", phrase, err)
		}
		rendered := Render(sch, "UTC")
		if strings.Contains(rendered, "starting at") || strings.Contains(rendered, " from ") {
			t.Errorf("Render(%q) synthesized an anchor clause: %q", phrase, rendered)
		}
	}
}

// TestRender_SuppliesNothingRatherThanGuessing pins FR-003: outside the phrase
// vocabulary, Render returns "" instead of inventing wording.
func TestRender_SuppliesNothingRatherThanGuessing(t *testing.T) {
	cases := []struct {
		name string
		sch  domain.Schedule
	}{
		{"one-off has no phrase", NewOneOff(now)},
		{"event kind", domain.Schedule{Kind: domain.ScheduleEvent, TriggerID: "t1"}},
		{"empty rrule", domain.Schedule{Kind: domain.ScheduleRecurring}},
		{"unsupported frequency", domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=YEARLY;INTERVAL=1"}},
		{"unsupported byday set", domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=WEEKLY;BYDAY=MO,WE,FR"}},
		{"monthly by month-day", domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MONTHLY;BYMONTHDAY=15"}},
		{"garbage", domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "not-an-rrule"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Render(tc.sch, "UTC"); got != "" {
				t.Errorf("Render = %q, want empty (must not guess)", got)
			}
		})
	}
}

// TestParse_RetainsExpression pins the other half of the recovery path: what
// the user typed is kept verbatim (trimmed) on the schedule.
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
