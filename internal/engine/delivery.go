package engine

import (
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// eligibleCompletionTask reloads the exact durable automatic request and current
// task/group eligibility at each admission boundary. It never changes manual runs.
func (e *Engine) eligibleCompletionTask(id string) (domain.Task, bool) {
	delivery, task, ready, reason, err := e.store.GetCompletionAdmission(id)
	if err != nil {
		e.log.Error("engine: read delivery eligibility", "delivery", id, "err", err)
		return domain.Task{}, false
	}
	if !ready {
		e.resolveCompletionDelivery(id, delivery.JoinInstanceID, reason)
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
	default:
		return ""
	}
}
