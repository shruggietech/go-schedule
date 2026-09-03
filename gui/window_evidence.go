package gui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
)

const windowEvidencePathEnv = "GOSCHEDULE_WINDOW_EVIDENCE_PATH"

type windowEvidence struct {
	SchemaVersion int       `json:"schema_version"`
	ProcessID     int       `json:"process_id"`
	CapturedAt    time.Time `json:"captured_at"`
	ContentWidth  float32   `json:"content_width"`
	ContentHeight float32   `json:"content_height"`
	CanvasScale   float32   `json:"canvas_scale"`
}

// writeWindowEvidence emits exact in-process Fyne canvas measurements only
// when the attended release harness supplies an absolute, unused output path.
func writeWindowEvidence(size fyne.Size, scale float32) error {
	target := os.Getenv(windowEvidencePathEnv)
	if target == "" {
		return nil
	}
	if !filepath.IsAbs(target) {
		return fmt.Errorf("%s must be an absolute path", windowEvidencePathEnv)
	}
	parent := filepath.Dir(target)
	info, err := os.Stat(parent)
	if err != nil {
		return fmt.Errorf("window evidence directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("window evidence parent is not a directory: %s", parent)
	}
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("window evidence path already exists: %s", target)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect window evidence path: %w", err)
	}

	temporary, err := os.CreateTemp(parent, ".goschedule-window-evidence-*")
	if err != nil {
		return fmt.Errorf("create temporary window evidence: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return fmt.Errorf("protect temporary window evidence: %w", err)
	}
	record := windowEvidence{
		SchemaVersion: 1,
		ProcessID:     os.Getpid(),
		CapturedAt:    time.Now().UTC(),
		ContentWidth:  size.Width,
		ContentHeight: size.Height,
		CanvasScale:   scale,
	}
	encoder := json.NewEncoder(temporary)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(record); err != nil {
		temporary.Close()
		return fmt.Errorf("encode window evidence: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("sync window evidence: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close window evidence: %w", err)
	}
	if err := os.Link(temporaryName, target); err != nil {
		return fmt.Errorf("publish window evidence without overwrite: %w", err)
	}
	return nil
}
