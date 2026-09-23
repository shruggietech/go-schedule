//go:build linux

package service

import "testing"

func TestParseSystemdState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, raw string
		want      State
		wantErr   bool
	}{
		{"running", "LoadState=loaded\nActiveState=active\nSubState=running\n", StateRunning, false},
		{"stopped", "LoadState=loaded\nActiveState=inactive\nSubState=dead\n", StateStopped, false},
		{"failed", "LoadState=loaded\nActiveState=failed\nSubState=failed\n", StateStopped, false},
		{"missing", "LoadState=not-found\nActiveState=inactive\n", StateNotInstalled, false},
		{"starting", "LoadState=loaded\nActiveState=activating\n", StateStarting, false},
		{"stopping", "LoadState=loaded\nActiveState=deactivating\n", StateStopping, false},
		{"malformed", "ActiveState=active\n", StateUnknown, true},
		{"unknown", "LoadState=loaded\nActiveState=mystery\n", StateUnknown, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseSystemdState(tt.raw)
			if got != tt.want || (err != nil) != tt.wantErr {
				t.Fatalf("parseSystemdState() = %q, %v; want %q, error=%t", got, err, tt.want, tt.wantErr)
			}
		})
	}
}
