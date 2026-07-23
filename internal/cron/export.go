package cron

import (
	"fmt"
	"strings"

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
	switch {
	case !task.Enabled || task.State == domain.TaskDisabled:
		return "", Unsupported{Reason: "the task is disabled and cron has no disabled state"}, false
	case sch.Kind == domain.ScheduleOneOff:
		return "", Unsupported{Reason: "cron cannot express a schedule that fires exactly once"}, false
	case sch.Kind != domain.ScheduleRecurring:
		return "", Unsupported{Reason: "only a recurring schedule can be expressed as cron"}, false
	}

	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		return "", Unsupported{Reason: "the stored recurrence could not be read"}, false
	}
	interval := opt.Interval
	if interval < 1 {
		interval = 1
	}

	// A non-default missing-date policy changes which dates the task runs on,
	// and cron has no notion of it. Saying so is the whole point of the export.
	if task.MissingDatePolicy != "" && task.MissingDatePolicy != domain.MissingDateSkip {
		if datebearing(opt) {
			return "", Unsupported{Reason: fmt.Sprintf(
				"the task's missing-date policy (%s) has no cron equivalent — cron would silently skip the periods this task runs in",
				task.MissingDatePolicy)}, false
		}
	}

	minute, hour := timeFields(opt)

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
		if len(opt.Byweekday) > 0 {
			return "", Unsupported{Reason: "an ordinal-weekday rule (\"the 3rd Wednesday\") has no standard cron equivalent"}, false
		}
		if len(opt.Bymonthday) != 1 || opt.Bymonthday[0] < 1 {
			return "", Unsupported{Reason: "only a single, positive day of the month can be expressed as cron"}, false
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

// datebearing reports whether the rule addresses a date that can be absent from
// a period, which is when the missing-date policy actually changes run times.
func datebearing(opt *rrule.ROption) bool {
	if len(opt.Bymonthday) == 1 && opt.Bymonthday[0] > 28 {
		return true
	}
	return len(opt.Byweekday) == 1 && opt.Byweekday[0].N() >= 5
}

// timeFields recovers the rule's minute and hour, falling back to the anchor.
func timeFields(opt *rrule.ROption) (minute, hour int) {
	minute, hour = opt.Dtstart.Minute(), opt.Dtstart.Hour()
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
	nums := make([]string, 0, len(opt.Byweekday))
	for _, w := range opt.Byweekday {
		if w.N() != 0 {
			return "", false // an ordinal weekday is not a weekly rule
		}
		// rrule numbers days from Monday=0; cron from Sunday=0.
		nums = append(nums, fmt.Sprint((w.Day()+1)%7))
	}
	return strings.Join(nums, ","), true
}
