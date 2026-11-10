package mcpobserve

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestJoinRecoveryObservationKeepsSeparateExactStages(t *testing.T) {
	d := joinDefinitionData(domain.ParallelJoin{ID: "join", RecoveryTaskID: "repair", RecoveryOn: []domain.JoinRecoveryCondition{}, CleanupTaskID: "clean"})
	raw, _ := json.Marshal(d)
	if !strings.Contains(string(raw), "\"recovery_on\":[]") || d["recovery_task_id"] != "repair" || d["cleanup_task_id"] != "clean" {
		t.Fatalf("configuration=%s", raw)
	}
	v := joinInstanceData(domain.ParallelJoinInstance{ID: "instance", SelectedRoute: domain.JoinRouteRecovery, RecoveryTaskID: "repair", RecoveryDeliveryID: "repair-delivery", RecoveryRunID: "repair-run", RecoverySettlement: domain.JoinSettlementFailure, RecoveryReason: "actual failure", CleanupTaskID: "clean", CleanupDeliveryID: "cleanup-delivery", CleanupSettlement: domain.JoinSettlementNotStarted, CleanupReason: "disabled", Branches: []domain.ParallelJoinInvocation{{ID: "invocation", Settlement: domain.JoinSettlementSkipped}}})
	for k, want := range map[string]any{"selected_route": domain.JoinRouteRecovery, "recovery_run_id": "repair-run", "recovery_delivery_id": "repair-delivery", "recovery_reason": "actual failure", "cleanup_delivery_id": "cleanup-delivery", "cleanup_reason": "disabled"} {
		if v[k] != want {
			t.Errorf("%s=%v want %v", k, v[k], want)
		}
	}
	raw, _ = json.Marshal(v)
	if !strings.Contains(string(raw), "\"settlement\":\"skipped\"") || !strings.Contains(string(raw), "\"cleanup_settlement\":\"not_started\"") {
		t.Fatalf("settlement=%s", raw)
	}
}
