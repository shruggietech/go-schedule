package integration

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/commandexample"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/engine"
	"github.com/shruggietech/go-schedule/internal/executor"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestGuidedPlatformExampleCreatesRunsAndRecordsRecognizableActivity(t *testing.T) {
	suggestion, ok := commandexample.ForPlatform(runtime.GOOS)
	if !ok {
		t.Skip("no guided example for this platform")
	}
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task := domain.Task{Name: "platform check", Command: suggestion.Program, Args: suggestion.Args, Enabled: false, Timezone: "Local", OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne, State: domain.TaskActive}
	if err := st.CreateTask(&task); err != nil {
		t.Fatal(err)
	}
	eng := engine.New(st, clock.NewReal(), executor.New(64<<10), quietLogger(), 1)
	runs := make(chan domain.Run, 1)
	eng.SetOnRun(func(run domain.Run) { runs <- run })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- eng.Start(ctx) }()
	select {
	case <-eng.Ready():
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not become ready")
	}
	if err := eng.RunNow(task.ID); err != nil {
		t.Fatal(err)
	}
	select {
	case run := <-runs:
		if run.Outcome != domain.OutcomeSuccess || run.Trigger != domain.TriggerManual || !suggestion.Recognizes(run.Output) {
			t.Fatalf("activity=%+v output=%q", run, run.Output)
		}
		stored, err := st.GetRun(run.ID)
		if err != nil || stored.Outcome != domain.OutcomeSuccess || !suggestion.Recognizes(stored.Output) {
			t.Fatalf("stored activity=%+v err=%v", stored, err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("guided task did not finish")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("engine did not stop")
	}
}
