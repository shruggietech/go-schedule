package scheduleinput

import (
	"strings"
	"testing"
	"time"

	"github.com/teambition/rrule-go"
)

func TestParseCompositeCronCompilesWithoutHumanPhrase(t *testing.T) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	got, err := Parse("0 9,17 * * MON,WED,FRI", SyntaxCron, "UTC", now)
	if err != nil {
		t.Fatal(err)
	}
	if got.Syntax != SyntaxCron || got.Expression != "0 9,17 * * MON,WED,FRI" {
		t.Fatalf("identity=%+v", got)
	}
	if got.Schedule.Expression != got.Expression || !strings.Contains(strings.ToLower(got.Schedule.HumanSummary), "hours") {
		t.Fatalf("schedule source/summary=%q / %q", got.Schedule.Expression, got.Schedule.HumanSummary)
	}
	opt, err := rrule.StrToROption(got.Schedule.RRULE)
	if err != nil {
		t.Fatal(err)
	}
	if opt.Freq != rrule.DAILY || len(opt.Byhour) != 2 || len(opt.Byweekday) != 3 {
		t.Fatalf("compiled recurrence=%+v", opt)
	}
}

func TestParseCompositeCronRefusalDoesNotFallback(t *testing.T) {
	_, err := Parse("0 0 13 * 5", "", "UTC", time.Now())
	if err == nil || !strings.Contains(err.Error(), "either") {
		t.Fatalf("error=%v, want DOM/DOW refusal", err)
	}
}
