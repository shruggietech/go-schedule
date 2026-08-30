package schedule

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestParseSchedulerStartupHasNoClockOccurrence(t *testing.T) {
	sch, err := Parse(" At Scheduler Startup ", "America/New_York", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !IsStartup(sch) || sch.HumanSummary != StartupSummary || sch.Expression != "At Scheduler Startup" {
		t.Fatalf("startup schedule = %+v", sch)
	}
	runs, err := UpcomingRunsWithPolicy(sch, "America/New_York", domain.SchedulePolicy{}.Effective(), time.Now(), 5)
	if err != nil || len(runs) != 0 {
		t.Fatalf("upcoming startup runs = %v, err=%v", runs, err)
	}
}
