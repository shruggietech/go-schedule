package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/cron"
)

func TestCompositeCronScheduleSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goschedule.db")
	anchor := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	sch, bad, err := cron.Compile("*/10 9-17 * * MON,WED,FRI", "UTC", anchor)
	if err != nil || bad.Reason != "" {
		t.Fatalf("Compile: refusal=%q err=%v", bad.Reason, err)
	}

	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.CreateSchedule(&sch); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetSchedule(sch.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.RRULE != sch.RRULE || loaded.Expression != sch.Expression || loaded.HumanSummary != sch.HumanSummary ||
		loaded.Anchor == nil || !loaded.Anchor.Equal(*sch.Anchor) {
		t.Fatalf("loaded=%+v, want %+v", loaded, sch)
	}
}
