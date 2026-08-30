package cron

import (
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func TestConvertCronToHuman(t *testing.T) {
	tests := []struct {
		name  string
		input string
		to    Syntax
		want  string
	}{
		{name: "automatic", input: "0 9 * * 1-5", want: "weekdays at 09:00"},
		{name: "surrounding whitespace", input: "  0 9 * * 1-5  ", want: "weekdays at 09:00"},
		{name: "forced", input: "0 9 * * 1-5", to: SyntaxHuman, want: "weekdays at 09:00"},
		{name: "shorthand", input: "@daily", want: "every day at 00:00"},
		{name: "last weekday", input: "0 9 * * 5L", want: "last friday of the month at 09:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.input, tt.to)
			if err != nil {
				t.Fatal(err)
			}
			if got.InputSyntax != SyntaxCron || got.OutputSyntax != SyntaxHuman {
				t.Fatalf("syntax = %q -> %q, want cron -> human", got.InputSyntax, got.OutputSyntax)
			}
			if got.Input != strings.TrimSpace(tt.input) {
				t.Errorf("normalized input = %q, want %q", got.Input, strings.TrimSpace(tt.input))
			}
			if got.Output != tt.want || got.RefusalReason != "" {
				t.Errorf("conversion = %+v, want output %q", got, tt.want)
			}
		})
	}
}

func TestConvertCronRefusalNeverFallsBackToHuman(t *testing.T) {
	tests := []struct {
		input      string
		wantReason string
	}{
		{input: "61 9 * * *", wantReason: "minute"},
		{input: "0 9 * * FUNDAY", wantReason: "day-of-week"},
	}
	for _, tt := range tests {
		got, err := Convert(tt.input, "")
		if err != nil {
			t.Fatal(err)
		}
		if got.InputSyntax != SyntaxCron {
			t.Errorf("%q classified as %q, want cron", tt.input, got.InputSyntax)
		}
		if got.Output != "" || !strings.Contains(strings.ToLower(got.RefusalReason), tt.wantReason) {
			t.Errorf("conversion of %q = %+v, want refusal containing %q", tt.input, got, tt.wantReason)
		}
	}
}

func TestConvertRejectsUnknownDestination(t *testing.T) {
	if _, err := Convert("0 9 * * *", Syntax("yaml")); err == nil {
		t.Fatal("unknown destination was accepted")
	}
}

func TestConvertAutoDetectionKeepsFiveFieldHumanGrammar(t *testing.T) {
	got, err := Convert("3rd wednesday monthly at 14:00", "")
	if err != nil {
		t.Fatal(err)
	}
	if got.InputSyntax != SyntaxHuman || got.Output != "0 14 * * 3#3" || got.RefusalReason != "" {
		t.Fatalf("five-field human phrase was misclassified: %+v", got)
	}
}

func TestDetectSyntaxMatchesConverterClassification(t *testing.T) {
	tests := []struct {
		input string
		want  Syntax
	}{
		{input: "  0 9 * * 1-5  ", want: SyntaxCron},
		{input: "@daily", want: SyntaxCron},
		{input: "61 9 * * *", want: SyntaxCron},
		{input: "weekdays at 09:00", want: SyntaxHuman},
		{input: "every 15 minutes from 9am", want: SyntaxHuman},
		{input: "3rd wednesday monthly at 14:00", want: SyntaxHuman},
		{input: "last wednesday of the month at 14:00", want: SyntaxHuman},
	}
	for _, tt := range tests {
		if got := DetectSyntax(tt.input); got != tt.want {
			t.Errorf("DetectSyntax(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestConvertHumanToCanonicalCron(t *testing.T) {
	tests := []struct {
		name  string
		input string
		to    Syntax
		want  string
	}{
		{name: "automatic weekdays", input: "weekdays at 09:00", want: "0 9 * * 1-5"},
		{name: "forced weekdays", input: "weekdays at 09:00", to: SyntaxCron, want: "0 9 * * 1-5"},
		{name: "monthly date", input: "on the 31st of every month at 09:00", want: "0 9 31 * *"},
		{name: "aligned interval", input: "every 15 minutes starting at 00:00", want: "*/15 * * * *"},
		{name: "five-field human interval", input: "every 15 minutes from 9am", want: "*/15 * * * *"},
		{name: "hourly nonzero minute phase", input: "every hour starting at 00:30", want: "30 * * * *"},
		{name: "multi-hour nonzero minute phase", input: "every 2 hours starting at 08:30", want: "30 */2 * * *"},
		{name: "ordinal weekday", input: "3rd wednesday monthly at 14:00", want: "0 14 * * 3#3"},
		{name: "last weekday", input: "last wednesday of the month at 14:00", want: "0 14 * * 3L"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.input, tt.to)
			if err != nil {
				t.Fatal(err)
			}
			if got.InputSyntax != SyntaxHuman || got.OutputSyntax != SyntaxCron {
				t.Fatalf("syntax = %q -> %q, want human -> cron", got.InputSyntax, got.OutputSyntax)
			}
			if got.Output != tt.want || got.RefusalReason != "" {
				t.Errorf("conversion = %+v, want output %q", got, tt.want)
			}
		})
	}
}

func TestConvertHumanRefusesImplicitOrLossyTiming(t *testing.T) {
	tests := []struct {
		input      string
		wantReason string
	}{
		{input: "every 15 minutes", wantReason: "starting at"},
		{input: "every 15 minutes starting at 00:05", wantReason: "phase"},
		{input: "every 7 minutes starting at 00:00", wantReason: "divide"},
		{input: "every 2 hours starting at 09:30", wantReason: "phase"},
		{input: "every day", wantReason: "time of day"},
		{input: "tomorrow at 09:00", wantReason: "could not understand"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Convert(tt.input, SyntaxCron)
			if err != nil {
				t.Fatal(err)
			}
			if got.Output != "" || !strings.Contains(strings.ToLower(got.RefusalReason), tt.wantReason) {
				t.Errorf("conversion = %+v, want refusal containing %q", got, tt.wantReason)
			}
		})
	}
}

func TestConvertHumanRoundTripCrossesDSTAndMonthBoundary(t *testing.T) {
	const tz = "America/New_York"
	start := time.Date(2026, 2, 25, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)
	for _, input := range []string{
		"weekdays at 09:00",
		"on the 31st of every month at 09:00",
		"last friday of the month at 02:30",
	} {
		t.Run(input, func(t *testing.T) {
			converted, err := Convert(input, SyntaxCron)
			if err != nil || converted.RefusalReason != "" {
				t.Fatalf("Convert(%q): err=%v result=%+v", input, err, converted)
			}
			phrase, bad, err := Explain(converted.Output)
			if err != nil || bad.Reason != "" {
				t.Fatalf("Explain(%q): err=%v refusal=%q", converted.Output, err, bad.Reason)
			}
			original, err := schedule.Parse(input, tz, start)
			if err != nil {
				t.Fatal(err)
			}
			roundTrip, err := schedule.Parse(phrase, tz, start)
			if err != nil {
				t.Fatal(err)
			}
			a := runsBetween(t, original, tz, start, end)
			b := runsBetween(t, roundTrip, tz, start, end)
			if len(a) == 0 || len(a) != len(b) {
				t.Fatalf("run count changed: %d -> %d", len(a), len(b))
			}
			for i := range a {
				if !a[i].Equal(b[i]) {
					t.Fatalf("run %d changed: %v -> %v", i, a[i], b[i])
				}
			}
		})
	}
}

func runsBetween(t *testing.T, sch domain.Schedule, tz string, from, to time.Time) []time.Time {
	t.Helper()
	var out []time.Time
	for cursor, i := from, 0; i < 5000; i++ {
		next, ok, err := schedule.NextRun(sch, tz, domain.MissingDateSkip, cursor)
		if err != nil {
			t.Fatal(err)
		}
		if !ok || !next.Before(to) {
			break
		}
		out = append(out, next)
		cursor = next
	}
	return out
}
