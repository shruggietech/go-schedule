package automation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type fakeBackend struct {
	Backend
	tasks        []server.TaskResponse
	chains       []domain.CompletionChain
	triggers     []server.TriggerResponse
	sets         []server.TriggerSetResponse
	watchers     []server.FilesystemWatcherResponse
	firedKey     string
	updatedChain bool
	fail         error
}

func (f *fakeBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	return f.tasks, f.fail
}
func (f *fakeBackend) ListChains(context.Context) ([]domain.CompletionChain, error) {
	return f.chains, nil
}
func (f *fakeBackend) ListTriggers(context.Context) ([]server.TriggerResponse, error) {
	return f.triggers, nil
}
func (f *fakeBackend) ListTriggerSets(context.Context) ([]server.TriggerSetResponse, error) {
	return f.sets, nil
}
func (f *fakeBackend) ListFilesystemWatchers(context.Context) ([]server.FilesystemWatcherResponse, error) {
	return f.watchers, nil
}
func (f *fakeBackend) GetChain(context.Context, string) (domain.CompletionChain, error) {
	return f.chains[0], nil
}
func (f *fakeBackend) UpdateChain(context.Context, string, server.ChainUpdateRequest) (domain.CompletionChain, error) {
	f.updatedChain = true
	return f.chains[0], nil
}
func (f *fakeBackend) RevealTrigger(context.Context, string) (server.TriggerSecretResponse, error) {
	return server.TriggerSecretResponse{Trigger: f.triggers[0], Key: "raw-secret", Command: "gosched trigger fire raw-secret"}, nil
}
func (f *fakeBackend) FireTrigger(_ context.Context, key string) error { f.firedKey = key; return nil }

func TestWorkspaceIsCompleteSecretFreeAndHonest(t *testing.T) {
	now := time.Now().UTC()
	fake := &fakeBackend{
		tasks:    []server.TaskResponse{{Task: domain.Task{ID: "task-1", Name: "Build"}}},
		chains:   []domain.CompletionChain{{ID: "chain-1", SourceTaskID: "gone", TargetTaskID: "task-1", TargetTaskName: "Build", OnOutcome: domain.CompletionOnSuccess, UpdatedAt: now}},
		triggers: []server.TriggerResponse{{ID: "trigger-1", Name: "Deploy", TargetTaskID: "gone", Readiness: "target_missing", Reason: "Target task is missing.", UpdatedAt: now}, {ID: "member-1", Name: "Member", TargetTaskID: "task-1", TargetTaskName: "Build", SetID: "set-1", SetName: "Fleet", SetPosition: 1, Enabled: true, Readiness: "ready", UpdatedAt: now}},
		sets:     []server.TriggerSetResponse{{ID: "set-1", Name: "Fleet", TargetTaskID: "task-1", TargetTaskName: "Build", MemberCount: 2, EnabledCount: 1, UpdatedAt: now, Members: []server.TriggerResponse{{ID: "member-1", Name: "Member", TargetTaskID: "task-1", TargetTaskName: "Build", SetPosition: 1, Enabled: true, Readiness: "ready", UpdatedAt: now}, {ID: "member-2", Name: "Paused", TargetTaskID: "task-1", TargetTaskName: "Build", SetPosition: 2, Enabled: false, Readiness: "disabled", Reason: "Trigger is disabled.", UpdatedAt: now}}}},
		watchers: []server.FilesystemWatcherResponse{{ID: "watcher-1", Name: "Inbox", Kind: domain.WatcherDirectory, Path: `C:\very\long\inbox`, TargetTaskID: "task-1", TargetTaskName: "Build", Enabled: true, Health: domain.WatcherHealth{State: domain.WatcherDegraded, Reason: "permission denied"}, Readiness: "ready", UpdatedAt: now}},
	}
	result := NewService(fake).Workspace(context.Background())
	if result.Outcome != "accepted" || result.Workspace == nil {
		t.Fatalf("result=%+v", result)
	}
	if len(result.Workspace.Triggers) != 2 || result.Workspace.Triggers[1].SetID != "set-1" {
		t.Fatalf("triggers=%+v", result.Workspace.Triggers)
	}
	if result.Workspace.Chains[0].Readiness != "source_missing" || result.Workspace.TriggerSets[0].Readiness != "ready" || result.Workspace.Watchers[0].Health != "degraded" || result.Workspace.Watchers[0].HealthReason != "permission denied" {
		t.Fatalf("workspace=%+v", result.Workspace)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "raw-secret") || strings.Contains(string(encoded), "gosched trigger fire") {
		t.Fatalf("ordinary workspace leaked a secret: %s", encoded)
	}
}

func TestWorkspaceRejectsPartialSnapshots(t *testing.T) {
	result := NewService(&fakeBackend{fail: errors.New("offline")}).Workspace(context.Background())
	if result.Outcome != "unavailable" || result.Workspace != nil {
		t.Fatalf("result=%+v", result)
	}
}

func TestStaleChainRequiresExplicitOverwrite(t *testing.T) {
	now := time.Now().UTC()
	fake := &fakeBackend{chains: []domain.CompletionChain{{ID: "chain-1", SourceTaskID: "a", TargetTaskID: "b", OnOutcome: domain.CompletionOnSuccess, UpdatedAt: now}}}
	service := NewService(fake)
	draft := ChainDraft{ID: "chain-1", SourceTaskID: "a", TargetTaskID: "b", OnOutcome: "success", OriginalUpdatedAt: now.Add(-time.Minute).Format(time.RFC3339Nano)}
	if result := service.SaveChain(context.Background(), draft); result.Outcome != "stale" || fake.updatedChain {
		t.Fatalf("result=%+v updated=%t", result, fake.updatedChain)
	}
	draft.OverwriteStale = true
	result := service.SaveChain(context.Background(), draft)
	if !fake.updatedChain || result.EntityID != "chain-1" {
		t.Fatalf("result=%+v updated=%t", result, fake.updatedChain)
	}
}

func TestFireTriggerKeepsSecretInsideService(t *testing.T) {
	fake := &fakeBackend{triggers: []server.TriggerResponse{{ID: "trigger-1", Name: "Deploy"}}}
	result := NewService(fake).FireTrigger(context.Background(), "trigger-1")
	if fake.firedKey != "raw-secret" || result.Outcome != "accepted" {
		t.Fatalf("result=%+v key=%q", result, fake.firedKey)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "raw-secret") {
		t.Fatalf("fire result leaked key: %s", encoded)
	}
}

func TestWatcherDurationBoundsAreValidatedBeforeSubmission(t *testing.T) {
	service := NewService(&fakeBackend{})
	base := WatcherDraft{Name: "Inbox", Kind: "directory", Path: "inbox", TargetTaskID: "task-1", Debounce: "25ms", Stability: "1h", IsNew: true}
	for _, test := range []struct{ field, value string }{{"debounce", "0s"}, {"debounce", "-1s"}, {"debounce", "2h"}, {"stability", "24ms"}, {"stability", "invalid"}} {
		draft := base
		if test.field == "debounce" {
			draft.Debounce = test.value
		} else {
			draft.Stability = test.value
		}
		result := service.SaveWatcher(context.Background(), draft)
		if result.Outcome != "rejected" || result.Field != test.field {
			t.Fatalf("%s=%q result=%+v", test.field, test.value, result)
		}
	}
}
