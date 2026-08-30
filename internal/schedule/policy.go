package schedule

import (
	"fmt"
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// DescribePolicy renders the effective recurrence clock and transition choices
// for human-facing task detail and previews.
func DescribePolicy(policy domain.SchedulePolicy) string {
	policy = policy.Effective()
	switch policy.TimeBasis {
	case domain.TimeBasisElapsed:
		return "Fixed elapsed interval (DST transition choices stored but inactive)"
	case domain.TimeBasisUTC:
		return "UTC clock (DST transition choices stored but inactive)"
	default:
		gap := strings.ReplaceAll(string(policy.DSTGap), "_", " ")
		return fmt.Sprintf("Local wall clock; spring gap: %s; fall overlap: %s", gap, policy.DSTOverlap)
	}
}
