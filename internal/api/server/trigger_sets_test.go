package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestTriggerSetCreateRevealAndOrdinaryRedaction(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{Command: "echo"})
	rec := doJSON(t, s, http.MethodPost, "/v1/trigger-sets", TriggerSetCreateRequest{Name: "Callers", TargetTaskID: task.Task.ID, Count: 3})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created TriggerSetSecretResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if len(created.Members) != 3 || created.Members[0].Position != 1 || created.Members[2].Position != 3 {
		t.Fatalf("created=%+v", created)
	}
	keys := map[string]bool{}
	for _, member := range created.Members {
		if keys[member.Key] || !strings.Contains(member.Command, member.Key) {
			t.Fatalf("secret member=%+v", member)
		}
		keys[member.Key] = true
	}
	rec = doJSON(t, s, http.MethodGet, "/v1/trigger-sets", nil)
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), "gst_") || strings.Contains(rec.Body.String(), `"key"`) {
		t.Fatalf("ordinary list status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/trigger-sets/"+created.TriggerSet.ID+"/reveal", nil)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), created.Members[0].Key) {
		t.Fatalf("reveal status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestTriggerSetValidationAndBroadLifecycle(t *testing.T) {
	s := newTestServer(t)
	first := newTaskFor(t, s, TaskCreateRequest{Command: "echo"})
	second := newTaskFor(t, s, TaskCreateRequest{Command: "echo"})
	for _, count := range []int{0, 100} {
		rec := doJSON(t, s, http.MethodPost, "/v1/trigger-sets", TriggerSetCreateRequest{Name: "Bad", TargetTaskID: first.Task.ID, Count: count})
		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), `"field":"count"`) {
			t.Fatalf("count %d status=%d body=%s", count, rec.Code, rec.Body.String())
		}
	}
	rec := doJSON(t, s, http.MethodPost, "/v1/trigger-sets", TriggerSetCreateRequest{Name: "Pair", TargetTaskID: first.Task.ID, Count: 2})
	var created TriggerSetSecretResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	id := created.TriggerSet.ID
	rec = doJSON(t, s, http.MethodPatch, "/v1/trigger-sets/"+id, TriggerSetRetargetRequest{TargetTaskID: second.Task.ID})
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), second.Task.ID) {
		t.Fatalf("retarget status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/trigger-sets/"+id+"/disable", nil)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"enabled_count":0`) {
		t.Fatalf("disable status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/trigger-sets/"+id+"/rotate", nil)
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), created.Members[0].Key) {
		t.Fatalf("rotate status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodDelete, "/v1/trigger-sets/"+id, nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestTriggerSetMemberCannotBeRetargetedIndividually(t *testing.T) {
	s := newTestServer(t)
	first := newTaskFor(t, s, TaskCreateRequest{Command: "echo"})
	second := newTaskFor(t, s, TaskCreateRequest{Command: "echo"})
	rec := doJSON(t, s, http.MethodPost, "/v1/trigger-sets", TriggerSetCreateRequest{Name: "One", TargetTaskID: first.Task.ID, Count: 1})
	var created TriggerSetSecretResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	rec = doJSON(t, s, http.MethodPatch, "/v1/triggers/"+created.Members[0].TriggerID, TriggerUpdateRequest{TargetTaskID: &second.Task.ID})
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "Trigger Set retarget") {
		t.Fatalf("member retarget status=%d body=%s", rec.Code, rec.Body.String())
	}
}
