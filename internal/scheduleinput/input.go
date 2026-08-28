// Package scheduleinput compiles supported human or cron source expressions
// into the scheduler's single RRULE/anchor execution model.
package scheduleinput

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/cron"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

// Syntax is the source grammar selected for a recurring expression.
type Syntax = cron.Syntax

const (
	// SyntaxCron selects the supported cron dialect without fallback.
	SyntaxCron = cron.SyntaxCron
	// SyntaxHuman selects the human schedule grammar without fallback.
	SyntaxHuman = cron.SyntaxHuman
)

// ErrInvalidSyntax identifies a request hint outside human or cron.
var ErrInvalidSyntax = errors.New("schedule syntax must be human or cron")

// Input is one successfully compiled recurring source expression.
type Input struct {
	Expression string
	Syntax     Syntax
	Schedule   domain.Schedule
}

// Parse selects exactly one source grammar, compiles it through the existing
// recurrence model, and retains the normalized submitted expression for later
// editing. An empty hint enables structural detection.
func Parse(input string, hint Syntax, tz string, now time.Time) (Input, error) {
	expression := strings.TrimSpace(input)
	syntax, err := selectSyntax(expression, hint)
	if err != nil {
		return Input{}, err
	}

	if syntax == SyntaxHuman {
		sch, err := schedule.Parse(expression, tz, now)
		if err != nil {
			return Input{}, err
		}
		return Input{Expression: expression, Syntax: syntax, Schedule: sch}, nil
	}

	phrase, bad, err := cron.Explain(expression)
	if err != nil {
		return Input{}, fmt.Errorf("schedule input: %w", err)
	}
	if bad.Reason != "" {
		return Input{}, fmt.Errorf("schedule input: cron: %s", bad.Reason)
	}
	sch, err := schedule.Parse(phrase, tz, now)
	if err != nil {
		return Input{}, fmt.Errorf("schedule input: parse converted cron: %w", err)
	}
	sch.Expression = expression
	return Input{Expression: expression, Syntax: syntax, Schedule: sch}, nil
}

func selectSyntax(input string, hint Syntax) (Syntax, error) {
	normalized := Syntax(strings.ToLower(strings.TrimSpace(string(hint))))
	switch normalized {
	case "":
		return cron.DetectSyntax(input), nil
	case SyntaxHuman, SyntaxCron:
		return normalized, nil
	default:
		return "", fmt.Errorf("%w, got %q", ErrInvalidSyntax, strings.TrimSpace(string(hint)))
	}
}

// SourceSyntax derives response identity from a retained recurring expression.
// It returns empty for one-offs, event schedules, and legacy expressionless
// rows. The value is response metadata and is never an execution input.
func SourceSyntax(sch domain.Schedule) Syntax {
	if sch.Kind != domain.ScheduleRecurring || strings.TrimSpace(sch.Expression) == "" {
		return ""
	}
	return cron.DetectSyntax(sch.Expression)
}
