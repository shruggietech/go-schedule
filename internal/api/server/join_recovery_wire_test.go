package server

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS145JoinDTORecoveryWireAndBaseCompatibility(t *testing.T) {
	for _, update := range []bool{false, true} {
		t.Run(map[bool]string{false: "create", true: "update"}[update], func(t *testing.T) {
			raw := []byte(`{"name":"block","entry_task_id":"entry","target_task_id":"publish","branches":[{"id":"a","name":"A","task_id":"export"}],"recovery_task_id":"recover","recovery_on":[],"cleanup_task_id":"cleanup"}`)
			var target any = &JoinCreateRequest{}
			if update {
				target = &JoinUpdateRequest{}
			}
			if err := json.Unmarshal(raw, target); err != nil {
				t.Fatalf("complete routing request rejected: %v", err)
			}
			got, err := json.Marshal(target)
			if err != nil {
				t.Fatal(err)
			}
			for _, needle := range []string{`"recovery_task_id":"recover"`, `"recovery_on":[]`, `"cleanup_task_id":"cleanup"`} {
				if !strings.Contains(string(got), needle) {
					t.Fatalf("lost %s in %s", needle, got)
				}
			}
		})
	}
	for _, target := range []any{&JoinCreateRequest{Name: "base"}, &JoinUpdateRequest{}} {
		got, err := json.Marshal(target)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(got), "recovery") || strings.Contains(string(got), "cleanup") {
			t.Fatalf("base packet breaks older strict server: %s", got)
		}
	}
	var clear JoinUpdateRequest
	if err := json.Unmarshal([]byte(`{"recovery_task_id":"","cleanup_task_id":""}`), &clear); err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(clear)
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{`"recovery_task_id":""`, `"cleanup_task_id":""`} {
		if !strings.Contains(string(got), needle) {
			t.Fatalf("lost explicit clearing %s", got)
		}
	}
}

func TestS145MergeExactRecoveryCleanupAndTerminalPrecedence(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	instance := domain.ParallelJoinInstance{ID: "instance", RecoveryTaskID: "recover", RecoveryDeliveryID: "rd", RecoveryDeliveryState: domain.DeliveryClaimed, CleanupTaskID: "cleanup", CleanupDeliveryID: "cd", CleanupDeliveryState: domain.DeliveryClaimed}
	active := map[string]domain.Run{
		"rd": {ID: "rr", TaskID: "recover", JoinInstanceID: "instance", JoinRole: "recovery", StartedAt: &now, CancellationStatus: "failed", CancellationError: strings.Repeat("x", 100)},
		"cd": {ID: "cr", TaskID: "cleanup", JoinInstanceID: "instance", JoinRole: "cleanup", StartedAt: &now},
	}
	mergeJoinActive(&instance, active, 16)
	if instance.RecoveryRunID != "rr" || instance.CleanupRunID != "cr" || instance.RecoveryCancellationStatus != "failed" || len(instance.RecoveryCancellationError) != 16 {
		t.Fatalf("stage live metadata missing: %+v", instance)
	}
	if instance.RecoverySettlement != "" || instance.CleanupSettlement != "" {
		t.Fatal("live observation fabricated settlement")
	}
	instance.RecoveryOutcome = domain.OutcomeFailure
	instance.RecoverySettlement = domain.JoinSettlementFailure
	instance.RecoveryRunID = "durable"
	mergeJoinActive(&instance, active, 16)
	if instance.RecoveryRunID != "durable" {
		t.Fatal("sampled live result replaced durable terminal evidence")
	}
	instance.CleanupRunID = ""
	active["cd"] = domain.Run{ID: "wrong", TaskID: "cleanup", JoinInstanceID: "other", JoinRole: "cleanup"}
	mergeJoinActive(&instance, active, 16)
	if instance.CleanupRunID != "" {
		t.Fatal("other instance contributed cleanup")
	}
}
