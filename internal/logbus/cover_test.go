package logbus

import (
	"bytes"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// TestHandler_WithGroupAndWithAttrs exercises the slog WithGroup/With paths
// (handler source + accumulated attrs, plus correlation-attr promotion).
func TestHandler_WithGroupAndWithAttrs(t *testing.T) {
	ring := NewRing(10)
	base := slog.New(NewHandler(slog.LevelInfo, ring, nil, nil))
	log := base.WithGroup("engine").With("component", "executor", "task", "tsk_7")
	log.Info("hello", "k", "v")

	rec := ring.Snapshot("", 1)[0]
	if rec.Source != "engine" {
		t.Errorf("source = %q, want engine", rec.Source)
	}
	if rec.TaskID != "tsk_7" {
		t.Errorf("task promoted = %q, want tsk_7", rec.TaskID)
	}
	if rec.Attrs["component"] != "executor" || rec.Attrs["k"] != "v" {
		t.Errorf("attrs = %v", rec.Attrs)
	}
}

// TestHandler_NilLevelDefaultsInfo covers the nil-level default in NewHandler.
func TestHandler_NilLevelDefaultsInfo(t *testing.T) {
	ring := NewRing(4)
	log := slog.New(NewHandler(nil, ring, nil, nil))
	log.Debug("debug dropped")
	log.Info("info kept")
	if got := ring.Snapshot("", 0); len(got) != 1 || got[0].Message != "info kept" {
		t.Fatalf("nil level should default to Info: %+v", got)
	}
}

// TestMarshalLine_FileBranches covers source/task_id/run_id in the JSONL output.
func TestMarshalLine_FileBranches(t *testing.T) {
	ring := NewRing(4)
	var file bytes.Buffer
	base := slog.New(NewHandler(slog.LevelInfo, ring, &file, nil))
	base.WithGroup("executor").With("task", "t1", "run", "r1").Error("boom", "exit_code", 2)

	out := file.String()
	for _, want := range []string{`"source":"executor"`, `"task_id":"t1"`, `"run_id":"r1"`, `"exit_code":2`, `"severity":"error"`} {
		if !strings.Contains(out, want) {
			t.Errorf("file line missing %s: %s", want, out)
		}
	}
}

// TestRotatingWriter_DefaultsAndDoubleClose covers the default-arg branches and
// idempotent Close.
func TestRotatingWriter_DefaultsAndDoubleClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "l.log")
	w, err := NewRotatingWriter(path, 0, 0) // defaults: 10MiB / 5 files
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("x\n")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	// Second close must be safe (f already closed -> no panic, error tolerated).
	_ = w.Close()
}

// TestRing_AddFromSlogRecordShape is a light guard that a zero-value severity
// record round-trips through the ring.
func TestRing_AddZeroSeverity(t *testing.T) {
	r := NewRing(2)
	r.Add(domain.LogRecord{ID: "z", Message: "m"})
	if got := r.Snapshot("", 0); len(got) != 1 || got[0].ID != "z" {
		t.Fatalf("ring = %+v", got)
	}
}
