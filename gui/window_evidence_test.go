package gui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
)

func TestWriteWindowEvidenceRequiresOptIn(t *testing.T) {
	t.Setenv(windowEvidencePathEnv, "")
	if err := writeWindowEvidence(fyne.NewSize(1280, 800), 1.5); err != nil {
		t.Fatalf("writeWindowEvidence() = %v", err)
	}
}

func TestWriteWindowEvidenceRecordsExactCanvasMetrics(t *testing.T) {
	path := filepath.Join(t.TempDir(), "window.json")
	t.Setenv(windowEvidencePathEnv, path)

	if err := writeWindowEvidence(fyne.NewSize(1280, 800), 1.5); err != nil {
		t.Fatalf("writeWindowEvidence() = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got windowEvidence
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ContentWidth != 1280 || got.ContentHeight != 800 || got.CanvasScale != 1.5 {
		t.Fatalf("window evidence = %+v", got)
	}
	if got.ProcessID != os.Getpid() {
		t.Fatalf("process_id = %d, want %d", got.ProcessID, os.Getpid())
	}
	if got.CapturedAt.IsZero() {
		t.Fatal("captured_at is zero")
	}
}

func TestWriteWindowEvidenceRejectsRelativePath(t *testing.T) {
	t.Setenv(windowEvidencePathEnv, "relative.json")
	if err := writeWindowEvidence(fyne.NewSize(1280, 800), 1); err == nil {
		t.Fatal("writeWindowEvidence() unexpectedly accepted relative path")
	}
}
