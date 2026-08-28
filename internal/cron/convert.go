package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

// Syntax identifies one side of a string conversion.
type Syntax string

const (
	// SyntaxCron identifies the supported five-field cron dialect.
	SyntaxCron Syntax = "cron"
	// SyntaxHuman identifies go-schedule's human-readable schedule language.
	SyntaxHuman Syntax = "human"
)

// Conversion is the stable result of one local schedule-string conversion.
// Exactly one of Output and RefusalReason is populated.
type Conversion struct {
	InputSyntax   Syntax `json:"input_syntax"`
	OutputSyntax  Syntax `json:"output_syntax"`
	Input         string `json:"input"`
	Output        string `json:"output"`
	RefusalReason string `json:"refusal_reason"`
}

// Convert translates one schedule string to destination. An empty destination
// enables structural detection. Parse and fidelity failures are represented as
// refusals; an error is reserved for an unknown destination value.
func Convert(input string, destination Syntax) (Conversion, error) {
	normalized := strings.TrimSpace(input)
	inputSyntax, outputSyntax, err := conversionSyntax(normalized, destination)
	if err != nil {
		return Conversion{}, err
	}
	result := Conversion{
		InputSyntax:  inputSyntax,
		OutputSyntax: outputSyntax,
		Input:        normalized,
	}

	if inputSyntax == SyntaxHuman {
		convertHuman(&result)
		return result, nil
	}

	phrase, bad, err := Explain(normalized)
	switch {
	case err != nil:
		result.RefusalReason = err.Error()
	case bad.Reason != "":
		result.RefusalReason = bad.Reason
	default:
		result.Output = phrase
	}
	return result, nil
}

var conversionAnchor = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

func convertHuman(result *Conversion) {
	lower := strings.ToLower(result.Input)
	if schedule.IsSubDailyInterval(result.Input) &&
		!strings.Contains(lower, "starting at") && !strings.Contains(lower, " from ") {
		result.RefusalReason = "sub-daily intervals require an explicit starting at time for faithful cron conversion"
		return
	}

	sch, err := schedule.Parse(result.Input, "UTC", conversionAnchor)
	if err != nil {
		result.RefusalReason = err.Error()
		return
	}
	opt, err := rrule.StrToROption(sch.RRULE)
	if err != nil {
		result.RefusalReason = fmt.Sprintf("schedule recurrence could not be read: %v", err)
		return
	}
	if reason := implicitTimingRefusal(sch, opt); reason != "" {
		result.RefusalReason = reason
		return
	}

	expr, bad, ok := ExportSchedule(sch, domain.MissingDateSkip)
	if !ok {
		result.RefusalReason = bad.Reason
		return
	}
	result.Output = expr
}

func implicitTimingRefusal(sch domain.Schedule, opt *rrule.ROption) string {
	interval := opt.Interval
	if interval < 1 {
		interval = 1
	}
	if reason := subDailyPhaseRefusal(sch, opt, interval); reason != "" {
		return reason
	}
	switch opt.Freq {
	case rrule.SECONDLY:
		return "cron has no sub-minute resolution"
	case rrule.DAILY, rrule.WEEKLY, rrule.MONTHLY, rrule.YEARLY:
		if len(opt.Byhour) != 1 || len(opt.Byminute) != 1 {
			return "the schedule must state an explicit time of day for faithful cron conversion"
		}
	}
	return ""
}

func conversionSyntax(input string, destination Syntax) (Syntax, Syntax, error) {
	switch destination {
	case SyntaxCron:
		return SyntaxHuman, SyntaxCron, nil
	case SyntaxHuman:
		return SyntaxCron, SyntaxHuman, nil
	case "":
		if DetectSyntax(input) == SyntaxCron {
			return SyntaxCron, SyntaxHuman, nil
		}
		return SyntaxHuman, SyntaxCron, nil
	default:
		return "", "", fmt.Errorf("destination must be cron or human, got %q", destination)
	}
}

// DetectSyntax classifies one normalized or raw schedule string using the
// structural rule shared by conversion and task authoring. Classification does
// not validate the selected grammar and never falls back after a parse failure.
func DetectSyntax(input string) Syntax {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "@") || looksLikeCron(input) {
		return SyntaxCron
	}
	return SyntaxHuman
}

func looksLikeCron(input string) bool {
	fields := strings.Fields(input)
	if len(fields) != 5 || fields[0] == "" {
		return false
	}
	for _, r := range fields[0] {
		if r >= '0' && r <= '9' {
			continue
		}
		switch r {
		case '*', ',', '-', '/', '?':
			continue
		default:
			return false
		}
	}
	return true
}
