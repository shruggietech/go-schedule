package cron

import (
	"fmt"
	"strings"
)

// Phrase renders a Spec as a concise human-readable phrase when the natural
// language grammar can express the cron schedule without losing information.
func Phrase(s Spec) (string, Unsupported, bool) {
	if len(s.Second.Values) != 1 || s.Second.Values[0] != 0 {
		return "", Unsupported{Reason: "seconds precision uses a field-complete description rather than the human schedule grammar"}, false
	}
	// Sub-hourly: "*/n" in the minute field with every hour.
	if s.Minute.Wildcard && s.Minute.Step > 1 {
		if !s.Hour.EveryValue() || !s.DOM.EveryValue() || !s.Month.EveryValue() || !s.DOW.EveryValue() {
			return "", Unsupported{Reason: "a minute step combined with other restrictions has no phrase equivalent"}, false
		}
		if 60%s.Minute.Step != 0 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"a step of %d does not divide the hour evenly: cron restarts the sequence at :00, which a fixed interval does not reproduce",
				s.Minute.Step)}, false
		}
		return fmt.Sprintf("every %d minutes starting at 00:00", s.Minute.Step), Unsupported{}, true
	}

	// Every minute.
	if s.Minute.EveryValue() {
		if s.Hour.EveryValue() && s.DOM.EveryValue() && s.Month.EveryValue() && s.DOW.EveryValue() {
			return "every minute starting at 00:00", Unsupported{}, true
		}
		return "", Unsupported{Reason: "an every-minute rule restricted by hour, date, or weekday has no phrase equivalent"}, false
	}

	minute, ok := s.Minute.Single()
	if !ok {
		return "", Unsupported{Reason: "a minute list has no phrase equivalent; only a single minute or an evenly dividing step is expressible"}, false
	}

	// Hourly at a fixed minute: "0 * * * *".
	if s.Hour.EveryValue() {
		if !s.DOM.EveryValue() || !s.Month.EveryValue() || !s.DOW.EveryValue() {
			return "", Unsupported{Reason: "an hourly rule restricted by date or weekday has no phrase equivalent"}, false
		}
		if minute != 0 {
			return "", Unsupported{Reason: "an hourly rule at a minute other than :00 has no phrase equivalent"}, false
		}
		return "every hour starting at 00:00", Unsupported{}, true
	}

	// A step in the hour field: "0 */6 * * *".
	if s.Hour.Wildcard && s.Hour.Step > 1 {
		if !s.DOM.EveryValue() || !s.Month.EveryValue() || !s.DOW.EveryValue() || minute != 0 {
			return "", Unsupported{Reason: "an hour step combined with other restrictions has no phrase equivalent"}, false
		}
		if 24%s.Hour.Step != 0 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"a step of %d does not divide the day evenly: cron restarts the sequence at midnight, which a fixed interval does not reproduce",
				s.Hour.Step)}, false
		}
		return fmt.Sprintf("every %d hours starting at 00:00", s.Hour.Step), Unsupported{}, true
	}

	hour, ok := s.Hour.Single()
	if !ok {
		return "", Unsupported{Reason: "an hour list has no phrase equivalent; only a single hour or an evenly dividing step is expressible"}, false
	}
	at := fmt.Sprintf(" at %02d:%02d", hour, minute)
	if s.DOM.calendarSelector != calendarNone {
		if !s.Month.EveryValue() || !s.DOW.EveryValue() {
			return "", Unsupported{Reason: "a monthly day-of-month selector requires unrestricted month and day-of-week fields"}, false
		}
		switch s.DOM.calendarSelector {
		case calendarLastDay:
			return "last day of every month" + at, Unsupported{}, true
		case calendarLastWeekday:
			return "last weekday of every month" + at, Unsupported{}, true
		case calendarNearestWeekday:
			day, ok := s.DOM.Single()
			if !ok {
				return "", Unsupported{Reason: "only one nearest-weekday date has a phrase equivalent"}, false
			}
			return fmt.Sprintf("nearest weekday to the %s of every month%s", ordinal(day), at), Unsupported{}, true
		}
		return "", Unsupported{Reason: "that monthly day-of-month selector has no phrase equivalent"}, false
	}
	if s.DOW.Ordinal != 0 {
		if !s.DOM.EveryValue() {
			return "", Unsupported{Reason: "a monthly weekday selector combined with a day-of-month restriction has no phrase equivalent"}, false
		}
		if !s.Month.EveryValue() {
			return "", Unsupported{Reason: "a monthly weekday selector restricted to particular months has no phrase equivalent"}, false
		}
		day, ok := s.DOW.Single()
		if !ok {
			return "", Unsupported{Reason: "only one monthly weekday selector has a phrase equivalent"}, false
		}
		if s.DOW.Ordinal == -1 {
			return fmt.Sprintf("last %s of the month%s", dayName(day), at), Unsupported{}, true
		}
		return fmt.Sprintf("%s %s monthly%s", ordinal(s.DOW.Ordinal), dayName(day), at), Unsupported{}, true
	}

	// From here the time of day is fixed. What remains is which days.
	switch {
	case s.DOM.EveryValue() && s.Month.EveryValue() && s.DOW.EveryValue():
		return "every day" + at, Unsupported{}, true

	case s.DOM.EveryValue() && s.Month.EveryValue(): // a weekday restriction
		phrase, ok := weekdayPhrase(s.DOW)
		if !ok {
			return "", Unsupported{Reason: "that combination of weekdays has no phrase equivalent"}, false
		}
		return phrase + at, Unsupported{}, true

	case s.DOW.EveryValue() && s.Month.EveryValue(): // a day-of-month restriction
		day, ok := s.DOM.Single()
		if !ok {
			return "", Unsupported{Reason: "a day-of-month list has no phrase equivalent; only a single date is expressible"}, false
		}
		return fmt.Sprintf("on the %s of every month%s", ordinal(day), at), Unsupported{}, true

	case s.DOW.EveryValue(): // a specific month and date - a yearly rule
		day, dayOK := s.DOM.Single()
		month, monthOK := s.Month.Single()
		if !dayOK || !monthOK {
			return "", Unsupported{Reason: "a month or date list has no phrase equivalent; only a single month and date are expressible"}, false
		}
		return fmt.Sprintf("every year on %s %d%s", monthName(month), day, at), Unsupported{}, true
	}

	return "", Unsupported{Reason: "that combination of fields has no phrase equivalent"}, false
}

// weekdayPhrase renders a day-of-week set, covering the three shapes the
// grammar has words for: the weekday set, the weekend set, and a single day.
func weekdayPhrase(f Field) (string, bool) {
	set := map[int]bool{}
	for _, v := range f.Values {
		set[v] = true
	}
	switch {
	case len(set) == 5 && set[1] && set[2] && set[3] && set[4] && set[5]:
		return "weekdays", true
	case len(set) == 2 && set[0] && set[6]:
		return "weekends", true
	case len(set) == 1:
		for v := range set {
			return "every " + dayName(v), true
		}
	}
	return "", false
}

func dayName(v int) string {
	return [...]string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}[v%7]
}

func monthName(v int) string {
	return [...]string{"", "january", "february", "march", "april", "may", "june",
		"july", "august", "september", "october", "november", "december"}[v]
}

// ordinal renders an English ordinal for a day of the month.
func ordinal(n int) string {
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
	return fmt.Sprintf("%d%s", n, suffix)
}

// Explain is the whole conversion in one call: an expression in, a phrase out,
// or a named refusal. It is what the explain and import paths both use, so the
// two cannot diverge.
func Explain(expr string) (phrase string, bad Unsupported, err error) {
	res, err := Parse(expr)
	if err != nil {
		return "", Unsupported{}, err
	}
	if !res.OK {
		return "", res.Bad, nil
	}
	phrase, bad, ok := Phrase(res.Spec)
	if !ok {
		if res.Spec.DOM.calendarSelector == calendarNone && res.Spec.DOW.Ordinal == 0 {
			return describeSpec(res.Spec), Unsupported{}, nil
		}
		bad.Input = strings.TrimSpace(expr)
		return "", bad, nil
	}
	return phrase, Unsupported{}, nil
}
