package main

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/desktop/automation"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type recoveryCapabilityBackend struct {
	capabilityAutomationBackend
	creates int
}

func (b *recoveryCapabilityBackend) CreateJoin(context.Context, server.JoinCreateRequest) (domain.ParallelJoin, error) {
	b.creates++
	return domain.ParallelJoin{ID: "created"}, nil
}
func TestJoinRecoveryFacadeGatesOnlyNewOptions(t *testing.T) {
	for _, supported := range []bool{false, true} {
		t.Run(map[bool]string{false: "base-only", true: "recovery"}[supported], func(t *testing.T) {
			b := &recoveryCapabilityBackend{}
			ready := make(chan struct{}, 1)
			capabilities := []string{"parallel-joins"}
			if supported {
				capabilities = append(capabilities, "parallel-join-recovery")
			}
			app := newApp(cancellationAppBackend{permissions: []string{"manage"}, capabilities: capabilities}, cancellationAppEmitter{ready}, nil, appServices{automation: automation.NewService(b)})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			app.startup(ctx)
			defer app.shutdown(context.Background())
			select {
			case <-ready:
			case <-time.After(time.Second):
				t.Fatal("no current scheduler")
			}
			draft := automation.JoinDraft{ExpectedDaemonID: "owned-daemon", IsNew: true, Name: "X", EntryTaskID: "entry", TargetTaskID: "target", Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: "branch"}}}
			if r := app.SaveJoin(draft); r.Outcome != "accepted" || b.creates != 1 {
				t.Fatalf("base-only mutation=%+v creates=%d", r, b.creates)
			}
			draft.RecoveryTaskID = "repair"
			draft.RecoveryOn = []domain.JoinRecoveryCondition{}
			draft.CleanupTaskID = "clean"
			r := app.SaveJoin(draft)
			if supported {
				if r.Outcome != "accepted" || b.creates != 2 {
					t.Fatalf("supported mutation=%+v creates=%d", r, b.creates)
				}
			} else if r.Outcome != "rejected" || b.creates != 1 {
				t.Fatalf("unsupported mutation=%+v creates=%d", r, b.creates)
			}
		})
	}
}
