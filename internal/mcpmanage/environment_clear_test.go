package mcpmanage

import (
	"encoding/json"
	"testing"
)

func TestClearSensitiveInputsPreservesEnvironmentClearOnWire(t *testing.T) {
	request := (TaskUpdate{ClearSensitive: true}).request()
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	if wire["clear_env"] != true || wire["stdin"] != "" {
		t.Fatalf("sensitive-input clear lost its transport intent: %s", data)
	}
}
