package schedule

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Render reconstructs a human-readable schedule phrase from a stored schedule,
// inverting Parse. It exists for schedules persisted before the phrase was
// retained (store migration v4): those rows carry an RRULE but no Expression,
// and without a phrase an editing client has nothing to show the user.
//
// The returned phrase describes the same recurrence as sch — re-parsing it
// yields an identical RRULE — but it is canonical wording, not necessarily the
// wording the user originally typed. Prefer sch.Expression when it is non-empty.
//
// Render is deliberately partial. It covers exactly the shapes Parse can
// produce and returns "" for anything else (one-off and event schedules,
// hand-written rules, unsupported frequencies) rather than guessing: an empty
// schedule field leaves the stored schedule untouched, whereas a wrong phrase
// would silently rewrite it.
//
// It never emits a "starting at" clause. A stored Anchor cannot be
// distinguished from the creation timestamp Parse assigns when no anchor was
// given, so rendering one would put a time into the user's hands that they
// never chose. tzName is accepted for symmetry with Parse and future
// zone-dependent wording; the shapes rendered today are zone-independent.
func Render(sch domain.Schedule, tzName string) string {
	_ = tzName
	if sch.Kind != domain.ScheduleRecurring || sch.RRULE == "" {
		return ""
	}
	parts := rruleParts(sch.RRULE)
	freq := parts["FREQ"]
	if freq == "" {
		return ""
	}

	tod, hasTOD := renderTimeOfDay(parts)

	switch freq {
	case "MONTHLY":
		return renderMonthly(parts, tod, hasTOD)
	case "WEEKLY":
		if byday := parts["BYDAY"]; byday != "" {
			return renderWeeklyByDay(byday, tod, hasTOD)
		}
		return renderInterval(parts, "week", tod, hasTOD)
	case "DAILY":
		return renderInterval(parts, "day", tod, hasTOD)
	case "HOURLY":
		return renderInterval(parts, "hour", tod, hasTOD)
	case "MINUTELY":
		return renderInterval(parts, "minute", tod, hasTOD)
	case "SECONDLY":
		return renderInterval(parts, "second", tod, hasTOD)
	}
	return ""
}

// rruleParts splits an RRULE into its KEY=VALUE pairs. Keys are upper-cased;
// a malformed segment yields no entry, so callers see a missing key rather than
// a wrong one.
func rruleParts(rule string) map[string]string {
	out := map[string]string{}
	for _, seg := range strings.Split(rule, ";") {
		k, v, ok := strings.Cut(seg, "=")
		if !ok {
			continue
		}
		out[strings.ToUpper(strings.TrimSpace(k))] = strings.TrimSpace(v)
	}
	return out
}

// renderTimeOfDay rebuilds the "at HH:MM" clause from BYHOUR/BYMINUTE. Parse
// always writes both together, so a lone BYHOUR is treated as minute zero.
func renderTimeOfDay(parts map[string]string) (string, bool) {
	hs, ok := parts["BYHOUR"]
	if !ok {
		return "", false
	}
	h, err := strconv.Atoi(hs)
	if err != nil || h < 0 || h > 23 {
		return "", false
	}
	m := 0
	if ms, ok := parts["BYMINUTE"]; ok {
		if m, err = strconv.Atoi(ms); err != nil || m < 0 || m > 59 {
			return "", false
		}
	}
	return clock(h, m), true
}

// interval returns the INTERVAL value, defaulting to 1 when absent. The bool is
// false for a malformed or non-positive interval, which Parse cannot produce.
func interval(parts map[string]string) (int, bool) {
	s, ok := parts["INTERVAL"]
	if !ok {
		return 1, true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, false
	}
	return n, true
}

// renderInterval produces "every N <unit>" with an optional time-of-day, the
// inverse of parseInterval.
func renderInterval(parts map[string]string, unit string, tod string, hasTOD bool) string {
	n, ok := interval(parts)
	if !ok {
		return ""
	}
	// Sub-daily units reject an "at <time>" clause in Parse, so a rule carrying
	// one is not something Parse built; refuse rather than emit an unparseable
	// phrase.
	if hasTOD && (unit == "hour" || unit == "minute" || unit == "second") {
		return ""
	}
	out := "every " + plural(n, unit)
	if hasTOD {
		out += " at " + tod
	}
	return out
}

// renderWeeklyByDay inverts parseDayset and parseEveryWeekday. Only the three
// BYDAY shapes Parse can create are recognized.
func renderWeeklyByDay(byday, tod string, hasTOD bool) string {
	var out string
	switch byday {
	case "MO,TU,WE,TH,FR":
		out = "weekdays"
	case "SA,SU":
		out = "weekends"
	default:
		name, ok := weekdayName[byday]
		if !ok {
			return ""
		}
		out = "every " + name
	}
	if hasTOD {
		out += " at " + tod
	}
	return out
}

// renderMonthly inverts parseOrdinal: BYDAY=+3WE becomes "3rd wednesday
// monthly", BYDAY=-1FR becomes "last friday monthly".
func renderMonthly(parts map[string]string, tod string, hasTOD bool) string {
	byday := parts["BYDAY"]
	if len(byday) < 3 {
		return ""
	}
	code := byday[len(byday)-2:]
	name, ok := weekdayName[code]
	if !ok {
		return ""
	}
	n, err := strconv.Atoi(strings.TrimPrefix(byday[:len(byday)-2], "+"))
	if err != nil {
		return ""
	}
	var ordinal string
	switch {
	case n == -1:
		ordinal = "last"
	case n >= 1 && n <= 5:
		ordinal = ordinalWord(n)
	default:
		return ""
	}
	out := fmt.Sprintf("%s %s monthly", ordinal, name)
	if hasTOD {
		out += " at " + tod
	}
	return out
}

// weekdayName is the inverse of weekdayCode, in the lower-case spelling the
// parser's regexps expect.
var weekdayName = map[string]string{
	"MO": "monday", "TU": "tuesday", "WE": "wednesday", "TH": "thursday",
	"FR": "friday", "SA": "saturday", "SU": "sunday",
}
