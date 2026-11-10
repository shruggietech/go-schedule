package cli

import (
	"reflect"
	"testing"
)

func TestJoinRecoveryFlagsPreserveExplicitConditionsAndClearing(t *testing.T) {
	body, _, calls, err := executeAuthoring(t, []string{"join", "create", "--name", "Exports", "--entry", "entry", "--target", "target", "--branches-file", "-", "--recovery", "repair", "--recovery-on", "failure,canceled,skipped,not_started", "--cleanup", "cleanup", "--json"}, "[{\"id\":\"a\",\"name\":\"A\",\"task_id\":\"branch\"}]")
	if err != nil || calls != 1 {
		t.Fatalf("new authoring flags err=%v calls=%d", err, calls)
	}
	if body["recovery_task_id"] != "repair" || body["cleanup_task_id"] != "cleanup" || !reflect.DeepEqual(body["recovery_on"], []any{"failure", "canceled", "skipped", "not_started"}) {
		t.Fatalf("routing fields=%#v", body)
	}
	body, _, calls, err = executeAuthoring(t, []string{"join", "update", "join", "--recovery=", "--cleanup=", "--recovery-on", "none", "--json"}, "")
	if err != nil || calls != 1 || body["recovery_task_id"] != "" || body["cleanup_task_id"] != "" || !reflect.DeepEqual(body["recovery_on"], []any{}) {
		t.Fatalf("clear packet=%#v calls=%d err=%v", body, calls, err)
	}
}
