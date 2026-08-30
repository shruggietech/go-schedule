package cron

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseSecondsAndQuartzDaySemantics(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		seconds []int
		days    []int
	}{
		{"five field defaults to zero", "0 12 * * 1", []int{0}, []int{1}},
		{"wildcard step", "*/30 * * * * *", []int{0, 30}, []int{0, 1, 2, 3, 4, 5, 6}},
		{"value start step", "5/15 * * * * *", []int{5, 20, 35, 50}, []int{0, 1, 2, 3, 4, 5, 6}},
		{"quartz monday numeric", "0 0 12 ? * 2", []int{0}, []int{1}},
		{"quartz monday name", "0 0 12 ? * MON", []int{0}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Parse(tt.expr)
			if err != nil || !res.OK {
				t.Fatalf("Parse(%q) = %+v, %v", tt.expr, res, err)
			}
			if !reflect.DeepEqual(res.Spec.Second.Values, tt.seconds) || !reflect.DeepEqual(res.Spec.DOW.Values, tt.days) {
				t.Fatalf("seconds=%v days=%v", res.Spec.Second.Values, res.Spec.DOW.Values)
			}
		})
	}
}

func TestParseQuartzQuestionAndYearFailures(t *testing.T) {
	for _, tt := range []struct {
		expr     string
		contains string
	}{
		{"? 0 12 * * MON", "seconds"},
		{"0 0 12 ? * ?", "both"},
		{"0 0 12 ?/2 * MON", "complete"},
		{"0 12 ? * MON", "six-field"},
		{"0 0 12 ? * MON 2026", "year"},
		{"60 * * * * *", "seconds"},
	} {
		res, err := Parse(tt.expr)
		message := ""
		if err != nil {
			message = err.Error()
		} else {
			message = res.Bad.Reason
		}
		if !strings.Contains(strings.ToLower(message), tt.contains) {
			t.Errorf("Parse(%q) message %q, want %q", tt.expr, message, tt.contains)
		}
	}
}

func TestCompileSecondsUsesDurableRecurrence(t *testing.T) {
	anchor := time.Date(2026, 1, 5, 11, 59, 50, 0, time.UTC)
	sch, bad, err := Compile("*/30 * * * * *", "UTC", anchor)
	if err != nil || bad.Reason != "" {
		t.Fatalf("Compile = bad=%+v err=%v", bad, err)
	}
	if !strings.Contains(sch.RRULE, "BYSECOND=0,30") || sch.Expression != "*/30 * * * * *" {
		t.Fatalf("schedule = %+v", sch)
	}
	phrase, bad, err := Explain("5/15 * * * * *")
	if err != nil || bad.Reason != "" || !strings.Contains(phrase, "seconds 5, 20, 35, and 50") {
		t.Fatalf("Explain = %q bad=%+v err=%v", phrase, bad, err)
	}
}

func TestDetectSyntaxRecognizesSixFields(t *testing.T) {
	if got := DetectSyntax("0 0 12 ? * MON"); got != SyntaxCron {
		t.Fatalf("syntax = %q", got)
	}
}
