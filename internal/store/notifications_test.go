package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
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

func TestPublishedWebhookSchemaAcceptsConditionAndHealthPayloads(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "specs", "067-webhook-notifications", "contracts", "webhook-v1.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var document any
	if err := json.NewDecoder(file).Decode(&document); err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("webhook-v1.schema.json", document); err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile("webhook-v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	created := "2026-09-21T12:00:00Z"
	samples := []string{
		`{"schema":"go-schedule.webhook.v1","event":"run.completed","delivery":{"id":"delivery-1","created_at":"` + created + `"},"daemon":{"version":"v1.5.0"},"task":{"id":"task-1","name":"Backup"},"run":{"id":"run-1","outcome":"failure","trigger":"manual","scheduled_for":"` + created + `"},"condition":{"kind":"failure_to_start","summary":"Task process could not be started.","streak":1,"threshold":1}}`,
		`{"schema":"go-schedule.webhook.v1","event":"daemon.health","delivery":{"id":"delivery-2","created_at":"` + created + `"},"daemon":{"version":"v1.5.0","status":"healthy","next_expected_at":"2026-09-21T12:05:00Z"},"task":null,"condition":{"kind":"daemon_health","summary":"Daemon is healthy."}}`,
	}
	for _, sample := range samples {
		var value any
		if err := json.Unmarshal([]byte(sample), &value); err != nil {
			t.Fatal(err)
		}
		if err := schema.Validate(value); err != nil {
			t.Fatalf("published schema rejected payload %s: %v", sample, err)
		}
	}
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

func TestNotificationConditionsThresholdReminderRecoveryAndPrecedence(t *testing.T) {
	st := openMem(t)
	task := createNotificationTestTask(t, st, "")
	channel := createNotificationTestChannel(t, st, "conditions")
	assignment := domain.NotificationAssignment{ChannelID: channel.ID, OnFailure: true, FailureThreshold: 2, OnFailureToStart: true, DurationThresholdSeconds: 60, OnRecovery: true, ReminderIntervalSeconds: 120, QuietPeriodSeconds: 180}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{assignment}); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	record := func(id string, outcome domain.RunOutcome, startFailed bool, ended time.Time, duration time.Duration) {
		t.Helper()
		started := ended.Add(-duration)
		run := domain.Run{ID: id, TaskID: task.ID, ScheduledFor: started, StartedAt: &started, EndedAt: &ended, Outcome: outcome, Trigger: domain.TriggerManual, StartFailed: startFailed}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	record("failure-1", domain.OutcomeFailure, false, base, time.Second)
	record("failure-2", domain.OutcomeFailure, false, base.Add(time.Minute), time.Second)
	record("failure-3", domain.OutcomeFailure, false, base.Add(2*time.Minute), time.Second)
	record("failure-4", domain.OutcomeFailure, false, base.Add(3*time.Minute), time.Second)
	record("failure-5", domain.OutcomeFailure, false, base.Add(4*time.Minute), time.Second)
	record("start-failed", domain.OutcomeFailure, true, base.Add(5*time.Minute), time.Second)
	record("recovered", domain.OutcomeSuccess, false, base.Add(6*time.Minute), time.Second)
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(deliveries) != 4 {
		t.Fatalf("deliveries=%d, want threshold, reminder, start precedence, recovery: %+v", len(deliveries), deliveries)
	}
	want := []domain.NotificationConditionKind{domain.NotificationConditionRecovery, domain.NotificationConditionFailureToStart, domain.NotificationConditionConsecutiveFailure, domain.NotificationConditionConsecutiveFailure}
	for i, delivery := range deliveries {
		if delivery.ConditionKind != want[i] {
			t.Fatalf("delivery %d condition=%s, want %s", i, delivery.ConditionKind, want[i])
		}
		var event domain.WebhookEvent
		if err := json.Unmarshal(delivery.Payload, &event); err != nil || event.Condition == nil || event.Condition.Kind != want[i] {
			t.Fatalf("delivery %d event=%+v err=%v", i, event, err)
		}
	}
}

func TestNotificationSuccessContinuesWhenRecoveryDeliveryIsDisabled(t *testing.T) {
	st := openMem(t)
	task := createNotificationTestTask(t, st, "")
	channel := createNotificationTestChannel(t, st, "success without recovery")
	assignment := domain.NotificationAssignment{ChannelID: channel.ID, OnSuccess: true, OnFailure: true}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{assignment}); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	for index, outcome := range []domain.RunOutcome{domain.OutcomeFailure, domain.OutcomeSuccess} {
		ended := base.Add(time.Duration(index) * time.Minute)
		run := domain.Run{ID: fmt.Sprintf("transition-%d", index), TaskID: task.ID, ScheduledFor: ended, EndedAt: &ended, Outcome: outcome, Trigger: domain.TriggerManual}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID})
	if err != nil || len(deliveries) != 2 || deliveries[0].ConditionKind != domain.NotificationConditionSuccess || deliveries[1].ConditionKind != domain.NotificationConditionFailure {
		t.Fatalf("deliveries=%+v err=%v", deliveries, err)
	}
}

func TestNotificationPolicyOverrideRoundTripResetsInheritedState(t *testing.T) {
	st := openMem(t)
	group := domain.Group{Name: "Inherited policy", Enabled: true}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	task := createNotificationTestTask(t, st, group.ID)
	inherited := createNotificationTestChannel(t, st, "inherited")
	direct := createNotificationTestChannel(t, st, "direct")
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeGroup, group.ID, []domain.NotificationAssignment{{ChannelID: inherited.ID, OnFailure: true, FailureThreshold: 2}}); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	recordFailure := func(id string, at time.Time) {
		t.Helper()
		run := domain.Run{ID: id, TaskID: task.ID, ScheduledFor: at, EndedAt: &at, Outcome: domain.OutcomeFailure, Trigger: domain.TriggerManual}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	recordFailure("inherited-first", base)
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: direct.ID, OnFailure: true}}); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, nil); err != nil {
		t.Fatal(err)
	}
	recordFailure("inherited-after-override", base.Add(time.Minute))
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID})
	if err != nil || len(deliveries) != 0 {
		t.Fatalf("stale inherited state created delivery: %+v err=%v", deliveries, err)
	}
}

func TestNotificationSuccessQuietPeriodAndDaemonHeartbeat(t *testing.T) {
	st := openMem(t)
	task := createNotificationTestTask(t, st, "")
	channel := createNotificationTestChannel(t, st, "quiet and health")
	channel.HealthIntervalSeconds = 60
	if err := st.UpdateNotificationChannel(channel); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: channel.ID, OnSuccess: true, QuietPeriodSeconds: 120}}); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	for i, offset := range []time.Duration{0, time.Minute, 2 * time.Minute} {
		ended := base.Add(offset)
		run := domain.Run{ID: fmt.Sprintf("success-%d", i), TaskID: task.ID, ScheduledFor: ended, EndedAt: &ended, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
		if err := st.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID})
	if err != nil || len(deliveries) != 2 {
		t.Fatalf("quiet deliveries=%+v err=%v", deliveries, err)
	}
	created, err := st.CreateDueDaemonHealthDeliveries(base)
	if err != nil || created != 1 {
		t.Fatalf("created heartbeat=%d err=%v", created, err)
	}
	created, err = st.CreateDueDaemonHealthDeliveries(base.Add(30 * time.Second))
	if err != nil || created != 0 {
		t.Fatalf("early heartbeat=%d err=%v", created, err)
	}
	heartbeats, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{ChannelID: channel.ID, Limit: 20})
	if err != nil || len(heartbeats) != 3 {
		t.Fatalf("heartbeat deliveries=%+v err=%v", heartbeats, err)
	}
	foundHeartbeat := false
	for _, delivery := range heartbeats {
		foundHeartbeat = foundHeartbeat || (delivery.EventKind == domain.NotificationEventDaemonHealth && delivery.ConditionKind == domain.NotificationConditionDaemonHealth)
	}
	if !foundHeartbeat {
		t.Fatalf("daemon heartbeat not found: %+v", heartbeats)
	}
}

func TestNotificationConditionStateSurvivesRestartAndPolicyChangeResetsIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notifications.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	task := createNotificationTestTask(t, st, "")
	channel := createNotificationTestChannel(t, st, "restart")
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: channel.ID, OnFailure: true, FailureThreshold: 2}}); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	recordFailure := func(store *Store, id string, at time.Time) {
		t.Helper()
		run := domain.Run{ID: id, TaskID: task.ID, ScheduledFor: at, EndedAt: &at, Outcome: domain.OutcomeFailure, Trigger: domain.TriggerManual}
		if err := store.RecordRunAndCreateDeliveries(&run, ""); err != nil {
			t.Fatal(err)
		}
	}
	recordFailure(st, "before-restart", base)
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	recordFailure(st, "after-restart", base.Add(time.Minute))
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID})
	if err != nil || len(deliveries) != 1 || deliveries[0].ConditionKind != domain.NotificationConditionConsecutiveFailure {
		t.Fatalf("restart deliveries=%+v err=%v", deliveries, err)
	}
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeTask, task.ID, []domain.NotificationAssignment{{ChannelID: channel.ID, OnFailure: true, FailureThreshold: 2}}); err != nil {
		t.Fatal(err)
	}
	recordFailure(st, "after-policy-change", base.Add(2*time.Minute))
	deliveries, err = st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{TaskID: task.ID})
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("policy reset created delivery: %+v err=%v", deliveries, err)
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

func TestNotificationRecoveryExhaustsThirdClaimAndErasesSecrets(t *testing.T) {
	st := openMem(t)
	channel := createNotificationTestChannel(t, st, "recovery exhaustion")
	delivery, err := st.CreateTestNotificationDelivery(channel.ID)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	for attempt := 1; attempt <= 3; attempt++ {
		claimed, err := st.ClaimNotificationDeliveries(1, now.Add(time.Second))
		if err != nil || len(claimed) != 1 || claimed[0].Attempts != attempt {
			t.Fatalf("attempt %d claim=%+v err=%v", attempt, claimed, err)
		}
		if attempt < 3 {
			if err := st.RetryNotificationDelivery(delivery.ID, now, 503, "retry"); err != nil {
				t.Fatal(err)
			}
		}
	}
	if recovered, err := st.RecoverNotificationDeliveries(now); err != nil || recovered != 1 {
		t.Fatalf("recovered=%d err=%v", recovered, err)
	}
	failed, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{State: domain.NotificationDeliveryFailed})
	if err != nil || len(failed) != 1 || failed[0].Endpoint != "" || failed[0].Authorization != "" {
		t.Fatalf("failed=%+v err=%v", failed, err)
	}
}

func TestNotificationOperationsReportClosedStore(t *testing.T) {
	st := openMem(t)
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	checks := []struct {
		name string
		run  func() error
	}{
		{"delete group", func() error { return st.DeleteGroup("missing") }},
		{"delete task", func() error { return st.DeleteTask("missing") }},
		{"recover deliveries", func() error { _, err := st.RecoverNotificationDeliveries(now); return err }},
		{"claim deliveries", func() error { _, err := st.ClaimNotificationDeliveries(1, now); return err }},
		{"complete delivery", func() error { return st.CompleteNotificationDelivery("missing", false, 0, "", now) }},
		{"list deliveries", func() error { _, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{}); return err }},
	}
	for _, check := range checks {
		if err := check.run(); err == nil {
			t.Fatalf("%s succeeded on closed store", check.name)
		}
	}
}

func TestDeletingNotificationSourcesRemovesOnlyUnfinishedWork(t *testing.T) {
	st := openMem(t)
	group := domain.Group{Name: "source group", Enabled: true}
	if err := st.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	task := createNotificationTestTask(t, st, group.ID)
	channel := createNotificationTestChannel(t, st, "source lifecycle")
	if err := st.ReplaceNotificationAssignments(domain.NotificationScopeGroup, group.ID, []domain.NotificationAssignment{{ChannelID: channel.ID, OnSuccess: true}}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	terminalRun := domain.Run{TaskID: task.ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&terminalRun, ""); err != nil {
		t.Fatal(err)
	}
	claimed, err := st.ClaimNotificationDeliveries(1, now.Add(time.Second))
	if err != nil || len(claimed) != 1 {
		t.Fatalf("claimed=%+v err=%v", claimed, err)
	}
	if err := st.CompleteNotificationDelivery(claimed[0].ID, true, 204, "", now); err != nil {
		t.Fatal(err)
	}
	pendingRun := domain.Run{TaskID: task.ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&pendingRun, ""); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteTask(task.ID); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteTask(task.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing source task error=%v", err)
	}
	deliveries, err := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{})
	if err != nil || len(deliveries) != 1 || deliveries[0].State != domain.NotificationDeliverySucceeded {
		t.Fatalf("deliveries after task deletion=%+v err=%v", deliveries, err)
	}

	groupTask := createNotificationTestTask(t, st, group.ID)
	pendingGroupRun := domain.Run{TaskID: groupTask.ID, ScheduledFor: now, Outcome: domain.OutcomeSuccess, Trigger: domain.TriggerManual}
	if err := st.RecordRunAndCreateDeliveries(&pendingGroupRun, ""); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteGroup(group.ID); err != nil {
		t.Fatal(err)
	}
	deliveries, err = st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{RunID: pendingGroupRun.ID})
	if err != nil || len(deliveries) != 0 {
		t.Fatalf("deliveries after group deletion=%+v err=%v", deliveries, err)
	}
	if err := st.DeleteGroup(group.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing source group error=%v", err)
	}
}
