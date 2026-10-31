package clievaluate

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
	"github.com/shruggietech/go-schedule/v2/internal/catchup"
	"github.com/shruggietech/go-schedule/v2/internal/clock"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/engine"
	"github.com/shruggietech/go-schedule/v2/internal/schedule"
	"github.com/shruggietech/go-schedule/v2/internal/store"
	"github.com/shruggietech/go-schedule/v2/internal/userprefs"
)

func policyResult(observed any, passed bool, err error) execution {
	detail := ""
	if err != nil {
		passed = false
		detail = err.Error()
	}
	return execution{Code: 0, Checks: []bool{passed, err == nil}, State: mustObservation(map[string]any{"observation": observed, "diagnostic": detail})}
}

func calendarPolicy(name string) execution {
	anchor := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	observations := []NamedValue{}
	good := true
	type boundary struct {
		expression, zone, after, want string
		policy                        domain.SchedulePolicy
	}
	var cases []boundary
	if name == "dst-gap-fold" {
		cases = []boundary{{"every day at 02:30", "America/New_York", "2026-03-07T08:00:00Z", "2026-03-08T07:00:00Z", domain.SchedulePolicy{DSTGap: domain.DSTGapNextValid}}, {"every day at 02:30", "America/New_York", "2026-03-07T08:00:00Z", "2026-03-09T06:30:00Z", domain.SchedulePolicy{DSTGap: domain.DSTGapSkip}}, {"every day at 01:30", "America/New_York", "2026-11-01T05:00:00Z", "2026-11-01T05:30:00Z", domain.SchedulePolicy{DSTOverlap: domain.DSTOverlapFirst}}, {"every day at 01:30", "America/New_York", "2026-11-01T05:00:00Z", "2026-11-01T06:30:00Z", domain.SchedulePolicy{DSTOverlap: domain.DSTOverlapLast}}, {"every day at 01:30", "America/New_York", "2026-11-01T05:45:00Z", "2026-11-01T06:30:00Z", domain.SchedulePolicy{DSTOverlap: domain.DSTOverlapBoth}}}
	} else {
		cases = []boundary{{"last day of every month at 09:00", "UTC", "2028-02-01T00:00:00Z", "2028-02-29T09:00:00Z", domain.SchedulePolicy{}}, {"nearest weekday to the 31st of every month at 09:00", "UTC", "2026-04-01T00:00:00Z", "2026-05-29T09:00:00Z", domain.SchedulePolicy{MissingDate: domain.MissingDateSkip}}, {"nearest weekday to the 31st of every month at 09:00", "UTC", "2026-04-01T00:00:00Z", "2026-04-30T09:00:00Z", domain.SchedulePolicy{MissingDate: domain.MissingDateLastValid}}, {"nearest weekday to the 31st of every month at 09:00", "UTC", "2026-04-01T00:00:00Z", "2026-05-01T09:00:00Z", domain.SchedulePolicy{MissingDate: domain.MissingDateNextValid}}}
	}
	for _, boundary := range cases {
		sch, err := schedule.Parse(boundary.expression, boundary.zone, anchor)
		if err != nil {
			return policyResult(observations, false, err)
		}
		after, err := time.Parse(time.RFC3339, boundary.after)
		if err != nil {
			return policyResult(observations, false, err)
		}
		actual, ok, err := schedule.NextRunWithPolicy(sch, boundary.zone, boundary.policy, after)
		observations = append(observations, NamedValue{Name: boundary.expression + " " + boundary.after + " expected " + boundary.want, Value: actual.Format(time.RFC3339)})
		good = good && ok && err == nil && actual.Format(time.RFC3339) == boundary.want
	}
	return policyResult(observations, good, nil)
}

func bundlePolicy() execution {
	doc := bundle.Document{Schema: bundle.SchemaV2, Tasks: []bundle.Task{{PortableID: "task", Name: "QA", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "none", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}, Watchers: []bundle.Watcher{{PortableID: "watcher", Name: "QA", Kind: "directory", TargetTaskID: "missing", Debounce: "10ms", Stability: "10ms"}}}
	bad := bundle.Validate(doc)
	doc.Watchers[0].TargetTaskID = "task"
	valid := bundle.Validate(doc)
	return policyResult(map[string]any{"dangling_reference": bad, "repaired_reference": valid}, len(bad) > 0 && len(valid) == 0, nil)
}

func expiredPairingPolicy() execution {
	st, err := store.Open(":memory:")
	if err != nil {
		return policyResult(nil, false, err)
	}
	defer func() { _ = st.Close() }()
	now := time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)
	pair := domain.PairingSession{ID: uuid.NewString(), DisplayName: "QA expired", Kind: domain.ActorKindCLI, Capability: domain.CapabilityObserve, State: domain.PairingActive, AttemptsRemaining: 5, CreatedAt: now.Add(-time.Hour), ExpiresAt: now.Add(-time.Minute)}
	salt := make([]byte, 16)
	verifier := make([]byte, 32)
	if err := st.CreatePairing(pair, salt, verifier); err != nil {
		return policyResult(nil, false, err)
	}
	_, _, exchangeErr := st.ExchangePairing(pair.ID, verifier, domain.Actor{ID: uuid.NewString()}, domain.ClientCredential{ID: uuid.NewString()}, make([]byte, 32), now)
	credentials, err := st.ListCredentials()
	return policyResult(map[string]any{"expired_rejected": errors.Is(exchangeErr, store.ErrEnrollmentRejected), "credential_count": len(credentials)}, errors.Is(exchangeErr, store.ErrEnrollmentRejected) && len(credentials) == 0, err)
}

func popupPolicy(root string) execution {
	prefs := userprefs.NewStore(userprefs.Options{Path: root})
	enabled := true
	conditions := []string{"failure"}
	severities := []string{"error"}
	ids := []string{"qa-source"}
	saved, err := prefs.Patch(userprefs.Patch{PopupEnabled: &enabled, PopupConditions: &conditions, PopupSeverities: &severities, PopupDaemonIDs: &ids})
	if err != nil {
		return policyResult(nil, false, err)
	}
	loaded, err := prefs.Load()
	invalid := userprefs.ValidatePopups(userprefs.PopupPreferences{Enabled: true, Conditions: []string{"unknown"}, Severities: severities})
	return policyResult(map[string]any{"saved": saved, "loaded": loaded, "invalid_filter_rejected": invalid != nil}, equalObservation(saved, loaded) && invalid != nil, err)
}

type controlledRunner struct {
	started  chan domain.Run
	release  chan struct{}
	mu       sync.Mutex
	count    int
	now      time.Time
	clockNow func() time.Time
}

func (r *controlledRunner) Run(ctx context.Context, task domain.Task, at time.Time, trigger domain.RunTrigger) domain.Run {
	run := domain.Run{TaskID: task.ID, ScheduledFor: at, Trigger: trigger, Outcome: domain.OutcomeSuccess}
	select {
	case r.started <- run:
	case <-ctx.Done():
		return run
	}
	select {
	case <-r.release:
	case <-ctx.Done():
	}
	r.mu.Lock()
	r.count++
	r.mu.Unlock()
	end := r.now
	if r.clockNow != nil {
		end = r.clockNow()
	}
	run.EndedAt = &end
	return run
}

func enginePolicy(ctx context.Context, name string) execution {
	observations := []any{}
	policies := []domain.OverlapPolicy{domain.OverlapQueueOne, domain.OverlapSkip, domain.OverlapAllowConcurrent}
	if name == "catchup-startup" {
		policies = policies[:1]
	}
	for _, policy := range policies {
		st, err := store.Open(":memory:")
		if err != nil {
			return policyResult(observations, false, err)
		}
		now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
		clk := clock.NewFake(now)
		runner := &controlledRunner{started: make(chan domain.Run, 8), release: make(chan struct{}, 8), now: now}
		sch := schedule.NewStartup("@reboot")
		if name != "catchup-startup" {
			sch = domain.Schedule{Kind: domain.ScheduleRecurring, RRULE: "FREQ=HOURLY"}
		}
		err = st.CreateSchedule(&sch)
		task := domain.Task{Name: "QA engine", Command: "fixture", Enabled: true, Timezone: "UTC", State: domain.TaskActive, ScheduleID: sch.ID, OverlapPolicy: policy, CatchupPolicy: domain.CatchupOne}
		if err == nil {
			err = st.CreateTask(&task)
		}
		if err != nil {
			_ = st.Close()
			return policyResult(observations, false, err)
		}
		if name == "catchup-startup" {
			group := domain.Group{Name: "QA disabled ancestor", Enabled: false}
			if err := st.CreateGroup(&group); err != nil {
				_ = st.Close()
				return policyResult(observations, false, err)
			}
			for _, definition := range []struct {
				enabled bool
				group   string
			}{{false, ""}, {true, group.ID}} {
				inactive := domain.Task{Name: "QA inactive startup", Command: "fixture", Enabled: definition.enabled, GroupID: definition.group, Timezone: "UTC", State: domain.TaskActive, ScheduleID: sch.ID}
				if err := st.CreateTask(&inactive); err != nil {
					_ = st.Close()
					return policyResult(observations, false, err)
				}
			}
		}
		e := engine.New(st, clk, runner, slog.New(slog.NewTextHandler(io.Discard, nil)), 2)
		child, cancel := context.WithCancel(ctx)
		done := make(chan error, 1)
		go func() { done <- e.Start(child) }()
		select {
		case <-e.Ready():
		case <-ctx.Done():
			cancel()
			<-done
			_ = st.Close()
			return policyResult(observations, false, ctx.Err())
		}
		if name != "catchup-startup" {
			err = e.RunNow(task.ID)
		}
		if err == nil {
			select {
			case <-runner.started:
			case <-ctx.Done():
				err = ctx.Err()
			}
		}
		if err == nil {
			if name == "catchup-startup" {
				e.Reload()
			} else {
				err = e.RunNow(task.ID)
				if err == nil {
					err = e.RunNow(task.ID)
				}
			}
		}
		for range 4 {
			runner.release <- struct{}{}
		}
		expected := 2
		if policy == domain.OverlapAllowConcurrent {
			expected = 3
		}
		if policy == domain.OverlapSkip || name == "catchup-startup" {
			expected = 1
		}
		if err == nil {
			err = waitUntil(ctx, func() bool {
				runs, e := st.ListRuns(task.ID, 0)
				count := 0
				for _, run := range runs {
					if run.Outcome == domain.OutcomeSuccess {
						count++
					}
				}
				return e == nil && count >= expected
			})
		}
		runs, readErr := st.ListRuns(task.ID, 0)
		allRuns, allErr := st.ListRuns("", 0)
		readErr = errors.Join(readErr, allErr)
		if name == "catchup-startup" && len(allRuns) != len(runs) {
			readErr = errors.Join(readErr, errors.New("disabled startup task or ancestor dispatched unexpectedly"))
		}
		cancel()
		stopErr := <-done
		if errors.Is(stopErr, context.Canceled) {
			stopErr = nil
		}
		_ = st.Close()
		count := 0
		for _, run := range runs {
			if run.Outcome == domain.OutcomeSuccess {
				count++
			}
		}
		observations = append(observations, map[string]any{"policy": policy, "runs": runs, "expected_successes": expected})
		if err != nil || readErr != nil || stopErr != nil || count != expected {
			return policyResult(observations, false, errors.Join(err, readErr, stopErr))
		}
	}
	if name == "catchup-startup" {
		now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
		sch, err := schedule.Parse("every hour", "UTC", now.Add(-24*time.Hour))
		if err != nil {
			return policyResult(observations, false, err)
		}
		for _, policy := range []domain.CatchupPolicy{domain.CatchupOne, domain.CatchupNone} {
			decision, err := catchup.EvaluateWithPolicy(sch, "UTC", now.Add(-5*time.Hour), true, policy, domain.SchedulePolicy{}, now)
			observations = append(observations, map[string]any{"catchup_policy": policy, "decision": decision})
			if err != nil || decision.ShouldCatchUp != (policy == domain.CatchupOne) {
				return policyResult(observations, false, err)
			}
		}
	}
	return policyResult(observations, true, nil)
}

func oneOffPolicy(ctx context.Context) execution {
	st, err := store.Open(":memory:")
	if err != nil {
		return policyResult(nil, false, err)
	}
	defer func() { _ = st.Close() }()
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	clk := clock.NewFake(now)
	at := now.Add(time.Minute)
	sch := schedule.NewOneOff(at)
	if err := st.CreateSchedule(&sch); err != nil {
		return policyResult(nil, false, err)
	}
	task := domain.Task{Name: "QA one-off", Command: "fixture", Enabled: true, Timezone: "UTC", State: domain.TaskActive, ScheduleID: sch.ID, OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupNone}
	if err := st.CreateTask(&task); err != nil {
		return policyResult(nil, false, err)
	}
	runner := &controlledRunner{started: make(chan domain.Run, 4), release: make(chan struct{}, 4), clockNow: clk.Now}
	e := engine.New(st, clk, runner, slog.New(slog.NewTextHandler(io.Discard, nil)), 1)
	child, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- e.Start(child) }()
	select {
	case <-e.Ready():
	case <-ctx.Done():
		cancel()
		<-done
		return policyResult(nil, false, ctx.Err())
	}
	err = waitUntil(ctx, func() bool { return clk.Waiters() > 0 })
	if err == nil {
		clk.Advance(time.Minute)
		select {
		case run := <-runner.started:
			if !run.ScheduledFor.Equal(at) {
				err = errors.New("one-off dispatch used the wrong scheduled instant")
			}
		case <-ctx.Done():
			err = ctx.Err()
		}
	}
	runner.release <- struct{}{}
	if err == nil {
		err = waitUntil(ctx, func() bool { stored, e := st.GetTask(task.ID); return e == nil && stored.State == domain.TaskCompleted })
	}
	clk.Advance(time.Hour)
	cancel()
	stopErr := <-done
	if errors.Is(stopErr, context.Canceled) {
		stopErr = nil
	}
	runs, readErr := st.ListRuns(task.ID, 0)
	stored, taskErr := st.GetTask(task.ID)
	return policyResult(map[string]any{"expected_at": at, "task": stored, "runs": runs}, len(runs) == 1 && stored.State == domain.TaskCompleted, errors.Join(err, stopErr, readErr, taskErr))
}
