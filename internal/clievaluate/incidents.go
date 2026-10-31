package clievaluate

import (
	"context"
	"errors"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/incident"
)

// incidentWatchdog collects while a slow case is alive, independently of its API.
// Its finish capability always joins the owned goroutine before session cleanup.
func (s *session) incidentWatchdog(ctx context.Context, d *executionDiagnostics, expectedExit ...int) func(execution) execution {
	expected := 0
	if len(expectedExit) > 0 {
		expected = expectedExit[0]
	}
	if s.incidentError != nil {
		return func(result execution) execution {
			return attachIncident(result, map[string]any{"kind": "go-schedule.incident", "version": 1, "category": "owned_collection_setup", "original": incidentOriginal(ctx, d, result), "collection": incident.Collection{Status: "unavailable", Detail: "The owned incident directory could not be prepared. Retain the run and inspect directory access before another diagnostic attempt. The original case outcome is unchanged."}})
		}
	}
	if s.incidents == nil || s.daemon == nil {
		return func(result execution) execution { return result }
	}
	delay := 5 * time.Second
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline)/2 < delay {
		delay = time.Until(deadline) / 2
	}
	if delay < time.Millisecond {
		delay = time.Millisecond
	}
	owned, stop := context.WithCancel(context.Background())
	done := make(chan map[string]any, 1)
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-owned.Done():
			done <- nil
			return
		case <-timer.C:
		}
		trigger := "live"
		if ctx.Err() != nil {
			trigger = "after_cancellation"
		}
		done <- s.collectIncident(owned, d, trigger, delay)
	}()
	return func(result execution) execution {
		stop()
		entry := <-done
		if entry == nil {
			failed := ctx.Err() != nil || result.Blocked != "" || result.Code != expected
			for _, passed := range result.Checks {
				failed = failed || !passed
			}
			if !failed {
				return result
			}
			trigger := "after_case_failure"
			if ctx.Err() != nil {
				trigger = "after_cancellation"
			}
			entry = s.collectIncident(context.Background(), d, trigger, delay)
		}
		entry["original"] = incidentOriginal(ctx, d, result)
		return attachIncident(result, entry)
	}
}

func incidentOriginal(ctx context.Context, d *executionDiagnostics, result execution) map[string]any {
	return map[string]any{"command_ran": d.CommandRan, "native_exit": result.Code, "inspection_complete": d.AfterComplete, "canceled": ctx.Err() != nil}
}

func (s *session) collectIncident(parent context.Context, d *executionDiagnostics, trigger string, threshold time.Duration) map[string]any {
	d.mu.Lock()
	phase := d.current
	d.mu.Unlock()
	if phase == "" {
		phase = "not_started"
	}
	started := time.Now()
	collection, cancel := context.WithTimeout(parent, 2*time.Second)
	snapshot, err := s.incidents.Collect(collection, s.daemon.PID, phase)
	cancel()
	status, detail := "collected", ""
	if err != nil {
		status = "failed"
		if errors.Is(err, context.DeadlineExceeded) {
			status = "timed_out"
		}
		if errors.Is(err, context.Canceled) {
			status = "canceled"
		}
		detail = "Owned daemon collection: " + err.Error() + ". Retain the private incident directory and inspect daemon logs. The original case outcome is unchanged."
		if status == "canceled" {
			detail = "Incident collection stopped when the case finished. The original case outcome is unchanged."
		}
	}
	entry := map[string]any{"kind": "go-schedule.incident", "version": 1, "category": "owned_case", "phase": phase, "trigger_state": trigger, "threshold_ms": threshold.Milliseconds(), "collection": incident.Collection{Status: status, Detail: detail, DurationMS: time.Since(started).Milliseconds()}, "private_directory": s.incidents.Dir}
	if err == nil {
		entry["snapshot"] = snapshot
	}
	return entry
}

func attachIncident(result execution, entry map[string]any) execution {
	state := observationState(result.State)
	state["incident"] = entry
	result.State = mustObservation(state)
	return result
}
