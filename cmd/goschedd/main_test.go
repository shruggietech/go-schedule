package main

import (
	"log/slog"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/ipc"
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

func TestLogIPCAccessRestricted(t *testing.T) {
	ring := logbus.NewRing(10)
	log := slog.New(logbus.NewHandler(slog.LevelInfo, ring, nil, nil))
	logIPCAccess(log, ipc.AccessInfo{Mode: ipc.AccessModeRestricted, AdminGroup: "goschedadmin"})

	records := ring.Snapshot("", 0)
	if len(records) != 1 || records[0].Message != "IPC access configured" {
		t.Fatalf("records = %+v", records)
	}
	if records[0].Severity != domain.SeverityInfo {
		t.Fatalf("severity = %q, want info", records[0].Severity)
	}
	if records[0].Attrs["access_mode"] != string(ipc.AccessModeRestricted) || records[0].Attrs["admin_group"] != "goschedadmin" {
		t.Fatalf("attrs = %+v", records[0].Attrs)
	}
}

func TestLogIPCAccessCompatibilityWarnsOnce(t *testing.T) {
	ring := logbus.NewRing(10)
	log := slog.New(logbus.NewHandler(slog.LevelInfo, ring, nil, nil))
	logIPCAccess(log, ipc.AccessInfo{Mode: ipc.AccessModeCompatibility})

	records := ring.Snapshot("", 0)
	if len(records) != 1 || records[0].Message != "IPC compatibility mode enabled" {
		t.Fatalf("records = %+v", records)
	}
	if records[0].Severity != domain.SeverityWarning {
		t.Fatalf("severity = %q, want warning", records[0].Severity)
	}
	if records[0].Attrs["access_mode"] != string(ipc.AccessModeCompatibility) || records[0].Attrs["admin_group"] != "" {
		t.Fatalf("attrs = %+v", records[0].Attrs)
	}
}

func TestIPCAccessEvidencePrecedesReadyRecord(t *testing.T) {
	ring := logbus.NewRing(10)
	log := slog.New(logbus.NewHandler(slog.LevelInfo, ring, nil, nil))
	logIPCAccess(log, ipc.AccessInfo{Mode: ipc.AccessModeRestricted, AdminGroup: "goschedadmin"})
	logDaemonReady(log, "endpoint", "db", "log")

	records := ring.Snapshot("", 0)
	if len(records) != 2 || records[0].Message != "daemon startup complete" || records[1].Message != "IPC access configured" {
		t.Fatalf("newest-first startup records = %+v", records)
	}
}
