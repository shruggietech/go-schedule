package schedule

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// These tests pin the missing-date policy to real calendar dates rather than
// synthetic offsets, per the constitution's testing standards: 2026 has 30-day
// months and 2027 is a common year, so both anomalies are reachable without
// inventing a calendar.

// runsFrom collects n run instants in tz for a phrase under a policy, anchored
// at the given moment.
func runsFrom(t *testing.T, phrase, tz string, policy domain.MissingDatePolicy, anchor time.Time, n int) []time.Time {
	t.Helper()
	sch, err := Parse(phrase, tz, anchor)
	if err != nil {
		t.Fatalf("Parse(%q): %v", phrase, err)
	}
	runs, err := UpcomingRuns(sch, tz, policy, anchor, n)
	if err != nil {
		t.Fatalf("UpcomingRuns(%q): %v", phrase, err)
	}
	return runs
}

func ymd(ts []time.Time) []string {
	out := make([]string, len(ts))
	for i, x := range ts {
		out[i] = x.UTC().Format("2006-01-02 15:04")
	}
	return out
}

func eqDates(t *testing.T, got []time.Time, want []string) {
	t.Helper()
	g := ymd(got)
	if len(g) != len(want) {
		t.Fatalf("got %d runs %v, want %d %v", len(g), g, len(want), want)
	}
	for i := range want {
		if g[i] != want[i] {
			t.Fatalf("run %d = %s, want %s (all: %v)", i, g[i], want[i], g)
		}
	}
}

// TestMissingDate_MonthDay31 is the headline case from issue #8: "the 31st of
// every month" fires seven months in twelve under the default, and the other two
// policies say what to do about it.
func TestMissingDate_MonthDay31(t *testing.T) {
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	const phrase = "on the 31st of every month at 09:00"

	t.Run("skip", func(t *testing.T) {
		// Jan, Mar, May, Jul, Aug — February, April and June have no 31st.
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateSkip, anchor, 5), []string{
			"2026-01-31 09:00", "2026-03-31 09:00", "2026-05-31 09:00",
			"2026-07-31 09:00", "2026-08-31 09:00",
		})
	})

	t.Run("last_valid", func(t *testing.T) {
		// Every month, clamped to the last day that exists.
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateLastValid, anchor, 5), []string{
			"2026-01-31 09:00", "2026-02-28 09:00", "2026-03-31 09:00",
			"2026-04-30 09:00", "2026-05-31 09:00",
		})
	})

	t.Run("next_valid", func(t *testing.T) {
		// February rolls to 1 March; that does not displace March's own 31st,
		// which still appears in its own period (FR-019a).
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateNextValid, anchor, 5), []string{
			"2026-01-31 09:00", "2026-03-01 09:00", "2026-03-31 09:00",
			"2026-05-01 09:00", "2026-05-31 09:00",
		})
	})
}

// TestMissingDate_LeapDay covers the yearly form across a common year. 2028 is a
// leap year; 2027 is not.
func TestMissingDate_LeapDay(t *testing.T) {
	anchor := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	const phrase = "every year on february 29 at 09:00"

	t.Run("skip", func(t *testing.T) {
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateSkip, anchor, 2), []string{
			"2028-02-29 09:00", "2032-02-29 09:00",
		})
	})

	t.Run("last_valid", func(t *testing.T) {
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateLastValid, anchor, 3), []string{
			"2027-02-28 09:00", "2028-02-29 09:00", "2029-02-28 09:00",
		})
	})

	t.Run("next_valid", func(t *testing.T) {
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateNextValid, anchor, 3), []string{
			"2027-03-01 09:00", "2028-02-29 09:00", "2029-03-01 09:00",
		})
	})
}

// TestMissingDate_Day30InFebruary is the case where fall-back and roll-forward
// differ by more than a day.
func TestMissingDate_Day30InFebruary(t *testing.T) {
	anchor := time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)
	const phrase = "on the 30th of every month at 12:00"

	eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateLastValid, anchor, 2), []string{
		"2027-02-28 12:00", "2027-03-30 12:00",
	})
	eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateNextValid, anchor, 2), []string{
		"2027-03-01 12:00", "2027-03-30 12:00",
	})
}

// TestMissingDate_FifthWeekday covers the rule whose existing behavior issue #8
// documents as wrong: "5th friday monthly" skips two thirds of the calendar and
// its summary claims otherwise.
func TestMissingDate_FifthWeekday(t *testing.T) {
	anchor := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	const phrase = "5th friday monthly at 09:00"

	t.Run("skip keeps the documented behavior", func(t *testing.T) {
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateSkip, anchor, 4), []string{
			"2026-05-29 09:00", "2026-07-31 09:00", "2026-10-30 09:00", "2027-01-29 09:00",
		})
	})

	t.Run("last_valid falls back to the last Friday", func(t *testing.T) {
		// June 2026 has four Fridays; the last is the 26th.
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateLastValid, anchor, 3), []string{
			"2026-05-29 09:00", "2026-06-26 09:00", "2026-07-31 09:00",
		})
	})

	t.Run("next_valid rolls into the next month", func(t *testing.T) {
		// June has no fifth Friday, so it rolls to the first Friday of July (3rd).
		eqDates(t, runsFrom(t, phrase, "UTC", domain.MissingDateNextValid, anchor, 3), []string{
			"2026-05-29 09:00", "2026-07-03 09:00", "2026-07-31 09:00",
		})
	})
}

// TestMissingDate_InertForDatelessRules is FR-024: a rule with no date component
// produces identical run times under every policy. If this fails, the policy has
// leaked into schedules it has no business touching.
func TestMissingDate_InertForDatelessRules(t *testing.T) {
	anchor := time.Date(2026, 2, 1, 8, 0, 0, 0, time.UTC)
	for _, phrase := range []string{
		"every 15 minutes",
		"every day at 09:00",
		"weekdays at 09:00",
		"every monday at 09:00",
		"last friday of the month at 09:00",
	} {
		t.Run(phrase, func(t *testing.T) {
			base := ymd(runsFrom(t, phrase, "UTC", domain.MissingDateSkip, anchor, 6))
			for _, p := range []domain.MissingDatePolicy{domain.MissingDateLastValid, domain.MissingDateNextValid} {
				got := ymd(runsFrom(t, phrase, "UTC", p, anchor, 6))
				if len(got) != len(base) {
					t.Fatalf("policy %s changed run count: %v vs %v", p, got, base)
				}
				for i := range base {
					if got[i] != base[i] {
						t.Fatalf("policy %s changed run %d: %s vs %s", p, i, got[i], base[i])
					}
				}
			}
		})
	}
}

// TestMissingDate_DSTStillApplies is FR-025: this feature resolves *which date*
// a run lands on and changes nothing about how a wall-clock time is resolved
// against a transition. 2026-03-08 is the US spring-forward.
func TestMissingDate_DSTStillApplies(t *testing.T) {
	const tz = "America/New_York"
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	// The 31st under last_valid puts February's run on the 28th; March's run on
	// the 31st is after the transition and must read 09:00 local (13:00 UTC),
	// while January's is before it (14:00 UTC).
	runs := runsFrom(t, "on the 31st of every month at 09:00", tz, domain.MissingDateLastValid, anchor, 3)
	eqDates(t, runs, []string{
		"2026-01-31 14:00", // EST, UTC-5
		"2026-02-28 14:00", // EST
		"2026-03-31 13:00", // EDT, UTC-4 — the offset moved, the local reading did not
	})

	// A run landing exactly on the transition date itself still resolves through
	// the same wall-clock rules.
	onTransition := runsFrom(t, "on the 8th of every month at 09:00", tz, domain.MissingDateLastValid,
		time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), 1)
	eqDates(t, onTransition, []string{"2026-03-08 13:00"})
}

// TestDescribe_NamesThePolicy is FR-023 / SC-005: no description may assert that
// a rule fires in every period when it does not.
func TestDescribe_NamesThePolicy(t *testing.T) {
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	byDate, err := Parse("on the 31st of every month at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		policy domain.MissingDatePolicy
		want   string
	}{
		{domain.MissingDateSkip, "The 31st of every month at 09:00, skipped in months that have no such date"},
		{domain.MissingDateLastValid, "The 31st of every month at 09:00, or the last day of the month when there is no such date"},
		{domain.MissingDateNextValid, "The 31st of every month at 09:00, rolling into the next period in months that have no such date"},
	} {
		if got := Describe(byDate, tc.policy); got != tc.want {
			t.Errorf("Describe(%s) = %q, want %q", tc.policy, got, tc.want)
		}
	}

	// A rule that cannot miss a period reads exactly as it always has — the
	// clause must not appear where it would be noise.
	plain, err := Parse("weekdays at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	if got := Describe(plain, domain.MissingDateLastValid); got != plain.HumanSummary {
		t.Errorf("Describe added a clause to a date-less rule: %q", got)
	}

	// The ordinal-weekday rule is the one whose old label was actually false.
	ordinal, err := Parse("5th friday monthly at 09:00", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	if got := Describe(ordinal, domain.MissingDateSkip); got == ordinal.HumanSummary {
		t.Errorf("Describe left the false 'every month' label intact: %q", got)
	}
}
