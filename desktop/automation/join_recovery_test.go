package automation

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type recoveryDraftBackend struct {
	fakeBackend
	current domain.ParallelJoin
	packet  []byte
	calls   int
}

func (b *recoveryDraftBackend) GetJoinDefinition(context.Context, string) (domain.ParallelJoin, error) {
	return b.current, nil
}
func (b *recoveryDraftBackend) UpdateJoin(_ context.Context, _ string, r server.JoinUpdateRequest) (domain.ParallelJoin, error) {
	b.calls++
	b.packet, _ = json.Marshal(r)
	return b.current, nil
}
func TestJoinRecoveryDraftSendsChangedFieldsAndExplicitEmpty(t *testing.T) {
	b := &recoveryDraftBackend{current: domain.ParallelJoin{ID: "join", Name: "Exports", EntryTaskID: "entry", TargetTaskID: "target", Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: "branch"}}, UpdatedAt: time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)}}
	var d JoinDraft
	if err := json.Unmarshal([]byte(`{"id":"join","name":"Exports","entryTaskId":"entry","targetTaskId":"target","recoveryTaskId":"repair","recoveryOn":[],"cleanupTaskId":"cleanup","branches":[{"id":"a","name":"A","task_id":"branch"}],"originalUpdatedAt":"2026-11-10T00:00:00Z"}`), &d); err != nil {
		t.Fatal(err)
	}
	r := NewService(b).SaveJoin(context.Background(), d)
	if r.Outcome != "accepted" || b.calls != 1 {
		t.Fatalf("result=%+v calls=%d", r, b.calls)
	}
	var p map[string]json.RawMessage
	if err := json.Unmarshal(b.packet, &p); err != nil {
		t.Fatal(err)
	}
	if string(p["recovery_task_id"]) != `"repair"` || string(p["cleanup_task_id"]) != `"cleanup"` || string(p["recovery_on"]) != "[]" {
		t.Fatalf("packet=%s", b.packet)
	}
	if _, ok := p["branches"]; ok {
		t.Fatalf("untouched branches replaced: %s", b.packet)
	}
}

func TestJoinRecoverySparseSavePreservesUntouchedPolicy(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(map[bool]string{false: "configured-causes", true: "explicit-empty"}[empty], func(t *testing.T) {
			conditions := []domain.JoinRecoveryCondition{domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnCanceled}
			if empty {
				conditions = []domain.JoinRecoveryCondition{}
			}
			b := &recoveryDraftBackend{current: domain.ParallelJoin{ID: "join", Name: "Exports", EntryTaskID: "entry", TargetTaskID: "target", RecoveryTaskID: "repair", RecoveryOn: conditions, CleanupTaskID: "clean", Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: "branch"}}, UpdatedAt: time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)}}
			d := JoinDraft{ID: "join", Name: "Exports", EntryTaskID: "entry", TargetTaskID: "changed", RecoveryTaskID: "repair", RecoveryOn: conditions, CleanupTaskID: "clean", Branches: b.current.Branches, OriginalUpdatedAt: "2026-11-10T00:00:00Z"}
			if r := NewService(b).SaveJoin(context.Background(), d); r.Outcome != "accepted" {
				t.Fatal(r)
			}
			var packet map[string]json.RawMessage
			json.Unmarshal(b.packet, &packet)
			if len(packet) != 2 || packet["target_task_id"] == nil || packet["expected_updated_at"] == nil {
				t.Fatalf("target-only packet=%s", b.packet)
			}
			d.TargetTaskID = "target"
			b.calls = 0
			if r := NewService(b).SaveJoin(context.Background(), d); r.Outcome != "accepted" || b.calls != 0 {
				t.Fatalf("no-op result=%+v calls=%d", r, b.calls)
			}
			d.RecoveryTaskID = ""
			d.RecoveryOn = []domain.JoinRecoveryCondition{}
			d.CleanupTaskID = ""
			if r := NewService(b).SaveJoin(context.Background(), d); r.Outcome != "accepted" {
				t.Fatal(r)
			}
			json.Unmarshal(b.packet, &packet)
			if string(packet["recovery_task_id"]) != `""` || string(packet["cleanup_task_id"]) != `""` {
				t.Fatalf("clear=%s", b.packet)
			}
		})
	}
}
