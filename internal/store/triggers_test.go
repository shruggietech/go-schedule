package store

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestExternalTriggerLifecycleAndSecretRedaction(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	trigger := &domain.ExternalTrigger{Name: "hook", TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(trigger.Key, "gst_") || len(trigger.Key) < 40 {
		t.Fatalf("generated key has unexpected shape")
	}
	encoded, err := json.Marshal(trigger)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), trigger.Key) || strings.Contains(string(encoded), `"key"`) {
		t.Fatalf("ordinary JSON disclosed trigger key: %s", encoded)
	}
	loaded, err := st.GetExternalTriggerByKey(trigger.Key)
	if err != nil || loaded.ID != trigger.ID || loaded.TargetTaskName != "target" {
		t.Fatalf("loaded trigger = %+v, err=%v", loaded, err)
	}
	oldKey := trigger.Key
	newKey, err := st.RotateExternalTrigger(trigger.ID)
	if err != nil || newKey == oldKey {
		t.Fatalf("rotate key=%q err=%v", newKey, err)
	}
	if _, err := st.GetExternalTriggerByKey(oldKey); err != ErrNotFound {
		t.Fatalf("old key error = %v, want ErrNotFound", err)
	}
}

func TestExternalTriggerAdministrationAndRetargeting(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	first := &domain.Task{Name: "first", Command: "echo", State: domain.TaskActive}
	second := &domain.Task{Name: "second", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(first); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateTask(second); err != nil {
		t.Fatal(err)
	}
	trigger := &domain.ExternalTrigger{Name: "original", TargetTaskID: first.ID, Enabled: true}
	if err := st.CreateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(first.ID, true); err != nil {
		t.Fatal(err)
	}
	trigger.Name = "renamed"
	trigger.TargetTaskID = second.ID
	if err := st.UpdateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	loaded, err := st.GetExternalTrigger(trigger.ID)
	if err != nil || loaded.Name != "renamed" || loaded.TargetTaskID != second.ID || loaded.Key != trigger.Key {
		t.Fatalf("updated trigger = %+v, err=%v", loaded, err)
	}
	oldTarget, _ := st.GetTask(first.ID)
	if oldTarget.Enabled {
		t.Fatal("retargeting final source did not disable old target")
	}
	items, err := st.ListExternalTriggers()
	if err != nil || len(items) != 1 || items[0].ID != trigger.ID {
		t.Fatalf("list = %+v, err=%v", items, err)
	}
	if err := st.SetTaskEnabled(second.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := st.SetExternalTriggerEnabled(trigger.ID, false); err != nil {
		t.Fatal(err)
	}
	loaded, _ = st.GetExternalTrigger(trigger.ID)
	secondTarget, _ := st.GetTask(second.ID)
	if loaded.Enabled || secondTarget.Enabled {
		t.Fatalf("disable trigger=%t target=%t", loaded.Enabled, secondTarget.Enabled)
	}
	if err := st.SetExternalTriggerEnabled(trigger.ID, true); err != nil {
		t.Fatal(err)
	}
	if has, err := st.TaskHasEnabledTrigger(second.ID); err != nil || !has {
		t.Fatalf("enabled source=%t err=%v", has, err)
	}
}

func TestExternalTriggerValidationAndTaskCascade(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if err := st.CreateExternalTrigger(&domain.ExternalTrigger{}); !errors.Is(err, ErrInvalidTrigger) {
		t.Fatalf("blank trigger error = %v", err)
	}
	if err := st.CreateExternalTrigger(&domain.ExternalTrigger{Name: "missing", TargetTaskID: "missing", Enabled: true}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing target error = %v", err)
	}
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	trigger := &domain.ExternalTrigger{Name: "hook", TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteTask(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetExternalTrigger(trigger.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cascaded trigger error = %v", err)
	}
	if err := st.DeleteExternalTrigger("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing delete error = %v", err)
	}
	if _, err := st.RotateExternalTrigger("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing rotate error = %v", err)
	}
}

func TestDeletingOneOfMultipleExternalSourcesKeepsTaskEnabled(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	first := &domain.ExternalTrigger{Name: "one", TargetTaskID: task.ID, Enabled: true}
	second := &domain.ExternalTrigger{Name: "two", TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateExternalTrigger(first); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateExternalTrigger(second); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(task.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteExternalTrigger(first.ID); err != nil {
		t.Fatal(err)
	}
	loaded, _ := st.GetTask(task.ID)
	if !loaded.Enabled {
		t.Fatal("task was disabled while another enabled trigger remained")
	}
}

func TestExternalTriggerIsAutomaticSourceAndDeletionDeactivatesFinalSource(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	trigger := &domain.ExternalTrigger{Name: "hook", TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(task.ID, true); err != nil {
		t.Fatalf("enable trigger-backed task: %v", err)
	}
	if err := st.DeleteExternalTrigger(trigger.ID); err != nil {
		t.Fatal(err)
	}
	loaded, err := st.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Enabled {
		t.Fatal("task remained enabled after its final automatic source was deleted")
	}
}
