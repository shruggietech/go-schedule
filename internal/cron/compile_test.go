package cron

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestCompileCompositeCronToDurableRecurrence(t *testing.T) {
	now := time.Date(2026, 8, 28, 12, 34, 56, 0, time.UTC)
	sch, bad, err := Compile("30 8-17 * * MON,WED,FRI", "America/New_York", now)
	if err != nil || bad.Reason != "" {
		t.Fatalf("Compile: err=%v refusal=%q", err, bad.Reason)
	}
	if sch.Kind != domain.ScheduleRecurring || sch.Anchor == nil || !sch.Anchor.Equal(now) {
		t.Fatalf("schedule identity=%+v", sch)
	}
	if sch.Expression != "30 8-17 * * MON,WED,FRI" || sch.HumanSummary == "" {
		t.Fatalf("source/summary=%q / %q", sch.Expression, sch.HumanSummary)
	}
	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		t.Fatal(err)
	}
	if opt.Freq != rrule.DAILY || !reflect.DeepEqual(opt.Byminute, []int{30}) ||
		!reflect.DeepEqual(opt.Byhour, rangeOf(8, 17)) || !reflect.DeepEqual(opt.Bysecond, []int{0}) {
		t.Fatalf("compiled time fields=%+v", opt)
	}
	wantDays := []rrule.Weekday{rrule.MO, rrule.WE, rrule.FR}
	if !reflect.DeepEqual(opt.Byweekday, wantDays) {
		t.Fatalf("compiled weekdays=%v, want %v", opt.Byweekday, wantDays)
	}
}

func TestCompilePreservesFocusedCalendarSelectors(t *testing.T) {
	for _, expr := range []string{"0 9 * * 5#3", "0 9 * * 5L", "0 9 L * *", "0 9 15W * *", "0 9 LW * *"} {
		sch, bad, err := Compile(expr, "UTC", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil || bad.Reason != "" || sch.Kind != domain.ScheduleRecurring {
			t.Errorf("Compile(%q): schedule=%+v refusal=%q err=%v", expr, sch, bad.Reason, err)
		}
	}
}

func TestCompileReturnsRefusalWithoutSchedule(t *testing.T) {
	sch, bad, err := Compile("0 0 13 * 5", "UTC", time.Now())
	if err != nil || bad.Reason == "" || sch.Kind != "" {
		t.Fatalf("schedule=%+v refusal=%q err=%v", sch, bad.Reason, err)
	}
}

func TestCompileRetainsCalendarWildcardStepSets(t *testing.T) {
	for expression, token := range map[string]string{
		"0 9 */2 * *": "BYMONTHDAY=1,3,5,7,9,11,13,15,17,19,21,23,25,27,29,31",
		"0 9 * */2 *": "BYMONTH=1,3,5,7,9,11",
		"0 9 * * */2": "BYDAY=SU,TU,TH,SA",
	} {
		sch, bad, err := Compile(expression, "UTC", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil || bad.Reason != "" || !strings.Contains(sch.RRULE, token) {
			t.Fatalf("Compile(%q) RRULE=%q refusal=%q err=%v, want %q", expression, sch.RRULE, bad.Reason, err, token)
		}
	}
}
