package mcpmanage

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJoinManageRecoveryConditionsPreserveIntent(t *testing.T) {
	for _, text := range []string{`{"name":"X","recovery_task_id":"repair","recovery_on":[],"cleanup_task_id":"clean"}`, `{"name":"X","recovery_task_id":"repair","recovery_on":["failure","canceled","skipped","not_started"],"cleanup_task_id":"clean"}`} {
		var input JoinCreate
		if err := json.Unmarshal([]byte(text), &input); err != nil {
			t.Fatal(err)
		}
		request := input.request()
		if request.RecoveryOn == nil || request.RecoveryTaskID != "repair" || request.CleanupTaskID != "clean" {
			t.Fatalf("request=%+v", request)
		}
		raw, _ := json.Marshal(request)
		if strings.Contains(text, `"recovery_on":[]`) && !strings.Contains(string(raw), `"recovery_on":[]`) {
			t.Fatalf("lost explicit empty: %s", raw)
		}
	}
	var update JoinUpdate
	if err := json.Unmarshal([]byte(`{"recovery_task_id":"","recovery_on":[],"cleanup_task_id":""}`), &update); err != nil {
		t.Fatal(err)
	}
	r := update.request()
	if r.RecoveryTaskID == nil || *r.RecoveryTaskID != "" || r.CleanupTaskID == nil || *r.CleanupTaskID != "" || r.RecoveryOn == nil || len(*r.RecoveryOn) != 0 {
		t.Fatalf("clear=%+v", r)
	}
	var omitted JoinCreate
	var empty JoinCreate
	json.Unmarshal([]byte(`{"recovery_task_id":"repair"}`), &omitted)
	json.Unmarshal([]byte(`{"recovery_task_id":"repair","recovery_on":[]}`), &empty)
	if fingerprint(omitted) == fingerprint(empty) {
		t.Fatal("dedup fingerprint conflated default and explicit no causes")
	}
}
