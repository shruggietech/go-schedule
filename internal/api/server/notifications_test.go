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
	rec := doJSON(t, s, http.MethodPost, "/v1/notification-channels", NotificationChannelCreateRequest{Name: "Build alerts", Endpoint: "https://user:pass@example.test/hook?token=secret", Authorization: "Bearer secret", HealthIntervalSeconds: 300})
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
	if channel.EndpointSummary != "https://example.test" || !channel.HasAuthorization || channel.HealthIntervalSeconds != 300 {
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
	rec = doJSON(t, s, http.MethodPut, "/v1/tasks/"+task.Task.ID+"/notifications", NotificationAssignmentsRequest{Assignments: []NotificationAssignmentInput{{ChannelID: channel.ID, OnSuccess: true, OnFailure: true, FailureThreshold: 3, OnFailureToStart: true, DurationThresholdSeconds: 600, OnRecovery: true, ReminderIntervalSeconds: 3600, QuietPeriodSeconds: 600}}})
	if rec.Code != http.StatusOK {
		t.Fatalf("task assignment status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodGet, "/v1/tasks/"+task.Task.ID+"/notifications/effective", nil)
	if err := json.Unmarshal(rec.Body.Bytes(), &policy); err != nil {
		t.Fatal(err)
	}
	if policy.SourceScopeType != domain.NotificationScopeTask || !policy.Assignments[0].OnSuccess || policy.Assignments[0].FailureThreshold != 3 || !policy.Assignments[0].OnFailureToStart || policy.Assignments[0].DurationThresholdSeconds != 600 || !policy.Assignments[0].OnRecovery || policy.Assignments[0].ReminderIntervalSeconds != 3600 || policy.Assignments[0].QuietPeriodSeconds != 600 {
		t.Fatalf("policy=%+v", policy)
	}
	rec = doJSON(t, s, http.MethodPut, "/v1/tasks/"+task.Task.ID+"/notifications", NotificationAssignmentsRequest{Assignments: []NotificationAssignmentInput{{ChannelID: channel.ID, OnFailure: true, ReminderIntervalSeconds: 1}}})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid reminder status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodPut, "/v1/tasks/"+task.Task.ID+"/notifications", NotificationAssignmentsRequest{Assignments: []NotificationAssignmentInput{{ChannelID: channel.ID, DurationThresholdSeconds: 1}}})
	if rec.Code != http.StatusOK {
		t.Fatalf("sub-minute duration status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestSMTPChannelAuthoringAndRedaction(t *testing.T) {
	s := newTestServer(t)
	config := domain.SMTPConfig{Host: "mail.example.test", Port: 587, TLSMode: "starttls", From: "sender@example.test", Recipients: []string{"ops@example.test"}, TimeoutSeconds: 10}
	rec := doJSON(t, s, http.MethodPost, "/v1/notification-channels", NotificationChannelCreateRequest{Name: "Mail", Kind: domain.NotificationChannelSMTP, SMTP: &config, SMTPAuth: &domain.SMTPAuth{Username: "secret-user", Password: "secret-password"}})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create %d %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "secret-user") || strings.Contains(rec.Body.String(), "secret-password") {
		t.Fatal("credentials leaked")
	}
	var channel domain.NotificationChannel
	if err := json.Unmarshal(rec.Body.Bytes(), &channel); err != nil {
		t.Fatal(err)
	}
	if channel.SMTP == nil || !channel.HasAuthorization {
		t.Fatalf("channel %+v", channel)
	}
	config.Host = "next.example.test"
	rec = doJSON(t, s, http.MethodPatch, "/v1/notification-channels/"+channel.ID, NotificationChannelUpdateRequest{SMTP: &config})
	if rec.Code != http.StatusOK {
		t.Fatalf("update %d %s", rec.Code, rec.Body.String())
	}
	endpoint := "https://example.test"
	rec = doJSON(t, s, http.MethodPatch, "/v1/notification-channels/"+channel.ID, NotificationChannelUpdateRequest{Endpoint: &endpoint})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("wrong kind %d", rec.Code)
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/notification-channels/"+channel.ID+"/rotate", NotificationChannelRotateRequest{SMTPAuth: &domain.SMTPAuth{}})
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), "\"has_authorization\":true") {
		t.Fatalf("clear %d %s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, s, http.MethodPost, "/v1/notification-channels", NotificationChannelCreateRequest{Name: "Wrong", Endpoint: "https://example.test", SMTP: &config})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("webhook SMTP fields accepted: %d", rec.Code)
	}
}
