package cron

import (
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/schedule"
)

func TestCompositeCronCanonicalExport(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0 9,17 * * *", "0 9,17 * * *"},
		{"30 8-17 * * MON-FRI", "30 8-17 * * 1-5"},
		{"*/10 9-17 * * *", "*/10 9-17 * * *"},
		{"10-20/2 * * * *", "10,12,14,16,18,20 * * * *"},
		{"0 0 1,15 JAN,MAR *", "0 0 1,15 1,3 *"},
		{"0 12 * * 7,MON,WED,FRI", "0 12 * * 0,1,3,5"},
		{"0-59 0-23 1-31 1-12 *", "* * * * *"},
		{"*/30 * * * * *", "*/30 * * * * *"},
		{"5/15 0 9 ? * 2", "5,20,35,50 0 9 * * 2"},
	}
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sch, bad, err := Compile(tt.input, "UTC", anchor)
			if err != nil || bad.Reason != "" {
				t.Fatalf("Compile: refusal=%q err=%v", bad.Reason, err)
			}
			got, bad, ok := ExportSchedule(sch, domain.MissingDateSkip)
			if !ok || bad.Reason != "" || got != tt.want {
				t.Fatalf("Export=%q ok=%v refusal=%q, want %q", got, ok, bad.Reason, tt.want)
			}
			roundTrip, bad, err := Compile(got, "UTC", anchor)
			if err != nil || bad.Reason != "" {
				t.Fatalf("round trip refusal=%q err=%v", bad.Reason, err)
			}
			before, err := schedule.UpcomingRuns(sch, "UTC", domain.MissingDateSkip, anchor, 20)
			if err != nil {
				t.Fatal(err)
			}
			after, err := schedule.UpcomingRuns(roundTrip, "UTC", domain.MissingDateSkip, anchor, 20)
			if err != nil {
				t.Fatal(err)
			}
			for i := range before {
				if i >= len(after) || !before[i].Equal(after[i]) {
					t.Fatalf("round trip runs=%v, want %v", after, before)
				}
			}
		})
	}
}

func TestCompositeCronExportIsSourceIndependent(t *testing.T) {
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sch, bad, err := Compile("0 9,17 * * *", "UTC", anchor)
	if err != nil || bad.Reason != "" {
		t.Fatal(err)
	}
	sch.Expression = "malicious or stale display text"
	got, bad, ok := ExportSchedule(sch, domain.MissingDateSkip)
	if !ok || bad.Reason != "" || got != "0 9,17 * * *" {
		t.Fatalf("Export=%q ok=%v refusal=%q", got, ok, bad.Reason)
	}
}

func TestCompositeCronExportRefusesBehavioralMissingDatePolicy(t *testing.T) {
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sch, _, err := Compile("0 9 1,31 * *", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	if _, bad, ok := ExportSchedule(sch, domain.MissingDateLastValid); ok || !strings.Contains(bad.Reason, "missing-date") {
		t.Fatalf("ok=%v refusal=%q, want missing-date refusal", ok, bad.Reason)
	}

	sch, _, err = Compile("0 9 1,15 * *", "UTC", anchor)
	if err != nil {
		t.Fatal(err)
	}
	if got, bad, ok := ExportSchedule(sch, domain.MissingDateLastValid); !ok || bad.Reason != "" || got != "0 9 1,15 * *" {
		t.Fatalf("policy-inert export=%q ok=%v refusal=%q", got, ok, bad.Reason)
	}
}

func TestCompositeCronExportRejectsInvalidStoredFieldValue(t *testing.T) {
	sch := domain.Schedule{
		Kind:  domain.ScheduleRecurring,
		RRULE: "FREQ=DAILY;INTERVAL=1;BYHOUR=24;BYMINUTE=0;BYSECOND=0",
	}
	if _, bad, ok := ExportSchedule(sch, domain.MissingDateSkip); ok || bad.Reason == "" {
		t.Fatalf("ok=%v refusal=%q, want invalid recurrence refusal", ok, bad.Reason)
	}
}
