package main

import (
	"log/slog"
	"testing"

	"github.com/shruggietech/go-schedule/internal/logbus"
)

func TestLogDaemonReady(t *testing.T) {
	ring := logbus.NewRing(10)
	log := slog.New(logbus.NewHandler(slog.LevelInfo, ring, nil, nil))
	logDaemonReady(log, `\\.\pipe\go-schedule`, `C:\Schedule Data\schedule.db`, `C:\Schedule Data\日志\goschedule.log`)

	records := ring.Snapshot("", 0)
	if len(records) != 1 {
		t.Fatalf("startup records = %d, want 1", len(records))
	}
	got := records[0]
	if got.Message != "daemon startup complete" {
		t.Fatalf("message = %q", got.Message)
	}
	wantAttrs := map[string]string{
		"endpoint": `\\.\pipe\go-schedule`,
		"db":       `C:\Schedule Data\schedule.db`,
		"log_path": `C:\Schedule Data\日志\goschedule.log`,
	}
	for key, want := range wantAttrs {
		if value, ok := got.Attrs[key].(string); !ok || value != want {
			t.Errorf("%s = %#v, want %q", key, got.Attrs[key], want)
		}
	}
}
