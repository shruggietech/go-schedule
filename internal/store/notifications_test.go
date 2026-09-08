package store

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func createNotificationTestTask(t *testing.T, st *Store, groupID string) domain.Task {
	t.Helper()
	task := domain.Task{Name: "notify me", GroupID: groupID, Command: "true", Enabled: true, Timezone: "UTC", State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	return task
}

func TestNotificationChannelLifecycleValidationAndFiltering(t *testing.T) {
	st := openMem(t)
	channels, err := st.ListNotificationChannels()
	if err != nil || len(channels) != 0 {
		t.Fatalf("empty channels=%+v err=%v", channels, err)
	}
	task := createNotificationTestTask(t, st, "")
	channel := createNotificationTestChannel(t, st, "original")
	duplicateChannel := channel
	if err := st.CreateNotificationChannel(&duplicateChannel); err == nil {
		t.Fatal("created duplicate notification channel")
	}
	if _, err := st.GetNotificationChannel("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get missing channel error=%v", err)
	}
	channels, err = st.ListNotificationChannels()
	if err != nil || len(channels) != 1 {
		t.Fatalf("channels=%+v err=%v", channels, err)
	}
	channel.Name, channel.Endpoint, channel.EndpointSummary = "updated", "https://new.example.test/hook", "https://new.example.test"
	if err := st.UpdateNotificationChannel(channel); err != nil {
		t.Fatal(err)
	}
	if err := st.RotateNotificationChannelAuthorization(channel.ID, "Bearer replacement"); err != nil {
		t.Fatal(err)
	}
	got, err := st.GetNotificationChannel(channel.ID)
	if err != nil || got.Name != "updated" || got.Endpoint != channel.Endpoint || got.Authorization != "Bearer replacement" {
		t.Fatalf("updated channel=%+v err=%v", got, err)
	}
	if err := st.SetNotificationChannelEnabled(channel.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTestNotificationDelivery(channel.ID); !errors.Is(err, ErrNotificationChannelDisabled) {
		t.Fatalf("disabled test error=%v", err)
	}
	if err := st.SetNotificationChannelEnabled(channel.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeNone, task.ID, nil); err == nil {
		t.Fatal("accepted invalid scope")
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, "missing", nil); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing scope error=%v", err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: "missing", OnFailure: true}}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing channel error=%v", err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: channel.ID}}); err == nil {
		t.Fatal("accepted assignment without outcome")
	}
	duplicate := []domain.NotificationAssignment{{ChannelID: channel.ID, OnFailure: true}, {ChannelID: channel.ID, OnSuccess: true}}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, duplicate); err == nil {
		t.Fatal("accepted duplicate channel assignment")
	}
	want := []domain.NotificationAssignment{{ChannelID: channel.ID, OnSuccess: true, OnFailure: true}}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, want); err != nil {
		t.Fatal(err)
	}
	assignments, err := st.ListNotificationAssignments(domain.NotificationScopeTask, task.ID)
	if err != nil || len(assignments) != 1 || !assignments[0].OnSuccess || !assignments[0].OnFailure {
		t.Fatalf("assignments=%+v err=%v", assignments, err)
	}
	if _, err := st.ListNotificationAssignments(domain.NotificationScopeNone, task.ID); err == nil {
		t.Fatal("listed invalid scope")
	}
	if _, err := st.EffectiveNotificationPolicy("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing effective policy error=%v", err)
	}
	isolated := createNotificationTestTask(t, st, "")
	policy, err := st.EffectiveNotificationPolicy(isolated.ID)
	if err != nil || policy.SourceScopeType != domain.NotificationScopeNone || len(policy.Assignments) != 0 {
		t.Fatalf("empty effective policy=%+v err=%v", policy, err)
	}
	delivery, err := st.CreateTestNotificationDelivery(channel.ID)
	if err != nil {
		t.Fatal(err)
	}
	claimed, err := st.ClaimNotificationDeliveries(0, time.Now().UTC().Add(time.Second))
	if err != nil || len(claimed) != 1 {
		t.Fatalf("claimed=%+v err=%v", claimed, err)
	}
	if err := st.CompleteNotificationDelivery(delivery.ID, true, 204, strings.Repeat("x", 600), time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	completed, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{ChannelID: channel.ID, State: domain.NotificationDeliverySucceeded, Limit: 2000})
	if err != nil || len(completed) != 1 || completed[0].LastStatus != 204 || len(completed[0].LastError) != 512 || completed[0].Endpoint != "" {
		t.Fatalf("completed=%+v err=%v", completed, err)
	}
	if err := st.RotateNotificationChannelAuthorization(channel.ID, ""); err != nil {
		t.Fatal(err)
	}
	got, err = st.GetNotificationChannel(channel.ID)
	if err != nil || got.HasAuthorization {
		t.Fatalf("cleared authorization channel=%+v err=%v", got, err)
	}
	missing := channel
	missing.ID = "missing"
	if err := st.UpdateNotificationChannel(missing); !errors.Is(err, ErrNotFound) {
		t.Fatalf("update missing channel error=%v", err)
	}
	if err := st.SetNotificationChannelEnabled(missing.ID, true); !errors.Is(err, ErrNotFound) {
		t.Fatalf("enable missing channel error=%v", err)
	}
	if err := st.RotateNotificationChannelAuthorization(missing.ID, "replacement"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("rotate missing channel error=%v", err)
	}
	if err := st.DeleteNotificationChannel(missing.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing channel error=%v", err)
	}
}

func TestDisabledChannelAndNonterminalRunCreateNoDelivery(t *testing.T) {
	st := openMem(t)
	group := domain.Group{Name: "scope", Enabled: true}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	task := createNotificationTestTask(t, st, group.ID)
	channel := createNotificationTestChannel(t, st, "disabled")
	assignment := []domain.NotificationAssignment{{ChannelID: channel.ID, OnSuccess: true}}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeGroup, group.ID, assignment); err != nil {
		t.Fatal(err)
	}
	configured, err := st.ListNotificationAssignments(domain.NotificationScopeGroup, group.ID)
	if err != nil || len(configured) != 1 {
		t.Fatalf("configured=%+v err=%v", configured, err)
	}
	if err := st.SetNotificationChannelEnabled(channel.ID, false); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	for _, outcome := range []domain.RunOutcome{domain.OutcomeSuccess, domain.OutcomeSkipped} {
		run := domain.Run{TaskID: task.ID, ScheduledFor: now, Outcome: outcome, Trigger: domain.TriggerManual}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID})
	if err != nil || len(deliveries) != 0 {
		t.Fatalf("deliveries=%+v err=%v", deliveries, err)
	}
	if err := st.SetNotificationChannelEnabled(channel.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeGroup, group.ID, []domain.NotificationAssignment{{ChannelID: channel.ID, OnFailure: true}}); err != nil {
		t.Fatal(err)
	}
	failure := domain.Run{TaskID: task.ID, ScheduledFor: now, Outcome: domain.OutcomeFailure, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&failure, ""); err != nil {
		t.Fatal(err)
	}
	deliveries, err = st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID, RunID: failure.ID, State: domain.NotificationDeliveryPending})
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("failure deliveries=%+v err=%v", deliveries, err)
	}
}

func createNotificationTestChannel(t *testing.T, st *Store, name string) domain.NotificationChannel {
	t.Helper()
	channel := domain.NotificationChannel{Name: name, Kind: domain.NotificationChannelWebhook, Endpoint: "https://example.test/hook?secret=value", EndpointSummary: "https://example.test", Authorization: "Bearer secret", Enabled: true}
	if err := st.CreateNotificationChannel(&channel); err != nil {
		t.Fatal(err)
	}
	return channel
}

func TestNotificationChannelAssignmentPrecedenceAndAtomicRunDelivery(t *testing.T) {
	st := openMem(t)
	parent := domain.Group{Name: "parent", Enabled: true}
	if err := st.CreateGroup(&parent); err != nil {
		t.Fatal(err)
	}
	child := domain.Group{Name: "child", ParentID: parent.ID, Enabled: true}
	if err := st.CreateGroup(&child); err != nil {
		t.Fatal(err)
	}
	task := createNotificationTestTask(t, st, child.ID)
	parentChannel := createNotificationTestChannel(t, st, "parent channel")
	childChannel := createNotificationTestChannel(t, st, "child channel")
	taskChannel := createNotificationTestChannel(t, st, "task channel")
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeGroup, parent.ID, []domain.NotificationAssignment{{ChannelID: parentChannel.ID, OnFailure: true}}); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeGroup, child.ID, []domain.NotificationAssignment{{ChannelID: childChannel.ID, OnFailure: true}}); err != nil {
		t.Fatal(err)
	}
	policy, err := st.EffectiveNotificationPolicy(task.ID)
	if err != nil || policy.SourceScopeID != child.ID || len(policy.Assignments) != 1 {
		t.Fatalf("child policy=%+v err=%v", policy, err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: taskChannel.ID, OnSuccess: true}}); err != nil {
		t.Fatal(err)
	}
	policy, err = st.EffectiveNotificationPolicy(task.ID)
	if err != nil || policy.SourceScopeType != domain.NotificationScopeTask || policy.Assignments[0].ChannelID != taskChannel.ID {
		t.Fatalf("task policy=%+v err=%v", policy, err)
	}

	now := time.Now().UTC()
	run := domain.Run{ID: "run-1", TaskID: task.ID, ScheduledFor: now, StartedAt: &now, EndedAt: &now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
		t.Fatal(err)
	}
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{RunID: run.ID})
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("deliveries=%+v err=%v", deliveries, err)
	}
	delivery := deliveries[0]
	if delivery.ChannelID != taskChannel.ID || delivery.State != domain.NotificationDeliveryPending || delivery.Authorization != "Bearer secret" {
		t.Fatalf("delivery=%+v", delivery)
	}
	var event domain.WebhookEvent
	if err := json.Unmarshal(delivery.Payload, &event); err != nil {
		t.Fatal(err)
	}
	if event.Task == nil || event.Task.Name != task.Name || event.Run == nil || event.Run.Outcome != domain.OutcomeSuccess || event.Delivery.ID != delivery.ID {
		t.Fatalf("event=%+v", event)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, nil); err != nil {
		t.Fatal(err)
	}
	policy, _ = st.EffectiveNotificationPolicy(task.ID)
	if policy.SourceScopeID != child.ID {
		t.Fatalf("inheritance did not resume: %+v", policy)
	}
}

func TestNotificationDeliveryRetryRecoveryCompletionAndChannelRemoval(t *testing.T) {
	st := openMem(t)
	channel := createNotificationTestChannel(t, st, "ops")
	delivery, err := st.CreateTestNotificationDelivery(channel.ID)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	claimed, err := st.ClaimNotificationDeliveries(1, now.Add(time.Second))
	if err != nil || len(claimed) != 1 || claimed[0].Attempts != 1 {
		t.Fatalf("claimed=%+v err=%v", claimed, err)
	}
	if err := st.RetryNotificationDelivery(delivery.ID, now, 503, "receiver returned HTTP 503"); err != nil {
		t.Fatal(err)
	}
	claimed, err = st.ClaimNotificationDeliveries(1, now.Add(time.Second))
	if err != nil || claimed[0].Attempts != 2 {
		t.Fatalf("second claim=%+v err=%v", claimed, err)
	}
	if recovered, err := st.RecoverNotificationDeliveries(now); err != nil || recovered != 1 {
		t.Fatalf("recovered=%d err=%v", recovered, err)
	}
	claimed, _ = st.ClaimNotificationDeliveries(1, now.Add(time.Second))
	if claimed[0].Attempts != 3 {
		t.Fatalf("attempts=%d", claimed[0].Attempts)
	}
	if err := st.CompleteNotificationDelivery(delivery.ID, false, 503, "receiver returned HTTP 503", now); err != nil {
		t.Fatal(err)
	}
	terminal, _ := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{ChannelID: channel.ID})
	if len(terminal) != 1 || terminal[0].Endpoint != "" || terminal[0].Authorization != "" || terminal[0].State != domain.NotificationDeliveryFailed {
		t.Fatalf("terminal=%+v", terminal)
	}
	if err := st.DeleteNotificationChannel(channel.ID); err != nil {
		t.Fatal(err)
	}
	terminal, _ = st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
	if len(terminal) != 1 || terminal[0].ChannelID != "" || terminal[0].ChannelName != "ops" {
		t.Fatalf("preserved=%+v", terminal)
	}
}
