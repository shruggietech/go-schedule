package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/ipc"
)

func TestCompatibilityWarningExplainsEffectiveAccessAndConfiguration(t *testing.T) {
	var out bytes.Buffer
	logIPCAccess(slog.New(slog.NewJSONHandler(&out, nil)), ipc.AccessInfo{Mode: ipc.AccessModeCompatibility, AdminGroup: ""})
	var entry map[string]any
	if err := json.Unmarshal(out.Bytes(), &entry); err != nil {
		t.Fatal(err)
	}
	message, _ := entry["msg"].(string)
	for _, want := range []string{"local", "admin_group", "config"} {
		if !strings.Contains(message, want) {
			t.Errorf("warning lacks %q: %s", want, message)
		}
	}
	if entry["access_mode"] != "compatibility" || entry["admin_group"] != "" || entry["level"] != "WARN" {
		t.Fatal(entry)
	}
}
