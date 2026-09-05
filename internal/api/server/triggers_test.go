package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestTriggerLifecycleRedactsOrdinaryResponsesAndRotates(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{Command: "echo"})
	rec := doJSON(t, s, http.MethodPost, "/v1/triggers", TriggerCreateRequest{Name: "hook", TargetTaskID: task.Task.ID})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created TriggerSecretResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(created.Key, "gst_") || !strings.Contains(created.Command, created.Key) {
		t.Fatalf("secret response = %+v", created)
	}
	rec = doJSON(t, s, http.MethodGet, "/v1/triggers", nil)
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), created.Key) || strings.Contains(rec.Body.String(), `"key"`) {
		t.Fatalf("list disclosed secret: status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/triggers/"+created.Trigger.ID+"/rotate", nil)
	var rotated TriggerSecretResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &rotated); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || rotated.Key == created.Key {
		t.Fatalf("rotate status=%d response=%+v", rec.Code, rotated)
	}
}

type triggerSchedulerStub struct{ err error }

func (triggerSchedulerStub) Reload()                                      {}
func (triggerSchedulerStub) RunNow(string) error                          { return nil }
func (s triggerSchedulerStub) FireExternalTrigger(string) (string, error) { return "trigger-id", s.err }

func TestFireTriggerUsesStableResponseAndDoesNotEchoKey(t *testing.T) {
	s := newTestServer(t)
	s.sched = triggerSchedulerStub{}
	key := "gst_sensitive_value"
	rec := doJSON(t, s, http.MethodPost, "/v1/triggers/fire", map[string]string{"key": key})
	if rec.Code != http.StatusAccepted || strings.Contains(rec.Body.String(), key) || !strings.Contains(rec.Body.String(), "trigger-id") {
		t.Fatalf("fire status=%d body=%s", rec.Code, rec.Body.String())
	}
}
