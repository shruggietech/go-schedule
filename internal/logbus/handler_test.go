package logbus

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type capturePub struct{ recs []domain.LogRecord }

func (c *capturePub) PublishLog(r domain.LogRecord) { c.recs = append(c.recs, r) }

func TestHandler_TeesAndMapsSeverity(t *testing.T) {
	ring := NewRing(10)
	var file bytes.Buffer
	pub := &capturePub{}
	log := slog.New(NewHandler(slog.LevelInfo, ring, &file, pub))

	log.Info("info msg")
	log.Warn("warn msg")
	log.Error("err msg", "exit_code", 1)

	// Ring has all three, newest first.
	snap := ring.Snapshot("", 0)
	if len(snap) != 3 {
		t.Fatalf("ring len = %d, want 3", len(snap))
	}
	if snap[0].Severity != domain.SeverityError {
		t.Errorf("newest severity = %q, want error", snap[0].Severity)
	}

	// Severity mapping across levels.
	if got := ring.Snapshot(domain.SeverityWarning, 0); len(got) != 1 || got[0].Message != "warn msg" {
		t.Errorf("warning filter = %+v", got)
	}

	// Publisher saw every record.
	if len(pub.recs) != 3 {
		t.Errorf("publisher saw %d, want 3", len(pub.recs))
	}

	// File got JSONL lines including the structured attr.
	out := file.String()
	if strings.Count(out, "\n") != 3 {
		t.Errorf("file line count: %q", out)
	}
	if !strings.Contains(out, `"exit_code":1`) {
		t.Errorf("attr not inlined in file: %q", out)
	}
}

func TestHandler_PromotesCorrelationAttrs(t *testing.T) {
	ring := NewRing(10)
	log := slog.New(NewHandler(slog.LevelInfo, ring, nil, nil))
	log.Error("task run failed", "task", "tsk_1", "run", "run_9", "error", "boom")

	rec := ring.Snapshot("", 1)[0]
	if rec.TaskID != "tsk_1" || rec.RunID != "run_9" {
		t.Errorf("correlation = task %q run %q", rec.TaskID, rec.RunID)
	}
	if rec.Attrs["error"] != "boom" {
		t.Errorf("attrs = %v", rec.Attrs)
	}
}

func TestHandler_RespectsLevel(t *testing.T) {
	ring := NewRing(10)
	h := NewHandler(slog.LevelWarn, ring, nil, nil)
	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("info should be disabled at warn level")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("error should be enabled at warn level")
	}
}
