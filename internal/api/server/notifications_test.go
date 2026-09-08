package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type notificationWake struct{ count int }

func (w *notificationWake) Wake() { w.count++ }

func TestNotificationChannelAPIIsRedactedAndQueuesTest(t *testing.T) {
	s := newTestServer(t)
	wake := &notificationWake{}
	s.SetNotificationDispatcher(wake)
	rec := doJSON(t, s, http.MethodPost, "/v1/notification-channels", NotificationChannelCreateRequest{Name: "Build alerts", Endpoint: "https://user:pass@example.test/hook?token=secret", Authorization: "Bearer secret"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	for _, secret := range []string{"user:pass", "token=secret", "Bearer secret"} {
		if strings.Contains(rec.Body.String(), secret) {
			t.Fatalf("response disclosed %q: %s", secret, rec.Body.String())
		}
	}
	var channel domain.NotificationChannel
	if err := json.Unmarshal(rec.Body.Bytes(), &channel); err != nil {
		t.Fatal(err)
	}
	if channel.EndpointSummary != "https://example.test" || !channel.HasAuthorization {
		t.Fatalf("channel=%+v", channel)
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/notification-channels/"+channel.ID+"/test", nil)
	if rec.Code != http.StatusAccepted || wake.count != 1 || strings.Contains(rec.Body.String(), "secret") {
		t.Fatalf("status=%d wakes=%d body=%s", rec.Code, wake.count, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodGet, "/v1/notification-deliveries?channel="+channel.ID+"&limit=1", nil)
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), "Bearer secret") || !strings.Contains(rec.Body.String(), "go-schedule.webhook.v1") {
		t.Fatalf("history status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestNotificationAssignmentAPIExplainsTaskPrecedence(t *testing.T) {
	s := newTestServer(t)
	groupRec := doJSON(t, s, http.MethodPost, "/v1/groups", GroupCreateRequest{Name: "ops"})
	var group domain.Group
	if err := json.Unmarshal(groupRec.Body.Bytes(), &group); err != nil {
		t.Fatal(err)
	}
	task := newTaskFor(t, s, TaskCreateRequest{Name: "task", Command: "echo", GroupID: group.ID})
	channelRec := doJSON(t, s, http.MethodPost, "/v1/notification-channels", NotificationChannelCreateRequest{Name: "hook", Endpoint: "https://example.test/hook"})
	var channel domain.NotificationChannel
	if err := json.Unmarshal(channelRec.Body.Bytes(), &channel); err != nil {
		t.Fatal(err)
	}
	rec := doJSON(t, s, http.MethodPut, "/v1/groups/"+group.ID+"/notifications", NotificationAssignmentsRequest{Assignments: []NotificationAssignmentInput{{ChannelID: channel.ID, OnFailure: true}}})
	if rec.Code != http.StatusOK {
		t.Fatalf("group assignment status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodGet, "/v1/tasks/"+task.Task.ID+"/notifications/effective", nil)
	var policy domain.EffectiveNotificationPolicy
	if err := json.Unmarshal(rec.Body.Bytes(), &policy); err != nil {
		t.Fatal(err)
	}
	if policy.SourceScopeType != domain.NotificationScopeGroup || policy.SourceScopeID != group.ID || len(policy.Assignments) != 1 {
		t.Fatalf("policy=%+v", policy)
	}
	rec = doJSON(t, s, http.MethodPut, "/v1/tasks/"+task.Task.ID+"/notifications", NotificationAssignmentsRequest{Assignments: []NotificationAssignmentInput{{ChannelID: channel.ID, OnSuccess: true}}})
	if rec.Code != http.StatusOK {
		t.Fatalf("task assignment status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodGet, "/v1/tasks/"+task.Task.ID+"/notifications/effective", nil)
	if err := json.Unmarshal(rec.Body.Bytes(), &policy); err != nil {
		t.Fatal(err)
	}
	if policy.SourceScopeType != domain.NotificationScopeTask || !policy.Assignments[0].OnSuccess {
		t.Fatalf("policy=%+v", policy)
	}
}
