// Package trigger evaluates task-completion events (the v1 event source) and
// fires dependent tasks. Delivery is at-least-once and de-duplicated within a
// per-trigger window/key so a single logical completion causes a single
// execution. The engine calls OnCompletion after each run is recorded.
package trigger

import (
	"log/slog"
	"time"

	"github.com/shruggietech/go-scheduler/internal/domain"
	"github.com/shruggietech/go-scheduler/internal/store"
)

// Fire dispatches the target task of a trigger (the engine supplies this so the
// run is recorded with the event trigger source).
type Fire func(targetTaskID string)

// Dispatcher matches completion events to triggers and fires targets, with
// durable de-duplication via the store's dedup ledger.
type Dispatcher struct {
	store *store.Store
	fire  Fire
	log   *slog.Logger
}

// New builds a Dispatcher.
func New(st *store.Store, fire Fire, log *slog.Logger) *Dispatcher {
	return &Dispatcher{store: st, fire: fire, log: log}
}

// OnCompletion is invoked when a source task's run completes. eventKey
// identifies the logical event (the source run ID); repeats of the same key
// within a trigger's window are de-duplicated.
func (d *Dispatcher) OnCompletion(sourceTaskID string, outcome domain.RunOutcome, eventKey string, now time.Time) {
	triggers, err := d.store.ListTriggersBySource(sourceTaskID)
	if err != nil {
		d.log.Error("trigger: list by source", "task", sourceTaskID, "err", err)
		return
	}
	for _, tr := range triggers {
		if !outcomeMatches(tr.OnOutcome, outcome) {
			continue
		}
		key := tr.DedupKey
		if key == "" {
			key = eventKey
		}
		claimed, err := d.store.ClaimEvent(tr.ID, key, tr.DedupWindow, now)
		if err != nil {
			d.log.Error("trigger: claim event", "trigger", tr.ID, "err", err)
			continue
		}
		if !claimed {
			d.log.Debug("trigger: duplicate event deduplicated", "trigger", tr.ID, "key", key)
			continue
		}
		d.fire(tr.TargetTaskID)
		if err := d.store.MarkExecuted(tr.ID, key); err != nil {
			d.log.Error("trigger: mark executed", "trigger", tr.ID, "err", err)
		}
	}
}

// RecoverPending re-fires any claimed-but-unexecuted events after a restart,
// providing at-least-once delivery across daemon restarts.
func (d *Dispatcher) RecoverPending() {
	pending, err := d.store.PendingClaims()
	if err != nil {
		d.log.Error("trigger: pending claims", "err", err)
		return
	}
	for _, p := range pending {
		tr, err := d.store.GetTrigger(p.TriggerID)
		if err != nil {
			continue
		}
		d.log.Warn("trigger: recovering pending event after restart", "trigger", tr.ID, "key", p.DedupKey)
		d.fire(tr.TargetTaskID)
		if err := d.store.MarkExecuted(tr.ID, p.DedupKey); err != nil {
			d.log.Error("trigger: mark executed (recovery)", "trigger", tr.ID, "err", err)
		}
	}
}

func outcomeMatches(on domain.TriggerOutcome, outcome domain.RunOutcome) bool {
	switch on {
	case domain.OnAny:
		return outcome == domain.OutcomeSuccess || outcome == domain.OutcomeFailure
	case domain.OnFailure:
		return outcome == domain.OutcomeFailure
	default: // OnSuccess
		return outcome == domain.OutcomeSuccess
	}
}
