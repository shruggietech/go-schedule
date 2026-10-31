package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func isolatedCommandBackend(t *testing.T) *client.Client {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Error(err)
		}
	})
	backend := server.New(st, nil, nil, nil, "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	return client.NewInProcess(backend.Handler())
}

func runFixtureCommand(c *client.Client, input string, args ...string) (string, error) {
	var out, stderr bytes.Buffer
	root := newRootWithDependencies(Dependencies{Client: c, In: strings.NewReader(input), Out: &out, Err: &stderr})
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func TestStructuredNotificationPolicyRoundTripAndInheritance(t *testing.T) {
	c := isolatedCommandBackend(t)
	ctx := context.Background()
	group, err := c.CreateGroup(ctx, server.GroupCreateRequest{Name: "fixture"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := c.CreateTask(ctx, server.TaskCreateRequest{Name: "fixture", GroupID: group.ID})
	if err != nil {
		t.Fatal(err)
	}
	channels := make([]domain.NotificationChannel, 2)
	for i, name := range []string{"success", "problems"} {
		channels[i], err = c.CreateNotificationChannel(ctx, server.NotificationChannelCreateRequest{Name: name, Endpoint: "https://example.test/fixture"})
		if err != nil {
			t.Fatal(err)
		}
	}
	mixed := server.NotificationAssignmentsRequest{Assignments: []server.NotificationAssignmentInput{
		{ChannelID: channels[0].ID, OnSuccess: true},
		{ChannelID: channels[1].ID, OnFailure: true, FailureThreshold: 3, OnFailureToStart: true, DurationThresholdSeconds: 600, OnRecovery: true, ReminderIntervalSeconds: 3600, QuietPeriodSeconds: 600},
	}}
	input, err := json.Marshal(mixed)
	if err != nil {
		t.Fatal(err)
	}
	for _, scope := range []struct{ kind, id string }{{"group", group.ID}, {"task", task.Task.ID}} {
		if _, err := runFixtureCommand(c, string(input), "notification", scope.kind, "set", scope.id, "--file", "-", "--json"); err != nil {
			t.Fatal(err)
		}
		exported, err := runFixtureCommand(c, "", "notification", scope.kind, "show", scope.id, "--input")
		if err != nil {
			t.Fatal(err)
		}
		var canonical server.NotificationAssignmentsRequest
		if err := json.Unmarshal([]byte(exported), &canonical); err != nil {
			t.Fatal(err)
		}
		byID := map[string]server.NotificationAssignmentInput{}
		for _, item := range canonical.Assignments {
			byID[item.ChannelID] = item
		}
		if len(canonical.Assignments) != 2 || !byID[channels[0].ID].OnSuccess || byID[channels[1].ID].FailureThreshold != 3 || !byID[channels[1].ID].OnRecovery || byID[channels[1].ID].QuietPeriodSeconds != 600 {
			t.Fatalf("mixed policy lost: %s", exported)
		}
		filename := filepath.Join(t.TempDir(), "policy.json")
		if err := os.WriteFile(filename, []byte(exported), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := runFixtureCommand(c, "", "notification", scope.kind, "set", scope.id, "--file", filename); err != nil {
			t.Fatal(err)
		}
	}
	before, err := c.ListTaskNotificationAssignments(ctx, task.Task.ID)
	if err != nil {
		t.Fatal(err)
	}
	duplicate := `{"assignments":[{"channel_id":"` + channels[0].ID + `","on_success":true},{"channel_id":"` + channels[0].ID + `","on_failure":true}]}`
	if _, err := runFixtureCommand(c, duplicate, "notification", "task", "set", task.Task.ID, "--file", "-"); err == nil {
		t.Fatal("duplicate replacement accepted")
	}
	after, err := c.ListTaskNotificationAssignments(ctx, task.Task.ID)
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatalf("invalid replacement changed policy: %v", err)
	}
	if _, err := runFixtureCommand(c, `{"assignments":[]}`, "notification", "task", "set", task.Task.ID, "--file", "-"); err != nil {
		t.Fatal(err)
	}
	policy, err := c.EffectiveTaskNotificationPolicy(ctx, task.Task.ID)
	if err != nil || policy.SourceScopeType != domain.NotificationScopeGroup || policy.SourceScopeID != group.ID || len(policy.Assignments) != 2 {
		t.Fatalf("inheritance not restored: %+v err=%v", policy, err)
	}
}

func TestNotificationStructuredInputRejectsMalformedAndConflictingForms(t *testing.T) {
	for _, input := range []string{`{}`, `null`, `{"assignments":null}`, `{"assignments":[] } {}`, `{"assignments":[],"extra":true}`, `{"assignments":[{"channel_id":"x","extra":true}]}`, string([]byte{0xff}), strings.Repeat(" ", 1024*1024+1)} {
		_, _, calls, err := executeAuthoring(t, []string{"notification", "task", "set", "fixture", "--file", "-"}, input)
		if err == nil || calls != 0 {
			t.Fatalf("invalid input accepted/requested: bytes=%d calls=%d err=%v", len(input), calls, err)
		}
	}
	for _, flag := range []string{"--channel=x", "--on=failure", "--failure-threshold=1", "--failure-to-start=false", "--duration-threshold=0s", "--recovery=false", "--reminder=0s", "--quiet-period=0s"} {
		_, _, calls, err := executeAuthoring(t, []string{"notification", "group", "set", "fixture", "--file", "-", flag}, `{"assignments":[]}`)
		if err == nil || calls != 0 {
			t.Fatalf("conflict %s accepted: calls=%d err=%v", flag, calls, err)
		}
	}
}
