package clievaluate

import "testing"

func TestHumanChannelCreationObservesIndependentPersistedIdentity(t *testing.T) {
	before := map[string]any{"channels": map[string]any{"prior": map[string]any{"id": "prior"}}, "tasks": map[string]any{}}
	after := map[string]any{"channels": map[string]any{"prior": map[string]any{"id": "prior"}, "created-channel": map[string]any{"id": "created-channel"}}, "tasks": map[string]any{"unrelated-task": map[string]any{"id": "unrelated-task"}}}
	id, err := newCommandIdentity("notification channel add", before, after)
	if err != nil || id != "created-channel" {
		t.Fatalf("human channel identity=%q err=%v", id, err)
	}
	if _, err := newCommandIdentity("notification channel add", before, before); err == nil {
		t.Fatal("unchanged state was accepted as channel creation")
	}
	after["channels"].(map[string]any)["second-channel"] = map[string]any{"id": "second-channel"}
	if _, err := newCommandIdentity("notification channel add", before, after); err == nil {
		t.Fatal("ambiguous channel creation was accepted")
	}
}
