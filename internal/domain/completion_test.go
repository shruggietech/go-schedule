package domain

import (
	"encoding/json"
	"testing"
)

func TestCompletionDomainValuesAndOptionalCorrelationJSON(t *testing.T) {
	for _, outcome := range []CompletionOutcome{CompletionOnSuccess, CompletionOnFailure, CompletionOnAny} {
		if outcome == "" {
			t.Fatal("completion outcome must have a stable wire value")
		}
	}
	for _, state := range []DeliveryState{DeliveryPending, DeliveryClaimed, DeliveryCompleted, DeliveryResolved} {
		if state == "" {
			t.Fatal("delivery state must have a stable wire value")
		}
	}
	plain, err := json.Marshal(Run{Trigger: TriggerSchedule})
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) == "" || containsJSONKey(plain, "source_task_id") || containsJSONKey(plain, "source_run_id") {
		t.Fatalf("plain run exposed empty correlation: %s", plain)
	}
	correlated, err := json.Marshal(Run{Trigger: TriggerCompletion, SourceTaskID: "source", SourceRunID: "run"})
	if err != nil {
		t.Fatal(err)
	}
	if !containsJSONKey(correlated, "source_task_id") || !containsJSONKey(correlated, "source_run_id") {
		t.Fatalf("completion run omitted correlation: %s", correlated)
	}
}

func containsJSONKey(data []byte, key string) bool {
	var object map[string]any
	if err := json.Unmarshal(data, &object); err != nil {
		return false
	}
	_, ok := object[key]
	return ok
}
