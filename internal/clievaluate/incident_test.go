package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/incident"
)

func TestIncidentWatchdogFailurePreservesOriginalCase(t *testing.T) {
	channel, err := incident.NewChannel(canonicalIncidentTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	s := &session{incidents: channel, daemon: &clifixture.Process{PID: 123}}
	d := &executionDiagnostics{}
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	d.phase(ctx, "after")
	finish := s.incidentWatchdog(ctx, d)
	<-ctx.Done()
	original := execution{Code: 0, Out: []byte("command succeeded"), State: []byte(`{"partial_after":{"tasks":{"owned":{}}},"after_complete":false}`), Cleanup: errors.New("original cleanup error")}
	got := finish(original)
	if got.Code != 0 || string(got.Out) != "command succeeded" || got.Cleanup != original.Cleanup {
		t.Fatal("collector changed primary result")
	}
	var state map[string]json.RawMessage
	if err := json.Unmarshal(got.State, &state); err != nil {
		t.Fatal(err)
	}
	if state["incident"] == nil || state["partial_after"] == nil || string(state["after_complete"]) != "false" {
		t.Fatalf("lost independent evidence: %s", got.State)
	}
}

func TestIncidentWatchdogDisabledLeavesStateUntouched(t *testing.T) {
	s := &session{}
	original := execution{Code: 7, State: []byte(`["original"]`)}
	got := s.incidentWatchdog(context.Background(), &executionDiagnostics{})(original)
	if string(got.State) != string(original.State) || got.Code != 7 {
		t.Fatal("disabled collector altered result")
	}
}

func TestIncidentSetupUnavailablePreservesCaseOutcome(t *testing.T) {
	s := &session{incidentError: errors.New("owned directory denied")}
	got := s.incidentWatchdog(context.Background(), &executionDiagnostics{})(execution{Code: 0, Out: []byte("original output"), State: []byte(`{"after_complete":false}`)})
	if got.Code != 0 || string(got.Out) != "original output" || !json.Valid(got.State) {
		t.Fatalf("collector setup changed original: %+v", got)
	}
}

func TestFastCaseFailureRetainsExplicitIncidentAvailability(t *testing.T) {
	channel, err := incident.NewChannel(canonicalIncidentTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	s := &session{incidents: channel, daemon: &clifixture.Process{PID: 123}}
	d := &executionDiagnostics{}
	d.phase(context.Background(), "before")
	got := s.incidentWatchdog(context.Background(), d)(execution{Blocked: "observation_unavailable", Detail: "original inspection failure", State: []byte(`{"before_complete":false}`)})
	var state map[string]json.RawMessage
	if err := json.Unmarshal(got.State, &state); err != nil {
		t.Fatal(err)
	}
	if state["incident"] == nil || got.Blocked != "observation_unavailable" || got.Detail != "original inspection failure" {
		t.Fatalf("fast failure lost original or collector disposition: %s", got.State)
	}
}

func TestExpectedFastRefusalDoesNotTriggerIncidentCollection(t *testing.T) {
	channel, err := incident.NewChannel(canonicalIncidentTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	s := &session{incidents: channel, daemon: &clifixture.Process{PID: 123}}
	got := s.incidentWatchdog(context.Background(), &executionDiagnostics{}, 7)(execution{Code: 7, Checks: []bool{true}, State: []byte(`{"expected_refusal":true}`)})
	var state map[string]json.RawMessage
	if err := json.Unmarshal(got.State, &state); err != nil {
		t.Fatal(err)
	}
	if state["incident"] != nil {
		t.Fatalf("expected refusal treated as an unexpected failure: %s", got.State)
	}
}

func TestResumeRejectsChangedIncidentSelection(t *testing.T) {
	prior := captureIdentity{Incidents: true}
	if err := compareResumeInputs(prior, captureIdentity{Incidents: false}); err == nil {
		t.Fatal("changed resumed incident mode")
	}
	if err := compareResumeInputs(prior, prior); err != nil {
		t.Fatal(err)
	}
}

func TestIncidentPreservesNonObjectAndMalformedState(t *testing.T) {
	for _, raw := range []string{`["observed"]`, `null`, `{"partial":`} {
		got := attachIncident(execution{State: []byte(raw), Code: 7}, map[string]any{"collection": "unavailable"})
		if got.Code != 7 || !json.Valid(got.State) {
			t.Fatalf("lost original state %q: %s", raw, got.State)
		}
		var state map[string]json.RawMessage
		if err := json.Unmarshal(got.State, &state); err != nil {
			t.Fatal(err)
		}
		if state["observation"] == nil && state["observation_raw"] == nil {
			t.Fatalf("original missing: %s", got.State)
		}
	}
}

func canonicalIncidentTempDir(t *testing.T) string {
	t.Helper()
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return parent
}
