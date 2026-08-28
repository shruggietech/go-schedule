package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestWriteLogsJSONPreservesBareArray(t *testing.T) {
	logs := []domain.LogRecord{{
		ID: "log-1", Time: time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC),
		Severity: domain.SeverityInfo, Source: "daemon", Message: "ready",
	}}
	var out bytes.Buffer
	if err := writeLogs(&out, logs, true); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(out.String())
	if !strings.HasPrefix(got, "[") || !strings.HasSuffix(got, "]") {
		t.Fatalf("JSON output = %q, want bare array", got)
	}
	if strings.Contains(got, `"logs"`) || strings.Contains(got, `"log_path"`) {
		t.Fatalf("JSON output leaked response envelope: %s", got)
	}
}

func TestWriteLogsHumanPreservesTable(t *testing.T) {
	logs := []domain.LogRecord{{
		Time:     time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC),
		Severity: domain.SeverityWarning, Source: "daemon", Message: "check path",
	}}
	var out bytes.Buffer
	if err := writeLogs(&out, logs, false); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"TIME", "SEVERITY", "SOURCE", "MESSAGE", "2026-08-28T12:00:00Z", "warning", "daemon", "check path"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("human output %q does not contain %q", out.String(), want)
		}
	}
}
