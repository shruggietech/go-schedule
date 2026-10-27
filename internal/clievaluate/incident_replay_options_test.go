package clievaluate

import (
	"testing"
	"time"
)

func TestIncidentReplayOptions(t *testing.T) {
	for _, test := range []struct {
		name, timeout, incidents string
		wantTimeout              time.Duration
		wantIncidents            bool
		wantError                bool
	}{
		{name: "supported defaults", wantTimeout: 120 * time.Second, wantIncidents: true},
		{name: "historical diagnostic", timeout: "15s", incidents: "0", wantTimeout: 15 * time.Second},
		{name: "explicit current", timeout: "120s", incidents: "1", wantTimeout: 120 * time.Second, wantIncidents: true},
		{name: "invalid duration", timeout: "later", wantError: true},
		{name: "zero duration", timeout: "0s", wantError: true},
		{name: "negative duration", timeout: "-1s", wantError: true},
		{name: "submillisecond duration", timeout: "1ns", wantError: true},
		{name: "fractional millisecond", timeout: "1500us", wantError: true},
		{name: "overflow duration", timeout: "999999999999999h", wantError: true},
		{name: "invalid collection", incidents: "maybe", wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			options, err := parseReplayOptions(test.timeout, test.incidents)
			if (err != nil) != test.wantError {
				t.Fatalf("options error=%v, want error=%t", err, test.wantError)
			}
			if err == nil && (options.caseTimeout != test.wantTimeout || options.incidents != test.wantIncidents) {
				t.Fatalf("options=%+v, want duration=%s collection=%t", options, test.wantTimeout, test.wantIncidents)
			}
		})
	}
}
