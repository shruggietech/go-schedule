package scheduleinput

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

var inputAnchor = time.Date(2026, 2, 25, 12, 0, 0, 0, time.UTC)

func TestParseAutoDetectsAndRetainsSource(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Syntax
	}{
		{name: "cron", input: "  0 9 * * 1-5  ", want: SyntaxCron},
		{name: "shorthand", input: "@daily", want: SyntaxCron},
		{name: "human", input: "weekdays at 09:00", want: SyntaxHuman},
		{name: "five-field human", input: "3rd wednesday monthly at 14:00", want: SyntaxHuman},
		{name: "last-weekday human", input: "last wednesday of the month at 14:00", want: SyntaxHuman},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input, "", "UTC", inputAnchor)
			if err != nil {
				t.Fatal(err)
			}
			if got.Syntax != tt.want {
				t.Errorf("Syntax = %q, want %q", got.Syntax, tt.want)
			}
			if got.Expression != strings.TrimSpace(tt.input) || got.Schedule.Expression != got.Expression {
				t.Errorf("retained expression = %q / %q, want %q",
					got.Expression, got.Schedule.Expression, strings.TrimSpace(tt.input))
			}
			if got.Schedule.RRULE == "" || got.Schedule.Anchor == nil {
				t.Errorf("compiled schedule is incomplete: %+v", got.Schedule)
			}
		})
	}
}

func TestParseExplicitHintSelectsOneParser(t *testing.T) {
	got, err := Parse("0 9 * * 1-5", Syntax(" CRON "), "UTC", inputAnchor)
	if err != nil || got.Syntax != SyntaxCron {
		t.Fatalf("normalized cron hint: result=%+v err=%v", got, err)
	}

	if _, err := Parse("weekdays at 09:00", SyntaxCron, "UTC", inputAnchor); err == nil ||
		!strings.Contains(strings.ToLower(err.Error()), "field") {
		t.Fatalf("forced cron fell back to human: %v", err)
	}
	if _, err := Parse("0 9 * * 1-5", SyntaxHuman, "UTC", inputAnchor); err == nil ||
		!strings.Contains(strings.ToLower(err.Error()), "could not understand") {
		t.Fatalf("forced human fell back to cron: %v", err)
	}
	if _, err := Parse("weekdays at 09:00", Syntax("yaml"), "UTC", inputAnchor); !errors.Is(err, ErrInvalidSyntax) {
		t.Fatalf("invalid hint error = %v, want ErrInvalidSyntax", err)
	}
}

func TestParseCronRefusalsNeverFallBack(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "61 9 * * *", want: "minute"},
		{input: "@reboot", want: "boot"},
		{input: "0 9 1 * 1", want: "either"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := Parse(tt.input, "", "UTC", inputAnchor)
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), tt.want) {
				t.Fatalf("Parse error = %v, want reason containing %q", err, tt.want)
			}
		})
	}
}

func TestParseHumanAndCronHaveIdenticalRuns(t *testing.T) {
	tests := []struct {
		name  string
		cron  string
		human string
	}{
		{name: "DST weekdays", cron: "0 9 * * 1-5", human: "weekdays at 09:00"},
		{name: "month boundary", cron: "0 9 31 * *", human: "on the 31st of every month at 09:00"},
		{name: "ordinal weekday", cron: "30 2 * * 5#5", human: "5th friday monthly at 02:30"},
		{name: "last weekday", cron: "30 2 * * 5L", human: "last friday of the month at 02:30"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cronInput, err := Parse(tt.cron, SyntaxCron, "America/New_York", inputAnchor)
			if err != nil {
				t.Fatal(err)
			}
			humanInput, err := Parse(tt.human, SyntaxHuman, "America/New_York", inputAnchor)
			if err != nil {
				t.Fatal(err)
			}
			// Cron and human are independent source grammars. Their internal rules
			// may use different faithful shapes; execution parity is the contract.
			a, err := schedule.UpcomingRuns(cronInput.Schedule, "America/New_York", domain.MissingDateSkip, inputAnchor, 20)
			if err != nil {
				t.Fatal(err)
			}
			b, err := schedule.UpcomingRuns(humanInput.Schedule, "America/New_York", domain.MissingDateSkip, inputAnchor, 20)
			if err != nil {
				t.Fatal(err)
			}
			if len(a) == 0 || len(a) != len(b) {
				t.Fatalf("run counts differ: %d != %d", len(a), len(b))
			}
			for i := range a {
				if !a[i].Equal(b[i]) {
					t.Fatalf("run %d differs: %s != %s", i, a[i], b[i])
				}
			}
		})
	}
}

func TestSourceSyntaxOmitsNonRecurringOrExpressionlessSchedules(t *testing.T) {
	if got := SourceSyntax(schedule.NewOneOff(inputAnchor.Add(time.Hour))); got != "" {
		t.Errorf("one-off source syntax = %q, want empty", got)
	}
	if got := SourceSyntax(domain.Schedule{Kind: domain.ScheduleRecurring}); got != "" {
		t.Errorf("expressionless source syntax = %q, want empty", got)
	}
	if got := SourceSyntax(domain.Schedule{Kind: domain.ScheduleRecurring, Expression: "0 9 * * *"}); got != SyntaxCron {
		t.Errorf("cron source syntax = %q, want cron", got)
	}
}
