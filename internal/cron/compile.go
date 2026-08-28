package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

// Compile parses one supported cron expression into the same durable schedule
// model used by human input. A recognizable but unsupported expression returns
// a refusal; malformed syntax or an invalid timezone returns an error.
func Compile(input, tzName string, now time.Time) (domain.Schedule, Unsupported, error) {
	expression := strings.TrimSpace(input)
	res, err := Parse(expression)
	if err != nil {
		return domain.Schedule{}, Unsupported{}, err
	}
	if !res.OK {
		return domain.Schedule{}, res.Bad, nil
	}

	// Focused calendar modifiers already have purpose-built human grammar and
	// schedule mappings. Keep those mappings rather than flattening their extra
	// semantics into ordinary field sets.
	if res.Spec.DOM.calendarSelector != calendarNone || res.Spec.DOW.Ordinal != 0 {
		phrase, bad, ok := Phrase(res.Spec)
		if !ok {
			bad.Input = expression
			return domain.Schedule{}, bad, nil
		}
		sch, err := schedule.Parse(phrase, tzName, now)
		if err != nil {
			return domain.Schedule{}, Unsupported{}, fmt.Errorf("cron: compile focused selector: %w", err)
		}
		sch.Expression = expression
		return sch, Unsupported{}, nil
	}

	loc, err := timezone.Resolve(tzName)
	if err != nil {
		return domain.Schedule{}, Unsupported{}, err
	}
	opt := rrule.ROption{
		Freq:     rrule.DAILY,
		Interval: 1,
		Dtstart:  now.In(loc),
		Byhour:   append([]int(nil), res.Spec.Hour.Values...),
		Byminute: append([]int(nil), res.Spec.Minute.Values...),
		Bysecond: []int{0},
	}
	if res.Spec.DOM.Restricted() {
		opt.Bymonthday = append([]int(nil), res.Spec.DOM.Values...)
	}
	if res.Spec.Month.Restricted() {
		opt.Bymonth = append([]int(nil), res.Spec.Month.Values...)
	}
	if res.Spec.DOW.Restricted() {
		opt.Byweekday = cronWeekdays(res.Spec.DOW.Values)
	}
	if _, err := rrule.NewRRule(opt); err != nil {
		return domain.Schedule{}, Unsupported{}, fmt.Errorf("cron: build recurrence: %w", err)
	}

	description := describeSpec(res.Spec)
	if phrase, _, ok := Phrase(res.Spec); ok {
		description = phrase
	}
	anchor := now.UTC()
	return domain.Schedule{
		Kind:         domain.ScheduleRecurring,
		RRULE:        opt.RRuleString(),
		Anchor:       &anchor,
		HumanSummary: sentenceCase(description),
		Expression:   expression,
	}, Unsupported{}, nil
}

func cronWeekdays(values []int) []rrule.Weekday {
	seen := [7]bool{}
	for _, value := range values {
		seen[value%7] = true
	}
	out := make([]rrule.Weekday, 0, len(values))
	for cronDay, present := range seen {
		if !present {
			continue
		}
		rruleDay := (cronDay + 6) % 7 // cron Sunday=0; rrule Monday=0.
		out = append(out, [...]rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR, rrule.SA, rrule.SU}[rruleDay])
	}
	return out
}

func sentenceCase(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func describeSpec(s Spec) string {
	parts := []string{describeTime(s.Minute, s.Hour)}
	switch {
	case s.DOM.Restricted():
		parts = append(parts, "on "+pluralValues("date", s.DOM.Values))
	case s.DOW.Restricted():
		parts = append(parts, "on "+humanWeekdays(s.DOW.Values))
	default:
		parts = append(parts, "every day")
	}
	if s.Month.Restricted() {
		parts = append(parts, "in "+humanMonths(s.Month.Values))
	}
	return strings.Join(parts, " ")
}

func describeTime(minute, hour Field) string {
	if minute.Wildcard && minute.Step > 1 {
		if 60%minute.Step != 0 {
			result := "at " + pluralValues("minute", minute.Values)
			if hour.EveryValue() {
				return result + " of each hour"
			}
			return result + " during " + pluralValues("hour", hour.Values)
		}
		result := fmt.Sprintf("every %d minutes", minute.Step)
		if !hour.EveryValue() {
			result += " during " + pluralValues("hour", hour.Values)
		}
		return result
	}
	if minute.EveryValue() {
		if hour.EveryValue() {
			return "every minute"
		}
		return "every minute during " + pluralValues("hour", hour.Values)
	}
	if len(minute.Values) == 1 && len(hour.Values) == 1 {
		return fmt.Sprintf("at %02d:%02d", hour.Values[0], minute.Values[0])
	}
	if len(minute.Values) == 1 {
		return fmt.Sprintf("at minute %02d during %s", minute.Values[0], pluralValues("hour", hour.Values))
	}
	result := "at " + pluralValues("minute", minute.Values)
	if hour.EveryValue() {
		return result + " every hour"
	}
	return result + " during " + pluralValues("hour", hour.Values)
}

func pluralValues(noun string, values []int) string {
	if len(values) == 1 {
		return noun + " " + fmt.Sprint(values[0])
	}
	return noun + "s " + humanNumbers(values)
}

func humanNumbers(values []int) string {
	if len(values) >= 3 && consecutive(values) {
		return fmt.Sprintf("%d through %d", values[0], values[len(values)-1])
	}
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprint(value)
	}
	return humanList(parts)
}

func consecutive(values []int) bool {
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1]+1 {
			return false
		}
	}
	return true
}

func humanWeekdays(values []int) string {
	parts := make([]string, 0, len(values))
	seen := [7]bool{}
	for _, value := range values {
		seen[value%7] = true
	}
	for day, present := range seen {
		if present {
			parts = append(parts, sentenceCase(dayName(day)))
		}
	}
	return humanList(parts)
}

func humanMonths(values []int) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = sentenceCase(monthName(value))
	}
	return humanList(parts)
}

func humanList(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
	}
}
