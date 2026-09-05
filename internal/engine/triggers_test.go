package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestFireExternalTriggerRecordsProvenance(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive, OverlapPolicy: domain.OverlapAllowConcurrent}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	trigger := &domain.ExternalTrigger{Name: "hook", TargetTaskID: task.ID, Enabled: true}
	if err := st.CreateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	if err := st.SetTaskEnabled(task.ID, true); err != nil {
		t.Fatal(err)
	}
	e := newEngine(st, instantRunner{})
	if id, err := e.FireExternalTrigger(trigger.Key); err != nil || id != trigger.ID {
		t.Fatalf("fire id=%q err=%v", id, err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		runs, err := st.ListRuns(task.ID, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(runs) == 1 {
			if runs[0].Trigger != domain.TriggerExternal || runs[0].SourceTriggerID != trigger.ID {
				t.Fatalf("run provenance = %+v", runs[0])
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("triggered run was not recorded")
}

func TestFireExternalTriggerRejectsUnknownDisabledAndBlocked(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	e := newEngine(st, instantRunner{})
	if _, err := e.FireExternalTrigger("gst_unknown"); !errors.Is(err, ErrTriggerUnknown) {
		t.Fatalf("unknown error = %v", err)
	}
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	trigger := &domain.ExternalTrigger{Name: "hook", TargetTaskID: task.ID, Enabled: false}
	if err := st.CreateExternalTrigger(trigger); err != nil {
		t.Fatal(err)
	}
	if _, err := e.FireExternalTrigger(trigger.Key); !errors.Is(err, ErrTriggerDisabled) {
		t.Fatalf("disabled error = %v", err)
	}
}

func BenchmarkFireExternalTriggerDecision(b *testing.B) {
	st, err := store.Open(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer st.Close()
	task := &domain.Task{Name: "target", Command: "echo", State: domain.TaskActive, OverlapPolicy: domain.OverlapAllowConcurrent}
	_ = st.CreateTask(task)
	trigger := &domain.ExternalTrigger{Name: "hook", TargetTaskID: task.ID, Enabled: true}
	_ = st.CreateExternalTrigger(trigger)
	_ = st.SetTaskEnabled(task.ID, true)
	e := newEngine(st, instantRunner{})
	b.SetParallelism(5)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.FireExternalTrigger(trigger.Key)
		}
	})
	b.StopTimer()
	e.runWG.Wait()
}
