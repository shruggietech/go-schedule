package trigger

import (
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-scheduler/internal/domain"
	"github.com/shruggietech/go-scheduler/internal/store"
)

func quiet() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

// fireCounter records target dispatches.
type fireCounter struct {
	mu    sync.Mutex
	fired []string
}

func (f *fireCounter) Fire(target string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fired = append(f.fired, target)
}
func (f *fireCounter) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.fired)
}

// setup creates two tasks (A source, B target) and a trigger A.success -> B.
func setup(t *testing.T, window time.Duration) (*store.Store, *domain.Task, *domain.Task, *domain.Trigger) {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	mkTask := func(name string) *domain.Task {
		sch := &domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=MINUTELY;INTERVAL=1"}
		_ = st.CreateSchedule(sch)
		task := &domain.Task{Name: name, Command: "x", Enabled: true, Timezone: "UTC", ScheduleID: sch.ID, State: domain.TaskActive}
		if err := st.CreateTask(task); err != nil {
			t.Fatal(err)
		}
		return task
	}
	a, b := mkTask("A"), mkTask("B")
	tr := &domain.Trigger{SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.OnSuccess, DedupWindow: window}
	if err := st.CreateTrigger(tr); err != nil {
		t.Fatal(err)
	}
	return st, a, b, tr
}

func TestOnCompletion_FiresTargetOnce(t *testing.T) {
	st, a, b, _ := setup(t, 5*time.Minute)
	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())

	now := time.Now().UTC()
	d.OnCompletion(a.ID, domain.OutcomeSuccess, "run-1", now)

	if fc.count() != 1 || fc.fired[0] != b.ID {
		t.Fatalf("expected target fired once, got %v", fc.fired)
	}
}

func TestOnCompletion_DedupWithinWindow(t *testing.T) {
	st, a, _, _ := setup(t, 5*time.Minute)
	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())

	now := time.Now().UTC()
	d.OnCompletion(a.ID, domain.OutcomeSuccess, "evt-1", now)
	d.OnCompletion(a.ID, domain.OutcomeSuccess, "evt-1", now.Add(time.Minute)) // duplicate within window

	if fc.count() != 1 {
		t.Fatalf("duplicate within window should not re-fire; got %d fires", fc.count())
	}
}

func TestOnCompletion_NewEventAfterWindow(t *testing.T) {
	st, a, _, _ := setup(t, time.Minute)
	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())

	now := time.Now().UTC()
	d.OnCompletion(a.ID, domain.OutcomeSuccess, "evt-1", now)
	d.OnCompletion(a.ID, domain.OutcomeSuccess, "evt-1", now.Add(2*time.Minute)) // past window → new event

	if fc.count() != 2 {
		t.Fatalf("event after window should re-fire; got %d fires", fc.count())
	}
}

func TestOnCompletion_OutcomeFilter(t *testing.T) {
	st, a, _, _ := setup(t, time.Minute) // trigger is OnSuccess
	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())

	d.OnCompletion(a.ID, domain.OutcomeFailure, "evt-1", time.Now().UTC())
	if fc.count() != 0 {
		t.Fatal("failure outcome should not fire a success-only trigger")
	}
}

func TestRecoverPending_RefiresUnexecuted(t *testing.T) {
	st, a, b, tr := setup(t, 5*time.Minute)

	// Simulate a crash: claim the event but never mark it executed.
	claimed, err := st.ClaimEvent(tr.ID, "evt-1", tr.DedupWindow, time.Now().UTC())
	if err != nil || !claimed {
		t.Fatalf("claim failed: ok=%v err=%v", claimed, err)
	}
	_ = a

	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())
	d.RecoverPending()

	if fc.count() != 1 || fc.fired[0] != b.ID {
		t.Fatalf("recovery should re-fire the unexecuted event once, got %v", fc.fired)
	}
	// After recovery the claim is executed → a second recovery does nothing.
	d.RecoverPending()
	if fc.count() != 1 {
		t.Fatalf("recovery should be idempotent after marking executed, got %d", fc.count())
	}
}

func TestRecoverPending_SkipsClaimWithMissingTrigger(t *testing.T) {
	st, _, _, tr := setup(t, 5*time.Minute)
	// Claim an event, then delete the trigger before recovery runs.
	if _, err := st.ClaimEvent(tr.ID, "orphan", tr.DedupWindow, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteTrigger(tr.ID); err != nil {
		t.Fatal(err)
	}

	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())
	d.RecoverPending() // GetTrigger fails → entry skipped, no panic, no fire

	if fc.count() != 0 {
		t.Fatalf("recovery should skip a claim whose trigger no longer exists, got %d", fc.count())
	}
}

func TestDispatcher_StoreErrorsHandledGracefully(t *testing.T) {
	st, a, _, _ := setup(t, time.Minute)
	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())

	// Closing the store makes subsequent queries fail; the dispatcher must log
	// and return rather than panic or fire.
	_ = st.Close()
	d.OnCompletion(a.ID, domain.OutcomeSuccess, "evt", time.Now().UTC())
	d.RecoverPending()

	if fc.count() != 0 {
		t.Fatalf("no targets should fire when the store errors, got %d", fc.count())
	}
}

func TestOnCompletion_OnAnyMatchesFailure(t *testing.T) {
	st, a, b, tr := setup(t, time.Minute)
	// Switch the trigger to OnAny by recreating it.
	_ = st.DeleteTrigger(tr.ID)
	any := &domain.Trigger{SourceTaskID: a.ID, TargetTaskID: b.ID, OnOutcome: domain.OnAny, DedupWindow: time.Minute}
	if err := st.CreateTrigger(any); err != nil {
		t.Fatal(err)
	}
	fc := &fireCounter{}
	d := New(st, fc.Fire, quiet())
	d.OnCompletion(a.ID, domain.OutcomeFailure, "evt", time.Now().UTC())
	if fc.count() != 1 {
		t.Fatalf("OnAny should fire on failure, got %d", fc.count())
	}
}
