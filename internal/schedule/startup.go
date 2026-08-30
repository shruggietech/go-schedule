package schedule

import (
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	StartupPhrase  = "at scheduler startup"
	StartupSummary = "At scheduler startup"
)

// NewStartup constructs the one supported event schedule. The expression is
// retained for editing and never participates in execution.
func NewStartup(expression string) domain.Schedule {
	return domain.Schedule{
		Kind:         domain.ScheduleEvent,
		TriggerID:    domain.StartupEventID,
		HumanSummary: StartupSummary,
		Expression:   strings.TrimSpace(expression),
	}
}

// IsStartup reports whether a schedule represents the daemon-start event.
func IsStartup(sch domain.Schedule) bool {
	return sch.Kind == domain.ScheduleEvent && sch.TriggerID == domain.StartupEventID
}
