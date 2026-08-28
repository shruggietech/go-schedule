package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Export renders a task's schedule as a crontab timing expression, or refuses by
// name. It is a pure function of the stored schedule: nothing about the export
// depends on the daemon's state, which is what makes it usable as a diff between
// two machines.
//
// It declines rather than approximates. A schedule cron cannot carry is worth
// more as a visible refusal than as a line that runs at the wrong time.
func Export(task domain.Task, sch domain.Schedule) (expr string, bad Unsupported, ok bool) {
	if !task.Enabled || task.State == domain.TaskDisabled {
		return "", Unsupported{Reason: "the task is disabled and cron has no disabled state"}, false
	}
	return ExportSchedule(sch, task.MissingDatePolicy)
}

// ExportSchedule renders a schedule as a cron timing expression without task
// state or command policy. It is the pure recurrence mapping shared by task
// export and string conversion.
func ExportSchedule(sch domain.Schedule, policy domain.MissingDatePolicy) (expr string, bad Unsupported, ok bool) {
	switch {
	case sch.Kind == domain.ScheduleOneOff:
		return "", Unsupported{Reason: "cron cannot express a schedule that fires exactly once"}, false
	case sch.Kind != domain.ScheduleRecurring:
		return "", Unsupported{Reason: "only a recurring schedule can be expressed as cron"}, false
	}

	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		return "", Unsupported{Reason: "the stored recurrence could not be read"}, false
	}
	if sch.CalendarAdjustment != "" && sch.CalendarAdjustment != domain.CalendarAdjustmentNearestWeekday {
		return "", Unsupported{Reason: fmt.Sprintf("the stored calendar adjustment %q has no cron equivalent", sch.CalendarAdjustment)}, false
	}
	if sch.CalendarAdjustment == domain.CalendarAdjustmentNearestWeekday {
		if opt.Freq != rrule.MONTHLY {
			return "", Unsupported{Reason: "nearest_weekday requires one unbounded monthly day-of-month rule"}, false
		}
		if _, ok := nearestWeekdayField(opt); !ok {
			return "", Unsupported{Reason: "nearest_weekday requires one unbounded monthly day-of-month rule"}, false
		}
	}
	interval := opt.Interval
	if interval < 1 {
		interval = 1
	}
	if reason := subDailyPhaseRefusal(sch, opt, interval); reason != "" {
		return "", Unsupported{Reason: reason}, false
	}
	if reason := calendarAnchorRefusal(sch, opt); reason != "" {
		return "", Unsupported{Reason: reason}, false
	}

	// A non-default missing-date policy changes which dates the task runs on,
	// and cron has no notion of it. Saying so is the whole point of the export.
	if policy != "" && policy != domain.MissingDateSkip {
		if datebearing(opt) {
			return "", Unsupported{Reason: fmt.Sprintf(
				"the task's missing-date policy (%s) has no cron equivalent — cron would silently skip the periods this task runs in",
				policy)}, false
		}
	}

	minute, hour := timeFields(opt, sch.Anchor)

	switch opt.Freq {
	case rrule.SECONDLY:
		return "", Unsupported{Reason: "cron has no sub-minute resolution"}, false

	case rrule.MINUTELY:
		if 60%interval != 0 || interval > 30 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"an interval of %d minutes does not divide the hour evenly, so cron cannot reproduce it", interval)}, false
		}
		if interval == 1 {
			return "* * * * *", Unsupported{}, true
		}
		return fmt.Sprintf("*/%d * * * *", interval), Unsupported{}, true

	case rrule.HOURLY:
		if 24%interval != 0 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"an interval of %d hours does not divide the day evenly, so cron cannot reproduce it", interval)}, false
		}
		if interval == 1 {
			return fmt.Sprintf("%d * * * *", minute), Unsupported{}, true
		}
		return fmt.Sprintf("%d */%d * * *", minute, interval), Unsupported{}, true

	case rrule.DAILY:
		if interval != 1 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"an every-%d-days rule has no cron equivalent — cron repeats by calendar position, not by elapsed days", interval)}, false
		}
		return fmt.Sprintf("%d %d * * *", minute, hour), Unsupported{}, true

	case rrule.WEEKLY:
		if interval != 1 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"an every-%d-weeks rule has no cron equivalent", interval)}, false
		}
		days, ok := weekdayField(opt)
		if !ok {
			return "", Unsupported{Reason: "that weekday selection has no cron equivalent"}, false
		}
		return fmt.Sprintf("%d %d * * %s", minute, hour, days), Unsupported{}, true

	case rrule.MONTHLY:
		if interval != 1 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"an every-%d-months rule has no cron equivalent", interval)}, false
		}
		if sch.CalendarAdjustment == domain.CalendarAdjustmentNearestWeekday {
			day, ok := nearestWeekdayField(opt)
			if !ok {
				return "", Unsupported{Reason: "nearest_weekday requires one unbounded monthly day-of-month rule"}, false
			}
			return fmt.Sprintf("%d %d %dW * *", minute, hour, day), Unsupported{}, true
		}
		if lastWeekdayOfMonth(opt) {
			return fmt.Sprintf("%d %d LW * *", minute, hour), Unsupported{}, true
		}
		if len(opt.Byweekday) > 0 {
			weekday, occurrence, ok := monthlyWeekdayField(opt)
			if !ok {
				return "", Unsupported{Reason: "only one first-through-fifth or last weekday has a cron equivalent"}, false
			}
			if occurrence == -1 {
				return fmt.Sprintf("%d %d * * %dL", minute, hour, weekday), Unsupported{}, true
			}
			return fmt.Sprintf("%d %d * * %d#%d", minute, hour, weekday, occurrence), Unsupported{}, true
		}
		if len(opt.Bymonthday) == 1 && opt.Bymonthday[0] == -1 && plainMonthlySelector(opt) {
			return fmt.Sprintf("%d %d L * *", minute, hour), Unsupported{}, true
		}
		if len(opt.Bymonthday) != 1 || opt.Bymonthday[0] < 1 || !plainMonthlySelector(opt) {
			return "", Unsupported{Reason: "only one supported day-of-month selector can be expressed as cron"}, false
		}
		return fmt.Sprintf("%d %d %d * *", minute, hour, opt.Bymonthday[0]), Unsupported{}, true

	case rrule.YEARLY:
		if interval != 1 {
			return "", Unsupported{Reason: fmt.Sprintf(
				"an every-%d-years rule has no cron equivalent", interval)}, false
		}
		if len(opt.Bymonth) != 1 || len(opt.Bymonthday) != 1 {
			return "", Unsupported{Reason: "only a single month and date can be expressed as cron"}, false
		}
		return fmt.Sprintf("%d %d %d %d *", minute, hour, opt.Bymonthday[0], opt.Bymonth[0]), Unsupported{}, true
	}

	return "", Unsupported{Reason: "that recurrence has no cron equivalent"}, false
}

func plainMonthlySelector(opt *rrule.ROption) bool {
	return len(opt.Byweekday) == 0 && len(opt.Bymonth) == 0 && len(opt.Bysetpos) == 0 &&
		len(opt.Byyearday) == 0 && len(opt.Byweekno) == 0 && len(opt.Byeaster) == 0 &&
		opt.Count == 0 && opt.Until.IsZero() && len(opt.Byhour) <= 1 && len(opt.Byminute) <= 1 &&
		(len(opt.Bysecond) == 0 || len(opt.Bysecond) == 1 && opt.Bysecond[0] == 0)
}

func nearestWeekdayField(opt *rrule.ROption) (int, bool) {
	if !plainMonthlySelector(opt) || len(opt.Bymonthday) != 1 || opt.Bymonthday[0] < 1 || opt.Bymonthday[0] > 31 {
		return 0, false
	}
	return opt.Bymonthday[0], true
}

func lastWeekdayOfMonth(opt *rrule.ROption) bool {
	if len(opt.Byweekday) != 5 || len(opt.Bymonthday) != 0 || len(opt.Bymonth) != 0 ||
		len(opt.Bysetpos) != 1 || opt.Bysetpos[0] != -1 || len(opt.Byyearday) != 0 ||
		len(opt.Byweekno) != 0 || len(opt.Byeaster) != 0 || opt.Count != 0 || !opt.Until.IsZero() ||
		len(opt.Byhour) > 1 || len(opt.Byminute) > 1 || len(opt.Bysecond) > 1 ||
		len(opt.Bysecond) == 1 && opt.Bysecond[0] != 0 {
		return false
	}
	days := map[int]bool{}
	for _, day := range opt.Byweekday {
		if day.N() != 0 {
			return false
		}
		days[day.Day()] = true
	}
	return days[rrule.MO.Day()] && days[rrule.TU.Day()] && days[rrule.WE.Day()] && days[rrule.TH.Day()] && days[rrule.FR.Day()]
}

func monthlyWeekdayField(opt *rrule.ROption) (weekday, occurrence int, ok bool) {
	if len(opt.Byweekday) != 1 || len(opt.Bymonthday) != 0 || len(opt.Bymonth) != 0 ||
		len(opt.Bysetpos) != 0 || len(opt.Byyearday) != 0 || len(opt.Byweekno) != 0 || len(opt.Byeaster) != 0 ||
		opt.Count != 0 || !opt.Until.IsZero() || len(opt.Byhour) > 1 || len(opt.Byminute) > 1 ||
		len(opt.Bysecond) > 1 || len(opt.Bysecond) == 1 && opt.Bysecond[0] != 0 {
		return 0, 0, false
	}
	w := opt.Byweekday[0]
	occurrence = w.N()
	if occurrence != -1 && (occurrence < 1 || occurrence > 5) {
		return 0, 0, false
	}
	return (w.Day() + 1) % 7, occurrence, true
}

func calendarAnchorRefusal(sch domain.Schedule, opt *rrule.ROption) string {
	switch opt.Freq {
	case rrule.DAILY, rrule.WEEKLY, rrule.MONTHLY, rrule.YEARLY:
	default:
		return ""
	}
	if len(opt.Byhour) == 1 && len(opt.Byminute) == 1 {
		return ""
	}
	if sch.Anchor == nil {
		return "the recurrence has no anchor from which to recover its time of day"
	}
	if sch.Anchor.Second() != 0 || sch.Anchor.Nanosecond() != 0 {
		return "the schedule phase is below cron's one-minute resolution"
	}
	return ""
}

func subDailyPhaseRefusal(sch domain.Schedule, opt *rrule.ROption, interval int) string {
	if opt.Freq != rrule.MINUTELY && opt.Freq != rrule.HOURLY {
		return ""
	}
	if sch.Anchor == nil || sch.Anchor.Second() != 0 || sch.Anchor.Nanosecond() != 0 {
		return "the interval phase is below cron's one-minute resolution"
	}
	if opt.Freq == rrule.MINUTELY && sch.Anchor.Minute()%interval != 0 {
		return "the interval phase does not align with cron's minute step"
	}
	if opt.Freq == rrule.HOURLY && sch.Anchor.Hour()%interval != 0 {
		return "the interval phase does not align with cron's hour step"
	}
	return ""
}

// datebearing reports whether the rule addresses a date that can be absent from
// a period, which is when the missing-date policy actually changes run times.
func datebearing(opt *rrule.ROption) bool {
	if len(opt.Bymonthday) == 1 && opt.Bymonthday[0] > 28 {
		return true
	}
	return len(opt.Byweekday) == 1 && opt.Byweekday[0].N() >= 5
}

// timeFields recovers the rule's minute and hour, falling back to the anchor.
func timeFields(opt *rrule.ROption, anchor *time.Time) (minute, hour int) {
	minute, hour = opt.Dtstart.Minute(), opt.Dtstart.Hour()
	if anchor != nil {
		minute, hour = anchor.Minute(), anchor.Hour()
	}
	if len(opt.Byminute) == 1 {
		minute = opt.Byminute[0]
	}
	if len(opt.Byhour) == 1 {
		hour = opt.Byhour[0]
	}
	return minute, hour
}

// weekdayField renders a weekly rule's days as a cron day-of-week field.
func weekdayField(opt *rrule.ROption) (string, bool) {
	if len(opt.Byweekday) == 0 {
		return "", false
	}
	var days [7]bool
	count := 0
	for _, w := range opt.Byweekday {
		if w.N() != 0 {
			return "", false // an ordinal weekday is not a weekly rule
		}
		// rrule numbers days from Monday=0; cron from Sunday=0.
		day := (w.Day() + 1) % 7
		if !days[day] {
			days[day] = true
			count++
		}
	}
	if count == 5 && days[1] && days[2] && days[3] && days[4] && days[5] {
		return "1-5", true
	}
	nums := make([]string, 0, count)
	for day, present := range days {
		if present {
			nums = append(nums, fmt.Sprint(day))
		}
	}
	return strings.Join(nums, ","), true
}
