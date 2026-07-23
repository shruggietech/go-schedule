package schedule

import (
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// This file resolves recurrences whose target date does not exist in every
// period — the 31st in a 30-day month, 29 February in a common year, the fifth
// Friday in a month with only four.
//
// rrule-go already implements one of the three policies: it omits such periods
// entirely, which is exactly domain.MissingDateSkip. That path is left
// untouched, so the default behavior is not a reimplementation of the old
// behavior but literally the old code. The other two policies need the *intent*
// (which day of the month, or which ordinal weekday, the operator asked for),
// which a generated occurrence no longer carries — so they are resolved here by
// walking periods and computing the target date directly.

// Describe renders a schedule's human summary with its missing-date policy named
// — but only when the rule can actually miss a period, so ordinary schedules read
// exactly as they always have.
//
// It is deliberately a render-time function rather than something baked into the
// stored summary: the policy lives on the task and can change without the phrase
// changing, so a stored sentence naming the policy would go stale the moment an
// operator switched it. That staleness is the exact defect this discharges —
// "The 5th Friday of every month" for a rule that fires four times a year
// (FR-023, SC-005).
func Describe(sch domain.Schedule, policy domain.MissingDatePolicy) string {
	if sch.Kind != domain.ScheduleRecurring || sch.HumanSummary == "" {
		return sch.HumanSummary
	}
	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		return sch.HumanSummary
	}
	in, ok := dateIntent(opt)
	if !ok {
		return sch.HumanSummary
	}
	return sch.HumanSummary + ", " + policyClause(in, policy)
}

// policyClause is the sentence fragment naming what happens in a period with no
// matching date.
func policyClause(in targetDate, policy domain.MissingDatePolicy) string {
	period := "months"
	if in.yearly {
		period = "years"
	}
	switch policy {
	case domain.MissingDateLastValid:
		if in.kind == intentOrdinal {
			return "or the last one when there is none"
		}
		return "or the last day of the month when there is no such date"
	case domain.MissingDateNextValid:
		return "rolling into the next period in " + period + " that have no such date"
	default:
		return "skipped in " + period + " that have no such date"
	}
}

// maxPeriodWalk bounds the search for the next occurrence. A monthly rule needs
// at most a couple of periods and a yearly rule at most a few, so this is a
// backstop against a malformed rule spinning the dispatch loop rather than a
// working limit: exhausting it reports "no further run", the same as a genuinely
// exhausted recurrence.
const maxPeriodWalk = 500

// intentKind distinguishes the two shapes of rule that can miss a period.
type intentKind int

const (
	intentMonthDay intentKind = iota // BYMONTHDAY=<n>
	intentOrdinal                    // BYDAY=+<n><weekday>
)

// targetDate is the operator's target within a period, recovered from the rule.
type targetDate struct {
	kind intentKind
	// yearly is true when the rule advances by years rather than months.
	yearly bool
	// interval is the rule's INTERVAL (periods to advance per step).
	interval int
	// month is the target month for a yearly rule; zero for a monthly rule.
	month time.Month
	// day is the target day-of-month for intentMonthDay.
	day int
	// nth and weekday describe intentOrdinal (nth is 1-5; negative ordinals are
	// "last", which can never miss a period and so never reach here).
	nth     int
	weekday time.Weekday
}

// dateIntent recovers the target from a parsed rule, reporting false whenever
// the rule *cannot* miss a period. That covers more shapes than it first
// appears:
//
//   - interval and plain weekday rules, which address no date at all;
//   - an ordinal counting back from the end of the month ("last Friday"), which
//     every month has;
//   - days 1–28, which every month has;
//   - ordinals 1–4, since every month contains at least four of each weekday.
//
// Only the 29th–31st and a fifth weekday can be absent. Excluding the rest is
// what makes the policy inert rather than merely harmless (FR-024) — and it is
// also what keeps "The 3rd Wednesday of every month" from being annotated with a
// policy that could never apply to it.
func dateIntent(opt *rrule.ROption) (targetDate, bool) {
	interval := opt.Interval
	if interval < 1 {
		interval = 1
	}
	switch opt.Freq {
	case rrule.MONTHLY, rrule.YEARLY:
	default:
		return targetDate{}, false
	}
	yearly := opt.Freq == rrule.YEARLY

	if len(opt.Bymonthday) == 1 && opt.Bymonthday[0] > 28 {
		in := targetDate{kind: intentMonthDay, yearly: yearly, interval: interval, day: opt.Bymonthday[0]}
		if yearly {
			if len(opt.Bymonth) != 1 {
				return targetDate{}, false
			}
			in.month = time.Month(opt.Bymonth[0])
		}
		return in, true
	}

	// A fifth weekday. Negative ordinals ("last"), bare weekdays, and the first
	// through fourth cannot miss a period.
	if len(opt.Byweekday) == 1 && !yearly {
		wd := opt.Byweekday[0]
		n := wd.N()
		if n >= 5 {
			return targetDate{
				kind: intentOrdinal, interval: interval,
				nth: n, weekday: weekdayOf(wd),
			}, true
		}
	}
	return targetDate{}, false
}

// weekdayOf maps an rrule weekday to a time.Weekday. rrule numbers days from
// Monday=0; time numbers them from Sunday=0.
func weekdayOf(w rrule.Weekday) time.Weekday {
	return time.Weekday((w.Day() + 1) % 7)
}

// resolveMissingDate returns the first occurrence strictly after `after`,
// applying policy to periods that have no matching date. The returned time is a
// wall-clock time in loc; the caller applies DST normalization to it, so this
// function never reasons about offsets.
func resolveMissingDate(in targetDate, opt *rrule.ROption, loc *time.Location, policy domain.MissingDatePolicy, after time.Time) (time.Time, bool) {
	h, mi, s := timeOfDay(opt)
	// Start from the anchor's period so INTERVAL is counted from where the
	// operator anchored the rule, then jump forward by whole intervals to just
	// before `after`. Without the jump, a rule anchored years ago would spend
	// the whole walk budget replaying periods that are already past.
	period := periodStart(opt.Dtstart.In(loc), in.yearly, in.month, loc)
	if skip := wholeIntervalsBefore(period, after, in); skip > 0 {
		period = advanceBy(period, in, skip)
	}

	for i := 0; i < maxPeriodWalk; i++ {
		if occ, ok := occurrenceIn(in, period, h, mi, s, policy, loc); ok && occ.After(after) {
			return occ, true
		}
		period = advance(period, in)
	}
	return time.Time{}, false
}

// wholeIntervalsBefore counts how many whole INTERVAL steps separate the anchor
// period from `after`, less one so the walk starts a period early and cannot
// step over the answer. It never returns a negative count, so an anchor in the
// future is left alone.
func wholeIntervalsBefore(period, after time.Time, in targetDate) int {
	var elapsed int
	if in.yearly {
		elapsed = after.Year() - period.Year()
	} else {
		elapsed = (after.Year()-period.Year())*12 + int(after.Month()) - int(period.Month())
	}
	steps := elapsed/in.interval - 1
	if steps < 0 {
		return 0
	}
	return steps
}

// advanceBy steps n whole periods forward.
func advanceBy(period time.Time, in targetDate, n int) time.Time {
	if in.yearly {
		return period.AddDate(in.interval*n, 0, 0)
	}
	return period.AddDate(0, in.interval*n, 0)
}

// periodStart normalizes an instant to the first day of the period it belongs
// to: the month for a monthly rule, or the target month of that year for a
// yearly rule.
func periodStart(t time.Time, yearly bool, month time.Month, loc *time.Location) time.Time {
	if yearly {
		return time.Date(t.Year(), month, 1, 0, 0, 0, 0, loc)
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
}

// advance steps to the start of the next period, honoring INTERVAL.
func advance(period time.Time, in targetDate) time.Time {
	if in.yearly {
		return period.AddDate(in.interval, 0, 0)
	}
	return period.AddDate(0, in.interval, 0)
}

// occurrenceIn computes the occurrence for one period, applying the policy when
// the target date does not exist. It reports false only for the skip policy,
// which the caller does not use — kept explicit so the three-way decision is
// readable in one place.
func occurrenceIn(in targetDate, period time.Time, h, mi, s int, policy domain.MissingDatePolicy, loc *time.Location) (time.Time, bool) {
	switch in.kind {
	case intentMonthDay:
		last := daysIn(period)
		switch {
		case in.day <= last:
			return time.Date(period.Year(), period.Month(), in.day, h, mi, s, 0, loc), true
		case policy == domain.MissingDateLastValid:
			return time.Date(period.Year(), period.Month(), last, h, mi, s, 0, loc), true
		case policy == domain.MissingDateNextValid:
			// The first day of the following period. This is a distinct
			// occurrence from that period's own target date, which the walk
			// still produces on its own turn (FR-019a).
			next := period.AddDate(0, 1, 0)
			return time.Date(next.Year(), next.Month(), 1, h, mi, s, 0, loc), true
		}
		return time.Time{}, false

	case intentOrdinal:
		day, ok := nthWeekday(period, in.nth, in.weekday)
		switch {
		case ok:
			return time.Date(period.Year(), period.Month(), day, h, mi, s, 0, loc), true
		case policy == domain.MissingDateLastValid:
			// The last occurrence of that weekday in this month — what "the
			// fifth Friday, or the last one when there is none" means.
			return time.Date(period.Year(), period.Month(), lastWeekday(period, in.weekday), h, mi, s, 0, loc), true
		case policy == domain.MissingDateNextValid:
			next := period.AddDate(0, 1, 0)
			return time.Date(next.Year(), next.Month(), firstWeekday(next, in.weekday), h, mi, s, 0, loc), true
		}
		return time.Time{}, false
	}
	return time.Time{}, false
}

// daysIn returns the number of days in the month containing period.
func daysIn(period time.Time) int {
	first := time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, period.Location())
	return first.AddDate(0, 1, -1).Day()
}

// nthWeekday returns the day-of-month of the nth given weekday in period's
// month, reporting false when the month has no such occurrence.
func nthWeekday(period time.Time, nth int, wd time.Weekday) (int, bool) {
	day := firstWeekday(period, wd) + (nth-1)*7
	if day > daysIn(period) {
		return 0, false
	}
	return day, true
}

// firstWeekday returns the day-of-month of the first given weekday in period's
// month. Every month contains every weekday, so this always succeeds.
func firstWeekday(period time.Time, wd time.Weekday) int {
	first := time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, period.Location())
	return 1 + (int(wd)-int(first.Weekday())+7)%7
}

// lastWeekday returns the day-of-month of the last given weekday in period's
// month.
func lastWeekday(period time.Time, wd time.Weekday) int {
	day := firstWeekday(period, wd)
	for last := daysIn(period); day+7 <= last; {
		day += 7
	}
	return day
}

// timeOfDay recovers the rule's fixed time-of-day, falling back to the anchor's
// own clock reading when the rule carries none.
func timeOfDay(opt *rrule.ROption) (h, mi, s int) {
	h, mi, s = opt.Dtstart.Hour(), opt.Dtstart.Minute(), opt.Dtstart.Second()
	if len(opt.Byhour) == 1 {
		h = opt.Byhour[0]
	}
	if len(opt.Byminute) == 1 {
		mi = opt.Byminute[0]
	}
	if len(opt.Bysecond) == 1 {
		s = opt.Bysecond[0]
	}
	return h, mi, s
}
