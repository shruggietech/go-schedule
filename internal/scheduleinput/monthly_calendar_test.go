package scheduleinput

import (
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func TestMonthlyCalendarCronAndHumanParity(t *testing.T) {
	tests := []struct {
		cron, human string
		adjustment  domain.CalendarAdjustment
	}{
		{"0 9 L * *", "last day of every month at 09:00", ""},
		{"0 9 15W * *", "nearest weekday to the 15th of every month at 09:00", domain.CalendarAdjustmentNearestWeekday},
		{"0 9 LW * *", "last weekday of every month at 09:00", ""},
	}
	for _, tt := range tests {
		cronInput, err := Parse(tt.cron, SyntaxCron, "America/New_York", inputAnchor)
		if err != nil {
			t.Fatal(err)
		}
		humanInput, err := Parse(tt.human, SyntaxHuman, "America/New_York", inputAnchor)
		if err != nil {
			t.Fatal(err)
		}
		if cronInput.Schedule.RRULE != humanInput.Schedule.RRULE || cronInput.Schedule.CalendarAdjustment != tt.adjustment {
			t.Fatalf("cron=%+v human=%+v", cronInput.Schedule, humanInput.Schedule)
		}
		if cronInput.Schedule.Expression != tt.cron || cronInput.Syntax != SyntaxCron {
			t.Fatalf("cron identity lost: %+v", cronInput)
		}
		cronRuns, err := schedule.UpcomingRuns(cronInput.Schedule, "America/New_York", domain.MissingDateSkip, inputAnchor, 12)
		if err != nil {
			t.Fatal(err)
		}
		humanRuns, err := schedule.UpcomingRuns(humanInput.Schedule, "America/New_York", domain.MissingDateSkip, inputAnchor, 12)
		if err != nil || len(cronRuns) != len(humanRuns) {
			t.Fatalf("run counts: cron=%d human=%d err=%v", len(cronRuns), len(humanRuns), err)
		}
		for i := range cronRuns {
			if !cronRuns[i].Equal(humanRuns[i]) {
				t.Fatalf("run %d: cron=%v human=%v", i, cronRuns[i], humanRuns[i])
			}
		}
	}
}
