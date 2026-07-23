package schedule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// anchorTOD is a parsed "starting at"/"from" anchor time-of-day for a sub-daily
// interval schedule. It is resolved to a concrete UTC instant in finish.
type anchorTOD struct{ h, mi int }

// Parse turns a human-readable schedule phrase into a recurring Schedule with an
// RRULE, anchor, and plain-language summary. It never requires cron syntax.
//
// Supported forms (case-insensitive):
//
//	every <N> <unit>            e.g. "every 15 minutes", "every 30s", "every 2 hours"
//	every <unit>               e.g. "every day", "every week"
//	... [at <time>]            day-or-coarser rules accept a time-of-day
//	every <N> <sub-daily> [(starting at|from) <time>]  e.g. "every 15 minutes starting at 09:00"
//	weekdays|weekends [at ...]  e.g. "weekdays at 09:00"
//	every <weekday> [at ...]    e.g. "every monday at 9am"
//	<ordinal> <weekday> monthly e.g. "3rd wednesday monthly at 14:00", "last friday of the month"
//
// For sub-daily intervals an optional "starting at"/"from" clause sets the first-cycle anchor so
// the interval aligns to a chosen phase (e.g. :00/:15/:30/:45) instead of the creation moment.
func Parse(input, tzName string, now time.Time) (domain.Schedule, error) {
	s := strings.ToLower(strings.TrimSpace(input))
	if s == "" {
		return domain.Schedule{}, fmt.Errorf("schedule: empty schedule expression")
	}
	// The trimmed input is retained on the returned Schedule so a client can
	// show the user their own phrase again when editing (domain.Schedule.
	// Expression). It never feeds back into evaluation.
	expr := strings.TrimSpace(input)

	if sch, ok, err := parseOrdinal(s); ok || err != nil {
		return finish(sch, tzName, now, nil, expr, err)
	}
	if sch, ok, err := parseByDate(s); ok || err != nil {
		return finish(sch, tzName, now, nil, expr, err)
	}
	if sch, ok, err := parseYearly(s); ok || err != nil {
		return finish(sch, tzName, now, nil, expr, err)
	}
	if sch, ok, err := parseDayset(s); ok || err != nil {
		return finish(sch, tzName, now, nil, expr, err)
	}
	if sch, ok, err := parseEveryWeekday(s); ok || err != nil {
		return finish(sch, tzName, now, nil, expr, err)
	}
	if sch, anchor, ok, err := parseInterval(s); ok || err != nil {
		return finish(sch, tzName, now, anchor, expr, err)
	}
	return domain.Schedule{}, fmt.Errorf("schedule: could not understand %q (try forms like \"every 15 minutes\", \"weekdays at 09:00\", \"3rd wednesday monthly at 14:00\")", input)
}

// finish validates the constructed RRULE, sets anchor/kind, and returns. When anchor is non-nil
// (an explicit "starting at"/"from" clause), the anchor instant is that wall time in the task
// timezone on the current day; otherwise the anchor defaults to now (creation-aligned).
// expr is the user's original phrase, retained for round-tripping only.
func finish(sch domain.Schedule, tzName string, now time.Time, anchor *anchorTOD, expr string, err error) (domain.Schedule, error) {
	if err != nil {
		return domain.Schedule{}, err
	}
	if _, perr := rrule.StrToROption(sch.RRULE); perr != nil {
		return domain.Schedule{}, fmt.Errorf("schedule: built invalid rule %q: %w", sch.RRULE, perr)
	}
	sch.Kind = domain.ScheduleRecurring
	sch.Expression = expr
	if anchor == nil {
		a := now.UTC()
		sch.Anchor = &a
		return sch, nil
	}
	loc, lerr := timezone.Resolve(tzName)
	if lerr != nil {
		return domain.Schedule{}, lerr
	}
	n := now.In(loc)
	a := time.Date(n.Year(), n.Month(), n.Day(), anchor.h, anchor.mi, 0, 0, loc).UTC()
	sch.Anchor = &a
	return sch, nil
}

var (
	reInterval = regexp.MustCompile(`^every\s+(?:(\d+)\s*)?(second|seconds|sec|secs|s|minute|minutes|min|mins|m|hour|hours|hr|hrs|h|day|days|d|week|weeks|w|month|months|mo|year|years|yr|yrs|y)(?:\s+(at|starting\s+at|from)\s+(.+))?$`)
	reDayset   = regexp.MustCompile(`^(weekdays|weekends)(?:\s+at\s+(.+))?$`)
	reEveryDay = regexp.MustCompile(`^every\s+(monday|tuesday|wednesday|thursday|friday|saturday|sunday)(?:\s+at\s+(.+))?$`)
	reOrdinal  = regexp.MustCompile(`^(1st|2nd|3rd|4th|5th|last|first|second|third|fourth|fifth)\s+(monday|tuesday|wednesday|thursday|friday|saturday|sunday)\s+(?:of\s+(?:the|each|every)\s+month|monthly)(?:\s+at\s+(.+))?$`)
	// By-date monthly: "on the 15th of every month", "the 31st monthly".
	reByDate = regexp.MustCompile(`^(?:on\s+)?the\s+(\d{1,2})(?:st|nd|rd|th)\s+(?:of\s+(?:the|each|every)\s+month|monthly)(?:\s+at\s+(.+))?$`)
	// Yearly by date: "every year on february 29", "annually on 29 february".
	reYearly = regexp.MustCompile(`^(?:every\s+year|annually|yearly)\s+on\s+(?:([a-z]+)\s+(\d{1,2})|(\d{1,2})\s+([a-z]+))(?:\s+at\s+(.+))?$`)
)

// monthNumber maps a month name or three-letter abbreviation to its number.
var monthNumber = map[string]int{
	"january": 1, "jan": 1, "february": 2, "feb": 2, "march": 3, "mar": 3,
	"april": 4, "apr": 4, "may": 5, "june": 6, "jun": 6, "july": 7, "jul": 7,
	"august": 8, "aug": 8, "september": 9, "sep": 9, "sept": 9, "october": 10, "oct": 10,
	"november": 11, "nov": 11, "december": 12, "dec": 12,
}

var monthTitle = [...]string{
	"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

// daysInMonth is the greatest valid day for each month, taking February as 29 so
// a leap-day rule is expressible. Whether it fires in a common year is the
// missing-date policy's business, not the grammar's.
var daysInMonth = [...]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

var weekdayCode = map[string]string{
	"monday": "MO", "tuesday": "TU", "wednesday": "WE", "thursday": "TH",
	"friday": "FR", "saturday": "SA", "sunday": "SU",
}

var weekdayTitle = map[string]string{
	"monday": "Monday", "tuesday": "Tuesday", "wednesday": "Wednesday", "thursday": "Thursday",
	"friday": "Friday", "saturday": "Saturday", "sunday": "Sunday",
}

var ordinalNum = map[string]int{
	"1st": 1, "first": 1, "2nd": 2, "second": 2, "3rd": 3, "third": 3,
	"4th": 4, "fourth": 4, "5th": 5, "fifth": 5, "last": -1,
}

// IsSubDailyInterval reports whether input is a fixed-interval schedule with a
// sub-daily unit (seconds/minutes/hours) — i.e. one that accepts an optional
// "starting at"/"from" anchor. It is used by the GUI to decide whether to offer
// the anchor input.
func IsSubDailyInterval(input string) bool {
	m := reInterval.FindStringSubmatch(strings.ToLower(strings.TrimSpace(input)))
	if m == nil {
		return false
	}
	_, _, subDaily := unitToFreq(m[2])
	return subDaily
}

func parseInterval(s string) (domain.Schedule, *anchorTOD, bool, error) {
	m := reInterval.FindStringSubmatch(s)
	if m == nil {
		return domain.Schedule{}, nil, false, nil
	}
	n := 1
	if m[1] != "" {
		var err error
		if n, err = strconv.Atoi(m[1]); err != nil || n < 1 {
			return domain.Schedule{}, nil, true, fmt.Errorf("schedule: invalid interval %q", m[1])
		}
	}
	freq, unitName, subDaily := unitToFreq(m[2])
	keyword := m[3]
	clause := strings.TrimSpace(m[4])

	parts := []string{"FREQ=" + freq, "INTERVAL=" + strconv.Itoa(n)}
	summary := "Every " + plural(n, unitName)

	isAnchor := keyword == "from" || strings.HasPrefix(keyword, "starting")
	switch {
	case isAnchor:
		if !subDaily {
			return domain.Schedule{}, nil, true, fmt.Errorf("schedule: 'starting at' only applies to interval schedules (seconds/minutes/hours)")
		}
		h, mi, ok := parseTimeOfDay(clause)
		if !ok {
			return domain.Schedule{}, nil, true, fmt.Errorf("schedule: invalid time-of-day %q (try 09:00, 9:00 AM, 9am)", clause)
		}
		summary += " starting at " + clock(h, mi)
		return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: summary}, &anchorTOD{h, mi}, true, nil
	case subDaily && keyword == "at":
		return domain.Schedule{}, nil, true, fmt.Errorf("schedule: %q does not support an 'at <time>' clause; use 'starting at' to set a first-run time", m[2])
	default:
		// Daily-or-coarser: an optional 'at <time>' sets the time-of-day.
		h, mi, withTime, err := maybeTime(clause)
		if err != nil {
			return domain.Schedule{}, nil, true, err
		}
		if withTime {
			parts = append(parts, byTime(h, mi)...)
			summary += " at " + clock(h, mi)
		}
	}
	return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: summary}, nil, true, nil
}

func parseDayset(s string) (domain.Schedule, bool, error) {
	m := reDayset.FindStringSubmatch(s)
	if m == nil {
		return domain.Schedule{}, false, nil
	}
	var byday, label string
	if m[1] == "weekdays" {
		byday, label = "MO,TU,WE,TH,FR", "Every weekday"
	} else {
		byday, label = "SA,SU", "Every weekend day"
	}
	parts := []string{"FREQ=WEEKLY", "BYDAY=" + byday}
	h, mi, withTime, err := maybeTime(strings.TrimSpace(m[2]))
	if err != nil {
		return domain.Schedule{}, true, err
	}
	if withTime {
		parts = append(parts, byTime(h, mi)...)
		label += " at " + clock(h, mi)
	}
	return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: label}, true, nil
}

func parseEveryWeekday(s string) (domain.Schedule, bool, error) {
	m := reEveryDay.FindStringSubmatch(s)
	if m == nil {
		return domain.Schedule{}, false, nil
	}
	parts := []string{"FREQ=WEEKLY", "BYDAY=" + weekdayCode[m[1]]}
	label := "Every " + weekdayTitle[m[1]]
	h, mi, withTime, err := maybeTime(strings.TrimSpace(m[2]))
	if err != nil {
		return domain.Schedule{}, true, err
	}
	if withTime {
		parts = append(parts, byTime(h, mi)...)
		label += " at " + clock(h, mi)
	}
	return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: label}, true, nil
}

func parseOrdinal(s string) (domain.Schedule, bool, error) {
	m := reOrdinal.FindStringSubmatch(s)
	if m == nil {
		return domain.Schedule{}, false, nil
	}
	n := ordinalNum[m[1]]
	day := weekdayCode[m[2]]
	sign := "+"
	if n < 0 {
		sign = ""
	}
	parts := []string{"FREQ=MONTHLY", fmt.Sprintf("BYDAY=%s%d%s", sign, n, day)}
	label := fmt.Sprintf("The %s %s of every month", ordinalWord(n), weekdayTitle[m[2]])
	h, mi, withTime, err := maybeTime(strings.TrimSpace(m[3]))
	if err != nil {
		return domain.Schedule{}, true, err
	}
	if withTime {
		parts = append(parts, byTime(h, mi)...)
		label += " at " + clock(h, mi)
	}
	return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: label}, true, nil
}

// parseByDate handles a monthly rule addressed by calendar date, e.g. "on the
// 15th of every month at 09:00". A day of 29–31 is accepted: whether it fires in
// a month that has no such date is the task's missing-date policy to decide, and
// the summary says which way that falls.
func parseByDate(s string) (domain.Schedule, bool, error) {
	m := reByDate.FindStringSubmatch(s)
	if m == nil {
		return domain.Schedule{}, false, nil
	}
	day, err := strconv.Atoi(m[1])
	if err != nil || day < 1 || day > 31 {
		return domain.Schedule{}, true, fmt.Errorf("schedule: day of month must be between 1 and 31, got %q", m[1])
	}
	parts := []string{"FREQ=MONTHLY", "BYMONTHDAY=" + strconv.Itoa(day)}
	label := fmt.Sprintf("The %s of every month", ordinalWord(day))
	h, mi, withTime, terr := maybeTime(strings.TrimSpace(m[2]))
	if terr != nil {
		return domain.Schedule{}, true, terr
	}
	if withTime {
		parts = append(parts, byTime(h, mi)...)
		label += " at " + clock(h, mi)
	}
	return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: label}, true, nil
}

// parseYearly handles a rule addressed by month and date once a year, e.g.
// "every year on february 29 at 09:00". Both orderings of month and day are
// accepted because both read naturally and neither is ambiguous.
func parseYearly(s string) (domain.Schedule, bool, error) {
	m := reYearly.FindStringSubmatch(s)
	if m == nil {
		return domain.Schedule{}, false, nil
	}
	name, dayStr := m[1], m[2]
	if name == "" {
		name, dayStr = m[4], m[3]
	}
	month, ok := monthNumber[name]
	if !ok {
		return domain.Schedule{}, true, fmt.Errorf("schedule: %q is not a month name", name)
	}
	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > daysInMonth[month] {
		return domain.Schedule{}, true, fmt.Errorf("schedule: %s has no day %s", monthTitle[month], dayStr)
	}
	parts := []string{"FREQ=YEARLY", "BYMONTH=" + strconv.Itoa(month), "BYMONTHDAY=" + strconv.Itoa(day)}
	label := fmt.Sprintf("Every year on %s %d", monthTitle[month], day)
	h, mi, withTime, terr := maybeTime(strings.TrimSpace(m[5]))
	if terr != nil {
		return domain.Schedule{}, true, terr
	}
	if withTime {
		parts = append(parts, byTime(h, mi)...)
		label += " at " + clock(h, mi)
	}
	return domain.Schedule{RRULE: strings.Join(parts, ";"), HumanSummary: label}, true, nil
}

// ---- helpers ------------------------------------------------------------

func unitToFreq(u string) (freq, name string, subDaily bool) {
	switch u {
	case "second", "seconds", "sec", "secs", "s":
		return "SECONDLY", "second", true
	case "minute", "minutes", "min", "mins", "m":
		return "MINUTELY", "minute", true
	case "hour", "hours", "hr", "hrs", "h":
		return "HOURLY", "hour", true
	case "day", "days", "d":
		return "DAILY", "day", false
	case "week", "weeks", "w":
		return "WEEKLY", "week", false
	case "month", "months", "mo":
		return "MONTHLY", "month", false
	case "year", "years", "yr", "yrs", "y":
		return "YEARLY", "year", false
	}
	return "", "", false
}

func byTime(h, mi int) []string {
	return []string{"BYHOUR=" + strconv.Itoa(h), "BYMINUTE=" + strconv.Itoa(mi), "BYSECOND=0"}
}

// maybeTime parses an optional time-of-day clause. Returns withTime=false when
// the clause is empty.
func maybeTime(s string) (h, mi int, withTime bool, err error) {
	if s == "" {
		return 0, 0, false, nil
	}
	h, mi, ok := parseTimeOfDay(s)
	if !ok {
		return 0, 0, false, fmt.Errorf("schedule: invalid time-of-day %q (try 09:00, 9:00 AM, 9am)", s)
	}
	return h, mi, true, nil
}

var reTOD = regexp.MustCompile(`^(\d{1,2})(?::(\d{2}))?\s*(am|pm)?$`)

// parseTimeOfDay accepts "14:00", "9:00", "9:00 am", "9am", "9".
func parseTimeOfDay(s string) (h, mi int, ok bool) {
	m := reTOD.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, 0, false
	}
	h, _ = strconv.Atoi(m[1])
	if m[2] != "" {
		mi, _ = strconv.Atoi(m[2])
	}
	switch m[3] {
	case "am":
		if h == 12 {
			h = 0
		}
	case "pm":
		if h != 12 {
			h += 12
		}
	}
	if h > 23 || mi > 59 {
		return 0, 0, false
	}
	return h, mi, true
}

func plural(n int, unit string) string {
	if n == 1 {
		return unit
	}
	return strconv.Itoa(n) + " " + unit + "s"
}

func clock(h, mi int) string { return fmt.Sprintf("%02d:%02d", h, mi) }

// ordinalWord renders an ordinal in English. The by-date monthly form reaches
// days up to 31, so the suffix is computed rather than tabulated: 21st and 31st
// take "st" while 11th, 12th and 13th take "th" despite their final digit.
func ordinalWord(n int) string {
	if n == -1 {
		return "last"
	}
	suffix := "th"
	if n%100 < 11 || n%100 > 13 {
		switch n % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return strconv.Itoa(n) + suffix
}
