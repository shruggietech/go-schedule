package cron

import (
	"fmt"
	"strings"
)

// Phrase renders a Spec as the human-readable phrase a user would have typed,
// or refuses by name. The phrase is the only route from cron into a schedule:
// callers pass it to schedule.Parse exactly as they would a phrase the operator
// typed themselves. That is what makes an import preview trustworthy — the
// string shown is the string parsed — and what keeps cron from becoming a second
// authoring path into the engine.
func Phrase(s Spec) (string, Unsupported, bool) {
	// Sub-hourly: "*/n" in the minute field with every hour.
	if s.Minute.Wildcard && s.Minute.Step > 1 {
		if !s.Hour.EveryValue() || !s.DOM.Wildcard || !s.Month.Wildcard || !s.DOW.Wildcard {
			return "", Unsupported{Reason: "a minute step combined with other restrictions has no phrase equivalent"}, false
		}
		if 60%s.Minute.Step != 0 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"a step of %d does not divide the hour evenly: cron restarts the sequence at :00, which a fixed interval does not reproduce",
				s.Minute.Step)}, false
		}
		return fmt.Sprintf("every %d minutes", s.Minute.Step), Unsupported{}, true
	}

	// Every minute.
	if s.Minute.EveryValue() {
		if s.Hour.EveryValue() && s.DOM.Wildcard && s.Month.Wildcard && s.DOW.Wildcard {
			return "every minute", Unsupported{}, true
		}
		return "", Unsupported{Reason: "an every-minute rule restricted by hour, date, or weekday has no phrase equivalent"}, false
	}

	minute, ok := s.Minute.Single()
	if !ok {
		return "", Unsupported{Reason: "a minute list has no phrase equivalent; only a single minute or an evenly dividing step is expressible"}, false
	}

	// Hourly at a fixed minute: "0 * * * *".
	if s.Hour.EveryValue() {
		if !s.DOM.Wildcard || !s.Month.Wildcard || !s.DOW.Wildcard {
			return "", Unsupported{Reason: "an hourly rule restricted by date or weekday has no phrase equivalent"}, false
		}
		if minute != 0 {
			return "", Unsupported{Reason: "an hourly rule at a minute other than :00 has no phrase equivalent"}, false
		}
		return "every hour", Unsupported{}, true
	}

	// A step in the hour field: "0 */6 * * *".
	if s.Hour.Wildcard && s.Hour.Step > 1 {
		if !s.DOM.Wildcard || !s.Month.Wildcard || !s.DOW.Wildcard || minute != 0 {
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

	// From here the time of day is fixed. What remains is which days.
	switch {
	case s.DOM.Wildcard && s.Month.Wildcard && s.DOW.Wildcard:
		return "every day" + at, Unsupported{}, true

	case s.DOM.Wildcard && s.Month.Wildcard: // a weekday restriction
		phrase, ok := weekdayPhrase(s.DOW)
		if !ok {
			return "", Unsupported{Reason: "that combination of weekdays has no phrase equivalent"}, false
		}
		return phrase + at, Unsupported{}, true

	case s.DOW.Wildcard && s.Month.Wildcard: // a day-of-month restriction
		day, ok := s.DOM.Single()
		if !ok {
			return "", Unsupported{Reason: "a day-of-month list has no phrase equivalent; only a single date is expressible"}, false
		}
		return fmt.Sprintf("on the %s of every month%s", ordinal(day), at), Unsupported{}, true

	case s.DOW.Wildcard: // a specific month and date — a yearly rule
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
		bad.Input = strings.TrimSpace(expr)
		return "", bad, nil
	}
	return phrase, Unsupported{}, nil
}
