package clievaluate

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestJoinPostconditionChecksRecoveryAndCleanup(t *testing.T) {
	valid := domain.ParallelJoin{ID: "join", Name: "QA parallel", EntryTaskID: "entry", TargetTaskID: "publish", RecoveryTaskID: "recover", CleanupTaskID: "publish", RecoveryOn: []domain.JoinRecoveryCondition{domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnCanceled}, Branches: []domain.ParallelJoinBranch{{ID: "export-a", Name: "Export A", TaskID: "branch"}, {ID: "export-b", Name: "Export B", TaskID: "branch"}}, BranchCount: 2}
	current := valid
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if err := json.NewEncoder(w).Encode(current); err != nil {
			t.Error(err)
		}
	}))
	s := &session{api: api, vars: map[string]string{"join": "join", "base_task": "entry", "output_task": "publish", "second_task": "recover", "failure_task": "branch"}}
	for _, tc := range []struct {
		name   string
		change func(*domain.ParallelJoin)
		want   bool
	}{
		{"configured", func(*domain.ParallelJoin) {}, true},
		{"wrong_recovery", func(v *domain.ParallelJoin) { v.RecoveryTaskID = "wrong" }, false},
		{"wrong_cleanup", func(v *domain.ParallelJoin) { v.CleanupTaskID = "wrong" }, false},
		{"lost_causes", func(v *domain.ParallelJoin) { v.RecoveryOn = []domain.JoinRecoveryCondition{} }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			current = valid
			tc.change(&current)
			got, _, err := s.joinPost(context.Background(), step{Spec: CaseSpec{ID: "join.create"}})
			if err != nil || got != tc.want {
				t.Fatalf("independent postcondition accepted=%v want=%v err=%v", got, tc.want, err)
			}
		})
	}
	current = valid
	current.Name = "QA updated parallel"
	current.RecoveryOn = []domain.JoinRecoveryCondition{}
	if good, _, err := s.joinPost(context.Background(), step{Spec: CaseSpec{ID: "join.update"}}); err != nil || !good {
		t.Fatalf("explicit empty update=%v err=%v", good, err)
	}
	current.RecoveryOn = nil
	if good, _, err := s.joinPost(context.Background(), step{Spec: CaseSpec{ID: "join.update"}}); err != nil || good {
		t.Fatalf("missing versus empty policy not distinguished: %v %v", good, err)
	}
}
