package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/bundle"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestBundleApplyRequiresPreviewAndCreatesDisabledDraft(t *testing.T) {
	s := newTestServer(t)
	doc := bundle.Document{Schema: bundle.SchemaV1, Groups: []bundle.Group{{PortableID: "g-parent", Name: "Parent"}, {PortableID: "g-child", Name: "Child", ParentPortableID: "g-parent"}}, Tasks: []bundle.Task{{PortableID: "t-task", Name: "Portable task", GroupPortableID: "g-child", Timezone: "UTC", Schedule: "weekdays at 09:00", ScheduleSyntax: "human", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}}

	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	var plan bundle.Plan
	if err := json.Unmarshal(preview.Body.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if apply.Code != http.StatusOK {
		t.Fatalf("apply status=%d body=%s", apply.Code, apply.Body.String())
	}
	var result BundleApplyResponse
	if err := json.Unmarshal(apply.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	tasks, err := s.store.ListTasks("", "")
	if err != nil || len(tasks) != 1 {
		t.Fatalf("tasks=%+v err=%v outcomes=%+v", tasks, err, result.Outcomes)
	}
	if tasks[0].Enabled || tasks[0].Command != "" {
		t.Fatalf("imported task must remain a command-free disabled draft: %+v", tasks[0])
	}
	if _, err := s.store.ObjectIDForPortableID("group", "g-parent"); err != nil {
		t.Fatalf("parent identity: %v", err)
	}
	if _, err := s.store.ObjectIDForPortableID("group", "g-child"); err != nil {
		t.Fatalf("child identity: %v", err)
	}
	if second := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint}); second.Code != http.StatusConflict {
		t.Fatalf("reused preview status=%d body=%s", second.Code, second.Body.String())
	}
}

func TestBundleApplyRejectsTargetDrift(t *testing.T) {
	s := newTestServer(t)
	doc := bundle.Document{Schema: bundle.SchemaV1, Groups: []bundle.Group{{PortableID: "group", Name: "Portable"}}}
	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	var plan bundle.Plan
	if preview.Code != http.StatusOK || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	newGroupFor(t, s, "Changed target")
	apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if apply.Code != http.StatusConflict {
		t.Fatalf("drifted apply status=%d body=%s", apply.Code, apply.Body.String())
	}
}

func TestBundleCompareCannotAuthorizeApply(t *testing.T) {
	s := newTestServer(t)
	doc := bundle.Document{Schema: bundle.SchemaV1, Groups: []bundle.Group{{PortableID: "group", Name: "Portable"}}}
	compare := doJSON(t, s, http.MethodPost, "/v1/bundles/compare", BundleRequest{Bundle: doc})
	if compare.Code != http.StatusOK {
		t.Fatalf("compare status=%d body=%s", compare.Code, compare.Body.String())
	}
	var plan bundle.Plan
	if err := json.Unmarshal(compare.Body.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if apply.Code != http.StatusConflict {
		t.Fatalf("compare plan applied with status=%d body=%s", apply.Code, apply.Body.String())
	}
}

func TestBundleExportDoesNotExposeExecutionInputs(t *testing.T) {
	s := newTestServer(t)
	secret := "bundle-secret-canary"
	newTaskFor(t, s, TaskCreateRequest{Name: "safe export", Command: secret, WorkingDir: secret, Env: map[string]string{"TOKEN": secret}, Stdin: secret})
	export := doJSON(t, s, http.MethodGet, "/v1/bundles/export", nil)
	if export.Code != http.StatusOK {
		t.Fatalf("export status=%d body=%s", export.Code, export.Body.String())
	}
	if got := export.Body.String(); strings.Contains(got, secret) {
		t.Fatalf("bundle leaked execution input: %s", got)
	}
}

func TestBundleRequestRejectsUnknownFields(t *testing.T) {
	s := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/bundles/validate", bytes.NewBufferString(`{"bundle":{"schema":"go-schedule.bundle/v1","command":"unsafe"}}`))
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestBundleUpdatePreservesRequestedEnabledStateWhenLocallyReady(t *testing.T) {
	s := newTestServer(t)
	doc := bundle.Document{Schema: bundle.SchemaV1, Tasks: []bundle.Task{{PortableID: "portable-task", Name: "Source name", Enabled: true, Timezone: "UTC", Schedule: "every day at 09:00", ScheduleSyntax: "human", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}}
	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	var plan bundle.Plan
	if preview.Code != http.StatusOK || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("initial preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	if apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint}); apply.Code != http.StatusOK {
		t.Fatalf("initial apply status=%d body=%s", apply.Code, apply.Body.String())
	}
	tasks, err := s.store.ListTasks("", "")
	if err != nil || len(tasks) != 1 {
		t.Fatalf("tasks=%+v err=%v", tasks, err)
	}
	task := tasks[0]
	task.Command, task.Enabled = "echo", true
	if err := s.store.UpdateTask(&task); err != nil {
		t.Fatal(err)
	}
	doc.Tasks[0].Name = "Updated name"
	preview = doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	if preview.Code != http.StatusOK || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("update preview status=%d body=%s", preview.Code, preview.Body.String())
	}
	if apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint}); apply.Code != http.StatusOK {
		t.Fatalf("update apply status=%d body=%s", apply.Code, apply.Body.String())
	}
	updated, err := s.store.GetTask(task.ID)
	if err != nil || !updated.Enabled {
		t.Fatalf("updated task=%+v err=%v", updated, err)
	}
}

func TestBundleV2SourcesRequireBindingsAndApplyAsDisabled(t *testing.T) {
	s := newTestServer(t)
	doc := bundle.Document{Schema: bundle.SchemaV2, Tasks: []bundle.Task{{PortableID: "task", Name: "Portable", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}, ExternalTriggers: []bundle.ExternalTrigger{{PortableID: "trigger", Name: " API ", TargetTaskID: "task"}}, TriggerSets: []bundle.TriggerSet{{PortableID: "set", Name: " Batch ", TargetTaskID: "task", MemberCount: 2}}, Watchers: []bundle.Watcher{{PortableID: "watcher", Name: " Files ", Kind: "directory", Pattern: "*.txt", Debounce: "0.25s", Stability: "0.5s", TargetTaskID: "task"}}, NotificationPolicies: []bundle.NotificationPolicy{{ScopeType: "task", ScopePortableID: "task", Assignments: []bundle.NotificationAssignment{{ChannelName: "ops", OnFailure: true, FailureThreshold: 1}}}}}
	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview=%d %s", preview.Code, preview.Body.String())
	}
	var plan bundle.Plan
	if err := json.Unmarshal(preview.Body.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	for _, item := range plan.Items {
		if item.Kind == "watcher" && item.Action != bundle.ActionConflict {
			t.Fatalf("unbound watcher=%+v", item)
		}
		if item.Kind == "notification_policy" && item.Action != bundle.ActionConflict {
			t.Fatalf("unbound policy=%+v", item)
		}
	}
	channel := domain.NotificationChannel{Name: "ops", Kind: domain.NotificationChannelWebhook, Endpoint: "https://example.test/secret", Authorization: "Bearer private", Enabled: true}
	if err := s.store.CreateNotificationChannel(&channel); err != nil {
		t.Fatal(err)
	}
	path := t.TempDir()
	foreignPath := `C:\source\watcher`
	if runtime.GOOS == "windows" {
		foreignPath = "/source/watcher"
	}
	foreign := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc, WatcherPaths: map[string]string{"watcher": foreignPath}})
	if foreign.Code != http.StatusOK || json.Unmarshal(foreign.Body.Bytes(), &plan) != nil {
		t.Fatalf("foreign preview=%d %s", foreign.Code, foreign.Body.String())
	}
	for _, item := range plan.Items {
		if item.Kind == "watcher" && item.Action != bundle.ActionConflict {
			t.Fatalf("foreign path accepted on %s: %+v", runtime.GOOS, item)
		}
	}
	preview = doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc, WatcherPaths: map[string]string{"watcher": path}})
	if preview.Code != http.StatusOK || json.Unmarshal(preview.Body.Bytes(), &plan) != nil {
		t.Fatalf("bound preview=%d %s", preview.Code, preview.Body.String())
	}
	for _, item := range plan.Items {
		if item.Action == bundle.ActionConflict || item.Action == bundle.ActionInvalid {
			t.Fatalf("unexpected conflict: %+v", item)
		}
	}
	apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if apply.Code != http.StatusOK {
		t.Fatalf("apply=%d %s", apply.Code, apply.Body.String())
	}
	var result BundleApplyResponse
	if err := json.Unmarshal(apply.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	for _, item := range result.Outcomes {
		if item.Action != bundle.ActionApplied {
			t.Fatalf("item not applied: %+v", item)
		}
	}
	triggers, err := s.store.ListExternalTriggers()
	if err != nil || len(triggers) != 3 {
		t.Fatalf("triggers=%+v err=%v", triggers, err)
	}
	for _, trigger := range triggers {
		if trigger.Enabled || !strings.HasPrefix(trigger.Key, "gst_") {
			t.Fatalf("unsafe imported trigger: %+v", trigger)
		}
	}
	watchers, err := s.store.ListFilesystemWatchers()
	if err != nil || len(watchers) != 1 || watchers[0].Enabled || watchers[0].Path != path {
		t.Fatalf("watchers=%+v err=%v", watchers, err)
	}
	export := doJSON(t, s, http.MethodGet, "/v1/bundles/export", nil)
	if export.Code != http.StatusOK {
		t.Fatalf("export=%d %s", export.Code, export.Body.String())
	}
	for _, forbidden := range []string{path, "https://example.test/secret", "Bearer private", triggers[0].Key} {
		if strings.Contains(export.Body.String(), forbidden) {
			t.Fatalf("export leaked %q", forbidden)
		}
	}
	if !strings.Contains(export.Body.String(), bundle.SchemaV2) {
		t.Fatal("export did not emit v2")
	}
	second := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	if second.Code != http.StatusOK || json.Unmarshal(second.Body.Bytes(), &plan) != nil {
		t.Fatalf("second preview=%d %s", second.Code, second.Body.String())
	}
	for _, item := range plan.Items {
		if item.Action != bundle.ActionUnchanged {
			t.Fatalf("round trip drift: %+v", item)
		}
	}
}

func TestBundlePolicyRejectsChannelReplacementAfterPreview(t *testing.T) {
	s := newTestServer(t)
	doc := bundle.Document{Schema: bundle.SchemaV2, Tasks: []bundle.Task{{PortableID: "task", Name: "Portable", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}, NotificationPolicies: []bundle.NotificationPolicy{{ScopeType: "task", ScopePortableID: "task", Assignments: []bundle.NotificationAssignment{{ChannelName: "ops", OnFailure: true, FailureThreshold: 1}}}}}
	channel := domain.NotificationChannel{Name: "ops", Kind: domain.NotificationChannelWebhook, Endpoint: "https://example.test/first", Enabled: true}
	if err := s.store.CreateNotificationChannel(&channel); err != nil {
		t.Fatal(err)
	}
	preview := doJSON(t, s, http.MethodPost, "/v1/bundles/preview", BundleRequest{Bundle: doc})
	if preview.Code != http.StatusOK {
		t.Fatalf("preview=%d %s", preview.Code, preview.Body.String())
	}
	var plan bundle.Plan
	if err := json.Unmarshal(preview.Body.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	channel.Name = "renamed"
	if err := s.store.UpdateNotificationChannel(channel); err != nil {
		t.Fatal(err)
	}
	replacement := domain.NotificationChannel{Name: "ops", Kind: domain.NotificationChannelWebhook, Endpoint: "https://example.test/replacement", Enabled: true}
	if err := s.store.CreateNotificationChannel(&replacement); err != nil {
		t.Fatal(err)
	}
	apply := doJSON(t, s, http.MethodPost, "/v1/bundles/apply", BundleApplyRequest{PlanID: plan.ID, BundleDigest: plan.BundleDigest, TargetDaemonID: plan.TargetDaemonID, TargetFingerprint: plan.TargetFingerprint})
	if apply.Code != http.StatusOK {
		t.Fatalf("apply=%d %s", apply.Code, apply.Body.String())
	}
	var result BundleApplyResponse
	if err := json.Unmarshal(apply.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	policyFailed := false
	for _, item := range result.Outcomes {
		if item.Kind == "notification_policy" && item.Action == bundle.ActionFailed && strings.Contains(item.Message, "preview again") {
			policyFailed = true
		}
	}
	if !policyFailed {
		t.Fatalf("changed channel was accepted: %+v", result.Outcomes)
	}
	taskID, err := s.store.ObjectIDForPortableID("task", "task")
	if err != nil {
		t.Fatal(err)
	}
	assignments, err := s.store.ListNotificationAssignments(domain.NotificationScopeTask, taskID)
	if err != nil || len(assignments) != 0 {
		t.Fatalf("unexpected policy assignments=%+v err=%v", assignments, err)
	}
}
