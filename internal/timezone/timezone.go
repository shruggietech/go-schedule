// Package timezone resolves IANA timezones and converts intended local
// wall-clock times into concrete instants, applying the project's Daylight
// Saving Time rules: a time that falls in a skipped hour (spring-forward) runs
// at the next valid instant; a time in a repeated hour (fall-back) runs once,
// on the first occurrence. Storage and scheduling use UTC throughout.
package timezone

import (
	"fmt"
	"sort"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Resolve returns the *time.Location for an IANA name. "Local" and "" map to the
// host's local zone.
func Resolve(name string) (*time.Location, error) {
	if name == "" || name == "Local" {
		return time.Local, nil
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("timezone: %q is not a valid IANA zone: %w", name, err)
	}
	return loc, nil
}

// WallTime resolves an intended local wall-clock time (y/mo/d h:mi:s in loc) to
// a concrete instant located in loc, applying the DST rules described in the
// package doc. Callers convert the result to UTC for storage/dispatch.
func WallTime(loc *time.Location, y int, mo time.Month, d, h, mi, s int) time.Time {
	resolved := ResolveWallTime(loc, y, mo, d, h, mi, s, domain.DSTGapNextValid, domain.DSTOverlapFirst)
	if len(resolved) != 0 {
		return resolved[0]
	}
	return time.Date(y, mo, d, h, mi, s, 0, loc)
}

// ResolveWallTime maps one intended local reading to zero, one, or two concrete
// instants according to the selected DST policies. It discovers the offsets on
// either side of the date instead of assuming a one-hour transition.
func ResolveWallTime(loc *time.Location, y int, mo time.Month, d, h, mi, s int, gap domain.DSTGapPolicy, overlap domain.DSTOverlapPolicy) []time.Time {
	gap = effectiveGap(gap)
	overlap = effectiveOverlap(overlap)
	requested := time.Date(y, mo, d, h, mi, s, 0, time.UTC)
	if exact := exactWallInstants(loc, requested); len(exact) != 0 {
		return selectOverlap(exact, overlap)
	}
	if gap == domain.DSTGapSkip {
		return nil
	}
	// Preserve the historical minute-granularity next-valid behavior. Schedule
	// occurrences are second-aligned, so the requested second remains stable.
	for add := 1; add <= 24*60; add++ {
		if exact := exactWallInstants(loc, requested.Add(time.Duration(add)*time.Minute)); len(exact) != 0 {
			return selectOverlap(exact, overlap)
		}
	}
	return nil
}

func exactWallInstants(loc *time.Location, intended time.Time) []time.Time {
	approx := time.Date(intended.Year(), intended.Month(), intended.Day(), intended.Hour(), intended.Minute(), intended.Second(), intended.Nanosecond(), loc)
	offsets := make(map[int]struct{}, 5)
	for _, delta := range []time.Duration{-48 * time.Hour, -24 * time.Hour, 0, 24 * time.Hour, 48 * time.Hour} {
		_, offset := approx.Add(delta).Zone()
		offsets[offset] = struct{}{}
	}
	instants := make([]time.Time, 0, 2)
	for offset := range offsets {
		candidate := intended.Add(-time.Duration(offset) * time.Second).In(loc)
		if wallMatches(candidate, intended) {
			instants = append(instants, candidate)
		}
	}
	sort.Slice(instants, func(i, j int) bool { return instants[i].Before(instants[j]) })
	return instants
}

func selectOverlap(instants []time.Time, policy domain.DSTOverlapPolicy) []time.Time {
	if len(instants) < 2 || policy == domain.DSTOverlapBoth {
		return instants
	}
	if policy == domain.DSTOverlapLast {
		return instants[len(instants)-1:]
	}
	return instants[:1]
}

func effectiveGap(p domain.DSTGapPolicy) domain.DSTGapPolicy {
	if p == "" {
		return domain.DSTGapNextValid
	}
	return p
}

func effectiveOverlap(p domain.DSTOverlapPolicy) domain.DSTOverlapPolicy {
	if p == "" {
		return domain.DSTOverlapFirst
	}
	return p
}

func wallMatches(cand, ref time.Time) bool {
	return cand.Year() == ref.Year() && cand.Month() == ref.Month() && cand.Day() == ref.Day() &&
		cand.Hour() == ref.Hour() && cand.Minute() == ref.Minute() && cand.Second() == ref.Second()
}
