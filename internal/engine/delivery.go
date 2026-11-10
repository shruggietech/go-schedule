package engine

import (
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// eligibleCompletionTask reloads the exact durable automatic request and current
// task/group eligibility at each admission boundary. A nonnil reservation is
// the caller's one reserved execution, which must be retired on rejection before
// its no-start result can publish successors. The initial drain passes nil.
// Manual runs are unchanged.
func (e *Engine) eligibleCompletionTask(id string, reservation *domain.Task) (domain.Task, bool) {
	delivery, task, ready, reason, err := e.store.GetCompletionAdmission(id)
	stopping := e.runCtx.Err() != nil
	if (err != nil || !ready || stopping) && reservation != nil {
		e.finish(*reservation)
	}
	if err != nil {
		e.log.Error("engine: read delivery eligibility", "delivery", id, "err", err)
		return domain.Task{}, false
	}
	if !ready {
		e.resolveCompletionDelivery(id, delivery.JoinInstanceID, reason)
		return domain.Task{}, false
	}
	if stopping {
		return domain.Task{}, false
	}
	return task, true
}

func completionJoinRole(role domain.CompletionDeliveryRole) string {
	switch role {
	case domain.DeliveryRoleJoinBranch:
		return "branch"
	case domain.DeliveryRoleJoinTarget:
		return "target"
	case domain.DeliveryRoleJoinRecovery:
		return "recovery"
	case domain.DeliveryRoleJoinCleanup:
		return "cleanup"
	default:
		return ""
	}
}

// wakeCompletionDeliveries coalesces committed automatic work into the existing
// event loop. In particular, no-start resolutions can create new work without
// a subsequent native completion to wake it.
func (e *Engine) wakeCompletionDeliveries() {
	if e.runCtx.Err() != nil {
		return
	}
	select {
	case e.completionWake <- struct{}{}:
	default:
	}
}
