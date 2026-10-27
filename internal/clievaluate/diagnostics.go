package clievaluate

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
)

type phaseTiming struct {
	Name       string `json:"name"`
	DurationMS int64  `json:"duration_ms"`
}
type requestTiming struct {
	Phase string `json:"phase"`
	client.RequestObservation
}
type executionDiagnostics struct {
	mu             sync.Mutex
	phases         []phaseTiming
	requests       []requestTiming
	dropped        int
	current        string
	started        time.Time
	BeforeComplete bool `json:"before_complete"`
	AfterComplete  bool `json:"after_complete"`
	CommandRan     bool `json:"command_ran"`
}

func (d *executionDiagnostics) phase(ctx context.Context, name string) context.Context {
	d.mu.Lock()
	if d.current != "" {
		d.phases = append(d.phases, phaseTiming{d.current, time.Since(d.started).Milliseconds()})
	}
	d.current = name
	d.started = time.Now()
	d.mu.Unlock()
	return client.WithRequestObserver(ctx, func(o client.RequestObservation) {
		d.mu.Lock()
		defer d.mu.Unlock()
		if len(d.requests) < 256 {
			d.requests = append(d.requests, requestTiming{name, o})
		} else {
			d.dropped++
		}
	})
}
func (d *executionDiagnostics) finish(result execution) execution {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.current != "" {
		d.phases = append(d.phases, phaseTiming{d.current, time.Since(d.started).Milliseconds()})
	}
	state := observationState(result.State)
	state["diagnostics"] = map[string]any{"phases": d.phases, "requests": d.requests, "requests_omitted": d.dropped, "before_complete": d.BeforeComplete, "after_complete": d.AfterComplete, "command_ran": d.CommandRan, "blocked_reason": result.Blocked, "execution_diagnostic": result.Detail, "process_cleanup_error": errorText(result.Cleanup)}
	result.State = mustObservation(state)
	return result
}

func observationState(raw []byte) map[string]any {
	state := map[string]any{}
	if len(raw) > 0 {
		var object map[string]json.RawMessage
		switch {
		case !json.Valid(raw):
			state["observation_raw"] = string(raw)
			state["state_decode_error"] = "The original case observation could not be decoded."
		case json.Unmarshal(raw, &object) == nil && object != nil:
			for name, value := range object {
				state[name] = value
			}
		default:
			state["observation"] = json.RawMessage(raw)
		}
	}
	return state
}
