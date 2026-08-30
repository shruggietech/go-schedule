// Package schedule compiles normalized human-readable scheduling intent into a
// stored representation (RFC 5545 RRULE for recurrence, or a single instant for
// one-off) and computes concrete next-run times in UTC. Source-syntax selection
// is owned by internal/scheduleinput; RRULE, an optional typed calendar
// adjustment, and anchor remain authoritative.
package schedule

import (
	"fmt"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// NewOneOff builds a one-off schedule that fires once at runAt (stored UTC).
func NewOneOff(runAt time.Time) domain.Schedule {
	u := runAt.UTC()
	return domain.Schedule{
		Kind:         domain.ScheduleOneOff,
		RunAt:        &u,
		HumanSummary: "Once at " + u.Format("2006-01-02 15:04 MST"),
	}
}

// NextRun returns the next run instant (UTC) strictly after `after` for the
// schedule evaluated in timezone tzName. The bool is false when there is no
// further run (exhausted one-off, or event schedule with no time component).
//
// policy decides what a date-bearing recurrence does in a period that has no
// matching date; it is inert for every other schedule shape. Pass
// domain.MissingDateSkip (or the zero value) for the historical behavior.
func NextRun(sch domain.Schedule, tzName string, policy domain.MissingDatePolicy, after time.Time) (time.Time, bool, error) {
	return NextRunWithPolicy(sch, tzName, domain.SchedulePolicy{MissingDate: policy}, after)
}

// NextRunWithPolicy returns the next run under the task's complete scheduling
// policy. Zero policy values use the compatibility defaults.
func NextRunWithPolicy(sch domain.Schedule, tzName string, policy domain.SchedulePolicy, after time.Time) (time.Time, bool, error) {
	policy = policy.Effective()
	if err := ValidatePolicy(sch, policy); err != nil {
		return time.Time{}, false, err
	}
	switch sch.Kind {
	case domain.ScheduleOneOff:
		if sch.RunAt == nil {
			return time.Time{}, false, fmt.Errorf("schedule: one-off missing run_at")
		}
		if sch.RunAt.After(after) {
			return sch.RunAt.UTC(), true, nil
		}
		return time.Time{}, false, nil
	case domain.ScheduleEvent:
		return time.Time{}, false, nil
	case domain.ScheduleRecurring:
		switch policy.TimeBasis {
		case domain.TimeBasisElapsed:
			return nextElapsed(sch, tzName, after)
		case domain.TimeBasisUTC:
			policy.TimeBasis = domain.TimeBasisWallClock
			return nextRecurring(sch, "UTC", policy, after)
		default:
			return nextRecurring(sch, tzName, policy, after)
		}
	default:
		return time.Time{}, false, fmt.Errorf("schedule: unknown kind %q", sch.Kind)
	}
}

// ValidatePolicy validates enum values and schedule/basis compatibility.
func ValidatePolicy(sch domain.Schedule, policy domain.SchedulePolicy) error {
	policy = policy.Effective()
	switch policy.MissingDate {
	case domain.MissingDateSkip, domain.MissingDateLastValid, domain.MissingDateNextValid:
	default:
		return fmt.Errorf("schedule: invalid missing-date policy %q", policy.MissingDate)
	}
	switch policy.TimeBasis {
	case domain.TimeBasisWallClock, domain.TimeBasisElapsed, domain.TimeBasisUTC:
	default:
		return fmt.Errorf("schedule: invalid time basis %q", policy.TimeBasis)
	}
	switch policy.DSTGap {
	case domain.DSTGapNextValid, domain.DSTGapSkip:
	default:
		return fmt.Errorf("schedule: invalid DST gap policy %q", policy.DSTGap)
	}
	switch policy.DSTOverlap {
	case domain.DSTOverlapFirst, domain.DSTOverlapBoth, domain.DSTOverlapLast:
	default:
		return fmt.Errorf("schedule: invalid DST overlap policy %q", policy.DSTOverlap)
	}
	if policy.TimeBasis == domain.TimeBasisElapsed && sch.Kind == domain.ScheduleRecurring {
		if _, _, err := fixedDuration(sch); err != nil {
			return err
		}
	}
	return nil
}

// PrepareForPolicy binds any durable schedule state required by policy. An
// elapsed schedule receives one absolute epoch, derived in the authoring
// timezone exactly once; subsequent timezone edits cannot move that phase.
func PrepareForPolicy(sch *domain.Schedule, tzName string, policy domain.SchedulePolicy) error {
	if err := ValidatePolicy(*sch, policy); err != nil {
		return err
	}
	policy = policy.Effective()
	if policy.TimeBasis != domain.TimeBasisElapsed || sch.Kind != domain.ScheduleRecurring || sch.ElapsedEpoch != nil {
		return nil
	}
	epoch, err := deriveElapsedEpoch(*sch, tzName)
	if err != nil {
		return err
	}
	sch.ElapsedEpoch = &epoch
	return nil
}

func nextRecurring(sch domain.Schedule, tzName string, policy domain.SchedulePolicy, after time.Time) (time.Time, bool, error) {
	loc, err := timezone.Resolve(tzName)
	if err != nil {
		return time.Time{}, false, err
	}
	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("schedule: parse rrule %q: %w", sch.RRULE, err)
	}
	anchor := after
	if sch.Anchor != nil {
		anchor = *sch.Anchor
	}
	opt.Dtstart = anchor.In(loc)

	switch sch.CalendarAdjustment {
	case "":
	case domain.CalendarAdjustmentNearestWeekday:
		occ, ok, err := resolveNearestWeekday(opt, loc, policy, after)
		if err != nil || !ok {
			return time.Time{}, ok, err
		}
		return occ, true, nil
	default:
		return time.Time{}, false, fmt.Errorf("schedule: unknown calendar adjustment %q", sch.CalendarAdjustment)
	}

	if compositeDailySet(opt) {
		if policy.MissingDate != domain.MissingDateSkip && compositeDateSet(opt) {
			occ, ok := resolveCompositeDateSet(opt, loc, policy, after)
			if !ok {
				return time.Time{}, false, nil
			}
			return occ, true, nil
		}
		if len(opt.Bysecond) != 1 || opt.Bysecond[0] != 0 {
			occ, ok := resolveCompositeDailySet(opt, loc, policy, after)
			return occ, ok, nil
		}
	}

	// A date-bearing rule under a non-skip policy is resolved by walking periods
	// rather than by asking rrule-go, because rrule-go's answer is precisely the
	// "skip" answer: it omits periods that have no matching date. The walk
	// returns a wall-clock occurrence in loc, which then takes the same DST
	// normalization below as any other day-or-coarser occurrence.
	// The policy comparison comes first so the default path does no extra work
	// at all: it is a string comparison against a constant, and everything below
	// it — including inspecting the rule for a date intent — is skipped entirely
	// for the schedules that make up the overwhelming majority of the hot path.
	if policy.MissingDate != domain.MissingDateSkip {
		if intent, ok := dateIntent(opt); ok {
			occ, ok := resolveMissingDate(intent, opt, loc, policy, after)
			if !ok {
				return time.Time{}, false, nil
			}
			return occ, true, nil
		}
	}

	return nextFloatingWall(opt, loc, policy, after)
}

// UpcomingRuns returns up to n future run instants (UTC) after `after`.
func UpcomingRuns(sch domain.Schedule, tzName string, policy domain.MissingDatePolicy, after time.Time, n int) ([]time.Time, error) {
	return UpcomingRunsWithPolicy(sch, tzName, domain.SchedulePolicy{MissingDate: policy}, after, n)
}

// UpcomingRunsWithPolicy returns up to n future UTC instants using policy.
func UpcomingRunsWithPolicy(sch domain.Schedule, tzName string, policy domain.SchedulePolicy, after time.Time, n int) ([]time.Time, error) {
	var out []time.Time
	cursor := after
	for i := 0; i < n; i++ {
		next, ok, err := NextRunWithPolicy(sch, tzName, policy, cursor)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		out = append(out, next)
		cursor = next
	}
	return out, nil
}

func nextElapsed(sch domain.Schedule, tzName string, after time.Time) (time.Time, bool, error) {
	duration, _, err := fixedDuration(sch)
	if err != nil {
		return time.Time{}, false, err
	}
	var anchor time.Time
	if sch.ElapsedEpoch != nil {
		anchor = sch.ElapsedEpoch.UTC()
	} else {
		// Compatibility for in-memory and pre-policy schedules. Persisted elapsed
		// schedules are prepared at their mutation boundary and take the branch
		// above, so task timezone changes never reach this derivation.
		anchor, err = deriveElapsedEpoch(sch, tzName)
		if err != nil {
			return time.Time{}, false, err
		}
	}
	if after.Before(anchor) {
		return anchor, true, nil
	}
	steps := after.Sub(anchor)/duration + 1
	return anchor.Add(steps * duration).UTC(), true, nil
}

func deriveElapsedEpoch(sch domain.Schedule, tzName string) (time.Time, error) {
	_, opt, err := fixedDuration(sch)
	if err != nil {
		return time.Time{}, err
	}
	loc, err := timezone.Resolve(tzName)
	if err != nil {
		return time.Time{}, err
	}
	opt.Dtstart = sch.Anchor.In(loc)
	r, err := rrule.NewRRule(*opt)
	if err != nil {
		return time.Time{}, fmt.Errorf("schedule: build elapsed rule: %w", err)
	}
	epoch := r.After(sch.Anchor.In(loc).Add(-time.Nanosecond), false)
	if epoch.IsZero() {
		return time.Time{}, fmt.Errorf("schedule: elapsed rule has no epoch")
	}
	return epoch.UTC(), nil
}

func fixedDuration(sch domain.Schedule) (time.Duration, *rrule.ROption, error) {
	if sch.Anchor == nil {
		return 0, nil, fmt.Errorf("schedule: elapsed basis requires an anchor")
	}
	if sch.CalendarAdjustment != "" {
		return 0, nil, fmt.Errorf("schedule: elapsed basis requires a fixed-duration interval")
	}
	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		return 0, nil, fmt.Errorf("schedule: parse rrule %q: %w", sch.RRULE, err)
	}
	if len(opt.Bymonth) != 0 || len(opt.Bymonthday) != 0 || len(opt.Byweekday) != 0 ||
		len(opt.Byyearday) != 0 || len(opt.Byweekno) != 0 || len(opt.Bysetpos) != 0 ||
		len(opt.Byhour) > 1 || len(opt.Byminute) > 1 || len(opt.Bysecond) > 1 ||
		opt.Count != 0 || !opt.Until.IsZero() {
		return 0, nil, fmt.Errorf("schedule: elapsed basis requires one fixed-duration interval occurrence")
	}
	interval := opt.Interval
	if interval < 1 {
		interval = 1
	}
	var unit time.Duration
	switch opt.Freq {
	case rrule.SECONDLY:
		unit = time.Second
	case rrule.MINUTELY:
		unit = time.Minute
	case rrule.HOURLY:
		unit = time.Hour
	case rrule.DAILY:
		unit = 24 * time.Hour
	case rrule.WEEKLY:
		unit = 7 * 24 * time.Hour
	default:
		return 0, nil, fmt.Errorf("schedule: elapsed basis requires a fixed-duration interval")
	}
	return time.Duration(interval) * unit, opt, nil
}

func nextFloatingWall(opt *rrule.ROption, loc *time.Location, policy domain.SchedulePolicy, after time.Time) (time.Time, bool, error) {
	anchor := opt.Dtstart.In(loc)
	opt.Dtstart = floating(anchor)
	r, err := rrule.NewRRule(*opt)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("schedule: build rrule: %w", err)
	}
	localCursor := floating(after.In(loc))
	search := localCursor.Add(-time.Nanosecond)
	best := nextRepeatedCandidate(r, loc, policy, after)
	for i := 0; i < 200000; i++ {
		intent := r.After(search, false)
		if intent.IsZero() {
			break
		}
		var next time.Time
		for _, candidate := range timezone.ResolveWallTime(loc, intent.Year(), intent.Month(), intent.Day(), intent.Hour(), intent.Minute(), intent.Second(), policy.DSTGap, policy.DSTOverlap) {
			candidate = candidate.UTC()
			if candidate.After(after) && (next.IsZero() || candidate.Before(next)) {
				next = candidate
			}
		}
		if !next.IsZero() {
			if best.IsZero() || next.Before(best) {
				return next, true, nil
			}
			return best, true, nil
		}
		search = intent
	}
	if !best.IsZero() {
		return best, true, nil
	}
	return time.Time{}, false, nil
}

// nextRepeatedCandidate finds the earliest second-fold occurrence that can
// still follow after. It jumps directly into the actual repeated wall interval
// instead of enumerating up to 52 hours of dense recurrence intents.
func nextRepeatedCandidate(r *rrule.RRule, loc *time.Location, policy domain.SchedulePolicy, after time.Time) time.Time {
	if policy.DSTOverlap != domain.DSTOverlapBoth && policy.DSTOverlap != domain.DSTOverlapLast {
		return time.Time{}
	}
	start, end, newOffset, ok := overlapWindow(loc, after)
	if !ok {
		return time.Time{}
	}
	threshold := floating(after.UTC().Add(time.Duration(newOffset) * time.Second))
	search := start.Add(-time.Nanosecond)
	if threshold.After(search) {
		search = threshold
	}
	intent := r.After(search, false)
	if intent.IsZero() || !intent.Before(end) {
		return time.Time{}
	}
	var best time.Time
	for _, candidate := range timezone.ResolveWallTime(loc, intent.Year(), intent.Month(), intent.Day(), intent.Hour(), intent.Minute(), intent.Second(), policy.DSTGap, policy.DSTOverlap) {
		candidate = candidate.UTC()
		if candidate.After(after) && (best.IsZero() || candidate.Before(best)) {
			best = candidate
		}
	}
	return best
}

// overlapWindow returns the repeated floating wall interval around a nearby
// backward offset transition, plus the offset used by the second fold.
func overlapWindow(loc *time.Location, around time.Time) (time.Time, time.Time, int, bool) {
	left := around.Add(-26 * time.Hour).Truncate(time.Second)
	right := around.Add(26 * time.Hour).Truncate(time.Second)
	_, previousOffset := left.In(loc).Zone()
	previousProbe := left
	for probe := left.Add(time.Hour); !probe.After(right); probe = probe.Add(time.Hour) {
		_, offset := probe.In(loc).Zone()
		if offset < previousOffset {
			lo, hi := previousProbe.Unix(), probe.Unix()
			for lo+1 < hi {
				mid := lo + (hi-lo)/2
				_, candidateOffset := time.Unix(mid, 0).In(loc).Zone()
				if candidateOffset < previousOffset {
					hi = mid
				} else {
					lo = mid
				}
			}
			transition := time.Unix(hi, 0).In(loc)
			_, newOffset := transition.Zone()
			start := floating(transition)
			end := start.Add(time.Duration(previousOffset-newOffset) * time.Second)
			return start, end, newOffset, true
		}
		previousOffset = offset
		previousProbe = probe
	}
	return time.Time{}, time.Time{}, 0, false
}

func floating(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}
